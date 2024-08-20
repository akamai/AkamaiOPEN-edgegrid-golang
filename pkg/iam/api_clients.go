package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

		// ListAPIClients lists API clients an administrator can manage
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-clients
		ListAPIClients(ctx context.Context, params ListAPIClientsRequest) (ListAPIClientsResponse, error)

		// GetAPIClient provides details about an API client. If `ClientID` is not provided, it returns details about your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-api-client and https://techdocs.akamai.com/iam-api/reference/get-api-client-self
		GetAPIClient(ctx context.Context, params GetAPIClientRequest) (*GetAPIClientResponse, error)

		// CreateAPIClient creates a new API client. Optionally, it can automatically assign a credential for the client when creating it
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-api-clients
		CreateAPIClient(ctx context.Context, params CreateAPIClientRequest) (*CreateAPIClientResponse, error)

		// UpdateAPIClient updates an API client. If `ClientID` is not provided, it updates your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-api-clients and https://techdocs.akamai.com/iam-api/reference/put-api-clients-self
		UpdateAPIClient(ctx context.Context, params UpdateAPIClientRequest) (*UpdateAPIClientResponse, error)

		// DeleteAPIClient permanently deletes the API client, breaking any API connections with the client.
		// If `ClientID` is not provided, it deletes your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-api-client and https://techdocs.akamai.com/iam-api/reference/delete-api-client-self
		DeleteAPIClient(ctx context.Context, params DeleteAPIClientRequest) error
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

	// ListAPIClientsRequest contains the request parameters for the ListAPIClients endpoint
	ListAPIClientsRequest struct {
		Actions bool
	}

	// ListAPIClientsResponse describes the response of the ListAPIClients endpoint
	ListAPIClientsResponse []ListAPIClientsItem

	// ListAPIClientsItem represents information returned by the ListAPIClients endpoint for a single API client
	ListAPIClientsItem struct {
		AccessToken             string                 `json:"accessToken"`
		Actions                 *ListAPIClientsActions `json:"actions"`
		ActiveCredentialCount   int64                  `json:"activeCredentialCount"`
		AllowAccountSwitch      bool                   `json:"allowAccountSwitch"`
		AuthorizedUsers         []string               `json:"authorizedUsers"`
		CanAutoCreateCredential bool                   `json:"canAutoCreateCredential"`
		ClientDescription       string                 `json:"clientDescription"`
		ClientID                string                 `json:"clientId"`
		ClientName              string                 `json:"clientName"`
		ClientType              ClientType             `json:"clientType"`
		CreatedBy               string                 `json:"createdBy"`
		CreatedDate             time.Time              `json:"createdDate"`
		IsLocked                bool                   `json:"isLocked"`
		NotificationEmails      []string               `json:"notificationEmails"`
		ServiceConsumerToken    string                 `json:"serviceConsumerToken"`
	}

	// ListAPIClientsActions specifies activities available for the API client
	ListAPIClientsActions struct {
		Delete        bool `json:"delete"`
		DeactivateAll bool `json:"deactivateAll"`
		Edit          bool `json:"edit"`
		Lock          bool `json:"lock"`
		Transfer      bool `json:"transfer"`
		Unlock        bool `json:"unlock"`
	}

	// GetAPIClientRequest contains the request parameters for the GetAPIClient endpoint
	GetAPIClientRequest struct {
		ClientID    string
		Actions     bool
		GroupAccess bool
		APIAccess   bool
		Credentials bool
		IPACL       bool
	}

	// CreateAPIClientResponse describes the response of the CreateAPIClient endpoint
	CreateAPIClientResponse struct {
		AccessToken             string                      `json:"accessToken"`
		Actions                 *APIClientActions           `json:"actions"`
		ActiveCredentialCount   int64                       `json:"activeCredentialCount"`
		AllowAccountSwitch      bool                        `json:"allowAccountSwitch"`
		APIAccess               APIAccess                   `json:"apiAccess"`
		AuthorizedUsers         []string                    `json:"authorizedUsers"`
		BaseURL                 string                      `json:"baseURL"`
		CanAutoCreateCredential bool                        `json:"canAutoCreateCredential"`
		ClientDescription       string                      `json:"clientDescription"`
		ClientID                string                      `json:"clientId"`
		ClientName              string                      `json:"clientName"`
		ClientType              ClientType                  `json:"clientType"`
		CreatedBy               string                      `json:"createdBy"`
		CreatedDate             time.Time                   `json:"createdDate"`
		Credentials             []CreateAPIClientCredential `json:"credentials"`
		GroupAccess             GroupAccess                 `json:"groupAccess"`
		IPACL                   IPACL                       `json:"ipAcl"`
		IsLocked                bool                        `json:"isLocked"`
		NotificationEmails      []string                    `json:"notificationEmails"`
		PurgeOptions            PurgeOptions                `json:"purgeOptions"`
		ServiceProviderID       int64                       `json:"serviceProviderId"`
	}

	// GetAPIClientResponse describes the response of the GetAPIClient endpoint
	GetAPIClientResponse struct {
		AccessToken             string                `json:"accessToken"`
		Actions                 *APIClientActions     `json:"actions"`
		ActiveCredentialCount   int64                 `json:"activeCredentialCount"`
		AllowAccountSwitch      bool                  `json:"allowAccountSwitch"`
		APIAccess               APIAccess             `json:"apiAccess"`
		AuthorizedUsers         []string              `json:"authorizedUsers"`
		BaseURL                 string                `json:"baseURL"`
		CanAutoCreateCredential bool                  `json:"canAutoCreateCredential"`
		ClientDescription       string                `json:"clientDescription"`
		ClientID                string                `json:"clientId"`
		ClientName              string                `json:"clientName"`
		ClientType              ClientType            `json:"clientType"`
		CreatedBy               string                `json:"createdBy"`
		CreatedDate             time.Time             `json:"createdDate"`
		Credentials             []APIClientCredential `json:"credentials"`
		GroupAccess             GroupAccess           `json:"groupAccess"`
		IPACL                   IPACL                 `json:"ipAcl"`
		IsLocked                bool                  `json:"isLocked"`
		NotificationEmails      []string              `json:"notificationEmails"`
		PurgeOptions            PurgeOptions          `json:"purgeOptions"`
		ServiceProviderID       int64                 `json:"serviceProviderId"`
	}

	// APIClientActions specifies activities available for the API client
	APIClientActions struct {
		Delete            bool `json:"delete"`
		DeactivateAll     bool `json:"deactivateAll"`
		Edit              bool `json:"edit"`
		EditAPIs          bool `json:"editApis"`
		EditAuth          bool `json:"editAuth"`
		EditGroups        bool `json:"editGroups"`
		EditIPAcl         bool `json:"editIpAcl"`
		EditSwitchAccount bool `json:"editSwitchAccount"`
		Lock              bool `json:"lock"`
		Transfer          bool `json:"transfer"`
		Unlock            bool `json:"unlock"`
	}

	// APIAccess represents the APIs the API client can access
	APIAccess struct {
		AllAccessibleAPIs bool  `json:"allAccessibleApis"`
		APIs              []API `json:"apis"`
	}

	// API represents single Application Programming Interface (API)
	API struct {
		AccessLevel      AccessLevel `json:"accessLevel"`
		APIID            int64       `json:"apiId"`
		APIName          string      `json:"apiName"`
		Description      string      `json:"description"`
		DocumentationURL string      `json:"documentationUrl"`
		Endpoint         string      `json:"endPoint"`
	}

	// APIClientCredential represents single Credential returned by APIClient interfaces
	APIClientCredential struct {
		Actions      CredentialActions `json:"actions"`
		ClientToken  string            `json:"clientToken"`
		CreatedOn    time.Time         `json:"createdOn"`
		CredentialID int64             `json:"credentialId"`
		Description  string            `json:"description"`
		ExpiresOn    time.Time         `json:"expiresOn"`
		Status       CredentialStatus  `json:"status"`
	}

	// CreateAPIClientCredential represents single Credential returned by CreateAPIClient method
	CreateAPIClientCredential struct {
		Actions      CredentialActions `json:"actions"`
		ClientToken  string            `json:"clientToken"`
		ClientSecret string            `json:"clientSecret"`
		CreatedOn    time.Time         `json:"createdOn"`
		CredentialID int64             `json:"credentialId"`
		Description  string            `json:"description"`
		ExpiresOn    time.Time         `json:"expiresOn"`
		Status       CredentialStatus  `json:"status"`
	}

	// GroupAccess specifies the API client's group access.
	GroupAccess struct {
		CloneAuthorizedUserGroups bool          `json:"cloneAuthorizedUserGroups"`
		Groups                    []ClientGroup `json:"groups"`
	}

	// ClientGroup represents a group the API client can access.
	ClientGroup struct {
		GroupID         int64         `json:"groupId"`
		GroupName       string        `json:"groupName"`
		IsBlocked       bool          `json:"isBlocked"`
		ParentGroupID   int64         `json:"parentGroupId"`
		RoleDescription string        `json:"roleDescription"`
		RoleID          int64         `json:"roleId"`
		RoleName        string        `json:"roleName"`
		Subgroups       []ClientGroup `json:"subgroups"`
	}

	// IPACL specifies the API client's IP list restriction
	IPACL struct {
		CIDR   []string `json:"cidr"`
		Enable bool     `json:"enable"`
	}

	//PurgeOptions specifies the API clients configuration for access to the Fast Purge API
	PurgeOptions struct {
		CanPurgeByCacheTag bool         `json:"canPurgeByCacheTag"`
		CanPurgeByCPCode   bool         `json:"canPurgeByCpcode"`
		CPCodeAccess       CPCodeAccess `json:"cpcodeAccess"`
	}

	// CPCodeAccess represents the CP codes the API client can purge
	CPCodeAccess struct {
		AllCurrentAndNewCPCodes bool    `json:"allCurrentAndNewCpcodes"`
		CPCodes                 []int64 `json:"cpcodes"`
	}

	// CreateAPIClientRequest contains the request parameters for the CreateAPIClient endpoint
	CreateAPIClientRequest struct {
		AllowAccountSwitch      bool          `json:"allowAccountSwitch"`
		APIAccess               APIAccess     `json:"apiAccess"`
		AuthorizedUsers         []string      `json:"authorizedUsers"`
		CanAutoCreateCredential bool          `json:"canAutoCreateCredential"`
		ClientDescription       string        `json:"clientDescription"`
		ClientName              string        `json:"clientName"`
		ClientType              ClientType    `json:"clientType"`
		CreateCredential        bool          `json:"createCredential"`
		GroupAccess             GroupAccess   `json:"groupAccess"`
		IPACL                   *IPACL        `json:"ipAcl,omitempty"`
		NotificationEmails      []string      `json:"notificationEmails"`
		PurgeOptions            *PurgeOptions `json:"purgeOptions,omitempty"`
	}

	// UpdateAPIClientRequest contains the request parameters for the UpdateAPIClient endpoint
	UpdateAPIClientRequest struct {
		ClientID string
		Body     UpdateAPIClientBody
	}

	// UpdateAPIClientBody represents body params for UpdateAPIClient
	UpdateAPIClientBody struct {
		AllowAccountSwitch      bool          `json:"allowAccountSwitch"`
		APIAccess               APIAccess     `json:"apiAccess"`
		AuthorizedUsers         []string      `json:"authorizedUsers"`
		CanAutoCreateCredential bool          `json:"canAutoCreateCredential"`
		ClientDescription       string        `json:"clientDescription"`
		ClientName              string        `json:"clientName"`
		ClientType              ClientType    `json:"clientType"`
		GroupAccess             GroupAccess   `json:"groupAccess"`
		IPACL                   *IPACL        `json:"ipAcl,omitempty"`
		NotificationEmails      []string      `json:"notificationEmails"`
		PurgeOptions            *PurgeOptions `json:"purgeOptions,omitempty"`
	}

	// UpdateAPIClientResponse describes the response of the UpdateAPIClientResponse endpoint
	UpdateAPIClientResponse GetAPIClientResponse

	// DeleteAPIClientRequest contains the request parameters for the DeleteAPIClient endpoint
	DeleteAPIClientRequest struct {
		ClientID string
	}

	// AccessLevel represents the access level for API
	AccessLevel string
)

const (
	// ReadWriteLevel is the `READ-WRITE` access level
	ReadWriteLevel AccessLevel = "READ-WRITE"
	// ReadOnlyLevel is the `READ-ONLY` access level
	ReadOnlyLevel AccessLevel = "READ-ONLY"
)

// Validate validates UnlockAPIClientRequest
func (r UnlockAPIClientRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ClientID": validation.Validate(r.ClientID, validation.Required),
	})
}

// Validate validates CreateAPIClientRequest
func (r CreateAPIClientRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIAccess":       validation.Validate(r.APIAccess, validation.Required),
		"AuthorizedUsers": validation.Validate(r.AuthorizedUsers, validation.Required, validation.Length(1, 0)),
		"ClientType": validation.Validate(r.ClientType, validation.Required, validation.In(ClientClientType, UserClientType).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT'", r.ClientType))),
		"GroupAccess":  validation.Validate(r.GroupAccess, validation.Required),
		"PurgeOptions": validation.Validate(r.PurgeOptions),
	})
}

// Validate validates APIAccess
func (a APIAccess) Validate() error {
	return validation.Errors{
		"APIs": validation.Validate(a.APIs, validation.When(!a.AllAccessibleAPIs, validation.Required)),
	}.Filter()
}

// Validate validates API
func (a API) Validate() error {
	return validation.Errors{
		"AccessLevel": validation.Validate(a.AccessLevel, validation.Required, validation.In(ReadOnlyLevel, ReadWriteLevel).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'READ-ONLY' or 'READ-WRITE'", a.AccessLevel))),
		"APIID": validation.Validate(a.APIID, validation.Required),
	}.Filter()
}

// Validate validates GroupAccess
func (ga GroupAccess) Validate() error {
	return validation.Errors{
		"Groups": validation.Validate(ga.Groups, validation.When(!ga.CloneAuthorizedUserGroups, validation.Required)),
	}.Filter()
}

// Validate validates ClientGroup
func (cg ClientGroup) Validate() error {
	return validation.Errors{
		"GroupID": validation.Validate(cg.GroupID, validation.Required),
		"RoleID":  validation.Validate(cg.RoleID, validation.Required),
	}.Filter()
}

// Validate validates UpdateAPIClientRequest
func (r UpdateAPIClientRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Body": validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates UpdateAPIClientBody
func (r UpdateAPIClientBody) Validate() error {
	return validation.Errors{
		"APIAccess":       validation.Validate(r.APIAccess, validation.Required),
		"AuthorizedUsers": validation.Validate(r.AuthorizedUsers, validation.Required, validation.Length(1, 0)),
		"ClientType": validation.Validate(r.ClientType, validation.Required, validation.In(ClientClientType, UserClientType).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CLIENT' or 'USER_CLIENT'", r.ClientType))),
		"GroupAccess":  validation.Validate(r.GroupAccess, validation.Required),
		"PurgeOptions": validation.Validate(r.PurgeOptions),
	}.Filter()
}

// Validate validates PurgeOptions
func (po PurgeOptions) Validate() error {
	return validation.Errors{
		"CPCodeAccess": validation.Validate(po.CPCodeAccess),
	}.Filter()
}

// Validate validates UpdateAPIClientBody
func (ca CPCodeAccess) Validate() error {
	return validation.Errors{
		"CPCodes": validation.Validate(ca.CPCodes, validation.When(!ca.AllCurrentAndNewCPCodes, validation.NotNil)),
	}.Filter()
}

var (
	// ErrLockAPIClient is returned when LockAPIClient fails
	ErrLockAPIClient = errors.New("lock api client")
	// ErrUnlockAPIClient is returned when UnlockAPIClient fails
	ErrUnlockAPIClient = errors.New("unlock api client")
	// ErrListAPIClients is returned when ListAPIClients fails
	ErrListAPIClients = errors.New("list api clients")
	// ErrGetAPIClient is returned when GetAPIClient fails
	ErrGetAPIClient = errors.New("get api client")
	// ErrCreateAPIClient is returned when CreateAPIClient fails
	ErrCreateAPIClient = errors.New("create api client")
	// ErrUpdateAPIClient is returned when UpdateAPIClient fails
	ErrUpdateAPIClient = errors.New("update api client")
	// ErrDeleteAPIClient is returned when DeleteAPIClient fails
	ErrDeleteAPIClient = errors.New("delete api client")
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

func (i *iam) ListAPIClients(ctx context.Context, params ListAPIClientsRequest) (ListAPIClientsResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAPIClients")

	uri, err := url.Parse("/identity-management/v3/api-clients")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAPIClients, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAPIClients, err)
	}

	var result ListAPIClientsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAPIClients, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAPIClients, i.Error(resp))
	}

	return result, nil

}

func (i *iam) GetAPIClient(ctx context.Context, params GetAPIClientRequest) (*GetAPIClientResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetAPIClient")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAPIClient, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	q.Add("groupAccess", strconv.FormatBool(params.GroupAccess))
	q.Add("apiAccess", strconv.FormatBool(params.APIAccess))
	q.Add("credentials", strconv.FormatBool(params.Credentials))
	q.Add("ipAcl", strconv.FormatBool(params.IPACL))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAPIClient, err)
	}

	var result GetAPIClientResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetAPIClient, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAPIClient, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) CreateAPIClient(ctx context.Context, params CreateAPIClientRequest) (*CreateAPIClientResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateAPIClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateAPIClient, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/identity-management/v3/api-clients", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateAPIClient, err)
	}

	var result CreateAPIClientResponse
	resp, err := i.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateAPIClient, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateAPIClient, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateAPIClient(ctx context.Context, params UpdateAPIClientRequest) (*UpdateAPIClientResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateAPIClient")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateAPIClient, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("/identity-management/v3/api-clients/%s", params.ClientID), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateAPIClient, err)
	}

	var result UpdateAPIClientResponse
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateAPIClient, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateAPIClient, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) DeleteAPIClient(ctx context.Context, params DeleteAPIClientRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeleteAPIClient")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/identity-management/v3/api-clients/%s", params.ClientID), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteAPIClient, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteAPIClient, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteAPIClient, i.Error(resp))
	}

	return nil
}
