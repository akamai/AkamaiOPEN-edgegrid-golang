package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// SortOrder represents SortOrder enum.
	SortOrder string

	// CertType represents CertType enum.
	CertType string

	// ListActivePropertyHostnamesRequest contains parameters required to list active property hostnames.
	ListActivePropertyHostnamesRequest struct {
		PropertyID        string
		Offset            int
		Limit             int
		Sort              SortOrder
		Hostname          string
		CnameTo           string
		Network           ActivationNetwork
		ContractID        string
		GroupID           string
		IncludeCertStatus bool
	}

	// GetActivePropertyHostnamesDiffRequest contains parameters required to list active property hostnames diff.
	GetActivePropertyHostnamesDiffRequest struct {
		// PropertyID is the unique identifier for the property.
		PropertyID string

		// Offset specifies the page of results you want to navigate to, starting from 0.
		Offset int

		// Limit specifies the number of hostnames objects to include on each page.
		Limit int

		// ContractID is the unique identifier for the contract.
		ContractID string

		// GroupID is the unique identifier for the group.
		GroupID string
	}

	// ListActiveAccountHostnamesRequest contains parameters required to list active property hostnames for an account.
	ListActiveAccountHostnamesRequest struct {
		// Offset specifies the page of results you want to navigate to, starting from 0.
		Offset int

		// Limit specifies the number of hostnames objects to include on each page.
		Limit int

		// Sort sorts the results based on the `cnameFrom` value, either `hostname:a` for ascending or `hostname:d` for descending order.
		// The default is `hostname:a`.
		Sort SortOrder

		// Hostname filters the results by `cnameFrom`. Supports wildcard matches with `*`.
		Hostname string

		// CnameTo filters the results by edge hostname. Supports wildcard matches with `*`.
		CnameTo string

		// Network is the network of activated hostnames, either `STAGING` or `PRODUCTION`.
		Network ActivationNetwork

		// ContractID is the unique identifier for the contract.
		ContractID string

		// GroupID is the unique identifier for the group.
		GroupID string
	}

	// ListActivePropertyHostnamesResponse contains information about each of the active property hostnames.
	ListActivePropertyHostnamesResponse struct {
		AccountID     string                 `json:"accountId"`
		AvailableSort []SortOrder            `json:"availableSort"`
		ContractID    string                 `json:"contractId"`
		CurrentSort   SortOrder              `json:"currentSort"`
		DefaultSort   SortOrder              `json:"defaultSort"`
		GroupID       string                 `json:"groupId"`
		PropertyID    string                 `json:"propertyId"`
		PropertyName  string                 `json:"propertyName"`
		Hostnames     HostnamesResponseItems `json:"hostnames"`
	}

	// GetActivePropertyHostnamesDiffResponse contains information about each of the active property hostnames that differ in staging and production networks.
	GetActivePropertyHostnamesDiffResponse struct {
		// AccountID identifies the prevailing account under which you requested the data.
		AccountID string `json:"accountId"`

		// ContractID identifies the prevailing contract under which you requested the data.
		ContractID string `json:"contractId"`

		// GroupID identifies the prevailing group under which you requested the data.
		GroupID string `json:"groupId"`

		// PropertyID is the unique identifier for a property.
		PropertyID string `json:"propertyId"`

		// Hostnames is the active property hostnames that differ in staging and production networks.
		Hostnames HostnamesDiffResponseItems `json:"hostnames"`
	}

	// ListActiveAccountHostnamesResponse contains information about each of the active property hostnames.
	ListActiveAccountHostnamesResponse struct {
		// AccountID identifies the prevailing account under which you requested the data.
		AccountID string `json:"accountId"`

		// AvailableSort contains available sorting options: `hostname:a` for ascending, and `hostname:d` for descending.
		AvailableSort []string `json:"availableSort"`

		// CurrentSort is the sorting criteria applied to the response, either `hostname:a` for ascending, or `hostname:d` for descending.
		CurrentSort string `json:"currentSort"`

		// DefaultSort shows the default `hostname:a` sorting criteria if you didn't specify any query parameters in the request.
		DefaultSort string `json:"defaultSort"`

		// Hostnames is the set of requested hostnames, available within an items array.
		Hostnames ActiveAccountHostnames `json:"hostnames"`
	}

	// HostnamesResponseItems contains the response body for ListActivePropertyHostnamesResponse.
	HostnamesResponseItems struct {
		Items            []HostnameItem `json:"items"`
		CurrentItemCount int            `json:"currentItemCount"`
		NextLink         *string        `json:"nextLink"`
		PreviousLink     *string        `json:"previousLink"`
		TotalItems       int            `json:"totalItems"`
	}

	// HostnamesDiffResponseItems contains the response body for GetActivePropertyHostnamesDiffResponse.
	HostnamesDiffResponseItems struct {
		// Items are the details of the active property hostnames that differ in staging and production networks.
		Items []HostnameDiffItem `json:"items"`

		// CurrentItemCount is the total count of items present in the current response body for requested criteria.
		CurrentItemCount int `json:"currentItemCount"`

		// NextLink is the link to next set of response items for paginated response.
		NextLink *string `json:"nextLink"`

		// PreviousLink is the link to previous set of response items for paginated response.
		PreviousLink *string `json:"previousLink"`

		// TotalItems is the total count of items for requested criteria.
		TotalItems int `json:"totalItems"`
	}

	// ActiveAccountHostnames contains the response body for ListActiveAccountHostnamesResponse.
	ActiveAccountHostnames struct {
		// Items contains each hostname.
		Items []ActiveAccountHostnameItem `json:"items"`

		// CurrentItemCount is the number of items present in the current response body view.
		CurrentItemCount int `json:"currentItemCount"`

		// NextLink is the link to the next set of response items for paginated responses.
		NextLink *string `json:"nextLink"`

		// PreviousLink is the link to the previous set of response items for paginated responses.
		PreviousLink *string `json:"previousLink"`

		// TotalItems is the total count of items returned for the requested criteria.
		TotalItems int `json:"totalItems"`
	}

	// HostnameItem contains information about each of the HostnamesResponseItems.
	HostnameItem struct {
		// CCMCertStatus is deployment status for the RSA and ECDSA certificates created with Cloud Certificate Manager (CCM).
		CCMCertStatus *CCMCertStatus `json:"ccmCertStatus,omitempty"`

		// CCMCertificates is certificate identifiers and links for the CCM-managed certificates.
		CCMCertificates *CCMCertificatesResp `json:"ccmCertificates,omitempty"`

		// CertStatus with the `includeCertStatus` parameter set to `true`,
		// determines whether a hostname is capable of serving secure content over the staging or production network.
		CertStatus *CertStatusItem `json:"certStatus"`

		// CnameFrom is hostname that your end users see, indicated by the `Host` header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// CnameType has one supported `EDGE_HOSTNAME` value.
		CnameType HostnameCnameType `json:"cnameType"`

		// MTLS is mutual TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		MTLS *MTLSResp `json:"mtls,omitempty"`

		// ProductionCertProvisioningType indicates the certificate's provisioning type.
		// Either `CPS_MANAGED` type for the certificates you create with the Certificate Provisioning System API (CPS),
		// `DEFAULT` for the Default Domain Validation (DV) certificates created automatically,
		// or `CCM` type for the third party certificates you create with the Cloud Certificate Manager.
		// Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.
		ProductionCertType CertType `json:"productionCertType"`

		// ProductionCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		// This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		ProductionCnameTo string `json:"productionCnameTo"`

		// ProductionEdgeHostnameID identifies each edge hostname.
		ProductionEdgeHostnameID string `json:"productionEdgeHostnameId"`

		// StagingCertType indicates the certificate's provisioning type.
		// Either `CPS_MANAGED` type for the certificates you create with the Certificate Provisioning System API (CPS),
		// `DEFAULT` for the Default Domain Validation (DV) certificates created automatically,
		// or `CCM` type for the third party certificates you create with the Cloud Certificate Manager.
		// Note that you can't specify the `DEFAULT` value if your property hostname uses the `akamaized.net` domain suffix.
		StagingCertType CertType `json:"stagingCertType"`

		// StagingCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		// This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		StagingCnameTo string `json:"StagingCnameTo"`

		// StagingEdgeHostnameID identifies each edge hostname.
		StagingEdgeHostnameID string `json:"stagingEdgeHostnameId"`

		// TLSConfiguration is optional TLS configuration settings applicable to the Cloud Certificate Manager (CCM) hostnames.
		TLSConfiguration *TLSConfiguration `json:"tlsConfiguration,omitempty"`
	}

	// HostnameDiffItem contains information about each of the HostnamesDiffResponseItems.
	HostnameDiffItem struct {
		// CnameFrom is the hostname that your end users see, indicated by the Host header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// ProductionCertProvisioningType indicates the type of the certificate used in the property hostname.
		ProductionCertProvisioningType CertType `json:"productionCertProvisioningType"`

		// ProductionCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		ProductionCnameTo string `json:"productionCnameTo"`

		// ProductionCnameType indicates the type of CNAME you used in the production network, either `EDGE_HOSTNAME` or `CUSTOM`.
		ProductionCnameType HostnameCnameType `json:"productionCnameType"`

		// ProductionEdgeHostnameID identifies the edge hostname you mapped your traffic to on the production network.
		ProductionEdgeHostnameID string `json:"productionEdgeHostnameId"`

		// StagingCertProvisioningType indicates the type of the certificate used in the property hostname.
		StagingCertProvisioningType CertType `json:"stagingCertProvisioningType"`

		// StagingCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		StagingCnameTo string `json:"stagingCnameTo"`

		// StagingCnameType indicates the type of CNAME you used in the staging network, either `EDGE_HOSTNAME` or `CUSTOM`.
		StagingCnameType HostnameCnameType `json:"stagingCnameType"`

		// StagingEdgeHostnameID identifies the edge hostname you mapped your traffic to on the production network.
		StagingEdgeHostnameID string `json:"stagingEdgeHostnameId"`
	}

	// ActiveAccountHostnameItem contains information about each of the account hostname.
	ActiveAccountHostnameItem struct {
		// CnameFrom is the hostname that your end users see, indicated by the Host header in end user requests.
		CnameFrom string `json:"cnameFrom"`

		// ContractID identifies the prevailing contract under which you requested the data.
		ContractID string `json:"contractId"`

		// GroupID identifies the prevailing group under which you requested the data.
		GroupID string `json:"groupId"`

		// LatestVersion specifies the most recent version of the property.
		LatestVersion int `json:"latestVersion"`

		// PropertyID is the unique identifier for the property.
		PropertyID string `json:"propertyId"`

		// PropertyName is a unique, descriptive name for the property. It's not modifiable after you initially create it on a POST request.
		PropertyName string `json:"propertyName"`

		// PropertyType specifies the type of the property.
		// Either `TRADITIONAL` for properties where you pair property hostnames with the property version,
		// or `HOSTNAME_BUCKET` where you manage property hostnames independently of the property version.
		PropertyType string `json:"propertyType"`

		// ProductionCertType indicates the type of the certificate used in the property hostname.
		// Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System API (CPS),
		// or `DEFAULT` for Default Domain Validation (DV) certificates deployed automatically.
		// Note that you can't specify the `DEFAULT` value if your account hostname uses the `akamaized.net` domain suffix.
		ProductionCertType *string `json:"productionCertType"`

		// ProductionCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		// This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		ProductionCnameTo *string `json:"productionCnameTo"`

		// ProductionCnameType indicates the type of CNAME you used in the production network, either `EDGE_HOSTNAME` or `CUSTOM`.
		ProductionCnameType *string `json:"productionCnameType"`

		// ProductionEdgeHostnameID identifies the edge hostname you mapped your traffic to on the production network.
		ProductionEdgeHostnameID *string `json:"productionEdgeHostnameId"`

		// ProductionProductID identifies the product association on the network.
		ProductionProductID *string `json:"productionProductId"`

		// StagingCertType indicates the type of the certificate used in the property hostname.
		// Either `CPS_MANAGED` for the certificates you create with the Certificate Provisioning System API (CPS),
		// or `DEFAULT` for Default Domain Validation (DV) certificates deployed automatically.
		// Note that you can't specify the `DEFAULT` value if your account hostname uses the `akamaized.net` domain suffix.
		StagingCertType *string `json:"stagingCertType"`

		// StagingCnameTo is the edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers.
		// This member corresponds to the edge hostname object's `edgeHostnameDomain` member.
		StagingCnameTo *string `json:"stagingCnameTo"`

		// StagingCnameType indicates the type of CNAME you used in the staging network, either `EDGE_HOSTNAME` or `CUSTOM`.
		StagingCnameType *string `json:"stagingCnameType"`

		// StagingEdgeHostnameID identifies the edge hostname you mapped your traffic to on the staging network.
		StagingEdgeHostnameID *string `json:"stagingEdgeHostnameId"`

		// StagingProductID identifies the product association on the network.
		StagingProductID *string `json:"stagingProductId"`
	}
)

const (
	// SortAscending represents ascending sorting by hostname.
	SortAscending SortOrder = "hostname:a"
	// SortDescending represents descending sorting by hostname.
	SortDescending SortOrder = "hostname:d"
	// CertTypeCPSManaged indicates that the certificate is provisioned using the Certificate Provisioning System (CPS).
	CertTypeCPSManaged CertType = "CPS_MANAGED"
	// CertTypeDefault indicates that the certificate is a Default Domain Validation (DV) certificate.
	CertTypeDefault CertType = "DEFAULT"
	// CertTypeCCM indicates that the certificate is a Cloud Controller Manager (CCM) certificate.
	CertTypeCCM CertType = "CCM"
	// maxHostnamesPerPage indicates the maximum possible value for 'limit' parameter for Get and List active property hostnames.
	maxHostnamesPerPage int = 999
)

var (
	// ErrListActivePropertyHostnames represents error when fetching active property hostnames fails.
	ErrListActivePropertyHostnames = errors.New("fetching active property hostnames")

	// ErrGetActivePropertyHostnamesDiff represents error when fetching active property hostnames diff fails.
	ErrGetActivePropertyHostnamesDiff = errors.New("fetching active property hostnames diff")

	// ErrListActiveAccountHostnames represents error when fetching active account hostnames fails.
	ErrListActiveAccountHostnames = errors.New("fetching active account hostnames")
)

// Validate validates SortOrder.
func (o SortOrder) Validate() validation.InRule {
	return validation.In(SortAscending, SortDescending).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'",
			o, SortAscending, SortDescending))
}

// Validate validates CertType.
func (t CertType) Validate() validation.InRule {
	return validation.In(CertTypeCPSManaged, CertTypeDefault, CertTypeCCM).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'",
			t, CertTypeCPSManaged, CertTypeDefault, CertTypeCCM))
}

// Validate validates ListActivePropertyHostnamesRequest.
func (r ListActivePropertyHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Network":    validation.Validate(r.Network, r.Network.Validate()),
		"Sort":       validation.Validate(r.Sort, r.Sort.Validate()),
		"Offset":     validation.Validate(r.Offset, validation.Min(0)),
		"Limit":      validation.Validate(r.Limit, validation.Min(1), validation.Max(maxHostnamesPerPage)),
	}.Filter()
}

// Validate validates GetActivePropertyHostnamesDiffRequest.
func (r GetActivePropertyHostnamesDiffRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(r.PropertyID, validation.Required),
		"Offset":     validation.Validate(r.Offset, validation.Min(0)),
		"Limit":      validation.Validate(r.Limit, validation.Min(1), validation.Max(maxHostnamesPerPage)),
	}.Filter()
}

// Validate validates ListActiveAccountHostnamesRequest.
func (r ListActiveAccountHostnamesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network, r.Network.Validate()),
		"Sort":    validation.Validate(r.Sort, r.Sort.Validate()),
		"Offset":  validation.Validate(r.Offset, validation.Min(0)),
		"Limit":   validation.Validate(r.Limit, validation.Min(1), validation.Max(maxHostnamesPerPage)),
	})
}

func (p *papi) ListActivePropertyHostnames(ctx context.Context, params ListActivePropertyHostnamesRequest) (*ListActivePropertyHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListActivePropertyHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListActivePropertyHostnames, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/properties/%s/hostnames", params.PropertyID).
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("sort", string(params.Sort), params.Sort != "").
		AddQueryParamIf("hostname", params.Hostname, params.Hostname != "").
		AddQueryParamIf("cnameTo", params.CnameTo, params.CnameTo != "").
		AddQueryParamIf("network", string(params.Network), params.Network != "").
		AddQueryParamIf("includeCertStatus", strconv.FormatBool(params.IncludeCertStatus), params.IncludeCertStatus).
		AddQueryParamIf("limit", strconv.Itoa(params.Limit), params.Limit != 0).
		AddQueryParamIf("offset", strconv.Itoa(params.Offset), params.Offset != 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListActivePropertyHostnames, err)
	}

	var hostnames ListActivePropertyHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListActivePropertyHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListActivePropertyHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) GetActivePropertyHostnamesDiff(ctx context.Context, params GetActivePropertyHostnamesDiffRequest) (*GetActivePropertyHostnamesDiffResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetActivePropertyHostnamesDiff")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetActivePropertyHostnamesDiff, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/properties/%s/hostnames/diff", params.PropertyID).
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("limit", strconv.Itoa(params.Limit), params.Limit != 0).
		AddQueryParamIf("offset", strconv.Itoa(params.Offset), params.Offset != 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetActivePropertyHostnamesDiff, err)
	}

	var hostnamesDiff GetActivePropertyHostnamesDiffResponse
	resp, err := p.Exec(req, &hostnamesDiff)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetActivePropertyHostnamesDiff, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetActivePropertyHostnamesDiff, p.Error(resp))
	}

	return &hostnamesDiff, nil
}

func (p *papi) ListActiveAccountHostnames(ctx context.Context, params ListActiveAccountHostnamesRequest) (*ListActiveAccountHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListActiveAccountHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListActiveAccountHostnames, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/papi/v1/hostnames").
		AddQueryParamIf("offset", strconv.Itoa(params.Offset), params.Offset != 0).
		AddQueryParamIf("limit", strconv.Itoa(params.Limit), params.Limit != 0).
		AddQueryParamIf("sort", string(params.Sort), params.Sort != "").
		AddQueryParamIf("hostname", params.Hostname, params.Hostname != "").
		AddQueryParamIf("cnameTo", params.CnameTo, params.CnameTo != "").
		AddQueryParamIf("network", string(params.Network), params.Network != "").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListActiveAccountHostnames, err)
	}

	var hostnames ListActiveAccountHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListActiveAccountHostnames, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListActiveAccountHostnames, p.Error(resp))
	}

	return &hostnames, nil
}
