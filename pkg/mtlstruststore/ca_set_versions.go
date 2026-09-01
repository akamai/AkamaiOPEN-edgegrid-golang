package mtlstruststore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/texts"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CreateCASetVersionRequest is used to request the creation of CA set version.
	CreateCASetVersionRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// Body is the body of the request.
		Body CreateCASetVersionRequestBody
	}

	// CreateCASetVersionRequestBody represents the body of a CreateCASetVersionRequest.
	CreateCASetVersionRequestBody struct {
		// AllowInsecureSHA1 permits SHA-1 signed certificates if set to true. Defaults to false.
		AllowInsecureSHA1 bool `json:"allowInsecureSha1"`

		// Description is an optional description for the can set.
		Description *string `json:"description,omitempty"`

		// Certificates is a list of valid root or intermediate certificates. At least one is required.
		Certificates []CertificateRequest `json:"certificates"`
	}

	// CloneCASetVersionRequest represents a request to clone a specific version of a CA Set.
	CloneCASetVersionRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64 `json:"version"`
	}

	// ListCASetVersionsRequest represents a request to retrieve a list of CA sets.
	ListCASetVersionsRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// IncludeCertificates includes certificates in the response if true. Defaults to false.
		IncludeCertificates bool

		// ActiveVersionsOnly includes only staging or production active versions if true. Defaults to false.
		ActiveVersionsOnly bool

		// CASetVersionStatuses filters CA sets by status values. Allowed values: `NOT_DELETED`, `DELETED`. Defaults to `NOT_DELETED`.
		CASetVersionStatuses []string
	}

	// GetCASetVersionRequest represents a request to retrieve details of a specific CA Set version.
	GetCASetVersionRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64 `json:"version"`
	}

	// UpdateCASetVersionRequest is used to request the update of an existing CA Set version.
	UpdateCASetVersionRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// Version is the version number within the CA Set.
		Version int64

		// Body is the body of the request.
		Body UpdateCASetVersionRequestBody
	}

	//UpdateCASetVersionRequestBody represent the body of a UpdateCASetVersionRequest.
	UpdateCASetVersionRequestBody struct {
		// Description is an optional description for the ca set.
		Description *string `json:"description,omitempty"`

		// AllowInsecureSHA1 indicates whether SHA-1 certificates are allowed.
		AllowInsecureSHA1 bool `json:"allowInsecureSha1"`

		// Certificates is a list of root or intermediate certificates in the version.
		Certificates []CertificateRequest `json:"certificates"`
	}

	// GetCASetVersionCertificatesRequest represents a request to retrieve certificates details of a specific CA Set version.
	GetCASetVersionCertificatesRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64 `json:"version"`

		// CertificateStatus filters by "EXPIRING", "EXPIRED", "EXPIRING,EXPIRED", "ACTIVE", or "ACTIVE,EXPIRED". Required if ExpiryThresholdInDays or ExpiryThresholdTimestamp are set.
		CertificateStatus *CertificateStatus

		// ExpiryThresholdInDays filters certificates expiring within or expired in past N days. Defaults to 30 if not set. Cannot be used with ExpiryThresholdTimestamp.
		ExpiryThresholdInDays *int

		// ExpiryThresholdTimestamp filters certificates expiring before this timestamp. Cannot be used with ExpiryThresholdInDays.
		ExpiryThresholdTimestamp time.Time
	}

	// CertificateRequest represents details of a certificate used in a CA Set version creation or update.
	CertificateRequest struct {
		// CertificatePEM is the PEM-encoded representation of the certificate.
		CertificatePEM string `json:"certificatePem"`

		// Description is an optional description of the certificate.
		Description *string `json:"description,omitempty"`
	}

	// DeleteCASetVersionRequest represents a request to mark a CA set version for future deletion.
	DeleteCASetVersionRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64
	}

	// CertificateResponse represents details of a certificate returned for a CA Set version.
	CertificateResponse struct {
		// Subject of the certificate.
		Subject string `json:"subject"`

		// Issuer of the certificate.
		Issuer string `json:"issuer"`

		// EndDate is the ISO-8601 date after which the certificate is not valid.
		EndDate time.Time `json:"endDate"`

		// StartDate is the ISO-8601 date before which the certificate is not valid.
		StartDate time.Time `json:"startDate"`

		// Fingerprint is the unique SHA-256 fingerprint of the certificate.
		Fingerprint string `json:"fingerprint"`

		// CertificatePEM is the PEM-encoded representation of the certificate.
		CertificatePEM string `json:"certificatePem"`

		// SerialNumber is the unique serial number of the certificate.
		SerialNumber string `json:"serialNumber"`

		// SignatureAlgorithm used to sign the certificate, e.g., SHA256WITHRSA.
		SignatureAlgorithm string `json:"signatureAlgorithm"`

		// CreatedDate is the ISO-8601 date the certificate was created.
		CreatedDate time.Time `json:"createdDate"`

		// CreatedBy is the user who created the certificate.
		CreatedBy string `json:"createdBy"`

		// Description is an optional description of the certificate.
		Description *string `json:"description"`
	}

	// Validation represents the validation results for a CA Set version, its certificates or activation.
	Validation struct {
		// Warnings is a list of validation warnings.
		Warnings []Warning `json:"warnings"`
	}

	// Warning represents a data about single validation warning.
	Warning struct {
		// ContextInfo provides additional context about the warning, such as the certificate's serial number or fingerprint.
		ContextInfo map[string]any `json:"contextInfo"`

		// Type is non-navigable URI to describe the problem type.
		Type string `json:"type"`

		// Title is a human-readable summary of the problem type.
		Title string `json:"title"`

		// Detail provides a human-readable description of the warning.
		Detail string `json:"detail"`

		// Pointer indicates the specific field or element related to the warning.
		Pointer string `json:"pointer"`
	}

	// CreateCASetVersionResponse represents the response returned after creating a new CA Set version.
	CreateCASetVersionResponse CASetVersion

	// CloneCASetVersionResponse represents the response returned after cloning an existing CA Set version.
	CloneCASetVersionResponse CASetVersion

	// GetCASetVersionResponse represents the response returned when fetching details of a specific CA Set version.
	GetCASetVersionResponse CASetVersion

	// UpdateCASetVersionResponse represents the response returned after updating an existing CA Set version.
	UpdateCASetVersionResponse CASetVersion

	// ListCASetVersionsResponse represents the response containing a list of CA Set versions.
	ListCASetVersionsResponse struct {
		Versions []CASetVersion `json:"versions"`
	}

	// GetCASetVersionCertificatesResponse represents the response with certificates details of a specific CA Set version.
	GetCASetVersionCertificatesResponse struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64 `json:"version"`

		// CASetName is a descriptive name for the set.
		CASetName string `json:"caSetName"`

		// Certificates is a list of valid root or intermediate certificates. At least one is required.
		Certificates []CertificateResponse `json:"certificates"`
	}

	// CASetVersion represents a single version of a CA Set.
	CASetVersion struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// Version is the version number within the CA Set, starting at 1 and incrementing sequentially.
		Version int64 `json:"version"`

		// CASetName is a descriptive name for the set.
		CASetName string `json:"caSetName"`

		// VersionLink is the hypermedia link to the version resource.
		VersionLink string `json:"versionLink"`

		// Description is an optional description for the version.
		Description *string `json:"description"`

		// AllowInsecureSHA1 indicates whether SHA-1 certificates are allowed.
		AllowInsecureSHA1 bool `json:"allowInsecureSha1"`

		// StagingStatus is "INACTIVE" initially, changes to "ACTIVE" when active on staging.
		StagingStatus string `json:"stagingStatus"`

		// ProductionStatus is "INACTIVE" initially, changes to "ACTIVE" when active on production.
		ProductionStatus string `json:"productionStatus"`

		// CreatedDate is the creation timestamp in ISO-8601 format.
		CreatedDate time.Time `json:"createdDate"`

		// CreatedBy is the user who created the version.
		CreatedBy string `json:"createdBy"`

		// ModifiedDate is the last update timestamp in ISO-8601 format or null for new versions.
		ModifiedDate *time.Time `json:"modifiedDate"`

		// ModifiedBy is the user who last updated the version or null for new versions.
		ModifiedBy *string `json:"modifiedBy"`

		// Certificates is a list of valid root or intermediate certificates. At least one is required.
		Certificates []CertificateResponse `json:"certificates"`

		// Validation contains any warnings or errors related to the CA Set version.
		Validation *Validation `json:"validation"`

		// RemovalDate is the time when the CA set will be permanently deleted from the system. The value is null when the CA set is not scheduled for deletion.
		RemovalDate *time.Time `json:"removalDate"`

		// CASetVersionStatus indicates the CA set version status, one of "NOT_DELETED" or "DELETED".
		CASetVersionStatus string `json:"caSetVersionStatus"`
	}
)

var (
	// ErrCreateCASetVersion represents an error when creating a CA set version fails.
	ErrCreateCASetVersion = errors.New("creating a CA set version")
	// ErrCloneCASetVersion represents an error when cloning a CA set version fails.
	ErrCloneCASetVersion = errors.New("cloning a CA set version")
	// ErrGetCASetVersion represents an error when fetching a CA set version fails.
	ErrGetCASetVersion = errors.New("fetching a CA set version")
	// ErrListCASetVersions represents an error when fetching CA set versions fails.
	ErrListCASetVersions = errors.New("fetching CA set versions")
	// ErrGetCASetVersionCertificates represents an error when fetching certificates for a CA set version fails.
	ErrGetCASetVersionCertificates = errors.New("fetching certificates for a CA set version")
	// ErrUpdateCASetVersion represents an error when updating a CA set version fails.
	ErrUpdateCASetVersion = errors.New("updating a CA set version")
	// ErrDeleteCASetVersion represents an error when marking a CA set version for deletion fails.
	ErrDeleteCASetVersion = errors.New("deleting a CA set version")
)

// CertificateStatus represents the state of certificates in a CA set version.
type CertificateStatus string

// VersionStatus represents the state of CA set version on the network.
type VersionStatus string

const (
	// ExpiringCert represents certificates that are about to expire within the provided threshold.
	ExpiringCert CertificateStatus = "EXPIRING"
	// ExpiredCert represents certificates that have already expired.
	ExpiredCert CertificateStatus = "EXPIRED"
	// ExpiredOrExpiringCert represents a status filter that matches certificates that are either expiring or expired.
	ExpiredOrExpiringCert CertificateStatus = "EXPIRING,EXPIRED"
	// ExpiringOrExpiredCert represents a status filter that matches certificates that are either expired or expiring.
	ExpiringOrExpiredCert CertificateStatus = "EXPIRED,EXPIRING"
	// ActiveCert represents certificates that are ACTIVE.
	ActiveCert CertificateStatus = "ACTIVE"
	// ActiveOrExpiredCert represents certificates that are active or already expired.
	ActiveOrExpiredCert CertificateStatus = "ACTIVE,EXPIRED"
	// ExpiredOrActiveCert represents certificates that are already expired or active.
	ExpiredOrActiveCert CertificateStatus = "EXPIRED,ACTIVE"
)

// Validate validates a CreateCASetVersionRequest.
func (v CreateCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID":      validation.Validate(v.CASetID, validation.Required),
		"Description":  validation.Validate(v.Body.Description, validation.NilOrNotEmpty, validation.Length(1, 255)),
		"Certificates": validation.Validate(v.Body.Certificates, validation.Required, validation.Each(certificateValidationRules())),
	})
}

// Validate validates a ListCASetVersionsRequest.
func (v ListCASetVersionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID":              validation.Validate(v.CASetID, validation.Required),
		"CASetVersionStatuses": validation.Validate(v.CASetVersionStatuses, validation.By(caSetVersionStatusRule)),
	})
}

// AllCASetVersionStatuses returns list of all allowed CA set version status values.
func AllCASetVersionStatuses() []string {
	return []string{CASetStatusNotDeleted, CASetStatusDeleted}
}

func caSetVersionStatusRule(value any) error {
	statuses, ok := value.([]string)
	if !ok {
		return fmt.Errorf("expected []string, got %T", value)
	}

	allCASetStatuses := AllCASetVersionStatuses()
	for _, s := range statuses {
		if !slices.Contains(allCASetStatuses, s) {
			return fmt.Errorf("list element '%s' is invalid. Each element must be one of: %s",
				s, "'"+strings.Join(allCASetStatuses, "', '")+"'")
		}
	}

	return nil
}

// Validate validates a UpdateCASetVersionRequest.
func (v UpdateCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID":      validation.Validate(v.CASetID, validation.Required),
		"Version":      validation.Validate(v.Version, validation.Required),
		"Description":  validation.Validate(v.Body.Description, validation.NilOrNotEmpty, validation.Length(1, 255)),
		"Certificates": validation.Validate(v.Body.Certificates, validation.Required, validation.Each(certificateValidationRules())),
	})
}

// certificateValidationRules defines validation rules for CA set certificates.
func certificateValidationRules() validation.Rule {
	return validation.By(func(val interface{}) error {
		cert, ok := val.(CertificateRequest)
		if !ok {
			return validation.NewError("validation", "invalid certificate type")
		}
		return validation.Errors{
			"CertificatePEM": validation.Validate(cert.CertificatePEM, validation.Required),
			"Description":    validation.Validate(cert.Description, validation.NilOrNotEmpty, validation.Length(1, 255)),
		}.Filter()
	})
}

// Validate validates a CloneCASetVersionRequest.
func (v CloneCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(v.CASetID, validation.Required),
		"Version": validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a GetCASetVersionRequest.
func (v GetCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(v.CASetID, validation.Required),
		"Version": validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates a GetCASetVersionCertificatesRequest.
func (v GetCASetVersionCertificatesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(v.CASetID, validation.Required),
		"Version": validation.Validate(v.Version, validation.Required),
		"CertificateStatus": validation.Validate(v.CertificateStatus,
			validation.When(
				v.CertificateStatus != nil,
				validation.In(ExpiringCert, ExpiredCert, ExpiredOrExpiringCert, ExpiringOrExpiredCert, ActiveCert, ActiveOrExpiredCert, ExpiredOrActiveCert).Error(fmt.Sprintf(
					"value must be one of: '%s', '%s', '%s', '%s', '%s', '%s', or '%s'",
					ExpiringCert,
					ExpiredCert,
					ExpiredOrExpiringCert,
					ExpiringOrExpiredCert,
					ActiveCert,
					ActiveOrExpiredCert,
					ExpiredOrActiveCert,
				)),
			),
		),
		"ExpiryThresholdInDays": validation.Validate(v.ExpiryThresholdInDays, validation.When(v.ExpiryThresholdInDays != nil,
			validation.By(validCertificateStatusForExpiryThreshold(v.CertificateStatus, []CertificateStatus{ExpiringCert, ExpiredCert, ExpiredOrExpiringCert})))),
		"ExpiryThresholdTimestamp": validation.Validate(v.ExpiryThresholdTimestamp, validation.When(!v.ExpiryThresholdTimestamp.IsZero(),
			validation.By(validCertificateStatusForExpiryThreshold(v.CertificateStatus, []CertificateStatus{ExpiringCert, ExpiredCert})),
			validation.By(bothExpiryThresholdProvided(v.ExpiryThresholdInDays))),
			validation.By(validCertificateStatusForDate(v.CertificateStatus, v.ExpiryThresholdTimestamp)),
		),
	})
}

func validCertificateStatusForExpiryThreshold(status *CertificateStatus, allowedStatuses []CertificateStatus) validation.RuleFunc {
	return func(_ interface{}) error {
		if status == nil {
			return fmt.Errorf("CertificateStatus must be provided with this field")
		}

		for _, allowedStatus := range allowedStatuses {
			if *status == allowedStatus {
				return nil
			}
		}
		return fmt.Errorf("with this field CertificateStatus must be one of: %v", allowedStatuses)
	}
}

func bothExpiryThresholdProvided(expiryThresholdInDays *int) validation.RuleFunc {
	return func(_ interface{}) error {
		if expiryThresholdInDays != nil {
			return fmt.Errorf("ExpiryThresholdInDays cannot be used with ExpiryThresholdTimestamp")
		}
		return nil
	}
}

func validCertificateStatusForDate(status *CertificateStatus, timestamp time.Time) validation.RuleFunc {
	return func(_ interface{}) error {
		if status == nil {
			// Is already validated by validCertificateStatusForExpiryThreshold
			return nil
		}
		if timestamp.IsZero() {
			// Init for other validation rules, skip checking
			return nil
		}
		if *status == ExpiredCert && timestamp.After(time.Now()) {
			return fmt.Errorf("ExpiryThresholdTimestamp cannot be in the future for 'EXPIRED' CertificateStatus")
		}
		if *status == ExpiringCert && timestamp.Before(time.Now()) {
			return fmt.Errorf("ExpiryThresholdTimestamp cannot be in the past for 'EXPIRING' CertificateStatus")
		}
		return nil
	}
}

func (m *mtlstruststore) CreateCASetVersion(ctx context.Context, params CreateCASetVersionRequest) (*CreateCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("CreateCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-edge-truststore/v2/ca-sets/%s/versions",
		params.CASetID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateCASetVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCASetVersion, err)
	}

	var result CreateCASetVersionResponse
	resp, err := m.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) CloneCASetVersion(ctx context.Context, params CloneCASetVersionRequest) (*CloneCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("CloneCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCloneCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/clone",
		params.CASetID,
		params.Version),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCloneCASetVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCloneCASetVersion, err)
	}

	var result CloneCASetVersionResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCloneCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) GetCASetVersion(ctx context.Context, params GetCASetVersionRequest) (*GetCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-edge-truststore/v2/ca-sets/%s/versions/%d",
		params.CASetID,
		params.Version),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCASetVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCASetVersion, err)
	}

	var result GetCASetVersionResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) ListCASetVersions(ctx context.Context, params ListCASetVersionsRequest) (*ListCASetVersionsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASetVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASetVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-edge-truststore/v2/ca-sets/%s/versions", params.CASetID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASetVersions, err)
	}

	query := url.Values{}
	if params.IncludeCertificates {
		query.Set("includeCertificates", strconv.FormatBool(params.IncludeCertificates))
	}

	if params.ActiveVersionsOnly {
		query.Set("activeVersionsOnly", strconv.FormatBool(params.ActiveVersionsOnly))
	}

	if len(params.CASetVersionStatuses) > 0 {
		query.Set("caSetVersionStatus", texts.JoinStringBased(params.CASetVersionStatuses, ","))
	}

	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASetVersions, err)
	}

	var result ListCASetVersionsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASetVersions, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) GetCASetVersionCertificates(ctx context.Context, params GetCASetVersionCertificatesRequest) (*GetCASetVersionCertificatesResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetCASetVersionCertificates")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCASetVersionCertificates, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-edge-truststore/v2/ca-sets/%s/versions/%d/certificates",
		params.CASetID,
		params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCASetVersionCertificates, err)
	}

	query := url.Values{}
	if params.CertificateStatus != nil {
		query.Set("certificateStatus", string(*params.CertificateStatus))
	}

	if params.ExpiryThresholdInDays != nil {
		query.Set("expiryThresholdInDays", strconv.Itoa(*params.ExpiryThresholdInDays))
	}

	if !params.ExpiryThresholdTimestamp.IsZero() {
		query.Set("expiryThresholdTimestamp", params.ExpiryThresholdTimestamp.Format(time.RFC3339Nano))
	}

	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCASetVersionCertificates, err)
	}

	var result GetCASetVersionCertificatesResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCASetVersionCertificates, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) UpdateCASetVersion(ctx context.Context, params UpdateCASetVersionRequest) (*UpdateCASetVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("UpdateCASetVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateCASetVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/versions/%d", params.CASetID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateCASetVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateCASetVersion, err)
	}

	var result UpdateCASetVersionResponse

	resp, err := m.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

// Validate validates a DeleteCASetVersionRequest.
func (v DeleteCASetVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(v.CASetID, validation.Required),
		"Version": validation.Validate(v.Version, validation.Required),
	})
}

func (m *mtlstruststore) DeleteCASetVersion(ctx context.Context, params DeleteCASetVersionRequest) error {
	logger := m.Log(ctx)
	logger.Debug("DeleteCASetVersion")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteCASetVersion, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/mtls-edge-truststore/v2/ca-sets/%s/versions/%d", params.CASetID, params.Version).
		Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteCASetVersion, err)
	}

	resp, err := m.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteCASetVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return m.Error(resp)
	}

	return nil
}
