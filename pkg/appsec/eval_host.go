package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The EvalHost interface supports retrieving and modifying list of evaluation hostnames for a configuration.
	// Deprecated: this interface will be removed in a future release. Use the WAPSelectedHostnames interface instead.
	EvalHost interface {
		// GetEvalHosts lists the evaluation hostnames for a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-selected-hostnames-eval-hostnames
		// Deprecated: this method will be removed in a future release. Use the GetWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		GetEvalHosts(ctx context.Context, params GetEvalHostsRequest) (*GetEvalHostsResponse, error)

		// GetEvalHost return the specified evaluation hostname for a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-selected-hostnames-eval-hostnames
		// Deprecated: this method will be removed in a future release. Use the GetWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		GetEvalHost(ctx context.Context, params GetEvalHostRequest) (*GetEvalHostResponse, error)

		// UpdateEvalHost updates the list of hostnames you want to evaluate for a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-selected-eval-hostnames
		// Deprecated: this method will be removed in a future release. Use the UpdateWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		UpdateEvalHost(ctx context.Context, params UpdateEvalHostRequest) (*UpdateEvalHostResponse, error)

		// RemoveEvalHost removed the specified evaluation hostname.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-selected-eval-hostnames
		// Deprecated: this method will be removed in a future release. Use the WAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
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
	logger := p.Log(ctx)
	logger.Debug("GetEvalHost")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalHost request: %w", err)
	}

	var result GetEvalHostResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval host request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetEvalHosts(ctx context.Context, params GetEvalHostsRequest) (*GetEvalHostsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalHosts")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalHosts request: %w", err)
	}

	var result GetEvalHostsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get eval hosts request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateEvalHost(ctx context.Context, params UpdateEvalHostRequest) (*UpdateEvalHostResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEvalHost")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalHost request: %w", err)
	}

	var result UpdateEvalHostResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update eval host request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveEvalHost(ctx context.Context, params RemoveEvalHostRequest) (*RemoveEvalHostResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveEvalHost")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames/eval-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveEvalHost request: %w", err)
	}

	var result RemoveEvalHostResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove eval host request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
