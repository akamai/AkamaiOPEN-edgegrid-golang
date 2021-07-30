package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AdvancedSettingsEvasivePathMatch represents a collection of AdvancedSettingsEvasivePathMatch
//
// See: AdvancedSettingsEvasivePathMatch.GetAdvancedSettingsEvasivePathMatch()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AdvancedSettingsEvasivePathMatch  contains operations available on AdvancedSettingsEvasivePathMatch  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getadvancedsettingsEvasivePathMatch
	AdvancedSettingsEvasivePathMatch interface {
		GetAdvancedSettingsEvasivePathMatch(ctx context.Context, params GetAdvancedSettingsEvasivePathMatchRequest) (*GetAdvancedSettingsEvasivePathMatchResponse, error)
		UpdateAdvancedSettingsEvasivePathMatch(ctx context.Context, params UpdateAdvancedSettingsEvasivePathMatchRequest) (*UpdateAdvancedSettingsEvasivePathMatchResponse, error)
		RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, params RemoveAdvancedSettingsEvasivePathMatchRequest) (*RemoveAdvancedSettingsEvasivePathMatchResponse, error)
	}

	GetAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetAdvancedSettingsEvasivePathMatchResponse struct {
		EnablePathMatch bool `json:"enablePathMatch"`
	}

	UpdateAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		EnablePathMatch bool   `json:"enablePathMatch"`
	}

	UpdateAdvancedSettingsEvasivePathMatchResponse struct {
		EnablePathMatch bool `json:"enablePathMatch"`
	}
	RemoveAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		EnablePathMatch bool   `json:"enablePathMatch"`
	}

	RemoveAdvancedSettingsEvasivePathMatchResponse struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		EnablePathMatch bool   `json:"enablePathMatch"`
	}
)

// Validate validates GetAdvancedSettingssEvasivePathMatchRequest
func (v GetAdvancedSettingsEvasivePathMatchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsEvasivePathMatchRequest
func (v UpdateAdvancedSettingsEvasivePathMatchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsEvasivePathMatchRequest
func (v RemoveAdvancedSettingsEvasivePathMatchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsEvasivePathMatch(ctx context.Context, params GetAdvancedSettingsEvasivePathMatchRequest) (*GetAdvancedSettingsEvasivePathMatchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsLoggings")

	var rval GetAdvancedSettingsEvasivePathMatchResponse
	var uri string

	if params.PolicyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/evasive-path-match",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/evasive-path-match",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getadvancedsettingsloggings request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getadvancedsettingsloggings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AdvancedSettingsEvasivePathMatch.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingsEvasivePathMatch

func (p *appsec) UpdateAdvancedSettingsEvasivePathMatch(ctx context.Context, params UpdateAdvancedSettingsEvasivePathMatchRequest) (*UpdateAdvancedSettingsEvasivePathMatchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsLogging")

	var putURL string
	if params.PolicyID != "" {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/evasive-path-match",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/evasive-path-match",
			params.ConfigID,
			params.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AdvancedSettingsLoggingrequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	var rval UpdateAdvancedSettingsEvasivePathMatchResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AdvancedSettingsLogging request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Remove will update a AdvancedSettingsEvasivePathMatch.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putadvancedsettingsEvasivePathMatch

func (p *appsec) RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, params RemoveAdvancedSettingsEvasivePathMatchRequest) (*RemoveAdvancedSettingsEvasivePathMatchResponse, error) {
	_, err := p.RemoveAdvancedSettingsEvasivePathMatch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveAdvancedSettingsEvasivePathMatch request failed: %w", err)
	}
	response := RemoveAdvancedSettingsEvasivePathMatchResponse{}
	return &response, nil
}
