package mtlstruststore

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ValidateCertificatesRequest holds request for ValidateCertificates.
	ValidateCertificatesRequest struct {
		// AllowInsecureSHA1 is if by default all the certificates in the version should use a signature algorithm of SHA-256 or better. If "allowInsecureSha1" is set to true, TCM will allow certificates with SHA-1 signature. False by default.
		AllowInsecureSHA1 bool `json:"allowInsecureSha1"`

		// Certificates is a list of root or intermediate certificates in PEM format. The list should have at least one certificate PEM format in it.
		Certificates []ValidateCertificate `json:"certificates"`
	}

	// ValidateCertificate hold details about one certificate to be used in ValidateCertificates.
	ValidateCertificate struct {
		// CertificatePEM is certificate in PEM format.
		CertificatePEM string `json:"certificatePem"`

		// Description is optional description for the certificate.
		Description *string `json:"description,omitempty"`
	}

	// ValidateCertificatesResponse holds response for ValidateCertificates.
	ValidateCertificatesResponse struct {
		// AllowInsecureSHA1 is if by default all the certificates in the version should use a signature algorithm of SHA-256 or better. If "allowInsecureSha1" is set to true, TCM will allow certificates with SHA-1 signature.
		AllowInsecureSHA1 bool `json:"allowInsecureSha1"`

		// Certificates is a collection of certificates where each element represents the details of one certificate and validationResults where each element represents the result of validation for a given certificate. validationResults is null in case of successful validation.
		Certificates []ValidateCertificateResponse `json:"certificates"`

		// Validation is the validation result of the request.
		Validation Validation `json:"validation"`
	}

	// ValidateCertificateResponse holds response about one certificate for ValidateCertificates.
	ValidateCertificateResponse struct {
		// CertificatePEM is certificate string in PEM format.
		CertificatePEM string `json:"certificatePem"`

		// EndDate is date after which the certificate is not valid. It is represented in ISO-8601 format.
		EndDate time.Time `json:"endDate"`

		// Fingerprint is unique SHA-256 fingerprint of the certificate.
		Fingerprint string `json:"fingerprint"`

		// Issuer of the certificate.
		Issuer string `json:"issuer"`

		// SerialNumber is unique serial number of the certificate.
		SerialNumber string `json:"serialNumber"`

		// SignatureAlgorithm is signature algorithm of the certificate - ex: SHA256WITHRSA.
		SignatureAlgorithm string `json:"signatureAlgorithm"`

		// StartDate is date before which the certificate is not valid. It is represented in ISO-8601 format.
		StartDate time.Time `json:"startDate"`

		// Subject of the certificate.
		Subject string `json:"subject"`

		// Description is optional description of the certificate.
		Description *string `json:"description"`
	}
)

// Validate validates ValidateCertificatesRequest.
func (r ValidateCertificatesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Certificates": validation.Validate(r.Certificates, validation.Required),
	})
}

// Validate validates ValidateCertificate.
func (v ValidateCertificate) Validate() error {
	return validation.Errors{
		"CertificatePEM": validation.Validate(v.CertificatePEM, validation.Required),
	}.Filter()
}

func (m *mtlstruststore) ValidateCertificates(ctx context.Context, params ValidateCertificatesRequest) (*ValidateCertificatesResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ValidateCertificates")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrValidateCertificates, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/mtls-edge-truststore/v2/certificates/validate", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrValidateCertificates, err)
	}

	var result ValidateCertificatesResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrValidateCertificates, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}
