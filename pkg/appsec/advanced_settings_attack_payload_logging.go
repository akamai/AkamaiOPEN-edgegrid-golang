package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsAttackPayloadLogging interface supports retrieving, updating or removing settings
	// related to Attack Payload logging.
	AdvancedSettingsAttackPayloadLogging interface {
		// GetAdvancedSettingsAttackPayloadLogging lists the attack payload logging settings for a configuration or policy. If
		// the request specifies a policy, then the settings for that policy will be returned, otherwise the
		// settings for the configuration will be returned.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-advanced-settings-attack-payload-logging
		GetAdvancedSettingsAttackPayloadLogging(ctx context.Context, params GetAdvancedSettingsAttackPayloadLoggingRequest) (*GetAdvancedSettingsAttackPayloadLoggingResponse, error)

		// UpdateAdvancedSettingsAttackPayloadLogging enables, disables, or updates the attack payload logging settings for a
		// configuration or policy. If the request specifies a policy, then the settings for that policy will be
		// updated, otherwise the settings for the configuration will be updated.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-advanced-settings-attack-payload-logging
		UpdateAdvancedSettingsAttackPayloadLogging(ctx context.Context, params UpdateAdvancedSettingsAttackPayloadLoggingRequest) (*UpdateAdvancedSettingsAttackPayloadLoggingResponse, error)

		// RemoveAdvancedSettingsAttackPayloadLogging disables attack payload logging for a configuration or policy. If the request
		// specifies a policy, then attack payload logging will be disabled for that policy, otherwise logging will be
		// disabled for the configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-advanced-settings-attack-payload-logging
		RemoveAdvancedSettingsAttackPayloadLogging(ctx context.Context, params RemoveAdvancedSettingsAttackPayloadLoggingRequest) (*RemoveAdvancedSettingsAttackPayloadLoggingResponse, error)
	}

	// GetAdvancedSettingsAttackPayloadLoggingRequest is used to retrieve the Attack Payload logging settings for a configuration or policy.
	GetAdvancedSettingsAttackPayloadLoggingRequest struct {
		ConfigID int
		Version  int
		PolicyID string
	}

	// GetAdvancedSettingsAttackPayloadLoggingResponse is returned from a call to GetAdvancedSettingsAttackPayloadLogging.
	GetAdvancedSettingsAttackPayloadLoggingResponse struct {
		Override     bool                             `json:"override"`
		Enabled      bool                             `json:"enabled"`
		RequestBody  AttackPayloadLoggingRequestBody  `json:"requestBody"`
		ResponseBody AttackPayloadLoggingResponseBody `json:"responseBody"`
	}

	// AttackPayloadLoggingRequestBody Type field represents whether attack payload is logged or not logged for RequestBody.
	AttackPayloadLoggingRequestBody struct {
		Type AttackPayloadType `json:"type"`
	}

	// AttackPayloadLoggingResponseBody Type field represents whether attack payload is logged or not logged for ResponseBody.
	AttackPayloadLoggingResponseBody struct {
		Type AttackPayloadType `json:"type"`
	}

	// UpdateAdvancedSettingsAttackPayloadLoggingRequest is used to update the Attack Payload logging settings for a configuration or policy.
	UpdateAdvancedSettingsAttackPayloadLoggingRequest struct {
		ConfigID       int
		Version        int
		PolicyID       string
		JSONPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateAdvancedSettingsAttackPayloadLoggingResponse is returned from a call to UpdateAdvancedSettingsAttackPayloadLogging.
	UpdateAdvancedSettingsAttackPayloadLoggingResponse struct {
		Override     bool                             `json:"override"`
		Enabled      bool                             `json:"enabled"`
		RequestBody  AttackPayloadLoggingRequestBody  `json:"requestBody"`
		ResponseBody AttackPayloadLoggingResponseBody `json:"responseBody"`
	}

	// RemoveAdvancedSettingsAttackPayloadLoggingRequest is used to disable Attack Payload logging for a configuration or policy.
	RemoveAdvancedSettingsAttackPayloadLoggingRequest struct {
		ConfigID     int
		Version      int
		PolicyID     string
		Override     bool                             `json:"override"`
		Enabled      bool                             `json:"enabled"`
		RequestBody  AttackPayloadLoggingRequestBody  `json:"requestBody"`
		ResponseBody AttackPayloadLoggingResponseBody `json:"responseBody"`
	}

	// RemoveAdvancedSettingsAttackPayloadLoggingResponse is returned from a call to RemoveAdvancedSettingsAttackPayloadLogging.
	RemoveAdvancedSettingsAttackPayloadLoggingResponse struct {
		Override     bool                             `json:"override"`
		Enabled      bool                             `json:"enabled"`
		RequestBody  AttackPayloadLoggingRequestBody  `json:"requestBody"`
		ResponseBody AttackPayloadLoggingResponseBody `json:"responseBody"`
	}

	// AttackPayloadType is used to create an "enum" of possible types ATTACK_PAYLOAD or NONE
	AttackPayloadType string
)

const (
	// AttackPayload AttackPayloadType
	AttackPayload AttackPayloadType = "ATTACK_PAYLOAD"
	// None AttackPayloadType
	None AttackPayloadType = "NONE"
)

// Validate validates a GetAdvancedSettingsAttackPayloadLoggingRequest.
func (v GetAdvancedSettingsAttackPayloadLoggingRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an UpdateAdvancedSettingsAttackPayloadLoggingRequest.
func (v UpdateAdvancedSettingsAttackPayloadLoggingRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a RemoveAdvancedSettingsAttackPayloadLoggingRequest.
func (v RemoveAdvancedSettingsAttackPayloadLoggingRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

func (a *appsec) GetAdvancedSettingsAttackPayloadLogging(ctx context.Context, params GetAdvancedSettingsAttackPayloadLoggingRequest) (*GetAdvancedSettingsAttackPayloadLoggingResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetAdvancedSettingsAttackPayloadLogging")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsAttackPayloadLogging request: %w", err)
	}

	var result GetAdvancedSettingsAttackPayloadLoggingResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings attack payload logging request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func (a *appsec) UpdateAdvancedSettingsAttackPayloadLogging(ctx context.Context, params UpdateAdvancedSettingsAttackPayloadLoggingRequest) (*UpdateAdvancedSettingsAttackPayloadLoggingResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsAttackPayloadLogging")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsAttackPayloadLogging request: %w", err)
	}

	var result UpdateAdvancedSettingsAttackPayloadLoggingResponse
	resp, err := a.Exec(req, &result, params.JSONPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings attack payload logging request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func (a *appsec) RemoveAdvancedSettingsAttackPayloadLogging(ctx context.Context, params RemoveAdvancedSettingsAttackPayloadLoggingRequest) (*RemoveAdvancedSettingsAttackPayloadLoggingResponse, error) {

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := getURI(params.ConfigID, params.Version, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveAdvancedSettingsAttackPayloadLogging request: %w", err)
	}

	var result RemoveAdvancedSettingsAttackPayloadLoggingResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove advanced settings attack payload logging request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, a.Error(resp)
	}

	return &result, nil
}

func getURI(configID, configVersion int, policyID string) string {
	var uri string
	if policyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/advanced-settings/logging/attack-payload", configID, configVersion, policyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/advanced-settings/logging/attack-payload", configID, configVersion)
	}
	return uri
}
