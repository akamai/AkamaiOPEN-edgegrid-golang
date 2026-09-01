package apidefinitions

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	validators "github.com/go-ozzo/ozzo-validation/v4/is"
)

type (

	// VerifyVersionRequest contains parameters for VerifyVersion operation
	VerifyVersionRequest struct {
		VersionNumber int64
		APIEndpointID int64
		Body          VerifyVersionRequestBody
	}

	// ActivateVersionRequest contains parameters for ActivateVersion operation
	ActivateVersionRequest struct {
		VersionNumber int64
		APIEndpointID int64
		Body          ActivationRequestBody
	}

	// VerifyVersionRequestBody contains body for VerifyVersion operation
	VerifyVersionRequestBody struct {
		Networks []NetworkType `json:"networks"`
	}

	// ActivationRequestBody contains body for ActivateVersion and DeactivateVersion operation
	ActivationRequestBody struct {
		Networks               []NetworkType `json:"networks"`
		Notes                  string        `json:"notes,omitempty"`
		NotificationRecipients []string      `json:"notificationRecipients,omitempty"`
	}

	// ActivateVersionResponse represents a response for ActivateVersion operation
	ActivateVersionResponse struct {
		Networks               []NetworkType
		Notes                  string
		NotificationRecipients []string
	}

	// DeactivateVersionRequest contains parameters for DeactivateVersion operation
	DeactivateVersionRequest struct {
		VersionNumber int64
		APIEndpointID int64
		Body          ActivationRequestBody
	}

	// DeactivateVersionResponse represents a response for DeactivateVersion operation
	DeactivateVersionResponse struct {
		Networks               []NetworkType
		Notes                  string
		NotificationRecipients []string
	}

	// VerifyVersionResponse represents a response for VerifyVersion operation
	VerifyVersionResponse []VerifyVersionAlert

	// ActivationType is an activation type value
	ActivationType string

	// ActivationStatus is an activation status value
	ActivationStatus string

	// NetworkType is the activation network value
	NetworkType string

	// VerifyVersionAlert represents Activation alerts
	VerifyVersionAlert struct {
		Detail   string   `json:"detail"`
		Severity Severity `json:"severity"`
	}

	// Severity represents the severity level of VerifyVersionAlert
	Severity string
)

const (
	// ActivationStatusPending is the pending status
	ActivationStatusPending ActivationStatus = "PENDING"
	// ActivationStatusActive is the active status
	ActivationStatusActive ActivationStatus = "ACTIVE"
	// ActivationStatusDeactivated is the deactivated status
	ActivationStatusDeactivated ActivationStatus = "DEACTIVATED"
	// ActivationStatusFailed is the failed status
	ActivationStatusFailed ActivationStatus = "FAILED"
	// ActivationNetworkStaging is the staging network
	ActivationNetworkStaging NetworkType = "STAGING"
	// ActivationNetworkProduction is the production network
	ActivationNetworkProduction NetworkType = "PRODUCTION"
	// SeverityError is the error severity
	SeverityError Severity = "ERROR"
	// SeverityWarning is the warning severity
	SeverityWarning Severity = "WARNING"
)

var (
	// ErrVerifyVersion is returned in case of error occurs in VerifyVersion operation
	ErrVerifyVersion = errors.New("verify version")
	// ErrActivateVersion is returned in case of error occurs in ActivateVersion operation
	ErrActivateVersion = errors.New("activate version")
	// ErrDeactivateVersion is returned in case of error occurs in DeactivateVersion operation
	ErrDeactivateVersion = errors.New("deactivate version")
)

// Validate validates NetworkType
func (n NetworkType) Validate() error {
	return validation.In(ActivationNetworkStaging, ActivationNetworkProduction).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' ", n, ActivationNetworkStaging, ActivationNetworkProduction)).
		Validate(n)
}

// Validate validates VerifyVersionRequest
func (r VerifyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
		"Body":          validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates ActivateVersionRequest
func (r ActivateVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
		"Body":          validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates DeactivateVersionRequest
func (r DeactivateVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
		"Body":          validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates VerifyVersionRequestBody
func (b VerifyVersionRequestBody) Validate() error {
	return validation.Errors{
		"Networks": validation.Validate(b.Networks, validation.Required),
	}.Filter()
}

// Validate validates ActivationRequestBody
func (b ActivationRequestBody) Validate() error {
	return validation.Errors{
		"Networks":               validation.Validate(b.Networks, validation.Required),
		"NotificationRecipients": validation.Validate(b.NotificationRecipients, validation.Each(validators.EmailFormat)),
	}.Filter()
}

func (a *apidefinitions) ActivateVersion(ctx context.Context, params ActivateVersionRequest) (*ActivateVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ActivateVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d/activate", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateVersion, err)
	}

	var result ActivateVersionResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrActivateVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) DeactivateVersion(ctx context.Context, params DeactivateVersionRequest) (*DeactivateVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("DeactivateVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d/deactivate", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateVersion, err)
	}

	var result DeactivateVersionResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) VerifyVersion(ctx context.Context, params VerifyVersionRequest) (VerifyVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("VerifyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrVerifyVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d/activate/verify", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrVerifyVersion, err)
	}

	var result VerifyVersionResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrVerifyVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrVerifyVersion, a.Error(resp))
	}

	return result, nil
}
