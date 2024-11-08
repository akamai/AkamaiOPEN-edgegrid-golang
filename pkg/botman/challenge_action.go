package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ChallengeAction interface supports creating, retrieving, modifying and removing challenge action as well as
	// updating google recaptcha secret key for challenge action for a configuration.
	ChallengeAction interface {
		// GetChallengeActionList https://techdocs.akamai.com/bot-manager/reference/get-challenge-actions
		GetChallengeActionList(ctx context.Context, params GetChallengeActionListRequest) (*GetChallengeActionListResponse, error)

		// GetChallengeAction https://techdocs.akamai.com/bot-manager/reference/get-challenge-action
		GetChallengeAction(ctx context.Context, params GetChallengeActionRequest) (map[string]interface{}, error)

		// CreateChallengeAction https://techdocs.akamai.com/bot-manager/reference/post-challenge-action
		CreateChallengeAction(ctx context.Context, params CreateChallengeActionRequest) (map[string]interface{}, error)

		// UpdateChallengeAction https://techdocs.akamai.com/bot-manager/reference/put-challenge-action
		UpdateChallengeAction(ctx context.Context, params UpdateChallengeActionRequest) (map[string]interface{}, error)

		// RemoveChallengeAction https://techdocs.akamai.com/bot-manager/reference/delete-challenge-action
		RemoveChallengeAction(ctx context.Context, params RemoveChallengeActionRequest) error

		// UpdateGoogleReCaptchaSecretKey https://techdocs.akamai.com/bot-manager/reference/put-google-recaptch-secret-key
		UpdateGoogleReCaptchaSecretKey(ctx context.Context, params UpdateGoogleReCaptchaSecretKeyRequest) error
	}

	// GetChallengeActionListRequest is used to retrieve challenge actions for a configuration.
	GetChallengeActionListRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// GetChallengeActionListResponse is used to retrieve challenge actions for a configuration.
	GetChallengeActionListResponse struct {
		ChallengeActions []map[string]interface{} `json:"challengeActions"`
	}

	// GetChallengeActionRequest is used to retrieve a specific challenge action
	GetChallengeActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// CreateChallengeActionRequest is used to create a new challenge action for a specific configuration.
	CreateChallengeActionRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateChallengeActionRequest is used to update an existing challenge action
	UpdateChallengeActionRequest struct {
		ConfigID    int64
		Version     int64
		ActionID    string
		JsonPayload json.RawMessage
	}

	// RemoveChallengeActionRequest is used to remove an existing challenge action
	RemoveChallengeActionRequest struct {
		ConfigID int64
		Version  int64
		ActionID string
	}

	// UpdateGoogleReCaptchaSecretKeyRequest is used to update google reCaptcha secret key
	UpdateGoogleReCaptchaSecretKeyRequest struct {
		ConfigID  int64  `json:"-"`
		Version   int64  `json:"-"`
		ActionID  string `json:"-"`
		SecretKey string `json:"googleReCaptchaSecretKey"`
	}
)

// Validate validates a GetChallengeActionRequest.
func (v GetChallengeActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

// Validate validates a GetChallengeActionListRequest.
func (v GetChallengeActionListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateChallengeActionRequest.
func (v CreateChallengeActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateChallengeActionRequest.
func (v UpdateChallengeActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"ActionID":    validation.Validate(v.ActionID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveChallengeActionRequest.
func (v RemoveChallengeActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ActionID": validation.Validate(v.ActionID, validation.Required),
	}.Filter()
}

// Validate validates a UpdateGoogleReCaptchaSecretKeyRequest.
func (v UpdateGoogleReCaptchaSecretKeyRequest) Validate() error {
	return validation.Errors{
		"ConfigID":  validation.Validate(v.ConfigID, validation.Required),
		"Version":   validation.Validate(v.Version, validation.Required),
		"ActionID":  validation.Validate(v.ActionID, validation.Required),
		"SecretKey": validation.Validate(v.SecretKey, validation.Required),
	}.Filter()
}

func (b *botman) GetChallengeAction(ctx context.Context, params GetChallengeActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetChallengeAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChallengeAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChallengeAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetChallengeActionList(ctx context.Context, params GetChallengeActionListRequest) (*GetChallengeActionListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetChallengeActionList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetChallengeActionListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetChallengeActionList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetChallengeActionListResponse
	if params.ActionID != "" {
		for _, val := range result.ChallengeActions {
			if val["actionId"].(string) == params.ActionID {
				filteredResult.ChallengeActions = append(filteredResult.ChallengeActions, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateChallengeAction(ctx context.Context, params UpdateChallengeActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateChallengeAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateChallengeAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateChallengeAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateChallengeAction(ctx context.Context, params CreateChallengeActionRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateChallengeAction")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateChallengeAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateChallengeAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveChallengeAction(ctx context.Context, params RemoveChallengeActionRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveChallengeAction")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions/%s",
		params.ConfigID,
		params.Version,
		params.ActionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveChallengeAction request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveChallengeAction request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}

func (b *botman) UpdateGoogleReCaptchaSecretKey(ctx context.Context, params UpdateGoogleReCaptchaSecretKeyRequest) error {
	logger := b.Log(ctx)
	logger.Debug("UpdateGoogleReCaptchaSecretKey")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/response-actions/challenge-actions/%s/google-recaptcha-secret-key",
		params.ConfigID,
		params.Version,
		params.ActionID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create UpdateGoogleReCaptchaSecretKey request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params)
	if err != nil {
		return fmt.Errorf("UpdateGoogleReCaptchaSecretKey request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
