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
	ReputationProfileAction interface {
		// GetReputationProfileActions returns a list of reputation profiles with their associated actions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-reputation-profiles-actions
		GetReputationProfileActions(ctx context.Context, params GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error)

		// GetReputationProfileAction returns the action a reputation profile takes when triggered.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-reputation-profile-action
		GetReputationProfileAction(ctx context.Context, params GetReputationProfileActionRequest) (*GetReputationProfileActionResponse, error)

		// UpdateReputationProfileAction updates what action to take when reputation profile's rule triggers.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-reputation-profile-action
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
	logger := p.Log(ctx)
	logger.Debug("GetReputationProfileAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetReputationProfileActionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation profile action request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetReputationProfileActions(ctx context.Context, params GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetReputationProfileActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-profiles",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetReputationProfileActions request: %w", err)
	}

	var result GetReputationProfileActionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation profile actions request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ReputationProfileID != 0 {
		var filteredResult GetReputationProfileActionsResponse
		for _, val := range result.ReputationProfiles {
			if val.ID == params.ReputationProfileID {
				filteredResult.ReputationProfiles = append(filteredResult.ReputationProfiles, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateReputationProfileAction(ctx context.Context, params UpdateReputationProfileActionRequest) (*UpdateReputationProfileActionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateReputationProfileAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/reputation-profiles/%d",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.ReputationProfileID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateReputationProfileAction request: %w", err)
	}

	var result UpdateReputationProfileActionResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update reputation profile action request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
