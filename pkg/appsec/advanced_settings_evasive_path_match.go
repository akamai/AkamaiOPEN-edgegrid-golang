package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsEvasivePathMatch interface supports retrieving or modifying the Evasive Path Match setting.
	AdvancedSettingsEvasivePathMatch interface {
		// GetAdvancedSettingsEvasivePathMatch retrieves the Evasive Path Match setting.
		// https://techdocs.akamai.com/application-security/reference/get-evasive-path-match-per-config
		GetAdvancedSettingsEvasivePathMatch(ctx context.Context, params GetAdvancedSettingsEvasivePathMatchRequest) (*GetAdvancedSettingsEvasivePathMatchResponse, error)
		// UpdateAdvancedSettingsEvasivePathMatch modifies the Evasive Path Match setting.
		// https://techdocs.akamai.com/application-security/reference/put-evasive-path-match-per-config
		UpdateAdvancedSettingsEvasivePathMatch(ctx context.Context, params UpdateAdvancedSettingsEvasivePathMatchRequest) (*UpdateAdvancedSettingsEvasivePathMatchResponse, error)
		// RemoveAdvancedSettingsEvasivePathMatch removes the Evasive Path Match setting.
		// Deprecated: this method will be removed in a future release. Use UpdateAdvancedSettingsEvasivePathMatch instead.
		RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, params RemoveAdvancedSettingsEvasivePathMatchRequest) (*RemoveAdvancedSettingsEvasivePathMatchResponse, error)
	}

	// GetAdvancedSettingsEvasivePathMatchRequest is used to retrieve the EvasivePathMatch setting
	GetAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// GetAdvancedSettingsEvasivePathMatchResponse returns the EvasivePathMatch setting
	GetAdvancedSettingsEvasivePathMatchResponse struct {
		EnablePathMatch bool `json:"enablePathMatch"`
	}

	// UpdateAdvancedSettingsEvasivePathMatchRequest is used to update the EvasivePathMatch setting
	UpdateAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		EnablePathMatch bool   `json:"enablePathMatch"`
	}

	// UpdateAdvancedSettingsEvasivePathMatchResponse returns the result of updating the EvasivePathMatch setting
	UpdateAdvancedSettingsEvasivePathMatchResponse struct {
		EnablePathMatch bool `json:"enablePathMatch"`
	}

	// RemoveAdvancedSettingsEvasivePathMatchRequest is used to clear the EvasivePathMatch setting
	RemoveAdvancedSettingsEvasivePathMatchRequest struct {
		ConfigID        int    `json:"-"`
		Version         int    `json:"-"`
		PolicyID        string `json:"-"`
		EnablePathMatch bool   `json:"enablePathMatch"`
	}

	// RemoveAdvancedSettingsEvasivePathMatchResponse returns the result of clearing the EvasivePathMatch setting
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

func (p *appsec) RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, params RemoveAdvancedSettingsEvasivePathMatchRequest) (*RemoveAdvancedSettingsEvasivePathMatchResponse, error) {
	request := UpdateAdvancedSettingsEvasivePathMatchRequest{
		ConfigID:        params.ConfigID,
		Version:         params.Version,
		PolicyID:        params.PolicyID,
		EnablePathMatch: false,
	}
	_, err := p.UpdateAdvancedSettingsEvasivePathMatch(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("UpdateAdvancedSettingsEvasivePathMatch request failed: %w", err)
	}
	response := RemoveAdvancedSettingsEvasivePathMatchResponse{
		ConfigID:        params.ConfigID,
		Version:         params.Version,
		PolicyID:        params.PolicyID,
		EnablePathMatch: false,
	}
	return &response, nil
}
