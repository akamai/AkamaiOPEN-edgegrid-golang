package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAPBypassNetworkLists interface supports listing or modifying which network lists are
	// used in the bypass network lists settings.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#bypassnetworklist
	WAPBypassNetworkLists interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getbypassnetworklistsforawapconfigversion
		GetWAPBypassNetworkLists(ctx context.Context, params GetWAPBypassNetworkListsRequest) (*GetWAPBypassNetworkListsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklistsforawapconfigversion
		UpdateWAPBypassNetworkLists(ctx context.Context, params UpdateWAPBypassNetworkListsRequest) (*UpdateWAPBypassNetworkListsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklistsforawapconfigversion
		RemoveWAPBypassNetworkLists(ctx context.Context, params RemoveWAPBypassNetworkListsRequest) (*RemoveWAPBypassNetworkListsResponse, error)
	}

	// GetWAPBypassNetworkListsRequest is used to list which network lists are used in the bypass network lists settings.
	GetWAPBypassNetworkListsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"policyId"`
	}

	// NetworkListsStruct is used to define a network list.
	NetworkListsStruct struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	// GetWAPBypassNetworkListsResponse is returned from a call to GetWAPBypassNetworkLists.
	GetWAPBypassNetworkListsResponse struct {
		NetworkLists []NetworkListsStruct `json:"networkLists"`
	}

	// UpdateWAPBypassNetworkListsRequest is used to modify which network lists are used in the bypass network lists settings.
	UpdateWAPBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		PolicyID     string   `json:"policyId"`
		NetworkLists []string `json:"networkLists"`
	}

	// IPNetworkListsList is used to define a list of IP network lists.
	IPNetworkListsList struct {
		NetworkList []string `json:"networkList"`
	}

	// GeoControlsList is used to define a list of blocked IP network lists.
	GeoControlsList struct {
		BlockedIPNetworkLists IPNetworkListsList `json:"networkList"`
	}

	// IPControlsLists is used to define lists of allowed and blocked IP network lists.
	IPControlsLists struct {
		AllowedIPNetworkLists IPNetworkListsList `json:"allowedIPNetworkLists"`
		BlockedIPNetworkLists IPNetworkListsList `json:"blockedIPNetworkLists"`
	}

	// UpdateWAPBypassNetworkListsResponse is returned from a call to UpdateWAPBypassNetworkLists.
	UpdateWAPBypassNetworkListsResponse struct {
		Block       string          `json:"block"`
		GeoControls GeoControlsList `json:"geoControls"`
		IPControls  IPControlsLists `json:"ipControls"`
	}

	// RemoveWAPBypassNetworkListsRequest is used to modify which network lists are used in the bypass network lists settings.
	// Deprecated: this struct will be removed in a future release.
	RemoveWAPBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		PolicyID     string   `json:"policyId"`
		NetworkLists []string `json:"networkLists"`
	}

	// RemoveWAPBypassNetworkListsResponse is returned from a call to RemoveWAPBypassNetworkLists.
	// Deprecated: this struct will be removed in a future release.
	RemoveWAPBypassNetworkListsResponse struct {
		NetworkLists []string `json:"networkLists"`
	}
)

// Validate validates a GetWAPBypassNetworkListsRequest.
func (v GetWAPBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateWAPBypassNetworkListsRequest.
func (v UpdateWAPBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveWAPBypassNetworkListsRequest.
// Deprecated: this method will be removed in a future release.
func (v RemoveWAPBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAPBypassNetworkLists(ctx context.Context, params GetWAPBypassNetworkListsRequest) (*GetWAPBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("GetWAPBypassNetworkLists(%+v)", params)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval GetWAPBypassNetworkListsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetWAPBypassNetworkLists request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetWAPBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateWAPBypassNetworkLists(ctx context.Context, params UpdateWAPBypassNetworkListsRequest) (*UpdateWAPBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("UpdateWAPBypassNetworkLists(%+v)", params)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateWAPBypassNetworkLists request: %w", err)
	}

	var rval UpdateWAPBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateWAPBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) RemoveWAPBypassNetworkLists(ctx context.Context, params RemoveWAPBypassNetworkListsRequest) (*RemoveWAPBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("RemoveWAPBypassNetworkLists(%+v)", params)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveWAPBypassNetworkLists request: %w", err)
	}

	var rval RemoveWAPBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("RemoveWAPBypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
