package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// WAFMode represents a collection of WAFMode
//
// See: WAFMode.GetWAFMode()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// WAFMode  contains operations available on WAFMode  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getwafmode
	WAFMode interface {
		GetWAFModes(ctx context.Context, params GetWAFModesRequest) (*GetWAFModesResponse, error)
		GetWAFMode(ctx context.Context, params GetWAFModeRequest) (*GetWAFModeResponse, error)
		UpdateWAFMode(ctx context.Context, params UpdateWAFModeRequest) (*UpdateWAFModeResponse, error)
	}

	GetWAFModesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	GetWAFModeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	GetWAFModesResponse struct {
		Current string `json:"current"`
		Mode    string `json:"mode"`
		Eval    string `json:"eval"`
	}

	GetWAFModeResponse struct {
		Current string `json:"current"`
		Mode    string `json:"mode"`
		Eval    string `json:"eval"`
	}

	UpdateWAFModeRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"-"`
		Mode     string `json:"mode"`
		Eval     string `json:"-"`
	}

	UpdateWAFModeResponse struct {
		Current string `json:"current"`
		Mode    string `json:"mode"`
	}
)

// Validate validates GetWAFModeRequest
func (v GetWAFModeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetWAFModesRequest
func (v GetWAFModesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateWAFModeRequest
func (v UpdateWAFModeRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAFMode(ctx context.Context, params GetWAFModeRequest) (*GetWAFModeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFMode")

	var rval GetWAFModeResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafmode request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafmode  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetWAFModes(ctx context.Context, params GetWAFModesRequest) (*GetWAFModesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFModes")

	var rval GetWAFModesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafmodes request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafmodes request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a WAFMode.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putwafmode

func (p *appsec) UpdateWAFMode(ctx context.Context, params UpdateWAFModeRequest) (*UpdateWAFModeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateWAFMode")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/mode",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create WAFModerequest: %w", err)
	}

	var rval UpdateWAFModeResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create WAFMode request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
