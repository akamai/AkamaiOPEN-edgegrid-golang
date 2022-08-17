package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsPrefetch interface supports retrieving or modifying the prefetch request settings
	// for a configuration.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#prefetchrequest
	AdvancedSettingsPrefetch interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprefetchrequestsforaconfiguration
		GetAdvancedSettingsPrefetch(ctx context.Context, params GetAdvancedSettingsPrefetchRequest) (*GetAdvancedSettingsPrefetchResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putprefetchrequestsforaconfiguration
		UpdateAdvancedSettingsPrefetch(ctx context.Context, params UpdateAdvancedSettingsPrefetchRequest) (*UpdateAdvancedSettingsPrefetchResponse, error)
	}

	// GetAdvancedSettingsPrefetchRequest is used to retrieve the prefetch request settings.
	GetAdvancedSettingsPrefetchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	// GetAdvancedSettingsPrefetchResponse is returned from a call to GetAdvancedSettingsPrefetch.
	GetAdvancedSettingsPrefetchResponse struct {
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	// UpdateAdvancedSettingsPrefetchRequest is used to modify the prefetch request settings.
	UpdateAdvancedSettingsPrefetchRequest struct {
		ConfigID           int      `json:"-"`
		Version            int      `json:"-"`
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	// UpdateAdvancedSettingsPrefetchResponse is returned from a call to UpdateAdvancedSettingsPrefetch.
	UpdateAdvancedSettingsPrefetchResponse struct {
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	// RemoveAdvancedSettingsPrefetchRequest is used to remove the prefetch request settings.
	RemoveAdvancedSettingsPrefetchRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
	}

	// RemoveAdvancedSettingsPrefetchResponse is returned from a call to RemoveAdvancedSettingsPrefetch.
	RemoveAdvancedSettingsPrefetchResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates a GetAdvancedSettingsPrefetchRequest.
func (v GetAdvancedSettingsPrefetchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAdvancedSettingsPrefetchRequest.
func (v UpdateAdvancedSettingsPrefetchRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsPrefetch(ctx context.Context, params GetAdvancedSettingsPrefetchRequest) (*GetAdvancedSettingsPrefetchResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPrefetch")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetAdvancedSettingsPrefetchResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/prefetch",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsPrefetch request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAdvancedSettingsPrefetch request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) UpdateAdvancedSettingsPrefetch(ctx context.Context, params UpdateAdvancedSettingsPrefetchRequest) (*UpdateAdvancedSettingsPrefetchResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsPrefetch")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/prefetch",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsPrefetch request: %w", err)
	}

	var result UpdateAdvancedSettingsPrefetchResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateAdvancedSettingsPrefetch request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
