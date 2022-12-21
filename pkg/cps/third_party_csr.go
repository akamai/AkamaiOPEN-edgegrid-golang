package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ThirdPartyCSR is a CPS API enabling management of third-party certificates
	ThirdPartyCSR interface {
		// GetChangeThirdPartyCSR gets certificate signing request
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeThirdPartyCSR(ctx context.Context, params GetChangeRequest) (*ThirdPartyCSRResponse, error)

		// UploadThirdPartyCertAndTrustChain uploads signed certificate and trust chain to cps
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		UploadThirdPartyCertAndTrustChain(context.Context, UploadThirdPartyCertAndTrustChainRequest) error
	}

	// ThirdPartyCSRResponse is a response object containing list of csrs
	ThirdPartyCSRResponse struct {
		CSRs []CertSigningRequest `json:"csrs"`
	}

	// CertSigningRequest holds CSR
	CertSigningRequest struct {
		CSR          string `json:"csr"`
		KeyAlgorithm string `json:"keyAlgorithm"`
	}

	// UploadThirdPartyCertAndTrustChainRequest contains parameters to upload certificates
	UploadThirdPartyCertAndTrustChainRequest struct {
		EnrollmentID int
		ChangeID     int
		Certificates ThirdPartyCertificates
	}

	// ThirdPartyCertificates contains certificates information
	ThirdPartyCertificates struct {
		CertificatesAndTrustChains []CertificateAndTrustChain `json:"certificatesAndTrustChains"`
	}

	// CertificateAndTrustChain contains single certificate with associated trust chain
	CertificateAndTrustChain struct {
		Certificate  string `json:"certificate"`
		TrustChain   string `json:"trustChain,omitempty"`
		KeyAlgorithm string `json:"keyAlgorithm"`
	}
)

// Validate validates UploadThirdPartyCertAndTrustChainRequest
func (r UploadThirdPartyCertAndTrustChainRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"EnrollmentID": validation.Validate(r.EnrollmentID, validation.Required),
		"ChangeID":     validation.Validate(r.ChangeID, validation.Required),
		"Certificates": validation.Validate(r.Certificates, validation.Required),
	})
}

// Validate validates ThirdPartyCertificates
func (r ThirdPartyCertificates) Validate() error {
	return validation.Errors{
		"CertificatesAndTrustChains": validation.Validate(r.CertificatesAndTrustChains),
	}.Filter()
}

// Validate validates CertificateAndTrustChain
func (r CertificateAndTrustChain) Validate() error {
	return validation.Errors{
		"Certificate": validation.Validate(r.Certificate, validation.Required),
		"KeyAlgorithm": validation.Validate(r.KeyAlgorithm, validation.Required, validation.In("RSA", "ECDSA").
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'RSA', 'ECDSA'", r.KeyAlgorithm))),
	}.Filter()
}

var (
	// ErrGetChangeThirdPartyCSR is returned when GetChangeThirdPartyCSR fails
	ErrGetChangeThirdPartyCSR = errors.New("get change third-party csr")
	// ErrUploadThirdPartyCertAndTrustChain is returned when UploadThirdPartyCertAndTrustChain fails
	ErrUploadThirdPartyCertAndTrustChain = errors.New("upload third-party cert and trust chain")
)

func (c *cps) GetChangeThirdPartyCSR(ctx context.Context, params GetChangeRequest) (*ThirdPartyCSRResponse, error) {
	c.Log(ctx).Debug("GetChangeThirdPartyCSR")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeThirdPartyCSR, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/info/third-party-csr",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeThirdPartyCSR, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.csr.v2+json")

	var result ThirdPartyCSRResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeThirdPartyCSR, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeThirdPartyCSR, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) UploadThirdPartyCertAndTrustChain(ctx context.Context, params UploadThirdPartyCertAndTrustChainRequest) error {
	c.Log(ctx).Debug("UploadThirdPartyCertAndTrustChain")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrUploadThirdPartyCertAndTrustChain, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/update/third-party-cert-and-trust-chain",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrUploadThirdPartyCertAndTrustChain, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.certificate-and-trust-chain.v2+json; charset=utf-8")

	resp, err := c.Exec(req, nil, params.Certificates)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrUploadThirdPartyCertAndTrustChain, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", ErrUploadThirdPartyCertAndTrustChain, c.Error(resp))
	}

	return nil
}
