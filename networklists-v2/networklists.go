package networklists

import (
	"fmt"
	"encoding/json"
	"bytes"
	"errors"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type ListNetworkListsResponse struct {
	NetworkLists []NetworkList `json:"networkLists"`
}

type NetworkList struct {
	Name            string   `json:"name"`
	UniqueID        string   `json:"uniqueId"`
	SyncPoint       int      `json:"syncPoint"`
	Type            string   `json:"type"`
	NetworkListType string   `json:"networkListType"`
	Description	string   `json:"description"`
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

type ActivationRequest struct {
	UniqueID string
	Network Network
	Comments               string   `json:"comments"`
	NotificationRecipients []string `json:"notificationRecipients"`
}

type ActivationStatus struct {
	ActivationID       int    `json:"activationId"`
	ActivationComments string `json:"activationComments"`
	ActivationStatus   string `json:"activationStatus"`
	SyncPoint          int    `json:"syncPoint"`
	UniqueID           string `json:"uniqueId"`
	Fast               bool   `json:"fast"`
	DispatchCount      int    `json:"dispatchCount"`
}

type Network string
const (
	Staging Network = "STAGING"
	Production Network = "PRODUCTION"
)



func ListNetworkLists() (*ListNetworkListsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/network-list/v2/network-lists?includeElements=false",
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

        r, err := json.Marshal(networklist)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		"/network-list/v2/network-lists",
		bytes.NewReader(r),
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

func UpdateNetworkList(networklist NetworkList) (*NetworkList, error) {

	id := networklist.UniqueID
	if id == "" {
		return nil, errors.New("Error: no UniqueID in NetworkList")
	}

        r, err := json.Marshal(networklist)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PUT",
		fmt.Sprintf("/network-list/v2/network-lists/%s?includeElements=false", id),
		bytes.NewReader(r),
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

func ActivateNetworkList(activationrequest ActivationRequest) (*ActivationStatus, error) {
	id := activationrequest.UniqueID
	if id == "" {
		return nil, errors.New("Error: no UniqueID in ActivationRequest")
	}

	network := activationrequest.Network
	if network == "" {
		return nil, errors.New("Error: no Network in ActivationRequest")
	}

        r, err := json.Marshal(activationrequest)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/activate", id, network),
		bytes.NewReader(r),
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

	var response ActivationStatus
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetActivationStatus(id string, network Network) (*ActivationStatus, error) {
	if id == "" {
		return nil, errors.New("Error: no UniqueID in ActivationRequest")
	}

	if network == "" {
		return nil, errors.New("Error: no Network in ActivationRequest")
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/status", id, network),
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

	var response ActivationStatus
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
