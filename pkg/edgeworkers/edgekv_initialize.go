package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// EdgeKVInitialize is EdgeKV Initialize API interface
type (
	EdgeKVInitialize interface {
		// InitializeEdgeKV Initialize the EdgeKV database
		//
		// See: https://techdocs.akamai.com/edgekv/reference/put_initialize
		InitializeEdgeKV(ctx context.Context) (*EdgeKVInitializationStatus, error)

		// GetEdgeKVInitializationStatus is used to check on the current initialization status
		//
		// See: https://techdocs.akamai.com/edgekv/reference/get_initialize
		GetEdgeKVInitializationStatus(ctx context.Context) (*EdgeKVInitializationStatus, error)
	}

	// EdgeKVInitializationStatus represents a response object returned by InitializeEdgeKV and GetEdgeKVInitializeStatus
	EdgeKVInitializationStatus struct {
		AccountStatus    string `json:"accountStatus"`
		CPCode           string `json:"cpcode"`
		ProductionStatus string `json:"productionStatus"`
		StagingStatus    string `json:"stagingStatus"`
	}
)

var (
	// ErrInitializeEdgeKV is returned in case an error occurs on InitializeEdgeKV operation
	ErrInitializeEdgeKV = errors.New("initialize EdgeKV")

	// ErrGetEdgeKVInitialize is returned in case an error occurs on GetEdgeKVInitializeStatus operation
	ErrGetEdgeKVInitialize = errors.New("get EdgeKV initialization status")
)

func (e *edgeworkers) InitializeEdgeKV(ctx context.Context) (*EdgeKVInitializationStatus, error) {
	logger := e.Log(ctx)
	logger.Debug("InitializeEdgeKV")

	uri := "/edgekv/v1/initialize"
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrInitializeEdgeKV, err)
	}

	var result EdgeKVInitializationStatus
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrInitializeEdgeKV, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrInitializeEdgeKV, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetEdgeKVInitializationStatus(ctx context.Context) (*EdgeKVInitializationStatus, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeKVInitializationStatus")

	uri := "/edgekv/v1/initialize"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeKVInitialize, err)
	}

	var result EdgeKVInitializationStatus
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeKVInitialize, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeKVInitialize, e.Error(resp))
	}

	return &result, nil
}
