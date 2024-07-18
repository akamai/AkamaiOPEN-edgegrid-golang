package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Support is a list of IAM supported objects API interfaces
	Support interface {
		// GetPasswordPolicy gets the password policy for the account.
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-password-policy
		GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResponse, error)

		// ListProducts lists products a user can subscribe to and receive notifications for on the account
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-notification-products
		ListProducts(context.Context) ([]string, error)

		// ListStates lists U.S. states or Canadian provinces
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-states
		ListStates(context.Context, ListStatesRequest) ([]string, error)

		// ListTimeoutPolicies lists all the possible session timeout policies
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-timeout-policies
		ListTimeoutPolicies(context.Context) ([]TimeoutPolicy, error)

		// SupportedContactTypes lists supported contact types
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-contact-types
		SupportedContactTypes(context.Context) ([]string, error)

		// SupportedCountries lists supported countries
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-countries
		SupportedCountries(context.Context) ([]string, error)

		// SupportedLanguages lists supported languages
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-languages
		SupportedLanguages(context.Context) ([]string, error)

		// SupportedTimezones lists supported timezones
		//
		// See: https://techdocs.akamai.com/iam-api/reference/get-common-timezones
		SupportedTimezones(context.Context) ([]Timezone, error)
	}

	// GetPasswordPolicyResponse holds the response data from GetPasswordPolicy.
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

	// TimeoutPolicy encapsulates the response of the list timeout policies endpoint
	TimeoutPolicy struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	// ListStatesRequest contains the country request parameter for the list states endpoint
	ListStatesRequest struct {
		Country string
	}

	// Timezone contains the response of the list supported timezones endpoint
	Timezone struct {
		Description string `json:"description"`
		Offset      string `json:"offset"`
		Posix       string `json:"posix"`
		Timezone    string `json:"timezone"`
	}
)

var (
	// ErrGetPasswordPolicy is returned when GetPasswordPolicy fails
	ErrGetPasswordPolicy = errors.New("get password policy")

	// ErrListProducts is returned when ListProducts fails
	ErrListProducts = errors.New("list products")

	// ErrListStates is returned when ListStates fails
	ErrListStates = errors.New("list states")

	// ErrListTimeoutPolicies is returned when ListTimeoutPolicies fails
	ErrListTimeoutPolicies = errors.New("list timeout policies")

	// ErrSupportedContactTypes is returned when SupportedContactTypes fails
	ErrSupportedContactTypes = errors.New("supported contact types")

	// ErrSupportedCountries is returned when SupportedCountries fails
	ErrSupportedCountries = errors.New("supported countries")

	// ErrSupportedLanguages is returned when SupportedLanguages fails
	ErrSupportedLanguages = errors.New("supported languages")

	// ErrSupportedTimezones is returned when SupportedTimezones fails
	ErrSupportedTimezones = errors.New("supported timezones")
)

// Validate validates ListStatesRequest
func (r ListStatesRequest) Validate() error {
	return validation.Errors{
		"Country": validation.Validate(r.Country, validation.Required),
	}.Filter()
}

func (i *iam) GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetPasswordPolicy")

	getURL := "/identity-management/v3/user-admin/common/password-policy"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPasswordPolicy, err)
	}

	var result GetPasswordPolicyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPasswordPolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPasswordPolicy, i.Error(resp))
	}

	return &result, nil
}

func (i *iam) ListProducts(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListProducts")

	getURL := "/identity-management/v3/user-admin/common/notification-products"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListProducts, err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListProducts, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListProducts, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) ListStates(ctx context.Context, params ListStatesRequest) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListStates")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListStates, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/identity-management/v3/user-admin/common/countries/%s/states", params.Country)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListStates, err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListStates, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListStates, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) ListTimeoutPolicies(ctx context.Context) ([]TimeoutPolicy, error) {
	logger := i.Log(ctx)
	logger.Debug("ListTimeoutPolicies")

	getURL := "/identity-management/v3/user-admin/common/timeout-policies"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListTimeoutPolicies, err)
	}

	var rval []TimeoutPolicy
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListTimeoutPolicies, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListTimeoutPolicies, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedContactTypes(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedContactTypes")

	getURL := "/identity-management/v3/user-admin/common/contact-types"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedContactTypes, err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedContactTypes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedContactTypes, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedCountries(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedCountries")

	getURL := "/identity-management/v3/user-admin/common/countries"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedCountries, err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedCountries, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedCountries, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedLanguages(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedLanguages")

	getURL := "/identity-management/v3/user-admin/common/supported-languages"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedLanguages, err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedLanguages, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedLanguages, i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedTimezones(ctx context.Context) ([]Timezone, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedTimezones")

	getURL := "/identity-management/v3/user-admin/common/timezones"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSupportedTimezones, err)
	}

	var rval []Timezone
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSupportedTimezones, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSupportedTimezones, i.Error(resp))
	}

	return rval, nil
}
