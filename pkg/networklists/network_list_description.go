package networklists

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The NetworkListDescription interface supports retrieving and updating a network list's description.
	NetworkListDescription interface {
		// GetNetworkListDescription retrieves network list with description.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/get-network-list
		GetNetworkListDescription(ctx context.Context, params GetNetworkListDescriptionRequest) (*GetNetworkListDescriptionResponse, error)

		// UpdateNetworkListDescription modifies network list description.
		//
		// See: https://techdocs.akamai.com/network-lists/reference/put-network-list-details
		UpdateNetworkListDescription(ctx context.Context, params UpdateNetworkListDescriptionRequest) (*UpdateNetworkListDescriptionResponse, error)
	}

	// GetNetworkListDescriptionRequest contains request parameters for GetNetworkListDescription method
	GetNetworkListDescriptionRequest struct {
		UniqueID    string `json:"uniqueId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// GetNetworkListDescriptionResponse contains response from GetNetworkListDescription method
	GetNetworkListDescriptionResponse struct {
		Name            string   `json:"name"`
		UniqueID        string   `json:"uniqueId"`
		Description     string   `json:"description"`
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

	// UpdateNetworkListDescriptionRequest contains request parameters for UpdateNetworkListDescription method
	UpdateNetworkListDescriptionRequest struct {
		UniqueID    string `json:"-"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// UpdateNetworkListDescriptionResponse contains response from UpdateNetworkListDescription method
	UpdateNetworkListDescriptionResponse struct {
		Empty string `json:"-"`
	}
)

// Validate validates GetNetworkListDescriptionRequest
func (v GetNetworkListDescriptionRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
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
		params.UniqueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklistdescription request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklistdescription  request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create create NetworkListDescriptionrequest: %s", err.Error())
	}

	var rval UpdateNetworkListDescriptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create NetworkListDescription request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
