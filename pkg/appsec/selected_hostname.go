package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The SelectedHostname interface supports retrieving and modifying the list of hostnames protected under
	// a configuration.
	// Deprecated: this interface will be removed in a future release. Use the WAPSelectedHostnames interface instead.
	SelectedHostname interface {
		// GetSelectedHostnames lists the hostnames that the configuration version selects as candidates of protected hostnames,
		// which you can use in match targets.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-selected-hostnames
		// Deprecated: this method will be removed in a future release. Use the GetWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		GetSelectedHostnames(ctx context.Context, params GetSelectedHostnamesRequest) (*GetSelectedHostnamesResponse, error)

		// GetSelectedHostname returns the hostname that the configuration version selects as a candidate of protected hostname,
		// which you can use in match targets.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-selected-hostnames
		// Deprecated: this method will be removed in a future release. Use the GetWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		GetSelectedHostname(ctx context.Context, params GetSelectedHostnameRequest) (*GetSelectedHostnameResponse, error)

		// UpdateSelectedHostname updates the selected hostname for a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-selected-hostnames
		// Deprecated: this method will be removed in a future release. Use the UpdateWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		UpdateSelectedHostname(ctx context.Context, params UpdateSelectedHostnameRequest) (*UpdateSelectedHostnameResponse, error)

		// UpdateSelectedHostnames updates the list of selected hostnames for a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-selected-hostnames
		// Deprecated: this method will be removed in a future release. Use the UpdateWAPSelectedHostnames method of the WAPSelectedHostnames interface instead.
		UpdateSelectedHostnames(ctx context.Context, params UpdateSelectedHostnamesRequest) (*UpdateSelectedHostnamesResponse, error)
	}

	// GetSelectedHostnamesRequest is used to retrieve the selected hostnames for a configuration.
	GetSelectedHostnamesRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	// GetSelectedHostnamesResponse is returned from a call to GetSelectedHostnames.
	GetSelectedHostnamesResponse struct {
		HostnameList []Hostname `json:"hostnameList,omitempty"`
	}

	// GetSelectedHostnameRequest is used to retrieve the selected hostnames for a configuration.
	// Deprecated: this struct will be removed in a future release.
	GetSelectedHostnameRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	// GetSelectedHostnameResponse is returned from a call to GetSelectedHostname.
	// Deprecated: this struct will be removed in a future release.
	GetSelectedHostnameResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}

	// UpdateSelectedHostnamesRequest is used to modify the selected hostnames for a configuration.
	UpdateSelectedHostnamesRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	// UpdateSelectedHostnamesResponse is returned from a call to UpdateSelectedHostnames.
	UpdateSelectedHostnamesResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}

	// UpdateSelectedHostnameRequest is used to modify the selected hostnames for a configuration.
	// Deprecated: this struct will be removed in a future release.
	UpdateSelectedHostnameRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	// UpdateSelectedHostnameResponse is returned from a call to UpdateSelectedHostname.
	// Deprecated: this struct will be removed in a future release.
	UpdateSelectedHostnameResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}

	// Hostname describes a hostname that may be protected.
	Hostname struct {
		Hostname string `json:"hostname"`
	}
)

// Validate validates a GetSelectedHostnameRequest.
// Deprecated: this method will be removed in a future release.
func (v GetSelectedHostnameRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a GetSelectedHostnamesRequest.
func (v GetSelectedHostnamesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSelectedHostnamesRequest.
func (v UpdateSelectedHostnamesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateSelectedHostnameRequest.
// Deprecated: this method will be removed in a future release.
func (v UpdateSelectedHostnameRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) GetSelectedHostname(ctx context.Context, params GetSelectedHostnameRequest) (*GetSelectedHostnameResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSelectedHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSelectedHostname request: %w", err)
	}

	var result GetSelectedHostnameResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get selected hostname request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetSelectedHostnames(ctx context.Context, params GetSelectedHostnamesRequest) (*GetSelectedHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSelectedHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSelectedHostnames request: %w", err)
	}

	var result GetSelectedHostnamesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get selected hostnames request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateSelectedHostnames(ctx context.Context, params UpdateSelectedHostnamesRequest) (*UpdateSelectedHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSelectedHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateSelectedHostnames request: %w", err)
	}

	var result UpdateSelectedHostnamesResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update selected hostnames request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) UpdateSelectedHostname(ctx context.Context, params UpdateSelectedHostnameRequest) (*UpdateSelectedHostnameResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateSelectedHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateSelectedHostname request: %w", err)
	}

	var result UpdateSelectedHostnameResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update selected hostname request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
