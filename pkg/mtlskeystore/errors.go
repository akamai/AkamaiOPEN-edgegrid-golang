// Package mtlskeystore provides access to the Akamai mTLS Keystore API.
package mtlskeystore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/errs"
)

type (
	// Error is a mtlskeystore error interface.
	Error struct {
		Title     string  `json:"title"`
		Type      string  `json:"type"`
		Detail    string  `json:"detail"`
		Status    int64   `json:"status"`
		ProblemID string  `json:"problemId"`
		Instance  string  `json:"instance"`
		Field     string  `json:"field,omitempty"`
		Parameter string  `json:"parameter,omitempty"`
		Value     string  `json:"value,omitempty"`
		Errors    []Error `json:"errors,omitempty"`
	}
)

// Error parses an error from the mTLS Keystore API response.
func (m *mtlskeystore) Error(r *http.Response) error {
	var e Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = int64(r.StatusCode)
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}
	if err := json.Unmarshal(body, &e); err != nil {
		m.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. mTLS Keystore API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}
	e.Status = int64(r.StatusCode)
	return &e
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s ", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons.
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}
	if e.Status != t.Status {
		return false
	}
	if e.Error() == t.Error() {
		return true
	}

	return false
}
