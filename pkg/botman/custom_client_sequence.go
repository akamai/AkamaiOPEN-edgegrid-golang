package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomClientSequence interface supports retrieving and updating the custom client sequence for a configuration
	CustomClientSequence interface {
		// GetCustomClientSequence is used to retrieve the custom client sequence for a config version
		// See https://techdocs.akamai.com/bot-manager/reference/get-custom-client-sequence
		GetCustomClientSequence(ctx context.Context, params GetCustomClientSequenceRequest) (*CustomClientSequenceResponse, error)

		// UpdateCustomClientSequence is used to update the existing custom client sequence for a config version
		// See https://techdocs.akamai.com/bot-manager/reference/put-custom-client-sequence
		UpdateCustomClientSequence(ctx context.Context, params UpdateCustomClientSequenceRequest) (*CustomClientSequenceResponse, error)
	}

	// GetCustomClientSequenceRequest is used to retrieve custom client sequence
	GetCustomClientSequenceRequest struct {
		ConfigID int64
		Version  int64
	}

	// UpdateCustomClientSequenceRequest is used to modify custom client sequence
	UpdateCustomClientSequenceRequest struct {
		ConfigID int64    `json:"-"`
		Version  int64    `json:"-"`
		Sequence []string `json:"sequence"`
	}

	// CustomClientSequenceResponse is used to represent custom client sequence
	CustomClientSequenceResponse struct {
		Sequence   []string           `json:"sequence"`
		Validation ValidationResponse `json:"validation"`
	}
)

// Validate validates a GetCustomClientSequenceRequest.
func (v GetCustomClientSequenceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates an UpdateCustomClientSequenceRequest.
func (v UpdateCustomClientSequenceRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"Sequence": validation.Validate(v.Sequence, validation.Required),
	})
}

func (b *botman) GetCustomClientSequence(ctx context.Context, params GetCustomClientSequenceRequest) (*CustomClientSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetCustomClientSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-client-sequence",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomClientSequence request: %w", err)
	}

	var result CustomClientSequenceResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetCustomClientSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}

func (b *botman) UpdateCustomClientSequence(ctx context.Context, params UpdateCustomClientSequenceRequest) (*CustomClientSequenceResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("UpdateCustomClientSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-client-sequence",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomClientSequence request: %w", err)
	}

	var result CustomClientSequenceResponse
	resp, err := b.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomClientSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	return &result, nil
}
