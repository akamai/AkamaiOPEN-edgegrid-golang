package accountprotection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/errs"
)

type (
	// Error is a accountProtection error interface.
	Error struct {
		// Type is the error type
		Type string `json:"type"`

		// Title is the error title
		Title string `json:"title"`

		// Detail is the error detail
		Detail string `json:"detail"`

		// Errors contains a list of errors in the response
		Errors []Error `json:"errors,omitempty"`

		// StatusCode is the HTTP status code of the response
		StatusCode int `json:"status,omitempty"`
	}
)

// Error extract error from response call
func (ap *accountProtection) Error(r *http.Response) error {
	var e Error
	var body []byte

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ap.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.StatusCode = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		ap.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Bot Manager API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.StatusCode = r.StatusCode

	return &e
}

// Error returns a string formatted using a given title, type, and detail information.
func (e *Error) Error() string {
	detail := e.Detail
	if len(e.Errors) > 0 {
		var childErrorDetails []string
		for _, err := range e.Errors {
			childErrorDetails = append(childErrorDetails, err.Detail)
		}
		detail += ": [" + strings.Join(childErrorDetails, ", ") + " ]"
	}
	return fmt.Sprintf("Title: %s; Type: %s; Detail: %s", e.Title, e.Type, detail)
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
