package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidateBundleRequest contains request bundle parameter to validate
	ValidateBundleRequest struct {
		Bundle
	}

	// ValidateBundleResponse represents a response object returned by ValidateBundle
	ValidateBundleResponse struct {
		Errors   []ValidationIssue `json:"errors"`
		Warnings []ValidationIssue `json:"warnings"`
	}

	// ValidationIssue represents a single error or warning
	ValidationIssue struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
)

var (
	// ErrValidateBundle is returned in case an error occurs on ValidateBundle operation
	ErrValidateBundle = errors.New("validate a bundle")
)

// Validate validates ValidateBundleRequest
func (r ValidateBundleRequest) Validate() error {
	return validation.Errors{
		"Bundle.Reader": validation.Validate(r.Bundle.Reader, validation.NotNil),
	}.Filter()
}

func (e *edgeworkers) ValidateBundle(ctx context.Context, params ValidateBundleRequest) (*ValidateBundleResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ValidateBundle")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrValidateBundle, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/edgeworkers/v1/validations").
		AddHeader("Content-Type", "application/gzip").
		WithBody(params.Bundle).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrValidateBundle, err)
	}

	var result ValidateBundleResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrValidateBundle, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrValidateBundle, e.Error(resp))
	}

	return &result, nil
}
