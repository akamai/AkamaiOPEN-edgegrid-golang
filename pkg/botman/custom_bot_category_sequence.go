package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomBotCategorySequence interface supports retrieving and updating custom bot category sequence
	CustomBotCategorySequence interface {
		// GetCustomBotCategorySequence https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-sequence
		GetCustomBotCategorySequence(ctx context.Context, params GetCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error)

		// UpdateCustomBotCategorySequence https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-sequence
		UpdateCustomBotCategorySequence(ctx context.Context, params UpdateCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error)
	}

	// GetCustomBotCategorySequenceRequest is used to retrieve custom bot category sequence
	GetCustomBotCategorySequenceRequest struct {
		ConfigID int64
		Version  int64
	}

	// CustomBotCategorySequenceResponse is used to retrieve custom bot category sequence
	CustomBotCategorySequenceResponse struct {
		Sequence []string `json:"sequence"`
	}

	// UpdateCustomBotCategorySequenceRequest is used to update custom bot category sequence
	UpdateCustomBotCategorySequenceRequest struct {
		ConfigID int64    `json:"-"`
		Version  int64    `json:"-"`
		Sequence []string `json:"sequence"`
	}
)

// Validate validates a GetCustomBotCategorySequenceRequest.
func (v GetCustomBotCategorySequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomBotCategorySequenceRequest.
func (v UpdateCustomBotCategorySequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"Sequence": validation.Validate(v.Sequence, validation.Required),
	}.Filter()
}

func (b *botman) GetCustomBotCategorySequence(ctx context.Context, params GetCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomBotCategorySequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-category-sequence",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomBotCategorySequence request: %w", err)
	}

	var result CustomBotCategorySequenceResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomBotCategorySequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) UpdateCustomBotCategorySequence(ctx context.Context, params UpdateCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomBotCategorySequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-bot-category-sequence",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomBotCategorySequence request: %w", err)
	}

	var result CustomBotCategorySequenceResponse
	resp, err := b.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomBotCategorySequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}
