package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Activations is an edgeworkers activations API interface
	Activations interface {
		// ListActivations lists all activations for an EdgeWorker
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-activations-1
		ListActivations(context.Context, ListActivationsRequest) (*ListActivationsResponse, error)

		// GetActivation fetches an EdgeWorker activation by id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/get-activation-1
		GetActivation(context.Context, GetActivationRequest) (*Activation, error)

		// ActivateVersion activates an EdgeWorker version on a given network
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/post-activations-1
		ActivateVersion(context.Context, ActivateVersionRequest) (*Activation, error)

		// CancelPendingActivation cancels pending activation with a given id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/cancel-activation
		CancelPendingActivation(context.Context, CancelActivationRequest) (*Activation, error)
	}

	// ListActivationsRequest contains parameters used to list activations
	ListActivationsRequest struct {
		EdgeWorkerID int
		Version      string
	}

	// ActivateVersionRequest contains path parameters and request body used to activate an edge worker
	ActivateVersionRequest struct {
		EdgeWorkerID int
		ActivateVersion
	}

	// ActivateVersion represents the request body used to activate a version
	ActivateVersion struct {
		Network ActivationNetwork `json:"network"`
		Version string            `json:"version"`
		Note    string            `json:"note,omitempty"`
	}

	// ActivationNetwork represents available activation network types
	ActivationNetwork string

	// GetActivationRequest contains path parameters used to fetch edge worker activation
	GetActivationRequest struct {
		EdgeWorkerID int
		ActivationID int
	}

	// ListActivationsResponse represents a response object returned when listing activations
	ListActivationsResponse struct {
		Activations []Activation `json:"activations"`
	}

	// CancelActivationRequest contains path parameters used to cancel edge worker activation
	CancelActivationRequest struct {
		EdgeWorkerID int
		ActivationID int
	}

	// Activation represents an activation object
	Activation struct {
		AccountID        string `json:"accountId"`
		ActivationID     int    `json:"activationId"`
		CreatedBy        string `json:"createdBy"`
		CreatedTime      string `json:"createdTime"`
		EdgeWorkerID     int    `json:"edgeWorkerId"`
		LastModifiedTime string `json:"lastModifiedTime"`
		Network          string `json:"network"`
		Status           string `json:"status"`
		Version          string `json:"version"`
		Note             string `json:"note"`
	}
)

const (
	// ActivationNetworkStaging is the staging network
	ActivationNetworkStaging ActivationNetwork = "STAGING"

	// ActivationNetworkProduction is the production network
	ActivationNetworkProduction ActivationNetwork = "PRODUCTION"
)

// Validate validates ListActivationsRequest
func (r ListActivationsRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates GetActivationRequest
func (r GetActivationRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
		"ActivationID": validation.Validate(r.ActivationID, validation.Required),
	}.Filter()
}

// Validate validates ActivateVersionRequest
func (r ActivateVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":    validation.Validate(r.EdgeWorkerID, validation.Required),
		"ActivateVersion": validation.Validate(&r.ActivateVersion, validation.Required),
	}.Filter()
}

// Validate validates ActivateVersion
func (r ActivateVersion) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(ActivationNetworkStaging, ActivationNetworkProduction).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, ActivationNetworkStaging, ActivationNetworkProduction))),
		"Version": validation.Validate(r.Version, validation.Required),
	}.Filter()
}

// Validate validates CancelActivationRequest
func (r CancelActivationRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(r.EdgeWorkerID, validation.Required),
		"ActivationID": validation.Validate(r.ActivationID, validation.Required),
	}.Filter()
}

var (
	// ErrListActivations is returned when ListActivations fails
	ErrListActivations = errors.New("list activations")
	// ErrGetActivation is returned when GetActivation fails
	ErrGetActivation = errors.New("get activation")
	// ErrActivateVersion is returned when ActivateVersion fails
	ErrActivateVersion = errors.New("activate version")
	// ErrCancelActivation is returned when CancelPendingActivation fails
	ErrCancelActivation = errors.New("cancel activation")
)

func (e edgeworkers) ListActivations(ctx context.Context, params ListActivationsRequest) (*ListActivationsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/activations", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListActivations, err)
	}

	q := uri.Query()
	if params.Version != "" {
		q.Add("version", params.Version)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListActivations, err)
	}

	var result ListActivationsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListActivations, e.Error(resp))
	}

	return &result, nil
}

func (e edgeworkers) GetActivation(ctx context.Context, params GetActivationRequest) (*Activation, error) {
	logger := e.Log(ctx)
	logger.Debug("GetActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/activations/%d", params.EdgeWorkerID, params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetActivation, err)
	}

	var result Activation
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetActivation, e.Error(resp))
	}

	return &result, nil
}

func (e edgeworkers) ActivateVersion(ctx context.Context, params ActivateVersionRequest) (*Activation, error) {
	logger := e.Log(ctx)
	logger.Debug("ActivateVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/activations", params.EdgeWorkerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateVersion, err)
	}

	var result Activation

	resp, err := e.Exec(req, &result, params.ActivateVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateVersion, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrActivateVersion, e.Error(resp))
	}

	return &result, nil
}

func (e edgeworkers) CancelPendingActivation(ctx context.Context, params CancelActivationRequest) (*Activation, error) {
	logger := e.Log(ctx)
	logger.Debug("CancelPendingActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCancelActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/activations/%d", params.EdgeWorkerID, params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCancelActivation, err)
	}

	var result Activation

	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCancelActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCancelActivation, e.Error(resp))
	}

	return &result, nil
}
