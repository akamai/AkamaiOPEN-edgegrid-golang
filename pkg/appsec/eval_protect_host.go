package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EvalProtectHost represents a collection of EvalProtectHost
//
// See: EvalProtectHost.GetEvalProtectHost()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// EvalProtectHost  contains operations available on EvalProtectHost  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalprotecthost
	EvalProtectHost interface {
		GetEvalProtectHosts(ctx context.Context, params GetEvalProtectHostsRequest) (*GetEvalProtectHostsResponse, error)
		GetEvalProtectHost(ctx context.Context, params GetEvalProtectHostRequest) (*GetEvalProtectHostResponse, error)
		UpdateEvalProtectHost(ctx context.Context, params UpdateEvalProtectHostRequest) (*UpdateEvalProtectHostResponse, error)
	}

	GetEvalProtectHostRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}
	GetEvalProtectHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	GetEvalProtectHostsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	GetEvalProtectHostsResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	UpdateEvalProtectHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}

	UpdateEvalProtectHostResponse struct {
		HostnameList []struct {
			Hostname string `json:"hostname"`
		} `json:"hostnameList"`
	}

	RemoveEvalProtectHostRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}
	RemoveEvalProtectHostResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetEvalProtectHostRequest
func (v GetEvalProtectHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetEvalProtectHostsRequest
func (v GetEvalProtectHostsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalProtectHostRequest
func (v UpdateEvalProtectHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalProtectHost(ctx context.Context, params GetEvalProtectHostRequest) (*GetEvalProtectHostResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalProtectHost")

	var rval GetEvalProtectHostResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalprotecthost request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalprotecthost  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvalProtectHosts(ctx context.Context, params GetEvalProtectHostsRequest) (*GetEvalProtectHostsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalProtectHosts")

	var rval GetEvalProtectHostsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalprotecthosts request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalprotecthosts request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a EvalProtectHost.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalprotecthost

func (p *appsec) UpdateEvalProtectHost(ctx context.Context, params UpdateEvalProtectHostRequest) (*UpdateEvalProtectHostResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalProtectHost")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/protect-eval-hostnames",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create EvalProtectHostrequest: %w", err)
	}

	var rval UpdateEvalProtectHostResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create EvalProtectHost request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
