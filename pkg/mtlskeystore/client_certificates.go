package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListClientCertificatesResponse is the response from the ListClientCertificates API.
	ListClientCertificatesResponse struct {
		// Certificates is the list of client certificates.
		Certificates []Certificate `json:"certificates"`
	}

	// Geography represents the geography type: `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.
	Geography string

	// CryptographicAlgorithm represents the cryptographic algorithm type: `RSA` or `ECDSA`.
	CryptographicAlgorithm string

	// SecureNetwork represents the secure network type: `STANDARD_TLS` or `ENHANCED_TLS`.
	SecureNetwork string

	// Signer represents the signer type: `AKAMAI` or `THIRD_PARTY`.
	Signer string

	// Certificate is a single client certificate.
	Certificate struct {
		// CertificateID is the unique ID of the client certificate.
		CertificateID int64 `json:"certificateId"`
		// CertificateName is the name of the client certificate.
		CertificateName string `json:"certificateName"`
		// CreatedBy is the user who created the client certificate.
		CreatedBy string `json:"createdBy"`
		// CreatedDate is the creation timestamp in ISO-8601 format.
		CreatedDate time.Time `json:"createdDate"`
		// Geography is the type of network to deploy the client certificate. Either `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.
		Geography Geography `json:"geography"`
		// KeyAlgorithm is the cryptographic algorithm used for key generation, either `RSA` or `ECDSA`.
		KeyAlgorithm CryptographicAlgorithm `json:"keyAlgorithm"`
		// NotificationEmails is the list of email addresses to notify for client certificate related issues.
		NotificationEmails []string `json:"notificationEmails"`
		// SecureNetwork is the indicator of the network deployment type. Either `STANDARD_TLS` or `ENHANCED_TLS`.
		SecureNetwork SecureNetwork `json:"secureNetwork"`
		// Signer is the signing entity of the client certificate. Either `AKAMAI` or `THIRD_PARTY`.
		Signer Signer `json:"signer"`
		// Subject of the client certificate.
		Subject string `json:"subject"`
	}

	// GetClientCertificateRequest is the request to get a client certificate.
	GetClientCertificateRequest struct {
		// CertificateID is the unique ID of the client certificate.
		CertificateID int64 `json:"certificateId"`
	}

	// GetClientCertificateResponse is the response from the GetClientCertificate API.
	GetClientCertificateResponse Certificate

	// CreateClientCertificateRequest is the request to create a client certificate.
	CreateClientCertificateRequest struct {
		// CertificateName is the name of the client certificate.
		CertificateName string `json:"certificateName"`
		// ContractID is the contract ID assigned to the client certificate.
		ContractID string `json:"contractId"`
		// Geography is the type of network to deploy the client certificate. Either `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.
		Geography Geography `json:"geography"`
		// GroupID is the group ID assigned to the client certificate.
		GroupID int64 `json:"groupId"`
		// KeyAlgorithm is the cryptographic algorithm used for key generation, either `RSA` or `ECDSA`. If not provided, the API returns 'RSA' by default. Optional.
		KeyAlgorithm *CryptographicAlgorithm `json:"keyAlgorithm,omitempty"`
		// NotificationEmails is the list of email addresses to notify client certificate related issues.
		NotificationEmails []string `json:"notificationEmails"`
		// PreferredCA is the common name of the account CA certificate selected to sign the client certificate.
		// This field could be `nil` if you want to add this later.
		PreferredCA *string `json:"preferredCa,omitempty"`
		// SecureNetwork is the indicator of the network deployment type. Either `STANDARD_TLS` or `ENHANCED_TLS`.
		SecureNetwork SecureNetwork `json:"secureNetwork"`
		// Signer is the signing entity of the client certificate. Either `AKAMAI` or `THIRD_PARTY`.
		Signer Signer `json:"signer"`
		// Subject of the certificate.
		// When null, the subject is constructed with this format: /C=US/O=Akamai Technologies, Inc./OU={vcdId} {contractId} {groupId}/CN={certificateName}/.
		Subject *string `json:"subject,omitempty"`
	}

	// CreateClientCertificateResponse is the response from the CreateClientCertificate API.
	CreateClientCertificateResponse Certificate

	// PatchClientCertificateRequest is the request to update the client certificate's name or notification emails.
	PatchClientCertificateRequest struct {
		// CertificateID is the unique ID of the client certificate.
		CertificateID int64 `json:"certificateId"`
		Body          PatchClientCertificateRequestBody
	}

	// PatchClientCertificateRequestBody is the body of the request to update the client certificate's name or notification emails.
	PatchClientCertificateRequestBody struct {
		// CertificateName is the name of the client certificate.
		// This field could be `nil` if you want to add this later.
		CertificateName *string `json:"certificateName"`
		// NotificationEmails is the list of email addresses to notify client certificate related issues.
		// This field could be `nil` if you want to add this later.
		NotificationEmails []string `json:"notificationEmails"`
	}
)

const (
	// GeographyCore represents the core geography.
	GeographyCore Geography = "CORE"
	// GeographyRussiaAndCore represents the Russia and core geography.
	GeographyRussiaAndCore Geography = "RUSSIA_AND_CORE"
	// GeographyChinaAndCore represents the China and core geography.
	GeographyChinaAndCore Geography = "CHINA_AND_CORE"
	// KeyAlgorithmRSA represents the RSA key algorithm.
	KeyAlgorithmRSA CryptographicAlgorithm = "RSA"
	// KeyAlgorithmECDSA represents the ECDSA key algorithm.
	KeyAlgorithmECDSA CryptographicAlgorithm = "ECDSA"
	// SecureNetworkStandardTLS represents the standard TLS secure network.
	SecureNetworkStandardTLS SecureNetwork = "STANDARD_TLS"
	// SecureNetworkEnhancedTLS represents the enhanced TLS secure network.
	SecureNetworkEnhancedTLS SecureNetwork = "ENHANCED_TLS"
	// SignerAkamai represents the Akamai signer.
	SignerAkamai Signer = "AKAMAI"
	// SignerThirdParty represents the third-party signer.
	SignerThirdParty Signer = "THIRD_PARTY"
)

var (
	// ErrListClientCertificates is the error returned when the ListClientCertificates API fails.
	ErrListClientCertificates = errors.New("list client certificates")
	// ErrGetClientCertificate is the error returned when the GetClientCertificate API fails.
	ErrGetClientCertificate = errors.New("get client certificate")
	// ErrCreateClientCertificate is the error returned when the CreateClientCertificate API fails.
	ErrCreateClientCertificate = errors.New("create client certificate")
	// ErrPatchClientCertificate is the error returned when the PatchClientCertificate API fails.
	ErrPatchClientCertificate = errors.New("patch client certificate")
)

// Validate validates the SecureNetwork.
func (n SecureNetwork) Validate() validation.InRule {
	return validation.In(SecureNetworkStandardTLS, SecureNetworkEnhancedTLS).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'",
			n, SecureNetworkStandardTLS, SecureNetworkEnhancedTLS))
}

// Validate validates the Geography.
func (g Geography) Validate() validation.InRule {
	return validation.In(GeographyCore, GeographyRussiaAndCore, GeographyChinaAndCore).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s'",
			g, GeographyCore, GeographyRussiaAndCore, GeographyChinaAndCore))
}

// Validate validates the Signer.
func (s Signer) Validate() validation.InRule {
	return validation.In(SignerAkamai, SignerThirdParty).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'",
			s, SignerAkamai, SignerThirdParty))
}

// Validate validates the GetClientCertificateRequest.
func (r GetClientCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
	})
}

// Validate validates the CreateClientCertificateRequest.
func (r CreateClientCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateName":    validation.Validate(r.CertificateName, validation.Required, validation.Length(1, 64)),
		"ContractID":         validation.Validate(r.ContractID, validation.Required),
		"Geography":          validation.Validate(r.Geography, validation.Required, r.Geography.Validate()),
		"GroupID":            validation.Validate(r.GroupID, validation.Required),
		"KeyAlgorithm":       validation.Validate(r.KeyAlgorithm, validation.By(keyAlgorithmValidate)),
		"NotificationEmails": validation.Validate(r.NotificationEmails, validation.Required),
		"SecureNetwork":      validation.Validate(r.SecureNetwork, validation.Required, r.SecureNetwork.Validate()),
		"Signer":             validation.Validate(r.Signer, validation.Required, r.Signer.Validate()),
	})
}

// Validate validates the PatchClientCertificateRequest.
func (r PatchClientCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CertificateID": validation.Validate(r.CertificateID, validation.Required),
		"Body":          validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates the PatchClientCertificateRequestBody.
func (r PatchClientCertificateRequestBody) Validate() error {
	if r.CertificateName == nil && r.NotificationEmails == nil {
		return fmt.Errorf("CertificateName or NotificationEmails must be provided")
	} else if r.CertificateName != nil {
		return validation.Errors{
			"CertificateName": validation.Validate(r.CertificateName, validation.Required.Error(
				"value is invalid"), validation.Length(1, 64).Error(
				fmt.Sprintf("value '%s' is invalid. Must be between 1 and 64 characters", *r.CertificateName))),
		}.Filter()
	}

	return nil
}

// keyAlgorithmValidate validates the KeyAlgorithm field.
func keyAlgorithmValidate(value interface{}) error {
	v, ok := value.(*CryptographicAlgorithm)
	if v == nil {
		return nil
	}
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	switch *v {
	case "RSA", "ECDSA":
		return nil
	default:
		return fmt.Errorf("value '%s' is invalid. Must be one of: 'RSA', 'ECDSA'", *v)
	}
}

func (m *mtlskeystore) ListClientCertificates(ctx context.Context) (*ListClientCertificatesResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListClientCertificates")

	uri, err := url.Parse("/mtls-origin-keystore/v1/client-certificates")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListClientCertificates, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListClientCertificates, err)
	}

	var result ListClientCertificatesResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListClientCertificates, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListClientCertificates, m.Error(resp))
	}

	return &result, nil
}

func (m *mtlskeystore) GetClientCertificate(ctx context.Context, params GetClientCertificateRequest) (*GetClientCertificateResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetClientCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetClientCertificate, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-origin-keystore/v1/client-certificates/%d", params.CertificateID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetClientCertificate, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetClientCertificate, err)
	}

	var result GetClientCertificateResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetClientCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetClientCertificate, m.Error(resp))
	}

	return &result, nil
}

func (m *mtlskeystore) CreateClientCertificate(ctx context.Context, params CreateClientCertificateRequest) (*CreateClientCertificateResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("CreateClientCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateClientCertificate, ErrStructValidation, err)
	}

	uri, err := url.Parse("/mtls-origin-keystore/v1/client-certificates")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateClientCertificate, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateClientCertificate, err)
	}

	var result CreateClientCertificateResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateClientCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateClientCertificate, m.Error(resp))
	}

	return &result, nil
}

func (m *mtlskeystore) PatchClientCertificate(ctx context.Context, params PatchClientCertificateRequest) error {
	logger := m.Log(ctx)
	logger.Debug("PatchClientCertificate")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrPatchClientCertificate, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-origin-keystore/v1/client-certificates/%d", params.CertificateID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrPatchClientCertificate, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrPatchClientCertificate, err)
	}

	resp, err := m.Exec(req, nil, params.Body)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrPatchClientCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", ErrPatchClientCertificate, m.Error(resp))
	}

	return nil
}
