package reportinggroups

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/errs"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Error represents an error returned by the Reporting Groups API.
	Error struct {
		Code       string           `json:"code"`
		Title      string           `json:"title"`
		IncidentID string           `json:"incidentId"`
		Details    []SecondaryError `json:"details"`
		HTTPStatus int              `json:"-"`
	}

	// SecondaryError is a representing inner details of error.
	SecondaryError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

// Error parses an error from the Reporting Groups API response.
func (m *reportinggroups) Error(r *http.Response) error {
	var e Error
	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Code = r.Status
		e.Title = "Failed to read error body"
		e.Details = []SecondaryError{
			{
				Code:    r.Status,
				Message: err.Error(),
			},
		}
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		m.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Reporting Groups API failed. Check details for more information."
		e.Details = []SecondaryError{
			{Message: errs.UnescapeContent(string(body))},
		}
	}

	e.HTTPStatus = r.StatusCode

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

	ignoreCode := t.Code == ""
	ignoreTitle := t.Title == ""
	ignoreDetails := len(t.Details) == 0
	ignoreHTTPStatus := t.HTTPStatus == 0

	matchCode := t.Code == e.Code
	matchTitle := t.Title == e.Title
	matchDetails := slices.Equal(t.Details, e.Details)
	matchHTTPStatus := t.HTTPStatus == e.HTTPStatus

	return (matchCode || ignoreCode) &&
		(matchTitle || ignoreTitle) &&
		(matchDetails || ignoreDetails) &&
		(matchHTTPStatus || ignoreHTTPStatus)
}
