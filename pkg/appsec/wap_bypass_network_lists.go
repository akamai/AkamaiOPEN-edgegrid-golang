package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The WAPBypassNetworkLists interface supports listing or modifying which network lists are
	// used in the bypass network lists settings.
	WAPBypassNetworkLists interface {
		// GetWAPBypassNetworkLists lists which network lists are used in the bypass network lists settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-bypass-network-lists
		GetWAPBypassNetworkLists(ctx context.Context, params GetWAPBypassNetworkListsRequest) (*GetWAPBypassNetworkListsResponse, error)

		// UpdateWAPBypassNetworkLists updates which network lists to use in the bypass network lists settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-bypass-network-lists
		UpdateWAPBypassNetworkLists(ctx context.Context, params UpdateWAPBypassNetworkListsRequest) (*UpdateWAPBypassNetworkListsResponse, error)

		// RemoveWAPBypassNetworkLists removes network lists in the bypass network lists settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-bypass-network-lists
		RemoveWAPBypassNetworkLists(ctx context.Context, params RemoveWAPBypassNetworkListsRequest) (*RemoveWAPBypassNetworkListsResponse, error)
	}

	// GetWAPBypassNetworkListsRequest is used to list which network lists are used in the bypass network lists settings.
	GetWAPBypassNetworkListsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"policyId"`
	}

	// NetworkList is used to define a network list.
	NetworkList struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	// GetWAPBypassNetworkListsResponse is returned from a call to GetWAPBypassNetworkLists.
	GetWAPBypassNetworkListsResponse struct {
		NetworkLists []NetworkList `json:"networkLists"`
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
	logger.Debugf("GetWAPBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetWAPBypassNetworkListsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get WAP bypass network lists request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateWAPBypassNetworkLists(ctx context.Context, params UpdateWAPBypassNetworkListsRequest) (*UpdateWAPBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("UpdateWAPBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateWAPBypassNetworkLists request: %w", err)
	}

	var result UpdateWAPBypassNetworkListsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update WAP bypass network lists request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) RemoveWAPBypassNetworkLists(ctx context.Context, params RemoveWAPBypassNetworkListsRequest) (*RemoveWAPBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("RemoveWAPBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/bypass-network-lists",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveWAPBypassNetworkLists request: %w", err)
	}

	var result RemoveWAPBypassNetworkListsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove WAP bypass network lists request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
