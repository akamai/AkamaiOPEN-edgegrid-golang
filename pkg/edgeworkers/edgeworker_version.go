package edgeworkers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetEdgeWorkerVersionRequest contains parameters used to get an EdgeWorkerVersion
	GetEdgeWorkerVersionRequest EdgeWorkerVersionRequest

	// EdgeWorkerVersion represents an EdgeWorkerVersion object
	EdgeWorkerVersion struct {
		EdgeWorkerID   int    `json:"edgeWorkerId"`
		Version        string `json:"version"`
		AccountID      string `json:"accountId"`
		Checksum       string `json:"checksum"`
		SequenceNumber int    `json:"sequenceNumber"`
		CreatedBy      string `json:"createdBy"`
		CreatedTime    string `json:"createdTime"`
	}

	// ListEdgeWorkerVersionsRequest contains query parameters used to list EdgeWorkerVersions
	ListEdgeWorkerVersionsRequest struct {
		EdgeWorkerID int
	}

	// ListEdgeWorkerVersionsResponse represents a response object returned by ListEdgeWorkerVersions
	ListEdgeWorkerVersionsResponse struct {
		EdgeWorkerVersions []EdgeWorkerVersion `json:"versions"`
	}

	// GetEdgeWorkerVersionContentRequest contains parameters used to get content bundle of an EdgeWorkerVersion
	GetEdgeWorkerVersionContentRequest EdgeWorkerVersionRequest

	// CreateEdgeWorkerVersionRequest contains parameters used to create EdgeWorkerVersion
	CreateEdgeWorkerVersionRequest struct {
		EdgeWorkerID  int
		ContentBundle Bundle
	}

	// DeleteEdgeWorkerVersionRequest contains parameters used to delete an EdgeWorkerVersion
	DeleteEdgeWorkerVersionRequest EdgeWorkerVersionRequest

	// EdgeWorkerVersionRequest contains request parameters used by GetEdgeWorkerVersion, GetEdgeWorkerVersionContent and DeleteEdgeWorkerVersion
	EdgeWorkerVersionRequest struct {
		EdgeWorkerID int
		Version      string
	}

	// Bundle is the type for content bundle of an Edgeworker Version
	Bundle struct {
		io.Reader
	}
)

// Validate validates GetEdgeWorkerVersionRequest
func (g GetEdgeWorkerVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(g.EdgeWorkerID, validation.Required),
		"Version":      validation.Validate(g.Version, validation.Required),
	}.Filter()
}

// Validate validates ListEdgeWorkerVersionsRequest
func (g ListEdgeWorkerVersionsRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(g.EdgeWorkerID, validation.Required),
	}.Filter()
}

// Validate validates CreateEdgeWorkerVersionRequest
func (g CreateEdgeWorkerVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID":  validation.Validate(g.EdgeWorkerID, validation.Required),
		"ContentBundle": validation.Validate(g.ContentBundle.Reader, validation.Required),
	}.Filter()
}

// Validate validates GetEdgeWorkerVersionContentRequest
func (g GetEdgeWorkerVersionContentRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(g.EdgeWorkerID, validation.Required),
		"Version":      validation.Validate(g.Version, validation.Required),
	}.Filter()
}

// Validate validates DeleteEdgeWorkerVersionRequest
func (g DeleteEdgeWorkerVersionRequest) Validate() error {
	return validation.Errors{
		"EdgeWorkerID": validation.Validate(g.EdgeWorkerID, validation.Required),
		"Version":      validation.Validate(g.Version, validation.Required),
	}.Filter()
}

var (
	// ErrGetEdgeWorkerVersion is returned in case an error occurs on GetEdgeWorkerVersion operation
	ErrGetEdgeWorkerVersion = errors.New("get an EdgeWorker Version")
	// ErrListEdgeWorkerVersions is returned in case an error occurs on ListEdgeWorkerVersions operation
	ErrListEdgeWorkerVersions = errors.New("list EdgeWorkers Versions")
	// ErrGetEdgeWorkerVersionContent is returned in case an error occurs on GetEdgeWorkerVersionContent operation
	ErrGetEdgeWorkerVersionContent = errors.New("get an EdgeWorker Version Content Bundle")
	// ErrCreateEdgeWorkerVersion is returned in case an error occurs on CreateEdgeWorkerVersion operation
	ErrCreateEdgeWorkerVersion = errors.New("create an EdgeWorker Version")
	// ErrDeleteEdgeWorkerVersion is returned in case an error occurs on DeleteEdgeWorkerVersion operation
	ErrDeleteEdgeWorkerVersion = errors.New("delete an EdgeWorker Version")
)

func (e *edgeworkers) GetEdgeWorkerVersion(ctx context.Context, params GetEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeWorkerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeWorkerVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/versions/%s", params.EdgeWorkerID, params.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeWorkerVersion, err)
	}

	var result EdgeWorkerVersion
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeWorkerVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeWorkerVersion, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListEdgeWorkerVersions(ctx context.Context, params ListEdgeWorkerVersionsRequest) (*ListEdgeWorkerVersionsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListEdgeWorkerVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListEdgeWorkerVersions, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/versions", params.EdgeWorkerID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEdgeWorkerVersions, err)
	}

	var result ListEdgeWorkerVersionsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEdgeWorkerVersions, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEdgeWorkerVersions, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetEdgeWorkerVersionContent(ctx context.Context, params GetEdgeWorkerVersionContentRequest) (*Bundle, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeWorkerVersionContent")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeWorkerVersionContent, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/versions/%s/content", params.EdgeWorkerID, params.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeWorkerVersionContent, err)
	}

	req.Header.Add("Accept", "application/gzip")
	resp, err := e.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeWorkerVersionContent, err)
	}
	defer session.CloseResponseBody(resp)

	var result Bundle
	data, err := ioutil.ReadAll(resp.Body)
	result.Reader = bytes.NewBuffer(data)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response body: %s", ErrGetEdgeWorkerVersionContent, err)
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeWorkerVersionContent, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) CreateEdgeWorkerVersion(ctx context.Context, params CreateEdgeWorkerVersionRequest) (*EdgeWorkerVersion, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateEdgeWorkerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateEdgeWorkerVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/versions", params.EdgeWorkerID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, ioutil.NopCloser(params.ContentBundle))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeWorkerVersion, err)
	}

	req.Header.Add("Content-Type", "application/gzip")
	var result EdgeWorkerVersion
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeWorkerVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateEdgeWorkerVersion, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeleteEdgeWorkerVersion(ctx context.Context, params DeleteEdgeWorkerVersionRequest) error {
	e.Log(ctx).Debug("DeleteEdgeWorkerVersion")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrDeleteEdgeWorkerVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/edgeworkers/v1/ids/%d/versions/%s", params.EdgeWorkerID, params.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeWorkerVersion, err)
	}

	resp, err := e.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeWorkerVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteEdgeWorkerVersion, e.Error(resp))
	}

	return nil
}
