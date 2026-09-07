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
	// The URLProtection interface supports creating, retrieving, updating and removing url protection policies.
	URLProtection interface {
		// ListURLProtectionPolicies returns url protection policies for a specific security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-url-protection-policies
		ListURLProtectionPolicies(ctx context.Context, params ListURLProtectionPoliciesRequest) (*ListURLProtectionPoliciesResponse, error)

		// GetURLProtectionPolicy returns the specified url protection policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-url-protection-policy
		GetURLProtectionPolicy(ctx context.Context, params GetURLProtectionPolicyRequest) (*GetURLProtectionPolicyResponse, error)

		// CreateURLProtectionPolicy creates a new url protection policy for a specific configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-url-protection-policies
		CreateURLProtectionPolicy(ctx context.Context, params CreateURLProtectionPolicyRequest) (*CreateURLProtectionPolicyResponse, error)

		// UpdateURLProtectionPolicy updates details for a specific url protection policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-url-protection-policy
		UpdateURLProtectionPolicy(ctx context.Context, params UpdateURLProtectionPolicyRequest) (*UpdateURLProtectionPolicyResponse, error)

		// RemoveURLProtectionPolicy deletes the specified url protection policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-url-protection-policy
		RemoveURLProtectionPolicy(ctx context.Context, params RemoveURLProtectionPolicyRequest) error
	}

	// HostnamePath is used to specify hostname and path combinations for URL protection.
	HostnamePath struct {

		// Hostname is the hostname to match on (e.g., "example.com")
		Hostname string `json:"hostname"`

		// Paths is the list of URL paths to match on for this hostname (e.g., ["/api", "/admin"])
		Paths []string `json:"paths"`
	}

	// APIDefinition is used to specify API definitions for URL protection policies.
	APIDefinition struct {

		// APIDefinitionID is the unique identifier for the API definition
		APIDefinitionID int64 `json:"apiDefinitionId"`

		// DefinedResources indicates whether to protect defined resources in the API definition
		DefinedResources bool `json:"definedResources"`

		// ResourceIDs lists the specific resource IDs to protect within the API definition
		ResourceIDs []int64 `json:"resourceIds"`

		// UndefinedResources indicates whether to protect undefined resources in the API definition
		UndefinedResources bool `json:"undefinedResources"`
	}

	// Category on which load shedding is performed when the origin traffic rate exceeds the load shedding threshold. If intelligentLoadShedding is set to true, specify one or more categories.
	Category struct {

		// Type determines the category type (e.g., "CLIENT_LIST","BOTS","CLOUD_PROVIDERS")
		Type string `json:"type"`

		// PositiveMatch indicates whether this is a positive match (true) or negative match (false)
		PositiveMatch *bool `json:"positiveMatch"`

		//ListIDs is the list of IDs associated with the category
		ListIDs []string `json:"listIds,omitempty"`
	}

	// AtomicCondition is used to add bypass conditions for URL protection policies.
	// Specify one or more types of conditions to match on. You can match on client lists, request headers, or both.
	AtomicCondition struct {

		// Type specifies the class name of the condition (e.g., "MatchConditionClientList", "MatchConditionHeader")
		Type string `json:"className"`

		// Name specifies the name(s) to match against (e.g., header names, query parameter names)
		Names []string `json:"name,omitempty"`

		// NameWildcard indicates whether wildcard matching is enabled for the Name field
		NameWildcard *bool `json:"nameWildcard,omitempty"`

		// Value specifies the value(s) to match against (e.g., header values, client list IDs)
		Values []string `json:"value"`

		// ValueCase indicates whether value matching is case-sensitive
		ValueCase *bool `json:"valueCase,omitempty"`

		// ValueWildcard indicates whether wildcard matching is enabled for the Value field
		ValueWildcard *bool `json:"valueWildcard,omitempty"`
	}

	// AtomicConditionResp is used to add bypass conditions for URL protection policies.
	// Specify one or more types of conditions to match on. You can match on client lists, request headers, or both.
	AtomicConditionResp struct {

		// CheckIPs specifies the IPs to check (e.g., "REMOTE_ADDR", "X_FORWARDED_FOR")
		CheckIPs *string `json:"checkIps"`

		// Type specifies the class name of the condition (e.g., "MatchConditionClientList", "MatchConditionHeader")
		Type string `json:"className"`

		// Index specifies the index position of this condition in the list
		Index int64 `json:"index"`

		// Name specifies the name(s) to match against (e.g., header names, query parameter names)
		Names []string `json:"name"`

		// NameWildcard indicates whether wildcard matching is enabled for the Name field
		NameWildcard *bool `json:"nameWildcard"`

		// PositiveMatch indicates whether this is a positive match (true) or negative match (false)
		PositiveMatch bool `json:"positiveMatch"`

		// Value specifies the value(s) to match against (e.g., header values, client list IDs)
		Values []string `json:"value"`

		// ValueCase indicates whether value matching is case-sensitive
		ValueCase *bool `json:"valueCase"`

		// ValueWildcard indicates whether wildcard matching is enabled for the Value field
		ValueWildcard *bool `json:"valueWildcard"`
	}

	// BypassCondition is used to define bypass conditions for URL protection policies.
	BypassCondition struct {

		// AtomicConditions lists the individual conditions to evaluate for bypass
		AtomicConditions []AtomicCondition `json:"atomicConditions,omitempty"`
	}

	// BypassConditionResp is used to define bypass conditions for URL protection policies.
	BypassConditionResp struct {

		// AtomicConditions lists the individual conditions to evaluate for bypass
		AtomicConditions []AtomicConditionResp `json:"atomicConditions"`
	}

	// URLProtectionPolicyRequestBody is used to create or update a url protection policy.
	URLProtectionPolicyRequestBody struct {

		// Name is the name of the URL protection policy
		Name string `json:"name"`

		// Description provides details about the URL protection policy
		Description *string `json:"description,omitempty"`

		// BypassCondition contains conditions under which the URL protection policy is bypassed
		BypassCondition *BypassCondition `json:"bypassCondition,omitempty"`

		// MaxRateThreshold specifies the maximum number of requests per second before rate limiting is applied
		MaxRateThreshold int64 `json:"rateThreshold"`

		// APIDefinitions lists the API definitions to which this policy applies
		APIDefinitions []APIDefinition `json:"apiDefinitions,omitempty"`

		// HostnamePaths lists the hostname and path combinations to which this policy applies
		HostnamePaths []HostnamePath `json:"hostnamePaths,omitempty"`

		// ProtectionType specifies the type of protection (e.g., "rate", "slowPost")
		ProtectionType *string `json:"protectionType,omitempty"`

		// IntelligentLoadShedding enables intelligent load shedding when traffic exceeds thresholds
		IntelligentLoadShedding bool `json:"intelligentLoadShedding"`

		// SheddingThresholdHitsPerSec specifies the threshold in hits per second for load shedding
		SheddingThresholdHitsPerSec *int64 `json:"sheddingThresholdHitsPerSec,omitempty"`

		// Categories lists the categories to which this policy applies
		Categories []Category `json:"categories,omitempty"`
	}

	// GetURLProtectionPolicyRequest is used to retrieve information about a specific url protection policy.
	GetURLProtectionPolicyRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64

		// URLProtectionPolicyID is the unique identifier of the URL protection policy to retrieve
		URLProtectionPolicyID int64
	}

	// GetURLProtectionPolicyResponse is returned from a call to GetURLProtectionPolicy.
	GetURLProtectionPolicyResponse struct {

		// URLProtectionPolicyID is the unique identifier for the URL protection policy
		URLProtectionPolicyID int64 `json:"policyId"`

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64 `json:"configId"`

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64 `json:"configVersion"`

		// Name is the name of the URL protection policy
		Name string `json:"name"`

		// Description provides details about the URL protection policy
		Description *string `json:"description"`

		// BypassCondition contains conditions under which the URL protection policy is bypassed
		BypassCondition *BypassConditionResp `json:"bypassCondition"`

		// MaxRateThreshold specifies the maximum number of requests per second before rate limiting is applied
		MaxRateThreshold int64 `json:"rateThreshold"`

		// APIDefinitions lists the API definitions to which this policy applies
		APIDefinitions []APIDefinition `json:"apiDefinitions"`

		// HostnamePaths lists the hostname and path combinations to which this policy applies
		HostnamePaths []HostnamePath `json:"hostnamePaths"`

		// ProtectionType specifies the type of protection (e.g., "rate", "slowPost")
		ProtectionType *string `json:"protectionType"`

		// IntelligentLoadShedding enables intelligent load shedding when traffic exceeds thresholds
		IntelligentLoadShedding bool `json:"intelligentLoadShedding"`

		// SheddingThresholdHitsPerSec specifies the threshold in hits per second for load shedding
		SheddingThresholdHitsPerSec *int64 `json:"sheddingThresholdHitsPerSec"`

		// Categories lists the categories to which this policy applies
		Categories []Category `json:"categories"`

		// Used indicates whether the policy is currently in use
		Used bool `json:"used"`

		// CreateDate is the timestamp when the policy was created
		CreateDate string `json:"createDate"`

		// CreatedBy is the user who created the policy
		CreatedBy string `json:"createdBy"`

		// UpdateDate is the timestamp when the policy was last updated
		UpdateDate string `json:"updateDate"`

		// UpdatedBy is the user who last updated the policy
		UpdatedBy string `json:"updatedBy"`
	}

	// ListURLProtectionPoliciesRequest is used to retrieve the url protection policies for a configuration.
	ListURLProtectionPoliciesRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64
	}

	// ListURLProtectionPoliciesResponse is returned from a call to ListURLProtectionPolicies.
	ListURLProtectionPoliciesResponse struct {

		// URLProtectionPolicies is the list of URL protection policies
		URLProtectionPolicies []GetURLProtectionPolicyResponse `json:"urlProtectionPolicies"`
	}

	// CreateURLProtectionPolicyRequest is used to create a url protection policy.
	CreateURLProtectionPolicyRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64

		// Body contains the details of the URL protection policy to create
		Body URLProtectionPolicyRequestBody
	}

	// CreateURLProtectionPolicyResponse is returned from a call to CreateURLProtectionPolicy.
	CreateURLProtectionPolicyResponse GetURLProtectionPolicyResponse

	// UpdateURLProtectionPolicyRequest is used to modify an existing url protection policy.
	UpdateURLProtectionPolicyRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64

		// URLProtectionPolicyID is the unique identifier for the URL protection policy
		URLProtectionPolicyID int64

		// Body contains the details of the URL protection policy to update
		Body URLProtectionPolicyRequestBody
	}

	// UpdateURLProtectionPolicyResponse is returned from a call to UpdateURLProtectionPolicy.
	UpdateURLProtectionPolicyResponse GetURLProtectionPolicyResponse

	// RemoveURLProtectionPolicyRequest is used to remove a url protection policy.
	RemoveURLProtectionPolicyRequest struct {

		// ConfigID is the unique identifier of the security configuration
		ConfigID int64

		// ConfigVersion is the version number of the security configuration
		ConfigVersion int64

		// URLProtectionPolicyID is the unique identifier of the URL protection policy to remove
		URLProtectionPolicyID int64
	}
)

// Validate validates a GetURLProtectionPolicyRequest.
func (v GetURLProtectionPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":              validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":         validation.Validate(v.ConfigVersion, validation.Required),
		"URLProtectionPolicyID": validation.Validate(v.URLProtectionPolicyID, validation.Required),
	})
}

// Validate validates a ListURLProtectionPoliciesRequest.
func (v ListURLProtectionPoliciesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	})
}

// Validate validates a Create/Update-URLProtectionPolicyRequest.
func (v URLProtectionPolicyRequestBody) Validate() error {
	return validation.Errors{
		"Name":             validation.Validate(v.Name, validation.Required),
		"MaxRateThreshold": validation.Validate(v.MaxRateThreshold, validation.Required),
	}.Filter()
}

// Validate validates a CreateURLProtectionPolicyRequest.
func (v CreateURLProtectionPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Body":          validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates an UpdateURLProtectionPolicyRequest.
func (v UpdateURLProtectionPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":              validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":         validation.Validate(v.ConfigVersion, validation.Required),
		"URLProtectionPolicyID": validation.Validate(v.URLProtectionPolicyID, validation.Required),
		"Body":                  validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a RemoveURLProtectionPolicyRequest.
func (v RemoveURLProtectionPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":              validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":         validation.Validate(v.ConfigVersion, validation.Required),
		"URLProtectionPolicyID": validation.Validate(v.URLProtectionPolicyID, validation.Required),
	})
}

func (p *appsec) GetURLProtectionPolicy(ctx context.Context, params GetURLProtectionPolicyRequest) (*GetURLProtectionPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetURLProtectionPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/url-protections/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.URLProtectionPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetURLProtectionPolicy request: %w", err)
	}

	var result GetURLProtectionPolicyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get url protection policy response failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) ListURLProtectionPolicies(ctx context.Context, params ListURLProtectionPoliciesRequest) (*ListURLProtectionPoliciesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListURLProtectionPolicies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/url-protections",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListURLProtectionPolicies request: %w", err)
	}

	var result ListURLProtectionPoliciesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get url protection policies request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateURLProtectionPolicy(ctx context.Context, params UpdateURLProtectionPolicyRequest) (*UpdateURLProtectionPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateURLProtectionPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/url-protections/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.URLProtectionPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create updateURLProtectionPolicy request: %w", err)
	}

	var result UpdateURLProtectionPolicyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update url protection policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateURLProtectionPolicy(ctx context.Context, params CreateURLProtectionPolicyRequest) (*CreateURLProtectionPolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateURLProtectionPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/url-protections",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateURLProtectionPolicy request: %w", err)
	}

	var result CreateURLProtectionPolicyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("create url protection policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveURLProtectionPolicy(ctx context.Context, params RemoveURLProtectionPolicyRequest) error {
	logger := p.Log(ctx)
	logger.Debug("RemoveURLProtectionPolicy")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/url-protections/%d", params.ConfigID, params.ConfigVersion, params.URLProtectionPolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveURLProtectionPolicy request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("remove url protection policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}
