// Package edgeClient is library for making Akamai OPEN Edgegrid API requests
package edgeClient

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	jsonhooks "github.com/RafPe/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	edgegrid "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

const (
	libraryVersion = "0.5.0"
)

var (
	UserAgent = "Akamai-Open-Edgegrid-golang/" + libraryVersion + " golang/" + strings.TrimPrefix(runtime.Version(), "go")
)

// A Client manages communication with the Akamai API.
type Client struct {
	// HTTP client used to communicate with the API.
	clientHTTP *http.Client

	// Configuration used for communication
	edgercConfig edgegrid.Config

	// User agent used when communicating with the GitLab API.
	UserAgent string
}

// NewClient returns a new edgeClient API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(c edgegrid.Config, httpClient *http.Client) *Client {
	return newClient(c, httpClient)
}

func newClient(c edgegrid.Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	apiClient := &Client{clientHTTP: httpClient, UserAgent: UserAgent, edgercConfig: c}

	return apiClient
}

// prepareHTTPRequest creates an HTTP request accordingly to requirements of Akamai API
func (cl *Client) prepareHTTPRequest(method, path string, body io.Reader) (*http.Request, error) {
	var (
		baseURL *url.URL
		err     error
	)

	if strings.HasPrefix(cl.edgercConfig.Host, "https://") {
		baseURL, err = url.Parse(cl.edgercConfig.Host)
	} else {
		baseURL, err = url.Parse("https://" + cl.edgercConfig.Host)
	}

	if err != nil {
		return nil, err
	}

	rel, err := url.Parse(strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}

	u := baseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", UserAgent)

	return req, nil
}

// RequestHTTP creates an HTTP request accordingly to requirements of Akamais API
// and execute it returning response object.
func (cl *Client) RequestHTTP(method, path string, body io.Reader) (*http.Response, error) {
	var (
		err error
	)

	req, err := cl.prepareHTTPRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	res, err := cl.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// RequestJSON creates an HTTP request that can be sent to the Akamai APIs with a JSON body
// The JSON body is encoded and the Content-Type/Accept headers are set automatically.
func (cl *Client) RequestJSON(method, path string, body interface{}) (*http.Response, error) {
	jsonBody, err := jsonhooks.Marshal(body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(jsonBody)
	req, err := cl.prepareHTTPRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json,*/*")

	res, err := cl.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ExecuteRequest performs a given HTTP Request, signed with the Akamai OPEN Edgegrid
// Authorization header. An edgegrid.Response or an error is returned.
func (cl *Client) ExecuteRequest(req *http.Request) (*http.Response, error) {
	cl.clientHTTP.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req = edgegrid.AddRequestHeader(cl.edgercConfig, req)
		return nil
	}

	req = edgegrid.AddRequestHeader(cl.edgercConfig, req)
	res, err := cl.clientHTTP.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// BodyJSON unmarshals the Response.Body into a given data structure
func (cl *Client) BodyJSON(r *http.Response, data interface{}) error {
	if data == nil {
		return errors.New("You must pass in an interface{}")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = jsonhooks.Unmarshal(body, data)

	return err
}
