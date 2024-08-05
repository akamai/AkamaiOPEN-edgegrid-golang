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
	// APIClientsCredentials is the IAM API clients credentials interface
	APIClientsCredentials interface {
		// CreateCredential creates a new credential for the API client.  If `ClientID` is not provided, it creates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credentials, https://techdocs.akamai.com/iam-api/reference/post-client-credentials
		CreateCredential(context.Context, CreateCredentialRequest) (*CreateCredentialResponse, error)

		// ListCredentials lists credentials for an API client. If `ClientID` is not provided, it lists credentials for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-self-credentials, https://techdocs.akamai.com/iam-api/reference/get-client-credentials
		ListCredentials(context.Context, ListCredentialsRequest) (ListCredentialsResponse, error)

		// GetCredential returns details about a specific credential for an API client. If `ClientID` is not provided, it gets credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-self-credential, https://techdocs.akamai.com/iam-api/reference/get-client-credential
		GetCredential(context.Context, GetCredentialRequest) (*GetCredentialResponse, error)

		// UpdateCredential updates a specific credential for an API client. If `ClientID` is not provided, it updates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/put-self-credential, https://techdocs.akamai.com/iam-api/reference/put-client-credential
		UpdateCredential(context.Context, UpdateCredentialRequest) (*UpdateCredentialResponse, error)

		// DeleteCredential deletes a specific credential from an API client. If `ClientID` is not provided, it deletes credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/delete-self-credential, https://techdocs.akamai.com/iam-api/reference/delete-client-credential
		DeleteCredential(context.Context, DeleteCredentialRequest) error

		// DeactivateCredential deactivates a specific credential for an API client. If `ClientID` is not provided, it deactivates credential for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credential-deactivate, https://techdocs.akamai.com/iam-api/reference/post-client-credential-deactivate
		DeactivateCredential(context.Context, DeactivateCredentialRequest) error

		// DeactivateCredentials deactivates all credentials for a specific API client. If `ClientID` is not provided, it deactivates all credentials for your API client.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/post-self-credentials-deactivate, https://techdocs.akamai.com/iam-api/reference/post-client-credentials-deactivate
		DeactivateCredentials(context.Context, DeactivateCredentialsRequest) error
	}

	// CreateCredentialRequest contains request parameters for CreateCredential operation
	CreateCredentialRequest struct {
		ClientID string
	}

	// ListCredentialsRequest contains request parameters for ListCredentials operation
	ListCredentialsRequest struct {
		ClientID string
		Actions  bool
	}

	// GetCredentialRequest contains request parameters for GetCredentials operation
	GetCredentialRequest struct {
		CredentialID int64
		ClientID     string
		Actions      bool
	}

	// UpdateCredentialRequest contains request parameters for UpdateCredential operation
	UpdateCredentialRequest struct {
		CredentialID int64
		ClientID     string
		RequestBody  UpdateCredentialRequestBody
	}

	// UpdateCredentialRequestBody contains request body parameters for UpdateCredential operation
	UpdateCredentialRequestBody struct {
		Description string           `json:"description,omitempty"`
		ExpiresOn   time.Time        `json:"expiresOn"`
		Status      CredentialStatus `json:"status"`
	}

	// DeleteCredentialRequest contains request parameters for DeleteCredential operation
	DeleteCredentialRequest struct {
		CredentialID int64
		ClientID     string
	}

	// DeactivateCredentialRequest contains request parameters for DeactivateCredential operation
	DeactivateCredentialRequest struct {
		CredentialID int64
		ClientID     string
	}

	// DeactivateCredentialsRequest contains request parameters for DeactivateCredentials operation
	DeactivateCredentialsRequest struct {
		ClientID string
	}

	// CreateCredentialResponse holds response from CreateCredentials operation
	CreateCredentialResponse struct {
		ClientSecret string           `json:"clientSecret"`
		ClientToken  string           `json:"clientToken"`
		CreatedOn    time.Time        `json:"createdOn"`
		CredentialID int64            `json:"credentialId"`
		Description  string           `json:"description"`
		ExpiresOn    time.Time        `json:"expiresOn"`
		Status       CredentialStatus `json:"status"`
	}

	// ListCredentialsResponse holds response from ListCredentials operation
	ListCredentialsResponse []Credential

	// Credential represents single credential information
	Credential struct {
		ClientToken      string             `json:"clientToken"`
		CreatedOn        time.Time          `json:"createdOn"`
		CredentialID     int64              `json:"credentialId"`
		Description      string             `json:"description"`
		ExpiresOn        time.Time          `json:"expiresOn"`
		Status           CredentialStatus   `json:"status"`
		MaxAllowedExpiry time.Time          `json:"maxAllowedExpiry"`
		Actions          *CredentialActions `json:"actions"`
	}

	// CredentialActions describes the actions that can be performed on the credential
	CredentialActions struct {
		Deactivate      bool `json:"deactivate"`
		Delete          bool `json:"delete"`
		Activate        bool `json:"activate"`
		EditDescription bool `json:"editDescription"`
		EditExpiration  bool `json:"editExpiration"`
	}

	// GetCredentialResponse holds response from GetCredential operation
	GetCredentialResponse Credential

	// UpdateCredentialResponse holds response from UpdateCredential operation
	UpdateCredentialResponse struct {
		Status      CredentialStatus `json:"status"`
		ExpiresOn   time.Time        `json:"expiresOn"`
		Description *string          `json:"description"`
	}

	// CredentialStatus represents the status of the credential
	CredentialStatus string
)

const (
	// CredentialActive represents active credential
	CredentialActive CredentialStatus = "ACTIVE"
	// CredentialInactive represents inactive credential
	CredentialInactive CredentialStatus = "INACTIVE"
	// CredentialDeleted represents deleted credential
	CredentialDeleted CredentialStatus = "DELETED"
)

// Validate validates GetCredentialRequest
func (r GetCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

// Validate validates UpdateCredentialRequest
func (r UpdateCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
		"RequestBody":  validation.Validate(r.RequestBody, validation.Required),
	})
}

// Validate validates UpdateCredentialRequestBody
func (r UpdateCredentialRequestBody) Validate() error {
	return validation.Errors{
		"ExpiresOn": validation.Validate(r.ExpiresOn, validation.Required),
		"Status":    validation.Validate(r.Status, validation.Required, validation.In(CredentialActive, CredentialInactive)),
	}.Filter()
}

// Validate validates DeleteCredentialRequest
func (r DeleteCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

// Validate validates DeactivateCredentialRequest
func (r DeactivateCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

var (
	// ErrCreateCredential is returned when CreateCredential fails
	ErrCreateCredential = errors.New("create credential")
	// ErrListCredentials is returned when ListCredentials fails
	ErrListCredentials = errors.New("list credentials")
	// ErrGetCredential is returned when GetCredential fails
	ErrGetCredential = errors.New("get credential")
	// ErrUpdateCredential is returned when UpdateCredential fails
	ErrUpdateCredential = errors.New("update credential")
	// ErrDeleteCredential is returned when DeleteCredential fails
	ErrDeleteCredential = errors.New("delete credential")
	// ErrDeactivateCredential is returned when DeactivateCredential fails
	ErrDeactivateCredential = errors.New("deactivate credential")
	// ErrDeactivateCredentials is returned when DeactivateCredentials fails
	ErrDeactivateCredentials = errors.New("deactivate credentials")
)

func (i *iam) CreateCredential(ctx context.Context, params CreateCredentialRequest) (*CreateCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("create credential")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateCredential, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCredential, err)
	}

	var result CreateCredentialResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCredential, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListCredentials(ctx context.Context, params ListCredentialsRequest) (ListCredentialsResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("list credentials")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCredentials, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCredentials, err)
	}

	var result ListCredentialsResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCredentials, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCredentials, i.Error(resp))
	}

	return result, nil
}

func (i *iam) GetCredential(ctx context.Context, params GetCredentialRequest) (*GetCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("get credential")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCredential, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials/%d", params.ClientID, params.CredentialID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCredential, err)
	}

	q := uri.Query()
	q.Add("actions", strconv.FormatBool(params.Actions))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCredential, err)
	}

	var result GetCredentialResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCredential, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateCredential(ctx context.Context, params UpdateCredentialRequest) (*UpdateCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("update credential")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateCredential, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	// Because API does not accept date without providing milliseconds, if there are no millisecond add a small duration to allow the request to
	// be processed. Only applicable when no milliseconds are provided, or they are equal to zero. Ticket for tracking: IDM-3347.
	if params.RequestBody.ExpiresOn.Nanosecond() == 0 {
		params.RequestBody.ExpiresOn = params.RequestBody.ExpiresOn.Add(time.Nanosecond)
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials/%d", params.ClientID, params.CredentialID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateCredential, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateCredential, err)
	}

	var result UpdateCredentialResponse
	resp, err := i.Exec(req, &result, params.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCredential, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) DeleteCredential(ctx context.Context, params DeleteCredentialRequest) error {
	logger := i.Log(ctx)
	logger.Debug("delete credential")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteCredential, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials/%d", params.ClientID, params.CredentialID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeleteCredential, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteCredential, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteCredential, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteCredential, i.Error(resp))
	}

	return nil
}

func (i *iam) DeactivateCredential(ctx context.Context, params DeactivateCredentialRequest) error {
	logger := i.Log(ctx)
	logger.Debug("deactivate credential")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeactivateCredential, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials/%d/deactivate", params.ClientID, params.CredentialID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeactivateCredential, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeactivateCredential, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeactivateCredential, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeactivateCredential, i.Error(resp))
	}

	return nil
}

func (i *iam) DeactivateCredentials(ctx context.Context, params DeactivateCredentialsRequest) error {
	logger := i.Log(ctx)
	logger.Debug("deactivate credentials")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/credentials/deactivate", params.ClientID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeactivateCredentials, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeactivateCredentials, err)
	}

	resp, err := i.Exec(req, nil, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeactivateCredentials, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeactivateCredentials, i.Error(resp))
	}

	return nil
}
