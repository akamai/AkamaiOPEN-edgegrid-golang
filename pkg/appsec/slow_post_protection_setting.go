package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SlowPostProtectionSetting interface supports retrieving and modifying the slow POST protection
	// settings for a specific configuration.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#slowpostprotection
	SlowPostProtectionSetting interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getslowpostprotectionsettings
		GetSlowPostProtectionSettings(ctx context.Context, params GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getslowpostprotectionsettings
		GetSlowPostProtectionSetting(ctx context.Context, params GetSlowPostProtectionSettingRequest) (*GetSlowPostProtectionSettingResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putslowpostprotectionsettings
		UpdateSlowPostProtectionSetting(ctx context.Context, params UpdateSlowPostProtectionSettingRequest) (*UpdateSlowPostProtectionSettingResponse, error)
	}

	// GetSlowPostProtectionSettingsRequest is used to retrieve the slow post protection settings for a configuration.
	GetSlowPostProtectionSettingsRequest struct {
		ConfigID          int                                         `json:"configId"`
		Version           int                                         `json:"version"`
		PolicyID          string                                      `json:"policyId"`
		Action            string                                      `json:"action"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	// GetSlowPostProtectionSettingsResponse is returned from a call to GetSlowPostProtectionSettings.
	GetSlowPostProtectionSettingsResponse struct {
		Action            string                                      `json:"action,omitempty"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	// GetSlowPostProtectionSettingRequest is used to retrieve the slow post protection settings for a configuration.
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

	// GetSlowPostProtectionSettingResponse is returned from a call to GetSlowPostProtectionSetting.
	GetSlowPostProtectionSettingResponse struct {
		Action            string                                      `json:"action,omitempty"`
		SlowRateThreshold *SlowPostProtectionSettingSlowRateThreshold `json:"slowRateThreshold,omitempty"`
		DurationThreshold *SlowPostProtectionSettingDurationThreshold `json:"durationThreshold,omitempty"`
	}

	// UpdateSlowPostProtectionSettingRequest is used to modify the slow post protection settings for a configuration.
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

	// UpdateSlowPostProtectionSettingResponse is returned from a call to UpdateSlowPostProtection.
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

	// SlowPostProtectionSettingSlowRateThreshold describes the average rate in bytes per second over a specified period of time before an action
	// (alert or abort) in the policy triggers.
	SlowPostProtectionSettingSlowRateThreshold struct {
		Rate   int `json:"rate"`
		Period int `json:"period"`
	}

	// SlowPostProtectionSettingDurationThreshold describes the length of time in seconds within which the first eight kilobytes of the POST body must be transferred to
	// avoid applying the action specified in the policy.
	SlowPostProtectionSettingDurationThreshold struct {
		Timeout int `json:"timeout"`
	}
)

// Validate validates a GetSlowPostProtectionSettingRequest.
func (v GetSlowPostProtectionSettingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetSlowPostProtectionSettingsRequest.
func (v GetSlowPostProtectionSettingsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSlowPostProtectionSettingRequest.
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
		return nil, fmt.Errorf("failed to create GetSlowPostProtectionSetting request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetSlowPostProtectionSetting request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetSlowPostProtectionSettings request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetSlowPostProtectionSettings request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create update UpdateSlowPostProtectionSetting request: %w", err)
	}

	var rval UpdateSlowPostProtectionSettingResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSlowPostProtectionSetting request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
