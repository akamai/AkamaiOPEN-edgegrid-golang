package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PropertyVersionHostnames contains operations available on PropertyVersionHostnames resource
	PropertyVersionHostnames interface {
		// GetPropertyVersionHostnames lists all the hostnames assigned to a property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version-hostnames
		GetPropertyVersionHostnames(context.Context, GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error)

		// UpdatePropertyVersionHostnames modifies the set of hostnames for a property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/patch-property-version-hostnames
		UpdatePropertyVersionHostnames(context.Context, UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error)
	}

	// GetPropertyVersionHostnamesRequest contains parameters required to list property version hostnames
	GetPropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		IncludeCertStatus bool
	}

	// GetPropertyVersionHostnamesResponse contains all property version hostnames associated to the given parameters
	GetPropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// HostnameResponseItems contains the response body for GetPropertyVersionHostnamesResponse
	HostnameResponseItems struct {
		Items []Hostname `json:"items"`
	}

	// Hostname contains information about each of the HostnameResponseItems
	Hostname struct {
		CnameType            HostnameCnameType `json:"cnameType"`
		EdgeHostnameID       string            `json:"edgeHostnameId,omitempty"`
		CnameFrom            string            `json:"cnameFrom"`
		CnameTo              string            `json:"cnameTo,omitempty"`
		CertProvisioningType string            `json:"certProvisioningType"`
		CertStatus           CertStatusItem    `json:"certStatus,omitempty"`
	}

	// CertStatusItem contains information about certificate status for specific Hostname
	CertStatusItem struct {
		ValidationCname ValidationCname `json:"validationCname,omitempty"`
		Staging         []StatusItem    `json:"staging,omitempty"`
		Production      []StatusItem    `json:"production,omitempty"`
	}

	// ValidationCname is the CNAME record used to validate the certificate’s domain
	ValidationCname struct {
		Hostname string `json:"hostname,omitempty"`
		Target   string `json:"target,omitempty"`
	}

	// StatusItem determines whether a hostname is capable of serving secure content over the staging or production network.
	StatusItem struct {
		Status string `json:"status,omitempty"`
	}

	// UpdatePropertyVersionHostnamesRequest contains parameters required to update the set of hostname entries for a property version
	UpdatePropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		IncludeCertStatus bool
		Hostnames         []Hostname
	}

	// UpdatePropertyVersionHostnamesResponse contains information about each of the HostnameRequestItems
	UpdatePropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// HostnameCnameType represents HostnameCnameType enum
	HostnameCnameType string
)

const (
	// HostnameCnameTypeEdgeHostname const
	HostnameCnameTypeEdgeHostname HostnameCnameType = "EDGE_HOSTNAME"
)

// Validate validates GetPropertyVersionHostnamesRequest
func (ph GetPropertyVersionHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ph.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ph.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates UpdatePropertyVersionHostnamesRequest
func (ch UpdatePropertyVersionHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ch.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ch.PropertyVersion, validation.Required),
	}.Filter()
}

var (
	// ErrGetPropertyVersionHostnames represents error when fetching hostnames fails
	ErrGetPropertyVersionHostnames = errors.New("fetching hostnames")
	// ErrUpdatePropertyVersionHostnames represents error when updating hostnames fails
	ErrUpdatePropertyVersionHostnames = errors.New("updating hostnames")
)

func (p *papi) GetPropertyVersionHostnames(ctx context.Context, params GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersionHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyVersionHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/hostnames?contractId=%s&groupId=%s&validateHostnames=%t&includeCertStatus=%t",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
		params.ValidateHostnames,
		params.IncludeCertStatus)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyVersionHostnames, err)
	}

	var hostnames GetPropertyVersionHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyVersionHostnames, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) UpdatePropertyVersionHostnames(ctx context.Context, params UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePropertyVersionHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("UpdatePropertyVersionHostnames")

	putURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%v/hostnames?contractId=%s&groupId=%s&validateHostnames=%t&includeCertStatus=%t",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
		params.ValidateHostnames,
		params.IncludeCertStatus,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePropertyVersionHostnames, err)
	}

	var hostnames UpdatePropertyVersionHostnamesResponse
	newHostnames := params.Hostnames
	if newHostnames == nil {
		newHostnames = []Hostname{}
	}
	resp, err := p.Exec(req, &hostnames, newHostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePropertyVersionHostnames, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}
