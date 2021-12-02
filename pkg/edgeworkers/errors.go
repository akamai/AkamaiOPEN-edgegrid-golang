package edgeworkers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Error is an edgeworkers error implementation
	// For details on possible error types, refer to: https://techdocs.akamai.com/edgeworkers/reference/api-errors
	Error struct {
		Type      string `json:"type,omitempty"`
		Title     string `json:"title,omitempty"`
		Detail    string `json:"detail,omitempty"`
		Instance  string `json:"instance,omitempty"`
		Status    int    `json:"status,omitempty"`
		ErrorCode string `json:"errorCode,omitempty"`
	}
)

// Error parses an error from the response
func (c *edgeworkers) Error(r *http.Response) error {
	var e Error
	var body []byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = string(body)
		e.Status = r.StatusCode
	}
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
