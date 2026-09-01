// Package domainownership provides access to the Domain Ownership Manager API.
package domainownership

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/errs"
)

type (
	// Error is a Domain Ownership error interface.
	Error struct {
		Type      string        `json:"type"`
		Title     string        `json:"title"`
		Detail    string        `json:"detail"`
		Instance  string        `json:"instance"`
		Status    int           `json:"status"`
		Errors    []ErrorDetail `json:"errors,omitempty"`
		ProblemID string        `json:"problemId"`
	}

	// ErrorDetail represents a list of specific problems in the error response.
	ErrorDetail struct {
		Type      string `json:"type"`
		Title     string `json:"title"`
		Detail    string `json:"detail"`
		ProblemID string `json:"problemId"`
		Field     string `json:"field"`
	}
)

// Error parses an error from the response.
func (d *domainownership) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		d.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Domain Ownership Manager API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.Status = r.StatusCode

	return &e
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons.
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
