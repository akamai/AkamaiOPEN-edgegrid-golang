package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ReputationProfile interface supports creating, retrieving, modifying and removing reputation
	// profiles for a specific security configuration version.
	ReputationProfile interface {
		// GetReputationProfiles returns reputation profiles for a specific security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-reputation-profiles
		GetReputationProfiles(ctx context.Context, params GetReputationProfilesRequest) (*GetReputationProfilesResponse, error)

		// GetReputationProfile returns the details for a specific reputation profile.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-reputation-profile
		GetReputationProfile(ctx context.Context, params GetReputationProfileRequest) (*GetReputationProfileResponse, error)

		// CreateReputationProfile creates a new reputation profile for a specific configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-reputation-profiles
		CreateReputationProfile(ctx context.Context, params CreateReputationProfileRequest) (*CreateReputationProfileResponse, error)

		// UpdateReputationProfile updates details for a specific reputation profile.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-reputation-profile
		UpdateReputationProfile(ctx context.Context, params UpdateReputationProfileRequest) (*UpdateReputationProfileResponse, error)

		// RemoveReputationProfile deletes a reputation profile.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-reputation-profile
		RemoveReputationProfile(ctx context.Context, params RemoveReputationProfileRequest) (*RemoveReputationProfileResponse, error)
	}

	atomicConditionsName []string

	// GetReputationProfilesRequest is used to retrieve the reputation profiles for a configuration.
	GetReputationProfilesRequest struct {
		ConfigID            int `json:"configId"`
		ConfigVersion       int `json:"configVersion"`
		ReputationProfileId int `json:"-"`
	}

	// GetReputationProfilesResponse is returned from a call to GetReputationProfiles.
	GetReputationProfilesResponse struct {
		ReputationProfiles []struct {
			Condition        *ReputationProfileCondition `json:"condition,omitempty"`
			Context          string                      `json:"context,omitempty"`
			ContextReadable  string                      `json:"-"`
			Enabled          bool                        `json:"-"`
			ID               int                         `json:"id,omitempty"`
			Name             string                      `json:"name,omitempty"`
			SharedIPHandling string                      `json:"sharedIpHandling,omitempty"`
			Threshold        int                         `json:"threshold,omitempty"`
		} `json:"reputationProfiles,omitempty"`
	}

	// GetReputationProfileRequest is used to retrieve the details for a specific reputation profile.
	GetReputationProfileRequest struct {
		ConfigID            int `json:"configId"`
		ConfigVersion       int `json:"configVersion"`
		ReputationProfileId int `json:"-"`
	}

	// GetReputationProfileResponse is returned from a call to GetReputationProfile.
	GetReputationProfileResponse struct {
		Condition        *GetReputationProfileResponseCondition `json:"condition,omitempty"`
		Context          string                                 `json:"context,omitempty"`
		ContextReadable  string                                 `json:"-"`
		Enabled          bool                                   `json:"-"`
		ID               int                                    `json:"-"`
		Name             string                                 `json:"name,omitempty"`
		SharedIPHandling string                                 `json:"sharedIpHandling,omitempty"`
		Threshold        int                                    `json:"threshold,omitempty"`
	}

	// CreateReputationProfileRequest is used to create a reputation profile.
	CreateReputationProfileRequest struct {
		ConfigID       int             `json:"-"`
		ConfigVersion  int             `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// CreateReputationProfileResponse is returned from a call to CreateReputationProfile.
	CreateReputationProfileResponse struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Context          string `json:"context"`
		Description      string `json:"description"`
		Threshold        int    `json:"threshold"`
		SharedIPHandling string `json:"sharedIpHandling"`
		Condition        struct {
			AtomicConditions []struct {
				CheckIps      string               `json:"checkIps,omitempty"`
				ClassName     string               `json:"className"`
				Index         int                  `json:"index"`
				PositiveMatch bool                 `json:"positiveMatch"`
				Value         []string             `json:"value,omitempty"`
				Name          atomicConditionsName `json:"name,omitempty"`
				NameCase      bool                 `json:"nameCase,omitempty"`
				NameWildcard  bool                 `json:"nameWildcard,omitempty"`
				ValueCase     bool                 `json:"valueCase,omitempty"`
				ValueWildcard bool                 `json:"valueWildcard,omitempty"`
				Host          []string             `json:"host,omitempty"`
			} `json:"atomicConditions"`
			PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
		} `json:"condition"`
		Enabled bool `json:"enabled"`
	}

	// UpdateReputationProfileRequest is used to modify an existing reputation profile.
	UpdateReputationProfileRequest struct {
		ConfigID            int             `json:"-"`
		ConfigVersion       int             `json:"-"`
		ReputationProfileId int             `json:"-"`
		JsonPayloadRaw      json.RawMessage `json:"-"`
	}

	// UpdateReputationProfileResponse is returned from a call to UpdateReputationProfile.
	UpdateReputationProfileResponse struct {
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

	// RemoveReputationProfileRequest is used to remove a reputation profile.
	RemoveReputationProfileRequest struct {
		ConfigID            int `json:"configId"`
		ConfigVersion       int `json:"configVersion"`
		ReputationProfileId int `json:"-"`
	}

	// RemoveReputationProfileResponse is returned from a call to RemoveReputationProfile.
	RemoveReputationProfileResponse struct {
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

	// ReputationProfileCondition is used as part of a reputation profile description.
	ReputationProfileCondition struct {
		AtomicConditions []struct {
			CheckIps      *json.RawMessage `json:"checkIps,omitempty"`
			ClassName     string           `json:"className,omitempty"`
			Index         int              `json:"index,omitempty"`
			PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
			Value         []string         `json:"value,omitempty"`
			Name          *json.RawMessage `json:"name,omitempty"`
			NameCase      bool             `json:"nameCase,omitempty"`
			NameWildcard  *json.RawMessage `json:"nameWildcard,omitempty"`
			ValueCase     bool             `json:"valueCase,omitempty"`
			ValueWildcard *json.RawMessage `json:"valueWildcard,omitempty"`
			Host          []string         `json:"host,omitempty"`
		} `json:"atomicConditions,omitempty"`
		PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
	}

	// GetReputationProfileResponseCondition is used as part of the response to GetReputationProfile.
	GetReputationProfileResponseCondition struct {
		AtomicConditions []struct {
			CheckIps      *json.RawMessage `json:"checkIps,omitempty"`
			ClassName     string           `json:"className,omitempty"`
			Index         int              `json:"index,omitempty"`
			PositiveMatch json.RawMessage  `json:"positiveMatch,omitempty"`
			Value         []string         `json:"value,omitempty"`
			Name          *json.RawMessage `json:"name,omitempty"`
			NameCase      bool             `json:"nameCase,omitempty"`
			NameWildcard  *json.RawMessage `json:"nameWildcard,omitempty"`
			ValueCase     bool             `json:"valueCase,omitempty"`
			ValueWildcard *json.RawMessage `json:"valueWildcard,omitempty"`
			Host          []string         `json:"host,omitempty"`
		} `json:"atomicConditions,omitempty"`
		PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
	}
)

func (c *atomicConditionsName) UnmarshalJSON(data []byte) error {
	var nums interface{}
	err := json.Unmarshal(data, &nums)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(nums)
	switch items.Kind() {
	case reflect.String:
		*c = append(*c, items.String())

	case reflect.Slice:
		*c = make(atomicConditionsName, 0, items.Len())
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				*c = append(*c, item.String())
			case reflect.Interface:
				*c = append(*c, item.Interface().(string))
			}
		}
	}
	return nil
}

// Validate validates a GetReputationProfileRequest.
func (v GetReputationProfileRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"RatePolicyID":  validation.Validate(v.ReputationProfileId, validation.Required),
	}.Filter()
}

// Validate validates a GetReputationProfilesRequest.
func (v GetReputationProfilesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates a CreateReputationProfileRequest.
func (v CreateReputationProfileRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates an UpdateReputationProfileRequest.
func (v UpdateReputationProfileRequest) Validate() error {
	return validation.Errors{
		"ConfigID":            validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":       validation.Validate(v.ConfigVersion, validation.Required),
		"ReputationProfileId": validation.Validate(v.ReputationProfileId, validation.Required),
	}.Filter()
}

// Validate validates a RemoveReputationProfileRequest.
func (v RemoveReputationProfileRequest) Validate() error {
	return validation.Errors{
		"ConfigID":            validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion":       validation.Validate(v.ConfigVersion, validation.Required),
		"ReputationProfileId": validation.Validate(v.ReputationProfileId, validation.Required),
	}.Filter()
}

func (p *appsec) GetReputationProfile(ctx context.Context, params GetReputationProfileRequest) (*GetReputationProfileResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetReputationProfile")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/reputation-profiles/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.ReputationProfileId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetReputationProfile request: %w", err)
	}

	var result GetReputationProfileResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation profile request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetReputationProfiles(ctx context.Context, params GetReputationProfilesRequest) (*GetReputationProfilesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetReputationProfiles")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/reputation-profiles",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetReputationProfiles request: %w", err)
	}

	var result GetReputationProfilesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get reputation profiles request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ReputationProfileId != 0 {
		var filteredResult GetReputationProfilesResponse
		for _, val := range result.ReputationProfiles {
			if val.ID == params.ReputationProfileId {
				filteredResult.ReputationProfiles = append(filteredResult.ReputationProfiles, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateReputationProfile(ctx context.Context, params UpdateReputationProfileRequest) (*UpdateReputationProfileResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateReputationProfile")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/reputation-profiles/%d",
		params.ConfigID,
		params.ConfigVersion,
		params.ReputationProfileId,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateReputationProfile request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result UpdateReputationProfileResponse
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update reputation profile request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateReputationProfile(ctx context.Context, params CreateReputationProfileRequest) (*CreateReputationProfileResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateReputationProfile")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/reputation-profiles",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateReputationProfile request: %w", err)
	}

	var result CreateReputationProfileResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create reputation profile request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveReputationProfile(ctx context.Context, params RemoveReputationProfileRequest) (*RemoveReputationProfileResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveReputationProfile")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/reputation-profiles/%d", params.ConfigID, params.ConfigVersion, params.ReputationProfileId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveReputationProfile request: %w", err)
	}

	var result RemoveReputationProfileResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("remove reputation profile request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
