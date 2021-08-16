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
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#ipgeofirewall
	IPGeo interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getipgeofirewall
		GetIPGeo(ctx context.Context, params GetIPGeoRequest) (*GetIPGeoResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putipgeofirewall
		UpdateIPGeo(ctx context.Context, params UpdateIPGeoRequest) (*UpdateIPGeoResponse, error)
	}

	// GetIPGeoRequest is used to retrieve the network lists used in IP/Geo firewall settings.
	GetIPGeoRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	// GetIPGeoResponse is returned from a call to GetIpGeo.
	GetIPGeoResponse struct {
		Block       string `json:"block,omitempty"`
		GeoControls struct {
			BlockedIPNetworkLists struct {
				NetworkList []string `json:"networkList,omitempty"`
			} `json:"blockedIPNetworkLists,omitempty"`
		} `json:"geoControls,omitempty"`
		IPControls struct {
			AllowedIPNetworkLists struct {
				NetworkList []string `json:"networkList,omitempty"`
			} `json:"allowedIPNetworkLists,omitempty"`
			BlockedIPNetworkLists struct {
				NetworkList []string `json:"networkList,omitempty"`
			} `json:"blockedIPNetworkLists,omitempty"`
		} `json:"ipControls,omitempty"`
	}

	// UpdateIPGeoRequest is used to update the method and which network lists are used for IP/Geo firewall blocking.
	UpdateIPGeoRequest struct {
		ConfigID    int    `json:"-"`
		Version     int    `json:"-"`
		PolicyID    string `json:"-"`
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

	// UpdateIPGeoResponse is returned from a call to UpdateIPGeo.
	UpdateIPGeoResponse struct {
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
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetIPGeo")

	var rval GetIPGeoResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ip-geo-firewall",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetIPGeo request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetIPGeo request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateIPGeo(ctx context.Context, params UpdateIPGeoRequest) (*UpdateIPGeoResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateIPGeo")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/ip-geo-firewall",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateIPGeo request: %w", err)
	}

	var rval UpdateIPGeoResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateIPGeo request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
