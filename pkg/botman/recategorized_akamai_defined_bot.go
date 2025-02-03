package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The RecategorizedAkamaiDefinedBot interface supports creating, retrieving, modifying and removing recategorized akamai defined bot for a configuration.
	RecategorizedAkamaiDefinedBot interface {
		// GetRecategorizedAkamaiDefinedBotList https://techdocs.akamai.com/bot-manager/reference/get-recategorized-akamai-defined-bots
		GetRecategorizedAkamaiDefinedBotList(ctx context.Context, params GetRecategorizedAkamaiDefinedBotListRequest) (*GetRecategorizedAkamaiDefinedBotListResponse, error)

		// GetRecategorizedAkamaiDefinedBot https://techdocs.akamai.com/bot-manager/reference/get-recategorized-akamai-defined-bot
		GetRecategorizedAkamaiDefinedBot(ctx context.Context, params GetRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error)

		// CreateRecategorizedAkamaiDefinedBot https://techdocs.akamai.com/bot-manager/reference/post-recategorized-akamai-defined-bot
		CreateRecategorizedAkamaiDefinedBot(ctx context.Context, params CreateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error)

		// UpdateRecategorizedAkamaiDefinedBot https://techdocs.akamai.com/bot-manager/reference/put-recategorized-akamai-defined-bot
		UpdateRecategorizedAkamaiDefinedBot(ctx context.Context, params UpdateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error)

		// RemoveRecategorizedAkamaiDefinedBot https://techdocs.akamai.com/bot-manager/reference/delete-recategorized-akamai-defined-bot
		RemoveRecategorizedAkamaiDefinedBot(ctx context.Context, params RemoveRecategorizedAkamaiDefinedBotRequest) error
	}

	// GetRecategorizedAkamaiDefinedBotListRequest is used to retrieve the recategorized akamai defined bots for a configuration.
	GetRecategorizedAkamaiDefinedBotListRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}

	// GetRecategorizedAkamaiDefinedBotRequest is used to retrieve a specific recategorized akamai defined bot
	GetRecategorizedAkamaiDefinedBotRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}

	// RecategorizedAkamaiDefinedBotResponse is used to retrieve a specific recategorized akamai defined bot
	RecategorizedAkamaiDefinedBotResponse struct {
		BotID      string `json:"botId"`
		CategoryID string `json:"customBotCategoryId"`
	}

	// GetRecategorizedAkamaiDefinedBotListResponse is used to retrieve the recategorized akamai defined bots for a configuration.
	GetRecategorizedAkamaiDefinedBotListResponse struct {
		Bots []RecategorizedAkamaiDefinedBotResponse `json:"recategorizedBots"`
	}

	// CreateRecategorizedAkamaiDefinedBotRequest is used to create a new recategorized akamai defined bot for a specific configuration.
	CreateRecategorizedAkamaiDefinedBotRequest struct {
		ConfigID   int64  `json:"-"`
		Version    int64  `json:"-"`
		BotID      string `json:"botId"`
		CategoryID string `json:"customBotCategoryId"`
	}

	// UpdateRecategorizedAkamaiDefinedBotRequest is used to update details for a specific recategorized akamai defined bot
	UpdateRecategorizedAkamaiDefinedBotRequest struct {
		ConfigID   int64  `json:"-"`
		Version    int64  `json:"-"`
		BotID      string `json:"botId"`
		CategoryID string `json:"customBotCategoryId"`
	}

	// RemoveRecategorizedAkamaiDefinedBotRequest is used to remove an existing recategorized akamai defined bot
	RemoveRecategorizedAkamaiDefinedBotRequest struct {
		ConfigID int64
		Version  int64
		BotID    string
	}
)

// Validate validates a GetRecategorizedAkamaiDefinedBotRequest.
func (v GetRecategorizedAkamaiDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"BotID":    validation.Validate(v.BotID, validation.Required),
	}.Filter()
}

// Validate validates a GetRecategorizedAkamaiDefinedBotsRequest.
func (v GetRecategorizedAkamaiDefinedBotListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateRecategorizedAkamaiDefinedBotRequest.
func (v CreateRecategorizedAkamaiDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"BotID":      validation.Validate(v.BotID, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateRecategorizedAkamaiDefinedBotRequest.
func (v UpdateRecategorizedAkamaiDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"BotID":      validation.Validate(v.BotID, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveRecategorizedAkamaiDefinedBotRequest.
func (v RemoveRecategorizedAkamaiDefinedBotRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"BotID":    validation.Validate(v.BotID, validation.Required),
	}.Filter()
}

func (b *botman) GetRecategorizedAkamaiDefinedBot(ctx context.Context, params GetRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetRecategorizedAkamaiDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/recategorized-akamai-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecategorizedAkamaiDefinedBot request: %w", err)
	}

	var result RecategorizedAkamaiDefinedBotResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecategorizedAkamaiDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) GetRecategorizedAkamaiDefinedBotList(ctx context.Context, params GetRecategorizedAkamaiDefinedBotListRequest) (*GetRecategorizedAkamaiDefinedBotListResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetRecategorizedAkamaiDefinedBotList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/recategorized-akamai-defined-bots",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetlustomDenyList request: %w", err)
	}

	var result GetRecategorizedAkamaiDefinedBotListResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetRecategorizedAkamaiDefinedBotList request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetRecategorizedAkamaiDefinedBotListResponse
	if params.BotID != "" {
		for _, val := range result.Bots {
			if val.BotID == params.BotID {
				filteredResult.Bots = append(filteredResult.Bots, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}

func (b *botman) UpdateRecategorizedAkamaiDefinedBot(ctx context.Context, params UpdateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateRecategorizedAkamaiDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/recategorized-akamai-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateRecategorizedAkamaiDefinedBot request: %w", err)
	}

	var result RecategorizedAkamaiDefinedBotResponse
	resp, err := b.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateRecategorizedAkamaiDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) CreateRecategorizedAkamaiDefinedBot(ctx context.Context, params CreateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("CreateRecategorizedAkamaiDefinedBot")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/recategorized-akamai-defined-bots",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateRecategorizedAkamaiDefinedBot request: %w", err)
	}

	var result RecategorizedAkamaiDefinedBotResponse
	resp, err := b.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("CreateRecategorizedAkamaiDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) RemoveRecategorizedAkamaiDefinedBot(ctx context.Context, params RemoveRecategorizedAkamaiDefinedBotRequest) error {
	logger := b.Log(ctx)
	logger.Debug("RemoveRecategorizedAkamaiDefinedBot")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/recategorized-akamai-defined-bots/%s",
		params.ConfigID,
		params.Version,
		params.BotID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create RemoveRecategorizedAkamaiDefinedBot request: %w", err)
	}

	var result map[string]interface{}
	resp, err := b.Exec(req, &result)
	if err != nil {
		return fmt.Errorf("RemoveRecategorizedAkamaiDefinedBot request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return b.Error(resp)
	}

	return nil
}
