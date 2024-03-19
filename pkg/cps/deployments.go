package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Deployments is a CPS deployments API interface
	Deployments interface {
		// ListDeployments fetches deployments for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployments
		ListDeployments(context.Context, ListDeploymentsRequest) (*ListDeploymentsResponse, error)

		// GetProductionDeployment fetches production deployment for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployments-production
		GetProductionDeployment(context.Context, GetDeploymentRequest) (*GetProductionDeploymentResponse, error)

		// GetStagingDeployment fetches staging deployment for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployment-staging
		GetStagingDeployment(context.Context, GetDeploymentRequest) (*GetStagingDeploymentResponse, error)
	}

	// ListDeploymentsRequest contains parameters for ListDeployments
	ListDeploymentsRequest struct {
		EnrollmentID int
	}

	// GetDeploymentRequest contains parameters for GetProductionDeployment and GetStagingDeployment
	GetDeploymentRequest struct {
		EnrollmentID int
	}

	// ListDeploymentsResponse contains response for ListDeployments
	ListDeploymentsResponse struct {
		Production *Deployment `json:"production"`
		Staging    *Deployment `json:"staging"`
	}

	// GetProductionDeploymentResponse contains response for GetProductionDeployment
	GetProductionDeploymentResponse Deployment

	// GetStagingDeploymentResponse contains response for GetStagingDeployment
	GetStagingDeploymentResponse Deployment

	// Deployment represents details of production or staging deployment
	Deployment struct {
		OCSPStapled              *bool                          `json:"ocspStapled"`
		OCSPURIs                 []string                       `json:"ocspUris"`
		NetworkConfiguration     DeploymentNetworkConfiguration `json:"networkConfiguration"`
		PrimaryCertificate       DeploymentCertificate          `json:"primaryCertificate"`
		MultiStackedCertificates []DeploymentCertificate        `json:"multiStackedCertificates"`
	}

	// DeploymentCertificate represents certificate in context of deployment operation
	DeploymentCertificate struct {
		Certificate        string `json:"certificate"`
		Expiry             string `json:"expiry"`
		KeyAlgorithm       string `json:"keyAlgorithm"`
		SignatureAlgorithm string `json:"signatureAlgorithm"`
		TrustChain         string `json:"trustChain"`
	}

	// DeploymentNetworkConfiguration represents network configuration in context of deployment operation
	DeploymentNetworkConfiguration struct {
		Geography             string   `json:"geography"`
		MustHaveCiphers       string   `json:"mustHaveCiphers"`
		OCSPStapling          string   `json:"ocspStapling"`
		PreferredCiphers      string   `json:"preferredCiphers"`
		QUICEnabled           bool     `json:"quicEnabled"`
		SecureNetwork         string   `json:"secureNetwork"`
		SNIOnly               bool     `json:"sniOnly"`
		DisallowedTLSVersions []string `json:"disallowedTlsVersions"`
		DNSNames              []string `json:"dnsNames"`
	}
)

// Validate validates ListDeploymentsRequest
func (c ListDeploymentsRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate validates GetDeploymentsRequest
func (c GetDeploymentRequest) Validate() error {
	return validation.Errors{
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
	}.Filter()
}

var (
	// ErrListDeployments is returned when ListDeployments fails
	ErrListDeployments = errors.New("list deployments")
	// ErrGetProductionDeployment is returned when GetProductionDeployment fails
	ErrGetProductionDeployment = errors.New("get production deployment")
	// ErrGetStagingDeployment is returned when GetStagingDeployment fails
	ErrGetStagingDeployment = errors.New("get staging deployment")
)

func (c *cps) ListDeployments(ctx context.Context, params ListDeploymentsRequest) (*ListDeploymentsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListDeployments")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListDeployments, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/deployments", params.EnrollmentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListDeployments, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.deployments.v8+json")

	var result ListDeploymentsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListDeployments, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListDeployments, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetProductionDeployment(ctx context.Context, params GetDeploymentRequest) (*GetProductionDeploymentResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetProductionDeployment")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProductionDeployment, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/deployments/production", params.EnrollmentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProductionDeployment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.deployment.v8+json")

	var result GetProductionDeploymentResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProductionDeployment, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProductionDeployment, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) GetStagingDeployment(ctx context.Context, params GetDeploymentRequest) (*GetStagingDeploymentResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetStagingDeployment")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetStagingDeployment, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/deployments/staging", params.EnrollmentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetStagingDeployment, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.deployment.v8+json")

	var result GetStagingDeploymentResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetStagingDeployment, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetStagingDeployment, c.Error(resp))
	}

	return &result, nil
}
