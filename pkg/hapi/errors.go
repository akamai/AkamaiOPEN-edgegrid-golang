package hapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/errs"
)

type (
	// Error is a hapi error interface
	Error struct {
		Type            string      `json:"type"`
		Title           string      `json:"title"`
		Detail          string      `json:"detail"`
		Instance        string      `json:"instance,omitempty"`
		RequestInstance string      `json:"requestInstance,omitempty"`
		Method          string      `json:"method,omitempty"`
		RequestTime     string      `json:"requestTime,omitempty"`
		BehaviorName    string      `json:"behaviorName,omitempty"`
		ErrorLocation   string      `json:"errorLocation,omitempty"`
		Status          int         `json:"status,omitempty"`
		DomainPrefix    string      `json:"domainPrefix,omitempty"`
		DomainSuffix    string      `json:"domainSuffix,omitempty"`
		Errors          []ErrorItem `json:"errors,omitempty"`
	}

	// ErrorItem represents single error item
	ErrorItem struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}
)

// Error parses an error from the response
func (h *hapi) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = fmt.Sprintf("Failed to read error body")
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		h.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = fmt.Sprintf("Failed to unmarshal error body. HAPI API failed. Check details for more information.")
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
	return e.Status == http.StatusNotFound
}
