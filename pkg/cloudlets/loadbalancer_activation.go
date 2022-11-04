package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LoadBalancerActivations is a cloudlets LoadBalancer Activation API interface
	LoadBalancerActivations interface {
		// ListLoadBalancerActivations fetches activations with the most recent listed first
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getloadbalancingconfigactivations
		ListLoadBalancerActivations(context.Context, ListLoadBalancerActivationsRequest) ([]LoadBalancerActivation, error)

		// ActivateLoadBalancerVersion activates the load balancing version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postloadbalancingconfigactivations
		ActivateLoadBalancerVersion(context.Context, ActivateLoadBalancerVersionRequest) (*LoadBalancerActivation, error)
	}

	// ListLoadBalancerActivationsRequest contains request parameters for ListLoadBalancerActivations
	ListLoadBalancerActivationsRequest struct {
		OriginID   string
		Network    LoadBalancerActivationNetwork
		LatestOnly bool
		PageSize   *int64
		Page       *int64
	}

	// ActivateLoadBalancerVersionRequest contains request parameters for LoadBalancer version activation
	ActivateLoadBalancerVersionRequest struct {
		OriginID string
		Async    bool
		LoadBalancerVersionActivation
	}

	// LoadBalancerVersionActivation contains request parameters for ActivateLoadBalancerVersionRequest
	LoadBalancerVersionActivation struct {
		Network LoadBalancerActivationNetwork `json:"network"`
		DryRun  bool                          `json:"dryrun,omitempty"`
		Version int64                         `json:"version"`
	}

	// LoadBalancerActivation contains response data for a single LB Version Activation
	LoadBalancerActivation struct {
		ActivatedBy   string                        `json:"activatedBy,omitempty"`
		ActivatedDate string                        `json:"activatedDate,omitempty"`
		Network       LoadBalancerActivationNetwork `json:"network"`
		OriginID      string                        `json:"originId,omitempty"`
		Status        LoadBalancerActivationStatus  `json:"status,omitempty"`
		DryRun        bool                          `json:"dryrun,omitempty"`
		Version       int64                         `json:"version"`
	}

	//LoadBalancerActivationNetwork is the activation network type for load balancer
	LoadBalancerActivationNetwork string

	// LoadBalancerActivationStatus is an activation status type for load balancer
	LoadBalancerActivationStatus string
)

const (
	// LoadBalancerActivationStatusActive is an activation that is currently active
	LoadBalancerActivationStatusActive LoadBalancerActivationStatus = "active"
	// LoadBalancerActivationStatusDeactivated is an activation that is deactivated
	LoadBalancerActivationStatusDeactivated LoadBalancerActivationStatus = "deactivated"
	// LoadBalancerActivationStatusInactive is an activation that is not active
	LoadBalancerActivationStatusInactive LoadBalancerActivationStatus = "inactive"
	// LoadBalancerActivationStatusPending is status of a pending activation
	LoadBalancerActivationStatusPending LoadBalancerActivationStatus = "pending"
	// LoadBalancerActivationStatusFailed is status of a failed activation
	LoadBalancerActivationStatusFailed LoadBalancerActivationStatus = "failed"

	// LoadBalancerActivationNetworkStaging is the staging network value for load balancer
	LoadBalancerActivationNetworkStaging LoadBalancerActivationNetwork = "STAGING"
	// LoadBalancerActivationNetworkProduction is the production network value for load balancer
	LoadBalancerActivationNetworkProduction LoadBalancerActivationNetwork = "PRODUCTION"

	// NetworkParamStaging is the staging network param value for ListLoadBalancerActivationsRequest
	NetworkParamStaging LoadBalancerActivationNetwork = "staging"
	// NetworkParamProduction is the production network param value for ListLoadBalancerActivationsRequest
	NetworkParamProduction LoadBalancerActivationNetwork = "prod"
)

var (
	// ErrListLoadBalancerActivations is returned when ListLoadBalancerActivations fails
	ErrListLoadBalancerActivations = errors.New("list load balancer activations")
	// ErrActivateLoadBalancerVersion is returned when ActivateLoadBalancerVersion fails
	ErrActivateLoadBalancerVersion = errors.New("activate load balancer version")
)

// Validate validates ActivateLoadBalancerVersionRequest
func (v ActivateLoadBalancerVersionRequest) Validate() error {
	errs := validation.Errors{
		"OriginID": validation.Validate(v.OriginID, validation.Required),
		"Params":   validation.Validate(v.LoadBalancerVersionActivation),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates ListLoadBalancerActivationsRequest
func (v ListLoadBalancerActivationsRequest) Validate() error {
	errs := validation.Errors{
		"OriginID": validation.Validate(v.OriginID, validation.Required),
		"Network": validation.Validate(v.Network, validation.In(NetworkParamStaging, NetworkParamProduction).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '' (empty)", v.Network, NetworkParamStaging, NetworkParamProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

//Validate validates LoadBalancerVersionActivation Struct
func (v LoadBalancerVersionActivation) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(v.Network, validation.Required, validation.In(LoadBalancerActivationNetworkStaging, LoadBalancerActivationNetworkProduction).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", v.Network))),
		"Version": validation.Validate(v.Version, validation.Min(0)),
	}.Filter()
}

// ListLoadBalancerActivations fetches activations with the most recent listed first
func (c *cloudlets) ListLoadBalancerActivations(ctx context.Context, params ListLoadBalancerActivationsRequest) ([]LoadBalancerActivation, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLoadBalancerActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListLoadBalancerActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/activations", params.OriginID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListLoadBalancerActivations, err)
	}

	q := uri.Query()
	if params.Network != "" {
		q.Add("network", fmt.Sprintf("%s", params.Network))
	}
	if params.PageSize != nil {
		q.Add("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}
	if params.Page != nil {
		q.Add("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.LatestOnly != false {
		q.Add("latestOnly", fmt.Sprintf("%s", strconv.FormatBool(params.LatestOnly)))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListLoadBalancerActivations, err)
	}

	var result []LoadBalancerActivation
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListLoadBalancerActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListLoadBalancerActivations, c.Error(resp))
	}

	return result, nil
}

// ActivateLoadBalancerVersion activates the load balacing version
func (c *cloudlets) ActivateLoadBalancerVersion(ctx context.Context, params ActivateLoadBalancerVersionRequest) (*LoadBalancerActivation, error) {
	logger := c.Log(ctx)
	logger.Debug("ActivateLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrActivateLoadBalancerVersion, ErrStructValidation, err)
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

	var result LoadBalancerActivation

	resp, err := c.Exec(req, &result, params.LoadBalancerVersionActivation)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateLoadBalancerVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrActivateLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}
