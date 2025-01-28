package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SiemSettings interface supports retrieving, modifying and removing the SIEM settings for a configuration.
	SiemSettings interface {
		// GetSiemSettings returns SIEM settings for a specific configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-siem
		GetSiemSettings(ctx context.Context, params GetSiemSettingsRequest) (*GetSiemSettingsResponse, error)

		// UpdateSiemSettings updates SIEM settings for a specific configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-siem
		UpdateSiemSettings(ctx context.Context, params UpdateSiemSettingsRequest) (*UpdateSiemSettingsResponse, error)

		// RemoveSiemSettings removes SIEM settings for a specific configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-siem
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
		EnableForAllPolicies    bool        `json:"enableForAllPolicies"`
		EnableSiem              bool        `json:"enableSiem"`
		EnabledBotmanSiemEvents bool        `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int         `json:"siemDefinitionId"`
		FirewallPolicyIDs       []string    `json:"firewallPolicyIds"`
		Exceptions              []Exception `json:"exceptions"`
	}

	// GetSiemSettingRequest is used to retrieve the SIEM settings for a configuration.
	GetSiemSettingRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// Exception is used to create exceptions list for SIEM events
	Exception struct {
		Protection  string   `json:"protection"`
		ActionTypes []string `json:"actionTypes"`
	}

	// GetSiemSettingResponse is returned from a call to GetSiemSettings.
	GetSiemSettingResponse struct {
		EnableForAllPolicies    bool        `json:"enableForAllPolicies"`
		EnableSiem              bool        `json:"enableSiem"`
		EnabledBotmanSiemEvents *bool       `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int         `json:"siemDefinitionId"`
		FirewallPolicyIDs       []string    `json:"firewallPolicyIds"`
		Exceptions              []Exception `json:"exceptions"`
	}

	// UpdateSiemSettingsRequest is used to modify the SIEM settings for a configuration.
	UpdateSiemSettingsRequest struct {
		ConfigID                int         `json:"-"`
		Version                 int         `json:"-"`
		EnableForAllPolicies    bool        `json:"enableForAllPolicies"`
		EnableSiem              bool        `json:"enableSiem"`
		EnabledBotmanSiemEvents *bool       `json:"enabledBotmanSiemEvents,omitempty"`
		SiemDefinitionID        int         `json:"siemDefinitionId"`
		FirewallPolicyIDs       []string    `json:"firewallPolicyIds"`
		Exceptions              []Exception `json:"exceptions,omitempty"`
	}

	// UpdateSiemSettingsResponse is returned from a call to UpdateSiemSettings.
	UpdateSiemSettingsResponse struct {
		EnableForAllPolicies    bool        `json:"enableForAllPolicies"`
		EnableSiem              bool        `json:"enableSiem"`
		EnabledBotmanSiemEvents *bool       `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int         `json:"siemDefinitionId"`
		FirewallPolicyIDs       []string    `json:"firewallPolicyIds"`
		Exceptions              []Exception `json:"exceptions"`
	}

	// RemoveSiemSettingsRequest is used to remove the SIEM settings for a configuration.
	// Deprecated: this struct will be removed in a future release.
	RemoveSiemSettingsRequest struct {
		ConfigID                int      `json:"-"`
		Version                 int      `json:"-"`
		EnableForAllPolicies    bool     `json:"-"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents *bool    `json:"-"`
		SiemDefinitionID        int      `json:"-"`
		FirewallPolicyIDs       []string `json:"-"`
	}

	// RemoveSiemSettingsResponse is returned from a call to RemoveSiemSettings.
	// Deprecated: this struct will be removed in a future release.
	RemoveSiemSettingsResponse struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents *bool    `json:"enabledBotmanSiemEvents"`
		SiemDefinitionID        int      `json:"siemDefinitionId"`
		FirewallPolicyIDs       []string `json:"firewallPolicyIds"`
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

// Validate validates an Exception struct.
func (v Exception) Validate() error {
	validActionTypes := []string{"alert", "deny", "all_custom", "abort", "allow", "delay", "ignore", "monitor", "slow", "tarpit", "*"}
	return validation.Errors{
		"Protection": validation.Validate(v.Protection, validation.Required, validation.In("botmanagement", "ipgeo", "rate", "urlProtection", "slowpost", "customrules", "waf", "apirequestconstraints", "clientrep", "malwareprotection", "aprProtection").
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'botmanagement', 'ipgeo', 'rate', 'urlProtection', 'slowpost', 'customrules', 'waf', 'apirequestconstraints', 'clientrep', 'malwareprotection', 'aprProtection'", v.Protection))),
		"ActionTypes": validation.ValidateStruct(&v,
			validation.Field(&v.ActionTypes, validation.Required, validation.By(func(value interface{}) error {
				actions, _ := value.([]string)
				for _, actionType := range actions {
					if !containsElement(validActionTypes, actionType) {
						return fmt.Errorf("value '%s' is invalid. Must be one of: %v", actionType, validActionTypes)
					}
				}
				return nil
			}))),
	}.Filter()
}

func containsElement(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
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

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/siem",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSiemSettings request: %w", err)
	}

	var result GetSiemSettingsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get siem settings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

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

	for _, exception := range params.Exceptions {
		if err := exception.Validate(); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
		}
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
		return nil, fmt.Errorf("update siem settings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

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
		return nil, fmt.Errorf("remove siem settings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
