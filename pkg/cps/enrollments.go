package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Enrollments is a CPS enrollments API interface
	Enrollments interface {
		// ListEnrollments fetches all enrollments with given contractId
		//
		// See https://techdocs.akamai.com/cps/reference/get-enrollments
		ListEnrollments(context.Context, ListEnrollmentsRequest) (*ListEnrollmentsResponse, error)

		// GetEnrollment fetches enrollment object with given ID
		//
		// See: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#getasingleenrollment
		GetEnrollment(context.Context, GetEnrollmentRequest) (*Enrollment, error)

		// CreateEnrollment creates a new enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/post-enrollment
		CreateEnrollment(context.Context, CreateEnrollmentRequest) (*CreateEnrollmentResponse, error)

		// UpdateEnrollment updates a single enrollment entry with given ID
		//
		// See: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#putasingleenrollment
		UpdateEnrollment(context.Context, UpdateEnrollmentRequest) (*UpdateEnrollmentResponse, error)

		// RemoveEnrollment removes an enrollment with given ID
		//
		// See: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#deleteasingleenrollment
		RemoveEnrollment(context.Context, RemoveEnrollmentRequest) (*RemoveEnrollmentResponse, error)
	}

	// ListEnrollmentsResponse represents list of CPS enrollment objects under given contractId. It is used as a response body while fetching enrollments by contractId
	ListEnrollmentsResponse struct {
		Enrollments []Enrollment `json:"enrollments"`
	}

	// Enrollment represents a CPS enrollment object. It is used both as a request body for enrollment creation and response body while fetching enrollment by ID
	Enrollment struct {
		AdminContact                   *Contact              `json:"adminContact"`
		AutoRenewalStartTime           string                `json:"autoRenewalStartTime,omitempty"`
		CertificateChainType           string                `json:"certificateChainType,omitempty"`
		CertificateType                string                `json:"certificateType"`
		ChangeManagement               bool                  `json:"changeManagement"`
		CSR                            *CSR                  `json:"csr"`
		EnableMultiStackedCertificates bool                  `json:"enableMultiStackedCertificates"`
		Location                       string                `json:"location,omitempty"`
		MaxAllowedSanNames             int                   `json:"maxAllowedSanNames,omitempty"`
		MaxAllowedWildcardSanNames     int                   `json:"maxAllowedWildcardSanNames,omitempty"`
		NetworkConfiguration           *NetworkConfiguration `json:"networkConfiguration"`
		Org                            *Org                  `json:"org"`
		PendingChanges                 []string              `json:"pendingChanges,omitempty"`
		RA                             string                `json:"ra"`
		SignatureAlgorithm             string                `json:"signatureAlgorithm,omitempty"`
		TechContact                    *Contact              `json:"techContact"`
		ThirdParty                     *ThirdParty           `json:"thirdParty,omitempty"`
		ValidationType                 string                `json:"validationType"`
	}

	// Contact contains contact information
	Contact struct {
		AddressLineOne   string `json:"addressLineOne,omitempty"`
		AddressLineTwo   string `json:"addressLineTwo,omitempty"`
		City             string `json:"city,omitempty"`
		Country          string `json:"country,omitempty"`
		Email            string `json:"email,omitempty"`
		FirstName        string `json:"firstName,omitempty"`
		LastName         string `json:"lastName,omitempty"`
		OrganizationName string `json:"organizationName,omitempty"`
		Phone            string `json:"phone,omitempty"`
		PostalCode       string `json:"postalCode,omitempty"`
		Region           string `json:"region,omitempty"`
		Title            string `json:"title,omitempty"`
	}

	// CSR is a Certificate Signing Request object
	CSR struct {
		C    string   `json:"c,omitempty"`
		CN   string   `json:"cn"`
		L    string   `json:"l,omitempty"`
		O    string   `json:"o,omitempty"`
		OU   string   `json:"ou,omitempty"`
		SANS []string `json:"sans,omitempty"`
		ST   string   `json:"st,omitempty"`
	}

	// NetworkConfiguration contains settings that specify any network information and TLS Metadata you want CPS to use to push the completed certificate to the network
	NetworkConfiguration struct {
		ClientMutualAuthentication *ClientMutualAuthentication `json:"clientMutualAuthentication,omitempty"`
		DisallowedTLSVersions      []string                    `json:"disallowedTlsVersions,omitempty"`
		DNSNameSettings            *DNSNameSettings            `json:"dnsNameSettings,omitempty"`
		Geography                  string                      `json:"geography,omitempty"`
		MustHaveCiphers            string                      `json:"mustHaveCiphers,omitempty"`
		OCSPStapling               OCSPStapling                `json:"ocspStapling,omitempty"`
		PreferredCiphers           string                      `json:"preferredCiphers,omitempty"`
		QuicEnabled                bool                        `json:"quicEnabled"`
		SecureNetwork              string                      `json:"secureNetwork,omitempty"`
		SNIOnly                    bool                        `json:"sniOnly"`
	}

	// ClientMutualAuthentication specifies the trust chain that is used to verify client certificates and some configuration options
	ClientMutualAuthentication struct {
		AuthenticationOptions *AuthenticationOptions `json:"authenticationOptions,omitempty"`
		SetID                 string                 `json:"setId,omitempty"`
	}

	// AuthenticationOptions contain the configuration options for the selected trust chain
	AuthenticationOptions struct {
		OCSP               *OCSP `json:"ocsp,omitempty"`
		SendCAListToClient *bool `json:"sendCaListToClient,omitempty"`
	}

	// OCSP specifies whether you want to enable ocsp stapling for client certificates
	OCSP struct {
		Enabled *bool `json:"enabled,omitempty"`
	}

	// DNSNameSettings contain DNS name setting in given network configuration
	DNSNameSettings struct {
		CloneDNSNames bool     `json:"cloneDnsNames"`
		DNSNames      []string `json:"dnsNames,omitempty"`
	}

	// Org represents organization information
	Org struct {
		AddressLineOne string `json:"addressLineOne,omitempty"`
		AddressLineTwo string `json:"addressLineTwo,omitempty"`
		City           string `json:"city,omitempty"`
		Country        string `json:"country,omitempty"`
		Name           string `json:"name,omitempty"`
		Phone          string `json:"phone,omitempty"`
		PostalCode     string `json:"postalCode,omitempty"`
		Region         string `json:"region,omitempty"`
	}

	// ThirdParty specifies that you want to use a third party certificate
	ThirdParty struct {
		ExcludeSANS bool `json:"excludeSans"`
	}

	// ListEnrollmentsRequest contains Contract ID of enrollments that are to be fetched with ListEnrollments
	ListEnrollmentsRequest struct {
		ContractID string
	}

	// GetEnrollmentRequest contains ID of an enrollment that is to be fetched with GetEnrollment
	GetEnrollmentRequest struct {
		EnrollmentID int
	}

	// CreateEnrollmentRequest contains request body and path parameters used to create an enrollment
	CreateEnrollmentRequest struct {
		Enrollment
		ContractID       string
		DeployNotAfter   string
		DeployNotBefore  string
		AllowDuplicateCN bool
	}

	// CreateEnrollmentResponse contains response body returned after successful enrollment creation
	CreateEnrollmentResponse struct {
		ID         int
		Enrollment string   `json:"enrollment"`
		Changes    []string `json:"changes"`
	}

	// UpdateEnrollmentRequest contains request body and path parameters used to update an enrollment
	UpdateEnrollmentRequest struct {
		Enrollment
		EnrollmentID              int
		AllowCancelPendingChanges *bool
		AllowStagingBypass        *bool
		DeployNotAfter            string
		DeployNotBefore           string
		ForceRenewal              *bool
		RenewalDateCheckOverride  *bool
	}

	// UpdateEnrollmentResponse contains response body returned after successful enrollment update
	UpdateEnrollmentResponse struct {
		ID         int
		Enrollment string   `json:"enrollment"`
		Changes    []string `json:"changes"`
	}

	// RemoveEnrollmentRequest contains parameters necessary to send a RemoveEnrollment request
	RemoveEnrollmentRequest struct {
		EnrollmentID              int
		AllowCancelPendingChanges *bool
		DeployNotAfter            string
		DeployNotBefore           string
	}

	// RemoveEnrollmentResponse contains response body returned after successful enrollment deletion
	RemoveEnrollmentResponse struct {
		Enrollment string   `json:"enrollment"`
		Changes    []string `json:"changes"`
	}

	// OCSPStapling is used to enable OCSP stapling for an enrollment
	OCSPStapling string
)

const (
	// OCSPStaplingOn parameter value
	OCSPStaplingOn OCSPStapling = "on"
	// OCSPStaplingOff parameter value
	OCSPStaplingOff OCSPStapling = "off"
	// OCSPStaplingNotSet parameter value
	OCSPStaplingNotSet OCSPStapling = "not-set"
)

// Validate performs validation on Enrollment
func (e Enrollment) Validate() error {
	return validation.Errors{
		"adminContact":         validation.Validate(e.AdminContact, validation.Required),
		"certificateType":      validation.Validate(e.CertificateType, validation.Required),
		"csr":                  validation.Validate(e.CSR, validation.Required),
		"networkConfiguration": validation.Validate(e.NetworkConfiguration, validation.Required),
		"org":                  validation.Validate(e.Org, validation.Required),
		"ra":                   validation.Validate(e.RA, validation.Required),
		"techContact":          validation.Validate(e.TechContact, validation.Required),
		"validationType":       validation.Validate(e.ValidationType, validation.Required),
		"thirdParty":           validation.Validate(e.ThirdParty),
	}.Filter()
}

// Validate performs validation on Enrollment
func (c CSR) Validate() error {
	return validation.Errors{
		"cn": validation.Validate(c.CN, validation.Required),
	}.Filter()
}

// Validate performs validation on NetworkConfiguration
func (n NetworkConfiguration) Validate() error {
	return validation.Errors{
		"ocspStapling": validation.Validate(n.OCSPStapling, validation.In(OCSPStaplingOn, OCSPStaplingOff, OCSPStaplingNotSet)),
	}.Filter()
}

// Validate performs validation on ListEnrollmentRequest
func (e ListEnrollmentsRequest) Validate() error {
	return validation.Errors{
		"contractId": validation.Validate(e.ContractID, validation.Required),
	}.Filter()
}

// Validate performs validation on GetEnrollmentRequest
func (e GetEnrollmentRequest) Validate() error {
	return validation.Errors{
		"enrollmentId": validation.Validate(e.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate performs validation on CreateEnrollmentRequest
func (e CreateEnrollmentRequest) Validate() error {
	return validation.Errors{
		"enrollment": validation.Validate(e.Enrollment, validation.Required),
		"contractId": validation.Validate(e.ContractID, validation.Required),
	}.Filter()
}

// Validate performs validation on UpdateEnrollmentRequest
func (e UpdateEnrollmentRequest) Validate() error {
	return validation.Errors{
		"enrollment":   validation.Validate(e.Enrollment, validation.Required),
		"enrollmentId": validation.Validate(e.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate performs validation on RemoveEnrollmentRequest
func (e RemoveEnrollmentRequest) Validate() error {
	return validation.Errors{
		"enrollmentId": validation.Validate(e.EnrollmentID, validation.Required),
	}.Filter()
}

var (
	// ErrListEnrollments is returned when ListEnrollments fails
	ErrListEnrollments = errors.New("fetching enrollments")
	// ErrGetEnrollment is returned when GetEnrollment fails
	ErrGetEnrollment = errors.New("fetching enrollment")
	// ErrCreateEnrollment is returned when CreateEnrollment fails
	ErrCreateEnrollment = errors.New("create enrollment")
	// ErrUpdateEnrollment is returned when UpdateEnrollment fails
	ErrUpdateEnrollment = errors.New("update enrollment")
	// ErrRemoveEnrollment is returned when RemoveEnrollment fails
	ErrRemoveEnrollment = errors.New("remove enrollment")
)

func (c *cps) ListEnrollments(ctx context.Context, params ListEnrollmentsRequest) (*ListEnrollmentsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListEnrollments, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("ListEnrollments")

	uri := fmt.Sprintf("/cps/v2/enrollments?contractId=%s", params.ContractID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEnrollments, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollments.v9+json")

	var result ListEnrollmentsResponse

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEnrollments, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEnrollments, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetEnrollment(ctx context.Context, params GetEnrollmentRequest) (*Enrollment, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEnrollment, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("GetEnrollment")

	uri := fmt.Sprintf("/cps/v2/enrollments/%d", params.EnrollmentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEnrollment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment.v9+json")

	var result Enrollment

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEnrollment, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEnrollment, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) CreateEnrollment(ctx context.Context, params CreateEnrollmentRequest) (*CreateEnrollmentResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEnrollment, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("CreateEnrollment")

	uri, err := url.Parse(fmt.Sprintf("/cps/v2/enrollments?contractId=%s", params.ContractID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrCreateEnrollment, err)
	}
	query := uri.Query()
	if params.DeployNotAfter != "" {
		query.Add("deploy-not-after", params.DeployNotAfter)
	}
	if params.DeployNotBefore != "" {
		query.Add("deploy-not-before", params.DeployNotBefore)
	}
	if params.AllowDuplicateCN {
		query.Add("allow-duplicate-cn", "true")
	}
	uri.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEnrollment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.enrollment.v9+json")

	var result CreateEnrollmentResponse

	resp, err := c.Exec(req, &result, params.Enrollment)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEnrollment, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrCreateEnrollment, c.Error(resp))
	}
	id, err := GetIDFromLocation(result.Enrollment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEnrollment, ErrInvalidLocation, err)
	}
	result.ID = id

	return &result, nil
}

func (c *cps) UpdateEnrollment(ctx context.Context, params UpdateEnrollmentRequest) (*UpdateEnrollmentResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEnrollment, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("UpdateEnrollment")

	uri, err := url.Parse(fmt.Sprintf("/cps/v2/enrollments/%d", params.EnrollmentID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrUpdateEnrollment, err)
	}
	query := uri.Query()
	if params.AllowCancelPendingChanges != nil {
		query.Add("allow-cancel-pending-changes", strconv.FormatBool(*params.AllowCancelPendingChanges))
	}
	if params.AllowStagingBypass != nil {
		query.Add("allow-staging-bypass", strconv.FormatBool(*params.AllowStagingBypass))
	}
	if params.DeployNotAfter != "" {
		query.Add("deploy-not-after", params.DeployNotAfter)
	}
	if params.DeployNotBefore != "" {
		query.Add("deploy-not-before", params.DeployNotBefore)
	}
	if params.ForceRenewal != nil {
		query.Add("force-renewal", strconv.FormatBool(*params.ForceRenewal))
	}
	if params.RenewalDateCheckOverride != nil {
		query.Add("renewal-date-check-override", strconv.FormatBool(*params.RenewalDateCheckOverride))
	}

	uri.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateEnrollment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.enrollment.v9+json")

	var result UpdateEnrollmentResponse

	resp, err := c.Exec(req, &result, params.Enrollment)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEnrollment, err)
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateEnrollment, c.Error(resp))
	}
	id, err := GetIDFromLocation(result.Enrollment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEnrollment, ErrInvalidLocation, err)
	}
	result.ID = id

	return &result, nil
}

func (c *cps) RemoveEnrollment(ctx context.Context, params RemoveEnrollmentRequest) (*RemoveEnrollmentResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrRemoveEnrollment, ErrStructValidation, err)
	}

	logger := c.Log(ctx)
	logger.Debug("RemoveEnrollment")

	uri, err := url.Parse(fmt.Sprintf("/cps/v2/enrollments/%d", params.EnrollmentID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrRemoveEnrollment, err)
	}
	query := uri.Query()
	if params.AllowCancelPendingChanges != nil {
		query.Add("allow-cancel-pending-changes", strconv.FormatBool(*params.AllowCancelPendingChanges))
	}
	if params.DeployNotAfter != "" {
		query.Add("deploy-not-after", params.DeployNotAfter)
	}
	if params.DeployNotBefore != "" {
		query.Add("deploy-not-before", params.DeployNotBefore)
	}

	uri.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRemoveEnrollment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")

	var result RemoveEnrollmentResponse

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRemoveEnrollment, err)
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrRemoveEnrollment, c.Error(resp))
	}

	return &result, nil
}
