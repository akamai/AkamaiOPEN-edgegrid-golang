package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ReputationProfileAction interface supports retrieving and modifying the action associated with
	// a specified reputation profile, or with all reputation profiles in a security policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#reputationprofileactiongroup
	ReputationProfileAction interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationprofileactions
		GetReputationProfileActions(ctx context.Context, params GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationprofileaction
		GetReputationProfileAction(ctx context.Context, params GetReputationProfileActionRequest) (*GetReputationProfileActionResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putreputationprofileaction
		UpdateReputationProfileAction(ctx context.Context, params UpdateReputationProfileActionRequest) (*UpdateReputationProfileActionResponse, error)
	}

	// GetReputationProfileActionsRequest is used to retrieve the list of reputation profiles and their associated actions.
	GetReputationProfileActionsRequest struct {
		ConfigID            int    `json:"configId"`
		Version             int    `json:"version"`
		PolicyID            string `json:"policyId"`
		ReputationProfileID int    `json:"id"`
		Action              string `json:"action"`
	}

	// GetReputationProfileActionsResponse is returned from a call to GetReputationProfileActions.
	GetReputationProfileActionsResponse struct {
		ReputationProfiles []struct {
			Action string `json:"action,omitempty"`
			ID     int    `json:"id,omitempty"`
		} `json:"reputationProfiles,omitempty"`
	}

	// GetReputationProfileActionRequest is used to retrieve the details for a specific reputation profile.
	GetReputationProfileActionRequest struct {
		ConfigID            int    `json:"configId"`
		Version             int    `json:"version"`
		PolicyID            string `json:"policyId"`
		ReputationProfileID int    `json:"id"`
		Action              string `json:"action"`
	}

	// GetReputationProfileActionResponse is returned from a call to GetReputationProfileAction.
	GetReputationProfileActionResponse struct {
		Action string `json:"action,omitempty"`
	}

	// UpdateReputationProfileActionRequest is used to modify the details for a specific reputation profile.
	UpdateReputationProfileActionRequest struct {
		ConfigID            int    `json:"-"`
		Version             int    `json:"-"`
		PolicyID            string `json:"-"`
		ReputationProfileID int    `json:"-"`
		Action              string `json:"action"`
	}

	// UpdateReputationProfileActionResponse is returned from a call to UpdateReputationProfileAction.
	UpdateReputationProfileActionResponse struct {
		Action string `json:"action"`
	}

	// ReputationProfileActionPost is currently unused.
	ReputationProfileActionPost struct {
		Action string `json:"action"`
	}
)

// Validate validates a GetReputationProfileActionRequest.
func (v GetReputationProfileActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetReputationProfileActionsRequest.
func (v GetReputationProfileActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateReputationProfileActionRequest.
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
		return nil, fmt.Errorf("failed to create GetReputationProfileAction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetReputationProfileAction request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetReputationProfileActions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetReputationProfileActions request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create UpdateReputationProfileAction request: %w", err)
	}

	var rval UpdateReputationProfileActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateReputationProfileAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
