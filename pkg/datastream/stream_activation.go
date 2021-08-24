package datastream

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"net/url"
)

type (
	// Activation is a ds stream activations API interface
	Activation interface {
		// ActivateStream activates stream with given ID
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#putactivate
		ActivateStream(context.Context, ActivateStreamRequest) (*ActivateStreamResponse, error)

		// DeactivateStream deactivates stream with given ID
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#putdeactivate
		DeactivateStream(context.Context, DeactivateStreamRequest) (*DeactivateStreamResponse, error)

		// GetActivationHistory returns a history of activation status changes for all versions of a stream
		//
		// See: https://developer.akamai.com/api/core_features/datastream2_config/v1.html#getactivationhistory
		GetActivationHistory(context.Context, GetActivationHistoryRequest) ([]ActivationHistoryEntry, error)
	}

	// ActivateStreamResponse contains response body returned after successful stream activation
	ActivateStreamResponse struct {
		StreamVersionKey StreamVersionKey `json:"streamVersionKey"`
	}

	// StreamVersionKey contains stream details
	StreamVersionKey struct {
		StreamID        int `json:"streamId"`
		StreamVersionID int `json:"streamVersionId"`
	}

	// ActivationHistoryEntry contains single ActivationHistory item
	ActivationHistoryEntry struct {
		CreatedBy       string `json:"createdBy"`
		CreatedDate     string `json:"createdDate"`
		IsActive        bool   `json:"isActive"`
		StreamID        int    `json:"streamId"`
		StreamVersionID int    `json:"streamVersionId"`
	}

	// DeactivateStreamResponse contains response body returned after successful stream activation
	DeactivateStreamResponse ActivateStreamResponse

	// ActivateStreamRequest contains parameters necessary to send a ActivateStream request
	ActivateStreamRequest struct {
		StreamID int
	}

	// DeactivateStreamRequest contains parameters necessary to send a DeactivateStream request
	DeactivateStreamRequest ActivateStreamRequest

	// GetActivationHistoryRequest contains parameters necessary to send a GetActivationHistory request
	GetActivationHistoryRequest ActivateStreamRequest
)

// Validate performs validation on ActivateStreamRequest
func (r ActivateStreamRequest) Validate() error {
	return validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
	}.Filter()
}

// Validate performs validation on DeactivateStreamRequest
func (r DeactivateStreamRequest) Validate() error {
	return validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
	}.Filter()
}

// Validate performs validation on DeactivateStreamRequest
func (r GetActivationHistoryRequest) Validate() error {
	return validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
	}.Filter()
}

var (
	// ErrActivateStream is returned when ActivateStream fails
	ErrActivateStream = errors.New("activate stream")
	// ErrDeactivateStream is returned when DeactivateStream fails
	ErrDeactivateStream = errors.New("deactivate stream")
	// ErrGetActivationHistory is returned when DeactivateStream fails
	ErrGetActivationHistory = errors.New("view activation history")
)

func (d *ds) ActivateStream(ctx context.Context, params ActivateStreamRequest) (*ActivateStreamResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateStream, ErrStructValidation, err)
	}

	logger := d.Log(ctx)
	logger.Debug("ActivateStream")

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d/activate",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrActivateStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateStream, err)
	}

	var rval ActivateStreamResponse
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateStream, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrActivateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) DeactivateStream(ctx context.Context, params DeactivateStreamRequest) (*DeactivateStreamResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateStream, ErrStructValidation, err)
	}

	logger := d.Log(ctx)
	logger.Debug("DeactivateStream")

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d/deactivate",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrDeactivateStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateStream, err)
	}

	var rval DeactivateStreamResponse
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateStream, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) GetActivationHistory(ctx context.Context, params GetActivationHistoryRequest) ([]ActivationHistoryEntry, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetActivationHistory, ErrStructValidation, err)
	}

	logger := d.Log(ctx)
	logger.Debug("GetActivationHistory")

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d/activationHistory",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrGetActivationHistory, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetActivationHistory, err)
	}

	var rval []ActivationHistoryEntry
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetActivationHistory, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetActivationHistory, d.Error(resp))
	}

	return rval, nil
}
