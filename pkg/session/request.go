package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// Exec will sign and execute the request using the client edgegrid.Config
func (s *session) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	log := s.Log(r.Context())

	if r.UserAgent() == "" {
		r.Header.Set("User-Agent", s.userAgent)
	}

	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json")
	}

	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}

	if r.URL.Host == "" {
		r.URL.Host = s.config.Host
	}

	if err := s.Sign(r); err != nil {
		return nil, err
	}

	if len(in) > 0 {
		data, err := json.Marshal(in[0])
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input: %w", err)
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	if s.trace {
		data, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			log.WithError(err).Error("Failed to dump request")
		} else {
			log.Debug(string(data))
		}
	}

	resp, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}

	if s.trace {
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.WithError(err).Error("Failed to dump response")
		} else {
			log.Debug(string(data))
		}
	}

	if out != nil {
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, out); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// Sign will only sign a request
func (s *session) Sign(r *http.Request) error {
	return s.config.SignRequest(r)
}
