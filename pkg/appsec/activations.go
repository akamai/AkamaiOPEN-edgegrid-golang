package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Activations represents a collection of Activations
//
// See: Activations.GetActivations()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// Activations  contains operations available on Activations  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getactivations
	Activations interface {
		//GetActivationss(ctx context.Context, params GetActivationssRequest) (*GetActivationssResponse, error)
		GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error)
		CreateActivations(ctx context.Context, params CreateActivationsRequest, acknowledgeWarnings bool) (*CreateActivationsResponse, error)
		RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error)
	}

	GetActivationsRequest struct {
		ActivationID int `json:"activationId"`
	}

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

	ActivationsPost struct {
		Action             string   `json:"action"`
		Network            string   `json:"network"`
		Note               string   `json:"note"`
		NotificationEmails []string `json:"notificationEmails"`
		ActivationConfigs  []struct {
			ConfigID      int `json:"configId"`
			ConfigVersion int `json:"configVersion"`
		} `json:"activationConfigs"`
	}

	ActivationConfigs struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
	}

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

// Validate validates GetActivationsRequest
func (v GetActivationsRequest) Validate() error {
	return validation.Errors{
		"activationid": validation.Validate(v.ActivationID, validation.Required),
	}.Filter()
}

// GetActivations populates  *Activations with it's related Activations
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func (p *appsec) GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetActivations")

	var rval GetActivationsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/activations/%d",
		params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getactivations request: %w", err)
	}

	resp, errp := p.Exec(req, &rval)
	if errp != nil {
		return nil, fmt.Errorf("getactivations request failed: %w", errp)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Save activates a given Configuration
//
// If acknowledgeWarnings is true and warnings are returned on the first attempt,
// a second attempt is made, acknowledging the warnings.
//
func (p *appsec) CreateActivations(ctx context.Context, params CreateActivationsRequest, acknowledgeWarnings bool) (*CreateActivationsResponse, error) {
	//func (activations *CreateActivationsResponse) SaveActivations(postpayload *ActivationsPost, acknowledgeWarnings bool, correlationid string) (*CreateActivationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateActivations")

	uri := "/appsec/v1/activations"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create activation request: %w", err)
	}

	var rval CreateActivationsResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create activationrequest failed: %w", err)
	}

	var rvalget CreateActivationsResponse

	uriget := fmt.Sprintf(
		"/appsec/v1/activations/%d",
		rval.ActivationID,
	)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, uriget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation request: %w", err)
	}

	resp, err = p.Exec(req, &rvalget)
	if err != nil {
		return nil, fmt.Errorf("get activation request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rvalget, nil

}

// Delete will delete a Activations
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteactivations
func (p *appsec) RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateRatePolicy")

	uri := "/appsec/v1/activations"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remove activation request: %w", err)
	}

	var rval RemoveActivationsResponse

	_, errp := p.Exec(req, &rval, params)
	if errp != nil {
		return nil, fmt.Errorf("remove activationrequest failed: %w", errp)
	}

	return &rval, nil
}

// ActivationValue is used to create an "enum" of possible Activation.ActivationType values
type ActivationValue string

// NetworkValue is used to create an "enum" of possible Activation.Network values
type NetworkValue string

// StatusValue is used to create an "enum" of possible Activation.Status values
type StatusValue string

const (
	// ActivationTypeActivate Activation.ActivationType value ACTIVATE
	ActivationTypeActivate ActivationValue = "ACTIVATE"
	// ActivationTypeDeactivate Activation.ActivationType value DEACTIVATE
	ActivationTypeDeactivate ActivationValue = "DEACTIVATE"

	// NetworkProduction Activation.Network value PRODUCTION
	NetworkProduction NetworkValue = "PRODUCTION"
	// NetworkStaging Activation.Network value STAGING
	NetworkStaging NetworkValue = "STAGING"

	// StatusActive Activation.Status value ACTIVE
	StatusActive StatusValue = "ACTIVATED"
	// StatusInactive Activation.Status value INACTIVE
	StatusInactive StatusValue = "INACTIVE"
	// StatusPending Activation.Status value RECEIVED
	StatusPending StatusValue = "RECEIVED"
	// StatusAborted Activation.Status value ABORTED
	StatusAborted StatusValue = "ABORTED"
	// StatusFailed Activation.Status value FAILED
	StatusFailed StatusValue = "FAILED"
	// StatusDeactivated Activation.Status value DEACTIVATED
	StatusDeactivated StatusValue = "DEACTIVATED"
	// StatusPendingDeactivation Activation.Status value PENDING_DEACTIVATION
	StatusPendingDeactivation StatusValue = "PENDING_DEACTIVATION"
	// StatusNew Activation.Status value NEW
	StatusNew StatusValue = "NEW"
)
