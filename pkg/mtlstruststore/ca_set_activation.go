package mtlstruststore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ActivateCASetVersionRequest holds request content for ActivateCASetVersion.
	ActivateCASetVersionRequest struct {
		// CASetID that needs to be activated on the network.
		CASetID string `json:"-"`

		// Version number of the CA set that needs to be activated on the network.
		Version int64 `json:"-"`

		// Network on which the CA set version needs to be activated. One of "STAGING" or "PRODUCTION".
		Network ActivationNetwork `json:"network"`
	}

	// DeactivateCASetVersionRequest holds request content for DeactivateCASetVersion.
	DeactivateCASetVersionRequest ActivateCASetVersionRequest

	// ActivateCASetVersionResponse contains response from ActivateCASetVersion.
	ActivateCASetVersionResponse struct {
		// ActivationID is a unique identifier representing the CA set activation.
		ActivationID int64 `json:"activationId"`

		// ActivationLink is the link to the CA set activation.
		ActivationLink string `json:"activationLink"`

		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// CASetName is the name of the CA set.
		CASetName string `json:"caSetName"`

		// CASetLink is the link to the CA set.
		CASetLink string `json:"caSetLink"`

		// CreatedBy is the user who created the CA set.
		CreatedBy string `json:"createdBy"`

		// CreatedDate is the date when the CA set was created.
		CreatedDate time.Time `json:"createdDate"`

		// FailureReason is the reason for failure, if any.
		FailureReason *string `json:"failureReason"`

		// ModifiedBy is the user who modified the CA set last time.
		ModifiedBy *string `json:"modifiedBy"`

		// ModifiedDate is the date when the CA set was modified.
		ModifiedDate *time.Time `json:"modifiedDate"`

		// Network is the network on which the CA set is activated. It can be one of "STAGING" or "PRODUCTION".
		Network string `json:"network"`

		// ActivationStatus is the status of the CA set activation. It can be one of the following: "IN_PROGRESS", "COMPLETE" or "FAILED".
		ActivationStatus string `json:"activationStatus"`

		// ActivationType is the type of activation. It can be one of the following: "ACTIVATE" or "DEACTIVATE".
		ActivationType string `json:"activationType"`

		// PercentComplete is the percentage of completion of the activation.
		PercentComplete int `json:"percentComplete"`

		// Version is the version number of the CA set.
		Version int64 `json:"version"`

		// VersionLink is the link to the CA set version.
		VersionLink string `json:"versionLink"`

		// RetryAfter is a time when CA set version activation/deactivation status can be checked again.
		// Usually 300 seconds after the activation request is made.
		// This header value is returned only if the CA set version activation/deactivation status is "IN_PROGRESS".
		RetryAfter time.Time

		// Validation contains validation information for the activation.
		Validation *Validation `json:"validation"`
	}

	// DeactivateCASetVersionResponse contains response from DeactivateCASetVersion.
	DeactivateCASetVersionResponse ActivateCASetVersionResponse

	// ActivationNetwork represents the network type: 'staging' or 'production'.
	ActivationNetwork string

	// GetCASetVersionActivationRequest holds request content for GetCASetVersionActivation.
	GetCASetVersionActivationRequest struct {
		// CASetID is the ID of the CA set to get activation details for.
		CASetID string

		// Version is the version number of the CA set to get activation details for.
		Version int64

		// ActivationID is the ID of the activation to get details for.
		ActivationID int64
	}

	// GetCASetVersionActivationResponse contains response from GetCASetVersionActivation.
	GetCASetVersionActivationResponse ActivateCASetVersionResponse

	// ListCASetVersionActivationsRequest holds request content for ListCASetActivations.
	ListCASetVersionActivationsRequest struct {
		// CASetID is the ID of the CA set to list activations for.
		CASetID string

		// Version is the optional version number of the CA set to list activations for.
		Version int64
	}

	// ListCASetVersionActivationsResponse contains response from ListCASetVersionActivations.
	ListCASetVersionActivationsResponse struct {
		// Activations is the list of activations for the CA set version.
		Activations []ActivateCASetVersionResponse `json:"activations"`
	}

	// ListCASetActivationsRequest holds request content for ListCASetActivations.
	ListCASetActivationsRequest struct {
		// CASetID is the ID of the CA set to list activations for.
		CASetID string
	}

	// ListCASetActivationsResponse contains response from ListCASetActivations.
	ListCASetActivationsResponse struct {
		// Activations is the list of activations for the CA set.
		Activations []ActivateCASetVersionResponse `json:"activations"`
	}
)

const (
	// ActivationNetworkStaging represents staging network.
	ActivationNetworkStaging ActivationNetwork = "STAGING"
	// ActivationNetworkProduction represents production network.
	ActivationNetworkProduction ActivationNetwork = "PRODUCTION"
)

var (
	// ErrActivateCASetVersion is returned when the request to activate a CA set version fails.
	ErrActivateCASetVersion = errors.New("activate ca set version failed")
	// ErrDeactivateCASetVersion is returned when the request to deactivate a CA set version fails.
	ErrDeactivateCASetVersion = errors.New("deactivate ca set version failed")
	// ErrGetCASetVersionActivation is returned when the request to get CA set version activation fails.
	ErrGetCASetVersionActivation = errors.New("get ca set version activation failed")
	// ErrListCASetVersionActivations is returned when the request to list CA set version activations fails.
	ErrListCASetVersionActivations = errors.New("list ca set version activations failed")
	// ErrListCASetActivations is returned when the request to list CA set activations fails.
	ErrListCASetActivations = errors.New("list ca set activations failed")
)

// Validate validates ActivateCASetVersionRequest.
func (r ActivateCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required),
		"Version": validation.Validate(r.Version, validation.Required),
		"Network": validation.Validate(r.Network, validation.Required, r.Network.Validate()),
	})
}

// Validate validates DeactivateCASetVersionRequest.
func (r DeactivateCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required),
		"Version": validation.Validate(r.Version, validation.Required),
		"Network": validation.Validate(r.Network, validation.Required, r.Network.Validate()),
	})
}

// Validate validates ActivationNetwork.
func (n ActivationNetwork) Validate() validation.InRule {
	return validation.In(ActivationNetworkStaging, ActivationNetworkProduction).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			n, ActivationNetworkStaging, ActivationNetworkProduction))
}

// Validate validates GetCASetVersionActivationRequest.
func (r GetCASetVersionActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID":      validation.Validate(r.CASetID, validation.Required),
		"Version":      validation.Validate(r.Version, validation.Required),
		"ActivationID": validation.Validate(r.ActivationID, validation.Required),
	})
}

// Validate validates ListCASetVersionActivationsRequest.
func (r ListCASetVersionActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required),
		"Version": validation.Validate(r.Version, validation.Required),
	})
}

// Validate validates ListCASetActivationsRequest.
func (r ListCASetActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required),
	})
}

func (m *mtlstruststore) ActivateCASetVersion(ctx context.Context, params ActivateCASetVersionRequest) (*ActivateCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ActivateCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/activate", params.CASetID, params.Version))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateCASetVersion, err)
	}

	var result ActivateCASetVersionResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) DeactivateCASetVersion(ctx context.Context, params DeactivateCASetVersionRequest) (*DeactivateCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("DeactivateCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/deactivate", params.CASetID, params.Version))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateCASetVersion, err)
	}

	var result DeactivateCASetVersionResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) GetCASetVersionActivation(ctx context.Context, params GetCASetVersionActivationRequest) (*GetCASetVersionActivationResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetCASetVersionActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCASetVersionActivation, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/activations/%d", params.CASetID, params.Version, params.ActivationID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCASetVersionActivation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCASetVersionActivation, err)
	}

	var result GetCASetVersionActivationResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCASetVersionActivation, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, m.Error(resp)
	}

	if resp.Header.Get("Retry-After") != "" {
		after, err := time.Parse(time.RFC1123, resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse Retry-After header: %s", ErrGetCASetVersionActivation, err)
		}
		result.RetryAfter = after
	}

	return &result, nil
}

func (m *mtlstruststore) ListCASetVersionActivations(ctx context.Context, params ListCASetVersionActivationsRequest) (*ListCASetVersionActivationsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASetVersionActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASetVersionActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/activations", params.CASetID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASetVersionActivations, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASetVersionActivations, err)
	}

	var result ListCASetVersionActivationsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASetVersionActivations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) ListCASetActivations(ctx context.Context, params ListCASetActivationsRequest) (*ListCASetActivationsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASetActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASetActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/activations", params.CASetID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASetActivations, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASetActivations, err)
	}

	var result ListCASetActivationsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASetActivations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}
