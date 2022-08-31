package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// PostVerification is a CPS API enabling management of post-verification-warnings
	PostVerification interface {
		// GetChangePostVerificationWarnings gets information about post verification warnings
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangePostVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PostVerificationWarnings, error)
	}

	// PostVerificationWarnings is a response object containing all warnings encountered during enrollment post-verification
	PostVerificationWarnings struct {
		Warnings string `json:"warnings"`
	}
)

var (
	// ErrGetChangePostVerificationWarnings is returned when GetChangePostVerificationWarnings fails
	ErrGetChangePostVerificationWarnings = errors.New("get post-verification-warnings")
)

func (c *cps) GetChangePostVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PostVerificationWarnings, error) {
	c.Log(ctx).Debug("GetChangePostVerificationWarnings")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangePostVerificationWarnings, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/info/post-verification-warnings",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangePostVerificationWarnings, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.warnings.v1+json")

	var result PostVerificationWarnings
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangePostVerificationWarnings, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangePostVerificationWarnings, c.Error(resp))
	}

	return &result, nil
}
