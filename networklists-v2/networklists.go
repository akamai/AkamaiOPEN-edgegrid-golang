package networklists

import (
	"fmt"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type ListNetworkListsResponse struct {
	NetworkLists []struct {
		NetworkListType    string `json:"networkListType"`
		AccessControlGroup string `json:"accessControlGroup,omitempty"`
		Name               string `json:"name"`
		ElementCount       int    `json:"elementCount"`
		ReadOnly           bool   `json:"readOnly"`
		Shared             bool   `json:"shared"`
		SyncPoint          int    `json:"syncPoint"`
		Type               string `json:"type"`
		UniqueID           string `json:"uniqueId"`
		Links              struct {
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
	} `json:"networkLists"`
}

type NetworkList struct {
	Name            string   `json:"name"`
	UniqueID        string   `json:"uniqueId"`
	SyncPoint       int      `json:"syncPoint"`
	Type            string   `json:"type"`
	NetworkListType string   `json:"networkListType"`
	ElementCount    int      `json:"elementCount"`
	ReadOnly        bool     `json:"readOnly"`
	Shared          bool     `json:"shared"`
	List            []string `json:"list"`
}

type Message struct {
	Status    int    `json:"status"`
	UniqueID  string `json:"uniqueId"`
	SyncPoint int    `json:"syncPoint"`
}



func ListNetworkLists() (*ListNetworkListsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/network-list/v2/network-lists?extended=true&includeElements=false",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListNetworkListsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateNetworkList(networklist NetworkList) (*NetworkList, error) {
	req, err := client.NewRequest(
		Config,
		"POST",
		"/network-list/v2/network-lists"
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response NetworkList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func UpdateNetworkList(id string, networklist NetworkList) (*NetworkList, error) {
	req, err := client.NewRequest(
		Config,
		"PUT",
		fmt.Sprintf("/network-list/v2/network-lists/%s?extended=true&includeElements=false", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response NetworkList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetNetworkList(id string) (*NetworkList, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/network-list/v2/network-lists/%s?extended=true&includeElements=false", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response NetworkList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func DeleteNetworkList(id string) (*Message, error) {
	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/network-list/v2/network-lists/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Message
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

