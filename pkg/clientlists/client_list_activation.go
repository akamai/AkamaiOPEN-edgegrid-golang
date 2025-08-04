package clientlists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ActivationParams contains activation general parameters
	ActivationParams struct {
		Action                 ActivationAction  `json:"action"`
		Comments               string            `json:"comments"`
		Network                ActivationNetwork `json:"network"`
		NotificationRecipients []string          `json:"notificationRecipients"`
		SiebelTicketID         string            `json:"siebelTicketId"`
	}

	// GetActivationRequest contains activation request param
	GetActivationRequest struct {
		ActivationID int64
	}

	// GetActivationResponse contains activation details
	GetActivationResponse struct {
		ActivationID      int64            `json:"activationId"`
		CreateDate        string           `json:"createDate"`
		CreatedBy         string           `json:"createdBy"`
		Fast              bool             `json:"fast"`
		InitialActivation bool             `json:"initialActivation"`
		ActivationStatus  ActivationStatus `json:"activationStatus"`
		ListID            string           `json:"listId"`
		Version           int64            `json:"version"`
		ActivationParams
	}

	// CreateActivationRequest contains activation request parameters for CreateActivation method
	CreateActivationRequest struct {
		ListID string
		ActivationParams
	}

	// CreateDeactivationRequest contains deactivation request parameters for CreateDeactivation method
	CreateDeactivationRequest CreateActivationRequest

	// CreateActivationResponse contains activation response
	CreateActivationResponse GetActivationStatusResponse

	// CreateDeactivationResponse contains deactivation response
	CreateDeactivationResponse CreateActivationResponse

	// GetActivationStatusRequest contains request params for GetActivationStatus
	GetActivationStatusRequest struct {
		ListID  string
		Network ActivationNetwork
	}

	// GetActivationStatusResponse contains activation status response
	GetActivationStatusResponse struct {
		Action                 ActivationAction  `json:"action"`
		ActivationID           int64             `json:"activationId"`
		ActivationStatus       ActivationStatus  `json:"activationStatus"`
		Comments               string            `json:"comments"`
		CreateDate             string            `json:"createDate"`
		CreatedBy              string            `json:"createdBy"`
		ListID                 string            `json:"listId"`
		Network                ActivationNetwork `json:"network"`
		NotificationRecipients []string          `json:"notificationRecipients"`
		SiebelTicketID         string            `json:"siebelTicketId"`
		Version                int64             `json:"version"`
	}

	// ActivationNetwork is a type for network field
	ActivationNetwork string

	// ActivationStatus is a type for activationStatus field
	ActivationStatus string

	// ActivationAction is a type for action field
	ActivationAction string
)

const (
	// Staging activation network value STAGING
	Staging ActivationNetwork = "STAGING"
	// Production activation network value PRODUCTION
	Production ActivationNetwork = "PRODUCTION"

	// Inactive activation status value INACTIVE
	Inactive ActivationStatus = "INACTIVE"
	// PendingActivation activation status value PENDING_ACTIVATION
	PendingActivation ActivationStatus = "PENDING_ACTIVATION"
	// Active activation status value ACTIVE
	Active ActivationStatus = "ACTIVE"
	// Deactivated activation status value DEACTIVATED
	Deactivated ActivationStatus = "DEACTIVATED"
	// Modified activation status value MODIFIED
	Modified ActivationStatus = "MODIFIED"
	// PendingDeactivation activation status value PENDING_DEACTIVATION
	PendingDeactivation ActivationStatus = "PENDING_DEACTIVATION"
	// Failed activation status value FAILED
	Failed ActivationStatus = "FAILED"

	// Activate action value ACTIVATE
	Activate ActivationAction = "ACTIVATE"
	// Deactivate action value DEACTIVATE
	Deactivate ActivationAction = "DEACTIVATE"
)

func (v GetActivationRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ActivationID": validation.Validate(v.ActivationID, validation.Required),
	})
}

func (v GetActivationStatusRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID":  validation.Validate(v.ListID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required),
	})
}

func (v CreateActivationRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID":  validation.Validate(v.ListID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required),
	})
}

func (v CreateDeactivationRequest) validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ListID":  validation.Validate(v.ListID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required),
	})
}

// Validate validates ActivationNetwork
func (v ActivationNetwork) Validate() error {
	return validation.In(Staging, Production).Validate(v)
}

func (p *clientlists) CreateActivation(ctx context.Context, params CreateActivationRequest) (*CreateActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("Create Activation")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/lists/%s/activations", params.ListID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create 'create activation' request failed: %s", err.Error())
	}

	var rval CreateActivationResponse

	resp, err := p.Exec(req, &rval, params.ActivationParams)
	if err != nil {
		return nil, fmt.Errorf("create activation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) CreateDeactivation(ctx context.Context, params CreateDeactivationRequest) (*CreateDeactivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("Create Deactivation")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/lists/%s/activations", params.ListID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create 'create deactivation' request failed: %s", err.Error())
	}

	var rval CreateDeactivationResponse

	resp, err := p.Exec(req, &rval, params.ActivationParams)
	if err != nil {
		return nil, fmt.Errorf("create deactivation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) GetActivationStatus(ctx context.Context, params GetActivationStatusRequest) (*GetActivationStatusResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("Get Activation Status")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/lists/%s/environments/%s/status", params.ListID, params.Network)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create get activation status request failed: %s", err.Error())
	}

	var rval GetActivationStatusResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("get activation status request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *clientlists) GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("Get Activation")

	if err := params.validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/client-list/v1/activations/%d", params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create get activation request failed: %s", err.Error())
	}

	var rval GetActivationResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("get activation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
