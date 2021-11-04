package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// VersionActivationNetwork is the activation network value
	VersionActivationNetwork string

	// PolicyActivation is returned by ListPolicyActivations
	PolicyActivation struct {
		Network      VersionActivationNetwork `json:"network,omitempty"`
		APIVersion   string                   `json:"apiVersion,omitempty"`
		PolicyInfo   PolicyInfo               `json:"policyInfo"`
		PropertyInfo PropertyInfo             `json:"propertyInfo"`
	}

	// ListPolicyActivationsRequest contains the request parameters for ListPolicyActivations
	ListPolicyActivationsRequest struct {
		PolicyID     int64
		Network      VersionActivationNetwork
		PropertyName string
	}

	// ActivatePolicyVersionRequest contains the request parameters for ActivatePolicyVersion
	ActivatePolicyVersionRequest struct {
		PolicyID    int64
		Async       bool
		Version     int64
		RequestBody ActivatePolicyVersionRequestBody
	}

	// ActivatePolicyVersionRequestBody is the body content for an ActivatePolicyVersionRequest
	ActivatePolicyVersionRequestBody struct {
		Network                 VersionActivationNetwork `json:"network"`
		AdditionalPropertyNames []string                 `json:"additionalPropertyNames,omitempty"`
	}

	// PolicyVersionActivation is a cloudlets PolicyVersionActivation API interface
	PolicyVersionActivation interface {

		// ListPolicyActivations returns the complete activation history for the selected policy in reverse chronological order.
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyactivations
		ListPolicyActivations(context.Context, ListPolicyActivationsRequest) ([]PolicyActivation, error)

		// ActivatePolicyVersion activates the selected cloudlet policy version.
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postpolicyversionactivations
		ActivatePolicyVersion(context.Context, ActivatePolicyVersionRequest) error
	}
)

var (
	// ErrListPolicyActivations is returned when ListPolicyActivations fails
	ErrListPolicyActivations = errors.New("list policy activations")
	// ErrActivatePolicyVersion is returned when ActivatePolicyVersion fails
	ErrActivatePolicyVersion = errors.New("activate policy version")
)

// Validate validates ListPolicyActivationsRequest
func (r ListPolicyActivationsRequest) Validate() error {
	return validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Network": validation.Validate(
			r.Network,
			validation.In(VersionActivationNetworkStaging, VersionActivationNetworkProduction),
		),
	}.Filter()
}

// Validate validates ActivatePolicyVersionRequest
func (r ActivatePolicyVersionRequest) Validate() error {
	return validation.Errors{
		"PolicyID":                            validation.Validate(r.PolicyID, validation.Required),
		"Version":                             validation.Validate(r.Version, validation.Required),
		"RequestBody.AdditionalPropertyNames": validation.Validate(r.RequestBody.AdditionalPropertyNames, validation.Required),
		"RequestBody.Network": validation.Validate(
			r.RequestBody.Network,
			validation.In(VersionActivationNetworkStaging, VersionActivationNetworkProduction),
		),
	}.Filter()
}

func (c *cloudlets) ListPolicyActivations(ctx context.Context, params ListPolicyActivationsRequest) ([]PolicyActivation, error) {
	c.Log(ctx).Debug("ListPolicyActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPolicyActivations, ErrStructValidation, err)
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
		return fmt.Errorf("%s: %w: %s", ErrActivatePolicyVersion, ErrStructValidation, err)
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

	response, err := c.Exec(req, nil, params.RequestBody)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrActivatePolicyVersion, err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %s", ErrActivatePolicyVersion, c.Error(response))
	}

	return nil
}

const (
	// VersionActivationNetworkStaging is the staging network
	VersionActivationNetworkStaging VersionActivationNetwork = "staging"

	// VersionActivationNetworkProduction is the production network
	VersionActivationNetworkProduction VersionActivationNetwork = "prod"
)
