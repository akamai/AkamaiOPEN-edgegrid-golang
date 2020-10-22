package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// MatchTargetSequence represents a collection of MatchTargetSequence
//
// See: MatchTargetSequence.GetMatchTargetSequence()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// MatchTargetSequence  contains operations available on MatchTargetSequence  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmatchtargetsequence
	MatchTargetSequence interface {
		GetMatchTargetSequences(ctx context.Context, params GetMatchTargetSequencesRequest) (*GetMatchTargetSequencesResponse, error)
		GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error)
		UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error)
	}

	// GetMatchTargetSequence is the argument for GetProperties
	GetMatchTargetSequenceRequest struct {
		ConfigID      int    `json:"configId"`
		ConfigVersion int    `json:"configVersion"`
		Type          string `json:"type"`
	}

	// GetMatchTargetsRequest is the argument for GetProperties
	GetMatchTargetSequencesRequest struct {
		ConfigID      int    `json:"configId"`
		ConfigVersion int    `json:"configVersion"`
		Type          string `json:"type"`
	}

	// GetMatchTargetsSequenceRequest is the argument for GetMatchTargetsSequenceRequest
	GetMatchTargetsSequenceRequest struct {
		ConfigID      int    `json:"configId"`
		ConfigVersion int    `json:"configVersion"`
		Type          string `json:"type"`
	}

	// UpdateMatchTargetRequest is the argument for GetProperties
	UpdateMatchTargetSequenceRequest struct {
		ConfigID       int    `json:"-"`
		ConfigVersion  int    `json:"-"`
		Type           string `json:"type"`
		TargetSequence []struct {
			TargetID int `json:"targetId"`
			Sequence int `json:"sequence"`
		} `json:"targetSequence"`
	}

	// BypassNetworkList ...
	TargetSequence struct {
		TargetID int `json:"targetId"`
		Sequence int `json:"sequence"`
	}

	//GetMatchTargetResponse ...
	GetMatchTargetSequenceResponse struct {
		Type                      string `json:"type"`
		ConfigID                  int    `json:"configId"`
		ConfigVersion             int    `json:"configVersion"`
		DefaultFile               string `json:"defaultFile"`
		EffectiveSecurityControls struct {
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		} `json:"effectiveSecurityControls"`
		Hostnames                    []string `json:"hostnames"`
		IsNegativeFileExtensionMatch bool     `json:"isNegativeFileExtensionMatch"`
		IsNegativePathMatch          bool     `json:"isNegativePathMatch"`
		FilePaths                    []string `json:"filePaths"`
		FileExtensions               []string `json:"fileExtensions"`
		SecurityPolicy               struct {
			PolicyID string `json:"policyId"`
		} `json:"securityPolicy"`
		Sequence           int `json:"sequence"`
		TargetID           int `json:"targetId"`
		BypassNetworkLists []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"bypassNetworkLists"`
	}

	// GetMatchTargetResponse ...
	GetMatchTargetSequencesResponse struct {
		MatchTargets struct {
			APITargets []struct {
				ConfigID                  int    `json:"configId"`
				ConfigVersion             int    `json:"configVersion"`
				Sequence                  int    `json:"sequence"`
				TargetID                  int    `json:"targetId"`
				Type                      string `json:"type"`
				EffectiveSecurityControls struct {
					ApplyAPIConstraints           bool `json:"applyApiConstraints"`
					ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
					ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
					ApplyRateControls             bool `json:"applyRateControls"`
					ApplyReputationControls       bool `json:"applyReputationControls"`
					ApplySlowPostControls         bool `json:"applySlowPostControls"`
				} `json:"effectiveSecurityControls"`
				SecurityPolicy struct {
					PolicyID string `json:"policyId"`
				} `json:"securityPolicy"`
				Apis []struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"apis"`
				BypassNetworkLists []struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"bypassNetworkLists"`
			} `json:"apiTargets"`
			WebsiteTargets []struct {
				ConfigID                     int           `json:"configId"`
				ConfigVersion                int           `json:"configVersion"`
				DefaultFile                  string        `json:"defaultFile"`
				IsNegativeFileExtensionMatch bool          `json:"isNegativeFileExtensionMatch"`
				IsNegativePathMatch          bool          `json:"isNegativePathMatch"`
				Sequence                     int           `json:"sequence"`
				TargetID                     int           `json:"targetId"`
				Type                         string        `json:"type"`
				FileExtensions               []string      `json:"fileExtensions"`
				FilePaths                    []string      `json:"filePaths"`
				Hostnames                    []interface{} `json:"hostnames"`
				EffectiveSecurityControls    struct {
					ApplyAPIConstraints           bool `json:"applyApiConstraints"`
					ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
					ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
					ApplyRateControls             bool `json:"applyRateControls"`
					ApplyReputationControls       bool `json:"applyReputationControls"`
					ApplySlowPostControls         bool `json:"applySlowPostControls"`
				} `json:"effectiveSecurityControls"`
				SecurityPolicy struct {
					PolicyID string `json:"policyId"`
				} `json:"securityPolicy"`
				BypassNetworkLists []struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"bypassNetworkLists"`
			} `json:"websiteTargets"`
		} `json:"matchTargets"`
	}

	// UpdateMatchTargetResponse ...
	UpdateMatchTargetSequenceResponse struct {
		Type           string `json:"type"`
		TargetSequence []struct {
			TargetID int `json:"targetId"`
			Sequence int `json:"sequence"`
		} `json:"targetSequence"`
	}
)

// Validate validates GetMatchTargetSequenceRequest
func (v GetMatchTargetSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates GetMatchTargetSequencesRequest
func (v GetMatchTargetSequencesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

// Validate validates UpdateMatchTargetSequenceRequest
func (v UpdateMatchTargetSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

func (p *appsec) GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetMatchTargetSequence")

	var rval GetMatchTargetSequenceResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getmatchtargetsequence request: %w", err)
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

func (p *appsec) GetMatchTargetSequences(ctx context.Context, params GetMatchTargetSequencesRequest) (*GetMatchTargetSequencesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetMatchTargetSequences")

	var rval GetMatchTargetSequencesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets?includeChildObjectName=true",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getmatchtargetsequences request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getmatchtargetsequences request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a MatchTargetSequence.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmatchtargetsequence

func (p *appsec) UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateMatchTargetSequence")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets/sequence",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create MatchTargetSequencerequest: %w", err)
	}

	var rval UpdateMatchTargetSequenceResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create MatchTargetSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
