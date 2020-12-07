package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RatePolicy represents a collection of RatePolicy
//
// See: RatePolicy.GetRatePolicy()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RatePolicy  contains operations available on RatePolicy  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getratepolicy
	RatePolicy interface {
		GetRatePolicies(ctx context.Context, params GetRatePoliciesRequest) (*GetRatePoliciesResponse, error)
		GetRatePolicy(ctx context.Context, params GetRatePolicyRequest) (*GetRatePolicyResponse, error)
		CreateRatePolicy(ctx context.Context, params CreateRatePolicyRequest) (*CreateRatePolicyResponse, error)
		UpdateRatePolicy(ctx context.Context, params UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error)
		RemoveRatePolicy(ctx context.Context, params RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error)
	}

	GetRatePolicyResponse struct {
		ID                    int    `json:"id,omitempty"`
		PolicyID              int    `json:"policyId,omitempty"`
		ConfigID              int    `json:"-"`
		ConfigVersion         int    `json:"-"`
		MatchType             string `json:"matchType,omitempty"`
		Type                  string `json:"type,omitempty"`
		Name                  string `json:"name,omitempty"`
		Description           string `json:"description,omitempty"`
		AverageThreshold      int    `json:"averageThreshold,omitempty"`
		BurstThreshold        int    `json:"burstThreshold,omitempty"`
		ClientIdentifier      string `json:"clientIdentifier,omitempty"`
		UseXForwardForHeaders bool   `json:"useXForwardForHeaders,omitempty"`
		RequestType           string `json:"requestType,omitempty"`
		SameActionOnIpv6      bool   `json:"sameActionOnIpv6,omitempty"`
		Path                  struct {
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			Values        []string `json:"values,omitempty"`
		} `json:"path,omitempty"`
		PathMatchType        string `json:"pathMatchType,omitempty"`
		PathURIPositiveMatch bool   `json:"pathUriPositiveMatch,omitempty"`
		FileExtensions       struct {
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			Values        []string `json:"values,omitempty"`
		} `json:"fileExtensions,omitempty"`
		Hostnames              []string `json:"hostNames,omitempty"`
		AdditionalMatchOptions []struct {
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			Type          string   `json:"type,omitempty"`
			Values        []string `json:"values,omitempty"`
		} `json:"additionalMatchOptions,omitempty"`
		QueryParameters []struct {
			Name          string   `json:"name,omitempty"`
			Values        []string `json:"values,omitempty"`
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			ValueInRange  bool     `json:"valueInRange,omitempty"`
		} `json:"queryParameters,omitempty"`
		CreateDate string `json:"createDate,omitempty"`
		UpdateDate string `json:"updateDate,omitempty"`
		Used       bool   `json:"used,omitempty"`
	}

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
		Hostnames              []string `json:"hostNames"`
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
		CreateDate string `json:"createDate"`
		UpdateDate string `json:"updateDate"`
		Used       bool   `json:"used"`
	}

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
		Hostnames              []string `json:"hostNames"`
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
		CreateDate string `json:"createDate"`
		UpdateDate string `json:"updateDate"`
		Used       bool   `json:"used"`
	}

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
		Hostnames              []string `json:"hostNames"`
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
		CreateDate string `json:"createDate"`
		UpdateDate string `json:"updateDate"`
		Used       bool   `json:"used"`
	}

	GetRatePoliciesRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	GetRatePolicyRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	CreateRatePolicyRequest struct {
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
		Hostnames              []string `json:"hostNames"`
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
		CreateDate string `json:"createDate"`
		UpdateDate string `json:"updateDate"`
		Used       bool   `json:"used"`
	}

	UpdateRatePolicyRequest struct {
		RatePolicyID          int    `json:"id"`
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
		Hostnames              []string `json:"hostNames"`
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
		CreateDate string `json:"createDate"`
		UpdateDate string `json:"updateDate"`
		Used       bool   `json:"used"`
	}

	RemoveRatePolicyRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
		RatePolicyID  int `json:"ratePolicyId"`
	}

	GetRatePoliciesResponse struct {
		RatePolicies []struct {
			ID                    int    `json:"id"`
			ConfigID              int    `json:"-"`
			ConfigVersion         int    `json:"-"`
			MatchType             string `json:"matchType,omitempty"`
			Type                  string `json:"type,omitempty"`
			Name                  string `json:"name,omitempty"`
			Description           string `json:"description,omitempty"`
			AverageThreshold      int    `json:"averageThreshold,omitempty"`
			BurstThreshold        int    `json:"burstThreshold,omitempty"`
			ClientIdentifier      string `json:"clientIdentifier,omitempty"`
			UseXForwardForHeaders bool   `json:"useXForwardForHeaders,omitempty"`
			RequestType           string `json:"requestType,omitempty"`
			SameActionOnIpv6      bool   `json:"sameActionOnIpv6,omitempty"`
			Path                  struct {
				PositiveMatch bool     `json:"positiveMatch,omitempty"`
				Values        []string `json:"values,omitempty"`
			} `json:"path,omitempty"`
			PathMatchType        string `json:"pathMatchType,omitempty"`
			PathURIPositiveMatch bool   `json:"pathUriPositiveMatch,omitempty"`
			FileExtensions       struct {
				PositiveMatch bool     `json:"positiveMatch,omitempty"`
				Values        []string `json:"values,omitempty"`
			} `json:"fileExtensions,omitempty"`
			Hostnames              []string `json:"hostnames,omitempty"`
			AdditionalMatchOptions []struct {
				PositiveMatch bool     `json:"positiveMatch,omitempty"`
				Type          string   `json:"type,omitempty"`
				Values        []string `json:"values,omitempty"`
			} `json:"additionalMatchOptions,omitempty"`
			QueryParameters []struct {
				Name          string   `json:"name,omitempty"`
				Values        []string `json:"values,omitempty"`
				PositiveMatch bool     `json:"positiveMatch,omitempty"`
				ValueInRange  bool     `json:"valueInRange,omitempty"`
			} `json:"queryParameters,omitempty"`
			CreateDate      string `json:"createDate,omitempty"`
			UpdateDate      string `json:"updateDate,omitempty"`
			EnableActions   bool   `json:"enableActions,omitempty"`
			Used            bool   `json:"used,omitempty"`
			SameActionOnIpv bool   `json:"sameActionOnIpv,omitempty"`
			APISelectors    []struct {
				APIDefinitionID int   `json:"apiDefinitionId,omitempty"`
				ResourceIds     []int `json:"resourceIds,omitempty"`
			} `json:"apiSelectors,omitempty"`
			BodyParameters []struct {
				Name          string   `json:"name,omitempty"`
				Values        []string `json:"values,omitempty"`
				PositiveMatch bool     `json:"positiveMatch,omitempty"`
				ValueInRange  bool     `json:"valueInRange,omitempty"`
			} `json:"bodyParameters,omitempty"`
		} `json:"ratePolicies,omitempty"`
	}
)

// Validate validates GetRatePolicyRequest
func (v GetRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetRatePolicysRequest
func (v GetRatePoliciesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates CreateRatePolicyRequest
func (v CreateRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates UpdateRatePolicyRequest
func (v UpdateRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

// Validate validates RemoveRatePolicyRequest
func (v RemoveRatePolicyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.RatePolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRatePolicy(ctx context.Context, params GetRatePolicyRequest) (*GetRatePolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRatePolicy")

	var rval GetRatePolicyResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getratepolicy request: %w", err)
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

func (p *appsec) GetRatePolicies(ctx context.Context, params GetRatePoliciesRequest) (*GetRatePoliciesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRatePolicys")

	var rval GetRatePoliciesResponse
	var rvalfiltered GetRatePoliciesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getratepolicies request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getratepolicies request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RatePolicyID != 0 {
		for _, val := range rval.RatePolicies {
			if val.ID == params.RatePolicyID {
				rvalfiltered.RatePolicies = append(rvalfiltered.RatePolicies, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a RatePolicy.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putratepolicy

func (p *appsec) UpdateRatePolicy(ctx context.Context, params UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRatePolicy")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create RatePolicyrequest: %w", err)
	}

	var rval UpdateRatePolicyResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RatePolicy request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Create will create a new ratepolicy.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postratepolicy
func (p *appsec) CreateRatePolicy(ctx context.Context, params CreateRatePolicyRequest) (*CreateRatePolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateRatePolicy")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ratepolicy request: %w", err)
	}

	var rval CreateRatePolicyResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create ratepolicyrequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a RatePolicy
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteratepolicy

func (p *appsec) RemoveRatePolicy(ctx context.Context, params RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveRatePolicyResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveRatePolicy")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/rate-policies/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.RatePolicyID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delratepolicy request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delratepolicy request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
