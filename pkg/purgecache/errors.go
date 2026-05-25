// Package purgecache provides access to the Purge Cache API.
package purgecache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/errs"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Error represents an error returned by the Purge Cache API.
	Error struct {
		HTTPStatus                  int    `json:"httpStatus"`
		Title                       string `json:"title"`
		Detail                      string `json:"detail"`
		DescribedBy                 string `json:"describedBy"`
		SupportID                   string `json:"supportId"`
		RateLimit                   int    `json:"rateLimit,omitempty"`
		RateLimitCurrentRequestSize int    `json:"rateLimitCurrentRequestSize,omitempty"`
		RateLimitRemaining          int    `json:"rateLimitRemaining,omitempty"`
	}
)

// Error parses an error from the Purge Cache API response.
func (r *purgecache) Error(resp *http.Response) error {
	var e Error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Log(resp.Request.Context()).Errorf("reading error response body: %s", err)
		e.HTTPStatus = resp.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		r.Log(resp.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.HTTPStatus = resp.StatusCode
		e.Title = "Failed to unmarshal error body. Purge Cache API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

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

	ignoreHTTPStatus := t.HTTPStatus == 0
	ignoreTitle := t.Title == ""
	ignoreDetail := t.Detail == ""
	ignoreDescribedBy := t.DescribedBy == ""

	matchHTTPStatus := t.HTTPStatus == e.HTTPStatus
	matchTitle := t.Title == e.Title
	matchDetail := t.Detail == e.Detail
	matchDescribedBy := t.DescribedBy == e.DescribedBy

	return (matchHTTPStatus || ignoreHTTPStatus) &&
		(matchTitle || ignoreTitle) &&
		(matchDetail || ignoreDetail) &&
		(matchDescribedBy || ignoreDescribedBy)
}
