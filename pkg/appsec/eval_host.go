package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EvalHost represents a collection of EvalHost
//
// See: EvalHost.GetEvalHost()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// EvalHost  contains operations available on EvalHost  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalhost
	EvalHost interface {
		GetEvalHosts(ctx context.Context, params GetEvalHostsRequest) (*GetEvalHostsResponse, error)
		GetEvalHost(ctx context.Context, params GetEvalHostRequest) (*GetEvalHostResponse, error)
		UpdateEvalHost(ctx context.Context, params UpdateEvalHostRequest) (*UpdateEvalHostResponse, error)
	}

	GetEvalHostRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	GetEvalHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	GetEvalHostsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	GetEvalHostsResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	UpdateEvalHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}
	
	UpdateEvalHostResponse struct {
		HostnameList []struct {
			Hostname string `json:"hostname"`
		} `json:"hostnameList"`
	}

	RemoveEvalHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}

	RemoveEvalHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}
)

// Validate validates GetEvalHostRequest
func (v GetEvalHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetEvalHostsRequest
func (v GetEvalHostsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalHostRequest
func (v UpdateEvalHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalHost(ctx context.Context, params GetEvalHostRequest) (*GetEvalHostResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalHost")

	var rval GetEvalHostResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalhost request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalhost  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvalHosts(ctx context.Context, params GetEvalHostsRequest) (*GetEvalHostsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalHosts")

	var rval GetEvalHostsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalhosts request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalhosts request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a EvalHost.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalhost

func (p *appsec) UpdateEvalHost(ctx context.Context, params UpdateEvalHostRequest) (*UpdateEvalHostResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalHost")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create EvalHostrequest: %w", err)
	}

	var rval UpdateEvalHostResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create EvalHost request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
