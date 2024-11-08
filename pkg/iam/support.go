package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPasswordPolicyResponse holds the response data from the GetPasswordPolicy endpoint.
	GetPasswordPolicyResponse struct {
		CaseDiff        int64  `json:"caseDif"`
		MaxRepeating    int64  `json:"maxRepeating"`
		MinDigits       int64  `json:"minDigits"`
		MinLength       int64  `json:"minLength"`
		MinLetters      int64  `json:"minLetters"`
		MinNonAlpha     int64  `json:"minNonAlpha"`
		MinReuse        int64  `json:"minReuse"`
		PwClass         string `json:"pwclass"`
		RotateFrequency int64  `json:"rotateFrequency"`
	}

	// TimeoutPolicy encapsulates the response from the ListTimeoutPolicies endpoint.
	TimeoutPolicy struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	// ListStatesRequest contains the country request parameter for the ListStates endpoint.
	ListStatesRequest struct {
		Country string
	}

	// ListAccountSwitchKeysRequest contains the request parameters for the ListAccountSwitchKeys endpoint.
	ListAccountSwitchKeysRequest struct {
		ClientID string
		Search   string
	}

	// ListAccountSwitchKeysResponse holds the response data from the ListAccountSwitchKeys endpoint.
	ListAccountSwitchKeysResponse []AccountSwitchKey

	// AccountSwitchKey contains information about account switch key.
	AccountSwitchKey struct {
		AccountName      string `json:"accountName"`
		AccountSwitchKey string `json:"accountSwitchKey"`
	}

	// Timezone contains the response from the ListSupportedTimezones endpoint.
	Timezone struct {
		Description string `json:"description"`
		Offset      string `json:"offset"`
		Posix       string `json:"posix"`
		Timezone    string `json:"timezone"`
	}
)

var (
	// ErrGetPasswordPolicy is returned when GetPasswordPolicy fails.
	ErrGetPasswordPolicy = errors.New("get password policy")

	// ErrListProducts is returned when ListProducts fails.
	ErrListProducts = errors.New("list products")

	// ErrListStates is returned when ListStates fails.
	ErrListStates = errors.New("list states")

	// ErrListTimeoutPolicies is returned when ListTimeoutPolicies fails.
	ErrListTimeoutPolicies = errors.New("list timeout policies")

	// ErrListAccountSwitchKeys is returned when ListAccountSwitchKeys fails.
	ErrListAccountSwitchKeys = errors.New("list account switch keys")

	// ErrSupportedContactTypes is returned when SupportedContactTypes fails.
	ErrSupportedContactTypes = errors.New("supported contact types")

	// ErrSupportedCountries is returned when SupportedCountries fails.
	ErrSupportedCountries = errors.New("supported countries")

	// ErrSupportedLanguages is returned when SupportedLanguages fails.
	ErrSupportedLanguages = errors.New("supported languages")

	// ErrSupportedTimezones is returned when SupportedTimezones fails.
	ErrSupportedTimezones = errors.New("supported timezones")
)

// Validate validates ListStatesRequest.
func (r ListStatesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Country": validation.Validate(r.Country, validation.Required),
	})
}

func (i *iam) GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetPasswordPolicy")

	uri := "/identity-management/v3/user-admin/common/password-policy"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPasswordPolicy, err)
	}

	var result GetPasswordPolicyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPasswordPolicy, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPasswordPolicy, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListProducts(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListProducts")

	uri := "/identity-management/v3/user-admin/common/notification-products"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListProducts, err)
	}

	var result []string
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListProducts, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListProducts, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListStates(ctx context.Context, params ListStatesRequest) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListStates")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListStates, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/identity-management/v3/user-admin/common/countries/%s/states", params.Country)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListStates, err)
	}

	var result []string
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListStates, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListStates, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListTimeoutPolicies(ctx context.Context) ([]TimeoutPolicy, error) {
	logger := i.Log(ctx)
	logger.Debug("ListTimeoutPolicies")

	uri := "/identity-management/v3/user-admin/common/timeout-policies"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListTimeoutPolicies, err)
	}

	var result []TimeoutPolicy
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListTimeoutPolicies, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListTimeoutPolicies, i.Error(resp))
	}

	return result, nil
}

func (i *iam) ListAccountSwitchKeys(ctx context.Context, params ListAccountSwitchKeysRequest) (ListAccountSwitchKeysResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListAccountSwitchKeys")

	if params.ClientID == "" {
		params.ClientID = "self"
	}

	uri, err := url.Parse(fmt.Sprintf("/identity-management/v3/api-clients/%s/account-switch-keys", params.ClientID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListAccountSwitchKeys, err)
	}

	if params.Search != "" {
		q := uri.Query()
		q.Add("search", params.Search)
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAccountSwitchKeys, err)
	}

	var result ListAccountSwitchKeysResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAccountSwitchKeys, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAccountSwitchKeys, i.Error(resp))
	}

	return result, nil
}

func (i *iam) SupportedContactTypes(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedContactTypes")

	uri := "/identity-management/v3/user-admin/common/contact-types"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedContactTypes, err)
	}

	var result []string
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedContactTypes, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedContactTypes, i.Error(resp))
	}

	return result, nil
}

func (i *iam) SupportedCountries(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedCountries")

	uri := "/identity-management/v3/user-admin/common/countries"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedCountries, err)
	}

	var result []string
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedCountries, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedCountries, i.Error(resp))
	}

	return result, nil
}

func (i *iam) SupportedLanguages(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedLanguages")

	uri := "/identity-management/v3/user-admin/common/supported-languages"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedLanguages, err)
	}

	var result []string
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedLanguages, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedLanguages, i.Error(resp))
	}

	return result, nil
}

func (i *iam) SupportedTimezones(ctx context.Context) ([]Timezone, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedTimezones")

	uri := "/identity-management/v3/user-admin/common/timezones"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedTimezones, err)
	}

	var result []Timezone
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedTimezones, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedTimezones, i.Error(resp))
	}

	return result, nil
}
