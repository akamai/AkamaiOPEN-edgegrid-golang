package networklists

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// NetworkListDescription represents a collection of NetworkListDescription
//
// See: NetworkListDescription.GetNetworkListDescription()
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html

type (
	// NetworkListDescription  contains operations available on NetworkListDescription  resource
	// See: // network_lists v2
	//
	// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#getnetworklistdescription
	NetworkListDescription interface {
		//GetNetworkListDescriptions(ctx context.Context, params GetNetworkListDescriptionsRequest) (*GetNetworkListDescriptionsResponse, error)
		GetNetworkListDescription(ctx context.Context, params GetNetworkListDescriptionRequest) (*GetNetworkListDescriptionResponse, error)
		UpdateNetworkListDescription(ctx context.Context, params UpdateNetworkListDescriptionRequest) (*UpdateNetworkListDescriptionResponse, error)
	}

	GetNetworkListDescriptionRequest struct {
		UniqueID    string `json:"uniqueId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	GetNetworkListDescriptionResponse struct {
		Name            string   `json:"name"`
		UniqueID        string   `json:"uniqueId"`
		SyncPoint       int      `json:"syncPoint"`
		Type            string   `json:"type"`
		NetworkListType string   `json:"networkListType"`
		ElementCount    int      `json:"elementCount"`
		ReadOnly        bool     `json:"readOnly"`
		Shared          bool     `json:"shared"`
		List            []string `json:"list"`
		Links           struct {
			ActivateInProduction struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"activateInProduction"`
			ActivateInStaging struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"activateInStaging"`
			AppendItems struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"appendItems"`
			Retrieve struct {
				Href string `json:"href"`
			} `json:"retrieve"`
			StatusInProduction struct {
				Href string `json:"href"`
			} `json:"statusInProduction"`
			StatusInStaging struct {
				Href string `json:"href"`
			} `json:"statusInStaging"`
			Update struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"update"`
		} `json:"links"`
	}

	UpdateNetworkListDescriptionRequest struct {
		UniqueID    string `json:"-"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UpdateNetworkListDescriptionResponse struct {
		Empty string `json:"-"`
	}
)

// Validate validates GetNetworkListDescriptionRequest
func (v GetNetworkListDescriptionRequest) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(v.Name, validation.Required),
	}.Filter()
}

// Validate validates UpdateNetworkListDescriptionRequest
func (v UpdateNetworkListDescriptionRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
	}.Filter()
}

func (p *networklists) GetNetworkListDescription(ctx context.Context, params GetNetworkListDescriptionRequest) (*GetNetworkListDescriptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetNetworkListDescription")

	var rval GetNetworkListDescriptionResponse

	uri := fmt.Sprintf(
		"/network-list/v2/network-lists/%s",
		params.Name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklistdescription request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklistdescription  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a NetworkListDescription.
//
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#putnetworklistdescription

func (p *networklists) UpdateNetworkListDescription(ctx context.Context, params UpdateNetworkListDescriptionRequest) (*UpdateNetworkListDescriptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkListDescription")

	putURL := fmt.Sprintf(
		"/network-list/v2/network-lists/%s/details",
		params.UniqueID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create NetworkListDescriptionrequest: %w", err)
	}

	var rval UpdateNetworkListDescriptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create NetworkListDescription request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
