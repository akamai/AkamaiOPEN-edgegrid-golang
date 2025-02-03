package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SlowPostProtectionSetting interface supports retrieving and updating the slow POST protection settings for a configuration and policy.
	SlowPostProtectionSetting interface {
		// GetSlowPostProtectionSettings retrieves the current SLOW post protection settings for a configuration and policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-slow-post
		GetSlowPostProtectionSettings(ctx context.Context, params GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error)

		// UpdateSlowPostProtectionSetting updates the SLOW post protection settings for a configuration and policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-slow-post
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

func (p *appsec) GetSlowPostProtectionSettings(ctx context.Context, params GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSlowPostProtectionSettings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/slow-post",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSlowPostProtectionSettings request: %w", err)
	}

	var result GetSlowPostProtectionSettingsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get slow post protection settings request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateSlowPostProtectionSetting(ctx context.Context, params UpdateSlowPostProtectionSettingRequest) (*UpdateSlowPostProtectionSettingResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSlowPostProtectionSetting")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/slow-post",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update UpdateSlowPostProtectionSetting request: %w", err)
	}

	var result UpdateSlowPostProtectionSettingResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update slow post protection setting request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
