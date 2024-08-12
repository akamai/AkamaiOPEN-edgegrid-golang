package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PropertyHostnames contains operations available on PropertyHostnames resource
	PropertyHostnames interface {
		// ListPropertyHostnames lists all the active hostnames for a property
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-hostnames
		ListPropertyHostnames(context.Context, ListPropertyHostnamesRequest) (*ListPropertyHostnamesResponse, error)

		// ListPropertyHostnames lists all hostname activations for a property
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-hostname-activations
		ListPropertyHostnameActivations(context.Context, ListPropertyHostnameActivationsRequest) (*ListPropertyHostnameActivationsResponse, error)

		// GetPropertyHostnames gets details about a specific property hostname activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-hostname-activation
		GetPropertyHostnameActivation(context.Context, GetPropertyHostnameActivationRequest) (*GetPropertyHostnameActivationResponse, error)

		// PatchPropertyHostnames modifies the set of hostnames for a property
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/patch-property-hostnames
		PatchPropertyHostnames(context.Context, PatchPropertyHostnamesRequest) (*PatchPropertyHostnamesResponse, error)
	}

	// ListPropertyHostnamesRequest contains parameters required to list property hostnames
	ListPropertyHostnamesRequest struct {
		PropertyID string
		ContractID string
		GroupID    string

		Offset   int
		Limit    int
		Sort     string
		Hostname string
		CnameTo  string
		Network  ActivationNetwork
	}

	// ListPropertyHostnamesResponse contains all property version hostnames associated to the given parameters
	ListPropertyHostnamesResponse struct {
		AccountID    string        `json:"accountId"`
		ContractID   string        `json:"contractId"`
		GroupID      string        `json:"groupId"`
		Hostnames    HostnamesList `json:"hostnames"`
		PropertyID   string        `json:"propertyId"`
		PropertyName string        `json:"propertyName"`
	}

	HostnamesList struct {
		Items        []BucketHostname `json:"items"`
		NextLink     string           `json:"nextLink"`
		PreviousLink string           `json:"previousLink"`
		TotalItems   int              `json:"totalItems"`
	}

	BucketHostname struct {
		CertStatus CertStatusItem `json:"certStatus,omitempty"`
		CnameFrom  string         `json:"cnameFrom"`

		ProductionCertType       string `json:"productionCertType"`
		ProductionCnameTo        string `json:"productionCnameTo"`
		ProductionEdgeHostnameID string `json:"productionEdgeHostnameId"`

		StagingCertType       string `json:"stagingCertType"`
		StagingCnameTo        string `json:"stagingCnameTo"`
		StagingEdgeHostnameID string `json:"stagingEdgeHostnameId"`
	}

	// ListPropertyHostnameActivationssRequest contains parameters required to list property hostname activations
	ListPropertyHostnameActivationsRequest struct {
		PropertyID string
		ContractID string
		GroupID    string

		Offset int
		Limit  int
	}

	// ListPropertyHostnameActivationsResponse
	ListPropertyHostnameActivationsResponse struct {
		AccountID  string `json:"accountId"`
		ContractID string `json:"contractId"`
		GroupID    string `json:"groupId"`

		CurrentItemCount int    `json:"currentItemCount"`
		NextLink         string `json:"nextLink"`
		PreviousLink     string `json:"previousLink"`
		TotalItems       int    `json:"totalItems"`

		HostnameActivations HostnameActivationsList `json:"hostnameActivations"`
	}

	HostnameActivationsList struct {
		NextLink     string               `json:"nextLink"`
		PreviousLink string               `json:"previousLink"`
		TotalItems   int                  `json:"totalItems"`
		Items        []HostnameActivation `json:"items"`
	}

	HostnameActivation struct {
		AccountID            string            `json:"accountId,omitempty"`
		ActivationType       ActivationType    `json:"activationType,omitempty"`
		GroupID              string            `json:"groupId,omitempty"`
		HostnameActivationID string            `json:"hostnameActivationId,omitempty"`
		Network              ActivationNetwork `json:"network"`
		Note                 string            `json:"note,omitempty"`
		NotifyEmails         []string          `json:"notifyEmails"`
		PropertyID           string            `json:"propertyId,omitempty"`
		PropertyName         string            `json:"propertyName,omitempty"`
		PropertyVersion      int               `json:"propertyVersion"`
		Status               ActivationStatus  `json:"status,omitempty"`
		SubmitDate           string            `json:"submitDate,omitempty"`
		UpdateDate           string            `json:"updateDate,omitempty"`
	}

	// GetPropertyHostnameActivationRequest
	GetPropertyHostnameActivationRequest struct {
		PropertyID           string
		HostnameActivationID string
		ContractID           string
		GroupID              string

		IncludeHostnames bool
	}

	// GetPropertyHostnameActivationResponse
	GetPropertyHostnameActivationResponse struct {
		AccountID  string `json:"accountId"`
		ContractID string `json:"contractId"`
		GroupID    string `json:"groupId"`

		HostnameActivations HostnameActivationsList `json:"hostnameActivations"`
	}

	// PatchPropertyHostnamesRequest contains parameters required to update the set of hostname entries for a property version
	PatchPropertyHostnamesRequest struct {
		PropertyID   string
		ContractID   string
		GroupID      string
		Network      ActivationNetwork `json:"network"`
		Note         string            `json:"note"`
		NotifyEmails []string          `json:"notifyEmails"`
		Add          []Hostname        `json:"add"`
		Remove       []string          `json:"remove"`
	}

	// PatchPropertyHostnamesResponse
	PatchPropertyHostnamesResponse struct {
		ActivationID   string     `json:"activationId"`
		ActivationLink string     `json:"activationLink"`
		Hostnames      []Hostname `json:"hostnames"`
	}
)

// Validate validates ListPropertyHostnamesRequest
func (ph ListPropertyHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(ph.PropertyID, validation.Required),
	}.Filter()
}

// Validate validates ListPropertyHostnameActivationsRequest
func (ph ListPropertyHostnameActivationsRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(ph.PropertyID, validation.Required),
	}.Filter()
}

// Validate validates GetPropertyHostnameActivationRequest
func (ph GetPropertyHostnameActivationRequest) Validate() error {
	return validation.Errors{
		"PropertyID":           validation.Validate(ph.PropertyID, validation.Required),
		"HostnameActivationID": validation.Validate(ph.HostnameActivationID, validation.Required),
	}.Filter()
}

// Validate validates PatchPropertyHostnamesRequest
func (ch PatchPropertyHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(ch.PropertyID, validation.Required),
	}.Filter()
}

var (
	// ErrListPropertyHostnames represents error when listing hostnames fails
	ErrListPropertyHostnames = errors.New("listing hostnames")
	// ErrListPropertyHostnameActivations represents error when listing hostname activations fails
	ErrListPropertyHostnameActivations = errors.New("listing hostname activations")
	// ErrGetPropertyHostnameActivation represents error when getting hostname activation fails
	ErrGetPropertyHostnameActivation = errors.New("getting hostname activation")
	// ErrPatchPropertyHostnames represents error when patching hostnames fails
	ErrPatchPropertyHostnames = errors.New("patching hostnames")
)

func (p *papi) ListPropertyHostnames(ctx context.Context, params ListPropertyHostnamesRequest) (*ListPropertyHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPropertyHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("ListPropertyHostnames")

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/hostnames", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse uri: %s", ErrListPropertyHostnames, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.Offset != 0 {
		q.Add("offset", fmt.Sprint(params.Offset))
	}
	if params.Limit != 0 {
		q.Add("limit", fmt.Sprint(params.Limit))
	}
	if params.Sort != "" {
		q.Add("sort", params.Sort)
	}
	if params.Hostname != "" {
		q.Add("hostname", params.Hostname)
	}
	if params.CnameTo != "" {
		q.Add("cnameTo", params.CnameTo)
	}
	if params.Network != "" {
		q.Add("network", string(params.Network))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPropertyHostnames, err)
	}

	var hostnames ListPropertyHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPropertyHostnames, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPropertyHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) ListPropertyHostnameActivations(ctx context.Context, params ListPropertyHostnameActivationsRequest) (*ListPropertyHostnameActivationsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPropertyHostnameActivations, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("ListPropertyHostnameActivations")

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/hostname-activations", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse uri: %s", ErrListPropertyHostnameActivations, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.Offset != 0 {
		q.Add("offset", fmt.Sprint(params.Offset))
	}
	if params.Limit != 0 {
		q.Add("limit", fmt.Sprint(params.Limit))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPropertyHostnameActivations, err)
	}

	var hostnames ListPropertyHostnameActivationsResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPropertyHostnameActivations, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPropertyHostnameActivations, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) GetPropertyHostnameActivation(ctx context.Context, params GetPropertyHostnameActivationRequest) (*GetPropertyHostnameActivationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyHostnameActivation, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyHostnameActivation")

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/hostname-activations", params.PropertyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse uri: %s", ErrGetPropertyHostnameActivation, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyHostnameActivation, err)
	}

	var hostnames GetPropertyHostnameActivationResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyHostnameActivation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyHostnameActivation, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) PatchPropertyHostnames(ctx context.Context, params PatchPropertyHostnamesRequest) (*PatchPropertyHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrPatchPropertyHostnames, ErrStructValidation, err)
	}

	patchURL := fmt.Sprintf(
		"/papi/v1/properties/%s/hostnames?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ContractID,
		params.GroupID,
	)

	logger := p.Log(ctx)
	logger.Debug("PatchPropertyHostnames")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, patchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrPatchPropertyHostnames, err)
	}

	var hostnames PatchPropertyHostnamesResponse
	resp, err := p.Exec(req, &hostnames, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrPatchPropertyHostnames, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrPatchPropertyHostnames, p.Error(resp))
	}

	return &hostnames, nil
}
