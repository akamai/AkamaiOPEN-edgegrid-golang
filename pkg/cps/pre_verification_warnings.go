package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	// PreVerification is a CPS API enabling management of pre-verification-warnings
	PreVerification interface {
		// GetChangePreVerificationWarnings gets detailed information about Domain Validation challenges
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangePreVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PreVerificationWarnings, error)

		// AcknowledgePreVerificationWarnings sends acknowledgement request to CPS informing that the warnings should be ignored
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgePreVerificationWarnings(context.Context, AcknowledgementRequest) error
	}

	// PreVerificationWarnings is a response object containing all warnings encountered during enrollment pre-verification
	PreVerificationWarnings struct {
		Warnings string `json:"warnings"`
	}
)

var (
	// ErrGetChangePreVerificationWarnings is returned when GetChangeLetsEncryptChallenges fails
	ErrGetChangePreVerificationWarnings = errors.New("fetching pre-verification-warnings")
	// ErrAcknowledgePreVerificationWarnings when AcknowledgeDVChallenges fails
	ErrAcknowledgePreVerificationWarnings = errors.New("acknowledging pre-verification-warnings")
)

func (c *cps) GetChangePreVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PreVerificationWarnings, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangePreVerificationWarnings, ErrStructValidation, err)
	}

	var rval PreVerificationWarnings

	logger := c.Log(ctx)
	logger.Debug("GetChangePreVerificationWarnings")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d/input/info/pre-verification-warnings",
		params.EnrollmentID,
		params.ChangeID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetChangePreVerificationWarnings, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangePreVerificationWarnings, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.warnings.v1+json")

	resp, err := c.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangePreVerificationWarnings, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangePreVerificationWarnings, c.Error(resp))
	}

	return &rval, nil
}

func (c *cps) AcknowledgePreVerificationWarnings(ctx context.Context, params AcknowledgementRequest) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrAcknowledgePreVerificationWarnings, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("AcknowledgePreVerificationWarnings")

	uri, err := url.Parse(fmt.Sprintf(
		"/cps/v2/enrollments/%d/changes/%d/input/update/pre-verification-warnings-ack",
		params.EnrollmentID, params.ChangeID))
	if err != nil {
		return fmt.Errorf("%w: parsing URL: %s", ErrAcknowledgePreVerificationWarnings, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrAcknowledgePreVerificationWarnings, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.acknowledgement.v1+json")

	resp, err := c.Exec(req, nil, params.Acknowledgement)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrAcknowledgePreVerificationWarnings, err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", ErrAcknowledgePreVerificationWarnings, c.Error(resp))
	}

	return nil
}
