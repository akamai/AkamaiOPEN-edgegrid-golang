// Package mtlskeystore provides access to the Akamai mTLS Keystore API.
package mtlskeystore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/errs"
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

const resourceNotFoundType = "resource-not-found"
const badRequestType = "bad-request"

// ErrClientCertificateNotFound is returned when the requested client certificate could not be found on the server.
var ErrClientCertificateNotFound = errors.New("the requested resource could not be found on the serve")

// ErrInvalidClientCertificate is returned when the client certificate is either invalid or cannot not be accepted.
var ErrInvalidClientCertificate = errors.New("certificate is either invalid or cannot not be accepted")

// ErrDuplicateClientCertificate is returned when certificate with same name already exists.
var ErrDuplicateClientCertificate = errors.New("certificate with same name already exists")

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
	if errors.Is(target, ErrClientCertificateNotFound) {
		return e.Status == http.StatusNotFound && e.Type == resourceNotFoundType
	}

	if errors.Is(target, ErrInvalidClientCertificate) || errors.Is(target, ErrDuplicateClientCertificate) {
		return e.Status == http.StatusBadRequest && e.Type == badRequestType
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
