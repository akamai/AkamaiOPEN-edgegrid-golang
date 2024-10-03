package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetAccessKeyStatusResponse contains response from GetAccessKeyStatus
	GetAccessKeyStatusResponse struct {
		AccessKey        *KeyLink            `json:"accessKey"`
		AccessKeyVersion *KeyVersion         `json:"accessKeyVersion"`
		ProcessingStatus ProcessingType      `json:"processingStatus"`
		Request          *RequestInformation `json:"request"`
		RequestDate      time.Time           `json:"requestDate"`
		RequestID        int64               `json:"requestId"`
		RequestedBy      string              `json:"requestedBy"`
	}

	// GetAccessKeyStatusRequest holds parameters for GetAccessKeyStatus
	GetAccessKeyStatusRequest struct {
		RequestID int64
	}

	// CreateAccessKeyRequest holds request body for CreateAccessKey
	CreateAccessKeyRequest struct {
		AccessKeyName        string        `json:"accessKeyName"`
		AuthenticationMethod string        `json:"authenticationMethod"`
		ContractID           string        `json:"contractId"`
		Credentials          Credentials   `json:"credentials"`
		GroupID              int64         `json:"groupId"`
		NetworkConfiguration SecureNetwork `json:"networkConfiguration"`
	}

	// Credentials holds information used to sign API requests
	Credentials struct {
		CloudAccessKeyID     string `json:"cloudAccessKeyId"`
		CloudSecretAccessKey string `json:"cloudSecretAccessKey"`
	}

	// CreateAccessKeyResponse contains response from CreateAccessKey
	CreateAccessKeyResponse struct {
		RequestID  int64 `json:"requestId,omitempty"`
		RetryAfter int64 `json:"retryAfter,omitempty"`
		Location   string
	}

	// AccessKeyRequest holds parameters for GetAccessKey, UpdateAccessKey, DeleteAccessKey
	AccessKeyRequest struct {
		AccessKeyUID int64
	}

	// AccessKeyResponse contains response from ListAccessKeys
	AccessKeyResponse struct {
		AccessKeyUID         int64          `json:"accessKeyUid"`
		AccessKeyName        string         `json:"accessKeyName"`
		AuthenticationMethod string         `json:"authenticationMethod"`
		NetworkConfiguration *SecureNetwork `json:"networkConfiguration"`
		LatestVersion        int64          `json:"latestVersion"`
		Groups               []Group        `json:"groups"`
		CreatedBy            string         `json:"createdBy"`
		CreatedTime          time.Time      `json:"createdTime"`
	}

	// GetAccessKeyResponse contains response from GetAccessKey
	GetAccessKeyResponse AccessKeyResponse

	// Group is an object to which the access key is assigned
	Group struct {
		ContractIDs []string `json:"contractIds"`
		GroupID     int64    `json:"groupId"`
		GroupName   *string  `json:"groupName"`
	}

	// ListAccessKeysRequest holds parameters for GetAccessKeys
	ListAccessKeysRequest struct {
		VersionGUID string
	}

	// ListAccessKeysResponse contains array of GetAccessKeyResponse
	ListAccessKeysResponse struct {
		AccessKeys []AccessKeyResponse `json:"accessKeys"`
	}

	// UpdateAccessKeyRequest holds request body for UpdateAccessKey
	UpdateAccessKeyRequest struct {
		AccessKeyName string `json:"accessKeyName,omitempty"`
	}

	// UpdateAccessKeyResponse contains response from UpdateAccessKey
	UpdateAccessKeyResponse struct {
		AccessKeyUID  int64  `json:"accessKeyUid"`
		AccessKeyName string `json:"accessKeyName"`
	}
)

// Validate validates GetAccessKeyStatusRequest
func (r GetAccessKeyStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"RequestID": validation.Validate(r.RequestID, validation.Required),
	})
}

// Validate validates CreateAccessKeyRequest
func (r CreateAccessKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyName":        validation.Validate(r.AccessKeyName, validation.Required),
		"AuthenticationMethod": validation.Validate(r.AuthenticationMethod, validation.Required),
		"ContractID":           validation.Validate(r.ContractID, validation.Required),
		"CloudAccessKeyID":     validation.Validate(r.Credentials.CloudAccessKeyID, validation.Required),
		"CloudSecretAccessKey": validation.Validate(r.Credentials.CloudSecretAccessKey, validation.Required),
		"GroupID":              validation.Validate(r.GroupID, validation.Required),
		"SecurityNetwork":      validation.Validate(r.NetworkConfiguration.SecurityNetwork, validation.Required),
	})
}

// Validate validates AccessKeyRequest
func (r AccessKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyUID": validation.Validate(r.AccessKeyUID, validation.Required),
	})
}

// Validate validates UpdateAccessKeyRequest
func (r UpdateAccessKeyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"AccessKeyName": validation.Validate(r.AccessKeyName, validation.Required, validation.Length(1, 50)),
	})
}

var (
	// ErrGetAccessKeyStatus is returned when GetAccessKeyStatus fails
	ErrGetAccessKeyStatus = errors.New("get the status of an access key")

	// ErrCreateAccessKey  is returned when CreateAccessKey fails
	ErrCreateAccessKey = errors.New("create an access key")

	// ErrGetAccessKey  is returned when GetAccessKey fails
	ErrGetAccessKey = errors.New("get an access key")

	// ErrUpdateAccessKey  is returned when UpdateAccessKey fails
	ErrUpdateAccessKey = errors.New("update an access key")

	// ErrDeleteAccessKey  is returned when DeleteAccessKey fails
	ErrDeleteAccessKey = errors.New("delete an access key")
)

func (c *cloudaccess) GetAccessKeyStatus(ctx context.Context, params GetAccessKeyStatusRequest) (*GetAccessKeyStatusResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKeyStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKeyStatus, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/cam/v1/access-key-create-requests/%d", params.RequestID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKeyStatus request: %w", err)
	}

	var result GetAccessKeyStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAccessKeyStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) CreateAccessKey(ctx context.Context, accessKey CreateAccessKeyRequest) (*CreateAccessKeyResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateAccessKey")

	if err := accessKey.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateAccessKey, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cam/v1/access-keys")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateAccessKey request: %w", err)
	}

	var result CreateAccessKeyResponse
	resp, err := c.Exec(req, &result, accessKey)
	if err != nil {
		return nil, fmt.Errorf("CreateAccessKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, c.Error(resp)
	}

	result.Location = resp.Header.Get("Location")

	return &result, nil
}

func (c *cloudaccess) GetAccessKey(ctx context.Context, params AccessKeyRequest) (*GetAccessKeyResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetAccessKey")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAccessKey, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/access-keys/%d", params.AccessKeyUID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKey request: %w", err)
	}

	var result GetAccessKeyResponse
	resp, err := c.Exec(req, &result)
	if err != nil {

		return nil, fmt.Errorf("GetAccessKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) ListAccessKeys(ctx context.Context, params ListAccessKeysRequest) (*ListAccessKeysResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListAccessKeys")

	uri, err := url.Parse("/cam/v1/access-keys")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	q := uri.Query()
	if params.VersionGUID != "" {
		q.Add("versionGuid", params.VersionGUID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListAccessKeys request: %w", err)
	}

	var result ListAccessKeysResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListAccessKeys request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) DeleteAccessKey(ctx context.Context, params AccessKeyRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteAccessKey")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteAccessKey, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/access-keys/%d", params.AccessKeyUID))
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteAccessKey request: %s", err.Error())
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteAccessKey request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusNoContent {
		return c.Error(resp)
	}

	return nil
}

func (c *cloudaccess) UpdateAccessKey(ctx context.Context, accessKey UpdateAccessKeyRequest, params AccessKeyRequest) (*UpdateAccessKeyResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateAccessKey")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateAccessKey, ErrStructValidation, err)
	}
	if err := accessKey.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateAccessKey, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cam/v1/access-keys/%d", params.AccessKeyUID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAccessKey request: %w", err)
	}

	var result UpdateAccessKeyResponse
	resp, err := c.Exec(req, &result, accessKey)
	if err != nil {
		return nil, fmt.Errorf("UpdateAccessKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil

}
