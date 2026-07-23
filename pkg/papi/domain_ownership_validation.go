package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidateDomainsOwnershipRequest contains parameters required to initiate domains ownership validation
	ValidateDomainsOwnershipRequest struct {
		RefreshToken bool
		Body         ValidateDomainsOwnershipRequestBody
	}

	// ValidateDomainsOwnershipRequestBody contains the list of hostnames to validate
	ValidateDomainsOwnershipRequestBody struct {
		Hostnames []string `json:"hostnames"`
	}

	// ValidateDomainsOwnershipResponse contains the response from domains ownership validation
	ValidateDomainsOwnershipResponse struct {
		AccountID string                      `json:"accountId"`
		Hostnames []HostnameValidationDetails `json:"hostnames"`
	}

	// HostnameValidationDetails contains validation details for a specific hostname
	HostnameValidationDetails struct {
		Hostname                 string           `json:"hostname"`
		DomainValidationStatus   string           `json:"domainValidationStatus"`
		ValidationScope          *string          `json:"validationScope"`
		ChallengeTokenExpiryDate *time.Time       `json:"challengeTokenExpiryDate"`
		ValidationCname          *ValidationCname `json:"validationCname"`
		ValidationTXT            *ValidationTXT   `json:"validationTxt"`
		ValidationHTTP           *ValidationHTTP  `json:"validationHttp"`
	}
)

var (
	// ErrValidateDomainsOwnership represents error when validating domains ownership fails
	ErrValidateDomainsOwnership = errors.New("validating domains ownership")
)

// Validate validates ValidateDomainsOwnershipRequest
func (r ValidateDomainsOwnershipRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Hostnames": validation.Validate(r.Body.Hostnames, validation.Required, validation.Each(validation.Required)),
	})
}

// ValidateDomainsOwnership initiates ownership validation of the specified domains
func (p *papi) ValidateDomainsOwnership(ctx context.Context, params ValidateDomainsOwnershipRequest) (*ValidateDomainsOwnershipResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ValidateDomainsOwnership")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrValidateDomainsOwnership, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/papi/v1/domain-challenges").
		AddQueryParam("refreshToken", strconv.FormatBool(params.RefreshToken)).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrValidateDomainsOwnership, err)
	}

	var result ValidateDomainsOwnershipResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrValidateDomainsOwnership, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrValidateDomainsOwnership, p.Error(resp))
	}

	return &result, nil
}
