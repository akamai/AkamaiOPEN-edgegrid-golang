package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ApiHostnameCoverageMatchTargets represents a collection of ApiHostnameCoverageMatchTargets
//
// See: ApiHostnameCoverageMatchTargets.GetApiHostnameCoverageMatchTargets()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiHostnameCoverageMatchTargets  contains operations available on ApiHostnameCoverageMatchTargets  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapihostnamecoveragematchtargets
	ApiHostnameCoverageMatchTargets interface {
		GetApiHostnameCoverageMatchTargets(ctx context.Context, params GetApiHostnameCoverageMatchTargetsRequest) (*GetApiHostnameCoverageMatchTargetsResponse, error)
	}
	GetApiHostnameCoverageMatchTargetsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	GetApiHostnameCoverageMatchTargetsResponse struct {
		MatchTargets struct {
			WebsiteTargets []struct {
				Type                         string                                                `json:"type"`
				BypassNetworkLists           *HostnameCoverageMatchTargetBypassNetworkLists        `json:"bypassNetworkLists,omitempty"`
				ConfigID                     int                                                   `json:"configId"`
				ConfigVersion                int                                                   `json:"configVersion"`
				DefaultFile                  string                                                `json:"defaultFile"`
				EffectiveSecurityControls    *HostnameCoverageMatchTargetEffectiveSecurityControls `json:"effectiveSecurityControls,omitempty"`
				FilePaths                    []string                                              `json:"filePaths"`
				Hostnames                    []string                                              `json:"hostnames"`
				IsNegativeFileExtensionMatch bool                                                  `json:"isNegativeFileExtensionMatch"`
				IsNegativePathMatch          bool                                                  `json:"isNegativePathMatch"`
				SecurityPolicy               struct {
					PolicyID string `json:"policyId"`
				} `json:"securityPolicy"`
				Sequence int `json:"sequence"`
				TargetID int `json:"targetId"`
			} `json:"websiteTargets"`
			APITargets []interface{} `json:"apiTargets"`
		} `json:"matchTargets"`
	}

	HostnameCoverageMatchTargetBypassNetworkLists []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	HostnameCoverageMatchTargetEffectiveSecurityControls struct {
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates GetApiHostnameCoverageMatchTargetsRequest
func (v GetApiHostnameCoverageMatchTargetsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiHostnameCoverageMatchTargets(ctx context.Context, params GetApiHostnameCoverageMatchTargetsRequest) (*GetApiHostnameCoverageMatchTargetsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverageMatchTargets")

	var rval GetApiHostnameCoverageMatchTargetsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/hostname-coverage/match-targets?hostname=%s",
		params.ConfigID,
		params.Version,
		params.Hostname)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapihostnamecoveragematchtargets request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapihostnamecoveragematchtargets  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
