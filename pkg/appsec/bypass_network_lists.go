package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BypassNetworkLists interface supports listing or modifying which network lists are
	// used in the bypass network lists settings.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#bypassnetworklist
	BypassNetworkLists interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getbypassnetworklistsforawapconfigversion
		GetBypassNetworkLists(ctx context.Context, params GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklistsforawapconfigversion
		UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklistsforawapconfigversion
		// Note: this method is DEPRECATED and will be removed in a future release.
		RemoveBypassNetworkLists(ctx context.Context, params RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error)
	}

	// GetBypassNetworkListsRequest is used to list which network lists are used in the bypass network lists settings.
	GetBypassNetworkListsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"policyId"`
	}

	// GetBypassNetworkListsResponse is returned from a call to GetBypassNetworkLists.
	GetBypassNetworkListsResponse struct {
		NetworkLists []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"networkLists"`
	}

	// UpdateBypassNetworkListsRequest is used to modify which network lists are used in the bypass network lists settings.
	UpdateBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		PolicyID     string   `json:"policyId"`
		NetworkLists []string `json:"networkLists"`
	}

	// UpdateBypassNetworkListsResponse is returned from a call to UpdateBypassNetworkLists.
	UpdateBypassNetworkListsResponse struct {
		Block       string `json:"block"`
		GeoControls struct {
			BlockedIPNetworkLists struct {
				NetworkList []string `json:"networkList"`
			} `json:"blockedIPNetworkLists"`
		} `json:"geoControls"`
		IPControls struct {
			AllowedIPNetworkLists struct {
				NetworkList []string `json:"networkList"`
			} `json:"allowedIPNetworkLists"`
			BlockedIPNetworkLists struct {
				NetworkList []string `json:"networkList"`
			} `json:"blockedIPNetworkLists"`
		} `json:"ipControls"`
	}

	// RemoveBypassNetworkListsRequest is used to modify which network lists are used in the bypass network lists settings.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	RemoveBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		PolicyID     string   `json:"policyId"`
		NetworkLists []string `json:"networkLists"`
	}

	// RemoveBypassNetworkListsResponse is returned from a call to RemoveBypassNetworkLists.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	RemoveBypassNetworkListsResponse struct {
		NetworkLists []string `json:"networkLists"`
	}
)

// Validate validates a GetBypassNetworkListsRequest.
func (v GetBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateBypassNetworkListsRequest.
func (v UpdateBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a RemoveBypassNetworkListsRequest.
// Note: this method is DEPRECATED and will be removed in a future release.
func (v RemoveBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetBypassNetworkLists(ctx context.Context, params GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debugf("GetBypassNetworkLists(%+v)", params)

	var rval GetBypassNetworkListsResponse

	var uri string
	if params.PolicyID == "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
			params.ConfigID,
			params.Version,
		)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBypassNetworkLists request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debugf("UpdateBypassNetworkLists(%+v)", params)

	var putURL string
	if params.PolicyID == "" {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
			params.ConfigID,
			params.Version,
		)
	} else {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBypassNetworkLists request: %w", err)
	}

	var rval UpdateBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Note: this method is DEPRECATED and will be removed in a future release.
func (p *appsec) RemoveBypassNetworkLists(ctx context.Context, params RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debugf("RemoveBypassNetworkLists(%+v)", params)

	var putURL string
	if params.PolicyID == "" {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
			params.ConfigID,
			params.Version,
		)
	} else {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveBypassNetworkLists request: %w", err)
	}

	var rval RemoveBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
