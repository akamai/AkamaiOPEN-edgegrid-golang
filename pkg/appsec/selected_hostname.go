package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SelectedHostname represents a collection of SelectedHostname
//
// See: SelectedHostname.GetSelectedHostname()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// SelectedHostname  contains operations available on SelectedHostname  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getselectedhostname
	SelectedHostname interface {
		GetSelectedHostnames(ctx context.Context, params GetSelectedHostnamesRequest) (*GetSelectedHostnamesResponse, error)
		GetSelectedHostname(ctx context.Context, params GetSelectedHostnameRequest) (*GetSelectedHostnameResponse, error)
		UpdateSelectedHostname(ctx context.Context, params UpdateSelectedHostnameRequest) (*UpdateSelectedHostnameResponse, error)
	}

	GetSelectedHostnamesRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	GetSelectedHostnameRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	GetSelectedHostnamesResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}

	UpdateSelectedHostnameRequest struct {
		ConfigID     int        `json:"configId"`
		Version      int        `json:"version"`
		HostnameList []Hostname `json:"hostnameList"`
	}

	UpdateSelectedHostnameResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}

	Hostname struct {
		Hostname string `json:"hostname"`
	}

	GetSelectedHostnameResponse struct {
		HostnameList []Hostname `json:"hostnameList"`
	}
)

// Validate validates GetSelectedHostnameRequest
func (v GetSelectedHostnameRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetSelectedHostnamesRequest
func (v GetSelectedHostnamesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateSelectedHostnameRequest
func (v UpdateSelectedHostnameRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetSelectedHostname(ctx context.Context, params GetSelectedHostnameRequest) (*GetSelectedHostnameResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSelectedHostname")

	var rval GetSelectedHostnameResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getselectedhostname request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetSelectedHostnames(ctx context.Context, params GetSelectedHostnamesRequest) (*GetSelectedHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetSelectedHostnames")

	var rval GetSelectedHostnamesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getselectedhostnames request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getselectedhostnames request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a SelectedHostname.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putselectedhostname

func (p *appsec) UpdateSelectedHostname(ctx context.Context, params UpdateSelectedHostnameRequest) (*UpdateSelectedHostnameResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateSelectedHostname")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/selected-hostnames",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create SelectedHostnamerequest: %w", err)
	}

	var rval UpdateSelectedHostnameResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create SelectedHostname request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
