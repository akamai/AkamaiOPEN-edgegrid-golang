package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// PostVerificationWarnings is a response object containing all warnings encountered during enrollment post-verification
	PostVerificationWarnings struct {
		Warnings string `json:"warnings"`
	}
)

var (
	// ErrGetChangePostVerificationWarnings is returned when GetChangePostVerificationWarnings fails
	ErrGetChangePostVerificationWarnings = errors.New("get post-verification-warnings")
	// ErrAcknowledgePostVerificationWarnings is returned when AcknowledgePostVerificationWarnings fails
	ErrAcknowledgePostVerificationWarnings = errors.New("acknowledging post-verification-warnings")
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangePostVerificationWarnings, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) AcknowledgePostVerificationWarnings(ctx context.Context, params AcknowledgementRequest) error {
	c.Log(ctx).Debug("AcknowledgePostVerificationWarnings")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrAcknowledgePostVerificationWarnings, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/update/post-verification-warnings-ack",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrAcknowledgePostVerificationWarnings, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.acknowledgement.v1+json; charset=utf-8")

	resp, err := c.Exec(req, nil, params.Acknowledgement)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrAcknowledgePostVerificationWarnings, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", ErrAcknowledgePostVerificationWarnings, c.Error(resp))
	}

	return nil
}
