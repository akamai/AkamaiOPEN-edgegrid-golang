package papi

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

type (
	GetEdgeHostnamesRequest struct {
		ContractID string
		GroupID    string
		Options    []string
	}

	GetEdgeHostnameRequest struct {
		EdgeHostnameID string
		ContractID     string
		GroupID        string
		Options        []string
	}

	GetEdgeHostnamesResponse struct {
		AccountID     string            `json:"accountId"`
		ContractID    string            `json:"contractId"`
		GroupID       string            `json:"groupId"`
		EdgeHostnames EdgeHostnameItems `json:"edgeHostnames"`
	}

	EdgeHostnameItems struct {
		Items []EdgeHostnameGetItem `json:"items"`
	}

	EdgeHostnameGetItem struct {
		ID                string    `json:"edgeHostnameId"`
		Domain            string    `json:"edgeHostnameDomain"`
		ProductID         string    `json:"productId"`
		DomainPrefix      string    `json:"domainPrefix"`
		DomainSuffix      string    `json:"domainSuffix"`
		Status            string    `json:"status,omitempty"`
		Secure            bool      `json:"secure"`
		IPVersionBehavior string    `json:"ipVersionBehavior"`
		UseCases          []UseCase `json:"useCases,omitempty"`
	}

	UseCase struct {
		Option  string `json:"option"`
		Type    string `json:"type"`
		UseCase string `json:"useCase"`
	}

	CreateEdgeHostnameRequest struct {
		ContractID   string
		GroupID      string
		Options      []string
		EdgeHostname EdgeHostnameCreate
	}

	EdgeHostnameCreate struct {
		ProductID         string    `json:"productId"`
		DomainPrefix      string    `json:"domainPrefix"`
		DomainSuffix      string    `json:"domainSuffix"`
		Secure            bool      `json:"secure,omitempty"`
		SecureNetwork     string    `json:"secureNetwork,omitempty"`
		SlotNumber        int       `json:"slotNumber,omitEmpty"`
		IPVersionBehavior string    `json:"ipVersionBehavior"`
		CertEnrollmentID  int       `json:"certEnrollmentId,omitempty"`
		UseCases          []UseCase `json:"useCases,omitempty"`
	}

	CreateEdgeHostnameResponse struct {
		EdgeHostnameLink string `json:"edgeHostnameLink"`
		EdgeHostnameID   string `json:"-"`
	}
)

const (
	EHSecureNetworkStandardTLS = "STANDARD_TLS"
	EHSecureNetworkSharedCert  = "SHARED_CERT"
	EHSecureNetworkEnhancedTLS = "ENHANCED_TLS"
	EHIPVersionV4              = "IPV4"
	EHIPVersionV6Compliance    = "IPV4"

	UseCaseGlobal = "GLOBAL"
)

func (eh CreateEdgeHostnameRequest) Validate() error {
	return validation.ValidateStruct(&eh,
		validation.Field(&eh.ContractID, validation.Required),
		validation.Field(&eh.GroupID, validation.Required),
		validation.Field(&eh.EdgeHostname.DomainPrefix, validation.Required),
		validation.Field(&eh.EdgeHostname.DomainSuffix, validation.Required,
			validation.When(eh.EdgeHostname.SecureNetwork == EHSecureNetworkStandardTLS, validation.In("edgesuite.net")),
			validation.When(eh.EdgeHostname.SecureNetwork == EHSecureNetworkSharedCert, validation.In("akamaized.net")),
			validation.When(eh.EdgeHostname.SecureNetwork == EHSecureNetworkEnhancedTLS, validation.In("edgekey.net")),
		),
		validation.Field(&eh.EdgeHostname.ProductID, validation.Required),
		validation.Field(&eh.EdgeHostname.CertEnrollmentID, validation.Required.When(eh.EdgeHostname.SecureNetwork == EHSecureNetworkEnhancedTLS)),
		validation.Field(&eh.EdgeHostname.IPVersionBehavior, validation.Required, validation.In(EHIPVersionV4, EHIPVersionV6Compliance)),
		validation.Field(&eh.EdgeHostname.SecureNetwork, validation.In(EHSecureNetworkStandardTLS, EHSecureNetworkSharedCert, EHSecureNetworkEnhancedTLS)),
		validation.Field(&eh.EdgeHostname.UseCases),
	)
}

func (uc UseCase) Validate() error {
	return validation.ValidateStruct(&uc,
		validation.Field(&uc.Option, validation.Required),
		validation.Field(&uc.Type, validation.Required, validation.In(UseCaseGlobal)),
		validation.Field(&uc.UseCase, validation.Required),
	)
}

func (eh GetEdgeHostnamesRequest) Validate() error {
	return validation.ValidateStruct(&eh,
		validation.Field(&eh.ContractID, validation.Required),
		validation.Field(&eh.GroupID, validation.Required),
	)
}

func (eh GetEdgeHostnameRequest) Validate() error {
	return validation.ValidateStruct(&eh,
		validation.Field(&eh.EdgeHostnameID, validation.Required),
		validation.Field(&eh.ContractID, validation.Required),
		validation.Field(&eh.GroupID, validation.Required),
	)
}

func (p *papi) GetEdgeHostnames(ctx context.Context, params GetEdgeHostnamesRequest) (*GetEdgeHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEdgeHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/edgehostnames?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID,
	)
	if len(params.Options) > 0 {
		getURL = fmt.Sprintf("%s&options=%s", getURL, strings.Join(params.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getedgehostnames request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var edgeHostnames GetEdgeHostnamesResponse
	resp, err := p.Exec(req, &edgeHostnames)
	if err != nil {
		return nil, fmt.Errorf("getedgehostnames request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &edgeHostnames, nil
}

func (p *papi) GetEdgeHostname(ctx context.Context, params GetEdgeHostnameRequest) (*GetEdgeHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEdgeHostname")

	getURL := fmt.Sprintf(
		"/papi/v1/edgehostnames/%s?contractId=%s&groupId=%s",
		params.EdgeHostnameID,
		params.ContractID,
		params.GroupID,
	)
	if len(params.Options) > 0 {
		getURL = fmt.Sprintf("%s&options=%s", getURL, strings.Join(params.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getedgehostname request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var edgeHostname GetEdgeHostnamesResponse
	resp, err := p.Exec(req, &edgeHostname)
	if err != nil {
		return nil, fmt.Errorf("getedgehostname request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &edgeHostname, nil
}

func (p *papi) CreateEdgeHostname(ctx context.Context, r CreateEdgeHostnameRequest) (*CreateEdgeHostnameResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateEdgeHostname")

	createURL := fmt.Sprintf(
		"/papi/v1/edgehostnames?contractId=%s&groupId=%s",
		r.ContractID,
		r.GroupID,
	)
	if len(r.Options) > 0 {
		createURL = fmt.Sprintf("%s&options=%s", createURL, strings.Join(r.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create createedgehostname request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var createResponse CreateEdgeHostnameResponse
	resp, err := p.Exec(req, &createResponse, r.EdgeHostname)
	if err != nil {
		return nil, fmt.Errorf("createedgehostname request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, session.NewAPIError(resp, logger)
	}
	id, err := tools.FetchIDFromLocation(createResponse.EdgeHostnameLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidLocation, err.Error())
	}
	createResponse.EdgeHostnameID = id
	return &createResponse, nil
}
