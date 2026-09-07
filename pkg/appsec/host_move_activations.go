package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// HostMoveActivations interface supports host move validation and activation with host move.
	HostMoveActivations interface {
		// GetHostMoveValidation validates if hosts need to be moved during activation
		GetHostMoveValidation(ctx context.Context, params GetHostMoveValidationRequest) (*GetHostMoveValidationResponse, error)

		// CreateActivationsWithHostMove activates a configuration with host move support
		CreateActivationsWithHostMove(ctx context.Context, params CreateActivationsWithHostMoveRequest) (*CreateActivationsWithHostMoveResponse, error)
	}

	// GetHostMoveValidationRequest is used to request host move validation for a configuration
	GetHostMoveValidationRequest struct {
		ConfigID      int
		ConfigVersion int
		Network       NetworkValue
	}

	// ConfigInfo represents configuration details in host move validation
	ConfigInfo struct {
		ConfigID      int      `json:"configId"`
		ConfigName    string   `json:"configName"`
		ConfigVersion int      `json:"configVersion"`
		InvalidHosts  []string `json:"invalidHosts"`
	}

	// HostToMove represents a host that needs to be moved between configurations
	HostToMove struct {
		Allowed    bool       `json:"allowed"`
		FromConfig ConfigInfo `json:"fromConfig"`
		Host       string     `json:"host"`
		ToConfig   ConfigInfo `json:"toConfig"`
	}

	// GetHostMoveValidationResponse is returned from a call to GetHostMoveValidation
	GetHostMoveValidationResponse struct {
		HostsToMove []HostToMove `json:"hostsToMove"`
		Network     string       `json:"network"`
	}

	// AcknowledgedInvalidHostsByConfig represents acknowledged invalid hosts by configuration
	AcknowledgedInvalidHostsByConfig struct {
		ConfigID     int      `json:"configId"`
		InvalidHosts []string `json:"invalidHosts"`
	}

	// CreateActivationsWithHostMoveRequest is used to request activation with host move
	CreateActivationsWithHostMoveRequest struct {
		ConfigID                         int                                `json:"-"`
		ConfigVersion                    int                                `json:"-"`
		Action                           string                             `json:"action"`
		Network                          NetworkValue                       `json:"network"`
		Note                             string                             `json:"note"`
		NotificationEmails               []string                           `json:"notificationEmails"`
		AcknowledgedInvalidHosts         []string                           `json:"acknowledgedInvalidHosts"`
		AcknowledgedInvalidHostsByConfig []AcknowledgedInvalidHostsByConfig `json:"acknowledgedInvalidHostsByConfig"`
		HostsToMove                      []HostToMove                       `json:"hostsToMove"`
		SupportID                        string                             `json:"supportId"`
	}

	// ActivationConfig represents configuration details in activation response
	ActivationConfig struct {
		ConfigID      int    `json:"configId"`
		ConfigName    string `json:"configName"`
		ConfigVersion int    `json:"configVersion"`
	}

	// CreateActivationsWithHostMoveResponse is returned from a call to CreateActivationsWithHostMove
	CreateActivationsWithHostMoveResponse struct {
		Action            string             `json:"action"`
		ActivationConfigs []ActivationConfig `json:"activationConfigs"`
		ActivationID      int                `json:"activationId"`
		CreateDate        time.Time          `json:"createDate"`
		CreatedBy         string             `json:"createdBy"`
		Network           string             `json:"network"`
		Status            string             `json:"status"`
	}
)

// Validate validates a GetHostMoveValidationRequest.
func (v GetHostMoveValidationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Network":       validation.Validate(v.Network, validation.Required, validation.In(NetworkProduction, NetworkStaging)),
	})
}

// Validate validates a CreateActivationsWithHostMoveRequest.
func (v CreateActivationsWithHostMoveRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Action":        validation.Validate(v.Action, validation.Required),
		"Network":       validation.Validate(v.Network, validation.Required, validation.In(NetworkProduction, NetworkStaging)),
	})
}

func (p *appsec) GetHostMoveValidation(ctx context.Context, params GetHostMoveValidationRequest) (*GetHostMoveValidationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetHostMoveValidation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/network/%s/host-move-validation",
		params.ConfigID,
		params.ConfigVersion,
		params.Network)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create GetHostMoveValidation request: %s", ErrRequestCreation, err.Error())
	}

	var result GetHostMoveValidationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: get host move validation request failed: %s", ErrAPICallFailure, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateActivationsWithHostMove(ctx context.Context, params CreateActivationsWithHostMoveRequest) (*CreateActivationsWithHostMoveResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateActivationsWithHostMove")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/activations-with-host-move",
		params.ConfigID,
		params.ConfigVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create CreateActivationsWithHostMove request: %s", ErrRequestCreation, err.Error())
	}

	var result CreateActivationsWithHostMoveResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: create activations with host move request failed: %s", ErrAPICallFailure, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
