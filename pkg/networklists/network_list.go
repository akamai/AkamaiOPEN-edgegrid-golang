package networklists

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation"
)

// NetworkList represents a collection of NetworkList
//
// See: NetworkList.GetNetworkList()
// API Docs: // network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html

type (
	// NetworkList  contains operations available on NetworkList  resource
	// See: // network_lists v2
	//
	// https://developer.akamai.com/api/cloud_security/network_lists/v2.html#getnetworklist
	NetworkList interface {
		GetNetworkLists(ctx context.Context, params GetNetworkListsRequest) (*GetNetworkListsResponse, error)
		GetNetworkList(ctx context.Context, params GetNetworkListRequest) (*GetNetworkListResponse, error)
		CreateNetworkList(ctx context.Context, params CreateNetworkListRequest) (*CreateNetworkListResponse, error)
		UpdateNetworkList(ctx context.Context, params UpdateNetworkListRequest) (*UpdateNetworkListResponse, error)
		RemoveNetworkList(ctx context.Context, params RemoveNetworkListRequest) (*RemoveNetworkListResponse, error)
	}

	GetNetworkListRequest struct {
		Name string `json:"name"`
	}

	GetNetworkListsRequest struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	// prevent these empty sections from being returned:
	/*
	{
		"links": {
			"create": {
				"href": "",
				"method": ""
			}
		}
	}
	*/
	GetNetworkListsResponse struct {
		Links *NetworkListsResponseLinks `json:"links,omitempty"`
		NetworkLists []struct {
			ElementCount int `json:"elementCount"`
			Links              *NetworkListsLinks `json:"links,omitempty"`
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

	GetNetworkListResponse struct {
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

	CreateNetworkListRequest struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		List        []string `json:"list"`
	}

	UpdateNetworkListRequest struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		SyncPoint   int      `json:"syncPoint"`
		List        []string `json:"list"`
	}

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

	RemoveNetworkListRequest struct {
		Name string `json:"name"`
	}

	RemoveNetworkListResponse struct {
		Status    int    `json:"status"`
		UniqueID  string `json:"uniqueId"`
		SyncPoint int    `json:"syncPoint"`
	}

	CreateNetworkListResponse struct {
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

	LinkInfo struct {
		Href   string `json:"href,omitempty"`
		Method string `json:"method,omitempty"`
	}

	NetworkListsResponseLinks struct {
		Create *LinkInfo `json:"create,omitempty"`
	}

	NetworkListsLinks        struct {
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
		"Name": validation.Validate(v.Name, validation.Required),
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
		"Name": validation.Validate(v.Name, validation.Required),
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
		params.Name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklist request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
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
		return nil, fmt.Errorf("failed to create getnetworklists request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklists request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Name != "" && params.Type != "" {
		for _, val := range rval.NetworkLists {
			if val.Name == params.Name && val.Type == params.Type {
				rvalfiltered.NetworkLists = append(rvalfiltered.NetworkLists, val)
			}
		}
	} else if params.Name != "" {
		for _, val := range rval.NetworkLists {
			if val.Name == params.Name {
				rvalfiltered.NetworkLists = append(rvalfiltered.NetworkLists, val)
			}
		}
	} else if params.Type != "" {
		for _, val := range rval.NetworkLists {
			if val.Type == params.Type {
				rvalfiltered.NetworkLists = append(rvalfiltered.NetworkLists, val)
			}
		}
	} else {
		rvalfiltered = rval
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
		params.Name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create NetworkListrequest: %w", err)
	}

	var rval UpdateNetworkListResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create NetworkList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Create will create a new networklist.
//
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
		return nil, fmt.Errorf("failed to create create networklist request: %w", err)
	}

	var rval CreateNetworkListResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create networklistrequest failed: %w", err)
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
		params.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delnetworklist request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delnetworklist request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
