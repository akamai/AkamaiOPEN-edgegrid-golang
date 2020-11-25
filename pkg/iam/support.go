package iam

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

type (
	// Support is a list of iam supported object methods
	Support interface {
		SupportedCountries(context.Context) ([]string, error)
		SupportedContactTypes(context.Context) ([]string, error)
		SupportedLanguages(context.Context) ([]string, error)
		ListProducts(context.Context) ([]string, error)
		ListTimeoutPolicies(context.Context) ([]TimeoutPolicy, error)
		ListStates(context.Context, ListStatesRequest) ([]string, error)
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
)

func (i *iam) SupportedCountries(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedCountries")

	getURL := path.Join(UserAdminEP, "common", "countries")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "SupportedCountries", err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "SupportedCountries", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "SupportedCountries", i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedContactTypes(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedContactTypes")

	getURL := path.Join(UserAdminEP, "common", "contact-types")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "SupportedContactTypes", err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "SupportedContactTypes", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "SupportedContactTypes", i.Error(resp))
	}

	return rval, nil
}

func (i *iam) SupportedLanguages(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("SupportedLanguages")

	getURL := path.Join(UserAdminEP, "common", "supported-languages")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "SupportedLanguages", err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "SupportedLanguages", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "SupportedLanguages", i.Error(resp))
	}

	return rval, nil
}

func (i *iam) ListProducts(ctx context.Context) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListProducts")

	getURL := path.Join(UserAdminEP, "common", "notification-products")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "ListProducts", err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "ListProducts", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "ListProducts", i.Error(resp))
	}

	return rval, nil
}

func (i *iam) ListTimeoutPolicies(ctx context.Context) ([]TimeoutPolicy, error) {
	logger := i.Log(ctx)
	logger.Debug("ListTimeoutPolicies")

	getURL := path.Join(UserAdminEP, "common", "timeout-policies")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "ListTimeoutPolicies", err)
	}

	var rval []TimeoutPolicy
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "ListTimeoutPolicies", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "ListTimeoutPolicies", i.Error(resp))
	}

	return rval, nil
}

func (i *iam) ListStates(ctx context.Context, params ListStatesRequest) ([]string, error) {
	logger := i.Log(ctx)
	logger.Debug("ListStates")

	getURL := path.Join(UserAdminEP, "common", "countries", params.Country, "states")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %s", "ListStates", err)
	}

	var rval []string
	resp, err := i.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %s", "ListStates", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "ListStates", i.Error(resp))
	}

	return rval, nil
}
