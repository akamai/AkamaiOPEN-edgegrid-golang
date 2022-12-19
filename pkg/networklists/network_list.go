package networklists

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The NetworkList interface supports creating, retrieving, modifying and removing network lists.
	//
	// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#networklist
	NetworkList interface {
		// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#getlists
		GetNetworkLists(ctx context.Context, params GetNetworkListsRequest) (*GetNetworkListsResponse, error)

		// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#getlist
		GetNetworkList(ctx context.Context, params GetNetworkListRequest) (*GetNetworkListResponse, error)

		// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#postlists
		CreateNetworkList(ctx context.Context, params CreateNetworkListRequest) (*CreateNetworkListResponse, error)

		// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#putlist
		UpdateNetworkList(ctx context.Context, params UpdateNetworkListRequest) (*UpdateNetworkListResponse, error)

		// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#deletelist
		RemoveNetworkList(ctx context.Context, params RemoveNetworkListRequest) (*RemoveNetworkListResponse, error)
	}

	// GetNetworkListRequest contains request parameters for GetNetworkList method
	GetNetworkListRequest struct {
		UniqueID string `json:"-"`
	}

	// GetNetworkListsRequest contains request parameters for GetNetworkLists method
	GetNetworkListsRequest struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	// GetNetworkListsResponse contains response from GetNetworkLists method
	GetNetworkListsResponse struct {
		Links        *NetworkListsResponseLinks           `json:"links,omitempty"`
		NetworkLists []GetNetworkListsResponseListElement `json:"networkLists"`
	}

	// GetNetworkListsResponseListElement contains information about a single network list
	GetNetworkListsResponseListElement struct {
		ElementCount       int                `json:"elementCount"`
		Links              *NetworkListsLinks `json:"links,omitempty"`
		Name               string             `json:"name"`
		NetworkListType    string             `json:"networkListType"`
		ReadOnly           bool               `json:"readOnly"`
		Shared             bool               `json:"shared"`
		SyncPoint          int                `json:"syncPoint"`
		Type               string             `json:"type"`
		UniqueID           string             `json:"uniqueId"`
		AccessControlGroup string             `json:"accessControlGroup,omitempty"`
		Description        string             `json:"description,omitempty"`
	}

	// GetNetworkListResponse contains response from GetNetworkList method
	GetNetworkListResponse struct {
		Name            string   `json:"name"`
		UniqueID        string   `json:"uniqueId"`
		ContractID      string   `json:"contractId"`
		GroupID         int      `json:"groupId"`
		SyncPoint       int      `json:"syncPoint"`
		Type            string   `json:"type"`
		Description     string   `json:"description,omitempty"`
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

	// CreateNetworkListRequest contains request parameters for CreateNetworkList method
	CreateNetworkListRequest struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		ContractID  string   `json:"contractId,omitempty"`
		GroupID     int      `json:"groupId,omitempty"`
		List        []string `json:"list"`
	}

	// UpdateNetworkListRequest contains request parameters for CreateNetworkLists method
	UpdateNetworkListRequest struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		ContractID  string   `json:"contractId,omitempty"`
		GroupID     int      `json:"groupId,omitempty"`
		SyncPoint   int      `json:"syncPoint"`
		List        []string `json:"list"`
		UniqueID    string   `json:"uniqueId"`
	}

	// UpdateNetworkListResponse contains response from CreateNetworkList method
	UpdateNetworkListResponse struct {
		Links struct {
			Create struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"create"`
		} `json:"links"`
		NetworkLists []struct {
			ElementCount int `json:"elementCount"`
			Links        struct {
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
			Name               string `json:"name"`
			NetworkListType    string `json:"networkListType"`
			ReadOnly           bool   `json:"readOnly"`
			Shared             bool   `json:"shared"`
			SyncPoint          int    `json:"syncPoint"`
			Type               string `json:"type"`
			UniqueID           string `json:"uniqueId"`
			AccessControlGroup string `json:"accessControlGroup,omitempty"`
			Description        string `json:"description,omitempty"`
		} `json:"networkLists"`
	}

	// RemoveNetworkListRequest contains request parameters for RemoveNetworkList method
	RemoveNetworkListRequest struct {
		UniqueID string `json:"-"`
	}

	// RemoveNetworkListResponse contains response from RemoveNetworkList method
	RemoveNetworkListResponse struct {
		Status    int    `json:"status"`
		UniqueID  string `json:"uniqueId"`
		SyncPoint int    `json:"syncPoint"`
	}

	// CreateNetworkListResponse contains response from CreateNetworkList method
	CreateNetworkListResponse struct {
		Name            string   `json:"name"`
		Description     string   `json:"description,omitempty"`
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

	// LinkInfo contains hypermedia link
	LinkInfo struct {
		Href   string `json:"href,omitempty"`
		Method string `json:"method,omitempty"`
	}

	// nolint:revive
	// NetworkListsResponseLinks contains LinkInfo
	NetworkListsResponseLinks struct {
		Create *LinkInfo `json:"create,omitempty"`
	}

	// nolint:revive
	// NetworkListsLinks encapsulates the set of API hypermedia
	NetworkListsLinks struct {
		ActivateInProduction *LinkInfo `json:"activateInProduction,omitempty"`
		ActivateInStaging    *LinkInfo `json:"activateInStaging,omitempty"`
		AppendItems          *LinkInfo `json:"appendItems,omitempty"`
		Retrieve             *LinkInfo `json:"retrieve,omitempty"`
		StatusInProduction   *LinkInfo `json:"statusInProduction,omitempty"`
		StatusInStaging      *LinkInfo `json:"statusInStaging,omitempty"`
		Update               *LinkInfo `json:"update,omitempty"`
	}
)

// Validate validates GetNetworkListRequest
func (v GetNetworkListRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
	}.Filter()
}

// Validate validates CreateNetworkListRequest
func (v CreateNetworkListRequest) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(v.Name, validation.Required),
	}.Filter()
}

// Validate validates UpdateNetworkListRequest
func (v UpdateNetworkListRequest) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(v.Name, validation.Required),
	}.Filter()
}

// Validate validates RemoveNetworkListRequest
func (v RemoveNetworkListRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
	}.Filter()
}

func (p *networklists) GetNetworkList(ctx context.Context, params GetNetworkListRequest) (*GetNetworkListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetNetworkList")

	var rval GetNetworkListResponse

	uri := fmt.Sprintf(
		"/network-list/v2/network-lists/%s",
		params.UniqueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklist request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklist request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *networklists) GetNetworkLists(ctx context.Context, params GetNetworkListsRequest) (*GetNetworkListsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetNetworkLists")

	var rval GetNetworkListsResponse
	var rvalfiltered GetNetworkListsResponse

	uri :=
		"/network-list/v2/network-lists"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklists request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklists request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Name == "" && params.Type == "" {
		return &rval, nil
	}

	rvalfiltered.Links = rval.Links
	for _, val := range rval.NetworkLists {
		if (params.Name == "" || params.Name == val.Name) && (params.Type == "" || params.Type == val.Type) {
			rvalfiltered.NetworkLists = append(rvalfiltered.NetworkLists, val)
		}
	}

	return &rvalfiltered, nil
}

// Update will update a NetworkList.
//
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#putnetworklist

func (p *networklists) UpdateNetworkList(ctx context.Context, params UpdateNetworkListRequest) (*UpdateNetworkListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkList")

	putURL := fmt.Sprintf(
		"/network-list/v2/network-lists/%s",
		params.UniqueID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update NetworkListrequest: %s", err.Error())
	}

	var rval UpdateNetworkListResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("update NetworkList request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// CreateNetworkList will create a new networklist.
//
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#postnetworklist
func (p *networklists) CreateNetworkList(ctx context.Context, params CreateNetworkListRequest) (*CreateNetworkListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateNetworkList")

	uri :=
		"/network-list/v2/network-lists"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create networklist request: %s", err.Error())
	}

	var rval CreateNetworkListResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create networklistrequest failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a NetworkList
//
//
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#deletenetworklist

func (p *networklists) RemoveNetworkList(ctx context.Context, params RemoveNetworkListRequest) (*RemoveNetworkListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveNetworkListResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveNetworkList")

	uri, err := url.Parse(fmt.Sprintf(
		"/network-list/v2/network-lists/%s",
		params.UniqueID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveNetworkList request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("RemoveNetworkList request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
