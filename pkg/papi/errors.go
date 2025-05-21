package papi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/errs"
)

type (
	// Error is a papi error interface
	Error struct {
		Type          string          `json:"type"`
		Title         string          `json:"title,omitempty"`
		Detail        string          `json:"detail"`
		Instance      string          `json:"instance,omitempty"`
		BehaviorName  string          `json:"behaviorName,omitempty"`
		ErrorLocation string          `json:"errorLocation,omitempty"`
		StatusCode    int             `json:"statusCode,omitempty"`
		Errors        json.RawMessage `json:"errors,omitempty"`
		Warnings      json.RawMessage `json:"warnings,omitempty"`
		LimitKey      string          `json:"limitKey,omitempty"`
		Limit         *int            `json:"limit,omitempty"`
		Remaining     *int            `json:"remaining,omitempty"`
	}

	// ActivationError represents errors returned in validation objects in include activation response
	ActivationError struct {
		Type      string                   `json:"type"`
		Title     string                   `json:"title"`
		Instance  string                   `json:"instance"`
		Status    int                      `json:"status"`
		Errors    []ActivationErrorMessage `json:"errors"`
		MessageID string                   `json:"messageId"`
		Result    string                   `json:"result"`
	}

	// ActivationErrorMessage represents detailed information about validation errors
	ActivationErrorMessage struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}
)

// Error parses an error from the response
func (p *papi) Error(r *http.Response) error {
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
		e.Title = "Failed to unmarshal error body. PAPI API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.StatusCode = r.StatusCode

	return &e
}

func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

func (e *ActivationError) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons
func (e *Error) Is(target error) bool {
	if errors.Is(target, ErrSBDNotEnabled) {
		return e.isErrSBDNotEnabled()
	}
	if errors.Is(target, ErrDefaultCertLimitReached) {
		return e.isErrDefaultCertLimitReached()
	}
	if errors.Is(target, ErrNotFound) {
		return e.isErrNotFound()
	}
	if errors.Is(target, ErrActivationTooFar) {
		return e.isActivationTooFar()
	}
	if errors.Is(target, ErrActivationAlreadyActive) {
		return e.isActivationAlreadyActive()
	}

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

// Is handles error comparisons for ActivationError type
func (e *ActivationError) Is(target error) bool {
	if errors.Is(target, ErrMissingComplianceRecord) {
		return e.MessageID == "missing_compliance_record"
	}

	var t *ActivationError
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

func (e *Error) isErrSBDNotEnabled() bool {
	return e.StatusCode == http.StatusForbidden && e.Type == "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable"
}

func (e *Error) isErrDefaultCertLimitReached() bool {
	return e.StatusCode == http.StatusTooManyRequests && e.LimitKey == "DEFAULT_CERTS_PER_CONTRACT" && e.Remaining != nil && *e.Remaining == 0
}

func (e *Error) isErrNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *Error) isActivationTooFar() bool {
	return e.StatusCode == http.StatusBadRequest && e.Title == "Error canceling Activation" && e.Detail == "cancellation_failed.error.activation.toofar"
}

func (e *Error) isActivationAlreadyActive() bool {
	return e.StatusCode == http.StatusUnprocessableEntity && e.Title == "Activation Unprocessable"
}
