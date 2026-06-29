package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
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

// Validate validates LoadBalancerVersionActivation struct
func (v LoadBalancerVersionActivation) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(v.Network, validation.Required, validation.In(LoadBalancerActivationNetworkStaging, LoadBalancerActivationNetworkProduction).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", v.Network))),
		"Version": validation.Validate(v.Version, validation.Min(0)),
	}.Filter()
}

func (c *cloudlets) ListLoadBalancerActivations(ctx context.Context, params ListLoadBalancerActivationsRequest) ([]LoadBalancerActivation, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLoadBalancerActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListLoadBalancerActivations, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cloudlets/api/v2/origins/%s/activations", params.OriginID).
		AddQueryParamIf("network", string(params.Network), params.Network != "").
		AddQueryParamIf("latestOnly", strconv.FormatBool(params.LatestOnly), params.LatestOnly).
		AddQueryParamFunc("pageSize", func() string {
			return fmt.Sprintf("%d", *params.PageSize)
		}, params.PageSize != nil).
		AddQueryParamFunc("page", func() string {
			return fmt.Sprintf("%d", *params.Page)
		}, params.Page != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListLoadBalancerActivations, err)
	}

	var result []LoadBalancerActivation
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListLoadBalancerActivations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListLoadBalancerActivations, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) ActivateLoadBalancerVersion(ctx context.Context, params ActivateLoadBalancerVersionRequest) (*LoadBalancerActivation, error) {
	logger := c.Log(ctx)
	logger.Debug("ActivateLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrActivateLoadBalancerVersion, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/cloudlets/api/v2/origins/%s/activations", params.OriginID).
		AddQueryParam("async", strconv.FormatBool(params.Async)).
		WithBody(params.LoadBalancerVersionActivation).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateLoadBalancerVersion, err)
	}

	var result LoadBalancerActivation

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateLoadBalancerVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrActivateLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}
