package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RatePolicy interface supports creating, retrieving, updating and removing rate policies.
	RatePolicy interface {
		// GetRatePolicies returns rate policies for a specific security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rate-policies
		GetRatePolicies(ctx context.Context, params GetRatePoliciesRequest) (*GetRatePoliciesResponse, error)

		// GetRatePolicy returns the specified rate policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-rate-policy
		GetRatePolicy(ctx context.Context, params GetRatePolicyRequest) (*GetRatePolicyResponse, error)

		// CreateRatePolicy creates a new rate policy for a specific configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-rate-policies
		CreateRatePolicy(ctx context.Context, params CreateRatePolicyRequest) (*CreateRatePolicyResponse, error)

		// UpdateRatePolicy updates details for a specific rate policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-rate-policy
		UpdateRatePolicy(ctx context.Context, params UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error)

		// RemoveRatePolicy deletes the specified rate policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-rate-policy
		RemoveRatePolicy(ctx context.Context, params RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error)
	}

	// CreateRatePolicyRequest is used to create a rate policy.
	CreateRatePolicyRequest struct {
		ID             int             `json:"-"`
		ConfigID       int             `json:"configId"`
		ConfigVersion  int             `json:"configVersion"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// CreateRatePolicyResponse is returned from a call to CreateRatePolicy.
	CreateRatePolicyResponse struct {
		ID                    int      `json:"id"`
		ConfigID              int      `json:"configId"`
		ConfigVersion         int      `json:"configVersion"`
		MatchType             string   `json:"matchType"`
		Type                  string   `json:"type"`
		Name                  string   `json:"name"`
		Description           string   `json:"description"`
		AverageThreshold      int      `json:"averageThreshold"`
		BurstThreshold        int      `json:"burstThreshold"`
		BurstWindow           int      `json:"burstWindow"`
		ClientIdentifiers     []string `json:"clientIdentifiers"`
		UseXForwardForHeaders bool     `json:"useXForwardForHeaders"`
		RequestType           string   `json:"requestType"`
		SameActionOnIpv6      bool     `json:"sameActionOnIpv6"`
		Path                  struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"path"`
		PathMatchType        string `json:"pathMatchType"`
		PathURIPositiveMatch bool   `json:"pathUriPositiveMatch"`
		FileExtensions       struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"fileExtensions"`
		Hosts                  *RatePoliciesHosts      `json:"hosts,omitempty"`
		Hostnames              []string                `json:"hostnames"`
		AdditionalMatchOptions []RatePolicyMatchOption `json:"additionalMatchOptions,omitempty"`
		Condition              *RatePolicyCondition    `json:"condition,omitempty"`
		QueryParameters        []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate         string          `json:"-"`
		UpdateDate         string          `json:"-"`
		Used               json.RawMessage `json:"used"`
		CounterType        string          `json:"counterType"`
		PenaltyBoxDuration string          `json:"penaltyBoxDuration"`
	}

	// UpdateRatePolicyRequest is used to modify an existing rate policy.
	UpdateRatePolicyRequest struct {
		RatePolicyID   int             `json:"id"`
		ConfigID       int             `json:"configId"`
		ConfigVersion  int             `json:"configVersion"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateRatePolicyResponse is returned from a call to UpdateRatePolicy.
	UpdateRatePolicyResponse struct {
		ID                    int      `json:"id"`
		ConfigID              int      `json:"configId"`
		ConfigVersion         int      `json:"configVersion"`
		MatchType             string   `json:"matchType"`
		Type                  string   `json:"type"`
		Name                  string   `json:"name"`
		Description           string   `json:"description"`
		AverageThreshold      int      `json:"averageThreshold"`
		BurstThreshold        int      `json:"burstThreshold"`
		BurstWindow           int      `json:"burstWindow"`
		ClientIdentifiers     []string `json:"clientIdentifiers"`
		UseXForwardForHeaders bool     `json:"useXForwardForHeaders"`
		RequestType           string   `json:"requestType"`
		SameActionOnIpv6      bool     `json:"sameActionOnIpv6"`
		Path                  struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"path"`
		PathMatchType        string `json:"pathMatchType"`
		PathURIPositiveMatch bool   `json:"pathUriPositiveMatch"`
		FileExtensions       struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"fileExtensions"`
		Hosts                  *RatePoliciesHosts      `json:"hosts,omitempty"`
		Hostnames              []string                `json:"hostnames"`
		AdditionalMatchOptions []RatePolicyMatchOption `json:"additionalMatchOptions,omitempty"`
		Condition              *RatePolicyCondition    `json:"condition,omitempty"`
		QueryParameters        []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate         string          `json:"-"`
		UpdateDate         string          `json:"-"`
		Used               json.RawMessage `json:"used"`
		CounterType        string          `json:"counterType"`
		PenaltyBoxDuration string          `json:"penaltyBoxDuration"`
	}

	// RemoveRatePolicyRequest is used to remove a rate policy.
	RemoveRatePolicyRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	// RemoveRatePolicyResponse is returned from a call to RemoveRatePolicy.
	RemoveRatePolicyResponse struct {
		ID                    int      `json:"id"`
		ConfigID              int      `json:"configId"`
		ConfigVersion         int      `json:"configVersion"`
		MatchType             string   `json:"matchType"`
		Type                  string   `json:"type"`
		Name                  string   `json:"name"`
		Description           string   `json:"description"`
		AverageThreshold      int      `json:"averageThreshold"`
		BurstThreshold        int      `json:"burstThreshold"`
		BurstWindow           int      `json:"burstWindow"`
		ClientIdentifiers     []string `json:"clientIdentifiers"`
		UseXForwardForHeaders bool     `json:"useXForwardForHeaders"`
		RequestType           string   `json:"requestType"`
		SameActionOnIpv6      bool     `json:"sameActionOnIpv6"`
		Path                  struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"path"`
		PathMatchType        string `json:"pathMatchType"`
		PathURIPositiveMatch bool   `json:"pathUriPositiveMatch"`
		FileExtensions       struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Values        []string `json:"values"`
		} `json:"fileExtensions"`
		Hosts                  *RatePoliciesHosts      `json:"hosts,omitempty"`
		Hostnames              []string                `json:"hostnames"`
		AdditionalMatchOptions []RatePolicyMatchOption `json:"additionalMatchOptions,omitempty"`
		Condition              *RatePolicyCondition    `json:"condition,omitempty"`
		QueryParameters        []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate         string          `json:"-"`
		UpdateDate         string          `json:"-"`
		Used               json.RawMessage `json:"used"`
		CounterType        string          `json:"counterType"`
		PenaltyBoxDuration string          `json:"penaltyBoxDuration"`
	}

	// GetRatePoliciesRequest is used to retrieve the rate policies for a configuration.
	GetRatePoliciesRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	// GetRatePoliciesResponse is returned from a call to GetRatePolicies.
	GetRatePoliciesResponse struct {
		RatePolicies []struct {
			ID                     int                        `json:"id"`
			ConfigID               int                        `json:"-"`
			ConfigVersion          int                        `json:"-"`
			MatchType              string                     `json:"matchType,omitempty"`
			Type                   string                     `json:"type,omitempty"`
			Name                   string                     `json:"name,omitempty"`
			Description            string                     `json:"description,omitempty"`
			AverageThreshold       int                        `json:"averageThreshold,omitempty"`
			BurstThreshold         int                        `json:"burstThreshold,omitempty"`
			BurstWindow            int                        `json:"burstWindow,omitempty"`
			ClientIdentifiers      []string                   `json:"clientIdentifiers,omitempty"`
			UseXForwardForHeaders  bool                       `json:"useXForwardForHeaders"`
			RequestType            string                     `json:"requestType,omitempty"`
			SameActionOnIpv6       bool                       `json:"sameActionOnIpv6"`
			Path                   *RatePolicyPath            `json:"path,omitempty"`
			PathMatchType          string                     `json:"pathMatchType,omitempty"`
			PathURIPositiveMatch   bool                       `json:"pathUriPositiveMatch"`
			FileExtensions         *RatePolicyFileExtensions  `json:"fileExtensions,omitempty"`
			Hosts                  *RatePoliciesHosts         `json:"hosts,omitempty"`
			Hostnames              []string                   `json:"hostnames,omitempty"`
			AdditionalMatchOptions []RatePolicyMatchOption    `json:"additionalMatchOptions,omitempty"`
			Condition              *RatePolicyCondition       `json:"condition,omitempty"`
			QueryParameters        *RatePolicyQueryParameters `json:"queryParameters,omitempty"`
			CreateDate             string                     `json:"-"`
			UpdateDate             string                     `json:"-"`
			Used                   json.RawMessage            `json:"used"`
			SameActionOnIpv        bool                       `json:"sameActionOnIpv"`
			APISelectors           *RatePolicyAPISelectors    `json:"apiSelectors,omitempty"`
			BodyParameters         *RatePolicyBodyParameters  `json:"bodyParameters,omitempty"`
			CounterType            string                     `json:"counterType"`
			PenaltyBoxDuration     string                     `json:"penaltyBoxDuration"`
		} `json:"ratePolicies,omitempty"`
	}

	// GetRatePolicyRequest is used to retrieve information about a specific rate policy.
	GetRatePolicyRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	// GetRatePolicyResponse is returned from a call to GetRatePolicy.
	GetRatePolicyResponse struct {
		ID                     int                        `json:"-"`
		ConfigID               int                        `json:"-"`
		ConfigVersion          int                        `json:"-"`
		MatchType              string                     `json:"matchType,omitempty"`
		Type                   string                     `json:"type,omitempty"`
		Name                   string                     `json:"name,omitempty"`
		Description            string                     `json:"description,omitempty"`
		AverageThreshold       int                        `json:"averageThreshold,omitempty"`
		BurstThreshold         int                        `json:"burstThreshold,omitempty"`
		BurstWindow            int                        `json:"burstWindow,omitempty"`
		ClientIdentifiers      []string                   `json:"clientIdentifiers,omitempty"`
		UseXForwardForHeaders  bool                       `json:"useXForwardForHeaders"`
		RequestType            string                     `json:"requestType,omitempty"`
		SameActionOnIpv6       bool                       `json:"sameActionOnIpv6"`
		Path                   *RatePolicyPath            `json:"path,omitempty"`
		PathMatchType          string                     `json:"pathMatchType,omitempty"`
		PathURIPositiveMatch   bool                       `json:"pathUriPositiveMatch"`
		FileExtensions         *RatePolicyFileExtensions  `json:"fileExtensions,omitempty"`
		Hosts                  *RatePoliciesHosts         `json:"hosts,omitempty"`
		Hostnames              []string                   `json:"hostnames,omitempty"`
		AdditionalMatchOptions []RatePolicyMatchOption    `json:"additionalMatchOptions,omitempty"`
		Condition              *RatePolicyCondition       `json:"condition,omitempty"`
		QueryParameters        *RatePolicyQueryParameters `json:"queryParameters,omitempty"`
		CreateDate             string                     `json:"-"`
		UpdateDate             string                     `json:"-"`
		Used                   bool                       `json:"-"`
		CounterType            string                     `json:"counterType"`
		PenaltyBoxDuration     string                     `json:"penaltyBoxDuration"`
	}

	// RatePolicyAPISelectors is used as part of a rate policy description.
	RatePolicyAPISelectors []struct {
		APIDefinitionID    int   `json:"apiDefinitionId,omitempty"`
		DefinedResources   *bool `json:"definedResources,omitempty"`
		ResourceIDs        []int `json:"resourceIds"`
		UndefinedResources *bool `json:"undefinedResources,omitempty"`
	}

	// RatePolicyBodyParameters is used as part of a rate policy description.
	RatePolicyBodyParameters []struct {
		Name          string   `json:"name,omitempty"`
		Values        []string `json:"values,omitempty"`
		PositiveMatch bool     `json:"positiveMatch"`
		ValueInRange  bool     `json:"valueInRange"`
	}

	// RatePolicyPath is used as part of a rate policy description.
	RatePolicyPath struct {
		PositiveMatch bool     `json:"positiveMatch"`
		Values        []string `json:"values,omitempty"`
	}

	// RatePolicyFileExtensions is used as part of a rate policy description.
	RatePolicyFileExtensions struct {
		PositiveMatch bool     `json:"positiveMatch"`
		Values        []string `json:"values,omitempty"`
	}

	// RatePolicyMatchOption is used as part of a rate policy description.
	RatePolicyMatchOption struct {
		PositiveMatch bool     `json:"positiveMatch"`
		Type          string   `json:"type,omitempty"`
		Values        []string `json:"values,omitempty"`
	}

	// RatePolicyQueryParameters is used as part of a rate policy description.
	RatePolicyQueryParameters []struct {
		Name          string   `json:"name,omitempty"`
		Values        []string `json:"values,omitempty"`
		PositiveMatch bool     `json:"positiveMatch"`
		ValueInRange  bool     `json:"valueInRange"`
	}

	// RatePoliciesHosts is used as part of a rate policy description.
	RatePoliciesHosts struct {
		Values        *[]string        `json:"values,omitempty"`
		PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
	}

	// RatePolicyCondition is used as part of a rate policy description.
	RatePolicyCondition struct {
		AtomicConditions []struct {
			Value            *json.RawMessage `json:"value,omitempty"`
			ClassName        string           `json:"className"`
			PositiveMatch    bool             `json:"positiveMatch"`
			Name             []string         `json:"name,omitempty"`
			NameCase         bool             `json:"nameCase,omitempty"`
			NameWildcard     bool             `json:"nameWildcard,omitempty"`
			ValueCase        bool             `json:"valueCase,omitempty"`
			ValueWildcard    bool             `json:"valueWildcard,omitempty"`
			SharedIpHandling string           `json:"sharedIpHandling,omitempty"`
		} `json:"atomicConditions,omitempty"`
	}
)

// Validate validates a GetRatePolicyRequest.
func (v GetRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetRatePoliciesRequest.
func (v GetRatePoliciesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates a CreateRatePolicyRequest.
func (v CreateRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRatePolicyRequest.
func (v UpdateRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveRatePolicyRequest.
func (v RemoveRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRatePolicy(ctx context.Context, params GetRatePolicyRequest) (*GetRatePolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicy request: %w", err)
	}

	var result GetRatePolicyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rate policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetRatePolicies(ctx context.Context, params GetRatePoliciesRequest) (*GetRatePoliciesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetRatePolicies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicies request: %w", err)
	}

	var result GetRatePoliciesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get rate policies request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RatePolicyID != 0 {
		var filteredResult GetRatePoliciesResponse
		for _, val := range result.RatePolicies {
			if val.ID == params.RatePolicyID {
				filteredResult.RatePolicies = append(filteredResult.RatePolicies, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateRatePolicy(ctx context.Context, params UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateRatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRatePolicy request: %w", err)
	}

	var result UpdateRatePolicyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update rate policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateRatePolicy(ctx context.Context, params CreateRatePolicyRequest) (*CreateRatePolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateRatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateRatePolicy request: %w", err)
	}

	var result CreateRatePolicyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create rate policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveRatePolicy(ctx context.Context, params RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveRatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/rate-policies/%d", params.ConfigID, params.ConfigVersion, params.RatePolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveRatePolicy request: %w", err)
	}

	var result RemoveRatePolicyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("remove rate policy request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
