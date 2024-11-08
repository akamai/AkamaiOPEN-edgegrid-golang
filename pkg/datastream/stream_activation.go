package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ActivationHistoryEntry contains single ActivationHistory item
	ActivationHistoryEntry struct {
		ModifiedBy    string       `json:"modifiedBy"`
		ModifiedDate  string       `json:"modifiedDate"`
		Status        StreamStatus `json:"status"`
		StreamID      int64        `json:"streamId"`
		StreamVersion int64        `json:"streamVersion"`
	}

	// ActivateStreamRequest contains parameters necessary to send a ActivateStream request
	ActivateStreamRequest struct {
		StreamID int64
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

func (d *ds) ActivateStream(ctx context.Context, params ActivateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("ActivateStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/streams/%d/activate",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrActivateStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateStream, err)
	}

	var rval DetailedStreamVersion
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateStream, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrActivateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) DeactivateStream(ctx context.Context, params DeactivateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("DeactivateStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/streams/%d/deactivate",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrDeactivateStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateStream, err)
	}

	var rval DetailedStreamVersion
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateStream, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) GetActivationHistory(ctx context.Context, params GetActivationHistoryRequest) ([]ActivationHistoryEntry, error) {
	logger := d.Log(ctx)
	logger.Debug("GetActivationHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetActivationHistory, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/streams/%d/activation-history",
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetActivationHistory, d.Error(resp))
	}

	return rval, nil
}
