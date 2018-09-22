package apiendpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

func (list *EndpointList) ListEndpoints(options *ListEndpointOptions) error {
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

func CreateEndpoint(options *CreateEndpointOptions) (*Endpoint, error) {
	var req *http.Request
	var err error
	if options.Format == "json" {
		file, err := os.Open(options.ImportFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		var ep Endpoint
		err = json.Unmarshal(bytes, &ep)
		if err != nil {
			return nil, err
		}

		req, err = client.NewJSONRequest(
			Config,
			"POST",
			"/api-definitions/v2/endpoints",
			ep,
		)

		if err != nil {
			return nil, err
		}
	} else {
		req, err = client.NewMultiPartFormDataRequest(
			Config,
			"/api-definitions/v2/endpoints/files",
			options.ImportFile,
			map[string]string{
				"contractId":       options.ContractId,
				"groupId":          options.GroupId,
				"importFileFormat": options.Format,
			},
		)

		if err != nil {
			return nil, err
		}
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	ep := &Endpoint{}
	if err = client.BodyJSON(res, ep); err != nil {
		return nil, err
	}

	return ep, nil
}

type Endpoints []Endpoint

type Endpoint struct {
	APICategoryIds             []int                 `json:"apiCategoryIds,omitempty"`
	APIEndPointHosts           []string              `json:"apiEndPointHosts"`
	APIEndPointID              int                   `json:"apiEndPointId,omitempty"`
	APIEndPointLocked          bool                  `json:"apiEndPointLocked,omitempty"`
	APIEndPointName            string                `json:"apiEndPointName"`
	APIEndPointScheme          string                `json:"apiEndPointScheme,omitempty"`
	APIResourceBaseInfo        []*ResourceBaseInfo   `json:"apiResourceBaseInfo,omitempty"`
	BasePath                   string                `json:"basePath,omitempty"`
	ClonedFromVersion          *int                  `json:"clonedFromVersion,omitempty"`
	ConsumeType                string                `json:"consumeType,omitempty"`
	ContractID                 string                `json:"contractId,omitempty"`
	CreateDate                 string                `json:"createDate,omitempty"`
	CreatedBy                  string                `json:"createdBy,omitempty"`
	Description                string                `json:"description,omitempty"`
	GroupID                    int                   `json:"groupId,omitempty"`
	ProductionVersion          *VersionSummary       `json:"productionVersion,omitempty"`
	ProtectedByAPIKey          bool                  `json:"protectedByApiKey,omitempty"`
	StagingVersion             *VersionSummary       `json:"stagingVersion,omitempty"`
	UpdateDate                 string                `json:"updateDate,omitempty"`
	UpdatedBy                  string                `json:"updatedBy,omitempty"`
	VersionNumber              int                   `json:"versionNumber,omitempty"`
	SecurityScheme             *SecurityScheme       `json:"securityScheme,omitempty"`
	AkamaiSecurityRestrictions *SecurityRestrictions `json:"akamaiSecurityRestrictions,omitempty"`
	APIResources               *Resources            `json:"apiResources,omitempty"`
}

type SecurityScheme struct {
	SecuritySchemeType   string `json:"securitySchemeType,omitempty"`
	SecuritySchemeDetail struct {
		APIKeyLocation string `json:"apiKeyLocation,omitempty"`
		APIKeyName     string `json:"apiKeyName,omitempty"`
	} `json:"securitySchemeDetail,omitempty"`
}

type SecurityRestrictions struct {
	MaxJsonxmlElement       int `json:"MAX_JSONXML_ELEMENT,omitempty"`
	MaxElementNameLength    int `json:"MAX_ELEMENT_NAME_LENGTH,omitempty"`
	MaxDocDepth             int `json:"MAX_DOC_DEPTH,omitempty"`
	PositiveSecurityEnabled int `json:"POSITIVE_SECURITY_ENABLED,omitempty"`
	MaxStringLength         int `json:"MAX_STRING_LENGTH,omitempty"`
	MaxBodySize             int `json:"MAX_BODY_SIZE,omitempty"`
	MaxIntegerValue         int `json:"MAX_INTEGER_VALUE,omitempty"`
}

type ListEndpointOptions struct {
	ContractId        string `url:"contractId,omitempty"`
	GroupId           int    `url:"groupId,omitempty"`
	Category          string `url:"category,omitempty"`
	Contains          string `url:"contains,omitempty"`
	Page              int    `url:"page,omitempty"`
	PageSize          int    `url:"pageSize,omitempty"`
	Show              string `url:show,omitempty`
	SortBy            string `url:"sortBy,omitempty"`
	SortOrder         string `url:"sortOrder,omitempty"`
	VersionPreference string `url:"versionPreference,omitempty"`
}

type CreateEndpointOptions struct {
	ContractId string `url:"contractId,omitempty"`
	GroupId    string `url:"groupId,omitempty"`
	ImportFile string `url:"importFile,omitempty"`
	Format     string
}
