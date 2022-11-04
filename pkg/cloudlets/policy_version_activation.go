package cloudlets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PolicyVersionActivations is a cloudlets PolicyVersionActivations API interface
	PolicyVersionActivations interface {
		// ListPolicyActivations returns the complete activation history for the selected policy in reverse chronological order.
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyactivations
		ListPolicyActivations(context.Context, ListPolicyActivationsRequest) ([]PolicyActivation, error)

		// ActivatePolicyVersion activates the selected cloudlet policy version.
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postpolicyversionactivations
		ActivatePolicyVersion(context.Context, ActivatePolicyVersionRequest) error
	}

	// ListPolicyActivationsRequest contains the request parameters for ListPolicyActivations
	ListPolicyActivationsRequest struct {
		PolicyID     int64
		Network      PolicyActivationNetwork
		PropertyName string
	}

	// ActivatePolicyVersionRequest contains the request parameters for ActivatePolicyVersion
	ActivatePolicyVersionRequest struct {
		PolicyID int64
		Async    bool
		Version  int64
		PolicyVersionActivation
	}

	// PolicyVersionActivation is the body content for an ActivatePolicyVersion request
	PolicyVersionActivation struct {
		Network                 PolicyActivationNetwork `json:"network"`
		AdditionalPropertyNames []string                `json:"additionalPropertyNames,omitempty"`
	}

	// PolicyActivationNetwork is the activation network type for policy
	PolicyActivationNetwork string
)

var (
	// ErrListPolicyActivations is returned when ListPolicyActivations fails
	ErrListPolicyActivations = errors.New("list policy activations")
	// ErrActivatePolicyVersion is returned when ActivatePolicyVersion fails
	ErrActivatePolicyVersion = errors.New("activate policy version")
)

const (
	// PolicyActivationNetworkStaging is the staging network for policy
	PolicyActivationNetworkStaging PolicyActivationNetwork = "staging"
	// PolicyActivationNetworkProduction is the production network for policy
	PolicyActivationNetworkProduction PolicyActivationNetwork = "prod"
)

// Validate validates ListPolicyActivationsRequest
func (r ListPolicyActivationsRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Network": validation.Validate(
			r.Network,
			validation.In(PolicyActivationNetworkStaging, PolicyActivationNetworkProduction).Error(
				fmt.Sprintf("value '%s' is invalid. Must be one of: 'staging', 'prod' or '' (empty)", (&r).Network)),
		),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates ActivatePolicyVersionRequest
func (r ActivatePolicyVersionRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":                            validation.Validate(r.PolicyID, validation.Required),
		"Version":                             validation.Validate(r.Version, validation.Required),
		"RequestBody.AdditionalPropertyNames": validation.Validate(r.PolicyVersionActivation.AdditionalPropertyNames, validation.Required),
		"RequestBody.Network": validation.Validate(
			r.PolicyVersionActivation.Network,
			validation.Required,
			validation.In(PolicyActivationNetworkStaging, PolicyActivationNetworkProduction).Error(
				fmt.Sprintf("value '%s' is invalid. Must be one of: 'staging' or 'prod'", (&r).PolicyVersionActivation.Network)),
		),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// UnmarshalJSON unifies json network field into well defined values
func (n *PolicyActivationNetwork) UnmarshalJSON(data []byte) error {
	d := bytes.Trim(data, "\"")

	switch string(d) {
	case "STAGING", "staging":
		*n = PolicyActivationNetworkStaging
	case "PRODUCTION", "production", "prod":
		*n = PolicyActivationNetworkProduction
	default:
		return fmt.Errorf("cannot unmarshall PolicyActivationNetwork: %q", d)
	}
	return nil
}

func (c *cloudlets) ListPolicyActivations(ctx context.Context, params ListPolicyActivationsRequest) ([]PolicyActivation, error) {
	c.Log(ctx).Debug("ListPolicyActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListPolicyActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/activations", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicyActivations, err)
	}

	q := uri.Query()
	if params.Network != "" {
		q.Set("network", string(params.Network))
	}
	if params.PropertyName != "" {
		q.Set("propertyName", params.PropertyName)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicyActivations, err)
	}

	var result []PolicyActivation
	response, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicyActivations, err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicyActivations, c.Error(response))
	}

	return result, nil
}

func (c *cloudlets) ActivatePolicyVersion(ctx context.Context, params ActivatePolicyVersionRequest) error {
	c.Log(ctx).Debug("ActivatePolicyVersion")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrActivatePolicyVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/cloudlets/api/v2/policies/%d/versions/%d/activations",
		params.PolicyID, params.Version))
	if err != nil {
		return fmt.Errorf("%w: failed to create POST URI: %s", ErrActivatePolicyVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create POST request: %s", ErrActivatePolicyVersion, err)
	}

	response, err := c.Exec(req, nil, params.PolicyVersionActivation)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrActivatePolicyVersion, err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %s", ErrActivatePolicyVersion, c.Error(response))
	}

	return nil
}
