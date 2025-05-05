package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListAccountCACertificatesRequest is a request to list account CA certificates.
	ListAccountCACertificatesRequest struct {
		// A list of statuses separated by commas used to filter account CA certificates.
		Status []CertificateStatus
	}

	// ListAccountCACertificatesResponse is a response from list account CA certificates operation.
	ListAccountCACertificatesResponse struct {
		// Certificates contains the account CA certificates.
		Certificates []AccountCACertificate `json:"certificates"`
	}

	// AccountCACertificate represents a CA certificate.
	AccountCACertificate struct {
		// AccountID is the account the CA certificate is under.
		AccountID string `json:"accountId"`
		// Certificate is the certificate block of the CA certificate.
		Certificate string `json:"certificate"`
		// CommonName is the common name of the CA certificate.
		CommonName string `json:"commonName"`
		// CreatedBy is the user who created the CA certificate.
		CreatedBy string `json:"createdBy"`
		// CreatedDate is the timestamp indicating the CA certificate's creation.
		CreatedDate time.Time `json:"createdDate"`
		// ExpiryDate is the timestamp indicating when the CA certificate expires.
		ExpiryDate time.Time `json:"expiryDate"`
		// ID is the unique identifier of the CA certificate.
		ID int64 `json:"id"`
		// IssuedDate is the timestamp indicating the CA certificate's availability.
		IssuedDate time.Time `json:"issuedDate"`
		// KeyAlgorithm identifies the CA certificate's encryption algorithm. The only currently supported value is `RSA`.
		KeyAlgorithm CryptographicAlgorithm `json:"keyAlgorithm"`
		// KeySizeInBytes is the private key length of the CA certificate.
		KeySizeInBytes int64 `json:"keySizeInBytes"`
		// QalificationDate is the timestamp indicating when the CA certificate's status moved from `QUALIFYING` to `CURRENT`.
		QualificationDate time.Time `json:"qualificationDate"`
		// SignatureAlgorithm specifies the algorithm that secures the data exchange between the edge server and origin.
		SignatureAlgorithm string `json:"signatureAlgorithm"`
		// Status is the status of the CA certificate. Either `QUALIFYING`, `CURRENT`, `PREVIOUS`, or `EXPIRED`.
		Status CertificateStatus `json:"status"`
		// Subject is the public key's entity stored in the CA certificate's subject public key field.
		Subject string `json:"subject"`
		// Version is the version of the CA certificate.
		Version int64 `json:"version"`
	}

	// CertificateStatus is the status of the CA certificate. Either `QUALIFYING`, `CURRENT`, `PREVIOUS`, or `EXPIRED`.
	CertificateStatus string
)

const (
	// CertificateStatusCurrent represents the certificate status "CURRENT".
	CertificateStatusCurrent CertificateStatus = "CURRENT"
	// CertificateStatusExpired represents the certificate status "EXPIRED".
	CertificateStatusExpired CertificateStatus = "EXPIRED"
	// CertificateStatusPrevious represents the certificate status "PREVIOUS".
	CertificateStatusPrevious CertificateStatus = "PREVIOUS"
	// CertificateStatusQualifying represents the certificate status "QUALIFYING".
	CertificateStatusQualifying CertificateStatus = "QUALIFYING"
)

// ErrListAccountCACertificates is returned when the ListAccountCACertificates operation fails.
var ErrListAccountCACertificates = errors.New("list account ca certificates")

// Validate validates the ListAccountCACertificatesRequest.
func (r ListAccountCACertificatesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Status": validation.Validate(r.Status, validation.By(statusRule)),
	})
}

func statusRule(value any) error {
	status, ok := value.([]CertificateStatus)
	if !ok {
		return fmt.Errorf("expected []CertificateStatus, got %T", value)
	}

	for _, s := range status {
		if s != CertificateStatusCurrent &&
			s != CertificateStatusExpired &&
			s != CertificateStatusPrevious &&
			s != CertificateStatusQualifying {
			return fmt.Errorf("list '%s' contains invalid element '%s'. Each element must be one of: '%s', '%s', '%s', or '%s'",
				status, s, CertificateStatusCurrent, CertificateStatusExpired, CertificateStatusPrevious, CertificateStatusQualifying)
		}
	}
	return nil

}

func (m *mtlskeystore) ListAccountCACertificates(ctx context.Context, params ListAccountCACertificatesRequest) (*ListAccountCACertificatesResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListAccountCACertificates")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListAccountCACertificates, ErrStructValidation, err)
	}

	uri, err := url.Parse("/mtls-origin-keystore/v1/ca-certificates")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListAccountCACertificates, err)
	}

	if len(params.Status) > 0 {
		uri.RawQuery = statusesToQueryString(params.Status)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAccountCACertificates, err)
	}

	var result ListAccountCACertificatesResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAccountCACertificates, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAccountCACertificates, m.Error(resp))
	}

	return &result, nil
}

func statusesToQueryString(statuses []CertificateStatus) string {
	chunks := make([]string, 0, len(statuses))
	for _, s := range statuses {
		chunks = append(chunks, string(s))
	}
	q := url.Values{}
	q.Set("status", strings.Join(chunks, ","))
	return q.Encode()
}
