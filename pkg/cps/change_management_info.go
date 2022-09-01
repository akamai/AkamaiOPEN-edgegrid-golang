package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// ChangeManagementInfo is a CPS API enabling change management
	ChangeManagementInfo interface {
		// GetChangeManagementInfo gets information about acknowledgement status,
		// and may include warnings about potential conflicts that may occur if you proceed with acknowledgement
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeManagementInfo(ctx context.Context, params GetChangeRequest) (*ChangeManagementInfoResponse, error)

		// GetChangeDeploymentInfo gets deployment currently deployed to the staging network
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeDeploymentInfo(ctx context.Context, params GetChangeRequest) (*ChangeDeploymentInfoResponse, error)
	}

	// ChangeManagementInfoResponse contains response from GetChangeManagementInfo
	ChangeManagementInfoResponse struct {
		AcknowledgementDeadline *string           `json:"acknowledgementDeadline"`
		ValidationResultHash    string            `json:"validationResultHash"`
		PendingState            PendingState      `json:"pendingState"`
		ValidationResult        *ValidationResult `json:"validationResult"`
	}

	// PendingState contains the snapshot of the pending state for the enrollment
	PendingState struct {
		PendingCertificates         []PendingCertificate        `json:"pendingCertificates"`
		PendingNetworkConfiguration PendingNetworkConfiguration `json:"pendingNetworkConfiguration"`
	}

	// PendingCertificate contains the snapshot of the pending certificate for the enrollment
	PendingCertificate struct {
		CertificateType    string   `json:"certificateType"`
		FullCertificate    string   `json:"fullCertificate"`
		OCSPStapled        string   `json:"ocspStapled"`
		OCSPURIs           []string `json:"ocspUris"`
		SignatureAlgorithm string   `json:"signatureAlgorithm"`
		KeyAlgorithm       string   `json:"keyAlgorithm"`
	}

	// PendingNetworkConfiguration contains the snapshot of the pending network configuration for the enrollment
	PendingNetworkConfiguration struct {
		DNSNameSettings       *DNSNameSettings `json:"dnsNameSettings"`
		MustHaveCiphers       string           `json:"mustHaveCiphers"`
		NetworkType           string           `json:"networkType"`
		OCSPStapling          string           `json:"ocspStapling"`
		PreferredCiphers      string           `json:"preferredCiphers"`
		QUICEnabled           string           `json:"quicEnabled"`
		SNIOnly               string           `json:"sniOnly"`
		DisallowedTLSVersions []string         `json:"disallowedTlsVersions"`
	}

	// ValidationResult contains validation errors and warnings messages
	ValidationResult struct {
		Errors   []ValidationMessage `json:"errors"`
		Warnings []ValidationMessage `json:"warnings"`
	}

	// ValidationMessage holds validation message
	ValidationMessage struct {
		Message     string `json:"message"`
		MessageCode string `json:"messageCode"`
	}

	// ChangeDeploymentInfoResponse contains response from GetChangeDeploymentInfo
	ChangeDeploymentInfoResponse Deployment
)

var (
	// ErrGetChangeManagementInfo is returned when GetChangeManagementInfo fails
	ErrGetChangeManagementInfo = errors.New("get change management info")
	// ErrGetChangeDeploymentInfo is returned when GetChangeDeploymentInfo fails
	ErrGetChangeDeploymentInfo = errors.New("get change deployment info")
)

func (c *cps) GetChangeManagementInfo(ctx context.Context, params GetChangeRequest) (*ChangeManagementInfoResponse, error) {
	c.Log(ctx).Debug("GetChangeManagementInfo")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeManagementInfo, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/info/change-management-info",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeManagementInfo, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-management-info.v5+json")

	var result ChangeManagementInfoResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeManagementInfo, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeManagementInfo, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetChangeDeploymentInfo(ctx context.Context, params GetChangeRequest) (*ChangeDeploymentInfoResponse, error) {
	c.Log(ctx).Debug("GetChangeDeploymentInfo")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetChangeDeploymentInfo, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/input/info/change-management-info",
		params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetChangeDeploymentInfo, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.deployment.v7+json")

	var result ChangeDeploymentInfoResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetChangeDeploymentInfo, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetChangeDeploymentInfo, c.Error(resp))
	}

	return &result, nil
}
