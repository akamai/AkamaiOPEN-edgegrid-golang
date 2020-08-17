package edgeworkers

import (
	"fmt"
	"io/ioutil"
	"time"
	"encoding/json"
	"bytes"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type ActivationList struct {
	Activations []Activation `json:"activations"`
}

type Activation struct {
	EdgeWorkerID     int       `json:"edgeWorkerId"`
	Version          int       `json:"version"`
	ActivationID     int       `json:"activationId"`
	AccountID        string    `json:"accountId"`
	Status           string    `json:"status"`
	Network          string    `json:"network"`
	CreatedBy        string    `json:"createdBy"`
	CreatedTime      time.Time `json:"createdTime"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
}
	
type Validation struct {
	Errors []struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"errors"`
}

type EdgeWorkerList struct {
	EdgeWorkerIDs []EdgeWorker `json:"edgeWorkerIds"`
}
	
type EdgeWorker struct {
	EdgeWorkerID     int       `json:"edgeWorkerId"`
	Name             string    `json:"name"`
	AccountID        string    `json:"accountId"`
	GroupID          int       `json:"groupId"`
	CreatedBy        string    `json:"createdBy"`
	CreatedTime      time.Time `json:"createdTime"`
	LastModifiedBy   string    `json:"lastModifiedBy"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
}

type Version struct {
	EdgeWorkerID int       `json:"edgeWorkerId"`
	Version      int       `json:"version"`
	AccountID    string    `json:"accountId"`
	Checksum     string    `json:"checksum"`
	CreatedBy    string    `json:"createdBy"`
	CreatedTime  time.Time `json:"createdTime"`
}
	
type Permissions struct {
	GroupID      int      `json:"groupId"`
	GroupName    string   `json:"groupName"`
	Capabilities []string `json:"capabilities"`
}

type PermissionsList struct {
	Groups []Permissions `json:"groups"`
}

func ListPermissionsGroups() (*PermissionsList, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/edgeworkers/v1/groups",
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

	var response PermissionsList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetPermissionsGroup(id int) (*Permissions, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/groups/%d", id),
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

	var response Permissions
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func ListEdgeWorkerIDs() (*EdgeWorkerList, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/edgeworkers/v1/ids",
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

	var response EdgeWorkerList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateEdgeWorker(edgeworker EdgeWorker) (*EdgeWorker, error) {

        r, err := json.Marshal(edgeworker)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		"/edgeworkers/v1/ids",
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

	var response EdgeWorker
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetEdgeWorker(id int) (*EdgeWorker, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/ids/%d", id),
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

	var response EdgeWorker
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func UpdateEdgeWorker(id int, edgeworker EdgeWorker) (*EdgeWorker, error) {

        r, err := json.Marshal(edgeworker)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PUT",
		fmt.Sprintf("/edgeworkers/v1/ids/%d", id),
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

	var response EdgeWorker
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func ValidateEdgeWorkerBundle(r []byte) (*Validation, error) {

	req, err := client.NewRequest(
		Config,
		"POST",
		"/edgeworkers/v1/validations",
		bytes.NewReader(r),
	)

	if err != nil {
		return nil, err
	}

        req.Header.Set("Content-Type", "application/gzip")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Validation
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateVersion(id int, r []byte) (*Version, error) {

	req, err := client.NewRequest(
		Config,
		"POST",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/versions", id),
		bytes.NewReader(r),
	)

	if err != nil {
		return nil, err
	}

        req.Header.Set("Content-Type", "application/gzip")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Version
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetVersion(id int, version string) (*Version, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/versions/%s", id, version),
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

	var response Version
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetBundle(id int, version string) ([]byte, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/versions/%s/content", id, version),
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ListActivations(id int) (*ActivationList, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/activations", id),
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

	var response ActivationList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateActivation(id int, activation Activation) (*Activation, error) {

        r, err := json.Marshal(activation)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/activations", id),
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

	var response Activation
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetActivation(id int, activationid int) (*Activation, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/edgeworkers/v1/ids/%d/activations/%d", id, activationid),
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

	var response Activation
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
