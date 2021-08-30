package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The EvalProtectHost interface supports retrieving the evaluation hostnames for a configuration and
	// moving hostnames from evaluating to protected status.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#evalhostname
	EvalProtectHost interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevaluationhostnames
		GetEvalProtectHosts(ctx context.Context, params GetEvalProtectHostsRequest) (*GetEvalProtectHostsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevaluationhostnames
		// Note: this method is DEPRECATED and will be removed in a future release.
		GetEvalProtectHost(ctx context.Context, params GetEvalProtectHostRequest) (*GetEvalProtectHostResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmoveevaluationhostnamestoprotection
		UpdateEvalProtectHost(ctx context.Context, params UpdateEvalProtectHostRequest) (*UpdateEvalProtectHostResponse, error)
	}

	// GetEvalProtectHostRequest is used to call GetEvalProtectHost.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetEvalProtectHostRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetEvalProtectHostResponse is returned from a call to GetEvalProtectHost.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetEvalProtectHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	// GetEvalProtectHostsRequest is used to call GetEvalProtectHosts.
	GetEvalProtectHostsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetEvalProtectHostsResponse is returned from a call to GetEvalProtectHosts.
	GetEvalProtectHostsResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	// UpdateEvalProtectHostRequest is used to call UpdateEvalProtectHost.
	UpdateEvalProtectHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}

	// UpdateEvalProtectHostResponse is returned from a call to UpdateEvalProtectHost.
	UpdateEvalProtectHostResponse struct {
		HostnameList []struct {
			Hostname string `json:"hostname"`
		} `json:"hostnameList"`
	}
)

// Validate validates a GetEvalProtectHostRequest.
// Note: this method is DEPRECATED and will be removed in a future release.
func (v GetEvalProtectHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a GetEvalProtectHostsRequest.
func (v GetEvalProtectHostsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateEvalProtectHostRequest.
func (v UpdateEvalProtectHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Note: this method is DEPRECATED and will be removed in a future release.
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
		return nil, fmt.Errorf("failed to create GetEvalProtectHost request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalProtectHost request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetEvalProtectHosts request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalProtectHosts request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateEvalProtectHost request: %w", err)
	}

	var rval UpdateEvalProtectHostResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEvalProtectHost request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
