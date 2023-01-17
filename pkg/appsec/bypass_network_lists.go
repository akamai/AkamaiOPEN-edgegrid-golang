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
	// Deprecated: this interface will be removed in a future release. Use the WAPBypassNetworkLists interface instead.
	BypassNetworkLists interface {
		// Deprecated: this method will be removed in a future release. Use the GetWAPBypassNetworkLists method of the WAPBypassNetworkLists interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-bypass-network-lists
		GetBypassNetworkLists(ctx context.Context, params GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error)

		// Deprecated: this method will be removed in a future release. Use the UpdateWAPBypassNetworkLists method of the WAPBypassNetworkLists interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-bypass-network-lists
		UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error)

		// Deprecated: this method will be removed in a future release. Use the UpdateWAPBypassNetworkLists method of the WAPBypassNetworkLists interface instead.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-bypass-network-lists
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
		NetworkLists []NetworkList `json:"networkLists"`
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
		NetworkLists []string `json:"networkLists"`
	}

	// RemoveBypassNetworkListsRequest is used to modify which network lists are used in the bypass network lists settings.
	// Deprecated: this struct will be removed in a future release.
	RemoveBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		PolicyID     string   `json:"policyId"`
		NetworkLists []string `json:"networkLists"`
	}

	// RemoveBypassNetworkListsResponse is returned from a call to RemoveBypassNetworkLists.
	// Deprecated: this struct will be removed in a future release.
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
// Deprecated: this method will be removed in a future release.
func (v RemoveBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetBypassNetworkLists(ctx context.Context, params GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("GetBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetBypassNetworkListsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get bypass network lists request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("UpdateBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateBypassNetworkLists request: %w", err)
	}

	var result UpdateBypassNetworkListsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update bypass network lists request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

// Deprecated: this method will be removed in a future release.
func (p *appsec) RemoveBypassNetworkLists(ctx context.Context, params RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error) {
	logger := p.Log(ctx)
	logger.Debugf("RemoveBypassNetworkLists")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveBypassNetworkLists request: %w", err)
	}

	var result RemoveBypassNetworkListsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove bypass network lists request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
