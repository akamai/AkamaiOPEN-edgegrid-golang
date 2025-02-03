package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CreateCredentialRequest contains request parameters for the CreateCredential endpoint.
	CreateCredentialRequest struct {
		ClientID string
	}

	// ListCredentialsRequest contains request parameters for the ListCredentials endpoint.
	ListCredentialsRequest struct {
		ClientID string
		Actions  bool
	}

	// GetCredentialRequest contains request parameters for the GetCredential endpoint.
	GetCredentialRequest struct {
		CredentialID int64
		ClientID     string
		Actions      bool
	}

	// UpdateCredentialRequest contains request parameters for the UpdateCredential endpoint.
	UpdateCredentialRequest struct {
		CredentialID int64
		ClientID     string
		Body         UpdateCredentialRequestBody
	}

	// UpdateCredentialRequestBody contains request body parameters for the UpdateCredential endpoint.
	UpdateCredentialRequestBody struct {
		Description string           `json:"description,omitempty"`
		ExpiresOn   time.Time        `json:"expiresOn"`
		Status      CredentialStatus `json:"status"`
	}

	// DeleteCredentialRequest contains request parameters for the DeleteCredential endpoint.
	DeleteCredentialRequest struct {
		CredentialID int64
		ClientID     string
	}

	// DeactivateCredentialRequest contains request parameters for the DeactivateCredential endpoint.
	DeactivateCredentialRequest struct {
		CredentialID int64
		ClientID     string
	}

	// DeactivateCredentialsRequest contains request parameters for the DeactivateCredentials endpoint.
	DeactivateCredentialsRequest struct {
		ClientID string
	}

	// CreateCredentialResponse holds response from the CreateCredentials endpoint.
	CreateCredentialResponse struct {
		ClientSecret string           `json:"clientSecret"`
		ClientToken  string           `json:"clientToken"`
		CreatedOn    time.Time        `json:"createdOn"`
		CredentialID int64            `json:"credentialId"`
		Description  string           `json:"description"`
		ExpiresOn    time.Time        `json:"expiresOn"`
		Status       CredentialStatus `json:"status"`
	}

	// ListCredentialsResponse holds response from the ListCredentials endpoint.
	ListCredentialsResponse []Credential

	// Credential represents single credential information.
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

	// CredentialActions describes the actions that can be performed on the credential.
	CredentialActions struct {
		Deactivate      bool `json:"deactivate"`
		Delete          bool `json:"delete"`
		Activate        bool `json:"activate"`
		EditDescription bool `json:"editDescription"`
		EditExpiration  bool `json:"editExpiration"`
	}

	// GetCredentialResponse holds response from the GetCredential endpoint.
	GetCredentialResponse Credential

	// UpdateCredentialResponse holds response from the UpdateCredential endpoint.
	UpdateCredentialResponse struct {
		Status      CredentialStatus `json:"status"`
		ExpiresOn   time.Time        `json:"expiresOn"`
		Description *string          `json:"description"`
	}

	// CredentialStatus represents the status of the credential.
	CredentialStatus string
)

const (
	// CredentialActive represents active credential.
	CredentialActive CredentialStatus = "ACTIVE"
	// CredentialInactive represents inactive credential.
	CredentialInactive CredentialStatus = "INACTIVE"
	// CredentialDeleted represents deleted credential.
	CredentialDeleted CredentialStatus = "DELETED"
)

// Validate validates CredentialStatus.
func (c CredentialStatus) Validate() error {
	return validation.In(CredentialActive, CredentialInactive, CredentialDeleted).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'",
			c, CredentialActive, CredentialInactive, CredentialDeleted)).
		Validate(c)
}

// Validate validates GetCredentialRequest.
func (r GetCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

// Validate validates UpdateCredentialRequest.
func (r UpdateCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
		"Body":         validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates UpdateCredentialRequestBody.
func (r UpdateCredentialRequestBody) Validate() error {
	return validation.Errors{
		"ExpiresOn": validation.Validate(r.ExpiresOn, validation.Required),
		"Status":    validation.Validate(r.Status, validation.Required),
	}.Filter()
}

// Validate validates DeleteCredentialRequest.
func (r DeleteCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

// Validate validates DeactivateCredentialRequest.
func (r DeactivateCredentialRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CredentialID": validation.Validate(r.CredentialID, validation.Required),
	})
}

var (
	// ErrCreateCredential is returned when CreateCredential fails.
	ErrCreateCredential = errors.New("create credential")
	// ErrListCredentials is returned when ListCredentials fails.
	ErrListCredentials = errors.New("list credentials")
	// ErrGetCredential is returned when GetCredential fails.
	ErrGetCredential = errors.New("get credential")
	// ErrUpdateCredential is returned when UpdateCredential fails.
	ErrUpdateCredential = errors.New("update credential")
	// ErrDeleteCredential is returned when DeleteCredential fails.
	ErrDeleteCredential = errors.New("delete credential")
	// ErrDeactivateCredential is returned when DeactivateCredential fails.
	ErrDeactivateCredential = errors.New("deactivate credential")
	// ErrDeactivateCredentials is returned when DeactivateCredentials fails.
	ErrDeactivateCredentials = errors.New("deactivate credentials")
)

func (i *iam) CreateCredential(ctx context.Context, params CreateCredentialRequest) (*CreateCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("CreateCredential")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListCredentials(ctx context.Context, params ListCredentialsRequest) (ListCredentialsResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListCredentials")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListCredentials, i.Error(resp))
	}

	return result, nil
}

func (i *iam) GetCredential(ctx context.Context, params GetCredentialRequest) (*GetCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetCredential")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) UpdateCredential(ctx context.Context, params UpdateCredentialRequest) (*UpdateCredentialResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdateCredential")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateCredential, ErrStructValidation, err)
	}

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	// Because API does not accept date without providing milliseconds, if there are no millisecond add a small duration to allow the request to
	// be processed. Only applicable when no milliseconds are provided, or they are equal to zero. Ticket for tracking: IDM-3347.
	if params.Body.ExpiresOn.Nanosecond() == 0 {
		params.Body.ExpiresOn = params.Body.ExpiresOn.Add(time.Nanosecond)
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
	resp, err := i.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateCredential, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateCredential, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) DeleteCredential(ctx context.Context, params DeleteCredentialRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeleteCredential")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteCredential, i.Error(resp))
	}

	return nil
}

func (i *iam) DeactivateCredential(ctx context.Context, params DeactivateCredentialRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeactivateCredential")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeactivateCredential, i.Error(resp))
	}

	return nil
}

func (i *iam) DeactivateCredentials(ctx context.Context, params DeactivateCredentialsRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeactivateCredentials")

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeactivateCredentials, i.Error(resp))
	}

	return nil
}
