package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// IPGeo represents a collection of IPGeo
//
// See: IPGeo.GetIPGeo()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// IPGeo  contains operations available on IPGeo  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getipgeo
	IPGeo interface {
		GetIPGeo(ctx context.Context, params GetIPGeoRequest) (*GetIPGeoResponse, error)
		UpdateIPGeo(ctx context.Context, params UpdateIPGeoRequest) (*UpdateIPGeoResponse, error)
	}

	GetIPGeoRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetIPGeoResponse struct {
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

// Validate validates GetIPGeoRequest
func (v GetIPGeoRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateIPGeoRequest
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
		return nil, fmt.Errorf("failed to create getipgeo request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getipgeo  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a IPGeo.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putipgeo

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
		return nil, fmt.Errorf("failed to create create IPGeorequest: %w", err)
	}

	var rval UpdateIPGeoResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create IPGeo request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
