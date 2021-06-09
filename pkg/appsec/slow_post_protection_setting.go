package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SlowPostProtectionSetting represents a collection of SlowPostProtectionSetting
//
// See: SlowPostProtectionSetting.GetSlowPostProtectionSetting()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// SlowPostProtectionSetting  contains operations available on SlowPostProtectionSetting  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getslowpostprotectionsetting
	SlowPostProtectionSetting interface {
		GetSlowPostProtectionSettings(ctx context.Context, params GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error)
		GetSlowPostProtectionSetting(ctx context.Context, params GetSlowPostProtectionSettingRequest) (*GetSlowPostProtectionSettingResponse, error)
		UpdateSlowPostProtectionSetting(ctx context.Context, params UpdateSlowPostProtectionSettingRequest) (*UpdateSlowPostProtectionSettingResponse, error)
	}

	GetSlowPostProtectionSettingsRequest struct {
		ConfigID          int                                         `json:"configId"`
		Version           int                                         `json:"version"`
		PolicyID          string                                      `json:"policyId"`
		Action            string                                      `json:"action"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	SlowPostProtectionSettingSlowRateThreshold struct {
		Rate   int `json:"rate"`
		Period int `json:"period"`
	}

	SlowPostProtectionSettingDurationThreshold struct {
		Timeout int `json:"timeout"`
	}

	GetSlowPostProtectionSettingsResponse struct {
		Action            string                                      `json:"action,omitempty"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	GetSlowPostProtectionSettingRequest struct {
		ConfigID          int    `json:"configId"`
		Version           int    `json:"version"`
		PolicyID          string `json:"policyId"`
		Action            string `json:"action"`
		SlowRateThreshold struct {
			Rate   int `json:"rate"`
			Period int `json:"period"`
		} `json:"slowRateThreshold"`
		DurationThreshold struct {
			Timeout int `json:"timeout"`
		} `json:"durationThreshold"`
	}

	GetSlowPostProtectionSettingResponse struct {
		Action            string                                      `json:"action,omitempty"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	UpdateSlowPostProtectionSettingRequest struct {
		ConfigID          int    `json:"configId"`
		Version           int    `json:"version"`
		PolicyID          string `json:"policyId"`
		Action            string `json:"action"`
		SlowRateThreshold struct {
			Rate   int `json:"rate"`
			Period int `json:"period"`
		} `json:"slowRateThreshold"`
		DurationThreshold struct {
			Timeout int `json:"timeout"`
		} `json:"durationThreshold"`
	}

	UpdateSlowPostProtectionSettingResponse struct {
		Action            string `json:"action"`
		SlowRateThreshold struct {
			Rate   int `json:"rate"`
			Period int `json:"period"`
		} `json:"slowRateThreshold"`
		DurationThreshold struct {
			Timeout int `json:"timeout"`
		} `json:"durationThreshold"`
	}
)

// Validate validates GetSlowPostProtectionSettingRequest
func (v GetSlowPostProtectionSettingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetSlowPostProtectionSettingsRequest
func (v GetSlowPostProtectionSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateSlowPostProtectionSettingRequest
func (v UpdateSlowPostProtectionSettingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetSlowPostProtectionSetting(ctx context.Context, params GetSlowPostProtectionSettingRequest) (*GetSlowPostProtectionSettingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtectionSetting")

	var rval GetSlowPostProtectionSettingResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/slow-post",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getslowpostprotectionsetting request: %w", err)
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

func (p *appsec) GetSlowPostProtectionSettings(ctx context.Context, params GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtectionSettings")

	var rval GetSlowPostProtectionSettingsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/slow-post",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getslowpostprotectionsettings request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getslowpostprotectionsettings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a SlowPostProtectionSetting.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putslowpostprotectionsetting

func (p *appsec) UpdateSlowPostProtectionSetting(ctx context.Context, params UpdateSlowPostProtectionSettingRequest) (*UpdateSlowPostProtectionSettingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateSlowPostProtectionSetting")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/slow-post",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create SlowPostProtectionSettingrequest: %w", err)
	}

	var rval UpdateSlowPostProtectionSettingResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create SlowPostProtectionSetting request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
