package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RatePolicy interface supports creating, retrieving, updating and removing rate policies.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#ratepolicy
	RatePolicy interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicies
		GetRatePolicies(ctx context.Context, params GetRatePoliciesRequest) (*GetRatePoliciesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicy
		GetRatePolicy(ctx context.Context, params GetRatePolicyRequest) (*GetRatePolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postratepolicies
		CreateRatePolicy(ctx context.Context, params CreateRatePolicyRequest) (*CreateRatePolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putratepolicy
		UpdateRatePolicy(ctx context.Context, params UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteratepolicy
		RemoveRatePolicy(ctx context.Context, params RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error)
	}

	// CreateRatePolicyRequest is used to create a rate policy.
	CreateRatePolicyRequest struct {
		ID             int             `json:"-"`
		PolicyID       int             `json:"-"`
		ConfigID       int             `json:"configId"`
		ConfigVersion  int             `json:"configVersion"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// CreateRatePolicyResponse is returned from a call to CreateRatePolicy.
	CreateRatePolicyResponse struct {
		ID                    int    `json:"id"`
		PolicyID              int    `json:"policyId"`
		ConfigID              int    `json:"configId"`
		ConfigVersion         int    `json:"configVersion"`
		MatchType             string `json:"matchType"`
		Type                  string `json:"type"`
		Name                  string `json:"name"`
		Description           string `json:"description"`
		AverageThreshold      int    `json:"averageThreshold"`
		BurstThreshold        int    `json:"burstThreshold"`
		ClientIdentifier      string `json:"clientIdentifier"`
		UseXForwardForHeaders bool   `json:"useXForwardForHeaders"`
		RequestType           string `json:"requestType"`
		SameActionOnIpv6      bool   `json:"sameActionOnIpv6"`
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
		Hosts                  *RatePoliciesHosts                 `json:"hosts,omitempty"`
		Hostnames              []string                           `json:"hostnames"`
		AdditionalMatchOptions []RatePolicyAdditionalMatchOptions `json:"additionalMatchOptions"`
		QueryParameters        []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate string          `json:"-"`
		UpdateDate string          `json:"-"`
		Used       json.RawMessage `json:"used"`
	}

	// UpdateRatePolicyRequest is used to modify an existing rate policy.
	UpdateRatePolicyRequest struct {
		RatePolicyID   int             `json:"id"`
		PolicyID       int             `json:"policyId"`
		ConfigID       int             `json:"configId"`
		ConfigVersion  int             `json:"configVersion"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateRatePolicyResponse is returned from a call to UpdateRatePolicy.
	UpdateRatePolicyResponse struct {
		ID                    int    `json:"id"`
		PolicyID              int    `json:"policyId"`
		ConfigID              int    `json:"configId"`
		ConfigVersion         int    `json:"configVersion"`
		MatchType             string `json:"matchType"`
		Type                  string `json:"type"`
		Name                  string `json:"name"`
		Description           string `json:"description"`
		AverageThreshold      int    `json:"averageThreshold"`
		BurstThreshold        int    `json:"burstThreshold"`
		ClientIdentifier      string `json:"clientIdentifier"`
		UseXForwardForHeaders bool   `json:"useXForwardForHeaders"`
		RequestType           string `json:"requestType"`
		SameActionOnIpv6      bool   `json:"sameActionOnIpv6"`
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
		Hosts                  *RatePoliciesHosts `json:"hosts,omitempty"`
		Hostnames              []string           `json:"hostnames"`
		AdditionalMatchOptions []struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Type          string   `json:"type"`
			Values        []string `json:"values"`
		} `json:"additionalMatchOptions"`
		QueryParameters []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate string          `json:"-"`
		UpdateDate string          `json:"-"`
		Used       json.RawMessage `json:"used"`
	}

	// RemoveRatePolicyRequest is used to remove a rate policy.
	RemoveRatePolicyRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	// RemoveRatePolicyResponse is returned from a call to RemoveRatePolicy.
	RemoveRatePolicyResponse struct {
		ID                    int    `json:"id"`
		PolicyID              int    `json:"policyId"`
		ConfigID              int    `json:"configId"`
		ConfigVersion         int    `json:"configVersion"`
		MatchType             string `json:"matchType"`
		Type                  string `json:"type"`
		Name                  string `json:"name"`
		Description           string `json:"description"`
		AverageThreshold      int    `json:"averageThreshold"`
		BurstThreshold        int    `json:"burstThreshold"`
		ClientIdentifier      string `json:"clientIdentifier"`
		UseXForwardForHeaders bool   `json:"useXForwardForHeaders"`
		RequestType           string `json:"requestType"`
		SameActionOnIpv6      bool   `json:"sameActionOnIpv6"`
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
		Hosts                  *RatePoliciesHosts `json:"hosts,omitempty"`
		Hostnames              []string           `json:"hostnames"`
		AdditionalMatchOptions []struct {
			PositiveMatch bool     `json:"positiveMatch"`
			Type          string   `json:"type"`
			Values        []string `json:"values"`
		} `json:"additionalMatchOptions"`
		QueryParameters []struct {
			Name          string   `json:"name"`
			Values        []string `json:"values"`
			PositiveMatch bool     `json:"positiveMatch"`
			ValueInRange  bool     `json:"valueInRange"`
		} `json:"queryParameters"`
		CreateDate string          `json:"-"`
		UpdateDate string          `json:"-"`
		Used       json.RawMessage `json:"used"`
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
			ID                     int                               `json:"id"`
			ConfigID               int                               `json:"-"`
			ConfigVersion          int                               `json:"-"`
			MatchType              string                            `json:"matchType,omitempty"`
			Type                   string                            `json:"type,omitempty"`
			Name                   string                            `json:"name,omitempty"`
			Description            string                            `json:"description,omitempty"`
			AverageThreshold       int                               `json:"averageThreshold,omitempty"`
			BurstThreshold         int                               `json:"burstThreshold,omitempty"`
			ClientIdentifier       string                            `json:"clientIdentifier,omitempty"`
			UseXForwardForHeaders  bool                              `json:"useXForwardForHeaders"`
			RequestType            string                            `json:"requestType,omitempty"`
			SameActionOnIpv6       bool                              `json:"sameActionOnIpv6"`
			Path                   *RatePolicyPath                   `json:"path,omitempty"`
			PathMatchType          string                            `json:"pathMatchType,omitempty"`
			PathURIPositiveMatch   bool                              `json:"pathUriPositiveMatch"`
			FileExtensions         *RatePolicyFileExtensions         `json:"fileExtensions,omitempty"`
			Hosts                  *RatePoliciesHosts                `json:"hosts,omitempty"`
			Hostnames              []string                          `json:"hostnames,omitempty"`
			AdditionalMatchOptions *RatePolicyAdditionalMatchOptions `json:"additionalMatchOptions,omitempty"`
			QueryParameters        *RatePolicyQueryParameters        `json:"queryParameters,omitempty"`
			CreateDate             string                            `json:"-"`
			UpdateDate             string                            `json:"-"`
			Used                   json.RawMessage                   `json:"used"`
			SameActionOnIpv        bool                              `json:"sameActionOnIpv"`
			APISelectors           *RatePolicyAPISelectors           `json:"apiSelectors,omitempty"`
			BodyParameters         *RatePolicyBodyParameters         `json:"bodyParameters,omitempty"`
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
		ID                     int                               `json:"-"`
		PolicyID               int                               `json:"policyId,omitempty"`
		ConfigID               int                               `json:"-"`
		ConfigVersion          int                               `json:"-"`
		MatchType              string                            `json:"matchType,omitempty"`
		Type                   string                            `json:"type,omitempty"`
		Name                   string                            `json:"name,omitempty"`
		Description            string                            `json:"description,omitempty"`
		AverageThreshold       int                               `json:"averageThreshold,omitempty"`
		BurstThreshold         int                               `json:"burstThreshold,omitempty"`
		ClientIdentifier       string                            `json:"clientIdentifier,omitempty"`
		UseXForwardForHeaders  bool                              `json:"useXForwardForHeaders"`
		RequestType            string                            `json:"requestType,omitempty"`
		SameActionOnIpv6       bool                              `json:"sameActionOnIpv6"`
		Path                   *RatePolicyPath                   `json:"path,omitempty"`
		PathMatchType          string                            `json:"pathMatchType,omitempty"`
		PathURIPositiveMatch   bool                              `json:"pathUriPositiveMatch"`
		FileExtensions         *RatePolicyFileExtensions         `json:"fileExtensions,omitempty"`
		Hosts                  *RatePoliciesHosts                `json:"hosts,omitempty"`
		Hostnames              []string                          `json:"hostnames,omitempty"`
		AdditionalMatchOptions *RatePolicyAdditionalMatchOptions `json:"additionalMatchOptions,omitempty"`
		QueryParameters        *RatePolicyQueryParameters        `json:"queryParameters,omitempty"`
		CreateDate             string                            `json:"-"`
		UpdateDate             string                            `json:"-"`
		Used                   bool                              `json:"-"`
	}

	// RatePolicyAPISelectors is used as part of a rate policy description.
	RatePolicyAPISelectors []struct {
		APIDefinitionID    int   `json:"apiDefinitionId,omitempty"`
		DefinedResources   *bool `json:"definedResources,omitempty"`
		ResourceIds        []int `json:"resourceIds"`
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

	// RatePolicyAdditionalMatchOptions is used as part of a rate policy description.
	RatePolicyAdditionalMatchOptions []struct {
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
)

// Validate validates a GetRatePolicyRequest.
func (v GetRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetRatePolicysRequest.
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

	var result GetRatePolicyResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicy request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRatePolicy request failed: %w", err)
	}

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

	var result GetRatePoliciesResponse
	var filteredResult GetRatePoliciesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRatePolicies request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRatePolicies request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RatePolicyID != 0 {
		for _, val := range result.RatePolicies {
			if val.ID == params.RatePolicyID {
				filteredResult.RatePolicies = append(filteredResult.RatePolicies, val)
			}
		}

	} else {
		filteredResult = result
	}

	return &filteredResult, nil

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
		return nil, fmt.Errorf("UpdateRatePolicy request failed: %w", err)
	}

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
		return nil, fmt.Errorf("CreateRatePolicy request failed: %w", err)
	}

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

	var result RemoveRatePolicyResponse

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/rate-policies/%d", params.ConfigID, params.ConfigVersion, params.RatePolicyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveRatePolicy request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("RemoveRatePolicy request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
