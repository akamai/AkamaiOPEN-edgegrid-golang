package cloudaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/errs"
)

type (
	// Error is a cloudaccess error interface
	// For details on possible error types, refer to: https://techdocs.akamai.com/cloud-access-mgr/reference/errors
	Error struct {
		Type         string      `json:"type"`
		Title        string      `json:"title"`
		Detail       string      `json:"detail"`
		Instance     string      `json:"instance"`
		Status       int         `json:"status"`
		AccessKeyUID int         `json:"accessKeyUid"`
		ProblemID    string      `json:"problemId"`
		Errors       []ErrorItem `json:"errors"`
	}

	// ErrorItem is a cloud access error's item
	ErrorItem struct {
		Detail string `json:"detail"`
		Title  string `json:"title"`
		Type   string `json:"type"`
	}
)

// Error parses an error from the response
func (c *cloudaccess) Error(r *http.Response) error {
	var e Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = fmt.Sprintf("Failed to read error body")
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Cloud Access Manager API failed. Check details for more information."
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
