package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CreateEdgeKVAccessTokenRequest contains parameters used to create EdgeKV access token
	CreateEdgeKVAccessTokenRequest struct {
		// Whether to allow this token access to the Akamai production network
		AllowOnProduction bool `json:"allowOnProduction"`
		// Whether to allow this token access to the Akamai staging network
		AllowOnStaging bool `json:"allowOnStaging"`
		// Friendly name of the token. Used when retrieving tokens by name.
		Name string `json:"name"`
		// A list of namespace identifiers the token should have access to, plus the associated read, write, delete permissions
		NamespacePermissions NamespacePermissions `json:"namespacePermissions"`
		// A set of EdgeWorker IDs authorized to access EdgeKV via the token. By default, if you omit this array, the token authorizes access for all IDs.
		RestrictToEdgeWorkerIDs []string `json:"restrictToEdgeWorkerIds"`
	}

	// NamespacePermissions represents mapping between namespaces and permissions
	NamespacePermissions map[string][]Permission

	// Permissions represents set of permissions for namespace
	Permissions []Permission

	// Permission has possible values: `r` for read access, `w` for write access, `d` for delete access
	Permission string

	// EdgeKVAccessToken contains response from EdgeKV access token creation
	EdgeKVAccessToken struct {
		// The expiry date
		Expiry string `json:"expiry"`
		// The name assigned to the access token. You can't modify an access token name.
		Name string `json:"name"`
		// Internally generated unique identifier for the access token
		UUID string `json:"uuid"`
		// The IN_PROGRESS status indicates token activation is still in progress,
		// and it's not yet possible to make successful EdgeKV requests from EdgeWorkers that use the token.
		// Once activation completes, status is COMPLETE.
		// Otherwise, a value of ERROR indicates a problem that prevented activation.
		TokenActivationStatus *string `json:"tokenActivationStatus"`
		// Initial token creation date in ISO 8601 format.
		IssueDate *string `json:"issueDate"`
		// Most recent token refresh date in ISO 8601 format. A null value indicates the token has not yet been refreshed.
		LatestRefreshDate *string `json:"latestRefreshDate"`
		// Next scheduled date of the token refresh in ISO 8601 format.
		NextScheduledRefreshDate *string `json:"nextScheduledRefreshDate"`
	}

	// CreateEdgeKVAccessTokenResponse contains response from EdgeKV access token creation
	CreateEdgeKVAccessTokenResponse struct {
		AllowOnProduction        bool                 `json:"allowOnProduction"`
		AllowOnStaging           bool                 `json:"allowOnStaging"`
		CPCode                   string               `json:"cpcode"`
		Expiry                   string               `json:"expiry"`
		IssueDate                string               `json:"issueDate"`
		LatestRefreshDate        *string              `json:"latestRefreshDate"`
		Name                     string               `json:"name"`
		NamespacePermissions     NamespacePermissions `json:"namespacePermissions"`
		NextScheduledRefreshDate string               `json:"nextScheduledRefreshDate"`
		RestrictToEdgeWorkerIDs  []string             `json:"restrictToEdgeWorkerIds"`
		TokenActivationStatus    string               `json:"tokenActivationStatus"`
		UUID                     string               `json:"uuid"`
	}

	// GetEdgeKVAccessTokenRequest represents an TokenName object
	GetEdgeKVAccessTokenRequest struct {
		TokenName string
	}

	// GetEdgeKVAccessTokenResponse contains response from EdgeKV access token retrieval
	GetEdgeKVAccessTokenResponse CreateEdgeKVAccessTokenResponse

	// ListEdgeKVAccessTokensRequest contains request parameters for ListEdgeKVAccessTokens
	ListEdgeKVAccessTokensRequest struct {
		IncludeExpired bool
	}

	// ListEdgeKVAccessTokensResponse contains list of EdgeKV access tokens
	ListEdgeKVAccessTokensResponse struct {
		Tokens []EdgeKVAccessToken `json:"tokens"`
	}

	// DeleteEdgeKVAccessTokenRequest contains name of the EdgeKV access token to remove
	DeleteEdgeKVAccessTokenRequest struct {
		TokenName string
	}

	// DeleteEdgeKVAccessTokenResponse contains response after removal of EdgeKV access token
	DeleteEdgeKVAccessTokenResponse struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	}
)

const (
	// PermissionRead represents read permission
	PermissionRead Permission = "r"
	// PermissionWrite represents write permission
	PermissionWrite Permission = "w"
	// PermissionDelete represents delete permission
	PermissionDelete Permission = "d"
)

// Validate validates CreateEdgeKVAccessTokenRequest
func (c CreateEdgeKVAccessTokenRequest) Validate() error {
	namespaces := make([]string, 0)
	for name := range c.NamespacePermissions {
		namespaces = append(namespaces, name)
	}

	return validation.Errors{
		"AllowOnProduction":          validation.Validate(c.AllowOnProduction, validation.Required.When(!c.AllowOnStaging).Error("at least one of AllowOnProduction or AllowOnStaging has to be provided")),
		"AllowOnStaging":             validation.Validate(c.AllowOnStaging, validation.Required.When(!c.AllowOnProduction).Error("at least one of AllowOnProduction or AllowOnStaging has to be provided")),
		"Name":                       validation.Validate(c.Name, validation.Required, validation.Length(1, 32)),
		"NamespacePermissions.Names": validation.Validate(namespaces, validation.Required, validation.Each(validation.Required)),
		"NamespacePermissions": validation.Validate(c.NamespacePermissions,
			validation.Required, validation.Each( // map value
				validation.Required,
				validation.Each( // array
					validation.Required,
					validation.In(PermissionRead, PermissionWrite, PermissionDelete)))),
	}.Filter()
}

// Validate validates GetEdgeKVAccessTokenRequest
func (g GetEdgeKVAccessTokenRequest) Validate() error {
	return validation.Errors{
		"TokenName": validation.Validate(g.TokenName, validation.Required, validation.Length(1, 32)),
	}.Filter()
}

// Validate validates DeleteEdgeKVAccessTokenRequest
func (d DeleteEdgeKVAccessTokenRequest) Validate() error {
	return validation.Errors{
		"TokenName": validation.Validate(d.TokenName, validation.Required, validation.Length(1, 32)),
	}.Filter()
}

var (
	// ErrCreateEdgeKVAccessToken is returned in case an error occurs on CreateEdgeKVAccessToken operation
	ErrCreateEdgeKVAccessToken = errors.New("create an EdgeKV access token")
	// ErrGetEdgeKVAccessToken is returned in case an error occurs on GetEdgeKVAccessToken operation
	ErrGetEdgeKVAccessToken = errors.New("get an EdgeKV access token")
	// ErrListEdgeKVAccessToken is returned in case an error occurs on ListEdgeKVAccessToken operation
	ErrListEdgeKVAccessToken = errors.New("list EdgeKV access tokens")
	// ErrDeleteEdgeKVAccessToken is returned in case an error occurs on DeleteEdgeKVAccessToken operation
	ErrDeleteEdgeKVAccessToken = errors.New("delete an EdgeKV access token")
)

func (e *edgeworkers) CreateEdgeKVAccessToken(ctx context.Context, params CreateEdgeKVAccessTokenRequest) (*CreateEdgeKVAccessTokenResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("CreateEdgeKVAccessToken")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEdgeKVAccessToken, ErrStructValidation, err)
	}

	uri := "/edgekv/v1/tokens"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeKVAccessToken, err)
	}

	var result CreateEdgeKVAccessTokenResponse
	resp, err := e.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeKVAccessToken, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCreateEdgeKVAccessToken, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) GetEdgeKVAccessToken(ctx context.Context, params GetEdgeKVAccessTokenRequest) (*GetEdgeKVAccessTokenResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("GetEdgeKVAccessToken")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeKVAccessToken, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgekv/v1/tokens/%s", params.TokenName))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetEdgeKVAccessToken, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeKVAccessToken, err)
	}

	var result GetEdgeKVAccessTokenResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeKVAccessToken, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeKVAccessToken, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListEdgeKVAccessTokens(ctx context.Context, params ListEdgeKVAccessTokensRequest) (*ListEdgeKVAccessTokensResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListEdgeKVAccessToken")

	uri, err := url.Parse("/edgekv/v1/tokens")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListEdgeKVAccessToken, err)
	}

	q := uri.Query()
	if params.IncludeExpired {
		q.Add("includeExpired", strconv.FormatBool(params.IncludeExpired))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEdgeKVAccessToken, err)
	}

	var result ListEdgeKVAccessTokensResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEdgeKVAccessToken, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEdgeKVAccessToken, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) DeleteEdgeKVAccessToken(ctx context.Context, params DeleteEdgeKVAccessTokenRequest) (*DeleteEdgeKVAccessTokenResponse, error) {
	e.Log(ctx).Debug("DeleteEdgeKVAccessToken")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrDeleteEdgeKVAccessToken, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgekv/v1/tokens/%s", params.TokenName))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrDeleteEdgeKVAccessToken, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeKVAccessToken, err)
	}

	var result DeleteEdgeKVAccessTokenResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeKVAccessToken, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeleteEdgeKVAccessToken, e.Error(resp))
	}

	return &result, nil
}
