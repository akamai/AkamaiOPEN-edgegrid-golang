package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAPSelectedHostnames interface supports retrieving and modifying the list of hostnames protected under
	// a configuration and security policy.
	WAPSelectedHostnames interface {
		GetWAPSelectedHostnames(ctx context.Context, params GetWAPSelectedHostnamesRequest) (*GetWAPSelectedHostnamesResponse, error)
		UpdateWAPSelectedHostnames(ctx context.Context, params UpdateWAPSelectedHostnamesRequest) (*UpdateWAPSelectedHostnamesResponse, error)
	}

	// GetWAPSelectedHostnamesRequest is used to retrieve the WAP selected hostnames and evaluated hostnames.
	GetWAPSelectedHostnamesRequest struct {
		ConfigID         int    `json:"configId"`
		Version          int    `json:"version"`
		SecurityPolicyID string `json:"securityPolicyID"`
	}

	// GetWAPSelectedHostnamesResponse is returned from a call to GetWAPSelectedHostnames.
	GetWAPSelectedHostnamesResponse struct {
		ProtectedHosts []string `json:"protectedHostnames,omitempty"`
		EvaluatedHosts []string `json:"evalHostnames,omitempty"`
	}

	// UpdateWAPSelectedHostnamesRequest is used to modify the WAP selected hostnames and evaluated hostnames.
	UpdateWAPSelectedHostnamesRequest struct {
		ConfigID         int      `json:"configId"`
		Version          int      `json:"version"`
		SecurityPolicyID string   `json:"securityPolicyID"`
		ProtectedHosts   []string `json:"protectedHostnames"`
		EvaluatedHosts   []string `json:"evalHostnames"`
	}

	// UpdateWAPSelectedHostnamesResponse is returned from a call to UpdateWAPSelectedHostnames.
	UpdateWAPSelectedHostnamesResponse struct {
		ProtectedHosts []string `json:"protectedHostnames"`
		EvaluatedHosts []string `json:"evalHostnames"`
	}
)

// Validate validates a GetWAPSelectedHostnamesRequest.
func (v GetWAPSelectedHostnamesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateWAPSelectedHostnamesRequest.
func (v UpdateWAPSelectedHostnamesRequest) Validate() error {
	return validation.Errors{
		"ConfigID":         validation.Validate(v.ConfigID, validation.Required),
		"Version":          validation.Validate(v.Version, validation.Required),
		"SecurityPolicyID": validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAPSelectedHostnames(ctx context.Context, params GetWAPSelectedHostnamesRequest) (*GetWAPSelectedHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetWAPSelectedHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/wap-selected-hostnames",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetWAPSelectedHostnames request: %w", err)
	}

	var result GetWAPSelectedHostnamesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get WAP selected hostnames request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateWAPSelectedHostnames(ctx context.Context, params UpdateWAPSelectedHostnamesRequest) (*UpdateWAPSelectedHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateWAPSelectedHostnames")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/wap-selected-hostnames",
		params.ConfigID,
		params.Version,
		params.SecurityPolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateWAPSelectedHostnames request: %w", err)
	}

	var result UpdateWAPSelectedHostnamesResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update WAP selected hostnames request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
