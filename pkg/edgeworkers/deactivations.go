package edgeworkers

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
	// Deactivation is the response returned by GetDeactivation, DeactivateVersion and ListDeactivation
	Deactivation struct {
		EdgeWorkerID     int               `json:"edgeWorkerId"`
		Version          string            `json:"version"`
		DeactivationID   int               `json:"deactivationId"`
		AccountID        string            `json:"accountId"`
		Status           string            `json:"status"`
		Network          ActivationNetwork `json:"network"`
		Note             string            `json:"note,omitempty"`
		CreatedBy        string            `json:"createdBy"`
		CreatedTime      string            `json:"createdTime"`
		LastModifiedTime string            `json:"lastModifiedTime"`
	}

	// ListDeactivationsRequest describes the parameters for the list deactivations request
	ListDeactivationsRequest struct {
		EdgeWorkerID int
		Version      string
	}

	// DeactivateVersionRequest describes the request parameters for DeactivateVersion
	DeactivateVersionRequest struct {
		EdgeWorkerID int
		DeactivateVersion
	}

	// GetDeactivationRequest describes the request parameters for GetDeactivation
	GetDeactivationRequest struct {
		EdgeWorkerID   int
		DeactivationID int
	}

	// DeactivateVersion represents the request body used to deactivate a version
	DeactivateVersion struct {
		Network ActivationNetwork `json:"network"`
		Note    string            `json:"note,omitempty"`
		Version string            `json:"version"`
	}

	// ListDeactivationsResponse describes the list deactivations response
	ListDeactivationsResponse struct {
		Deactivations []Deactivation `json:"deactivations"`
	}
)

// Validate validates ListDeactivationsRequest
func (r *ListDeactivationsRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates DeactivateVersionRequest
func (r *DeactivateVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":      validation.Validate(r.EdgeWorkerID, validation.Required),
		"DeactivateVersion": validation.Validate(&r.DeactivateVersion),
	}.Filter()
}

// Validate validates DeactivateVersion
func (r *DeactivateVersion) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(
			ActivationNetworkProduction, ActivationNetworkStaging,
		).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			r.Network, ActivationNetworkStaging, ActivationNetworkProduction))),
		"Version": validation.Validate(r.Version, validation.Required),
	}.Filter()
}

// Validate validates GetDeactivationRequest
func (r *GetDeactivationRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":   validation.Validate(r.EdgeWorkerID, validation.Required),
		"DeactivationID": validation.Validate(r.DeactivationID, validation.Required),
	}.Filter()
}

var (
	// ErrListDeactivations is returned when ListDeactivations fails
	ErrListDeactivations = errors.New("list deactivations")
	// ErrDeactivateVersion is returned when DeactivateVersion fails
	ErrDeactivateVersion = errors.New("deactivate version")
	// ErrGetDeactivation is returned when GetDeactivation fails
	ErrGetDeactivation = errors.New("get deactivation")
)

func (e *edgeworkers) ListDeactivations(ctx context.Context, params ListDeactivationsRequest) (*ListDeactivationsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListDeactivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListDeactivations, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrListDeactivations, err.Error())
	}

	q := uri.Query()
	if params.Version != "" {
		q.Add("version", params.Version)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListDeactivations, err.Error())
	}

	var result ListDeactivationsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListDeactivations, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListDeactivations, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeactivateVersion(ctx context.Context, params DeactivateVersionRequest) (*Deactivation, error) {
	logger := e.Log(ctx)
	logger.Debug("DeactivateVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateVersion, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrDeactivateVersion, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateVersion, err.Error())
	}

	var result Deactivation
	resp, err := e.Exec(req, &result, params.DeactivateVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateVersion, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateVersion, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetDeactivation(ctx context.Context, params GetDeactivationRequest) (*Deactivation, error) {
	logger := e.Log(ctx)
	logger.Debug("GetDeactivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDeactivation, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/deactivations/%d", params.EdgeWorkerID, params.DeactivationID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrGetDeactivation, err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDeactivation request: %w", err)
	}
	var result Deactivation
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDeactivation, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDeactivation, e.Error(resp))
	}

	return &result, nil
}
