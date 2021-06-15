package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// BypassNetworkLists represents a collection of BypassNetworkLists
//
// See: BypassNetworkLists.GetBypassNetworkLists()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// BypassNetworkLists  contains operations available on BypassNetworkLists  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getbypassnetworklists
	BypassNetworkLists interface {
		GetBypassNetworkLists(ctx context.Context, params GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error)
		UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error)
		RemoveBypassNetworkLists(ctx context.Context, params RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error)
	}

	GetBypassNetworkListsRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	GetBypassNetworkListsResponse struct {
		NetworkLists []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"networkLists"`
	}

	UpdateBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		NetworkLists []string `json:"networkLists"`
	}

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

	RemoveBypassNetworkListsRequest struct {
		ConfigID     int      `json:"-"`
		Version      int      `json:"-"`
		NetworkLists []string `json:"networkLists"`
	}
	RemoveBypassNetworkListsResponse struct {
		NetworkLists []string `json:"networkLists"`
	}
)

// Validate validates GetBypassNetworkListsRequest
func (v GetBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateBypassNetworkListsRequest
func (v UpdateBypassNetworkListsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates RemoveBypassNetworkListsRequest
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
	logger.Debug("GetBypassNetworkLists")

	var rval GetBypassNetworkListsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getbypassnetworklists request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getbypassnetworklists  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a BypassNetworkLists.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklists

func (p *appsec) UpdateBypassNetworkLists(ctx context.Context, params UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateBypassNetworkLists")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create BypassNetworkListsrequest: %w", err)
	}

	var rval UpdateBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create BypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Remove will Remove a BypassNetworkLists.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putbypassnetworklists

func (p *appsec) RemoveBypassNetworkLists(ctx context.Context, params RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveBypassNetworkLists")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/bypass-network-lists",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create BypassNetworkListsrequest: %w", err)
	}

	var rval RemoveBypassNetworkListsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("remove BypassNetworkLists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
