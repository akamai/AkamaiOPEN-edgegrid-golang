package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomDefinedBot interface supports creating, retrieving, modifying and removing custom defined bots for a
	// configuration.
	CustomDefinedBot interface {
		// GetCustomDefinedBotList https://techdocs.akamai.com/bot-manager/reference/get-custom-defined-bots
		GetCustomDefinedBotList(ctx context.Context, params GetCustomDefinedBotListRequest) (*GetCustomDefinedBotListResponse, error)

		// GetCustomDefinedBot https://techdocs.akamai.com/bot-manager/reference/get-custom-defined-bot
		GetCustomDefinedBot(ctx context.Context, params GetCustomDefinedBotRequest) (map[string]interface{}, error)

		// CreateCustomDefinedBot https://techdocs.akamai.com/bot-manager/reference/post-custom-defined-bot
		CreateCustomDefinedBot(ctx context.Context, params CreateCustomDefinedBotRequest) (map[string]interface{}, error)

		// UpdateCustomDefinedBot https://techdocs.akamai.com/bot-manager/reference/put-custom-defined-bot
		UpdateCustomDefinedBot(ctx context.Context, params UpdateCustomDefinedBotRequest) (map[string]interface{}, error)

		// RemoveCustomDefinedBot https://techdocs.akamai.com/bot-manager/reference/delete-custom-defined-bot
		RemoveCustomDefinedBot(ctx context.Context, params RemoveCustomDefinedBotRequest) error
	}

	// GetCustomDefinedBotListRequest is used to retrieve the custom defined bots for a configuration.
	GetCustomDefinedBotListRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}

	// GetCustomDefinedBotListResponse is used to retrieve the custom defined bots for a configuration.
	GetCustomDefinedBotListResponse struct {
		Bots []map[string]interface{} `json:"bots"`
	}

	// GetCustomDefinedBotRequest is used to retrieve a specific custom defined bot.
	GetCustomDefinedBotRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}

	// CreateCustomDefinedBotRequest is used to create a new custom defined bot for a specific configuration.
	CreateCustomDefinedBotRequest struct {
		ConfigID    int64
		Version     int64
		JsonPayload json.RawMessage
	}

	// UpdateCustomDefinedBotRequest is used to update details for a specific custom defined bot.
	UpdateCustomDefinedBotRequest struct {
		ConfigID    int64
		Version     int64
		BotID       string
		JsonPayload json.RawMessage
	}

	// RemoveCustomDefinedBotRequest is used to remove an existing custom defined bot.
	RemoveCustomDefinedBotRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}
)

// Validate validates a GetCustomDefinedBotRequest.
func (v GetCustomDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"BotID":    validation.Validate(v.BotID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomDefinedBotsRequest.
func (v GetCustomDefinedBotListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomDefinedBotRequest.
func (v CreateCustomDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomDefinedBotRequest.
func (v UpdateCustomDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID":    validation.Validate(v.ConfigID, validation.Required),
		"Version":     validation.Validate(v.Version, validation.Required),
		"BotID":       validation.Validate(v.BotID, validation.Required),
		"JsonPayload": validation.Validate(v.JsonPayload, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomDefinedBotRequest.
func (v RemoveCustomDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"BotID":    validation.Validate(v.BotID, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomDefinedBot(ctx context.Context, params GetCustomDefinedBotRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomDefinedBot request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) GetCustomDefinedBotList(ctx context.Context, params GetCustomDefinedBotListRequest) (*GetCustomDefinedBotListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomDefinedBotList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-defined-bots",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetCustomDefinedBotListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomDefinedBotList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetCustomDefinedBotListResponse
	if params.BotID != "" {
		for _, val := range result.Bots {
			if val["botId"].(string) == params.BotID {
				filteredResult.Bots = append(filteredResult.Bots, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateCustomDefinedBot(ctx context.Context, params UpdateCustomDefinedBotRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomDefinedBot request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) CreateCustomDefinedBot(ctx context.Context, params CreateCustomDefinedBotRequest) (map[string]interface{}, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateCustomDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-defined-bots",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomDefinedBot request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result, params.JsonPayload)
	if err != nil {
		return nil, fmt.Errorf("CreateCustomDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return result, nil
}

func (b *botman) RemoveCustomDefinedBot(ctx context.Context, params RemoveCustomDefinedBotRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveCustomDefinedBot")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/custom-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveCustomDefinedBot request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveCustomDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
