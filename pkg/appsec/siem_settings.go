package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SiemSettings interface supports retrieving, modifying and removing the SIEM settings for a configuration.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#siem
	SiemSettings interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemsettings
		GetSiemSettings(ctx context.Context, params GetSiemSettingsRequest) (*GetSiemSettingsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsiemsettings
		UpdateSiemSettings(ctx context.Context, params UpdateSiemSettingsRequest) (*UpdateSiemSettingsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsiemsettings
		// Deprecated: this method will be removed in a future release.
		RemoveSiemSettings(ctx context.Context, params RemoveSiemSettingsRequest) (*RemoveSiemSettingsResponse, error)
	}

	// GetSiemSettingsRequest is used to retrieve the SIEM settings for a configuration.
	GetSiemSettingsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetSiemSettingsResponse is returned from a call to GetSiemSettings.
	GetSiemSettingsResponse struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds"`
	}

	// GetSiemSettingRequest is used to retrieve the SIEM settings for a configuration.
	GetSiemSettingRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetSiemSettingResponse is returned from a call to GetSiemSettings.
	GetSiemSettingResponse struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds"`
	}

	// UpdateSiemSettingsRequest is used to modify the SIEM settings for a configuration.
	UpdateSiemSettingsRequest struct {
		ConfigID                int      `json:"-"`
		Version                 int      `json:"-"`
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds"`
	}

	// UpdateSiemSettingsResponse is returned from a call to UpdateSiemSettings.
	UpdateSiemSettingsResponse struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds"`
	}

	// RemoveSiemSettingsRequest is used to remove the SIEM settings for a configuration.
	// Deprecated: this struct will be removed in a future release.
	RemoveSiemSettingsRequest struct {
		ConfigID                int      `json:"-"`
		Version                 int      `json:"-"`
		EnableForAllPolicies    bool     `json:"-"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"-"`
		SiemDefinitionID        int      `json:"-"`
		FirewallPolicyIds       []string `json:"-"`
	}

	// RemoveSiemSettingsResponse is returned from a call to RemoveSiemSettings.
	// Deprecated: this struct will be removed in a future release.
	RemoveSiemSettingsResponse struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds"`
	}
)

// Validate validates a GetSiemSettingsRequest.
func (v GetSiemSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSiemSettingsRequest.
func (v UpdateSiemSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a RemoveSiemSettingsRequest.
// Deprecated: this method will be removed in a future release.
func (v RemoveSiemSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetSiemSettings(ctx context.Context, params GetSiemSettingsRequest) (*GetSiemSettingsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSiemSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetSiemSettingsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/siem",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSiemSettings request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetSiemSettings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateSiemSettings(ctx context.Context, params UpdateSiemSettingsRequest) (*UpdateSiemSettingsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSiemSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/siem",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateSiemSettings request: %w", err)
	}

	var result UpdateSiemSettingsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSiemSettings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) RemoveSiemSettings(ctx context.Context, params RemoveSiemSettingsRequest) (*RemoveSiemSettingsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveSiemSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/siem",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveSiemSettings request: %w", err)
	}

	var result RemoveSiemSettingsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveSiemSettings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
