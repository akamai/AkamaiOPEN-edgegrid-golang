package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetAccessKeyVersionStatusResponse contains response from GetAccessKeyVersionStatus
	GetAccessKeyVersionStatusResponse struct {
		AccessKeyVersion *KeyVersion    `json:"accessKeyVersion"`
		ProcessingStatus ProcessingType `json:"processingStatus"`
		RequestDate      string         `json:"requestDate"`
		RequestedBy      string         `json:"requestedBy"`
	}

	// GetAccessKeyVersionStatusRequest holds parameters for GetAccessKeyVersionStatus
	GetAccessKeyVersionStatusRequest struct {
		RequestID int64
	}
)

// Validate validates GetAccessKeyVersionStatusRequest
func (r GetAccessKeyVersionStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

var (
	// ErrGetAccessKeyVersionStatus is returned when GetAccessKeyVersionStatus fails
	ErrGetAccessKeyVersionStatus = errors.New("get the status of an access key version")
)

func (c *cloudaccess) GetAccessKeyVersionStatus(ctx context.Context, params GetAccessKeyVersionStatusRequest) (*GetAccessKeyVersionStatusResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKeyStatusVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKeyVersionStatus, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/cam/v1/access-key-version-create-requests/%d", params.RequestID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKeyStatusVersion request: %w", err)
	}

	var result GetAccessKeyVersionStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAccessKeyStatusVersion request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}
