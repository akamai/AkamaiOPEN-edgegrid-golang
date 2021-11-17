package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The EvalHost interface supports retrieving and modifying list of evaluation hostnames for a configuration.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#evalhostname
	EvalHost interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevaluationhostnames
		GetEvalHosts(ctx context.Context, params GetEvalHostsRequest) (*GetEvalHostsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevaluationhostnames
		// Deprecated: this method will be removed in a future release. Use GetEvalHosts instead.
		GetEvalHost(ctx context.Context, params GetEvalHostRequest) (*GetEvalHostResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevaluationhostnames
		UpdateEvalHost(ctx context.Context, params UpdateEvalHostRequest) (*UpdateEvalHostResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevaluationhostnames
		RemoveEvalHost(ctx context.Context, params RemoveEvalHostRequest) (*RemoveEvalHostResponse, error)
	}

	// GetEvalHostRequest is used to retrieve the evaluation hostnames for a configuration.
	// Deprecated: this struct will be removed in a future release.
	GetEvalHostRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetEvalHostResponse is returned from a call to GetEvalHost.
	// Deprecated: this struct will be removed in a future release.
	GetEvalHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	// GetEvalHostsRequest is used to retrieve the evaluation hostnames for a configuration.
	GetEvalHostsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetEvalHostsResponse is returned from a call to GetEvalHosts.
	GetEvalHostsResponse struct {
		Hostnames []string `json:"hostnames"`
	}

	// UpdateEvalHostRequest is used to modify the evaluation hostnames for a configuration.
	UpdateEvalHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}

	// UpdateEvalHostResponse is returned from a call to UpdateEvalHost.
	UpdateEvalHostResponse struct {
		HostnameList []struct {
			Hostname string `json:"hostname"`
		} `json:"hostnameList"`
	}

	// RemoveEvalHostRequest is used to remove the evaluation hostnames for a configuration.
	RemoveEvalHostRequest struct {
		ConfigID  int      `json:"-"`
		Version   int      `json:"-"`
		Hostnames []string `json:"hostnames"`
	}

	// RemoveEvalHostResponse is returned from a call to RemoveEvalHost.
	RemoveEvalHostResponse struct {
		Hostnames []string `json:"hostnames"`
	}
)

// Validate validates a GetEvalHostRequest.
// Deprecated: this method will be removed in a future release.
func (v GetEvalHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a GetEvalHostsRequest.
func (v GetEvalHostsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateEvalHostRequest.
func (v UpdateEvalHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a RemoveEvalHostRequest.
func (v RemoveEvalHostRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Deprecated: this method will be removed in a future release.
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
		return nil, fmt.Errorf("failed to create GetEvalHost request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalHost request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create GetEvalHosts request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetEvalHosts request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateEvalHost request: %w", err)
	}

	var rval UpdateEvalHostResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEvalHost request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *appsec) RemoveEvalHost(ctx context.Context, params RemoveEvalHostRequest) (*RemoveEvalHostResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveEvalHost")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveEvalHost request: %w", err)
	}

	var rval RemoveEvalHostResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveEvalHost request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
