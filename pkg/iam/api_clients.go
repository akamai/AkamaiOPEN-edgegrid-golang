package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// APIClients is the IAM API clients interface
	APIClients interface {
		// LockAPIClient locks an API client based on `ClientID` parameter. If `ClientID` is not provided, it locks your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-lock-api-client, https://techdocs.akamai.com/iam-api/reference/put-lock-api-client-self
		LockAPIClient(ctx context.Context, params LockAPIClientRequest) (*LockAPIClientResponse, error)

		// UnlockAPIClient unlocks an API client
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-unlock-api-client
		UnlockAPIClient(ctx context.Context, params UnlockAPIClientRequest) (*UnlockAPIClientResponse, error)
	}

	// LockAPIClientRequest contains the request parameters for the LockAPIClient operation
	LockAPIClientRequest struct {
		ClientID string
	}

	// UnlockAPIClientRequest contains the request parameters for the UnlockAPIClient endpoint
	UnlockAPIClientRequest struct {
		ClientID string
	}

	// LockAPIClientResponse holds the response data from LockAPIClient
	LockAPIClientResponse APIClient

	// UnlockAPIClientResponse holds the response data from UnlockAPIClient
	UnlockAPIClientResponse APIClient

	// APIClient contains information about the API client
	APIClient struct {
		AccessToken             string    `json:"accessToken"`
		ActiveCredentialCount   int64     `json:"activeCredentialCount"`
		AllowAccountSwitch      bool      `json:"allowAccountSwitch"`
		AuthorizedUsers         []string  `json:"authorizedUsers"`
		CanAutoCreateCredential bool      `json:"canAutoCreateCredential"`
		ClientDescription       string    `json:"clientDescription"`
		ClientID                string    `json:"clientId"`
		ClientName              string    `json:"clientName"`
		ClientType              string    `json:"clientType"`
		CreatedBy               string    `json:"createdBy"`
		CreatedDate             time.Time `json:"createdDate"`
		IsLocked                bool      `json:"isLocked"`
		NotificationEmails      []string  `json:"notificationEmails"`
		ServiceConsumerToken    string    `json:"serviceConsumerToken"`
	}
)

// Validate validates UnlockAPIClientRequest
func (r UnlockAPIClientRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ClientID": validation.Validate(r.ClientID, validation.Required),
	})
}

var (
	// ErrLockAPIClient is returned when LockAPIClient fails
	ErrLockAPIClient = errors.New("lock api client")
	// ErrUnlockAPIClient is returned when UnlockAPIClient fails
	ErrUnlockAPIClient = errors.New("unlock api client")
)

func (i *iam) LockAPIClient(ctx context.Context, params LockAPIClientRequest) (*LockAPIClientResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("LockAPIClient")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/lock", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrLockAPIClient, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrLockAPIClient, err)
	}

	var result LockAPIClientResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrLockAPIClient, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrLockAPIClient, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UnlockAPIClient(ctx context.Context, params UnlockAPIClientRequest) (*UnlockAPIClientResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("UnlockAPIClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUnlockAPIClient, ErrStructValidation, err)
	}

	u, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/unlock", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUnlockAPIClient, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUnlockAPIClient, err)
	}

	var result UnlockAPIClientResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUnlockAPIClient, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUnlockAPIClient, i.Error(resp))
	}

	return &result, nil
}
