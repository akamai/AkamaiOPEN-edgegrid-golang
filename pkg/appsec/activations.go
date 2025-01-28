package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Activations interface supports the activation and deactivation of security configurations.
	Activations interface {
		// GetActivations returns the status of an activation.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-activation
		GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error)

		// GetActivationHistory lists the activation history for a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-activation-history
		GetActivationHistory(ctx context.Context, params GetActivationHistoryRequest) (*GetActivationHistoryResponse, error)

		// CreateActivations activates a configuration. If acknowledgeWarnings is true and warnings are
		// returned on the first attempt, a second attempt is made acknowledging the warnings.
		//
		CreateActivations(ctx context.Context, params CreateActivationsRequest, acknowledgeWarnings bool) (*CreateActivationsResponse, error)

		// RemoveActivations deactivates a configuration.
		//
		RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error)
	}

	// GetActivationsRequest is used to request the status of an activation request.
	GetActivationsRequest struct {
		ActivationID int `json:"activationId"`
	}

	// GetActivationsResponse is returned from a call to GetActivations.
	GetActivationsResponse struct {
		DispatchCount     int          `json:"dispatchCount"`
		ActivationID      int          `json:"activationId"`
		Action            string       `json:"action"`
		Status            StatusValue  `json:"status"`
		Network           NetworkValue `json:"network"`
		Estimate          string       `json:"estimate"`
		CreatedBy         string       `json:"createdBy"`
		CreateDate        time.Time    `json:"createDate"`
		ActivationConfigs []struct {
			ConfigID              int    `json:"configId"`
			ConfigName            string `json:"configName"`
			ConfigVersion         int    `json:"configVersion"`
			PreviousConfigVersion int    `json:"previousConfigVersion"`
		} `json:"activationConfigs"`
	}

	// GetActivationHistoryRequest is used to request the activation history for a configuration.
	GetActivationHistoryRequest struct {
		ConfigID int `json:"configId"`
	}

	// GetActivationHistoryResponse lists the activation history for a configuration.
	GetActivationHistoryResponse struct {
		ConfigID          int          `json:"configId"`
		ActivationHistory []Activation `json:"activationHistory,omitempty"`
	}

	// Activation represents the status of a configuration activation.
	Activation struct {
		ActivationID       int       `json:"activationId"`
		Version            int       `json:"version"`
		Status             string    `json:"status"`
		Network            string    `json:"Network"`
		ActivatedBy        string    `json:"activatedBy"`
		ActivationDate     time.Time `json:"activationDate"`
		Notes              string    `json:"notes"`
		NotificationEmails []string  `json:"notificationEmails"`
	}

	// CreateActivationsRequest is used to request activation or deactivation of a configuration.
	CreateActivationsRequest struct {
		Action             string   `json:"action"`
		Network            string   `json:"network"`
		Note               string   `json:"note"`
		NotificationEmails []string `json:"notificationEmails"`
		ActivationConfigs  []struct {
			ConfigID      int `json:"configId"`
			ConfigVersion int `json:"configVersion"`
		} `json:"activationConfigs"`
	}

	// CreateActivationsResponse is returned from a call to CreateActivations.
	CreateActivationsResponse struct {
		DispatchCount     int          `json:"dispatchCount"`
		ActivationID      int          `json:"activationId"`
		Action            string       `json:"action"`
		Status            StatusValue  `json:"status"`
		Network           NetworkValue `json:"network"`
		Estimate          string       `json:"estimate"`
		CreatedBy         string       `json:"createdBy"`
		CreateDate        time.Time    `json:"createDate"`
		ActivationConfigs []struct {
			ConfigID              int    `json:"configId"`
			ConfigName            string `json:"configName"`
			ConfigVersion         int    `json:"configVersion"`
			PreviousConfigVersion int    `json:"previousConfigVersion"`
		} `json:"activationConfigs"`
	}

	// ActivationConfigs describes a specific configuration version to be activated or deactivated.
	ActivationConfigs struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
	}

	// RemoveActivationsRequest is used to request deactivation of one or more configurations.
	RemoveActivationsRequest struct {
		ActivationID       int      `json:"-"`
		Action             string   `json:"action"`
		Network            string   `json:"network"`
		Note               string   `json:"note"`
		NotificationEmails []string `json:"notificationEmails"`
		ActivationConfigs  []struct {
			ConfigID      int `json:"configId"`
			ConfigVersion int `json:"configVersion"`
		} `json:"activationConfigs"`
	}

	// RemoveActivationsResponse is returned from a call to RemoveActivations.
	RemoveActivationsResponse struct {
		DispatchCount     int          `json:"dispatchCount"`
		ActivationID      int          `json:"activationId"`
		Action            string       `json:"action"`
		Status            StatusValue  `json:"status"`
		Network           NetworkValue `json:"network"`
		Estimate          string       `json:"estimate"`
		CreatedBy         string       `json:"createdBy"`
		CreateDate        time.Time    `json:"createDate"`
		ActivationConfigs []struct {
			ConfigID              int    `json:"configId"`
			ConfigName            string `json:"configName"`
			ConfigVersion         int    `json:"configVersion"`
			PreviousConfigVersion int    `json:"previousConfigVersion"`
		} `json:"activationConfigs"`
	}
)

// Validate validates a GetActivationsRequest.
func (v GetActivationsRequest) Validate() error {
	return validation.Errors{
		"activationid": validation.Validate(v.ActivationID, validation.Required),
	}.Filter()
}

// Validate validates a GetActivationHistoryRequest.
func (v GetActivationHistoryRequest) Validate() error {
	return validation.Errors{
		"configId": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

func (p *appsec) GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/appsec/v1/activations/%d", params.ActivationID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err)
	}

	q := uri.Query()
	q.Add("updateLatestNetworkStatus", strconv.FormatBool(true))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetActivations request: %w", err)
	}

	var result GetActivationsResponse
	resp, errp := p.Exec(req, &result)
	if errp != nil {
		return nil, fmt.Errorf("get activations request failed: %w", errp)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetActivationHistory(ctx context.Context, params GetActivationHistoryRequest) (*GetActivationHistoryResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetActivationHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/activations",
		params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetActivationHistory request: %w", err)
	}

	var result GetActivationHistoryResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get activation history request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateActivations(ctx context.Context, params CreateActivationsRequest, _ bool) (*CreateActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateActivations")

	uri := "/appsec/v1/activations"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateActivations request: %w", err)
	}

	var result CreateActivationsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("create activations request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	uriget := fmt.Sprintf(
		"/appsec/v1/activations/%d",
		result.ActivationID,
	)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, uriget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetActivation request: %w", err)
	}

	resp, err = p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get activation request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveActivations")

	uri := "/appsec/v1/activations"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveActivations request: %w", err)
	}

	var result RemoveActivationsResponse
	resp, errp := p.Exec(req, &result, params)
	if errp != nil {
		return nil, fmt.Errorf("remove activations request failed: %w", errp)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// ActivationValue is used to create an "enum" of possible Activation.ActivationType values
type ActivationValue string

// NetworkValue is used to create an "enum" of possible Activation.Network values
type NetworkValue string

// StatusValue is used to create an "enum" of possible Activation.Status values
type StatusValue string

const (

	// ActivationTypeActivate is used to activate a configuration.
	ActivationTypeActivate ActivationValue = "ACTIVATE"

	// ActivationTypeDeactivate is used to deactivate a configuration.
	ActivationTypeDeactivate ActivationValue = "DEACTIVATE"

	// NetworkProduction is used to activate/deactivate a configuration in the production network.
	NetworkProduction NetworkValue = "PRODUCTION"

	// NetworkStaging is used to activate/deactivate a configuration in the staging network.
	NetworkStaging NetworkValue = "STAGING"

	// StatusActive indicates that a configuration has been activated.
	StatusActive StatusValue = "ACTIVATED"

	// StatusInactive indicates that a configuration is inactive.
	StatusInactive StatusValue = "INACTIVE"

	// StatusPending indicates that an activation/deactivation request has been received.
	StatusPending StatusValue = "RECEIVED"

	// StatusAborted indicates that an activation/deactivation request has been aborted.
	StatusAborted StatusValue = "ABORTED"

	// StatusFailed indicates that an activation/deactivation request has failed.
	StatusFailed StatusValue = "FAILED"

	// StatusDeactivated indicates that an configuration has been deactivated.
	StatusDeactivated StatusValue = "DEACTIVATED"

	// StatusPendingDeactivation indicates that a deactivation request is in progress.
	StatusPendingDeactivation StatusValue = "PENDING_DEACTIVATION"

	// StatusNew indicates that a deactivation request is new.
	StatusNew StatusValue = "NEW"
)
