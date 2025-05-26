package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Change contains change status information
	Change struct {
		AllowedInput []AllowedInput `json:"allowedInput"`
		StatusInfo   *StatusInfo    `json:"statusInfo"`
	}

	// AllowedInput contains the resource locations (path) of data inputs allowed by this Change
	AllowedInput struct {
		Info              string `json:"info"`
		RequiredToProceed bool   `json:"requiredToProceed"`
		Type              string `json:"type"`
		Update            string `json:"update"`
	}

	// StatusInfo contains the status for this Change at this time
	StatusInfo struct {
		DeploymentSchedule *DeploymentSchedule `json:"deploymentSchedule"`
		Description        string              `json:"description"`
		Error              *StatusInfoError    `json:"error,omitempty"`
		State              string              `json:"state"`
		Status             string              `json:"status"`
	}

	// StatusInfoError contains error information for this Change
	StatusInfoError struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Timestamp   string `json:"timestamp"`
	}

	// Certificate is a digital certificate object
	Certificate struct {
		Certificate string `json:"certificate"`
		TrustChain  string `json:"trustChain,omitempty"`
	}

	// GetChangeStatusRequest contains params required to perform GetChangeStatus
	GetChangeStatusRequest struct {
		EnrollmentID int
		ChangeID     int
	}

	// GetChangeRequest contains params required to fetch a specific change (e.g. DV challenges)
	// It is the same for all GET change requests
	GetChangeRequest struct {
		EnrollmentID int
		ChangeID     int
	}

	// CancelChangeRequest contains params required to send CancelChange request
	CancelChangeRequest struct {
		EnrollmentID int
		ChangeID     int
	}

	// CancelChangeResponse is a response object returned from CancelChange request
	CancelChangeResponse struct {
		Change string `json:"change"`
	}

	// AcknowledgementRequest contains params and body required to send acknowledgement. It is the same for all acknowledgement types (dv, pre-verification-warnings etc.)
	AcknowledgementRequest struct {
		Acknowledgement
		EnrollmentID int
		ChangeID     int
	}

	// Acknowledgement is a request body of acknowledgement request
	Acknowledgement struct {
		Acknowledgement string `json:"acknowledgement"`
	}
)

const (
	// AcknowledgementAcknowledge parameter value
	AcknowledgementAcknowledge = "acknowledge"
	// AcknowledgementDeny parameter value
	AcknowledgementDeny = "deny"
)

// Validate validates GetChangeRequest
func (c GetChangeRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
	}.Filter()
}

// Validate validates GetChangeStatusRequest
func (c GetChangeStatusRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
	}.Filter()
}

// Validate validates CancelChangeRequest
func (c CancelChangeRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
	}.Filter()
}

// Validate validates AcknowledgementRequest
func (a AcknowledgementRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID":    validation.Validate(a.EnrollmentID, validation.Required),
		"ChangeID":        validation.Validate(a.ChangeID, validation.Required),
		"Acknowledgement": validation.Validate(a.Acknowledgement),
	}.Filter()
}

// Validate validates Acknowledgement
func (a Acknowledgement) Validate() error {
	return validation.Errors{
		"Acknowledgement": validation.Validate(a.Acknowledgement, validation.Required, validation.In(AcknowledgementAcknowledge, AcknowledgementDeny)),
	}.Filter()
}

// Validate validates Certificate
func (c Certificate) Validate() error {
	return validation.Errors{
		"Certificate": validation.Validate(c.Certificate, validation.Required),
	}.Filter()
}

var (
	// ErrGetChangeStatus is returned when GetChangeStatus fails
	ErrGetChangeStatus = errors.New("fetching change")
	// ErrCancelChange is returned when CancelChange fails
	ErrCancelChange = errors.New("canceling change")
)

func (c *cps) GetChangeStatus(ctx context.Context, params GetChangeStatusRequest) (*Change, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeStatus, ErrStructValidation, err)
	}

	var rval Change

	logger := c.Log(ctx)
	logger.Debug("GetChangeStatus")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d",
		params.EnrollmentID,
		params.ChangeID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetChangeStatus, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeStatus, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change.v2+json")

	resp, err := c.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeStatus, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeStatus, c.Error(resp))
	}

	return &rval, nil
}

func (c *cps) CancelChange(ctx context.Context, params CancelChangeRequest) (*CancelChangeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCancelChange, ErrStructValidation, err)
	}

	var rval CancelChangeResponse

	logger := c.Log(ctx)
	logger.Debug("CancelChange")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d",
		params.EnrollmentID,
		params.ChangeID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCancelChange, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCancelChange, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")

	resp, err := c.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCancelChange, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCancelChange, c.Error(resp))
	}

	return &rval, nil
}
