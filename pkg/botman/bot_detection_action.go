package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotDetectionAction interface supports retrieving and updating the actions for bot detections of a configuration.
	//
	BotDetectionAction interface {
		// GetBotDetectionActionList todo: add link
		GetBotDetectionActionList(ctx context.Context, params GetBotDetectionActionListRequest) (*GetBotDetectionActionListResponse, error)
		// GetBotDetectionAction todo: add link
		GetBotDetectionAction(ctx context.Context, params GetBotDetectionActionRequest) (map[string]interface{}, error)
		// UpdateBotDetectionAction todo: add link
		UpdateBotDetectionAction(ctx context.Context, params UpdateBotDetectionActionRequest) (map[string]interface{}, error)
	}

	// GetBotDetectionActionListRequest is used to retrieve the bot detection actions for a configuration.
	GetBotDetectionActionListRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		DetectionID      string
	}

	// GetBotDetectionActionListResponse is used to retrieve the bot detection actions for a configuration.
	GetBotDetectionActionListResponse struct {
		Actions []map[string]interface{} `json:"actions"`
	}

	// GetBotDetectionActionRequest is used to retrieve the action for a bot detection.
	GetBotDetectionActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		DetectionID      string
	}

	// UpdateBotDetectionActionRequest is used to modify a bot detection action.
	UpdateBotDetectionActionRequest struct {
		ConfigID         int64
		Version          int64
		SecurityPolicyID string
		DetectionID      string
		JsonPayload      json.RawMessage
	}
)

// Validate validates a GetBotDetectionActionRequest.
func (v GetBotDetectionActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"DetectionID":      validation.Validate(v.DetectionID, validation.Required),
	}.Filter()
}

// Validate validates a GetBotDetectionActionListRequest.
func (v GetBotDetectionActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBotDetectionActionRequest.
func (v UpdateBotDetectionActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.SecurityPolicyID, validation.Required),
		"DetectionID":      validation.Validate(v.DetectionID, validation.Required),
		"JsonPayload":      validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

func (b *botman) GetBotDetectionAction(ctx context.Context, params GetBotDetectionActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotDetectionAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bot-detection-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.DetectionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotDetectionAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotDetectionAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetBotDetectionActionList(ctx context.Context, params GetBotDetectionActionListRequest) (*GetBotDetectionActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotDetectionActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bot-detection-actions",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotDetectionActionList request: %w", err)
	}

	var result GetBotDetectionActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotDetectionActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetBotDetectionActionListResponse
	if params.DetectionID != "" {
		for _, val := range result.Actions {
			if val["detectionId"].(string) == params.DetectionID {
				filteredResult.Actions = append(filteredResult.Actions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateBotDetectionAction(ctx context.Context, params UpdateBotDetectionActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateBotDetectionAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bot-detection-actions/%s",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
		params.DetectionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBotDetectionAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateBotDetectionAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}
