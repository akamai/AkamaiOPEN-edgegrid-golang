package netstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/errs"
)

type (
	InnerError struct {
		Type      string `json:"type"`
		Title     string `json:"title,omitempty"`
		Detail    string `json:"detail"`
		ProblemID string `json:"problemId,omitempty"`
	}

	// Error is a netstorage error interface
	Error struct {
		Type      string        `json:"type"`
		Title     string        `json:"title,omitempty"`
		Instance  string        `json:"instance,omitempty"`
		Status    int           `json:"status,omitempty"`
		Detail    string        `json:"detail"`
		ProblemID string        `json:"problemId,omitempty"`
		Errors    []*InnerError `json:"errors,omitempty"`
	}
)

// Error parses an error from the response
func (p *netstorage) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		p.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. NetStorage API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.Status = r.StatusCode

	return &e
}

func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons
func (e *Error) Is(target error) bool {
	if errors.Is(target, ErrNotFound) {
		return e.isErrNotFound()
	}

	var t *Error
	if !errors.As(target, &t) {
		return false
	}

	if e == t {
		return true
	}

	if e.Status != t.Status {
		return false
	}

	return e.Error() == t.Error()
}

func (e *Error) isErrNotFound() bool {
	if e.Status == http.StatusBadRequest {
		for _, innerErr := range e.Errors {
			if innerErr.Detail == "Unable to find the given storage group." {
				return true
			}
		}
	}
	return false
}
