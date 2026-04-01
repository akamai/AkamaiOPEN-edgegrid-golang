// Package reportinggroups provides access to the Reporting Groups API.
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

	// ErrCreateReportingGroup is returned when there is an error creating a reporting group.
	ErrCreateReportingGroup = errors.New("create reporting group")

	// ErrGetReportingGroup is returned when there is an error getting a reporting group.
	ErrGetReportingGroup = errors.New("get reporting group")

	// ErrListReportingGroups is returned when there is an error listing reporting groups.
	ErrListReportingGroups = errors.New("list reporting groups")

	// ErrUpdateReportingGroup is returned when there is an error updating a reporting group.
	ErrUpdateReportingGroup = errors.New("update reporting group")

	// ErrDeleteReportingGroup is returned when there is an error deleting a reporting group.
	ErrDeleteReportingGroup = errors.New("delete reporting group")

	// ErrListProducts is returned when there is an error listing products within a reporting group.
	ErrListProducts = errors.New("list products")

	// ErrGetReportingGroupsWaterMarkLimits is returned when there is an error getting reporting groups water-mark limits.
	ErrGetReportingGroupsWaterMarkLimits = errors.New("get reporting groups water-mark limits")

	// ErrListCPCodes is returned when there is an error listing CP codes.
	ErrListCPCodes = errors.New("list CP codes")

	// ErrGetCPCodesWaterMarkLimits is returned when there is an error getting CP codes water-mark limits.
	ErrGetCPCodesWaterMarkLimits = errors.New("get CP codes water-mark limits")
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
func (r *reportinggroups) Error(resp *http.Response) error {
	var e Error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Log(resp.Request.Context()).Errorf("reading error response body: %s", err)
		e.Code = resp.Status
		e.Title = "Failed to read error body"
		e.Details = []SecondaryError{
			{
				Code:    resp.Status,
				Message: err.Error(),
			},
		}
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		r.Log(resp.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. Reporting Groups API failed. Check details for more information."
		e.Details = []SecondaryError{
			{Message: errs.UnescapeContent(string(body))},
		}
	}

	e.HTTPStatus = resp.StatusCode

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
