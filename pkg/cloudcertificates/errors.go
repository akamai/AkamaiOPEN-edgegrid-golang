// Package cloudcertificates provides access to the Akamai Cloud Certificates Manager API.
package cloudcertificates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/errs"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")

	// ErrPatchCertificate is returned when patching certificate fails.
	ErrPatchCertificate = errors.New("patching certificate")

	// ErrUpdateCertificate is returned when updating certificate fails.
	ErrUpdateCertificate = errors.New("updating certificate")

	// ErrListCertificates is returned when listing certificates fails.
	ErrListCertificates = errors.New("listing certificates")

	// ErrCreateCertificate is returned when CreateCertificate fails.
	ErrCreateCertificate = errors.New("creating certificate")

	// ErrGetCertificate is returned when GetCertificate fails.
	ErrGetCertificate = errors.New("getting certificate")

	// ErrDeleteCertificate is returned when DeleteCertificate fails.
	ErrDeleteCertificate = errors.New("deleting certificate")

	// ErrListCertificateBindings is returned when ListCertificateBindings fails.
	ErrListCertificateBindings = errors.New("listing certificate bindings")

	// ErrListBindings is returned when ListBindings fails.
	ErrListBindings = errors.New("listing bindings")
)

var (
	// ErrCertificateNameInUse represents error when the certificate name is already in use.
	ErrCertificateNameInUse = &Error{Type: "/error-types/certificate-name-already-in-use"}

	// ErrCertificateNotFound represents error when the certificate is not found.
	ErrCertificateNotFound = &Error{Type: "/error-types/certificate-not-found"}

	// ErrCertificateResourceNotFound represents error when the certificate resource is not found.
	ErrCertificateResourceNotFound = &Error{Type: "/error-types/certificate-resource-not-found"}
)

type (
	// Error represents an error returned by the Cloud Certificates API.
	Error struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`

		CertificateIdentifier      string           `json:"certificateIdentifier,omitempty"`
		CertificateIdentifierValue string           `json:"certificateIdentifierValue,omitempty"`
		Data                       *ValidationData  `json:"data,omitempty"`
		Explanation                string           `json:"explanation,omitempty"`
		Errors                     []SecondaryError `json:"errors,omitempty"`
		InvalidParameterValue      string           `json:"invalidParameterValue,omitempty"`
		ParameterName              string           `json:"parameterName,omitempty"`
	}

	// SecondaryError represents detailed error information for validation failures.
	SecondaryError struct {
		Type     string `json:"type,omitempty"`
		Title    string `json:"title,omitempty"`
		Detail   string `json:"detail,omitempty"`
		Instance string `json:"instance,omitempty"`

		Explanation           string   `json:"explanation,omitempty"`
		InvalidParameterValue []string `json:"invalidParameterValue,omitempty"`
		ParameterName         string   `json:"parameterName,omitempty"`
	}

	// ValidationData contains details about certificate and trust chain validation.
	ValidationData struct {
		SignedCertificatePEM string            `json:"signedCertificatePem,omitempty"`
		SignedCertificates   []PEMValidation   `json:"signedCertificates,omitempty"`
		TrustChain           []PEMValidation   `json:"trustChain,omitempty"`
		TrustChainPEM        *string           `json:"trustChainPem,omitempty"`
		Validation           *ValidationResult `json:"validation,omitempty"`
	}

	// PEMValidation represents a PEM certificate with validation details.
	PEMValidation struct {
		CertificatePEM     string            `json:"certificatePem"`
		CreatedBy          *string           `json:"createdBy,omitempty"`
		CreatedDate        *string           `json:"createdDate,omitempty"`
		DisplayName        *string           `json:"displayName,omitempty"`
		EndDate            *string           `json:"endDate,omitempty"`
		Fingerprint        *string           `json:"fingerprint,omitempty"`
		Issuer             *string           `json:"issuer,omitempty"`
		SerialNumber       *string           `json:"serialNumber,omitempty"`
		SignatureAlgorithm *string           `json:"signatureAlgorithm,omitempty"`
		StartDate          *string           `json:"startDate,omitempty"`
		Subject            *Subject          `json:"subject,omitempty"`
		Validation         *ValidationResult `json:"validation,omitempty"`
	}

	// ValidationResult contains validation results for certificate operations.
	ValidationResult struct {
		Errors   []ValidationDetail `json:"errors,omitempty"`
		Notices  []ValidationDetail `json:"notices,omitempty"`
		Warnings []ValidationDetail `json:"warnings,omitempty"`
	}

	// ValidationDetail provides details about a validation message.
	ValidationDetail struct {
		Detail   string `json:"detail,omitempty"`
		Instance string `json:"instance,omitempty"`
		Message  string `json:"message,omitempty"`
		Name     string `json:"name,omitempty"`
		Status   *int   `json:"status,omitempty"`
		Title    string `json:"title,omitempty"`
		Type     string `json:"type,omitempty"`
	}
)

// Error parses an error from the CloudCertificates API response.
func (c *cloudcertificates) Error(r *http.Response) error {
	var e Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. CCM API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.Status = r.StatusCode

	return &e
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s ", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons.
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}

	ignoreType := t.Type == ""
	ignoreStatus := t.Status == 0
	ignoreTitle := t.Title == ""
	matchType := t.Type == e.Type
	matchStatus := t.Status == e.Status
	matchTitle := t.Title == e.Title

	return (matchType || ignoreType) &&
		(matchStatus || ignoreStatus) &&
		(matchTitle || ignoreTitle)
}
