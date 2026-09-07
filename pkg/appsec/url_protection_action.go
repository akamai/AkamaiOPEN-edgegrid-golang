package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The URLProtectionAction interface supports retrieving and modifying the actions associated with
	// a specified url protection policy, or with all url protection policies in a security policy.
	URLProtectionAction interface {
		// ListURLProtectionPoliciesActions returns a list of all url protections policies currently in use with the actions each policy takes.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-url-protection-policies-actions
		ListURLProtectionPoliciesActions(ctx context.Context, params ListURLProtectionPoliciesActionsRequest) (*ListURLProtectionPoliciesActionsResponse, error)

		// GetURLProtectionPolicyActions returns a specific url protections policy currently in use with the actions.
		GetURLProtectionPolicyActions(ctx context.Context, params GetURLProtectionPolicyActionsRequest) (*GetURLProtectionPolicyActionsResponse, error)

		// UpdateURLProtectionPolicyActions updates the actions for a url protection policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-url-protection-policy-action
		UpdateURLProtectionPolicyActions(ctx context.Context, params UpdateURLProtectionPolicyActionsRequest) (*UpdateURLProtectionPolicyActionsResponse, error)
	}

	// URLProtectionPolicyActions represents the actions associated with a specific URL protection policy.
	URLProtectionPolicyActions struct {

		// MaxRateThresholdAction specifies the action to take when the max rate threshold is exceeded (e.g., "alert", "deny", "none", "challenge_{id}", "deny_custom_{id}")
		MaxRateThresholdAction string `json:"action"`

		// LoadSheddingAction specifies the action to take for load shedding (e.g., "alert", "deny", "challenge_{id}", "deny_custom_{id}")
		LoadSheddingAction string `json:"loadSheddingAction,omitempty"`
	}

	// ListURLProtectionPoliciesActionsRequest is used to retrieve a configuration's url protection policies and their associated actions.
	ListURLProtectionPoliciesActionsRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the ConfigVersion number of the security configuration
		ConfigVersion int64

		// PolicyID is the unique identifier of the security policy
		PolicyID string
	}

	// URLProtectionActionResp URLProtectionAction represents the action associated with a specific URL protection policy.
	URLProtectionActionResp struct {

		// URLProtectionPolicyID is the unique identifier for the URL protection policy
		URLProtectionPolicyID int64 `json:"policyId"`

		// MaxRateThresholdAction specifies the action to take when the max rate threshold is exceeded (e.g., "alert", "deny", "none", "challenge_{id}", "deny_custom_{id}")
		MaxRateThresholdAction string `json:"action"`

		// LoadSheddingAction specifies the action to take for load shedding (e.g., "alert", "deny","challenge_{id}","deny_custom_{id}")
		LoadSheddingAction string `json:"loadSheddingAction"`
	}

	// ListURLProtectionPoliciesActionsResponse is returned from a call to ListURLProtectionPoliciesActions.
	ListURLProtectionPoliciesActionsResponse struct {

		// URLProtectionPoliciesActions is the list of URL protection policies with their associated actions
		URLProtectionPoliciesActions []URLProtectionActionResp `json:"urlProtectionActions"`
	}

	// GetURLProtectionPolicyActionsRequest is used to retrieve the actions associated with the particular url protection of specific policy in a config ConfigVersion.
	GetURLProtectionPolicyActionsRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the ConfigVersion number of the security configuration
		ConfigVersion int64

		// PolicyID is the unique identifier of the security policy
		PolicyID string

		// URLProtectionPolicyID is the unique identifier of the URL protection policy
		URLProtectionPolicyID int64
	}

	// GetURLProtectionPolicyActionsResponse is returned from a call to GetURLProtectionPolicyActions.
	GetURLProtectionPolicyActionsResponse URLProtectionPolicyActions

	// UpdateURLProtectionPolicyActionsRequest is used to update the actions for a url protection policy.
	UpdateURLProtectionPolicyActionsRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the ConfigVersion number of the security configuration
		ConfigVersion int64

		// PolicyID is the unique identifier of the security policy
		PolicyID string

		// URLProtectionPolicyID is the unique identifier of the URL protection policy
		URLProtectionPolicyID int64

		// Body contains the URL protection policy actions to be updated
		Body URLProtectionPolicyActions
	}

	// UpdateURLProtectionPolicyActionsResponse is returned from a call to UpdateURLProtectionPolicyActions.
	UpdateURLProtectionPolicyActionsResponse URLProtectionPolicyActions
)

// Validate validates a ListURLProtectionPoliciesActionsRequest.
func (v ListURLProtectionPoliciesActionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"PolicyID":      validation.Validate(v.PolicyID, validation.Required),
	})
}

// Validate validates a GetURLProtectionPolicyActionsRequest.
func (v GetURLProtectionPolicyActionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":              validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":         validation.Validate(v.ConfigVersion, validation.Required),
		"PolicyID":              validation.Validate(v.PolicyID, validation.Required),
		"URLProtectionPolicyID": validation.Validate(v.URLProtectionPolicyID, validation.Required),
	})
}

// Validate validates an UpdateURLProtectionPolicyActions.
func (v URLProtectionPolicyActions) Validate() error {
	return (validation.Errors{
		"MaxRateThresholdAction": validation.Validate(v.MaxRateThresholdAction, validation.Required),
	}).Filter()
}

// Validate validates an UpdateURLProtectionPolicyActionRequest.
func (v UpdateURLProtectionPolicyActionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":              validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":         validation.Validate(v.ConfigVersion, validation.Required),
		"PolicyID":              validation.Validate(v.PolicyID, validation.Required),
		"URLProtectionPolicyID": validation.Validate(v.URLProtectionPolicyID, validation.Required),
		"Body":                  validation.Validate(v.Body, validation.Required),
	})
}

func (p *appsec) ListURLProtectionPoliciesActions(ctx context.Context, params ListURLProtectionPoliciesActionsRequest) (*ListURLProtectionPoliciesActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListURLProtectionPoliciesActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/url-protections",
		params.ConfigID,
		params.ConfigVersion,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListURLProtectionPoliciesActions request: %w", err)
	}

	var result ListURLProtectionPoliciesActionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get url protection policies actions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetURLProtectionPolicyActions(ctx context.Context, params GetURLProtectionPolicyActionsRequest) (*GetURLProtectionPolicyActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetURLProtectionPolicyActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/url-protections",
		params.ConfigID,
		params.ConfigVersion,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetURLProtectionPolicyActions request: %w", err)
	}

	var result ListURLProtectionPoliciesActionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get url protection policy actions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	var filteredResult GetURLProtectionPolicyActionsResponse
	for _, val := range result.URLProtectionPoliciesActions {
		if val.URLProtectionPolicyID == params.URLProtectionPolicyID {
			filteredResult.MaxRateThresholdAction = val.MaxRateThresholdAction
			filteredResult.LoadSheddingAction = val.LoadSheddingAction
			return &filteredResult, nil
		}
	}
	return nil, fmt.Errorf("incorrect URLProtectionPolicyID %d or no actions found for the specified policy", params.URLProtectionPolicyID)
}

func (p *appsec) UpdateURLProtectionPolicyActions(ctx context.Context, params UpdateURLProtectionPolicyActionsRequest) (*UpdateURLProtectionPolicyActionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateURLProtectionPolicyActions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/url-protections/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.PolicyID,
		params.URLProtectionPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateURLProtectionPolicyActions request: %w", err)
	}

	var result UpdateURLProtectionPolicyActionsResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update url protection policy action request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
