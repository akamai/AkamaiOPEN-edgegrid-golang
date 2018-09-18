package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/google/go-querystring/query"
)

type EndpointList struct {
	APIEndPoints Endpoints `json:"apiEndPoints"`
	Links        Links     `json:"links"`
	Page         int       `json:"page"`
	PageSize     int       `json:"pageSize"`
	TotalSize    int       `json:"totalSize"`
}

func (list *EndpointList) GetEndpointsList(options ListEndpointOptions) error {
	q, err := query.Values(options)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"/api-definitions/v2/endpoints?%s",
		q.Encode(),
	)

	req, err := client.NewJSONRequest(Config, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, list); err != nil {
		return err
	}

	return nil
}

type Endpoints []Endpoint

type Endpoint struct {
	APICategoryIds      []int              `json:"apiCategoryIds"`
	APIEndPointHosts    []string           `json:"apiEndPointHosts"`
	APIEndPointID       int                `json:"apiEndPointId,omitempty"`
	APIEndPointLocked   bool               `json:"apiEndPointLocked,omitempty"`
	APIEndPointName     string             `json:"apiEndPointName"`
	APIEndPointScheme   string             `json:"apiEndPointScheme"`
	APIResourceBaseInfo []ResourceBaseInfo `json:"apiResourceBaseInfo"`
	BasePath            string             `json:"basePath"`
	ClonedFromVersion   *int               `json:"clonedFromVersion,omitempty"`
	ConsumeType         string             `json:"consumeType"`
	ContractID          string             `json:"contractId,omitempty"`
	CreateDate          string             `json:"createDate,omitempty"`
	CreatedBy           string             `json:"createdBy,omitempty"`
	Description         string             `json:"description"`
	GroupID             int                `json:"groupId,omitempty"`
	ProductionVersion   VersionSummary     `json:"productionVersion"`
	ProtectedByAPIKey   bool               `json:"protectedByApiKey"`
	StagingVersion      VersionSummary     `json:"stagingVersion"`
	UpdateDate          string             `json:"updateDate"`
	UpdatedBy           string             `json:"updatedBy"`
	VersionNumber       int                `json:"versionNumber"`

	SecurityScheme struct {
		SecuritySchemeType   string `json:"securitySchemeType"`
		SecuritySchemeDetail struct {
			APIKeyLocation string `json:"apiKeyLocation"`
			APIKeyName     string `json:"apiKeyName"`
		} `json:"securitySchemeDetail"`
	} `json:"securityScheme"`
	AkamaiSecurityRestrictions struct {
		MaxJsonxmlElement       int `json:"MAX_JSONXML_ELEMENT"`
		MaxElementNameLength    int `json:"MAX_ELEMENT_NAME_LENGTH"`
		MaxDocDepth             int `json:"MAX_DOC_DEPTH"`
		PositiveSecurityEnabled int `json:"POSITIVE_SECURITY_ENABLED"`
		MaxStringLength         int `json:"MAX_STRING_LENGTH"`
		MaxBodySize             int `json:"MAX_BODY_SIZE"`
		MaxIntegerValue         int `json:"MAX_INTEGER_VALUE"`
	} `json:"akamaiSecurityRestrictions"`
	APIResources Resources `json:"apiResources"`
}

type ListEndpointOptions struct {
	ContractId        string                 `url:"contractId,omitempty"`
	GroupId           string                 `url:"groupId,omitempty"`
	Category          string                 `url:"category,omitempty"`
	Contains          string                 `url:"contains,omitempty"`
	Page              int                    `url:"page,omitempty"`
	PageSize          int                    `url:"pageSize,omitempty"`
	SortBy            SortByValue            `url:"sortBy,omitempty"`
	SortOrder         SortOrderValue         `url:"sortOrder,omitempty"`
	VersionPreference VersionPreferenceValue `url:"versionPreference,omitempty"`
}
type SortByValue string
type SortOrderValue string
type VersionPreferenceValue string

const (
	SortByName       SortByValue = "name"
	SortByUpdateDate SortByValue = "updateDate"

	SortOrderAsc  SortOrderValue = "asc"
	SortOrderDesc SortOrderValue = "desc"

	VersionPreferenceActivatedFirst VersionPreferenceValue = "ACTIVATED_FIRST"
	VersionPreferenceLastUpdated    VersionPreferenceValue = "LAST_UPDATED"
)
