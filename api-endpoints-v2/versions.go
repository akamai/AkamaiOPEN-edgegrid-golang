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

func ListVersions(endpointId int) (*Versions, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions",
			endpointId,
		),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Versions{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func GetVersion(endpointId, version int) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/resources-detail",
			endpointId,
			version,
		),
		nil,
	)

	return call(req, err)
}

func ModifyVersion(endpoint *Endpoint) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d",
			endpoint.APIEndPointID,
			endpoint.VersionNumber,
		),
		endpoint,
	)

	return call(req, err)
}

func CloneVersion(endpointId, version int) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/cloneVersion",
			endpointId,
			version,
		),
		nil,
	)

	return call(req, err)
}

func RemoveVersion(endpointId, version int) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d",
			endpointId,
			version,
		),
		nil,
	)

	return call(req, err)
}

func GetLatestVersionNumber(endpointId int) (int, error) {
	versions, err := ListVersions(endpointId)
	if err != nil {
		return 0, err
	}

	loc := len(versions.APIVersions) - 1
	v := versions.APIVersions[loc]
	return v.VersionNumber, nil
}
