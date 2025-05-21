package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsPIILearning interface supports retrieving or modifying the PII Learning setting.
	AdvancedSettingsPIILearning interface {
		// GetAdvancedSettingsPIILearning retrieves the PII Learning setting.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-pii-learning
		GetAdvancedSettingsPIILearning(ctx context.Context, params GetAdvancedSettingsPIILearningRequest) (*AdvancedSettingsPIILearningResponse, error)

		// UpdateAdvancedSettingsPIILearning modifies the PII Learning setting.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-pii-learning
		UpdateAdvancedSettingsPIILearning(ctx context.Context, params UpdateAdvancedSettingsPIILearningRequest) (*AdvancedSettingsPIILearningResponse, error)
	}

	// ConfigVersion is used to specify a security configuration and version.
	ConfigVersion struct {
		ConfigID int64
		Version  int
	}

	// GetAdvancedSettingsPIILearningRequest is used to retrieve the PIILearning setting.
	GetAdvancedSettingsPIILearningRequest struct {
		ConfigVersion
	}

	// UpdateAdvancedSettingsPIILearningRequest is used to update the PIILearning setting.
	UpdateAdvancedSettingsPIILearningRequest struct {
		ConfigVersion
		EnablePIILearning bool `json:"enablePiiLearning"`
	}

	// AdvancedSettingsPIILearningResponse returns the result of updating the PIILearning setting
	AdvancedSettingsPIILearningResponse struct {
		EnablePIILearning bool `json:"enablePiiLearning"`
	}
)

// Validate validates GetAdvancedSettingssPIILearningRequest
func (v GetAdvancedSettingsPIILearningRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates UpdateAdvancedSettingsPIILearningRequest
func (v UpdateAdvancedSettingsPIILearningRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

func (p *appsec) GetAdvancedSettingsPIILearning(ctx context.Context, params GetAdvancedSettingsPIILearningRequest) (*AdvancedSettingsPIILearningResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsPIILearning")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/pii-learning",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequestCreation, err.Error())
	}

	var result AdvancedSettingsPIILearningResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAPICallFailure, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAdvancedSettingsPIILearning(ctx context.Context, params UpdateAdvancedSettingsPIILearningRequest) (*AdvancedSettingsPIILearningResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsPIILearning")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/pii-learning",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequestCreation, err.Error())
	}

	var result AdvancedSettingsPIILearningResponse
	resp, err := p.Exec(req, &result, struct {
		EnablePIILearning bool `json:"enablePiiLearning"`
	}{
		EnablePIILearning: params.EnablePIILearning})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAPICallFailure, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
