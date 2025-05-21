package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetAccessKeyVersionStatusResponse contains response from GetAccessKeyVersionStatus
	GetAccessKeyVersionStatusResponse struct {
		AccessKeyVersion *KeyVersion    `json:"accessKeyVersion"`
		ProcessingStatus ProcessingType `json:"processingStatus"`
		RequestDate      time.Time      `json:"requestDate"`
		RequestedBy      string         `json:"requestedBy"`
	}

	// GetAccessKeyVersionStatusRequest holds parameters for GetAccessKeyVersionStatus
	GetAccessKeyVersionStatusRequest struct {
		RequestID int64
	}

	// CreateAccessKeyVersionRequest holds parameters for CreateAccessKeyVersion
	CreateAccessKeyVersionRequest struct {
		AccessKeyUID int64
		Body         CreateAccessKeyVersionRequestBody
	}

	// CreateAccessKeyVersionRequestBody holds body parameters for CreateAccessKeyVersion
	CreateAccessKeyVersionRequestBody struct {
		CloudAccessKeyID     string `json:"cloudAccessKeyId"`
		CloudSecretAccessKey string `json:"cloudSecretAccessKey"`
	}

	// CreateAccessKeyVersionResponse contains response from CreateAccessKeyVersion
	CreateAccessKeyVersionResponse struct {
		RequestID  int64 `json:"requestId"`
		RetryAfter int64 `json:"retryAfter"`
	}

	// GetAccessKeyVersionRequest holds parameters for GetAccessKeyVersion
	GetAccessKeyVersionRequest struct {
		Version      int64
		AccessKeyUID int64
	}

	// GetAccessKeyVersionResponse contains response from GetAccessKeyVersion
	GetAccessKeyVersionResponse AccessKeyVersion

	// ListAccessKeyVersionsRequest holds parameters for ListAccessKeyVersion
	ListAccessKeyVersionsRequest struct {
		AccessKeyUID int64
	}

	// ListAccessKeyVersionsResponse contains response from ListAccessKeyVersions
	ListAccessKeyVersionsResponse struct {
		AccessKeyVersions []AccessKeyVersion `json:"accessKeyVersions"`
	}

	// DeleteAccessKeyVersionRequest hold parameters for DeleteAccessKeyVersion
	DeleteAccessKeyVersionRequest struct {
		Version      int64
		AccessKeyUID int64
	}

	// DeleteAccessKeyVersionResponse contains response from DeleteAccessKeyVersion
	DeleteAccessKeyVersionResponse AccessKeyVersion

	// AccessKeyVersion holds information about access key version
	AccessKeyVersion struct {
		AccessKeyUID     int64            `json:"accessKeyUid"`
		CloudAccessKeyID *string          `json:"cloudAccessKeyId"`
		CreatedBy        string           `json:"createdBy"`
		CreatedTime      time.Time        `json:"createdTime"`
		DeploymentStatus DeploymentStatus `json:"deploymentStatus"`
		Version          int64            `json:"version"`
		VersionGUID      string           `json:"versionGuid"`
	}

	// DeploymentStatus represents deployment information
	DeploymentStatus string
)

const (
	// PendingActivation represents pending activation deployment status of access key version
	PendingActivation DeploymentStatus = "PENDING_ACTIVATION"
	// Active represents activated deployment status of access key version
	Active DeploymentStatus = "ACTIVE"
	// PendingDeletion represents pending deletion deployment status of access key version
	PendingDeletion DeploymentStatus = "PENDING_DELETION"
)

// Validate validates GetAccessKeyVersionStatusRequest
func (r GetAccessKeyVersionStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates CreateAccessKeyVersionRequest
func (r CreateAccessKeyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
		"Body":         validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates CreateAccessKeyVersionRequestBody
func (r CreateAccessKeyVersionRequestBody) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CloudAccessKeyID":     validation.Validate(r.CloudAccessKeyID, validation.Required),
		"CloudSecretAccessKey": validation.Validate(r.CloudSecretAccessKey, validation.Required),
	})
}

// Validate validates GetAccessKeyVersionRequest
func (r GetAccessKeyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
		"Version":      validation.Validate(r.Version, validation.Required),
	})
}

// Validate validates ListAccessKeyVersionsRequest
func (r ListAccessKeyVersionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
	})
}

// Validate validates DeleteAccessKeyVersionRequest
func (r DeleteAccessKeyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
		"Version":      validation.Validate(r.Version, validation.Required),
	})
}

var (
	// ErrGetAccessKeyVersionStatus is returned when GetAccessKeyVersionStatus fails
	ErrGetAccessKeyVersionStatus = errors.New("get the status of an access key version")
	// ErrCreateAccessKeyVersion is returned when CreateAccessKeyVersion fails
	ErrCreateAccessKeyVersion = errors.New("create access key version")
	// ErrGetAccessKeyVersion is returned when GetAccessKeyVersion fails
	ErrGetAccessKeyVersion = errors.New("get access key version")
	// ErrListAccessKeyVersions is returned when ListAccessKeyVersions fails
	ErrListAccessKeyVersions = errors.New("list access key versions")
	// ErrDeleteAccessKeyVersion is returned when DeleteAccessKeyVersion fails
	ErrDeleteAccessKeyVersion = errors.New("delete access key version")
)

func (c *cloudaccess) GetAccessKeyVersionStatus(ctx context.Context, params GetAccessKeyVersionStatusRequest) (*GetAccessKeyVersionStatusResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKeyStatusVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKeyVersionStatus, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/cam/v1/access-key-version-create-requests/%d", params.RequestID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAccessKeyVersionStatus, err)
	}

	var result GetAccessKeyVersionStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetAccessKeyVersionStatus, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAccessKeyVersionStatus, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) CreateAccessKeyVersion(ctx context.Context, params CreateAccessKeyVersionRequest) (*CreateAccessKeyVersionResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateAccessKeyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateAccessKeyVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cam/v1/access-keys/%d/versions", params.AccessKeyUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateAccessKeyVersion, err)
	}

	var result CreateAccessKeyVersionResponse
	resp, err := c.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrCreateAccessKeyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrCreateAccessKeyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) GetAccessKeyVersion(ctx context.Context, params GetAccessKeyVersionRequest) (*GetAccessKeyVersionResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKeyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKeyVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cam/v1/access-keys/%d/versions/%d", params.AccessKeyUID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAccessKeyVersion, err)
	}

	var result GetAccessKeyVersionResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetAccessKeyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAccessKeyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) ListAccessKeyVersions(ctx context.Context, params ListAccessKeyVersionsRequest) (*ListAccessKeyVersionsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListAccessKeyVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListAccessKeyVersions, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cam/v1/access-keys/%d/versions", params.AccessKeyUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAccessKeyVersions, err)
	}

	var result ListAccessKeyVersionsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListAccessKeyVersions, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAccessKeyVersions, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudaccess) DeleteAccessKeyVersion(ctx context.Context, params DeleteAccessKeyVersionRequest) (*DeleteAccessKeyVersionResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("DeleteAccessKeyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteAccessKeyVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cam/v1/access-keys/%d/versions/%d", params.AccessKeyUID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteAccessKeyVersion, err)
	}

	var result DeleteAccessKeyVersionResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrDeleteAccessKeyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrDeleteAccessKeyVersion, c.Error(resp))
	}

	return &result, nil
}
