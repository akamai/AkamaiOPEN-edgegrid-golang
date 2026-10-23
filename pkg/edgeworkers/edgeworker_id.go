package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
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
		GroupID        int64
		ResourceTierID int
	}

	// ListEdgeWorkersIDResponse represents a response object returned by ListEdgeWorkersID
	ListEdgeWorkersIDResponse struct {
		EdgeWorkers []EdgeWorkerID `json:"edgeWorkerIds"`
	}

	// CreateEdgeWorkerIDRequest contains body parameters used to create EdgeWorkerID
	CreateEdgeWorkerIDRequest struct {
		Name           string `json:"name"`
		GroupID        int64  `json:"groupId"`
		ResourceTierID int    `json:"resourceTierId"`
	}

	// EdgeWorkerIDRequestBody contains body parameters used to update or clone EdgeWorkerID
	EdgeWorkerIDRequestBody struct {
		Name           string `json:"name"`
		GroupID        int64  `json:"groupId"`
		ResourceTierID int    `json:"resourceTierId"`
	}

	// UpdateEdgeWorkerIDRequest contains body and path parameters used to update EdgeWorkerID
	UpdateEdgeWorkerIDRequest struct {
		Body         EdgeWorkerIDRequestBody
		EdgeWorkerID int
	}

	// CloneEdgeWorkerIDRequest contains body and path parameters used to clone EdgeWorkerID
	CloneEdgeWorkerIDRequest struct {
		Body         EdgeWorkerIDRequestBody
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
		"Name":           validation.Validate(c.Body.Name, validation.Required),
		"GroupID":        validation.Validate(c.Body.GroupID, validation.Required),
		"ResourceTierID": validation.Validate(c.Body.ResourceTierID, validation.Required),
		"EdgeWorkerID":   validation.Validate(c.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates CloneEdgeWorkerIDRequest
func (c CloneEdgeWorkerIDRequest) Validate() error {
	return validation.Errors{
		"Name":           validation.Validate(c.Body.Name, validation.Required),
		"GroupID":        validation.Validate(c.Body.GroupID, validation.Required),
		"ResourceTierID": validation.Validate(c.Body.ResourceTierID, validation.Required),
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

	req, err := request.NewGet(ctx, "/edgeworkers/v1/ids/%d", params.EdgeWorkerID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeWorkerID, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeWorkerID, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListEdgeWorkersID(ctx context.Context, params ListEdgeWorkersIDRequest) (*ListEdgeWorkersIDResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListEdgeWorkersID")

	req, err := request.NewGet(ctx, "/edgeworkers/v1/ids").
		AddQueryParamIf("groupId", fmt.Sprintf("%d", params.GroupID), params.GroupID != 0).
		AddQueryParamIf("resourceTierId", fmt.Sprintf("%d", params.ResourceTierID), params.ResourceTierID != 0).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEdgeWorkersID, err)
	}

	var result ListEdgeWorkersIDResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEdgeWorkersID, err)
	}
	defer session.CloseResponseBody(resp)

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

	req, err := request.NewPost(ctx, "/edgeworkers/v1/ids").
		WithBody(params).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeWorkerID, err)
	}
	defer session.CloseResponseBody(resp)

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

	req, err := request.NewPut(ctx, "/edgeworkers/v1/ids/%d", params.EdgeWorkerID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEdgeWorkerID, err)
	}
	defer session.CloseResponseBody(resp)

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

	req, err := request.NewPost(ctx, "/edgeworkers/v1/ids/%d/clone", params.EdgeWorkerID).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCloneEdgeWorkerID, err)
	}

	var result EdgeWorkerID
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCloneEdgeWorkerID, err)
	}
	defer session.CloseResponseBody(resp)

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

	req, err := request.NewDelete(ctx, "/edgeworkers/v1/ids/%d", params.EdgeWorkerID).Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeWorkerID, err)
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeWorkerID, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteEdgeWorkerID, e.Error(resp))
	}

	return nil
}
