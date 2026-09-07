package domainownership

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidateDomain represents the structure for validating a domain.
	ValidateDomain struct {
		// DomainName is the name of the domain to search for.
		DomainName string `json:"domainName"`

		// ValidationMethod for instant validation of a domain, either DNS_CNAME, DNS_TXT, or HTTP.
		ValidationMethod ValidationMethod `json:"validationMethod"`

		// ValidationScope indicates the scope of the validation, either HOST, DOMAIN, or WILDCARD.
		ValidationScope ValidationScope `json:"validationScope"`
	}

	// ValidateDomainResponse represents the response structure for Validate and Invalidate Domain.
	ValidateDomainResponse struct {
		// DomainName is the name of the domain.
		DomainName string `json:"domainName"`

		// DomainStatus is the status of the domain, either REQUEST_ACCEPTED, VALIDATION_IN_PROGRESS, or VALIDATED.
		DomainStatus string `json:"domainStatus"`

		// ValidationScope indicates the scope of the validation, either HOST, DOMAIN, or WILDCARD.
		ValidationScope string `json:"validationScope"`
	}

	// ValidateDomainsRequest represents the request structure for ValidateDomains.
	ValidateDomainsRequest struct {
		// Domains is a list of domains to validate.
		Domains []ValidateDomain `json:"domains"`
	}

	// ValidateDomainsResponse represents the response structure for ValidateDomains.
	ValidateDomainsResponse struct {
		// Domains contains the list of validated domains.
		Domains []ValidateDomainResponse `json:"domains"`
	}

	// InvalidateDomainRequest represents the request structure for InvalidateDomain.
	InvalidateDomainRequest Domain

	// InvalidateDomainResponse represents the response structure for InvalidateDomain.
	InvalidateDomainResponse ValidateDomainResponse

	// InvalidateDomainsRequest represents the request structure for InvalidateDomains.
	InvalidateDomainsRequest struct {
		// Domains is a list of domains to invalidate.
		Domains []Domain `json:"domains"`
	}

	// InvalidateDomainsResponse represents the response structure for InvalidateDomains.
	InvalidateDomainsResponse struct {
		// Domains contains the list of invalidated domains.
		Domains []InvalidateDomainResponse `json:"domains"`
	}
)

var (
	// ErrInvalidateDomain is returned when there is an error invalidating a domain.
	ErrInvalidateDomain = errors.New("invalidate domain")

	// ErrInvalidateDomains is returned when there is an error invalidating domains.
	ErrInvalidateDomains = errors.New("invalidate domains")

	// ErrValidateDomains is returned when there is an error validating domains.
	ErrValidateDomains = errors.New("validate domains")
)

// Validate validates the InvalidateDomainsRequest parameters.
func (r InvalidateDomainsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Domains": validation.Validate(r.Domains, validation.Required, validation.Length(1, 0)),
	})
}

// Validate validates the ValidateDomainsRequest parameters.
func (r ValidateDomainsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Domains": validation.Validate(r.Domains, validation.Required, validation.Length(1, 0)),
	})
}

// Validate validates the ValidateDomain parameters.
func (d ValidateDomain) Validate() error {
	return validation.Errors{
		"DomainName":       domainNameValidation(d.DomainName),
		"ValidationScope":  scopeValidation(d.ValidationScope),
		"ValidationMethod": validateValidationMethod(d.ValidationMethod),
	}.Filter()
}

// Validate validates the InvalidateDomainRequest parameters.
func (r InvalidateDomainRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":      domainNameValidation(r.DomainName),
		"ValidationScope": scopeValidation(r.ValidationScope),
	})
}

func (d *domainownership) ValidateDomains(ctx context.Context, params ValidateDomainsRequest) (*ValidateDomainsResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("ValidateDomains")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrValidateDomains, ErrStructValidation, err)
	}

	uri, err := url.Parse("/domain-validation/v1/domains/validate-now")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrValidateDomains, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrValidateDomains, err)
	}

	var result ValidateDomainsResponse
	resp, err := d.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrValidateDomains, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrValidateDomains, d.Error(resp))
	}

	return &result, nil
}

func (d *domainownership) InvalidateDomain(ctx context.Context, params InvalidateDomainRequest) (*InvalidateDomainResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("InvalidateDomain")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrInvalidateDomain, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/domain-validation/v1/domains/invalidate/%s", params.DomainName))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrInvalidateDomain, err)
	}

	q := uri.Query()
	q.Add("validationScope", string(params.ValidationScope))

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrInvalidateDomain, err)
	}

	var result InvalidateDomainResponse
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrInvalidateDomain, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrInvalidateDomain, d.Error(resp))
	}

	return &result, nil
}

func (d *domainownership) InvalidateDomains(ctx context.Context, params InvalidateDomainsRequest) (*InvalidateDomainsResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("InvalidateDomains")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrInvalidateDomains, ErrStructValidation, err)
	}

	uri, err := url.Parse("/domain-validation/v1/domains/invalidate")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrInvalidateDomains, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrInvalidateDomains, err)
	}

	var result InvalidateDomainsResponse
	resp, err := d.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrInvalidateDomains, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrInvalidateDomains, d.Error(resp))
	}

	return &result, nil
}
