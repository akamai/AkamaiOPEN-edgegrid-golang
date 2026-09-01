package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	// GetPropertyVersionHostnamesRequest contains parameters required to list property version hostnames
	GetPropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		IncludeCertStatus bool
	}

	// GetPropertyVersionHostnamesResponse contains all property version hostnames associated to the given parameters
	GetPropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		PropertyName    string                `json:"propertyName"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// HostnameResponseItems contains the response body for GetPropertyVersionHostnamesResponse
	HostnameResponseItems struct {
		Items []HostnameResponseItem `json:"items"`
	}

	// Hostname contains information about each of the hostname that will be used in UpdatePropertyVersionHostnames
	Hostname struct {
		// CnameType is only one supported `EDGE_HOSTNAME` value.
		CnameType HostnameCnameType `json:"cnameType"`

		// EdgeHostnameID identifies each edge hostname.
		EdgeHostnameID string `json:"edgeHostnameId,omitempty"`

		// CnameFrom is the hostname that your end users see, indicated by the `Host` header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// CnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		CnameTo string `json:"cnameTo,omitempty"`

		// CertProvisioningType indicates the certificate's provisioning type. Either `CPS_MANAGED` type for the certificates you create with the Certificate Provisioning System API (CPS), `DEFAULT` for the Default Domain Validation (DV) certificates created automatically, or `CCM` type for the third party certificates you create with the Cloud Certificate Manager.
		CertProvisioningType string `json:"certProvisioningType"`

		// CCMCertificates is certificate identifiers and links for the CCM-managed certificates.
		CCMCertificates *CCMCertificates `json:"ccmCertificates,omitempty"`

		// MTLS is mutual TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		MTLS *MTLS `json:"mtls,omitempty"`

		// TLSConfiguration is optional TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		TLSConfiguration *TLSConfiguration `json:"tlsConfiguration,omitempty"`
	}

	// HostnameResponseItem contains information about each of the HostnameResponseItems
	HostnameResponseItem struct {
		// CnameType is only one supported `EDGE_HOSTNAME` value.
		CnameType string `json:"cnameType"`

		// EdgeHostnameID identifies each edge hostname.
		EdgeHostnameID string `json:"edgeHostnameId"`

		// CnameFrom is the hostname that your end users see, indicated by the `Host` header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// CnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		CnameTo string `json:"cnameTo"`

		// CertProvisioningType indicates the certificate's provisioning type. Either `CPS_MANAGED` type for the certificates you create with the Certificate Provisioning System API (CPS), `DEFAULT` for the Default Domain Validation (DV) certificates created automatically, or `CCM` type for the third party certificates you create with the Cloud Certificate Manager.
		CertProvisioningType string `json:"certProvisioningType"`

		// CertStatus determines whether a hostname is capable of serving secure content over the staging or production network.
		CertStatus CertStatusItem `json:"certStatus"`

		// CCMCertStatus is deployment status for the RSA and ECDSA certificates created with Cloud Certificate Manager (CCM).
		CCMCertStatus *CCMCertStatus `json:"ccmCertStatus"`

		// CCMCertificates is certificate identifiers and links for the CCM-managed certificates.
		CCMCertificates *CCMCertificatesResp `json:"ccmCertificates"`

		// MTLS is mutual TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		MTLS *MTLSResp `json:"mtls"`

		// TLSConfiguration is optional TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		TLSConfiguration *TLSConfiguration `json:"tlsConfiguration"`

		// DomainOwnershipVerification is optional domain ownership verification details for the hostname.
		// This field is returned only in responses and should not be populated in requests.
		DomainOwnershipVerification *DomainOwnershipVerification `json:"domainOwnershipVerification"`
	}

	// CCMCertStatus is status of CCM certificates in each environment.
	CCMCertStatus struct {
		// ECDSAProductionStatus is ECDSA production status.
		ECDSAProductionStatus string `json:"ecdsaProductionStatus,omitempty"`

		// ECDSAStagingStatus is ECDSA staging status.
		ECDSAStagingStatus string `json:"ecdsaStagingStatus,omitempty"`

		// RSAProductionStatus is RSA production status.
		RSAProductionStatus string `json:"rsaProductionStatus,omitempty"`

		// RSAStagingStatus is RSA staging status.
		RSAStagingStatus string `json:"rsaStagingStatus,omitempty"`
	}

	// CertStatusItem contains information about certificate status for specific Hostname.
	CertStatusItem struct {
		// ValidationCname is the CNAME record used to validate the certificate’s domain.
		ValidationCname ValidationCname `json:"validationCname,omitempty"`

		// Staging contains certificate status for the hostname on the staging network.
		Staging []StatusItem `json:"staging,omitempty"`

		// Production contains certificate status for the hostname on the production network.
		Production []StatusItem `json:"production,omitempty"`

		// Authorization contains details domain validation methods available for your certificate.
		Authorization *Authorization `json:"authorization"`
	}

	// Authorization contains details domain validation methods available for your certificate.
	Authorization struct {
		// DNS01 contains details on the manual DNS validation method.
		DNS01 *DNSAuthorization `json:"dns01"`

		// HTTP01 contains details on the manual HTTP validation method.
		HTTP01 *HTTPAuthorization `json:"http01"`

		// Status of the validation that proves you control the domains listed in the certificate request.
		Status string `json:"status"`

		// ValidUntil is an ISO 8601 timestamp indicating when the domain validation challenge expires.
		ValidUntil *time.Time `json:"validUntil"`
	}

	// DNSAuthorization contains details on the manual DNS validation method.
	DNSAuthorization struct {
		// Result provides details on the validation challenge generation.
		Result AuthorizationResult `json:"result"`

		// Value is the token you need to copy to the DNS TXT record.
		Value string `json:"value"`
	}

	// AuthorizationResult provides details on the validation challenge generation.
	AuthorizationResult struct {
		// Message is a descriptive message on the challenge generation process.
		Message string `json:"message"`

		// Source is the system that sent the result details, either the Certificate Authority (CA) server or Certificate Management System (CPS).
		Source string `json:"source"`

		// Timestamp is the ISO 8601 timestamp indicating when the result was generated.
		Timestamp time.Time `json:"timestamp"`
	}

	// HTTPAuthorization contains details on the manual DNS validation method.
	HTTPAuthorization struct {
		// Body is the token you need to copy to the file on your origin server.
		Body string `json:"body"`

		// Result provides details on the validation challenge generation.
		Result AuthorizationResult `json:"result"`

		// URL is the location on your origin server where you save the file with the token.
		URL string `json:"url"`
	}

	// CertStatusPatchBucketItem contains information about certificate status for specific Hostname that is returned for hostname buckets
	CertStatusPatchBucketItem struct {
		// ValidationCname is the CNAME record used to validate the certificate’s domain.
		ValidationCname ValidationCname `json:"validationCname"`

		// Staging contains certificate status for the hostname on the staging network.
		Staging []StatusItem `json:"staging"`

		// Production contains certificate status for the hostname on the production network.
		Production []StatusItem `json:"production"`
	}

	// CCMCertificates contains identifiers for the RSA and ECDSA certificates.
	CCMCertificates struct {
		// ECDSACertID is certificate ID for ECDSA.
		ECDSACertID string `json:"ecdsaCertId,omitempty"`

		// RSACertID is certificate ID for RSA.
		RSACertID string `json:"rsaCertId,omitempty"`
	}

	// CCMCertificatesResp contains identifiers for the RSA and ECDSA certificates.
	CCMCertificatesResp struct {
		CCMCertificates

		// ECDSACertLink is link to the ECSDA certificate.
		ECDSACertLink string `json:"ecdsaCertLink"`

		// RSACertLink is link to the RSA certificate.
		RSACertLink string `json:"rsaCertLink"`
	}

	// MTLS is mutual TLS configuration used only for the `CCM` provisioning.
	MTLS struct {
		// CASetID is ID of the Client CA set used for mutual TLS.
		CASetID string `json:"caSetId,omitempty"`

		// CheckClientOCSP specifies whether to check the OCSP status of the client certificate.
		CheckClientOCSP bool `json:"checkClientOcsp,omitempty"`

		// SendCASetClient specifies whether to send the CA set to the client during the TLS handshake.
		SendCASetClient bool `json:"sendCaSetClient,omitempty"`
	}

	// MTLSResp is mutual TLS configuration used only for the `CCM` provisioning.
	MTLSResp struct {
		MTLS

		// CASetLink is link of the Client CA set used for mutual TLS.
		CASetLink string `json:"caSetLink"`
	}

	// TLSConfiguration is optional TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
	TLSConfiguration struct {
		// CipherProfile is cipher profile name.
		CipherProfile string `json:"cipherProfile,omitempty"`

		// DisallowedTLSVersions is  list of TLS versions that are disallowed.
		DisallowedTLSVersions []string `json:"disallowedTlsVersions,omitempty"`

		// StapleServerOcspResponse specifies whether to staple the OCSP response for the server certificate.
		StapleServerOcspResponse bool `json:"stapleServerOcspResponse,omitempty"`

		// FIPSMode specifies whether to enable the FIPS mode.
		FIPSMode bool `json:"fipsMode,omitempty"`
	}

	// ValidationCname is the CNAME record used to validate the certificate’s domain
	ValidationCname struct {
		Hostname string `json:"hostname,omitempty"`
		Target   string `json:"target,omitempty"`
	}

	// StatusItem determines whether a hostname is capable of serving secure content over the staging or production network.
	StatusItem struct {
		Status string `json:"status,omitempty"`
	}

	// DomainOwnershipVerification contains domain ownership verification details for the hostname.
	DomainOwnershipVerification struct {
		Status                   string           `json:"status"`
		ChallengeTokenExpiryDate *time.Time       `json:"challengeTokenExpiryDate"`
		ValidationCname          *ValidationCname `json:"validationCname"`
		ValidationHTTP           *ValidationHTTP  `json:"validationHttp"`
		ValidationTXT            *ValidationTXT   `json:"validationTxt"`
	}

	// ValidationHTTP contains HTTP validation methods for domain ownership verification.
	ValidationHTTP struct {
		FileContentMethod FileContentMethod `json:"fileContentMethod"`
		RedirectMethod    RedirectMethod    `json:"redirectMethod"`
	}

	// FileContentMethod contains details for the file content method of validation.
	FileContentMethod struct {
		Body string `json:"body"`
		URL  string `json:"url"`
	}

	// RedirectMethod contains details for the HTTP redirect method of validation.
	RedirectMethod struct {
		HTTPRedirectFrom string `json:"httpRedirectFrom"`
		HTTPRedirectTo   string `json:"httpRedirectTo"`
	}

	// ValidationTXT contains TXT record validation details for domain ownership verification.
	ValidationTXT struct {
		ChallengeToken string `json:"challengeToken"`
		Hostname       string `json:"hostname"`
	}

	// UpdatePropertyVersionHostnamesRequest contains parameters required to update the set of hostname entries for a property version
	UpdatePropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		IncludeCertStatus bool
		Hostnames         []Hostname
	}

	// UpdatePropertyVersionHostnamesResponse contains information about each of the HostnameRequestItems
	UpdatePropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		PropertyName    string                `json:"propertyName"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// PatchPropertyVersionHostnamesRequest contains parameters for patch property version hostnames
	PatchPropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		IncludeCertStatus bool
		Body              PatchPropertyVersionHostnamesRequestBody
	}

	// PatchPropertyVersionHostnamesRequestBody contains the request body for patching property version hostnames
	PatchPropertyVersionHostnamesRequestBody struct {
		Add    []HostnameAdd `json:"add,omitempty"`
		Remove []string      `json:"remove,omitempty"`
	}

	// HostnameAdd contains information about the hostname to be added in patch property version hostnames
	HostnameAdd struct {
		// CnameFrom is the hostname that your end users see, indicated by the `Host` header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// CnameType is only one supported `EDGE_HOSTNAME` value.
		CnameType HostnameCnameType `json:"cnameType,omitempty"`

		// CnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		CnameTo string `json:"cnameTo,omitempty"`

		// CertProvisioningType indicates the certificate's provisioning type. Either `CPS_MANAGED` type for the certificates you create with the Certificate Provisioning System API (CPS), `DEFAULT` for the Default Domain Validation (DV) certificates created automatically, or `CCM` type for the third party certificates you create with the Cloud Certificate Manager.
		CertProvisioningType CertType `json:"certProvisioningType,omitempty"`

		// EdgeHostnameID identifies each edge hostname.
		EdgeHostnameID string `json:"edgeHostnameId,omitempty"`

		// MTLS is mutual TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		MTLS *MTLS `json:"mtls,omitempty"`

		// TLSConfiguration is optional TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		TLSConfiguration *TLSConfiguration `json:"tlsConfiguration,omitempty"`

		// CCMCertStatus is deployment status for the RSA and ECDSA certificates created with Cloud Certificate Manager (CCM).
		CCMCertificates *CCMCertificates `json:"ccmCertificates,omitempty"`
	}

	// PatchPropertyVersionHostnamesResponse contains response from patch property version hostnames
	PatchPropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		Etag            string                `json:"etag"`
		PropertyID      string                `json:"propertyId"`
		PropertyName    string                `json:"propertyName"`
		PropertyVersion int                   `json:"propertyVersion"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// GetAuditHistoryRequest contains parameters for getting audit history of property hostname
	GetAuditHistoryRequest struct {
		// The cnameFrom for the hostname your end users see, indicated by the Host header in end user requests.
		Hostname string
	}

	// GetAuditHistoryResponse contains the audit history for property hostname
	GetAuditHistoryResponse struct {
		Hostname string          `json:"hostname"`
		History  HostnameHistory `json:"history"`
	}

	// HostnameHistory contains the entries of changes made to the property hostname
	HostnameHistory struct {
		Items []HostnameHistoryItem `json:"items"`
	}

	// HostnameHistoryItem contains information about each of the entry in the hostname history
	HostnameHistoryItem struct {
		// Action is the type of action performed to the property hostname, either:
		// `ACTIVATE` if the hostname is currently serving traffic,
		// `DEACTIVATE` if the hostname isn't serving traffic,
		// `ADD` if the user requested to add the hostname to a property,
		// `REMOVE` if the user requested to remove the hostname from a property,
		// `MOVE` if the hostname was moved from one property to another,
		// `MODIFY` if the user changed the `edgeHostnameId` or `certProvisioningType` values for an already-activated hostname,
		// `ABORTED` when the user request to cancel the hostname activation,
		// `ERROR` if the hostname activation failed.
		Action string `json:"action"`

		// CertProvisioningType indicates the type of the certificate used in the property hostname.
		// Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System API (CPS),
		// `DEFAULT` for Default Domain Validation (DV) certificates deployed automatically,
		// or `CCM` for the third party certificates you create with the Cloud Certificate Manager.
		// Note that you can't specify the `DEFAULT` value if your account hostname uses the `akamaized.net` domain suffix.
		CertProvisioningType string `json:"certProvisioningType"`

		// CnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		// This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		CnameTo string `json:"cnameTo"`

		// ContractID identifies the prevailing contract under which you requested the data.
		ContractID string `json:"contractId"`

		// EdgeHostnameID identifies the edge hostname you mapped your traffic to.
		EdgeHostnameID string `json:"edgeHostnameId"`

		// GroupID identifies the group under which the property activated.
		GroupID string `json:"groupId"`

		// Network is the network of activated hostnames, either `STAGING` or `PRODUCTION`.
		Network string `json:"network"`

		// PropertyID is the unique identifier for the property.
		PropertyID string `json:"propertyId"`

		// Timestamp indicates when the action occurred.
		Timestamp string `json:"timestamp"`

		// User is the user who initiated the action.
		User string `json:"user"`
	}

	// HostnameCnameType represents HostnameCnameType enum
	HostnameCnameType string
)

const (
	// HostnameCnameTypeEdgeHostname const
	HostnameCnameTypeEdgeHostname HostnameCnameType = "EDGE_HOSTNAME"
)

// Validate validates HostnameCnameType.
func (t HostnameCnameType) Validate() validation.InRule {
	return validation.In(HostnameCnameTypeEdgeHostname).
		Error(fmt.Sprintf("value '%s' is invalid. There is only one supported value of: %s",
			t, HostnameCnameTypeEdgeHostname))
}

// Validate validates GetPropertyVersionHostnamesRequest
func (ph GetPropertyVersionHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ph.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ph.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates UpdatePropertyVersionHostnamesRequest
func (r UpdatePropertyVersionHostnamesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID":      validation.Validate(r.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(r.PropertyVersion, validation.Required),
		"Hostnames":       validation.Validate(r.Hostnames),
	})
}

// Validate validates Hostname
func (h Hostname) Validate() error {
	return validation.Errors{
		"MTLS":                validation.Validate(h.MTLS),
		"CCMCertificates":     validation.Validate(h.CCMCertificates),
		"ValidateCCMHostname": validateCCMHostname(h.CertProvisioningType, h.CCMCertificates, h.MTLS, h.TLSConfiguration),
	}.Filter()
}

// Validate validates MTLS
func (m MTLS) Validate() error {
	return validation.Errors{
		"CASetID": validation.Validate(m.CASetID, validation.Required, is.Digit),
	}.Filter()
}

// Validate validates CCMCertificates
func (c CCMCertificates) Validate() error {
	return validation.Errors{
		"RSACertID":   validation.Validate(c.RSACertID, validation.When(c.RSACertID != "", is.Digit)),
		"ECDSACertID": validation.Validate(c.ECDSACertID, validation.When(c.ECDSACertID != "", is.Digit)),
	}.Filter()
}

func validateCCMHostname(certType string, certs *CCMCertificates, mTLS *MTLS, tls *TLSConfiguration) error {
	if certType != string(CertTypeCCM) {
		if certs != nil {
			return errors.New("the CCM cert details are provided without `certProvisioningType` set to `CCM`")
		}
		if mTLS != nil {
			return errors.New("the mTLS configuration is provided without `certProvisioningType` set to `CCM`")
		}
		if tls != nil {
			return errors.New("the TLS configuration is provided without `certProvisioningType` set to `CCM`")
		}
		return nil
	}
	if certs == nil {
		return errors.New("when using `certProvisioningType` set to `CCM`, the request body must contain `ccmCertificates` with at least `rsaCertId` or `ecdsaCertId`")
	}
	if certs.RSACertID == "" && certs.ECDSACertID == "" {
		return errors.New("either RSACertID or ECDSACertID must be provided")
	}
	if tls != nil && len(tls.CipherProfile) == 0 {
		return errors.New("the cipher profile is empty in the TLS configuration")
	}
	return nil
}

// Validate validates PatchPropertyVersionHostnamesRequest
func (r PatchPropertyVersionHostnamesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PropertyID":      validation.Validate(r.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(r.PropertyVersion, validation.Required),
		"Body":            validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates PatchPropertyVersionHostnamesRequestBody
func (b PatchPropertyVersionHostnamesRequestBody) Validate() error {
	return validation.Errors{
		"Add": validation.Validate(b.Add),
	}.Filter()
}

// Validate validates GetAuditHistoryRequest
func (r GetAuditHistoryRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Hostname": validation.Validate(r.Hostname, validation.Required),
	})
}

// Validate validates HostnameAdd
func (h HostnameAdd) Validate() error {
	return validation.Errors{
		"CnameFrom": validation.Validate(h.CnameFrom, validation.Required),
		"CnameType": validation.Validate(h.CnameType,
			validation.When(h.CnameType != "", h.CnameType.Validate())),
		"CertProvisioningType": validation.Validate(h.CertProvisioningType,
			validation.When(h.CertProvisioningType != "", h.CertProvisioningType.Validate())),
		"MTLS":                validation.Validate(h.MTLS),
		"CCMCertificates":     validation.Validate(h.CCMCertificates),
		"ValidateCCMHostname": validateCCMHostname(string(h.CertProvisioningType), h.CCMCertificates, h.MTLS, h.TLSConfiguration),
		"required parameters": validation.Validate(nil,
			validation.By(func(interface{}) error {
				if h.CnameTo == "" && h.EdgeHostnameID == "" {
					return fmt.Errorf("either CnameTo or EdgeHostnameID must be provided")
				}
				return nil
			})),
	}.Filter()
}

var (
	// ErrGetPropertyVersionHostnames represents error when fetching hostnames fails
	ErrGetPropertyVersionHostnames = errors.New("fetching hostnames")
	// ErrUpdatePropertyVersionHostnames represents error when updating hostnames fails
	ErrUpdatePropertyVersionHostnames = errors.New("updating hostnames")
	// ErrPatchPropertyVersionHostnames represents error when patching hostnames fails
	ErrPatchPropertyVersionHostnames = errors.New("patching hostnames")
	// ErrGetAuditHistory represents error when getting audit history fails
	ErrGetAuditHistory = errors.New("getting audit history")
)

func (p *papi) GetPropertyVersionHostnames(ctx context.Context, params GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetPropertyVersionHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersionHostnames, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/properties/%s/versions/%d/hostnames", params.PropertyID, params.PropertyVersion).
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParam("validateHostnames", strconv.FormatBool(params.ValidateHostnames)).
		AddQueryParam("includeCertStatus", strconv.FormatBool(params.IncludeCertStatus)).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyVersionHostnames, err)
	}

	var hostnames GetPropertyVersionHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyVersionHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersionHostnames, ErrNotFound, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) UpdatePropertyVersionHostnames(ctx context.Context, params UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdatePropertyVersionHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePropertyVersionHostnames, ErrStructValidation, err)
	}

	var hostnames UpdatePropertyVersionHostnamesResponse
	newHostnames := params.Hostnames
	if newHostnames == nil {
		newHostnames = []Hostname{}
	}
	req, err := request.NewPut(ctx, "/papi/v1/properties/%s/versions/%v/hostnames", params.PropertyID, params.PropertyVersion).
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParam("validateHostnames", strconv.FormatBool(params.ValidateHostnames)).
		AddQueryParam("includeCertStatus", strconv.FormatBool(params.IncludeCertStatus)).
		WithBody(newHostnames).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePropertyVersionHostnames, err)
	}

	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePropertyVersionHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) PatchPropertyVersionHostnames(ctx context.Context, params PatchPropertyVersionHostnamesRequest) (*PatchPropertyVersionHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("PatchPropertyVersionHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrPatchPropertyVersionHostnames, ErrStructValidation, err)
	}

	req, err := request.NewPatch(ctx, "/papi/v1/properties/%s/versions/%v/hostnames", params.PropertyID, params.PropertyVersion).
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("validateHostnames", strconv.FormatBool(params.ValidateHostnames), params.ValidateHostnames).
		AddQueryParamIf("includeCertStatus", strconv.FormatBool(params.IncludeCertStatus), params.IncludeCertStatus).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrPatchPropertyVersionHostnames, err)
	}

	var result PatchPropertyVersionHostnamesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrPatchPropertyVersionHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrPatchPropertyVersionHostnames, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) GetAuditHistory(ctx context.Context, params GetAuditHistoryRequest) (*GetAuditHistoryResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAuditHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetAuditHistory, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/hostnames/%s/audit-history", params.Hostname).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetAuditHistory, err)
	}

	var history GetAuditHistoryResponse
	resp, err := p.Exec(req, &history)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetAuditHistory, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetAuditHistory, p.Error(resp))
	}

	return &history, nil
}
