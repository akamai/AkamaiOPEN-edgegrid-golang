package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ChangeOperations is a CPS change API interface
	ChangeOperations interface {
		// GetChangeStatus fetches change status for given enrollment and change ID
		//
		// See: https://techdocs.akamai.com/cps/reference/get-enrollment-change
		GetChangeStatus(context.Context, GetChangeStatusRequest) (*Change, error)

		// CancelChange cancels a pending change
		//
		// See: https://techdocs.akamai.com/cps/reference/delete-enrollment-change
		CancelChange(context.Context, CancelChangeRequest) (*CancelChangeResponse, error)

		// UpdateChange updates a pending change
		// Deprecated: this function will be removed in a future release. Use one of:
		// AcknowledgeChangeManagement(), AcknowledgePostVerificationWarnings(),
		// AcknowledgePreVerificationWarnings(), UploadThirdPartyCertAndTrustChain()
		// or AcknowledgeDVChallenges()
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		UpdateChange(context.Context, UpdateChangeRequest) (*UpdateChangeResponse, error)
	}

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

	// UpdateChangeRequest contains params and body required to send UpdateChange request
	UpdateChangeRequest struct {
		Certificate
		EnrollmentID          int
		ChangeID              int
		AllowedInputTypeParam AllowedInputType
	}

	// UpdateChangeResponse is a response object returned from UpdateChange request
	UpdateChangeResponse struct {
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

	// AllowedInputType represents allowedInputTypeParam used for fetching and updating changes
	AllowedInputType string
)

const (
	// AllowedInputTypeChangeManagementACK parameter value
	AllowedInputTypeChangeManagementACK AllowedInputType = "change-management-ack"
	// AllowedInputTypeLetsEncryptChallengesCompleted parameter value
	AllowedInputTypeLetsEncryptChallengesCompleted AllowedInputType = "lets-encrypt-challenges-completed"
	// AllowedInputTypePostVerificationWarningsACK parameter value
	AllowedInputTypePostVerificationWarningsACK AllowedInputType = "post-verification-warnings-ack"
	// AllowedInputTypePreVerificationWarningsACK parameter value
	AllowedInputTypePreVerificationWarningsACK AllowedInputType = "pre-verification-warnings-ack"
	// AllowedInputTypeThirdPartyCertAndTrustChain parameter value
	AllowedInputTypeThirdPartyCertAndTrustChain AllowedInputType = "third-party-cert-and-trust-chain"
)

const (
	// AcknowledgementAcknowledge parameter value
	AcknowledgementAcknowledge = "acknowledge"
	// AcknowledgementDeny parameter value
	AcknowledgementDeny = "deny"
)

// AllowedInputContentTypeHeader maps content type headers to specific allowed input type params
var AllowedInputContentTypeHeader = map[AllowedInputType]string{
	AllowedInputTypeChangeManagementACK:            "application/vnd.akamai.cps.acknowledgement-with-hash.v1+json",
	AllowedInputTypeLetsEncryptChallengesCompleted: "application/vnd.akamai.cps.acknowledgement.v1+json",
	AllowedInputTypePostVerificationWarningsACK:    "application/vnd.akamai.cps.acknowledgement.v1+json",
	AllowedInputTypePreVerificationWarningsACK:     "application/vnd.akamai.cps.acknowledgement.v1+json",
	AllowedInputTypeThirdPartyCertAndTrustChain:    "application/vnd.akamai.cps.certificate-and-trust-chain.v1+json",
}

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

// Validate validates UpdateChangeRequest
func (c UpdateChangeRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
		"AllowedInputTypeParam": validation.Validate(c.AllowedInputTypeParam, validation.In(
			AllowedInputTypeChangeManagementACK,
			AllowedInputTypeLetsEncryptChallengesCompleted,
			AllowedInputTypePostVerificationWarningsACK,
			AllowedInputTypePreVerificationWarningsACK,
			AllowedInputTypeThirdPartyCertAndTrustChain,
		)),
		"Certificate": validation.Validate(c.Certificate, validation.Required),
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
	// ErrUpdateChange is returned when UpdateChange fails
	ErrUpdateChange = errors.New("updating change")
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCancelChange, c.Error(resp))
	}

	return &rval, nil
}

func (c *cps) UpdateChange(ctx context.Context, params UpdateChangeRequest) (*UpdateChangeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateChange, ErrStructValidation, err)
	}

	var rval UpdateChangeResponse

	logger := c.Log(ctx)
	logger.Debug("UpdateChangeLetsEncryptChallenges")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d/input/update/%s",
		params.EnrollmentID,
		params.ChangeID,
		params.AllowedInputTypeParam),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateChange, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateChange, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", AllowedInputContentTypeHeader[params.AllowedInputTypeParam])

	resp, err := c.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateChange, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateChange, c.Error(resp))
	}

	return &rval, nil
}
