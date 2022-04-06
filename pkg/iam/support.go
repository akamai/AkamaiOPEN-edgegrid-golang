package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

type (
	// Support is a list of iam supported object methods
	Support interface {
		// ListProducts returns all products a user can subscribe to and receive notifications for on the account
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-notification-products
		ListProducts(context.Context) ([]string, error)

		// ListStates lists U.S. states or Canadian provinces
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-states
		ListStates(context.Context, ListStatesRequest) ([]string, error)

		// ListTimeoutPolicies lists all the possible session timeout policies that Akamai supports
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-timeout-policies
		ListTimeoutPolicies(context.Context) ([]TimeoutPolicy, error)

		// SupportedContactTypes lists all the possible contact types that Akamai supports
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-contact-types
		SupportedContactTypes(context.Context) ([]string, error)

		// SupportedCountries returns all the possible countries that Akamai supports
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-countries
		SupportedCountries(context.Context) ([]string, error)

		// SupportedLanguages lists all the possible languages Akamai supports
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-languages
		SupportedLanguages(context.Context) ([]string, error)

		// SupportedTimezones lists all time zones Akamai supports
		//
		// See: https://techdocs.akamai.com/iam-user-admin/reference/get-common-timezones
		SupportedTimezones(context.Context) ([]Timezone, error)
	}

	// TimeoutPolicy specifies session timeout policy options that can be assigned to each user
	TimeoutPolicy struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	// ListStatesRequest specifies the country for the requested states
	ListStatesRequest struct {
		Country string `json:"country"`
	}

	// Timezone is the object retured by the SupportedTimezones method
	Timezone struct {
		Description string `json:"description"`
		Offset      string `json:"offset"`
		Posix       string `json:"posix"`
		Timezone    string `json:"timezone"`
	}
)

var (
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

func (i *iam) ListProducts(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListProducts")

	getURL := path.Join(UserAdminEP, "common", "notification-products")

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

	getURL := path.Join(UserAdminEP, "common", "countries", params.Country, "states")

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

	getURL := path.Join(UserAdminEP, "common", "timeout-policies")

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

	getURL := path.Join(UserAdminEP, "common", "contact-types")

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

	getURL := path.Join(UserAdminEP, "common", "countries")

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

	getURL := path.Join(UserAdminEP, "common", "supported-languages")

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

	getURL := path.Join(UserAdminEP, "common", "timezones")

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
