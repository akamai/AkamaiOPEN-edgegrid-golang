package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// G2OAccessKeyIDMinLength is the minimum length for a G2O access key ID
	G2OAccessKeyIDMinLength = 1
	// G2OAccessKeyIDMaxLength is the maximum length for a G2O access key ID
	G2OAccessKeyIDMaxLength = 8
	// G2OSecretAccessKeyMinLength is the minimum length for a G2O secret access key
	G2OSecretAccessKeyMinLength = 32
	// G2OSecretAccessKeyMaxLength is the maximum length for a G2O secret access key
	G2OSecretAccessKeyMaxLength = 64
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
		// Deprecated: AuthenticationMethod is no longer returned by the API and will always be empty.
		// It will be removed in the nearest breaking release.
		AuthenticationMethod string `json:"-"`
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
		CloudAccessKeyID     string `json:"cloudAccessKeyId,omitempty"`
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
		Note                 *string        `json:"note"`
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
		"AuthenticationMethod": validation.Validate(r.AuthenticationMethod, validation.Required, validation.In(string(AuthAWS), string(AuthGOOG), string(AuthAOS), string(AuthAVMCloudinary), string(AuthVPQueueIt), string(AuthG2O))),
		"ContractID":           validation.Validate(r.ContractID, validation.Required),
		"CloudAccessKeyID":     validation.Validate(r.Credentials.CloudAccessKeyID, validation.When(r.AuthenticationMethod != string(AuthVPQueueIt) && r.AuthenticationMethod != string(AuthAVMCloudinary), validation.Required), validation.When(r.AuthenticationMethod == string(AuthG2O), validation.Length(G2OAccessKeyIDMinLength, G2OAccessKeyIDMaxLength))),
		"CloudSecretAccessKey": validation.Validate(r.Credentials.CloudSecretAccessKey, validation.Required, validation.When(r.AuthenticationMethod == string(AuthG2O), validation.Length(G2OSecretAccessKeyMinLength, G2OSecretAccessKeyMaxLength))),
		"GroupID":              validation.Validate(r.GroupID, validation.Required),
		"SecurityNetwork":      validation.Validate(r.NetworkConfiguration.SecurityNetwork, validation.Required),
		"AdditionalCDN":        validation.Validate(r.NetworkConfiguration.AdditionalCDN, validation.When(r.AuthenticationMethod == string(AuthVPQueueIt) || r.AuthenticationMethod == string(AuthAVMCloudinary), validation.Nil)),
	})
}

// Validate validates SecureNetwork
func (r SecureNetwork) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"SecurityNetwork": validation.Validate(r.SecurityNetwork, validation.In(NetworkEnhanced, NetworkStandard)),
		"AdditionalCDN":   validation.Validate(r.AdditionalCDN, validation.When(r.AdditionalCDN != nil, validation.In(ChinaCDN, RussiaCDN))),
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
	_ validation.Validatable = GetAccessKeyStatusRequest{}
	_ validation.Validatable = CreateAccessKeyRequest{}
	_ validation.Validatable = SecureNetwork{}
	_ validation.Validatable = AccessKeyRequest{}
	_ validation.Validatable = UpdateAccessKeyRequest{}
)

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
		return nil, fmt.Errorf("%w: %w: %w", ErrGetAccessKeyStatus, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cam/v1/access-key-create-requests/%d", params.RequestID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKeyStatus request: %w", err)
	}

	var result GetAccessKeyStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetAccessKeyStatus request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) CreateAccessKey(ctx context.Context, accessKey CreateAccessKeyRequest) (*CreateAccessKeyResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateAccessKey")

	if err := accessKey.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrCreateAccessKey, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/cam/v1/access-keys").
		WithBody(accessKey).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateAccessKey request: %w", err)
	}

	var result CreateAccessKeyResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("CreateAccessKey request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

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
		return nil, fmt.Errorf("%w: %w: %w", ErrGetAccessKey, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cam/v1/access-keys/%d", params.AccessKeyUID).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAccessKey request: %w", err)
	}

	var result GetAccessKeyResponse
	resp, err := c.Exec(req, &result)
	if err != nil {

		return nil, fmt.Errorf("GetAccessKey request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) ListAccessKeys(ctx context.Context, params ListAccessKeysRequest) (*ListAccessKeysResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListAccessKeys")

	req, err := request.NewGet(ctx, "/cam/v1/access-keys").
		AddQueryParamIf("versionGuid", params.VersionGUID, params.VersionGUID != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ListAccessKeys request: %w", err)
	}

	var result ListAccessKeysResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListAccessKeys request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil
}

func (c *cloudaccess) DeleteAccessKey(ctx context.Context, params AccessKeyRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteAccessKey")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDeleteAccessKey, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/cam/v1/access-keys/%d", params.AccessKeyUID).Build()
	if err != nil {
		return fmt.Errorf("failed to create DeleteAccessKey request: %w", err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteAccessKey request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return c.Error(resp)
	}

	return nil
}

func (c *cloudaccess) UpdateAccessKey(ctx context.Context, accessKey UpdateAccessKeyRequest, params AccessKeyRequest) (*UpdateAccessKeyResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateAccessKey")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUpdateAccessKey, ErrStructValidation, err)
	}
	if err := accessKey.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUpdateAccessKey, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cam/v1/access-keys/%d", params.AccessKeyUID).
		WithBody(accessKey).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAccessKey request: %w", err)
	}

	var result UpdateAccessKeyResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("UpdateAccessKey request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, c.Error(resp)
	}

	return &result, nil

}
