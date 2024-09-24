package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetEdgeWorkerIDRequest contains parameters used to get an EdgeWorkerID
	GetEdgeWorkerIDRequest struct {
		EdgeWorkerID int
	}

	// DeleteEdgeWorkerIDRequest contains parameters used to delete an EdgeWorkerID
	DeleteEdgeWorkerIDRequest struct {
		EdgeWorkerID int
	}

	// EdgeWorkerID represents an EdgeWorkerID object
	EdgeWorkerID struct {
		EdgeWorkerID       int    `json:"edgeWorkerId"`
		Name               string `json:"name"`
		AccountID          string `json:"accountId"`
		GroupID            int64  `json:"groupId"`
		ResourceTierID     int    `json:"resourceTierId"`
		SourceEdgeWorkerID int    `json:"sourceEdgeWorkerId,omitempty"`
		CreatedBy          string `json:"createdBy"`
		CreatedTime        string `json:"createdTime"`
		LastModifiedBy     string `json:"lastModifiedBy"`
		LastModifiedTime   string `json:"lastModifiedTime"`
	}

	// ListEdgeWorkersIDRequest contains query parameters used to list EdgeWorkerIDs
	ListEdgeWorkersIDRequest struct {
		GroupID        int
		ResourceTierID int
	}

	// ListEdgeWorkersIDResponse represents a response object returned by ListEdgeWorkersID
	ListEdgeWorkersIDResponse struct {
		EdgeWorkers []EdgeWorkerID `json:"edgeWorkerIds"`
	}

	// CreateEdgeWorkerIDRequest contains body parameters used to create EdgeWorkerID
	CreateEdgeWorkerIDRequest struct {
		Name           string `json:"name"`
		GroupID        int    `json:"groupId"`
		ResourceTierID int    `json:"resourceTierId"`
	}

	// EdgeWorkerIDBodyRequest contains body parameters used to update or clone EdgeWorkerID
	EdgeWorkerIDBodyRequest struct {
		Name           string `json:"name"`
		GroupID        int    `json:"groupId"`
		ResourceTierID int    `json:"resourceTierId"`
	}

	// UpdateEdgeWorkerIDRequest contains body and path parameters used to update EdgeWorkerID
	UpdateEdgeWorkerIDRequest struct {
		EdgeWorkerIDBodyRequest
		EdgeWorkerID int
	}

	// CloneEdgeWorkerIDRequest contains body and path parameters used to clone EdgeWorkerID
	CloneEdgeWorkerIDRequest struct {
		EdgeWorkerIDBodyRequest
		EdgeWorkerID int
	}
)

// Validate validates GetEdgeWorkerIDRequest
func (g GetEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(g.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates CreateEdgeWorkerIDRequest
func (c CreateEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"Name":           validation.Validate(c.Name, validation.Required),
		"GroupID":        validation.Validate(c.GroupID, validation.Required),
		"ResourceTierID": validation.Validate(c.ResourceTierID, validation.Required),
	}.Filter()
}

// Validate validates CreateEdgeWorkerIDRequest
func (c UpdateEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"Name":           validation.Validate(c.EdgeWorkerIDBodyRequest.Name, validation.Required),
		"GroupID":        validation.Validate(c.EdgeWorkerIDBodyRequest.GroupID, validation.Required),
		"ResourceTierID": validation.Validate(c.EdgeWorkerIDBodyRequest.ResourceTierID, validation.Required),
		"EdgeWorkerID":   validation.Validate(c.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates CloneEdgeWorkerIDRequest
func (c CloneEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"Name":           validation.Validate(c.EdgeWorkerIDBodyRequest.Name, validation.Required),
		"GroupID":        validation.Validate(c.EdgeWorkerIDBodyRequest.GroupID, validation.Required),
		"ResourceTierID": validation.Validate(c.EdgeWorkerIDBodyRequest.ResourceTierID, validation.Required),
		"EdgeWorkerID":   validation.Validate(c.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates DeleteEdgeWorkerIDRequest
func (d DeleteEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(d.EdgeWorkerID, validation.Required),
	}.Filter()
}

var (
	// ErrGetEdgeWorkerID is returned in case an error occurs on GetEdgeWorkerID operation
	ErrGetEdgeWorkerID = errors.New("get an EdgeWorker ID")
	// ErrListEdgeWorkersID is returned in case an error occurs on ListEdgeWorkersID operation
	ErrListEdgeWorkersID = errors.New("list EdgeWorkers IDs")
	// ErrCreateEdgeWorkerID is returned in case an error occurs on CreateEdgeWorkerID operation
	ErrCreateEdgeWorkerID = errors.New("create an EdgeWorker ID")
	// ErrUpdateEdgeWorkerID is returned in case an error occurs on UpdateEdgeWorkerID operation
	ErrUpdateEdgeWorkerID = errors.New("update an EdgeWorker ID")
	// ErrCloneEdgeWorkerID is returned in case an error occurs on CloneEdgeWroker operation
	ErrCloneEdgeWorkerID = errors.New("clone an EdgeWorker ID")
	// ErrDeleteEdgeWorkerID is returned in case an error occurs on DeleteEdgeWorkerID operation
	ErrDeleteEdgeWorkerID = errors.New("delete an EdgeWorker ID")
)

func (e *edgeworkers) GetEdgeWorkerID(ctx context.Context, params GetEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeWorkerID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeWorkerID, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d", params.EdgeWorkerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeWorkerID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeWorkerID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListEdgeWorkersID(ctx context.Context, params ListEdgeWorkersIDRequest) (*ListEdgeWorkersIDResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListEdgeWorkersID")

	uri, err := url.Parse("/edgeworkers/v1/ids")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListEdgeWorkersID, err)
	}
	q := uri.Query()
	if params.GroupID != 0 {
		q.Add("groupId", fmt.Sprintf("%d", params.GroupID))
	}
	if params.ResourceTierID != 0 {
		q.Add("resourceTierId", fmt.Sprintf("%d", params.ResourceTierID))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEdgeWorkersID, err)
	}

	var result ListEdgeWorkersIDResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEdgeWorkersID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEdgeWorkersID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) CreateEdgeWorkerID(ctx context.Context, params CreateEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateEdgeWorkerID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateEdgeWorkerID, ErrStructValidation, err)
	}

	uri, err := url.Parse("/edgeworkers/v1/ids")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateEdgeWorkerID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeWorkerID, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateEdgeWorkerID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) UpdateEdgeWorkerID(ctx context.Context, params UpdateEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	logger := e.Log(ctx)
	logger.Debug("UpdateEdgeWorkerID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateEdgeWorkerID, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateEdgeWorkerID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result, params.EdgeWorkerIDBodyRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEdgeWorkerID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateEdgeWorkerID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) CloneEdgeWorkerID(ctx context.Context, params CloneEdgeWorkerIDRequest) (*EdgeWorkerID, error) {
	logger := e.Log(ctx)
	logger.Debug("CloneEdgeWorkerID")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCloneEdgeWorkerID, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d/clone", params.EdgeWorkerID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCloneEdgeWorkerID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCloneEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result, params.EdgeWorkerIDBodyRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCloneEdgeWorkerID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCloneEdgeWorkerID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeleteEdgeWorkerID(ctx context.Context, params DeleteEdgeWorkerIDRequest) error {
	e.Log(ctx).Debug("DeleteEdgeWorkerID")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrDeleteEdgeWorkerID, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/ids/%d", params.EdgeWorkerID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeleteEdgeWorkerID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeWorkerID, err)
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeWorkerID, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteEdgeWorkerID, e.Error(resp))
	}

	return nil
}
