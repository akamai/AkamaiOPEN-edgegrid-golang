package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LoadBalancerActivation is a cloudlets LoadBalancer Activation API interface
	LoadBalancerActivation interface {
		// GetLoadBalancerActivations fetches activations with the most recent listed first
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getloadbalancingconfigactivations
		GetLoadBalancerActivations(context.Context, string) (ActivationsList, error)

		// ActivateLoadBalancerVersion activates the load balacing version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postloadbalancingconfigactivations
		ActivateLoadBalancerVersion(context.Context, ActivateLoadBalancerVersionRequest) (*ActivationResponse, error)
	}

	// ActivationRequestParams contains Request body parameters for ActivateLoadBalancerVersionRequest
	ActivationRequestParams struct {
		Network ActivationNetwork `json:"network"`
		DryRun  bool              `json:"dryrun,omitempty"`
		Version int64             `json:"version"`
	}

	// ActivateLoadBalancerVersionRequest contains request parameters for LoadBalancer version activation
	ActivateLoadBalancerVersionRequest struct {
		OriginID          string
		Async             bool
		ActivationRequest ActivationRequestParams
	}

	// ActivationResponse contains response data for a single LB Version Activation
	ActivationResponse struct {
		ActivatedBy   string            `json:"activatedBy,omitempty"`
		ActivatedDate string            `json:"activatedDate,omitempty"`
		Network       ActivationNetwork `json:"network"`
		OriginID      string            `json:"originId,omitempty"`
		Status        ActivationStatus  `json:"status,omitempty"`
		DryRun        bool              `json:"dryrun,omitempty"`
		Version       int64             `json:"version"`
	}

	// ActivationsList is the response for GetLoadBalancerActivations
	ActivationsList []ActivationResponse

	// ActivationStatus is an LoadBalancer activation status value
	ActivationStatus string

	//ActivationNetwork is the Network value
	ActivationNetwork string
)

const (
	// ActivationStatusActive is an activation that is currently active
	ActivationStatusActive ActivationStatus = "active"

	// ActivationStatusInactive is an activation that is deactivated
	ActivationStatusDeactivated ActivationStatus = "deactivated"

	// ActivationStatusInactive is an activation that is not active
	ActivationStatusInactive ActivationStatus = "inactive"

	// ActivationStatusPending is status of a pending activation
	ActivationStatusPending ActivationStatus = "pending"

	// ActivationStatusPending is status of a failed activation
	ActivationStatusFailed ActivationStatus = "failed"

	// ActivationNetworkStaging is the staging network
	ActivationNetworkStaging ActivationNetwork = "STAGING"

	// ActivationNetworkProduction is the production network
	ActivationNetworkProduction ActivationNetwork = "PRODUCTION"

	// ActivationNetworkProd is the production network
	ActivationNetworkProd ActivationNetwork = "prod"

	// ActivationNetworkProductionLowCase is the production network
	ActivationNetworkProductionLowCase ActivationNetwork = "production"

	// ActivationNetworkStagingLowCase is the staging network
	ActivationNetworkStagingLowCase ActivationNetwork = "staging"
)

var (
	// ErrGetLoadBalancerActivations is returned when GetLoadBalancerActivations fails
	ErrGetLoadBalancerActivations = errors.New("get load balancing activations")
	// ErrActivateLoadBalancerVersion is returned when ActivateLoadBalancerVersion fails
	ErrActivateLoadBalancerVersion = errors.New("activate load balancing version")
)

// Validate validates ActivateLoadBalancerVersionRequest
func (v ActivateLoadBalancerVersionRequest) Validate() error {
	return validation.Errors{
		"OriginID": validation.Validate(v.OriginID, validation.Required),
	}.Filter()
}

//Validate validates ActivationRequestParams Struct
func (v ActivationRequestParams) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(v.Network, validation.In(ActivationNetworkStaging, ActivationNetworkProduction,
			ActivationNetworkProd, ActivationNetworkProductionLowCase, ActivationNetworkStagingLowCase)),
		"Version": validation.Validate(v.Version, validation.Min(0)),
	}.Filter()
}

// GetLoadBalancerActivations fetches activations with the most recent listed first
func (c *cloudlets) GetLoadBalancerActivations(ctx context.Context, originID string) (ActivationsList, error) {
	logger := c.Log(ctx)
	logger.Debug("GetLoadBalancerActivations")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/activations", originID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetLoadBalancerActivations, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetLoadBalancerActivations, err)
	}

	var result ActivationsList
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetLoadBalancerActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetLoadBalancerActivations, c.Error(resp))
	}

	return result, nil
}

// ActivateLoadBalancerVersion activates the load balacing version
func (c *cloudlets) ActivateLoadBalancerVersion(ctx context.Context, params ActivateLoadBalancerVersionRequest) (*ActivationResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ActivateLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateLoadBalancerVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/activations", params.OriginID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrActivateLoadBalancerVersion, err)
	}

	q := uri.Query()
	q.Add("async", strconv.FormatBool(params.Async))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateLoadBalancerVersion, err)
	}

	var result ActivationResponse

	resp, err := c.Exec(req, &result, params.ActivationRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateLoadBalancerVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrActivateLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}
