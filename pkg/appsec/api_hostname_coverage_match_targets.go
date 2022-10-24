package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiHostnameCoverageMatchTargets interface supports retrieving the API and website
	// match targets that protect a hostname.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#hostnamecoverage
	ApiHostnameCoverageMatchTargets interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethostnamecoveragematchtargets
		GetApiHostnameCoverageMatchTargets(ctx context.Context, params GetApiHostnameCoverageMatchTargetsRequest) (*GetApiHostnameCoverageMatchTargetsResponse, error)
	}

	// GetApiHostnameCoverageMatchTargetsRequest is used to retrieve the API and website match targets that protect a hostname.
	GetApiHostnameCoverageMatchTargetsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	// GetApiHostnameCoverageMatchTargetsResponse is returned from a call to GetApiHostnameCoverageMatchTargets.
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

	// HostnameCoverageMatchTargetBypassNetworkLists describes a network list included in the list of bypass network lists.
	HostnameCoverageMatchTargetBypassNetworkLists []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// HostnameCoverageMatchTargetEffectiveSecurityControls describes the effective security controls for a website target.
	HostnameCoverageMatchTargetEffectiveSecurityControls struct {
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
)

// Validate validates a GetApiHostnameCoverageMatchTargetsRequest.
func (v GetApiHostnameCoverageMatchTargetsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiHostnameCoverageMatchTargets(ctx context.Context, params GetApiHostnameCoverageMatchTargetsRequest) (*GetApiHostnameCoverageMatchTargetsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverageMatchTargets")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/hostname-coverage/match-targets?hostname=%s",
		params.ConfigID,
		params.Version,
		params.Hostname)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetApiHostnameCoverageMatchTargets request: %w", err)
	}

	var result GetApiHostnameCoverageMatchTargetsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get API hostname coverage match targets request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
