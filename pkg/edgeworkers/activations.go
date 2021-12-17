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
		// See: https://techdocs.akamai.com/edgeworkers/reference/activations#get-activations-1
		ListActivations(context.Context, ListActivationsRequest) (*ListActivationsResponse, error)

		// GetActivation fetches an EdgeWorker activation by id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference-link/get-activation-1
		GetActivation(context.Context, GetActivationRequest) (*Activation, error)

		// CreateActivation activates an EdgeWorker on a given network
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/activations#post-activations-1
		CreateActivation(context.Context, CreateActivationRequest) (*Activation, error)

		// CancelActivation cancels activation with a given id
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/activations#cancel-activation
		CancelActivation(context.Context, CancelActivationRequest) (*Activation, error)
	}

	// ListActivationsRequest contains parameters used to list activations
	ListActivationsRequest struct {
		EdgeWorkerID int
		Version      string
	}

	// CreateActivationRequest contains path parameters and request body used to activate an edge worker
	CreateActivationRequest struct {
		EdgeWorkerID int
		CreateActivation
	}

	// CreateActivation is a request body of create activation API request
	CreateActivation struct {
		Network ActivationNetwork `json:"network"`
		Version string            `json:"version"`
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

// Validate validates CreateActivationRequest
func (r CreateActivationRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":     validation.Validate(r.EdgeWorkerID, validation.Required),
		"CreateActivation": validation.Validate(r.CreateActivation, validation.Required),
	}.Filter()
}

// Validate validates CreateActivation
func (r CreateActivation) Validate() error {
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
	ErrListActivations = errors.New("listing activations")
	// ErrGetActivation is returned when GetActivation fails
	ErrGetActivation = errors.New("getting activation")
	// ErrCreateActivation is returned when CreateActivation fails
	ErrCreateActivation = errors.New("creating activation")
	// ErrCancelActivation is returned when CancelActivation fails
	ErrCancelActivation = errors.New("canceling activation")
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

func (e edgeworkers) CreateActivation(ctx context.Context, params CreateActivationRequest) (*Activation, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/activations", params.EdgeWorkerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateActivation, err)
	}

	var result Activation

	resp, err := e.Exec(req, &result, params.CreateActivation)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateActivation, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateActivation, e.Error(resp))
	}

	return &result, nil
}

func (e edgeworkers) CancelActivation(ctx context.Context, params CancelActivationRequest) (*Activation, error) {
	logger := e.Log(ctx)
	logger.Debug("CancelActivation")

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
