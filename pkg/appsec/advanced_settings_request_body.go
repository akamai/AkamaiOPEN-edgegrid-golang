package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsRequestBody interface supports retrieving, updating or removing settings
	// related to Request Size Inspection Limit.
	AdvancedSettingsRequestBody interface {
		// GetAdvancedSettingsRequestBody lists the Request Size Inspection Limit settings for a configuration or policy. If
		// the request specifies a policy, then the settings for that policy will be returned, otherwise the
		// settings for the configuration will be returned.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-advanced-settings-request-body
		GetAdvancedSettingsRequestBody(ctx context.Context, params GetAdvancedSettingsRequestBodyRequest) (*GetAdvancedSettingsRequestBodyResponse, error)

		// UpdateAdvancedSettingsRequestBody updates the Request Size Inspection Limit settings for a
		// configuration or policy. If the request specifies a policy, then the settings for that policy will be
		// updated, otherwise the settings for the configuration will be updated.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-advanced-settings-request-body
		UpdateAdvancedSettingsRequestBody(ctx context.Context, params UpdateAdvancedSettingsRequestBodyRequest) (*UpdateAdvancedSettingsRequestBodyResponse, error)

		// RemoveAdvancedSettingsRequestBody updates the Request Size Inspection Limit settings to default for a
		// configuration or policy. If the request specifies a policy, then the settings for that policy will be
		// updated, otherwise the settings for the configuration will be updated.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-advanced-settings-request-body
		RemoveAdvancedSettingsRequestBody(ctx context.Context, params RemoveAdvancedSettingsRequestBodyRequest) (*RemoveAdvancedSettingsRequestBodyResponse, error)
	}

	// GetAdvancedSettingsRequestBodyRequest is used to retrieve the Request Size Inspection Limit settings for a configuration or policy.
	GetAdvancedSettingsRequestBodyRequest struct {
		ConfigID int
		Version  int
		PolicyID string
	}

	// GetAdvancedSettingsRequestBodyResponse is returned from a call to GetAdvancedSettingsRequestBody.
	GetAdvancedSettingsRequestBodyResponse struct {
		RequestBodyInspectionLimitInKB     RequestBodySizeLimit `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool                 `json:"override"`
	}

	// UpdateAdvancedSettingsRequestBodyRequest is used to update the Request body settings for a configuration or policy.
	UpdateAdvancedSettingsRequestBodyRequest struct {
		ConfigID                           int
		Version                            int
		PolicyID                           string
		RequestBodyInspectionLimitInKB     RequestBodySizeLimit `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool                 `json:"override"`
	}

	// UpdateAdvancedSettingsRequestBodyResponse is returned from a call to UpdateAdvancedSettingsRequestBody.
	UpdateAdvancedSettingsRequestBodyResponse struct {
		RequestBodyInspectionLimitInKB     RequestBodySizeLimit `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool                 `json:"override"`
	}

	// RemoveAdvancedSettingsRequestBodyRequest is used to reset the Request body settings for a configuration or policy.
	RemoveAdvancedSettingsRequestBodyRequest struct {
		ConfigID                           int
		Version                            int
		PolicyID                           string
		RequestBodyInspectionLimitInKB     RequestBodySizeLimit `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool                 `json:"override"`
	}

	// RemoveAdvancedSettingsRequestBodyResponse is returned from a call to RemoveAdvancedSettingsRequestBody.
	RemoveAdvancedSettingsRequestBodyResponse struct {
		RequestBodyInspectionLimitInKB     RequestBodySizeLimit `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool                 `json:"override"`
	}

	// RequestBodySizeLimit is used to create an "enum" of possible types default, 8, 16, 32
	RequestBodySizeLimit string
)

const (
	// Default RequestBodySize
	Default RequestBodySizeLimit = "default"
	// Limit8KB RequestBodySize
	Limit8KB RequestBodySizeLimit = "8"
	// Limit16KB RequestBodySize
	Limit16KB RequestBodySizeLimit = "16"
	// Limit32KB RequestBodySize
	Limit32KB RequestBodySizeLimit = "32"
)

// Validate validates a GetAdvancedSettingsRequestBodyRequest.
func (v GetAdvancedSettingsRequestBodyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an UpdateAdvancedSettingsRequestBodyRequest.
func (v UpdateAdvancedSettingsRequestBodyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an RemoveAdvancedSettingsRequestBodyRequest.
func (v RemoveAdvancedSettingsRequestBodyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

func (a *appsec) GetAdvancedSettingsRequestBody(ctx context.Context, params GetAdvancedSettingsRequestBodyRequest) (*GetAdvancedSettingsRequestBodyResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetAdvancedSettingsRequestBody")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRequestBodyURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsRequestBody request: %w", err)
	}

	var result GetAdvancedSettingsRequestBodyResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings request body request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func (a *appsec) UpdateAdvancedSettingsRequestBody(ctx context.Context, params UpdateAdvancedSettingsRequestBodyRequest) (*UpdateAdvancedSettingsRequestBodyResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsRequestBody")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRequestBodyURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsRequestBody request: %w", err)
	}

	var result UpdateAdvancedSettingsRequestBodyResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings request body request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func (a *appsec) RemoveAdvancedSettingsRequestBody(ctx context.Context, params RemoveAdvancedSettingsRequestBodyRequest) (*RemoveAdvancedSettingsRequestBodyResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("RemoveAdvancedSettingsRequestBody")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getRequestBodyURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsRequestBody request: %w", err)
	}

	var result RemoveAdvancedSettingsRequestBodyResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove advanced settings request body request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func getRequestBodyURI(configID, configVersion int, policyID string) string {
	var uri string
	if policyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/request-body", configID, configVersion, policyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/request-body", configID, configVersion)
	}
	return uri
}
