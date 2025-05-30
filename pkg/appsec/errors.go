package appsec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/errs"
)

var (
	// ErrBadRequest is returned when a required parameter is missing.
	ErrBadRequest = errors.New("missing argument")
)

type (
	// Error is an appsec error interface.
	Error struct {
		Type          string `json:"type"`
		Title         string `json:"title"`
		Detail        string `json:"detail"`
		Instance      string `json:"instance,omitempty"`
		BehaviorName  string `json:"behaviorName,omitempty"`
		ErrorLocation string `json:"errorLocation,omitempty"`
		StatusCode    int    `json:"-"`
	}
)

func (p *appsec) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.StatusCode = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		p.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Application Security API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.StatusCode = r.StatusCode

	return &e
}

// Error returns a string formatted using a given title, type, and detail information.
func (e *Error) Error() string {
	return fmt.Sprintf("Title: %s; Type: %s; Detail: %s", e.Title, e.Type, e.Detail)
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

	if e.StatusCode != t.StatusCode {
		return false
	}

	return e.Error() == t.Error()
}
