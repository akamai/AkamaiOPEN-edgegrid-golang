package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListEdgeKVNamespacesRequest contains path parameters used to list namespaces
	ListEdgeKVNamespacesRequest struct {
		Network NamespaceNetwork
		Details bool
	}

	// GetEdgeKVNamespaceRequest contains path parameters used to fetch a namespace
	GetEdgeKVNamespaceRequest struct {
		Network NamespaceNetwork
		Name    string
	}

	// CreateEdgeKVNamespaceRequest contains path parameter and request body used to create a namespace
	CreateEdgeKVNamespaceRequest struct {
		Network NamespaceNetwork
		Namespace
	}

	// UpdateEdgeKVNamespaceRequest contains path parameters and request body used to update a namespace
	UpdateEdgeKVNamespaceRequest struct {
		Network NamespaceNetwork
		UpdateNamespace
	}

	// ListEdgeKVNamespacesResponse represents a response object returned when listing namespaces
	ListEdgeKVNamespacesResponse struct {
		Namespaces []Namespace `json:"namespaces"`
	}

	// Namespace represents a namespace object and a request body used to create a namespace
	Namespace struct {
		Name        string `json:"namespace"`
		GeoLocation string `json:"geoLocation,omitempty"`
		Retention   *int   `json:"retentionInSeconds,omitempty"`
		GroupID     *int   `json:"groupId,omitempty"`
	}

	// UpdateNamespace represents a request body used to update a namespace
	UpdateNamespace struct {
		Name      string `json:"namespace"`
		Retention *int   `json:"retentionInSeconds"`
		GroupID   *int   `json:"groupId"`
	}

	// NamespaceNetwork represents available namespace network types
	NamespaceNetwork string

	// DeleteEdgeKVNamespaceRequest represents the request to delete a namespace.
	DeleteEdgeKVNamespaceRequest struct {
		// Network specifies the network environment to execute the API request on.
		Network NamespaceNetwork

		// Name is a unique identifier for each namespace.
		Name string

		// Sync specifies whether to delete the namespace synchronously or asynchronously.
		Sync bool
	}

	// DeleteEdgeKVNamespacesResponse represents a response object returned when deleting a namespace.
	DeleteEdgeKVNamespacesResponse struct {
		// ScheduledDeleteTime is the time when the namespace will be deleted for asynchronous deletion.
		// For synchronous deletion, it will be nil.
		ScheduledDeleteTime *time.Time `json:"scheduledDeleteTime"`
	}
)

const (
	// NamespaceStagingNetwork is the staging network
	NamespaceStagingNetwork NamespaceNetwork = "staging"
	// NamespaceProductionNetwork is the production network
	NamespaceProductionNetwork NamespaceNetwork = "production"
)

// Validate validates ListEdgeKVNamespacesRequest
func (r ListEdgeKVNamespacesRequest) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(NamespaceStagingNetwork, NamespaceProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, NamespaceStagingNetwork, NamespaceProductionNetwork))),
	}.Filter()
}

// Validate validates GetEdgeKVNamespaceRequest
func (r GetEdgeKVNamespaceRequest) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(NamespaceStagingNetwork, NamespaceProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, NamespaceStagingNetwork, NamespaceProductionNetwork))),
		"Name": validation.Validate(r.Name, validation.Required, validation.Length(1, 32)),
	}.Filter()
}

// Validate validates CreateEdgeKVNamespaceRequest
func (r CreateEdgeKVNamespaceRequest) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(NamespaceStagingNetwork, NamespaceProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, NamespaceStagingNetwork, NamespaceProductionNetwork))),
		"Name":      validation.Validate(r.Name, validation.Required, validation.Length(1, 32)),
		"Retention": validation.Validate(r.Retention, validation.By(validateRetention)),
		"GroupID":   validation.Validate(r.GroupID, validation.By(validateGroupID)),
	}.Filter()
}

// Validate validates UpdateEdgeKVNamespaceRequest
func (r UpdateEdgeKVNamespaceRequest) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(NamespaceStagingNetwork, NamespaceProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, NamespaceStagingNetwork, NamespaceProductionNetwork))),
		"Name":      validation.Validate(r.Name, validation.Required, validation.Length(1, 32)),
		"Retention": validation.Validate(r.Retention, validation.By(validateRetention)),
		"GroupID":   validation.Validate(r.GroupID, validation.By(validateGroupID)),
	}.Filter()
}

// Validate validates DeleteEdgeKVNamespaceRequest
func (r DeleteEdgeKVNamespaceRequest) Validate() error {
	return validation.Errors{
		"Network": validation.Validate(r.Network, validation.Required, validation.In(NamespaceStagingNetwork, NamespaceProductionNetwork).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'", r.Network, NamespaceStagingNetwork, NamespaceProductionNetwork))),
		"Name": validation.Validate(r.Name, validation.Required, validation.Length(1, 32)),
	}.Filter()
}

func validateRetention(value interface{}) error {
	v, ok := value.(*int)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be *int", value)
	}
	if v == nil {
		return fmt.Errorf("cannot be blank")
	}
	if (*v < 86400 && *v != 0) || *v > 315360000 {
		return fmt.Errorf("a non zero value specified for retention period cannot be less than 86400 or more than 315360000")
	}
	return nil
}

func validateGroupID(value interface{}) error {
	v, ok := value.(*int)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be *int", value)
	}
	if v == nil {
		return fmt.Errorf("cannot be blank")
	}
	if *v < 0 {
		return fmt.Errorf("cannot be less than 0")
	}
	return nil
}

var (
	// ErrListEdgeKVNamespace is returned when ListEdgeKVNamespaces fails
	ErrListEdgeKVNamespace = errors.New("list EdgeKV namespaces")
	// ErrGetEdgeKVNamespace is returned when GetEdgeKVNamespace fails
	ErrGetEdgeKVNamespace = errors.New("get an EdgeKV namespace")
	// ErrCreateEdgeKVNamespace is returned when CreateEdgeKVNamespace fails
	ErrCreateEdgeKVNamespace = errors.New("create an EdgeKV namespace")
	// ErrUpdateEdgeKVNamespace is returned when UpdateEdgeKVNamespace fails
	ErrUpdateEdgeKVNamespace = errors.New("update an EdgeKV namespace")
	// ErrDeleteEdgeKVNamespace is returned when DeleteEdgeKVNamespace fails
	ErrDeleteEdgeKVNamespace = errors.New("delete an EdgeKV namespace")
)

func (e *edgeworkers) ListEdgeKVNamespaces(ctx context.Context, params ListEdgeKVNamespacesRequest) (*ListEdgeKVNamespacesResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListEdgeKVNamespaces")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListEdgeKVNamespace, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgekv/v1/networks/%s/namespaces", params.Network))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListEdgeKVNamespace, err.Error())
	}

	if params.Details {
		q := uri.Query()
		q.Add("details", "on")
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEdgeKVNamespace, err.Error())
	}

	var result ListEdgeKVNamespacesResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEdgeKVNamespace, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEdgeKVNamespace, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetEdgeKVNamespace(ctx context.Context, params GetEdgeKVNamespaceRequest) (*Namespace, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeKVNamespace")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeKVNamespace, ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s", params.Network, params.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeKVNamespace, err.Error())
	}

	var result Namespace
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeKVNamespace, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeKVNamespace, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) CreateEdgeKVNamespace(ctx context.Context, params CreateEdgeKVNamespaceRequest) (*Namespace, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateEdgeKVNamespace")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEdgeKVNamespace, ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces", params.Network)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeKVNamespace, err.Error())
	}

	var result Namespace
	resp, err := e.Exec(req, &result, params.Namespace)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeKVNamespace, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCreateEdgeKVNamespace, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) UpdateEdgeKVNamespace(ctx context.Context, params UpdateEdgeKVNamespaceRequest) (*Namespace, error) {
	logger := e.Log(ctx)
	logger.Debug("UpdateEdgeKVNamespace")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateEdgeKVNamespace, ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s", params.Network, params.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateEdgeKVNamespace, err.Error())
	}

	var result Namespace
	resp, err := e.Exec(req, &result, params.UpdateNamespace)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEdgeKVNamespace, err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateEdgeKVNamespace, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeleteEdgeKVNamespace(ctx context.Context, params DeleteEdgeKVNamespaceRequest) (*DeleteEdgeKVNamespacesResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("DeleteEdgeKVNamespace")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteEdgeKVNamespace, ErrStructValidation, err.Error())
	}

	uri, err := url.Parse(fmt.Sprintf("/edgekv/v1/networks/%s/namespaces/%s", params.Network, params.Name))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrDeleteEdgeKVNamespace, err.Error())
	}

	if params.Sync {
		q := url.Values{}
		q.Add("sync", "true")
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeKVNamespace, err.Error())
	}

	var result DeleteEdgeKVNamespacesResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeKVNamespace, err.Error())
	}
	defer session.CloseResponseBody(resp)

	var expectedCode int
	if params.Sync {
		expectedCode = http.StatusOK
	} else {
		expectedCode = http.StatusAccepted
	}

	if resp.StatusCode != expectedCode {
		return nil, fmt.Errorf("%s: %w", ErrDeleteEdgeKVNamespace, e.Error(resp))
	}

	return &result, nil
}
