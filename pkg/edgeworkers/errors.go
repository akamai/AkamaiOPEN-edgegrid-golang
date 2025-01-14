package edgeworkers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/errs"
)

type (
	// Error is an edgeworkers error implementation
	// For details on possible error types, refer to: https://techdocs.akamai.com/edgeworkers/reference/api-errors
	Error struct {
		Type             string     `json:"type,omitempty"`
		Title            string     `json:"title,omitempty"`
		Detail           string     `json:"detail,omitempty"`
		Instance         string     `json:"instance,omitempty"`
		Status           int        `json:"status,omitempty"`
		ErrorCode        string     `json:"errorCode,omitempty"`
		Method           string     `json:"method,omitempty"`
		ServerIP         string     `json:"serverIp,omitempty"`
		ClientIP         string     `json:"clientIp,omitempty"`
		RequestID        string     `json:"requestId,omitempty"`
		RequestTime      string     `json:"requestTime,omitempty"`
		AuthzRealm       string     `json:"authzRealm,omitempty"`
		AdditionalDetail Additional `json:"additionalDetail,omitempty"`
	}

	// Additional holds request_id for edgekv errors
	Additional struct {
		RequestID string `json:"requestId,omitempty"`
	}
)

const (
	errorCodeNotFound                  = "EKV_9000"
	errorCodeVersionIsBeingDeactivated = "EW1031"
	errorCodeVersionAlreadyDeactivated = "EW1032"
)

var (
	// ErrNotFound is returned when edgeKV resource does not exist
	ErrNotFound = errors.New("specified edgeKV resource does not exist")
	// ErrVersionBeingDeactivated is returned when edgeworkers version is currently being deactivated
	ErrVersionBeingDeactivated = errors.New("version is being deactivated")
	// ErrVersionAlreadyDeactivated is returned when edgeworkers version is already deactivated
	ErrVersionAlreadyDeactivated = errors.New("version is already deactivated")
)

// Error parses an error from the response
func (e *edgeworkers) Error(r *http.Response) error {
	var result Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		e.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		result.Status = r.StatusCode
		result.Title = "Failed to read error body"
		result.Detail = err.Error()
		return &result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		e.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		result.Title = fmt.Sprintf("Failed to unmarshal error body. Edgeworkers API failed. Check details for more information.")
		result.Detail = errs.UnescapeContent(string(body))
		result.Status = r.StatusCode
	}
	return &result
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
		return e.Status == http.StatusNotFound && e.ErrorCode == errorCodeNotFound
	}
	if errors.Is(target, ErrVersionBeingDeactivated) {
		return e.ErrorCode == errorCodeVersionIsBeingDeactivated
	}
	if errors.Is(target, ErrVersionAlreadyDeactivated) {
		return e.ErrorCode == errorCodeVersionAlreadyDeactivated
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
