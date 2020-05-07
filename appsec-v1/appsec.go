package appsec

import (
	"bytes"
	"encoding/json"
	"fmt"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"time"
)

type ActivationConfig struct {
	ConfigID              int    `json:"configId"`
	ConfigName            string `json:"configName"`
	ConfigVersion         int    `json:"configVersion"`
	PreviousConfigVersion int    `json:"previousConfigVersion"`
}

type Activation struct {
	Action             string             `json:"action"`
	Network            string             `json:"network"`
	Note               string             `json:"note"`
	NotificationEmails []string           `json:"notificationEmails"`
	ActivationConfigs  []ActivationConfig `json:"activationConfigs"`
}

type ActivationResponse struct {
	ResponseCode            int
	ActivationStatus        ActivationStatus
	ActivationRequestStatus ActivationRequestStatusCreated
}

type ActivationStatus struct {
	DispatchCount     int                `json:"dispatchCount"`
	ActivationID      int                `json:"activationId"`
	Action            string             `json:"action"`
	Status            string             `json:"status"`
	Network           string             `json:"network"`
	Estimate          string             `json:"estimate"`
	CreatedBy         string             `json:"createdBy"`
	CreateDate        time.Time          `json:"createDate"`
	ActivationConfigs []ActivationConfig `json:"activationConfigs"`
}

type ActivationRequestStatusCreated struct {
	StatusID   string    `json:"statusId"`
	CreateDate time.Time `json:"createDate"`
	Links      struct {
		CheckStatus struct {
			Href string `json:"href"`
		} `json:"check-status"`
	} `json:"links"`
}

type ActivationRequestStatusResponse struct {
	ResponseCode                      int
	ActivationRequestStatusInProgress ActivationRequestStatusInProgress
	ActivationRequestStatusComplete   ActivationRequestStatusComplete
}

type ActivationRequestStatusInProgress struct {
	StatusID   string    `json:"statusId"`
	CreateDate time.Time `json:"createDate"`
}

type ActivationRequestStatusComplete struct {
	ActivationID int `json:"activationId"`
}

type ConfigurationClone struct {
	CreateFromVersion int  `json:"createFromVersion"`
	RuleUpdate        bool `json:"ruleUpdate"`
}

type HostnameList struct {
	Hostname string `json:"hostname"`
}

type SelectedHostnames struct {
	HostnameList []HostnameList `json:"hostnameList"`
}

type VersionList struct {
	TotalSize                   int    `json:"totalSize"`
	PageSize                    int    `json:"pageSize"`
	Page                        int    `json:"page"`
	ConfigID                    int    `json:"configId"`
	ConfigName                  string `json:"configName"`
	StagingExpediteRequestID    int    `json:"stagingExpediteRequestId"`
	ProductionExpediteRequestID int    `json:"productionExpediteRequestId"`
	ProductionActiveVersion     int    `json:"productionActiveVersion"`
	StagingActiveVersion        int    `json:"stagingActiveVersion"`
	LastCreatedVersion          int    `json:"lastCreatedVersion"`
}

type Version struct {
	ConfigID     int       `json:"configId"`
	ConfigName   string    `json:"configName"`
	Version      int       `json:"version"`
	VersionNotes string    `json:"versionNotes"`
	CreateDate   time.Time `json:"createDate"`
	CreatedBy    string    `json:"createdBy"`
	BasedOn      int       `json:"basedOn"`
	Production   struct {
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	} `json:"production"`
	Staging struct {
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	} `json:"staging"`
}

func ListConfigurationVersions(configid int) (*VersionList, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/configs/%d/versions", configid),
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

	var response VersionList
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func CloneConfigurationVersion(configid int, configurationclone ConfigurationClone) (*Version, error) {

	r, err := json.Marshal(configurationclone)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		fmt.Sprintf("/appsec/v1/configs/%d/versions", configid),
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

	var response Version
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func ActivateConfigurationVersion(activation Activation) (*ActivationResponse, error) {

	r, err := json.Marshal(activation)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"POST",
		"/appsec/v1/activations",
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

	var response ActivationResponse
	var activationresponse ActivationStatus
	var activationrequeststatus ActivationRequestStatusCreated

	// This pesky API call can return different responses!
	if res.StatusCode == 200 {
		if err = client.BodyJSON(res, &activationresponse); err != nil {
			return nil, err
		}
	} else if res.StatusCode == 202 {
		if err = client.BodyJSON(res, &activationrequeststatus); err != nil {
			return nil, err
		}
	}

	response.ResponseCode = res.StatusCode
	response.ActivationStatus = activationresponse
	response.ActivationRequestStatus = activationrequeststatus

	return &response, nil
}

func GetActivationRequestStatus(statusid string) (*ActivationRequestStatusResponse, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/activations/status/%s", statusid),
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

	var response ActivationRequestStatusResponse
	var activationrequeststatusinprogress ActivationRequestStatusInProgress
	var activationrequeststatuscomplete ActivationRequestStatusComplete

	// This pesky API call can return different responses!
	if res.StatusCode == 200 {
		if err = client.BodyJSON(res, &activationrequeststatusinprogress); err != nil {
			return nil, err
		}
	} else if res.StatusCode == 303 {
		if err = client.BodyJSON(res, &activationrequeststatuscomplete); err != nil {
			return nil, err
		}
	}

	response.ResponseCode = res.StatusCode
	response.ActivationRequestStatusInProgress = activationrequeststatusinprogress
	response.ActivationRequestStatusComplete = activationrequeststatuscomplete

	return &response, nil
}

func GetConfigurationVersion(configid int, version int) (*Version, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/configs/%d/versions/%d", configid, version),
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

func GetActivationStatus(activationid int) (*ActivationStatus, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/activations/%d", activationid),
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

func ListSelectedHostnames(configid int, version int) (*SelectedHostnames, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/selected-hostnames", configid, version),
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

	var response SelectedHostnames
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func UpdateSelectedHostnames(configid int, version int, selectedhostnames SelectedHostnames) (*SelectedHostnames, error) {

	r, err := json.Marshal(selectedhostnames)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PUT",
		fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/selected-hostnames", configid, version),
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

	var response SelectedHostnames
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
