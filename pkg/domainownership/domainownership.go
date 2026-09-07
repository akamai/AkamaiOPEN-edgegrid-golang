package domainownership

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// DomainOwnership is the interface for the Domain Ownership Manager that is used for Domain Validation.
	DomainOwnership interface {
		// AddDomains adds domains to validate.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/post-domains
		AddDomains(ctx context.Context, params AddDomainsRequest) (*AddDomainsResponse, error)

		// DeleteDomain deletes a domain.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/delete-domain
		DeleteDomain(ctx context.Context, params DeleteDomainRequest) error

		// DeleteDomains deletes a batch of domains.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/delete-domains
		DeleteDomains(ctx context.Context, params DeleteDomainsRequest) error

		// InvalidateDomain invalidates a domain.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/post-invalidate-domain
		InvalidateDomain(ctx context.Context, params InvalidateDomainRequest) (*InvalidateDomainResponse, error)

		// InvalidateDomains invalidates domains.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/post-invalidate-domains
		InvalidateDomains(ctx context.Context, params InvalidateDomainsRequest) (*InvalidateDomainsResponse, error)

		// ValidateDomains validates list of domains.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/post-validate-domains
		ValidateDomains(ctx context.Context, params ValidateDomainsRequest) (*ValidateDomainsResponse, error)

		// ListDomains returns the list of available domains.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/get-domains
		ListDomains(ctx context.Context, params ListDomainsRequest) (*ListDomainsResponse, error)

		// GetDomain gets a specific domain.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/get-domain
		GetDomain(ctx context.Context, params GetDomainRequest) (*GetDomainResponse, error)

		// SearchDomains returns the status of specified domains. For any nonexistent domains, it returns the closest matching domain status.
		//
		// See: https://techdocs.akamai.com/domain-validation/reference/post-search-domains
		SearchDomains(ctx context.Context, params SearchDomainsRequest) (*SearchDomainsResponse, error)
	}

	domainownership struct {
		session.Session
	}

	// Option is a function that configures the Domain Ownership.
	Option func(*domainownership)
)

// Client creates a new DomainOwnership client.
func Client(sess session.Session, opts ...Option) DomainOwnership {
	c := &domainownership{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
