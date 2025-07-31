package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RotateClientCertificateVersionRequest is used to request the creation of a new client certificate version.
	RotateClientCertificateVersionRequest struct {
		// CertificateID is a unique identifier representing the client certificate.
		CertificateID int64
	}

	// ListClientCertificateVersionsRequest is used to request the client certificate versions.
	ListClientCertificateVersionsRequest struct {
		// CertificateID is a unique identifier representing the client certificate.
		CertificateID int64

		// IncludeAssociatedProperties specifies whether to include associated properties in the response.
		IncludeAssociatedProperties bool
	}

	// DeleteClientCertificateVersionRequest is used to delete a client certificate version with the provided certificateID and version.
	DeleteClientCertificateVersionRequest struct {
		// CertificateID is a unique identifier representing the client certificate.
		CertificateID int64

		// Version identifies each client certificate version.
		Version int64
	}

	// UploadSignedClientCertificateRequest represents a request to upload a signed client certificate.
	UploadSignedClientCertificateRequest struct {
		// CertificateID is a unique identifier representing the client certificate.
		CertificateID int64

		// Version identifies each client certificate version.
		Version int64

		// AcknowledgeAllWarnings specifies whether to ignore all warnings during the signed certificate's upload.
		AcknowledgeAllWarnings *bool

		// Body holds the request body for the UploadSignedClientCertificate operation.
		Body UploadSignedClientCertificateRequestBody
	}

	// UploadSignedClientCertificateRequestBody represents a request body for UploadSignedClientCertificateRequest.
	UploadSignedClientCertificateRequestBody struct {
		// Certificate is a text representation of the client certificate in PEM format after it's signed.
		Certificate string `json:"certificate"`

		// TrustChain is a text representation of the trust chain in PEM format.
		TrustChain *string `json:"trustChain,omitempty"`
	}

	// ListClientCertificateVersionsResponse represents a list of client certificate versions.
	ListClientCertificateVersionsResponse struct {
		// Versions contains the client certificate versions.
		Versions []ClientCertificateVersion `json:"versions"`
	}

	// RotateClientCertificateVersionResponse represents the new client certificate version created.
	RotateClientCertificateVersionResponse struct {
		// Version identifies each client certificate version.
		Version int64 `json:"version"`

		// VersionGUID is a unique identifier for the client certificate version.
		VersionGUID string `json:"versionGuid"`

		// VersionAlias is an alias for the client certificate version, can be either `CURRENT` or `PREVIOUS`.
		VersionAlias *string `json:"versionAlias"`

		// CertificateBlock creates a Certificate Signing Request (CSR) for new THIRD_PARTY client certificates as the signer.
		CertificateBlock *CertificateBlock `json:"certificateBlock"`

		// CreatedBy represents the user who created the client certificate version.
		CreatedBy string `json:"createdBy"`

		// CreatedDate is an ISO 8601 timestamp indicating the client certificate's creation.
		CreatedDate time.Time `json:"createdDate"`

		// CSRBlock creates a Certificate Signing Request (CSR) for new THIRD_PARTY client certificates as the signer.
		CSRBlock *CSRBlock `json:"csrBlock"`

		// ExpiryDate is an ISO 8601 timestamp indicating when the client certificate version expires.
		ExpiryDate *time.Time `json:"expiryDate"`

		// IssuedDate is an ISO 8601 timestamp indicating the client certificate version's availability.
		IssuedDate *time.Time `json:"issuedDate"`

		// Issuer represents the signing entity of the client certificate version.
		Issuer *string `json:"issuer"`

		// KeyAlgorithm identifies the client certificate version's encryption algorithm. Supported values are RSA and ECDSA.
		KeyAlgorithm string `json:"keyAlgorithm"`

		// EllipticCurve specifies the key elliptic curve when key algorithm ECDSA is used.
		EllipticCurve *string `json:"ellipticCurve"`

		// KeySizeInBytes represents the private key length of the client certificate version when key algorithm RSA is used.
		KeySizeInBytes *string `json:"keySizeInBytes"`

		// SignatureAlgorithm specifies the algorithm that secures the data exchange between the origin server and origin.
		SignatureAlgorithm *string `json:"signatureAlgorithm"`

		// Status is the client certificate version status.
		Status string `json:"status"`

		// Subject represents the public key's entity stored in the client certificate version's subject public key field.
		Subject *string `json:"subject"`
	}

	// ClientCertificateVersion represents a version of a client certificate.
	ClientCertificateVersion struct {
		// Version identifies each client certificate version.
		Version int64 `json:"version"`

		// VersionGUID is a unique identifier for the client certificate version.
		VersionGUID string `json:"versionGuid"`

		// VersionAlias is an alias for the client certificate version, can be either `CURRENT` or `PREVIOUS`.
		VersionAlias *string `json:"versionAlias"`

		// CertificateBlock creates a Certificate Signing Request (CSR) for new THIRD_PARTY client certificates as the signer.
		CertificateBlock *CertificateBlock `json:"certificateBlock"`

		// CertificateSubmittedBy represents the user who uploaded the THIRD_PARTY client certificate version.
		CertificateSubmittedBy *string `json:"certificateSubmittedBy"`

		// CertificateSubmittedDate is an ISO 8601 timestamp indicating when the THIRD_PARTY signer client certificate version was upload.
		CertificateSubmittedDate *time.Time `json:"certificateSubmittedDate"`

		// CreatedBy represents the user who created the client certificate version.
		CreatedBy string `json:"createdBy"`

		// CreatedDate is an ISO 8601 timestamp indicating the client certificate's creation.
		CreatedDate time.Time `json:"createdDate"`

		// CSRBlock creates a Certificate Signing Request (CSR) for new THIRD_PARTY client certificates as the signer.
		CSRBlock *CSRBlock `json:"csrBlock"`

		// DeleteRequestedDate is an ISO 8601 timestamp indicating the client certificate version's deletion request. Appears as null if not specified.
		DeleteRequestedDate *time.Time `json:"deleteRequestedDate"`

		// ExpiryDate is an ISO 8601 timestamp indicating when the client certificate version expires.
		ExpiryDate *time.Time `json:"expiryDate"`

		// IssuedDate is an ISO 8601 timestamp indicating the client certificate version's availability.
		IssuedDate *time.Time `json:"issuedDate"`

		// Issuer represents the signing entity of the client certificate version.
		Issuer *string `json:"issuer"`

		// KeyAlgorithm identifies the client certificate version's encryption algorithm. Supported values are RSA and ECDSA.
		KeyAlgorithm string `json:"keyAlgorithm"`

		// EllipticCurve specifies the key elliptic curve when key algorithm ECDSA is used.
		EllipticCurve *string `json:"ellipticCurve"`

		// KeySizeInBytes represents the private key length of the client certificate version when key algorithm RSA is used.
		KeySizeInBytes *string `json:"keySizeInBytes"`

		// ScheduledDeleteDate is an ISO 8601 timestamp indicating client certificate version's deletion. Appears as null if not specified.
		ScheduledDeleteDate *time.Time `json:"scheduledDeleteDate"`

		// SignatureAlgorithm specifies the algorithm that secures the data exchange between the origin server and origin.
		SignatureAlgorithm *string `json:"signatureAlgorithm"`

		// Status is the client certificate version status.
		Status string `json:"status"`

		// Subject represents the public key's entity stored in the client certificate version's subject public key field.
		Subject *string `json:"subject"`

		// Validation checks the versions when uploading THIRD_PARTY signed client certificates to Mutual TLS Origin Keystore.
		Validation ValidationResult `json:"validation"`

		// AssociatedProperties represents the properties associated with the client certificate version.
		AssociatedProperties []AssociatedProperty `json:"properties"`
	}

	// CertificateBlock represents the  Certificate Signing Request (CSR) block for THIRD_PARTY client certificates.
	CertificateBlock struct {
		// Certificate is a text representation of the client certificate in PEM format.
		Certificate string `json:"certificate"`

		// KeyAlgorithm identifies the CA certificate's encryption algorithm.
		KeyAlgorithm string `json:"keyAlgorithm"`

		// TrustChain a text representation of the trust chain in PEM format.
		TrustChain string `json:"trustChain"`
	}

	// CSRBlock represents the Certificate Signing Request (CSR) for new THIRD_PARTY client certificates as the signer.
	CSRBlock struct {
		// CSR is a text representation of the certificate signing request.
		CSR string `json:"csr"`

		// KeyAlgorithm identifies the client certificate's encryption algorithm.
		KeyAlgorithm string `json:"keyAlgorithm"`
	}

	// ValidationResult holds validation errors and warnings.
	ValidationResult struct {
		// Errors indicates validation errors you need to resolve for the request to succeed.
		Errors []ValidationDetail `json:"errors"`

		// Warnings indicates validation warnings you can resolve.
		Warnings []ValidationDetail `json:"warnings"`
	}

	// ValidationDetail represents individual validation error or warning.
	ValidationDetail struct {
		// Message specifies the error or warning details.
		Message string `json:"message"`

		// Reason specifies the error or warning root cause.
		Reason string `json:"reason"`

		// Type specifies the error or warning category.
		Type string `json:"type"`
	}

	// DeleteClientCertificateVersionResponse represents the response of delete client certificate version request.
	DeleteClientCertificateVersionResponse struct {
		// Message indicates the client certificate version's scheduled deletion date, and specifies its reuse.
		Message string `json:"message"`
	}

	// AssociatedProperty represents the property associated with the client certificate version.
	AssociatedProperty struct {
		// AssetID is the unique identifier of the property.
		AssetID int64 `json:"assetId"`

		// GroupID is the unique identifier of the property group.
		GroupID int64 `json:"groupId"`

		// PropertyName is the name of the property.
		PropertyName string `json:"propertyName"`

		// PropertyVersion is the version of the property.
		PropertyVersion int64 `json:"propertyVersion"`
	}

	// CertificateVersionStatus represents the state of client certificate version.
	CertificateVersionStatus string
)

const (
	// CertificateVersionStatusAwaitingSigned represents client certificate versions that are waiting to be signed.
	CertificateVersionStatusAwaitingSigned CertificateVersionStatus = "AWAITING_SIGNED_CERTIFICATE"

	// CertificateVersionStatusDeploymentPending represents client certificate versions that are pending deployment.
	CertificateVersionStatusDeploymentPending CertificateVersionStatus = "DEPLOYMENT_PENDING"

	// CertificateVersionStatusDeployed represents client certificate versions that are deployed.
	CertificateVersionStatusDeployed CertificateVersionStatus = "DEPLOYED"

	// CertificateVersionStatusDeletePending represents client certificate versions that are pending deletion.
	CertificateVersionStatusDeletePending CertificateVersionStatus = "DELETE_PENDING"
)

var (
	// ErrRotateClientCertificateVersion represents error when rotating a client certificate version fails.
	ErrRotateClientCertificateVersion = errors.New("rotating client certificate version")

	// ErrListClientCertificateVersions represents error when fetching a client certificate versions fails.
	ErrListClientCertificateVersions = errors.New("fetching client certificate versions")

	// ErrDeleteClientCertificateVersion represents error when deleting a client certificate version fails.
	ErrDeleteClientCertificateVersion = errors.New("deleting client certificate version")

	// ErrUploadClientCertificateVersion represents error when uploading a client certificate version fails.
	ErrUploadClientCertificateVersion = errors.New("uploading client certificate version")
)

func validateCertificateID(id int64) error {
	return validation.Validate(id, validation.Required)
}

func validateVersion(version int64) error {
	return validation.Validate(version, validation.Required)
}

// Validate validates a RotateClientCertificateVersionRequest.
func (v RotateClientCertificateVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validateCertificateID(v.CertificateID),
	})
}

// Validate validates a ListClientCertificateVersionsRequest.
func (v ListClientCertificateVersionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validateCertificateID(v.CertificateID),
	})
}

// Validate validates a DeleteClientCertificateVersionRequest.
func (v DeleteClientCertificateVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validateCertificateID(v.CertificateID),
		"Version":       validateVersion(v.Version),
	})
}

// Validate validates a UploadSignedClientCertificateRequest.
func (v UploadSignedClientCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validateCertificateID(v.CertificateID),
		"Version":       validateVersion(v.Version),
		"Body":          validation.Validate(v.Body, validation.Required),
	})
}

// Validate validates a UploadSignedClientCertificateRequestBody.
func (v UploadSignedClientCertificateRequestBody) Validate() error {
	return validation.Errors{
		"Certificate": validation.Validate(v.Certificate, validation.Required, validation.Length(1, 0)),
		"TrustChain":  validation.Validate(v.TrustChain, validation.When(v.TrustChain != nil, validation.Required)),
	}.Filter()
}

func (m *mtlskeystore) RotateClientCertificateVersion(ctx context.Context, params RotateClientCertificateVersionRequest) (*RotateClientCertificateVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("Rotating client certificate versions")

	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: validation failed: %s", ErrRotateClientCertificateVersion, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-origin-keystore/v1/client-certificates/%d/versions",
		params.CertificateID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrRotateClientCertificateVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create HTTP request: %s", ErrRotateClientCertificateVersion, err)
	}

	var result RotateClientCertificateVersionResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %s", ErrRotateClientCertificateVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlskeystore) ListClientCertificateVersions(ctx context.Context, params ListClientCertificateVersionsRequest) (*ListClientCertificateVersionsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("Fetching client certificate versions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: validation failed: %s", ErrListClientCertificateVersions, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-origin-keystore/v1/client-certificates/%d/versions",
		params.CertificateID),
	)

	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrListClientCertificateVersions, err)
	}

	if params.IncludeAssociatedProperties {
		query := url.Values{}
		query.Set("includeAssociatedProperties", strconv.FormatBool(params.IncludeAssociatedProperties))
		uri.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create HTTP request: %s", ErrListClientCertificateVersions, err)
	}

	var result ListClientCertificateVersionsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %s", ErrListClientCertificateVersions, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlskeystore) DeleteClientCertificateVersion(ctx context.Context, params DeleteClientCertificateVersionRequest) (*DeleteClientCertificateVersionResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("Deleting client certificate version")

	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: validation failed: %s", ErrDeleteClientCertificateVersion, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-origin-keystore/v1/client-certificates/%d/versions/%d",
		params.CertificateID,
		params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse URL: %s", ErrDeleteClientCertificateVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create HTTP request: %s", ErrDeleteClientCertificateVersion, err)
	}

	var result DeleteClientCertificateVersionResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %s", ErrDeleteClientCertificateVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlskeystore) UploadSignedClientCertificate(ctx context.Context, params UploadSignedClientCertificateRequest) error {
	logger := m.Log(ctx)
	logger.Debug("Uploading signed client certificate")

	err := params.Validate()
	if err != nil {
		return fmt.Errorf("%w: validation failed: %s", ErrUploadClientCertificateVersion, err)
	}

	query := url.Values{}
	if params.AcknowledgeAllWarnings != nil {
		query.Set("acknowledgeAllWarnings", strconv.FormatBool(*params.AcknowledgeAllWarnings))
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/mtls-origin-keystore/v1/client-certificates/%d/versions/%d/certificate-block",
		params.CertificateID,
		params.Version))
	if err != nil {
		return fmt.Errorf("%w: failed to parse URL: %s", ErrUploadClientCertificateVersion, err)
	}

	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create HTTP request: %s", ErrUploadClientCertificateVersion, err)
	}

	resp, err := m.Exec(req, nil, params.Body)
	if err != nil {
		return fmt.Errorf("%w: request execution failed: %s", ErrUploadClientCertificateVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return m.Error(resp)
	}

	return nil
}
