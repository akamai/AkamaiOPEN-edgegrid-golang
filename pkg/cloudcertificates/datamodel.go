package cloudcertificates

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/texts"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	_ validation.Validatable = CreateCertificateRequest{}
	_ validation.Validatable = GetCertificateRequest{}
	_ validation.Validatable = UpdateCertificateRequest{}
	_ validation.Validatable = PatchCertificateRequest{}
	_ validation.Validatable = ListCertificatesRequest{}
	_ validation.Validatable = DeleteCertificateRequest{}
	_ validation.Validatable = ListCertificateBindingsRequest{}
	_ validation.Validatable = ListBindingsRequest{}
	_ validation.Validatable = Subject{}
	_ validation.Validatable = CryptographicAlgorithm("")
	_ validation.Validatable = KeySize("")
	_ validation.Validatable = SecureNetwork("")
)

const (
	// StatusActive indicates the certificate is bound to one or more property hostnames and is `ACTIVE` on one or more networks.
	StatusActive CertificateStatus = "ACTIVE"

	// StatusReadyForUse indicates the accepted signed certificate (and/or trust chain) for one or both key types and is ready for usage.
	StatusReadyForUse CertificateStatus = "READY_FOR_USE"

	// StatusCSRReady indicates the CSR generation is complete and are available for download.
	StatusCSRReady CertificateStatus = "CSR_READY"

	// CryptographicAlgorithmRSA indicates the `RSA` algorithm.
	CryptographicAlgorithmRSA CryptographicAlgorithm = "RSA"

	// CryptographicAlgorithmECDSA indicates the `ECDSA` algorithm.
	CryptographicAlgorithmECDSA CryptographicAlgorithm = "ECDSA"

	// SecureNetworkEnhancedTLS represents the `ENHANCED_TLS` secure network type.
	SecureNetworkEnhancedTLS SecureNetwork = "ENHANCED_TLS"

	// SecureNetworkStandardTLS represents the `STANDARD_TLS` secure network type.
	SecureNetworkStandardTLS SecureNetwork = "STANDARD_TLS"

	// KeySize2048 represents a key size of 2048 bits.
	KeySize2048 KeySize = "2048"

	// KeySizeP256 represents a key size of 256 bits.
	KeySizeP256 KeySize = "P-256"

	// KeySizeP384 represents a key size of 384 bits.
	KeySizeP384 KeySize = "P-384"

	// NetworkStaging represents staging network.
	NetworkStaging Network = "STAGING"

	// NetworkProduction represents production network.
	NetworkProduction Network = "PRODUCTION"

	// SortFieldPat is a regex pattern for validating sort fields.
	SortFieldPat = `[+\-]?(modifiedDate|expirationDate|createdDate|certificateName)`
)

type (
	// CreateCertificateRequest represents the parameters for creating a new certificate.
	CreateCertificateRequest struct {
		// ContractID under which this certificate will be created.
		ContractID string

		// GroupID is a unique identifier for the group. This is a required parameter when creating a new certificate.
		GroupID string

		// Body of the create certificate request containing certificate details.
		Body CreateCertificateRequestBody
	}

	// CreateCertificateResponse contains the response for a certificate creation request.
	CreateCertificateResponse struct {
		// Certificate contains the created certificate details.
		Certificate Certificate

		// ResourceLimits contains information about certificate limits.
		ResourceLimits ResourceLimitsMetadata
		RateLimits     RateLimitsMetadata
	}

	// GetCertificateRequest represents the parameters for retrieving a certificate.
	GetCertificateRequest struct {
		// Certificate identifier on which to perform the desired operation.
		CertificateID string
	}

	// GetCertificateResponse contains the response for retrieving a certificate.
	GetCertificateResponse struct {
		Certificate Certificate
		RateLimits  RateLimitsMetadata
	}

	// PatchCertificateRequest represents the parameters for patching a certificate.
	PatchCertificateRequest struct {
		// Certificate identifier on which to perform the desired operation.
		CertificateID string

		// Signed certificate PEM content.
		SignedCertificatePEM string

		// PEM content of trust chain.
		TrustChainPEM string

		// Name of the certificate. To reset the name to a default value, provide an empty string.
		CertificateName *string

		// Parameter to acknowledge warnings and retry certificate upload when the returned response
		// contains warnings for the uploaded certificate.
		AcknowledgeWarnings bool
	}

	// PatchCertificateResponse contains the response for patching a certificate.
	PatchCertificateResponse struct {
		Certificate Certificate
		RateLimits  RateLimitsMetadata
	}

	// UpdateCertificateRequest represents the parameters for updating a certificate.
	UpdateCertificateRequest struct {
		// Certificate identifier on which to perform the desired operation.
		CertificateID string `json:"-"`

		// Signed certificate PEM content.
		SignedCertificatePEM string `json:"signedCertificatePem,omitempty"`

		// PEM content of trust chain.
		TrustChainPEM string `json:"trustChainPem,omitempty"`

		// Name of the certificate. To reset the name to a default value, provide an empty string.
		CertificateName *string `json:"certificateName,omitempty"`

		// Parameter to acknowledge warnings and retry certificate upload when the returned response
		// contains warnings for the uploaded certificate.
		AcknowledgeWarnings bool `json:"-"`
	}

	// UpdateCertificateResponse contains the response for updating a certificate.
	UpdateCertificateResponse struct {
		Certificate Certificate
		RateLimits  RateLimitsMetadata
	}

	// DeleteCertificateRequest represents the parameters for deleting a certificate.
	DeleteCertificateRequest struct {
		// Certificate identifier on which to perform the desired operation.
		CertificateID string
	}

	// DeleteCertificateResponse contains the response for deleting a certificate.
	DeleteCertificateResponse = RateLimitsMetadata

	// ListCertificateBindingsRequest represents the parameters for retrieving certificate bindings.
	ListCertificateBindingsRequest struct {
		// CertificateID is an identifier on which to perform the desired operation.
		CertificateID string

		// PageSize is the number of items to return in the response. The default is 10, and the maximum is 100.
		PageSize int64

		// Page is the specific page from the result set. 1 by default.
		Page int64
	}

	// ListCertificateBindingsResponse contains the bindings for a certificate.
	ListCertificateBindingsResponse struct {
		// Bindings is a list of bindings filtered or paginated per request.
		Bindings []CertificateBinding `json:"bindings"`

		// Links contains pagination and navigation links.
		Links Links `json:"links"`

		RateLimits RateLimitsMetadata
	}

	// ListCertificatesRequest represents the parameters for listing certificates.
	ListCertificatesRequest struct {
		// Filter by contract identifier and only return CloudCertificates certificates associated with that contract.
		ContractID string

		// Unique identifier for the group.
		GroupID string

		// Filter by status of certificate. Can accept comma separated value list and will return all certificates in the given status list.
		// Valid values: `ACTIVE`, `READY_FOR_USE`, `CSR_READY`.
		CertificateStatus []CertificateStatus

		// If provided, return certificates where the expiration date is within the provided number of days from the date of request.
		// A value of 0 returns only expired certificates.
		ExpiringInDays *int64

		// If provided, search across SAN domains and Subject CN for the provided domain exact match or wildcard match. Comparison is case-insensitive.
		Domain string

		// Substring match on certificate name.
		CertificateName string

		// Key type to filter certificates. Either `RSA` or `ECDSA`.
		KeyType CryptographicAlgorithm

		// Substring match on the signed certificate issuer field to be able to quickly find certificates issued by certain intermediate or root CA.
		Issuer string

		// If true, returns full materials for each certificate (CSR, signed certificate, trust chain).
		IncludeCertificateMaterials bool

		// The number of items to return in the response. The default is 10, and the maximum is 100.
		// If you specify a value greater than 100, the API will return an error.
		PageSize int64

		// Specific page from the result set. 1 by default.
		Page int64

		// Certificate fields prefixed with optional `+` or `-` ( `+` to indicate ascending and `-` for descending order) of the fields.
		// Supported fields:
		// - modifiedDate
		// - expirationDate
		// - createdDate
		// - certificateName
		//
		//  If none provided, results are sorted by default in descending order of modifiedDate.
		//  The order (from left to right) of the fields dictates the order in which sorting logic is applied on the query results.
		Sort string
	}

	// ListCertificatesResponse contains a list of certificates and metadata.
	ListCertificatesResponse struct {
		// Certificates filtered or paginated per request.
		Certificates []Certificate `json:"certificates"`

		// Links contains pagination and navigation links.
		Links Links `json:"links"`

		// Metadata contains metadata about the list of certificates.
		Metadata ListMetadata `json:"metadata"`

		// RateLimits contains information about API rate limits.
		RateLimits RateLimitsMetadata
	}

	// ListBindingsRequest represents the parameters for listing all certificate bindings.
	ListBindingsRequest struct {
		// ContractID filters results to certificates created under the specified contract.
		ContractID string

		// GroupID filters results to certificates under the group specified by this unique id. This is an optional parameter when listing certificate bindings.
		GroupID string

		// Domain filters results to certificates that include the specified domain as a SAN or subject CN.
		// Matches are case-insensitive, and support wildcards.
		Domain string

		// Network filters results to bindings on the specified network. Valid values are `STAGING` and `PRODUCTION`.
		Network Network

		// ExpiringInDays, if provided, filters results to certificates where the expiration date is within the provided number of days from the date of request.
		// A value of 0 returns only expired certificates.
		ExpiringInDays *int64

		// PageSize is the number of items to return in the response. The default is 10, and the maximum is 100.
		PageSize int64

		// Page is the specific page from the result set. 1 by default.
		Page int64
	}

	// ListBindingsResponse contains a list of certificate bindings.
	ListBindingsResponse struct {
		// Bindings is a list of bindings filtered or paginated per request.
		Bindings []CertificateBinding `json:"bindings"`

		// Links contains pagination and navigation links.
		Links Links `json:"links"`

		RateLimits RateLimitsMetadata
	}

	// CreateCertificateRequestBody contains the details for requesting a certificate.
	CreateCertificateRequestBody struct {
		// The name of the certificate.
		CertificateName string `json:"certificateName,omitempty"`

		// The key type for a certificate. Valid values are `RSA` or `ECDSA`.
		KeyType CryptographicAlgorithm `json:"keyType"`

		// The key size for a certificate. Valid values for key type RSA: `2048`. Valid values for key type ECDSA: `P-256`, `P-384`.
		KeySize KeySize `json:"keySize"`

		// Secure network type to use for the certificate. Valid values are `ENHANCED_TLS` and `STANDARD_TLS`.
		SecureNetwork SecureNetwork `json:"secureNetwork"`

		// The list of Subject Alternative Names (SANs) for the certificate.
		SANs []string `json:"sans"`

		// Subject fields as defined in X.509 certificates (RFC 5280).
		Subject *Subject `json:"subject,omitempty"`
	}

	// Certificate represents a certificate and its metadata.
	Certificate struct {
		// Unique identifier assigned to the newly created CloudCertificates certificate.
		CertificateID string `json:"certificateId"`

		// Name of the certificate. User provided, if not autogenerated by the system.
		CertificateName string `json:"certificateName"`

		// List of SAN (Subject Alternative Name) domains.
		SANs []string `json:"sans"`

		// Provided subject with optional fields as applicable.
		Subject *Subject `json:"subject"`

		// Certificate type. Defaults to `THIRD_PARTY`.
		CertificateType string `json:"certificateType"`

		// Key type to use for CSR. Either `RSA` or `ECDSA`.
		KeyType CryptographicAlgorithm `json:"keyType"`

		// Key size to use for CSR.
		KeySize KeySize `json:"keySize"`

		// Secure network type to use for the certificate.
		SecureNetwork string `json:"secureNetwork"`

		// Contract identifier.
		ContractID string `json:"contractId"`

		// Account identifier associated with ContractID.
		AccountID string `json:"accountId"`

		// Date the certificate was created in UTC.
		CreatedDate time.Time `json:"createdDate"`

		// User who created the certificate.
		CreatedBy string `json:"createdBy"`

		// Date the certificate was last updated.
		ModifiedDate time.Time `json:"modifiedDate"`

		// User who last modified the certificate.
		ModifiedBy string `json:"modifiedBy"`

		// Status of the certificate.
		CertificateStatus string `json:"certificateStatus"`

		// PEM-encoded certificate signing request (CSR) generated by Akamai for your selected key type.
		CSRPEM string `json:"csrPem"`

		// Date when CSR will expire.
		CSRExpirationDate time.Time `json:"csrExpirationDate"`

		// PEM-encoded signed certificate you uploaded for your selected key type.
		SignedCertificatePEM *string `json:"signedCertificatePem"`

		// Expiration date of signed certificate.
		SignedCertificateNotValidAfterDate *time.Time `json:"signedCertificateNotValidAfterDate"`

		// Date before which the signed certificate is not valid.
		SignedCertificateNotValidBeforeDate *time.Time `json:"signedCertificateNotValidBeforeDate"`

		// Signed certificate serial number in hex format.
		SignedCertificateSerialNumber *string `json:"signedCertificateSerialNumber"`

		// SHA-256 fingerprint of signed certificate.
		SignedCertificateSHA256Fingerprint *string `json:"signedCertificateSHA256Fingerprint"`

		// Issuer field of the signed certificate.
		SignedCertificateIssuer *string `json:"signedCertificateIssuer"`

		// PEM content of Trust chain uploaded by end user.
		TrustChainPEM *string `json:"trustChainPem"`
	}

	// Subject contains the subject details for a certificate.
	Subject struct {
		// Fully qualified domain name (FQDN) or other name associated with the subject.
		CommonName string `json:"commonName,omitempty"`

		// Legal name of the organization.
		Organization string `json:"organization,omitempty"`

		// Two-letter ISO 3166 country code.
		Country string `json:"country,omitempty"`

		// Full name of the state or province.
		State string `json:"state,omitempty"`

		// City or locality name.
		Locality string `json:"locality,omitempty"`
	}

	// CertificateBinding represents a binding between a certificate and a resource.
	CertificateBinding struct {
		// CertificateID is the unique identifier of the certificate.
		CertificateID json.Number `json:"certificateId"`

		// Hostname on the Akamai CDN the certificate applies to.
		Hostname string `json:"hostname"`

		// Network is the network the certificate is bound to. Valid values are `STAGING` and `PRODUCTION`.
		Network string `json:"network"`

		// ResourceType is the type of the certificate is bound to. Currently, only `CDN_HOSTNAME` is available.
		ResourceType string `json:"resourceType"`
	}

	// Links contains pagination and navigation links.
	Links struct {
		// Link to the current page.
		Self string `json:"self"`

		// Link to the next page, if available.
		Next *string `json:"next"`

		// Link to the previous page, if available.
		Previous *string `json:"previous"`
	}

	// ListMetadata contains metadata for paginated lists.
	ListMetadata struct {
		// Total number of items in the list.
		TotalItems int64 `json:"totalItems"`

		// Total number of pages.
		TotalPages int64 `json:"totalPages"`
	}

	// ResourceLimitsMetadata contains information about certificate limits.
	ResourceLimitsMetadata struct {
		// Total number of certificates allowed.
		CertificateLimitTotal *int64

		// Number of remaining certificates that can be created.
		CertificateLimitRemaining *int64
	}

	// RateLimitsMetadata contains information about API rate limits.
	RateLimitsMetadata struct {
		// Limit is the maximum number of requests allowed in the current rate limit window.
		// Nil is returned if the header is not present.
		Limit *int64

		// Remaining is the number of requests remaining in the current rate limit window.
		// Nil is returned if the header is not present.
		Remaining *int64
	}

	// CertificateStatus represents the status of a certificate: `ACTIVE`, `READY_FOR_USE`, or `CSR_READY`.
	CertificateStatus string

	// CryptographicAlgorithm represents the cryptographic algorithm type: `RSA` or `ECDSA`.
	CryptographicAlgorithm string

	// SecureNetwork represents the type of secure network (e.g. `ENHANCED_TLS`).
	SecureNetwork string

	// KeySize represents the size of the key in bits.
	KeySize string

	// Network represents the network type, `STAGING` or `PRODUCTION`.
	Network string
)

var (
	certificateNameLengthRule = validation.Length(0, 270)
	certificateNameRegexRule  = validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 .\-_]+$`)).Error("the input can only contain digits (1-9), letters (a-z, A-Z), spaces, hyphens, periods, and underscores.")
)

// Validate validates PatchCertificateRequest.
func (r PatchCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
		"CertificateName": validation.Validate(r.CertificateName,
			certificateNameLengthRule,
			certificateNameRegexRule),
		"required parameters": validation.Validate(nil, validation.By(func(any) error {
			if r.SignedCertificatePEM == "" && r.CertificateName == nil {
				return errors.New("at least one of SignedCertificatePEM or CertificateName must be provided")
			}
			return nil
		})),
	})
}

// Validate validates ListCertificatesRequest.
func (r ListCertificatesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateStatus": validation.Validate(r.CertificateStatus, validation.By(certificateStatusRule)),
		"KeyType":           validation.Validate(r.KeyType),
		"PageSize": validation.Validate(r.PageSize,
			validation.Min(int64(1)).Error("must be 1 or greater"),
			validation.Max(int64(100)).Error("cannot be greater than 100"),
		),
		"Page": validation.Validate(r.Page,
			validation.Min(int64(1)).Error("must be 1 or greater"),
		),
		"Sort": validation.Validate(r.Sort, validation.When(r.Sort != "", validation.By(sortValidationRule))),
	})
}

// Validate validates CryptographicAlgorithm.
func (c CryptographicAlgorithm) Validate() error {
	return validation.In(CryptographicAlgorithmRSA, CryptographicAlgorithmECDSA).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'", c, CryptographicAlgorithmRSA, CryptographicAlgorithmECDSA)).
		Validate(c)
}

// Validate validates ListCertificateBindingsRequest.
func (r ListCertificateBindingsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
		"PageSize": validation.Validate(r.PageSize,
			validation.Min(int64(1)).Error("must be 1 or greater"),
			validation.Max(int64(100)).Error("cannot be greater than 100"),
		),
		"Page": validation.Validate(r.Page,
			validation.Min(int64(1)).Error("must be 1 or greater"),
		),
	})
}

// Validate validates ListBindingsRequest.
func (r ListBindingsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network, validation.When(r.Network != "", validation.In(NetworkStaging, NetworkProduction).
			Error(fmt.Sprintf("must be either '%s' or '%s'", NetworkStaging, NetworkProduction)))),
		"PageSize": validation.Validate(r.PageSize,
			validation.Min(int64(1)).Error("must be 1 or greater"),
			validation.Max(int64(100)).Error("cannot be greater than 100"),
		),
		"Page": validation.Validate(r.Page,
			validation.Min(int64(1)).Error("must be 1 or greater"),
		),
	})
}

func certificateStatusRule(value any) error {
	statuses, ok := value.([]CertificateStatus)
	if !ok {
		return fmt.Errorf("expected []CertificateStatus, got %T", value)
	}

	validStatuses := []CertificateStatus{
		StatusActive,
		StatusReadyForUse,
		StatusCSRReady,
	}

	for _, s := range statuses {
		if !slices.Contains(validStatuses, s) {
			return fmt.Errorf("list '%v' contains invalid element '%s'. Each element must be one of: '%s'",
				statuses, s, texts.JoinStringBased(validStatuses, "', '"))
		}
	}

	return nil
}

func sortValidationRule(value any) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	re := regexp.MustCompile(`^` + SortFieldPat + `(,` + SortFieldPat + `)*$`)
	if !re.MatchString(s) {
		return fmt.Errorf("must be a comma-separated list of fields, optionally prefixed by + or - (e.g. +createdDate,-certificateName)")
	}

	return nil
}

// Validate validates CreateCertificateRequest.
func (r CreateCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID":      validation.Validate(r.ContractID, validation.Required),
		"GroupID":         validation.Validate(r.GroupID, validation.Required),
		"CertificateName": validation.Validate(r.Body.CertificateName, certificateNameLengthRule, certificateNameRegexRule),
		"KeyType":         validation.Validate(r.Body.KeyType, validation.Required),
		"KeySize":         validation.Validate(r.Body.KeySize, validation.Required),
		"SecureNetwork":   validation.Validate(r.Body.SecureNetwork, validation.Required),
		"SANs":            validation.Validate(r.Body.SANs, validation.Required),
		"Subject":         validation.Validate(r.Body.Subject),
	})
}

// Validate validates UpdateCertificateRequest.
func (r UpdateCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
		"CertificateName": validation.Validate(r.CertificateName,
			certificateNameLengthRule,
			certificateNameRegexRule),
		"SignedCertificatePEM": validation.Validate(r.SignedCertificatePEM, validation.Required.When(r.TrustChainPEM != "")),
	})
}

// Validate validates Subject.
func (s Subject) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CommonName":   validation.Validate(s.CommonName, validation.Length(1, 64), validation.Match(regexp.MustCompile(`\S`))),
		"Organization": validation.Validate(s.Organization, validation.Length(1, 64), validation.Match(regexp.MustCompile(`\S`))),
		"Country":      validation.Validate(s.Country, validation.Length(2, 2), validation.Match(regexp.MustCompile(`\S`))),
		"State":        validation.Validate(s.State, validation.Length(1, 128), validation.Match(regexp.MustCompile(`\S`))),
		"Locality":     validation.Validate(s.Locality, validation.Length(1, 128), validation.Match(regexp.MustCompile(`\S`))),
	})
}

// Validate validates KeySize.
func (k KeySize) Validate() error {
	return validation.In(KeySize2048, KeySizeP256, KeySizeP384).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', or '%s'",
			k, KeySize2048, KeySizeP256, KeySizeP384)).
		Validate(k)
}

// Validate validates SecureNetwork.
func (s SecureNetwork) Validate() error {
	return validation.In(SecureNetworkEnhancedTLS, SecureNetworkStandardTLS).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'",
			s, SecureNetworkEnhancedTLS, SecureNetworkStandardTLS)).
		Validate(s)
}

// Validate validates GetCertificateRequest.
func (r GetCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
	})
}

// Validate validates DeleteCertificateRequest.
func (r DeleteCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
	})
}
