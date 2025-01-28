package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomBotCategoryItemSequence interface supports retrieving and updating custom bot category item sequence
	CustomBotCategoryItemSequence interface {
		// GetCustomBotCategoryItemSequence https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-item-sequence
		GetCustomBotCategoryItemSequence(ctx context.Context, params GetCustomBotCategoryItemSequenceRequest) (*GetCustomBotCategoryItemSequenceResponse, error)
		// UpdateCustomBotCategoryItemSequence https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-item-sequence
		UpdateCustomBotCategoryItemSequence(ctx context.Context, params UpdateCustomBotCategoryItemSequenceRequest) (*UpdateCustomBotCategoryItemSequenceResponse, error)
	}

	// GetCustomBotCategoryItemSequenceRequest is used to retrieve custom bot category sequence
	GetCustomBotCategoryItemSequenceRequest struct {
		ConfigID   int64
		Version    int64
		CategoryID string
	}

	// GetCustomBotCategoryItemSequenceResponse contains the sequence of botIds
	GetCustomBotCategoryItemSequenceResponse UUIDSequence

	// UpdateCustomBotCategoryItemSequenceResponse contains the sequence of botIds
	UpdateCustomBotCategoryItemSequenceResponse UUIDSequence

	// UpdateCustomBotCategoryItemSequenceRequest is used to update custom bot category item sequence
	UpdateCustomBotCategoryItemSequenceRequest struct {
		ConfigID   int64
		Version    int64
		CategoryID string
		Sequence   UUIDSequence
	}

	// UUIDSequence is nothing more than a sequence of UUIDs, in this case bots, but this could be reused.
	UUIDSequence struct {
		Sequence []string `json:"sequence"`
	}
)

// Validate validates a GetCustomBotCategoryItemSequenceRequest.
func (v GetCustomBotCategoryItemSequenceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
	})
}

// Validate validates an UpdateCustomBotCategoryItemSequenceRequest.
func (v UpdateCustomBotCategoryItemSequenceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID":   validation.Validate(v.ConfigID, validation.Required),
		"Version":    validation.Validate(v.Version, validation.Required),
		"CategoryID": validation.Validate(v.CategoryID, validation.Required),
		"Sequence":   validation.Validate(v.Sequence.Sequence, validation.Required),
	})
}

func (b *botman) GetCustomBotCategoryItemSequence(ctx context.Context, params GetCustomBotCategoryItemSequenceRequest) (*GetCustomBotCategoryItemSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategoryItemSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories/%s/custom-bot-category-item-sequence",
		params.ConfigID,
		params.Version,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomBotCategoryItemSequence request: %w", err)
	}

	var result GetCustomBotCategoryItemSequenceResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategoryItemSequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) UpdateCustomBotCategoryItemSequence(ctx context.Context, params UpdateCustomBotCategoryItemSequenceRequest) (*UpdateCustomBotCategoryItemSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomBotCategoryItemSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-categories/%s/custom-bot-category-item-sequence",
		params.ConfigID,
		params.Version,
		params.CategoryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomBotCategoryItemSequence request: %w", err)
	}

	var result UpdateCustomBotCategoryItemSequenceResponse
	resp, err := b.Exec(req, &result, params.Sequence)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomBotCategoryItemSequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}
