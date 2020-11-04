package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RatePolicyAction represents a collection of RatePolicyAction
//
// See: RatePolicyAction.GetRatePolicyAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RatePolicyAction  contains operations available on RatePolicyAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicyaction
	RatePolicyAction interface {
		GetRatePolicyActions(ctx context.Context, params GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error)
		GetRatePolicyAction(ctx context.Context, params GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error)
		UpdateRatePolicyAction(ctx context.Context, params UpdateRatePolicyActionRequest) (*UpdateRatePolicyActionResponse, error)
	}

	GetRatePolicyActionsRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyID   string `json:"policyId"`
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	GetRatePolicyActionsResponse struct {
		RatePolicyActions []struct {
			ID         int    `json:"id"`
			Ipv4Action string `json:"ipv4Action"`
			Ipv6Action string `json:"ipv6Action"`
		} `json:"ratePolicyActions"`
	}

	GetRatePolicyActionRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		PolicyID   string `json:"policyId"`
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	GetRatePolicyActionResponse struct {
		RatePolicyActions []struct {
			ID         int    `json:"id"`
			Ipv4Action string `json:"ipv4Action"`
			Ipv6Action string `json:"ipv6Action"`
		} `json:"ratePolicyActions"`
	}

	UpdateRatePolicyActionRequest struct {
		ConfigID     int    `json:"-"`
		Version      int    `json:"-"`
		PolicyID     string `json:"-"`
		RatePolicyID int    `json:"-"`
		Ipv4Action   string `json:"ipv4Action"`
		Ipv6Action   string `json:"ipv6Action"`
	}

	UpdateRatePolicyActionResponse struct {
		RatePolicyActions []struct {
			ID         int    `json:"-"`
			Ipv4Action string `json:"ipv4Action"`
			Ipv6Action string `json:"ipv6Action"`
		} `json:"ratePolicyActions"`
	}

	RatePolicyActionPost struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}
)

// Validate validates GetRatePolicyActionRequest
func (v GetRatePolicyActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetRatePolicyActionsRequest
func (v GetRatePolicyActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRatePolicyActionRequest
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
		return nil, fmt.Errorf("failed to create getratepolicyaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
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

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rate-policies",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getratepolicyactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getratepolicyactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a RatePolicyAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putratepolicyaction

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
		return nil, fmt.Errorf("failed to create create RatePolicyActionrequest: %w", err)
	}

	var rval UpdateRatePolicyActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RatePolicyAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
