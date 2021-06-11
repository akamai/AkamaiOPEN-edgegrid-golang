package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RatePolicyAction interface supports retrieving and modifying the action associated with
	// a specified rate policy, or with all rate policies in a security policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#ratepolicyaction
	RatePolicyAction interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicyactions
		GetRatePolicyActions(ctx context.Context, params GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicyactions
		GetRatePolicyAction(ctx context.Context, params GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putactionsperratepolicy
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
	GetRatePolicyActionRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyID   string `json:"policyId"`
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	// GetRatePolicyActionResponse is returned from a call to GetRatePolicyAction.
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

func (p *appsec) GetRatePolicyAction(ctx context.Context, params GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRatePolicyAction")

	var rval GetRatePolicyActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicyAction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRatePolicyAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetRatePolicyActions(ctx context.Context, params GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRatePolicyActions")

	var rval GetRatePolicyActionsResponse
	var rvalfiltered GetRatePolicyActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicyActions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetRatePolicyActions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RatePolicyID != 0 {
		for _, val := range rval.RatePolicyActions {
			if val.ID == params.RatePolicyID {
				rvalfiltered.RatePolicyActions = append(rvalfiltered.RatePolicyActions, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

func (p *appsec) UpdateRatePolicyAction(ctx context.Context, params UpdateRatePolicyActionRequest) (*UpdateRatePolicyActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRatePolicyAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RatePolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRatePolicyAction: %w", err)
	}

	var rval UpdateRatePolicyActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRatePolicyAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
