package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The IPGeo interface supports querying which network lists are used in the IP/Geo firewall settings,
	// as well as updating the method and which network lists are used for IP/Geo firewall blocking.
	IPGeo interface {
		// GetIPGeo lists which network lists are used in the IP/Geo Firewall settings.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-ip-geo-firewall
		GetIPGeo(ctx context.Context, params GetIPGeoRequest) (*GetIPGeoResponse, error)

		// UpdateIPGeo updates the method and which network lists to use for IP/Geo firewall blocking.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-policy-ip-geo-firewall
		UpdateIPGeo(ctx context.Context, params UpdateIPGeoRequest) (*UpdateIPGeoResponse, error)
	}

	// GetIPGeoRequest is used to retrieve the network lists used in IP/Geo firewall settings.
	GetIPGeoRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// IPGeoNetworkLists is used to specify IP or GEO network lists to be blocked or allowed.
	IPGeoNetworkLists struct {
		NetworkList []string `json:"networkList,omitempty"`
	}

	// IPGeoGeoControls is used to specify GEO network lists to be blocked.
	IPGeoGeoControls struct {
		BlockedIPNetworkLists *IPGeoNetworkLists `json:"blockedIPNetworkLists,omitempty"`
	}

	// IPGeoIPControls is used to specify IP or GEO network lists to be blocked or allowed.
	IPGeoIPControls struct {
		AllowedIPNetworkLists *IPGeoNetworkLists `json:"allowedIPNetworkLists,omitempty"`
		BlockedIPNetworkLists *IPGeoNetworkLists `json:"blockedIPNetworkLists,omitempty"`
	}

	// UpdateIPGeoRequest is used to update the method and which network lists are used for IP/Geo firewall blocking.
	UpdateIPGeoRequest struct {
		ConfigID    int               `json:"-"`
		Version     int               `json:"-"`
		PolicyID    string            `json:"-"`
		Block       string            `json:"block"`
		GeoControls *IPGeoGeoControls `json:"geoControls,omitempty"`
		IPControls  *IPGeoIPControls  `json:"ipControls,omitempty"`
	}

	// IPGeoFirewall is used to describe an IP/Geo firewall.
	IPGeoFirewall struct {
		Block       string            `json:"block"`
		GeoControls *IPGeoGeoControls `json:"geoControls,omitempty"`
		IPControls  *IPGeoIPControls  `json:"ipControls,omitempty"`
	}

	// GetIPGeoResponse is returned from a call to GetIPGeo
	GetIPGeoResponse IPGeoFirewall

	// UpdateIPGeoResponse is returned from a call to UpdateIPGeo
	UpdateIPGeoResponse IPGeoFirewall
)

// Validate validates a GetIPGeoRequest.
func (v GetIPGeoRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateIPGeoRequest.
func (v UpdateIPGeoRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetIPGeo(ctx context.Context, params GetIPGeoRequest) (*GetIPGeoResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIPGeo")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ip-geo-firewall",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetIPGeo request: %w", err)
	}

	var result GetIPGeoResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get IPGeo request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateIPGeo(ctx context.Context, params UpdateIPGeoRequest) (*UpdateIPGeoResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateIPGeo")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ip-geo-firewall",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateIPGeo request: %w", err)
	}

	var result UpdateIPGeoResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update IPGeo request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
