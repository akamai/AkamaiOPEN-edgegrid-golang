package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ListPolicyActivationsRequest contains request parameters for ListPolicyActivations.
type ListPolicyActivationsRequest struct {
	PolicyID int64
	Page     int
	Size     int
}

// GetPolicyActivationRequest contains request parameters for GetPolicyActivation.
type GetPolicyActivationRequest struct {
	PolicyID     int64
	ActivationID int64
}

// ActivatePolicyRequest contains request parameters for ActivatePolicy.
type ActivatePolicyRequest struct {
	PolicyID      int64
	Network       Network
	PolicyVersion int64
}

// DeactivatePolicyRequest contains request parameters for DeactivatePolicy.
type DeactivatePolicyRequest struct {
	PolicyID      int64
	Network       Network
	PolicyVersion int64
}

type policyActivationRequest struct {
	Operation     PolicyActivationOperation `json:"operation"`
	Network       Network                   `json:"network"`
	PolicyVersion int64                     `json:"policyVersion"`
}

// PolicyActivation represents a single policy activation.
type PolicyActivation struct {
	CreatedBy            string                    `json:"createdBy"`
	CreatedDate          time.Time                 `json:"createdDate"`
	FinishDate           *time.Time                `json:"finishDate"`
	ID                   int64                     `json:"id"`
	Network              Network                   `json:"network"`
	Operation            PolicyActivationOperation `json:"operation"`
	PolicyID             int64                     `json:"policyId"`
	Status               ActivationStatus          `json:"status"`
	PolicyVersion        int64                     `json:"policyVersion"`
	PolicyVersionDeleted bool                      `json:"policyVersionDeleted"`
	Links                []Link                    `json:"links"`
}

// PolicyActivations represents the response data from ListPolicyActivations.
type PolicyActivations struct {
	Page              Page               `json:"page"`
	PolicyActivations []PolicyActivation `json:"content"`
	Links             []Link             `json:"links"`
}

// PolicyActivationOperation is an enum for policy activation operation
type PolicyActivationOperation string

const (
	// OperationActivation represents an operation used for activating a policy
	OperationActivation PolicyActivationOperation = "ACTIVATION"
	// OperationDeactivation represents an operation used for deactivating a policy
	OperationDeactivation PolicyActivationOperation = "DEACTIVATION"
)

// ActivationStatus represents information about policy activation status.
type ActivationStatus string

const (
	// ActivationStatusInProgress informs that activation is in progress.
	ActivationStatusInProgress ActivationStatus = "IN_PROGRESS"
	// ActivationStatusSuccess informs that activation succeeded.
	ActivationStatusSuccess ActivationStatus = "SUCCESS"
	// ActivationStatusFailed informs that activation failed.
	ActivationStatusFailed ActivationStatus = "FAILED"
)

// Network represents network on which policy version or property can be activated on.
type Network string

const (
	// StagingNetwork represents staging network.
	StagingNetwork Network = "STAGING"
	// ProductionNetwork represents production network.
	ProductionNetwork Network = "PRODUCTION"
)

var (
	// ErrListPolicyActivations is returned when ListPolicyActivations fails.
	ErrListPolicyActivations = errors.New("list policy activations")
	// ErrActivatePolicy is returned when ActivatePolicy fails.
	ErrActivatePolicy = errors.New("activate policy")
	// ErrDeactivatePolicy is returned when DeactivatePolicy fails.
	ErrDeactivatePolicy = errors.New("deactivate policy")
	// ErrGetPolicyActivation is returned when GetPolicyActivation fails.
	ErrGetPolicyActivation = errors.New("get policy activation")
)

// Validate validates ListPolicyActivationsRequest.
func (r ListPolicyActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Page":     validation.Validate(r.Page, validation.Min(0)),
		"Size":     validation.Validate(r.Size, validation.Min(10)),
	})
}

// Validate validates GetPolicyActivationRequest.
func (r GetPolicyActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":     validation.Validate(r.PolicyID, validation.Required),
		"ActivationID": validation.Validate(r.ActivationID, validation.Required),
	})
}

// Validate validates ActivatePolicyRequest.
func (r ActivatePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":      validation.Validate(r.PolicyID, validation.Required),
		"PolicyVersion": validation.Validate(r.PolicyVersion, validation.Required),
		"Network": validation.Validate(r.Network, validation.Required, validation.In(StagingNetwork, ProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", r.Network))),
	})
}

// Validate validates DeactivatePolicyRequest.
func (r DeactivatePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":      validation.Validate(r.PolicyID, validation.Required),
		"PolicyVersion": validation.Validate(r.PolicyVersion, validation.Required),
		"Network": validation.Validate(r.Network, validation.Required, validation.In(StagingNetwork, ProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'STAGING' or 'PRODUCTION'", r.Network))),
	})
}

func (c *cloudlets) ListPolicyActivations(ctx context.Context, params ListPolicyActivationsRequest) (*PolicyActivations, error) {
	c.Log(ctx).Debug("ListPolicyActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPolicyActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/v3/policies/%d/activations", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicyActivations, err)
	}

	q := uri.Query()
	if params.Size != 0 {
		q.Add("size", strconv.Itoa(params.Size))
	}
	if params.Page != 0 {
		q.Add("page", strconv.Itoa(params.Page))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicyActivations, err)
	}

	var result PolicyActivations
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicyActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicyActivations, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) ActivatePolicy(ctx context.Context, params ActivatePolicyRequest) (*PolicyActivation, error) {
	c.Log(ctx).Debug("ActivatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivatePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/activations", params.PolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivatePolicy, err)
	}

	reqBody := policyActivationRequest{
		Network:       params.Network,
		PolicyVersion: params.PolicyVersion,
		Operation:     OperationActivation,
	}

	var result PolicyActivation
	resp, err := c.Exec(req, &result, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivatePolicy, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrActivatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) DeactivatePolicy(ctx context.Context, params DeactivatePolicyRequest) (*PolicyActivation, error) {
	c.Log(ctx).Debug("DeactivatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivatePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/activations", params.PolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivatePolicy, err)
	}

	reqBody := policyActivationRequest{
		Network:       params.Network,
		PolicyVersion: params.PolicyVersion,
		Operation:     OperationDeactivation,
	}

	var result PolicyActivation
	resp, err := c.Exec(req, &result, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivatePolicy, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrDeactivatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) GetPolicyActivation(ctx context.Context, params GetPolicyActivationRequest) (*PolicyActivation, error) {
	c.Log(ctx).Debug("GetPolicyActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPolicyActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/activations/%d", params.PolicyID, params.ActivationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyActivation, err)
	}

	var result PolicyActivation
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyActivation, c.Error(resp))
	}

	return &result, nil
}
