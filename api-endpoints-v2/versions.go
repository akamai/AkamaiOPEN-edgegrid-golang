package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Versions struct {
	APIEndPointID   int       `json:"apiEndPointId"`
	APIEndPointName string    `json:"apiEndPointName"`
	APIVersions     []Version `json:"apiVersions"`
}

type Version struct {
	CreatedBy            string       `json:"createdBy"`
	CreateDate           string       `json:"createDate"`
	UpdateDate           string       `json:"updateDate"`
	UpdatedBy            string       `json:"updatedBy"`
	APIEndPointVersionID int          `json:"apiEndPointVersionId"`
	BasePath             string       `json:"basePath"`
	Description          *string      `json:"description`
	BasedOn              *int         `json:"basedOn"`
	StagingStatus        *StatusValue `json:"stagingStatus"`
	ProductionStatus     *StatusValue `json:"productionStatus"`
	StagingDate          *string      `json:"stagingDate"`
	ProductionDate       *string      `json:"productionDate"`
	IsVersionLocked      bool         `json:"isVersionLocked"`
	AvailableActions     []string     `json:"availableActions"`
	VersionNumber        int          `json:"versionNumber"`
	LockVersion          int          `json:"lockVersion"`
}

type VersionSummary struct {
	Status        StatusValue `json:"status,omitempty"`
	VersionNumber int         `json:"versionNumber,omitempty"`
}

type StatusValue string

const (
	StatusPending     string = "PENDING"
	StatusActive      string = "ACTIVE"
	StatusDeactivated string = "DEACTIVATED"
	StatusFailed      string = "FAILED"
)

type GetVersionOptions struct {
	EndpointId string
	Version    string
}

func GetVersion(options *GetVersionOptions) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%s/versions/%s/resources-detail",
			options.EndpointId,
			options.Version,
		),
		options,
	)

	return call(req, err)
}

type ModifyVersionOptions struct {
	EndpointId  string   `json:"-"`
	Version     string   `json:"-"`
	Name        string   `json:"apiEndPointName,omitempty"`
	Description string   `json:"description,omitempty"`
	BasePath    string   `json:"basePath,omitempty"`
	Hostnames   []string `json:"apiEndPointHosts,omitempty"`
	Scheme      string   `json:"apiEndPointScheme,omitempty"`
}

func ModifyVersion(options *ModifyVersionOptions) (*Endpoint, error) {
	ep, err := GetVersion(&GetVersionOptions{
		options.EndpointId,
		options.Version,
	})

	if err != nil {
		return nil, err
	}

	if IsActive(ep, "production") || IsActive(ep, "staging") {

	}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%s/versions/%s",
			options.EndpointId,
			options.Version,
		),
		options,
	)

	return call(req, err)
}

type RemoveVersionOptions struct {
	EndpointId    int
	VersionNumber int
}

func RemoveVersion(options *RemoveVersionOptions) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d",
			options.EndpointId,
			options.VersionNumber,
		),
		nil,
	)

	return call(req, err)
}
