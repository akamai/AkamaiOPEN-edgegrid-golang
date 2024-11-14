package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RatePolicyAction interface supports retrieving and modifying the action associated with
	// a specified rate policy, or with all rate policies in a security policy.
	RatePolicyAction interface {
		// GetRatePolicyActions returns a list of all rate policies currently in use with the actions each policy takes when conditions are met.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rate-policies-actions
		GetRatePolicyActions(ctx context.Context, params GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error)

		// GetRatePolicyAction returns a specified rate policy currently in use with the action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rate-policies-actions
		// Deprecated: this method will be removed in a future release. Use GetRatePolicyActions instead.
		GetRatePolicyAction(ctx context.Context, params GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error)

		// UpdateRatePolicyAction
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rate-policy-action
		UpdateRatePolicyAction(ctx context.Context, params UpdateRatePolicyActionRequest) (*UpdateRatePolicyActionResponse, error)
	}

	// GetRatePolicyActionsRequest is used to retrieve a configuration's rate policies and their associated actions.
	GetRatePolicyActionsRequest struct {
		ConfigID     int    `json:"configId"`
		Version      int    `json:"version"`
		PolicyID     string `json:"policyId"`
		RatePolicyID int    `json:"id"`
		Ipv4Action   string `json:"ipv4Action"`
		Ipv6Action   string `json:"ipv6Action"`
	}

	// GetRatePolicyActionsResponse is returned from a call to GetRatePolicyActions.
	GetRatePolicyActionsResponse struct {
		RatePolicyActions []struct {
			ID         int    `json:"id"`
			Ipv4Action string `json:"ipv4Action,omitempty"`
			Ipv6Action string `json:"ipv6Action,omitempty"`
		} `json:"ratePolicyActions,omitempty"`
	}

	// GetRatePolicyActionRequest is used to retrieve a configuration's rate policies and their associated actions.
	// Deprecated: this struct will be removed in a future release.
	GetRatePolicyActionRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyID   string `json:"policyId"`
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	// GetRatePolicyActionResponse is returned from a call to GetRatePolicyAction.
	// Deprecated: this struct will be removed in a future release.
	GetRatePolicyActionResponse struct {
		RatePolicyActions []struct {
			ID         int    `json:"id"`
			Ipv4Action string `json:"ipv4Action,omitempty"`
			Ipv6Action string `json:"ipv6Action,omitempty"`
		} `json:"ratePolicyActions"`
	}

	// UpdateRatePolicyActionRequest is used to update the actions for a rate policy.
	UpdateRatePolicyActionRequest struct {
		ConfigID     int    `json:"-"`
		Version      int    `json:"-"`
		PolicyID     string `json:"-"`
		RatePolicyID int    `json:"-"`
		Ipv4Action   string `json:"ipv4Action"`
		Ipv6Action   string `json:"ipv6Action"`
	}

	// UpdateRatePolicyActionResponse is returned from a call to UpdateRatePolicy.
	UpdateRatePolicyActionResponse struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	// RatePolicyActionPost is used to describe actions that may be taken as part of a rate policy.
	RatePolicyActionPost struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}
)

// Validate validates a GetRatePolicyActionRequest.
// Deprecated: this method will be removed in a future release.
func (v GetRatePolicyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetRatePolicyActionsRequest.
func (v GetRatePolicyActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRatePolicyActionRequest.
func (v UpdateRatePolicyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":     validation.Validate(v.ConfigID, validation.Required),
		"Version":      validation.Validate(v.Version, validation.Required),
		"PolicyID":     validation.Validate(v.PolicyID, validation.Required),
		"RatePolicyID": validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) GetRatePolicyAction(ctx context.Context, params GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRatePolicyAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicyAction request: %w", err)
	}

	var result GetRatePolicyActionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rate policy action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetRatePolicyActions(ctx context.Context, params GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRatePolicyActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicyActions request: %w", err)
	}

	var result GetRatePolicyActionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rate policy actions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RatePolicyID != 0 {
		var filteredResult GetRatePolicyActionsResponse
		for _, val := range result.RatePolicyActions {
			if val.ID == params.RatePolicyID {
				filteredResult.RatePolicyActions = append(filteredResult.RatePolicyActions, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateRatePolicyAction(ctx context.Context, params UpdateRatePolicyActionRequest) (*UpdateRatePolicyActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRatePolicyAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RatePolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRatePolicyAction request: %w", err)
	}

	var result UpdateRatePolicyActionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update rate policy action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
