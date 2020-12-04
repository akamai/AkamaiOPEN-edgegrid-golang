package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ReputationProfileAction represents a collection of ReputationProfileAction
//
// See: ReputationProfileAction.GetReputationProfileAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ReputationProfileAction  contains operations available on ReputationProfileAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationprofileaction
	ReputationProfileAction interface {
		GetReputationProfileActions(ctx context.Context, params GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error)
		GetReputationProfileAction(ctx context.Context, params GetReputationProfileActionRequest) (*GetReputationProfileActionResponse, error)
		UpdateReputationProfileAction(ctx context.Context, params UpdateReputationProfileActionRequest) (*UpdateReputationProfileActionResponse, error)
	}

	GetReputationProfileActionsRequest struct {
		ConfigID            int    `json:"configId"`
		Version             int    `json:"version"`
		PolicyID            string `json:"policyId"`
		ReputationProfileID int    `json:"id"`
		Action              string `json:"action"`
	}

	GetReputationProfileActionsResponse struct {
		ReputationProfiles []struct {
			Action string `json:"action"`
			ID     int    `json:"id"`
		} `json:"reputationProfiles"`
	}

	GetReputationProfileActionRequest struct {
		ConfigID            int    `json:"configId"`
		Version             int    `json:"version"`
		PolicyID            string `json:"policyId"`
		ReputationProfileID int    `json:"id"`
		Action              string `json:"action"`
	}

	GetReputationProfileActionResponse struct {
		Action string `json:"action"`
	}

	UpdateReputationProfileActionRequest struct {
		ConfigID            int    `json:"-"`
		Version             int    `json:"-"`
		PolicyID            string `json:"-"`
		ReputationProfileID int    `json:"-"`
		Action              string `json:"action"`
	}

	UpdateReputationProfileActionResponse struct {
		Action string `json:"action"`
	}

	ReputationProfileActionPost struct {
		action string `json:"action"`
	}
)

// Validate validates GetReputationProfileActionRequest
func (v GetReputationProfileActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetReputationProfileActionsRequest
func (v GetReputationProfileActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateReputationProfileActionRequest
func (v UpdateReputationProfileActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":            validation.Validate(v.ConfigID, validation.Required),
		"Version":             validation.Validate(v.Version, validation.Required),
		"PolicyID":            validation.Validate(v.PolicyID, validation.Required),
		"ReputationProfileID": validation.Validate(v.ReputationProfileID, validation.Required),
	}.Filter()
}

func (p *appsec) GetReputationProfileAction(ctx context.Context, params GetReputationProfileActionRequest) (*GetReputationProfileActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetReputationProfileAction")

	var rval GetReputationProfileActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-profiles/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.ReputationProfileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getreputationprofileaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getreputationprofileaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetReputationProfileActions(ctx context.Context, params GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetReputationProfileActions")

	var rval GetReputationProfileActionsResponse
	var rvalfiltered GetReputationProfileActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-profiles",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getreputationprofileactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getreputationprofileactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ReputationProfileID != 0 {
		for _, val := range rval.ReputationProfiles {
			if val.ID == params.ReputationProfileID {
				rvalfiltered.ReputationProfiles = append(rvalfiltered.ReputationProfiles, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a ReputationProfileAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putreputationprofileaction

func (p *appsec) UpdateReputationProfileAction(ctx context.Context, params UpdateReputationProfileActionRequest) (*UpdateReputationProfileActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateReputationProfileAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-profiles/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.ReputationProfileID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ReputationProfileActionrequest: %w", err)
	}

	var rval UpdateReputationProfileActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create ReputationProfileAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
