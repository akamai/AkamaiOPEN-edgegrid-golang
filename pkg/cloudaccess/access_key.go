package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetAccessKeyStatusResponse contains response from GetAccessKeyStatus
	GetAccessKeyStatusResponse struct {
		AccessKey        *KeyLink            `json:"accessKey"`
		AccessKeyVersion *KeyVersion         `json:"accessKeyVersion"`
		ProcessingStatus ProcessingType      `json:"processingStatus"`
		Request          *RequestInformation `json:"request"`
		RequestDate      string              `json:"requestDate"`
		RequestID        int64               `json:"requestId"`
		RequestedBy      string              `json:"requestedBy"`
	}

	// GetAccessKeyStatusRequest holds parameters for GetAccessKeyStatus
	GetAccessKeyStatusRequest struct {
		RequestID int64
	}
)

// Validate validates GetAccessKeyStatusRequest
func (r GetAccessKeyStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

var (
	// ErrGetAccessKeyStatus is returned when GetAccessKeyStatus fails
	ErrGetAccessKeyStatus = errors.New("get the status of an access key")
)

func (c *cloudaccess) GetAccessKeyStatus(ctx context.Context, params GetAccessKeyStatusRequest) (*GetAccessKeyStatusResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKeyStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKeyStatus, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/cam/v1/access-key-create-requests/%d", params.RequestID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKeyStatus request: %w", err)
	}

	var result GetAccessKeyStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAccessKeyStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}
