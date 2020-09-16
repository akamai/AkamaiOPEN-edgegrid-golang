package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
)

type (
	// Activations contains operations available on Activation resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#propertyactivationsgroup
	Activations interface {
		// CreateActivation creates a new activation or deactivation request
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postpropertyactivations
		CreateActivation(context.Context, CreateActivationRequest) (*CreateActivationResponse, error)

		// GetActivation gets details about an activation
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyactivation
		GetActivation(context.Context, GetActivationRequest) (*GetActivationResponse, error)

		// CancelActivation allows for canceling an activation while it is still PENDING
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#deletepropertyactivation
		CancelActivation(context.Context, CancelActivationRequest) (*CancelActivationResponse, error)
	}

	// ActivationFallbackInfo encapsulates information about fast fallback, which may allow you to fallback to a previous activation when
	// POSTing an activation with useFastFallback enabled.
	ActivationFallbackInfo struct {
		FastFallbackAttempted      bool    `json:"fastFallbackAttempted"`
		FallbackVersion            int     `json:"fallbackVersion"`
		CanFastFallback            bool    `json:"canFastFallback"`
		SteadyStateTime            int     `json:"steadyStateTime"`
		FastFallbackExpirationTime int     `json:"fastFallbackExpirationTime"`
		FastFallbackRecoveryState  *string `json:"fastFallbackRecoveryState,omitempty"`
	}

	// Activation represents a property activation resource
	Activation struct {
		AccountID              string                  `json:"accountId,omitempty"`
		ActivationID           string                  `json:"activationId,omitempty"`
		ActivationType         ActivationType          `json:"activationType,omitempty"`
		UseFastFallback        bool                    `json:"useFastFallback"`
		FallbackInfo           *ActivationFallbackInfo `json:"fallbackInfo,omitempty"`
		AcknowledgeWarnings    []string                `json:"acknowledgeWarnings,omitempty"`
		AcknowledgeAllWarnings bool                    `json:"acknowledgeAllWarnings"`
		FastPush               bool                    `json:"fastPush,omitempty"`
		FMAActivationState     string                  `json:"fmaActivationState,omitempty"`
		GroupID                string                  `json:"groupId,omitempty"`
		IgnoreHTTPErrors       bool                    `json:"ignoreHttpErrors,omitempty"`
		PropertyName           string                  `json:"propertyName,omitempty"`
		PropertyID             string                  `json:"propertyId,omitempty"`
		PropertyVersion        int                     `json:"propertyVersion"`
		Network                ActivationNetwork       `json:"network"`
		Status                 ActivationStatus        `json:"status,omitempty"`
		SubmitDate             string                  `json:"submitDate,omitempty"`
		UpdateDate             string                  `json:"updateDate,omitempty"`
		Note                   string                  `json:"note,omitempty"`
		NotifyEmails           []string                `json:"notifyEmails"`
	}

	// CreateActivationRequest is the request parameters for a new activation or deactivation request
	CreateActivationRequest struct {
		PropertyID string
		ContractID string
		GroupID    string
		Activation Activation
	}

	// CreateActivationResponse is the response for a new activation or deactivation
	CreateActivationResponse struct {
		ActivationID   string
		ActivationLink string `json:"activationLink"`
	}

	// GetActivationRequest is the get activation request
	GetActivationRequest struct {
		PropertyID   string
		ContractID   string
		GroupID      string
		ActivationID string
	}

	// ActivationsItems are the activation items array from a response
	ActivationsItems struct {
		Items []*Activation `json:"items"`
	}

	// GetActivationResponse is the get activation response
	GetActivationResponse struct {
		AccountID  string `json:"accountId"`
		ContractID string `json:"contractId"`
		GroupID    string `json:"groupId"`

		Activations ActivationsItems `json:"activations"`

		// RetryAfter is the value of the Retry-After header.
		//  For activations whose status is PENDING, a Retry-After header provides an estimate for when it’s likely to change.
		RetryAfter int `json:"-"`
	}

	// CancelActivationRequest is used to delete a PENDING activation
	CancelActivationRequest struct {
		PropertyID   string
		ActivationID string
		ContractID   string
		GroupID      string
	}

	// CancelActivationResponse is a response from deleting a PENDING activation
	CancelActivationResponse struct {
		Activations ActivationsItems `json:"activations"`
	}

	// ActivationType is an activation type value
	ActivationType string

	// ActivationStatus is an activation status value
	ActivationStatus string

	// ActivationNetwork is the activation network value
	ActivationNetwork string
)

const (
	// ActivationTypeActivate is used for creating a new activation
	ActivationTypeActivate ActivationType = "ACTIVATE"

	// ActivationTypeDeactivate is used for creating a new de-activation
	ActivationTypeDeactivate ActivationType = "DEACTIVATE"

	// ActivationStatusActive is an activation that is currently serving traffic
	ActivationStatusActive ActivationStatus = "ACTIVE"

	// ActivationStatusInactive is an activation that has been superceded by another
	ActivationStatusInactive ActivationStatus = "INACTIVE"

	// ActivationStatusNew is a not yet active activation
	ActivationStatusNew ActivationStatus = "NEW"

	// ActivationStatusPending is the pending status
	ActivationStatusPending ActivationStatus = "PENDING"

	// ActivationStatusAborted is returned when a PENDING activation is successfully canceled
	ActivationStatusAborted ActivationStatus = "ABORTED"

	// ActivationStatusZone1 is not yet active
	ActivationStatusZone1 ActivationStatus = "ZONE_1"

	// ActivationStatusZone2 is not yet active
	ActivationStatusZone2 ActivationStatus = "ZONE_2"

	// ActivationStatusZone3 is not yet active
	ActivationStatusZone3 ActivationStatus = "ZONE_3"

	// ActivationStatusDeactivating is pending deactivation
	ActivationStatusDeactivating ActivationStatus = "PENDING_DEACTIVATION"

	// ActivationStatusDeactivated is deactivated
	ActivationStatusDeactivated ActivationStatus = "DEACTIVATED"

	// ActivationNetworkStaging is the staging network
	ActivationNetworkStaging ActivationNetwork = "STAGING"

	// ActivationNetworkProduction is the production network
	ActivationNetworkProduction ActivationNetwork = "PRODUCTION"
)

// Validate validates CreateActivationRequest
func (v CreateActivationRequest) Validate() error {
	return validation.Errors{
		"PropertyID":                    validation.Validate(v.PropertyID, validation.Required),
		"ContractID":                    validation.Validate(v.ContractID, validation.Required),
		"GroupID":                       validation.Validate(v.GroupID, validation.Required),
		"Activation.AccountID":          validation.Validate(v.Activation.AccountID, validation.Empty),
		"Activation.ActivationID":       validation.Validate(v.Activation.ActivationID, validation.Empty),
		"Activation.FallbackInfo":       validation.Validate(v.Activation.FallbackInfo, validation.Nil),
		"Activation.FMAActivationState": validation.Validate(v.Activation.FMAActivationState, validation.Empty),
		"Activation.GroupID":            validation.Validate(v.Activation.GroupID, validation.Empty),
		"Activation.Network":            validation.Validate(v.Activation.Network, validation.Required),
		"Activation.NotifyEmails":       validation.Validate(v.Activation.NotifyEmails, validation.Length(1, 0)),
		"Activation.PropertyID":         validation.Validate(v.Activation.PropertyID, validation.Empty),
		"Activation.PropertyName":       validation.Validate(v.Activation.PropertyName, validation.Empty),
		"Activation.Status":             validation.Validate(v.Activation.Status, validation.Empty),
		"Activation.SubmitDate":         validation.Validate(v.Activation.SubmitDate, validation.Empty),
		"Activation.UpdateDate":         validation.Validate(v.Activation.UpdateDate, validation.Empty),
	}.Filter()
}

// Validate validates GetActivationRequest
func (v GetActivationRequest) Validate() error {
	return validation.Errors{
		"PropertyID":   validation.Validate(v.PropertyID, validation.Required),
		"ContractID":   validation.Validate(v.ContractID, validation.Required),
		"GroupID":      validation.Validate(v.GroupID, validation.Required),
		"ActivationID": validation.Validate(v.ActivationID, validation.Required),
	}.Filter()
}

// Validate validate CancelActivationRequest
func (v CancelActivationRequest) Validate() error {
	return validation.Errors{
		"PropertyID":   validation.Validate(v.PropertyID, validation.Required),
		"ContractID":   validation.Validate(v.ContractID, validation.Required),
		"GroupID":      validation.Validate(v.GroupID, validation.Required),
		"ActivationID": validation.Validate(v.ActivationID, validation.Required),
	}.Filter()
}

func (p *papi) CreateActivation(ctx context.Context, params CreateActivationRequest) (*CreateActivationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateActivation")

	// explicitly set the activation type
	if params.Activation.ActivationType == "" {
		params.Activation.ActivationType = ActivationTypeActivate
	}

	uri := fmt.Sprintf(
		"/papi/v1/properties/%s/activations?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ContractID,
		params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create createactivation request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var rval CreateActivationResponse

	resp, err := p.Exec(req, &rval, params.Activation)
	if err != nil {
		return nil, fmt.Errorf("createactivation request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, session.NewAPIError(resp, logger)
	}

	id, err := tools.FetchIDFromLocation(rval.ActivationLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrInvalidLocation, err.Error())
	}
	rval.ActivationID = id

	return &rval, nil
}

func (p *papi) GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetActivation")

	uri := fmt.Sprintf(
		"/papi/v1/properties/%s/activations/%s?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ActivationID,
		params.ContractID,
		params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getactivation request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var rval GetActivationResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getactivation request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	// Get the Retry-After header to return the caller
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		rval.RetryAfter = cast.ToInt(retryAfter)
	}

	return &rval, nil
}

func (p *papi) CancelActivation(ctx context.Context, params CancelActivationRequest) (*CancelActivationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetActivation")

	uri := fmt.Sprintf(
		"/papi/v1/properties/%s/activations/%s?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ActivationID,
		params.ContractID,
		params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getactivation request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var rval CancelActivationResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getactivation request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil
}
