package apidefinitions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	// Endpoint holds configuration for an API
	Endpoint struct {
		APIEndpointID             int64                 `json:"apiEndPointId"`
		APIEndpointName           string                `json:"apiEndPointName"`
		Description               *string               `json:"description"`
		BasePath                  string                `json:"basePath"`
		ConsumeType               *ConsumeType          `json:"consumeType"`
		APIEndpointScheme         *APIEndpointScheme    `json:"apiEndPointScheme"`
		APIEndpointVersion        int64                 `json:"apiEndPointVersion"`
		ContractID                string                `json:"contractId"`
		GroupID                   int64                 `json:"groupId"`
		VersionNumber             int64                 `json:"versionNumber"`
		ClonedFromVersion         *int64                `json:"clonedFromVersion"`
		Locked                    bool                  `json:"locked"`
		StagingVersion            VersionState          `json:"stagingVersion"`
		ProductionVersion         VersionState          `json:"productionVersion"`
		ProtectedByAPIKey         bool                  `json:"protectedByApiKey"`
		APIGatewayEnabled         bool                  `json:"apiGatewayEnabled"`
		CaseSensitive             bool                  `json:"caseSensitive"`
		APIEndpointHosts          []string              `json:"apiEndPointHosts"`
		APICategoryIDs            []int64               `json:"apiCategoryIds"`
		APIResourceBaseInfo       []APIResourceBaseInfo `json:"apiResourceBaseInfo"`
		Source                    *Source               `json:"source"`
		APIVersionInfo            *APIVersionInfo       `json:"apiVersionInfo"`
		PositiveConstrainsEnabled bool                  `json:"positiveConstrainsEnabled"`
		VersionHidden             bool                  `json:"versionHidden"`
		EndpointHidden            bool                  `json:"endpointHidden"`
		IsGraphQL                 bool                  `json:"isGraphQL,omitempty"`
		MatchPathSegmentParam     bool                  `json:"matchPathSegmentParam"`
		AvailableActions          []string              `json:"availableActions"`
		APISource                 *string               `json:"apiSource"`
		LockVersion               int64                 `json:"lockVersion"`
		UpdatedBy                 string                `json:"updatedBy"`
		CreatedBy                 string                `json:"createdBy"`
		CreateDate                string                `json:"createDate"`
		UpdateDate                string                `json:"updateDate"`
	}

	// EndpointDetail holds configuration for an API
	EndpointDetail struct {
		Endpoint
		SecurityScheme             *SecurityScheme             `json:"securityScheme,omitempty"`
		AkamaiSecurityRestrictions *AkamaiSecurityRestrictions `json:"akamaiSecurityRestrictions,omitempty"`
	}

	// EndpointWithResources holds configuration for an API with resources
	EndpointWithResources struct {
		EndpointDetail
		APIResources     []APIResource `json:"apiResources"`
		DiscoveredPIIIDs []int64       `json:"discoveredPiiIds"`
		ProductionStatus *string       `json:"productionStatus"`
		StagingStatus    *string       `json:"stagingStatus"`
	}

	// APISourceDiff contains source diff
	APISourceDiff struct {
		Name        string `json:"name"`
		SourceValue string `json:"sourceValue"`
		SavedValue  string `json:"savedValue"`
	}

	// RegisterEndpointRequest holds the RegisterEndpoint request parameters
	RegisterEndpointRequest struct {
		SecurityScheme             *SecurityScheme             `json:"securityScheme,omitempty"`
		AkamaiSecurityRestrictions *AkamaiSecurityRestrictions `json:"akamaiSecurityRestrictions,omitempty"`
		VersionNumber              int64                       `json:"versionNumber,omitempty"`
		APIEndpointName            string                      `json:"apiEndPointName"`
		Description                string                      `json:"description,omitempty"`
		DiscoveredPIIIDs           []int64                     `json:"discoveredPiiIds,omitempty"`
		BasePath                   string                      `json:"basePath,omitempty"`
		APIEndpointScheme          APIEndpointScheme           `json:"apiEndPointScheme,omitempty"`
		ConsumeType                ConsumeType                 `json:"consumeType,omitempty"`
		APIEndpointHosts           []string                    `json:"apiEndPointHosts"`
		APICategoryIDs             []int64                     `json:"apiCategoryIds,omitempty"`
		CaseSensitive              *bool                       `json:"caseSensitive,omitempty"`
		MatchPathSegmentParam      *bool                       `json:"matchPathSegmentParam,omitempty"`
		IsGraphQL                  *bool                       `json:"isGraphQL,omitempty"`
		APIResources               []APIResource               `json:"apiResources,omitempty"`
		APIVersionInfo             *APIVersionInfo             `json:"apiVersionInfo,omitempty"`
		APISource                  APISource                   `json:"apiSource,omitempty"`
		APIGatewayEnabled          *bool                       `json:"apiGatewayEnabled,omitempty"`
		ContractID                 string                      `json:"contractId"`
		GroupID                    int64                       `json:"groupId"`
	}

	// ImportFileFormat format of the file, either raml or swagger
	ImportFileFormat string

	// ImportFileSource location of the file, either URL if you store the file on the web, or BODY_BASE64 if you encode the file contents in the request body
	ImportFileSource string

	// RegisterEndpointFromFileRequest holds the RegisterEndpointFromFile request parameters
	RegisterEndpointFromFileRequest struct {
		ContractID        string           `json:"contractId"`
		GroupID           int64            `json:"groupId"`
		ImportFileContent *string          `json:"importFileContent,omitempty"`
		ImportFileFormat  ImportFileFormat `json:"importFileFormat"`
		ImportFileSource  ImportFileSource `json:"importFileSource"`
		ImportURL         *string          `json:"importUrl,omitempty"`
		Root              *string          `json:"root,omitempty"`
	}

	// ShowEndpointRequest contains parameters used to show an endpoint
	ShowEndpointRequest struct {
		APIEndpointID int64
	}

	// ShowEndpointResponse represents a response object returned by ShowEndpoint
	ShowEndpointResponse struct {
		EndpointWithResources
	}

	// HideEndpointRequest contains parameters used to hide an endpoint
	HideEndpointRequest struct {
		APIEndpointID int64
	}

	// HideEndpointResponse represents a response object returned by HideEndpoint
	HideEndpointResponse struct {
		EndpointWithResources
	}

	// GetEndpointRequest contains parameters used to get an endpoint
	GetEndpointRequest struct {
		APIEndpointID int64
	}

	// GetEndpointResponse represents a response object returned by GetEndpoint
	GetEndpointResponse EndpointDetail

	// DeleteEndpointRequest represents a response object returned by DeleteEndpoint
	DeleteEndpointRequest struct {
		APIEndpointID int64
	}

	// EndpointResponse represents a response endpoint object
	EndpointResponse struct {
		ContractID                string             `json:"contractId"`
		GroupID                   int64              `json:"groupId"`
		APIEndpointID             int64              `json:"apiEndPointId"`
		APIEndpointVersion        int64              `json:"apiEndPointVersion"`
		VersionNumber             int64              `json:"versionNumber"`
		APIEndpointName           string             `json:"apiEndPointName"`
		Description               *string            `json:"description"`
		BasePath                  string             `json:"basePath"`
		ClonedFromVersion         *int64             `json:"clonedFromVersion"`
		APIEndpointLocked         bool               `json:"apiEndPointLocked"`
		APIEndpointScheme         *APIEndpointScheme `json:"apiEndPointScheme"`
		ConsumeType               *ConsumeType       `json:"consumeType"`
		APIEndpointHosts          []string           `json:"apiEndPointHosts"`
		APICategoryIDs            []int64            `json:"apiCategoryIds"`
		LockVersion               int64              `json:"lockVersion"`
		UpdatedBy                 string             `json:"updatedBy"`
		CreatedBy                 string             `json:"createdBy"`
		CreateDate                string             `json:"createDate"`
		UpdateDate                string             `json:"updateDate"`
		PositiveConstrainsEnabled bool               `json:"positiveConstrainsEnabled"`
		CaseSensitive             bool               `json:"caseSensitive"`
		MatchPathSegmentParam     bool               `json:"matchPathSegmentParam"`
		Source                    *Source            `json:"source"`
		StagingVersion            VersionState       `json:"stagingVersion"`
		ProductionVersion         VersionState       `json:"productionVersion"`
		ProtectedByAPIKey         bool               `json:"protectedByApiKey"`
		IsGraphQL                 bool               `json:"isGraphQL,omitempty"`
		AvailableActions          []string           `json:"availableActions"`
		VersionHidden             bool               `json:"versionHidden"`
		EndpointHidden            bool               `json:"endpointHidden"`
		APISource                 *string            `json:"apiSource"`
		APIGatewayEnabled         bool               `json:"apiGatewayEnabled"`
		APIVersionInfo            *APIVersionInfo    `json:"apiVersionInfo"`
	}

	// ListEndpointsRequest holds the ListEndpoints request parameters
	ListEndpointsRequest struct {
		PIIOnly           bool
		Page              int64
		PageSize          int64
		Category          string
		Contains          string
		SortBy            ListEndpointSortType
		SortOrder         SortOrderType
		VersionPreference VersionPreference
		Show              Visibility
		ContractID        string
		GroupID           int64
	}

	// SecurityScheme contains information about the key with which users may access the endpoint
	SecurityScheme struct {
		SecuritySchemeType   string               `json:"securitySchemeType"`
		SecuritySchemeDetail SecuritySchemeDetail `json:"securitySchemeDetail"`
	}

	// SecuritySchemeDetail contains information about the location of the API key
	SecuritySchemeDetail struct {
		APIKeyName     string         `json:"apiKeyName,omitempty"`
		APIKeyLocation APIKeyLocation `json:"apiKeyLocation"`
	}

	// RegisterEndpointResponse holds the response from RegisterEndpoint
	RegisterEndpointResponse struct {
		EndpointWithResources
	}

	// RegisterEndpointFromFileResponse holds the response from RegisterEndpointFromFile
	RegisterEndpointFromFileResponse RegisterEndpointResponse

	// ListEndpointsResponse holds the ListEndpoints response data
	ListEndpointsResponse struct {
		TotalSize    int64      `json:"totalSize"`
		Page         int64      `json:"page"`
		PageSize     int64      `json:"pageSize"`
		APIEndpoints []Endpoint `json:"apiEndPoints"`
	}

	// APIResourceBaseInfo holds basic endpoint information
	APIResourceBaseInfo struct {
		CreatedBy               *string `json:"createdBy,omitempty"`
		CreateDate              *string `json:"createDate,omitempty"`
		UpdateDate              *string `json:"updateDate,omitempty"`
		UpdatedBy               *string `json:"updatedBy,omitempty"`
		LockVersion             int64   `json:"lockVersion"`
		APIResourceID           int64   `json:"apiResourceId"`
		APIResourceName         string  `json:"apiResourceName"`
		ResourcePath            string  `json:"resourcePath"`
		Description             *string `json:"description,omitempty"`
		Link                    *string `json:"link,omitempty"`
		APIResourceClonedFromID int64   `json:"apiResourceClonedFromId"`
		APIResourceLogicID      int64   `json:"apiResourceLogicId"`
		Private                 bool    `json:"private"`
	}

	// APIMethod represents method for the APIResourceBase
	APIMethod struct {
		APIResourceMethodID      int64        `json:"apiResourceMethodId"`
		APIResourceMethodLogicID int64        `json:"apiResourceMethodLogicId"`
		APIResourceMethod        string       `json:"apiResourceMethod"`
		IsPrivate                bool         `json:"isPrivate"`
		StagingVersion           VersionState `json:"stagingVersion"`
		ProductionVersion        VersionState `json:"productionVersion"`
	}

	//SourceType specifies if the API's data comes from Swagger or Raml file
	SourceType string

	// Source contains information about an endpoint source
	Source struct {
		Type                 SourceType `json:"type"`
		APIVersion           string     `json:"apiVersion"`
		SpecificationVersion string     `json:"specificationVersion"`
	}

	// VersionState contains information about activation status on given network
	VersionState struct {
		VersionNumber *int64            `json:"versionNumber"`
		Status        *ActivationStatus `json:"status"`
		Timestamp     *string           `json:"timestamp"`
		LastError     *EndpointError    `json:"lastError"`
	}

	// EndpointError contains information about an endpoint error
	EndpointError struct {
		Timestamp     string `json:"timestamp"`
		Status        string `json:"status"`
		Type          string `json:"type"`
		VersionNumber int64  `json:"versionNumber"`
	}

	// AkamaiSecurityRestrictions contains information about the Kona Site Defender security restrictions that you apply to an API
	AkamaiSecurityRestrictions struct {
		MaxJSONXMLElement             *int64            `json:"MAX_JSONXML_ELEMENT,omitempty"`
		MaxElementNameLength          *int64            `json:"MAX_ELEMENT_NAME_LENGTH,omitempty"`
		MaxStringLength               *int64            `json:"MAX_STRING_LENGTH,omitempty"`
		MaxIntegerValue               *int64            `json:"MAX_INTEGER_VALUE,omitempty"`
		MaxDocDepth                   *int64            `json:"MAX_DOC_DEPTH,omitempty"`
		MaxBodySize                   *int64            `json:"MAX_BODY_SIZE,omitempty"`
		PositiveSecurityVersion       *int64            `json:"POSITIVE_SECURITY_VERSION,omitempty"`
		PositiveSecurityEnabled       *restrictionsBool `json:"POSITIVE_SECURITY_ENABLED,omitempty"`
		AllowUndefinedResources       *restrictionsBool `json:"ALLOW_UNDEFINED_RESOURCES,omitempty"`
		AllowOnlySpecUndefinedMethods *restrictionsBool `json:"ALLOW_ONLY_SPEC_UNDEFINED_METHODS,omitempty"`
		AllowUndefinedMethodGet       *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_GET,omitempty"`
		AllowUndefinedMethodPost      *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_POST,omitempty"`
		AllowUndefinedMethodPut       *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_PUT,omitempty"`
		AllowUndefinedMethodDelete    *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_DELETE,omitempty"`
		AllowUndefinedMethodHead      *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_HEAD,omitempty"`
		AllowUndefinedMethodOptions   *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_OPTIONS,omitempty"`
		AllowUndefinedMethodPatch     *restrictionsBool `json:"ALLOW_UNDEFINED_METHOD_PATCH,omitempty"`
		AllowUndefinedParams          *restrictionsBool `json:"ALLOW_UNDEFINED_PARAMS,omitempty"`
		AllowUndefinedSpecParams      *restrictionsBool `json:"ALLOW_UNDEFINED_SPEC_PARAMS,omitempty"`
		AllowUndefinedCookieParams    *restrictionsBool `json:"ALLOW_UNDEFINED_COOKIE_PARAMS,omitempty"`
		AllowUndefinedQueryParams     *restrictionsBool `json:"ALLOW_UNDEFINED_QUERY_PARAMS,omitempty"`
		AllowUndefinedBodyParams      *restrictionsBool `json:"ALLOW_UNDEFINED_BODY_PARAMS,omitempty"`
		AllowUndefinedHeaderParams    *restrictionsBool `json:"ALLOW_UNDEFINED_HEADER_PARAMS,omitempty"`
	}

	restrictionsBool bool

	// APIResource contains information about the API resource
	APIResource struct {
		APIResourceClonedFromID    *int64              `json:"apiResourceClonedFromId,omitempty"`
		APIResourceID              *int64              `json:"apiResourceId,omitempty"`
		APIResourceLogicID         *int64              `json:"apiResourceLogicId,omitempty"`
		APIResourceMethodNameLists []string            `json:"apiResourceMethodNameLists,omitempty"`
		APIResourceMethods         []APIResourceMethod `json:"apiResourceMethods,omitempty"`
		APIResourceName            string              `json:"apiResourceName"`
		CreateDate                 string              `json:"createDate,omitempty"`
		CreatedBy                  string              `json:"createdBy,omitempty"`
		Description                string              `json:"description,omitempty"`
		Link                       *string             `json:"link,omitempty"`
		LockVersion                *int64              `json:"lockVersion,omitempty"`
		Private                    *bool               `json:"private,omitempty"`
		ResourcePath               string              `json:"resourcePath"`
		UpdateDate                 string              `json:"updateDate,omitempty"`
		UpdatedBy                  string              `json:"updatedBy,omitempty"`
	}

	// APIResourceRes contains information about the API resource
	APIResourceRes struct {
		APIResourceClonedFromID    int64                  `json:"apiResourceClonedFromId"`
		APIResourceID              int64                  `json:"apiResourceId"`
		APIResourceLogicID         int64                  `json:"apiResourceLogicId"`
		APIResourceMethodNameLists []string               `json:"apiResourceMethodNameLists"`
		APIResourceMethodsRes      []APIResourceMethodRes `json:"apiResourceMethods"`
	}

	// APIParameterRes holds information about an API method parameter
	APIParameterRes struct {
		APIParameterID          int64                   `json:"apiParameterId"`
		APIParameterName        string                  `json:"apiParameterName"`
		APIParameterRequired    bool                    `json:"apiParameterRequired"`
		APIParameterLocation    APIParameterLocation    `json:"apiParameterLocation"`
		PathParamLocationID     int64                   `json:"pathParamLocationId"`
		APIParameterType        APIParameterType        `json:"apiParameterType"`
		Array                   bool                    `json:"array"`
		APIParamLogicID         int64                   `json:"apiParamLogicId"`
		APIResourceMethParamID  int64                   `json:"apiResourceMethParamId"`
		APIParameterNotes       string                  `json:"apiParameterNotes"`
		APIParameterRestriction APIParameterRestriction `json:"apiParameterRestriction"`
		APIChildParameters      []APIParameterRes       `json:"apiChildParameters"`
	}

	// APIResourceMethod holds configuration for an API resource method
	APIResourceMethod struct {
		APIResourceMethodID      *int64              `json:"apiResourceMethodId,omitempty"`
		APIResourceMethodLogicID *int64              `json:"apiResourceMethodLogicId,omitempty"`
		APIResourceMethod        APIResourceMethods  `json:"apiResourceMethod"`
		APIParameters            []APIParameter      `json:"apiParameters,omitempty"`
		MethodRestrictions       *MethodRestrictions `json:"methodRestrictions,omitempty"`
	}

	// APIResourceMethodRes holds configuration for an API resource method
	APIResourceMethodRes struct {
		APIResourceMethodID      int64              `json:"apiResourceMethodId"`
		APIResourceMethod        APIResourceMethods `json:"apiResourceMethod"`
		APIResourceMethodLogicID int64              `json:"apiResourceMethodLogicId"`
		APIParameters            []APIParameterRes  `json:"apiParameters"`
		APIResourceName          string             `json:"apiResourceName"`
		CreateDate               string             `json:"createDate"`
		CreateBy                 string             `json:"createdBy"`
		Description              string             `json:"description"`
		Link                     string             `json:"link"`
		LockVersion              int64              `json:"lockVersion"`
		Private                  bool               `json:"private"`
		ResourcePath             string             `json:"resourcePath"`
		UpdateDate               string             `json:"updateDate"`
		UpdatedBy                string             `json:"updatedBy"`
	}

	// MethodRestrictions holds configuration for an API method restrictions
	MethodRestrictions struct {
		AllowMethodUndefinedParameters *AllowMethodUndefinedParameters `json:"allowMethodUndefinedParameters,omitempty"`
	}

	// AllowMethodUndefinedParameters holds configuration for an API method allowed undefined parameters
	AllowMethodUndefinedParameters struct {
		Cookie bool `json:"cookie"`
		Header bool `json:"header"`
		Body   bool `json:"body"`
		Query  bool `json:"query"`
	}

	// APIParameter holds information about an API method parameter
	APIParameter struct {
		APIParameterName        string                   `json:"apiParameterName"`
		APIParameterRequired    bool                     `json:"apiParameterRequired"`
		APIParameterLocation    APIParameterLocation     `json:"apiParameterLocation"`
		PathParamLocationID     *int64                   `json:"pathParamLocationId,omitempty"`
		APIParameterType        APIParameterType         `json:"apiParameterType"`
		Array                   *bool                    `json:"array,omitempty"`
		APIParameterNotes       *string                  `json:"apiParameterNotes,omitempty"`
		APIParameterRestriction *APIParameterRestriction `json:"apiParameterRestriction,omitempty"`
		APIChildParameters      []APIParameter           `json:"apiChildParameters"`
		APIParameterID          *int64                   `json:"apiParameterId,omitempty"`
		APIParamLogicID         *int64                   `json:"apiParamLogicId,omitempty"`
		APIResourceMethParamID  *int64                   `json:"apiResourceMethParamId,omitempty"`
	}

	// APIVersionInfoLocation the location of the API version value in an incoming request. Either HEADER, BASE_PATH, or QUERY parameter.
	APIVersionInfoLocation string

	// APIVersionInfo contains information about a major API version
	APIVersionInfo struct {
		Location      APIVersionInfoLocation `json:"location"`
		ParameterName string                 `json:"parameterName"`
		Value         string                 `json:"value"`
	}

	// APIParameterRestriction contains information about restrictions and XML representation rules specified for the parameter
	APIParameterRestriction struct {
		LengthRestriction      *LengthRestriction      `json:"lengthRestriction,omitempty"`
		RangeRestriction       *RangeRestriction       `json:"rangeRestriction,omitempty"`
		NumberRangeRestriction *NumberRangeRestriction `json:"numberRangeRestriction,omitempty"`
		ArrayRestriction       *ArrayRestriction       `json:"arrayRestriction,omitempty"`
		XMLConversionRule      *XMLConversionRule      `json:"xmlConversionRule,omitempty"`
		ResponseRestriction    *ResponseRestriction    `json:"responseRestriction,omitempty"`
	}

	// XMLConversionRule contains information about an XML representation of a JSON-encoded parameter
	XMLConversionRule struct {
		Attribute bool   `json:"attribute"`
		Wrapped   bool   `json:"wrapped"`
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Prefix    string `json:"prefix,omitempty"`
	}

	// LengthRestriction contains information about length restrictions for string type parameters
	LengthRestriction struct {
		LengthMax int64 `json:"lengthMax"`
		LengthMin int64 `json:"lengthMin"`
	}

	// RangeRestriction contains information about range restrictions for integer type parameters
	RangeRestriction struct {
		RangeMin int64 `json:"rangeMin"`
		RangeMax int64 `json:"rangeMax"`
	}

	// NumberRangeRestriction contains information about range restrictions for number type parameters
	NumberRangeRestriction struct {
		NumberRangeMin float64 `json:"numberRangeMin"`
		NumberRangeMax float64 `json:"numberRangeMax"`
	}

	// ArrayRestriction contains information about array restrictions for array type parameters
	ArrayRestriction struct {
		MaxItems int64 `json:"maxItems"`
		MinItems int64 `json:"minItems"`
	}

	// ResponseRestriction contains information about response restrictions
	ResponseRestriction struct {
		StatusCodes []int64     `json:"statusCodes"`
		MaxBodySize MaxBodySize `json:"maxBodySize,omitempty"`
	}

	// ListUserEntitlementsResponse contains response from ListUserEntitlements request
	ListUserEntitlementsResponse []string

	// MaxBodySize represents MaxBodySize value
	MaxBodySize string

	// APIResourceMethods represents APIResourceMethods value
	APIResourceMethods string

	// APIKeyLocation represents APIKeyLocation value
	APIKeyLocation string

	// APIParameterLocation represents APIParameterLocation value
	APIParameterLocation string

	// APIParameterType represents APIParameterType value
	APIParameterType string

	// VersionPreference represents VersionPreference query param in ListEndpoint
	VersionPreference string

	// ListEndpointSortType represents the type of the sorting is based on
	ListEndpointSortType string

	// ConsumeType content type the endpoint exchanges
	ConsumeType string

	// APIEndpointScheme the URL scheme to which the endpoint may respond, either http, https, or http/https for both
	APIEndpointScheme string

	//APISource Specifies if the endpoint data comes from API_DISCOVERY or was provided by the USER
	APISource string
)

const (
	// APIResourceMethodsGET represents HTTP GET method
	APIResourceMethodsGET APIResourceMethods = "GET"
	// APIResourceMethodsPUT represents HTTP PUT method
	APIResourceMethodsPUT APIResourceMethods = "PUT"
	// APIResourceMethodsPOST represents HTTP POST method
	APIResourceMethodsPOST APIResourceMethods = "POST"
	// APIResourceMethodsDELETE represents HTTP DELETE method
	APIResourceMethodsDELETE APIResourceMethods = "DELETE"
	// APIResourceMethodsHEAD represents HTTP HEAD method
	APIResourceMethodsHEAD APIResourceMethods = "HEAD"
	// APIResourceMethodsPATCH represents HTTP PATCH method
	APIResourceMethodsPATCH APIResourceMethods = "PATCH"
	// APIResourceMethodsOPTIONS represents HTTP OPTIONS  method
	APIResourceMethodsOPTIONS APIResourceMethods = "OPTIONS"

	// MaxBodySizeSize6k represents MaxBodySize of value "SIZE_6K"
	MaxBodySizeSize6k MaxBodySize = "SIZE_6K"

	// MaxBodySizeSize8k represents MaxBodySize of value "SIZE_8K"
	MaxBodySizeSize8k MaxBodySize = "SIZE_8K"

	// MaxBodySizeSize12k represents MaxBodySize of value "SIZE_12K"
	MaxBodySizeSize12k MaxBodySize = "SIZE_12K"

	// MaxBodySizeSize16k represents MaxBodySize of value "SIZE_16K"
	MaxBodySizeSize16k MaxBodySize = "SIZE_16K"

	// MaxBodySizeNoLimit represents MaxBodySize of value "NO_LIMIT"
	MaxBodySizeNoLimit MaxBodySize = "NO_LIMIT"

	// MaxBodySizeNull represents MaxBodySize as null
	MaxBodySizeNull MaxBodySize = ""

	// APIParameterLocationQuery represents APIParameterLocation for "query"
	APIParameterLocationQuery APIParameterLocation = "query"

	// APIParameterLocationCookie represents APIParameterLocation for "cookie"
	APIParameterLocationCookie APIParameterLocation = "cookie"

	// APIParameterLocationHeader represents APIParameterLocation for "header"
	APIParameterLocationHeader APIParameterLocation = "header"

	// APIParameterLocationPath represents APIParameterLocation for "path"
	APIParameterLocationPath APIParameterLocation = "path"

	// APIParameterLocationBody represents APIParameterLocation for "body"
	APIParameterLocationBody APIParameterLocation = "body"

	// APIParameterTypeString represents APIParameterLocation for "string"
	APIParameterTypeString APIParameterType = "string"
	// APIParameterTypeInteger represents APIParameterLocation for "integer"
	APIParameterTypeInteger APIParameterType = "integer"
	// APIParameterTypeNumber represents APIParameterLocation for "number"
	APIParameterTypeNumber APIParameterType = "number"
	// APIParameterTypeBoolean represents APIParameterLocation for "boolean"
	APIParameterTypeBoolean APIParameterType = "boolean"
	// APIParameterTypeJSONXML represents APIParameterLocation for "json/xml"
	APIParameterTypeJSONXML APIParameterType = "json/xml"

	// APIKeyLocationCookie represents APIKeyLocationQuery for "cookie"
	APIKeyLocationCookie APIKeyLocation = "cookie"
	// APIKeyLocationHeader represents APIKeyLocationQuery for "header"
	APIKeyLocationHeader APIKeyLocation = "header"
	// APIKeyLocationQuery represents APIKeyLocationQuery for "query"
	APIKeyLocationQuery APIKeyLocation = "query"

	// ConsumeTypeJSON holds value for consume type json
	ConsumeTypeJSON ConsumeType = "json"
	// ConsumeTypeXML holds value for consume type xml
	ConsumeTypeXML ConsumeType = "xml"
	// ConsumeTypeJSONXML holds value for consume type json/xml
	ConsumeTypeJSONXML ConsumeType = "json/xml"
	// ConsumeTypeUrlencoded holds value for consume type urlencoded
	ConsumeTypeUrlencoded ConsumeType = "urlencoded"
	// ConsumeTypeJSONUrlencoded holds value for consume type json/urlencoded
	ConsumeTypeJSONUrlencoded ConsumeType = "json/urlencoded"
	// ConsumeTypeXMLUrlencoded holds value for consume type xml/urlencoded
	ConsumeTypeXMLUrlencoded ConsumeType = "xml/urlencoded"
	// ConsumeTypeJSONXMLUrlencoded holds value for consume type json/xml/urlencoded
	ConsumeTypeJSONXMLUrlencoded ConsumeType = "json/xml/urlencoded"
	// ConsumeTypeAny holds value for consume type any
	ConsumeTypeAny ConsumeType = "any"
	// ConsumeTypeNone holds value for consume type none
	ConsumeTypeNone ConsumeType = "none"

	// ImportFileFormatRaml holds value for raml file format
	ImportFileFormatRaml ImportFileFormat = "raml"
	// ImportFileFormatSwagger holds value for swagger file format
	ImportFileFormatSwagger ImportFileFormat = "swagger"

	// ImportFileSourceURL holds value for raml file format
	ImportFileSourceURL ImportFileSource = "URL"
	// ImportFileSourceBase64 holds value for swagger file format
	ImportFileSourceBase64 ImportFileSource = "BODY_BASE64"

	// APIEndpointSchemeHTTP holds value for http scheme
	APIEndpointSchemeHTTP APIEndpointScheme = "http"
	// APIEndpointSchemeHTTPS holds value for https scheme
	APIEndpointSchemeHTTPS APIEndpointScheme = "https"
	// APIEndpointSchemeHTTPHTTPS holds value for http/https scheme
	APIEndpointSchemeHTTPHTTPS APIEndpointScheme = "http/https"

	// APISourceUser holds value for USER APISource
	APISourceUser APISource = "USER"
	// APISourceAPIDiscovery holds value for API_DISCOVERY APISource
	APISourceAPIDiscovery APISource = "API_DISCOVERY"

	// SourceTypeSwagger holds value for SWAGGER SourceType
	SourceTypeSwagger SourceType = "SWAGGER"
	// SourceTypeRaml holds value for SWAGGER SourceType
	SourceTypeRaml SourceType = "RAML"

	// NameSort represents "name" SortType query for ListEndpoint
	NameSort ListEndpointSortType = "name"
	// UpdateEndpointDateSort represents "updateDate" SortType query for ListEndpoint
	UpdateEndpointDateSort ListEndpointSortType = "updateDate"

	// VersionPreferenceLastUpdated represents "LAST_UPDATED" VersionPreference query for ListEndpoint
	VersionPreferenceLastUpdated VersionPreference = "LAST_UPDATED"
	// VersionPreferenceActivatedFirst represents "ACTIVATED_FIRST" VersionPreference query for ListEndpoint
	VersionPreferenceActivatedFirst VersionPreference = "ACTIVATED_FIRST"

	// APIVersionLocationHeader represents "HEADER" APIVersionInfoLocation
	APIVersionLocationHeader APIVersionInfoLocation = "HEADER"
	// APIVersionLocationBasePath represents "BASE_PATH" APIVersionInfoLocation
	APIVersionLocationBasePath APIVersionInfoLocation = "BASE_PATH"
	// APIVersionLocationQuery represents "QUERY" APIVersionInfoLocation
	APIVersionLocationQuery APIVersionInfoLocation = "QUERY"
)

// Validate validates GetEndpointRequest
func (r GetEndpointRequest) Validate() interface{} {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
	})
}

// Validate validates ShowEndpointRequest
func (r ShowEndpointRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
	})
}

// Validate validates HideEndpointRequest
func (r HideEndpointRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
	})
}

// Validate validates DeleteEndpointRequest
func (r DeleteEndpointRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
	})
}

// Validate validates ListEndpointSortType
func (s ListEndpointSortType) Validate() error {
	return validation.In(NameSort, UpdateEndpointDateSort).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of: '%s' or '%s'", s, NameSort, UpdateEndpointDateSort)).
		Validate(s)
}

// Validate validates ListEndpointSortType
func (v VersionPreference) Validate() error {
	return validation.In(VersionPreferenceLastUpdated, VersionPreferenceActivatedFirst).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of: '%s', '%s' or '' (empty)", v, VersionPreferenceLastUpdated, VersionPreferenceActivatedFirst)).
		Validate(v)
}

// Validate validates ListEndpointsRequest
func (r ListEndpointsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"SortBy":            validation.Validate(r.SortBy),
		"SortOrder":         validation.Validate(r.SortOrder),
		"VersionPreference": validation.Validate(r.VersionPreference),
		"Show":              validation.Validate(r.Show),
	})
}

// Validate validates SecurityScheme
func (s SecurityScheme) Validate() error {
	return validation.Errors{
		"SecuritySchemeType":   validation.Validate(s.SecuritySchemeType, validation.Required, validation.In("apikey").Error(fmt.Sprintf(`value '%s' is not valid. Must be: 'apikey'`, s.SecuritySchemeType))),
		"SecuritySchemeDetail": validation.Validate(s.SecuritySchemeDetail, validation.Required),
	}.Filter()
}

// Validate validates SecurityScheme
func (l APIKeyLocation) Validate() error {
	return validation.In(APIKeyLocationCookie, APIKeyLocationHeader, APIKeyLocationQuery).
		Error(fmt.Sprintf(`value '%s' is not valid. Must be one of: '%s', '%s', '%s'`, l, APIKeyLocationCookie, APIKeyLocationHeader, APIKeyLocationQuery)).
		Validate(l)
}

// Validate validate SecuritySchemeDetail
func (s SecuritySchemeDetail) Validate() error {
	return validation.Errors{
		"APIKeyLocation": validation.Validate(s.APIKeyLocation, validation.Required),
		"APIKeyName":     validation.Validate(s.APIKeyName, validation.Required),
	}.Filter()
}

// Validate validates APIParameterRestriction
func (r APIParameterRestriction) Validate() error {
	return validation.Errors{
		"LengthRestriction":   validation.Validate(r.LengthRestriction),
		"ArrayRestriction":    validation.Validate(r.ArrayRestriction),
		"ResponseRestriction": validation.Validate(r.ResponseRestriction),
	}.Filter()
}

// Validate validates LengthRestriction
func (r LengthRestriction) Validate() error {
	return validation.Errors{
		"LengthMax": validation.Validate(r.LengthMax, validation.Min(0)),
		"LengthMin": validation.Validate(r.LengthMin, validation.Min(0)),
	}.Filter()
}

// Validate validates ResponseRestriction
func (m MaxBodySize) Validate() error {
	return validation.In(MaxBodySizeSize6k, MaxBodySizeSize8k, MaxBodySizeSize12k, MaxBodySizeSize16k, MaxBodySizeNoLimit, MaxBodySizeNull).
		Error(fmt.Sprintf(`value '%s' is not valid. Must be one of: '%s', '%s', '%s', '%s', '%s' or %s`, m, MaxBodySizeSize6k, MaxBodySizeSize8k, MaxBodySizeSize12k, MaxBodySizeSize16k, MaxBodySizeNoLimit, MaxBodySizeNull)).
		Validate(m)
}

// Validate validates ResponseRestriction
func (r ResponseRestriction) Validate() error {
	return validation.Errors{
		"MaxBodySize": validation.Validate(r.MaxBodySize),
	}.Filter()
}

// Validate validates ResponseRestriction
func (l APIParameterLocation) Validate() error {
	return validation.In(APIParameterLocationQuery, APIParameterLocationHeader, APIParameterLocationPath, APIParameterLocationCookie, APIParameterLocationBody).
		Error(fmt.Sprintf(`value '%s' is not valid. Must be one of: '%s', '%s', '%s', '%s', '%s'`, l, APIParameterLocationQuery, APIParameterLocationHeader, APIParameterLocationPath, APIParameterLocationCookie, APIParameterLocationBody)).
		Validate(l)
}

// Validate validates ResponseRestriction
func (t APIParameterType) Validate() error {
	return validation.In(APIParameterTypeString, APIParameterTypeInteger, APIParameterTypeNumber, APIParameterTypeBoolean, APIParameterTypeJSONXML).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of '%s', '%s', '%s', '%s', or '%s'", t, APIParameterTypeString, APIParameterTypeInteger, APIParameterTypeNumber, APIParameterTypeBoolean, APIParameterTypeJSONXML)).
		Validate(t)
}

// Validate validates APIParameter
func (p APIParameter) Validate() error {
	return validation.Errors{
		"APIParameterName":        validation.Validate(p.APIParameterName, validation.Required),
		"APIParameterLocation":    validation.Validate(p.APIParameterLocation, validation.Required),
		"APIParameterType":        validation.Validate(p.APIParameterType, validation.Required),
		"APIParameterRestriction": validation.Validate(p.APIParameterRestriction),
	}.Filter()
}

// Validate validates ResponseRestriction
func (m APIResourceMethods) Validate() error {
	return validation.In(APIResourceMethodsGET, APIResourceMethodsPUT, APIResourceMethodsPOST, APIResourceMethodsDELETE, APIResourceMethodsHEAD, APIResourceMethodsPATCH, APIResourceMethodsOPTIONS).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s', '%s'", m, APIResourceMethodsGET, APIResourceMethodsPUT, APIResourceMethodsPOST, APIResourceMethodsDELETE, APIResourceMethodsHEAD, APIResourceMethodsPATCH, APIResourceMethodsOPTIONS)).
		Validate(m)
}

// Validate validates APIResourceMethod
func (a APIResourceMethod) Validate() error {
	return validation.Errors{
		"APIResourceMethod": validation.Validate(a.APIResourceMethod, validation.Required),
		"APIParameters":     validation.Validate(a.APIParameters),
	}.Filter()
}

// Validate validates APIResource
func (a APIResource) Validate() error {
	return validation.Errors{
		"APIResourceName":    validation.Validate(a.APIResourceName, validation.Required),
		"ResourcePath":       validation.Validate(a.ResourcePath, validation.Required),
		"APIResourceMethods": validation.Validate(a.APIResourceMethods),
	}.Filter()
}

// Validate validates AkamaiSecurityRestrictions
func (r AkamaiSecurityRestrictions) Validate() error {
	return validation.Errors{
		"POSITIVE_SECURITY_VERSION": validation.Validate(r.PositiveSecurityVersion, validation.In(int64(1), int64(2)).Error(fmt.Sprintf("value %d is not valid. Must be one of 1, 2", r.PositiveSecurityVersion))),
	}.Filter()
}

// Validate validates APIVersionInfoLocation
func (l APIVersionInfoLocation) Validate() error {
	return validation.In(APIVersionLocationHeader, APIVersionLocationBasePath, APIVersionLocationQuery).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of '%s', '%s', '%s'", l, APIVersionLocationHeader, APIVersionLocationBasePath, APIVersionLocationQuery)).
		Validate(l)
}

// Validate validates APIVersionInfo
func (a APIVersionInfo) Validate() error {
	return validation.Errors{
		"Location": validation.Validate(a.Location, validation.Required),
	}.Filter()
}

// Validate validates RegisterEndpointRequest
func (r RegisterEndpointRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointName":            validation.Validate(r.APIEndpointName, validation.Required),
		"BasePath":                   validation.Validate(r.BasePath, validation.NewStringRuleWithError(func(s string) bool { return len(s) == 0 || !strings.HasSuffix(s, "/") }, validation.NewError("basePath", "basePath should not end with `/`"))),
		"APIEndpointScheme":          validation.Validate(r.APIEndpointScheme),
		"ConsumeType":                validation.Validate(r.ConsumeType),
		"APIEndpointHosts":           validation.Validate(r.APIEndpointHosts, validation.Required),
		"APISource":                  validation.Validate(r.APISource),
		"APIResources":               validation.Validate(r.APIResources),
		"APIVersionInfo":             validation.Validate(r.APIVersionInfo, validation.NilOrNotEmpty),
		"AkamaiSecurityRestrictions": validation.Validate(r.AkamaiSecurityRestrictions, validation.NilOrNotEmpty),
		"ContractID":                 validation.Validate(r.ContractID, validation.Required),
		"GroupID":                    validation.Validate(r.GroupID, validation.Required),
	})
}

// Validate validates ImportFileFormat
func (f ImportFileFormat) Validate() error {
	return validation.In(ImportFileFormatSwagger, ImportFileFormatRaml).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of: '%s', '%s'", f, ImportFileFormatSwagger, ImportFileFormatRaml)).
		Validate(f)
}

// Validate validates ImportFileSource
func (f ImportFileSource) Validate() error {
	return validation.In(ImportFileSourceURL, ImportFileSourceBase64).
		Error(fmt.Sprintf("value '%s' is not valid. Must be one of: '%s', '%s'", f, ImportFileSourceURL, ImportFileSourceBase64)).
		Validate(f)
}

// Validate validates RegisterEndpointFromFileRequest
func (r *RegisterEndpointFromFileRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ImportFileContent": validation.Validate(r.ImportFileContent,
			validation.When(r.ImportFileSource == "BODY_BASE64", validation.Required.Error("must be set when ImportFileSource=='BODY_BASE64'")),
			validation.When(r.ImportFileSource == "URL", validation.Nil.Error("must not be set when ImportFileSource=='URL'"))),
		"ContractID":       validation.Validate(r.ContractID, validation.Required),
		"GroupID":          validation.Validate(r.GroupID, validation.Required),
		"ImportFileFormat": validation.Validate(r.ImportFileFormat, validation.Required),
		"ImportFileSource": validation.Validate(r.ImportFileSource, validation.Required),
		"ImportURL": validation.Validate(r.ImportURL,
			validation.When(r.ImportFileSource == "URL",
				validation.Required.Error("required field when ImportFileSource=='URL'"),
				is.URL.Error("must be a valid URL when ImportFileSource=='URL'"),
			).Else(
				validation.Nil.Error("should not be set when ImportFileSource=='BODY_BASE64'"),
			),
		),
		// it is not feasible to validate the root parameter from client side, because the file type is unknown until receiving the response
	})
}

// IsActive returns true if API is active on given network
func (n VersionState) IsActive() bool {
	return n.Status != nil && *n.Status == ActivationStatusActive
}

var (
	// ErrGetEndpoint is returned in case an error occurs on ShowEndpoint operation
	ErrGetEndpoint = errors.New("get endpoint")
	// ErrListEndpoints is returned in case an error occurs in ListEndpoints operation
	ErrListEndpoints = errors.New("list endpoints")
	// ErrRegisterEndpoint is returned in case an error occurs in RegisterEndpoint operation
	ErrRegisterEndpoint = errors.New("register endpoint")
	// ErrRegisterEndpointFromFile is returned in case an error occurs in RegisterEndpointFromFile operation
	ErrRegisterEndpointFromFile = errors.New("register endpoint from file")
	// ErrListUserEntitlements is returned in case an error occurs in ListUserEntitlements operation
	ErrListUserEntitlements = errors.New("list user entitlements")
	// ErrShowEndpoint is returned in case an error occurs on ShowEndpoint operation
	ErrShowEndpoint = errors.New("show endpoint")
	// ErrHideEndpoint is returned in case an error occurs on HideEndpoint operation
	ErrHideEndpoint = errors.New("hide endpoint")
	// ErrDeleteEndpoint is returned in case an error occurs on DeleteEndpoint operation
	ErrDeleteEndpoint = errors.New("delete endpoint")
)

func (a *apidefinitions) RegisterEndpoint(ctx context.Context, params RegisterEndpointRequest) (*RegisterEndpointResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("RegisterEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrRegisterEndpoint, ErrStructValidation, err)
	}

	uri := "/api-definitions/v2/endpoints"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRegisterEndpoint, err)
	}

	var result RegisterEndpointResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRegisterEndpoint, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %s", ErrRegisterEndpoint, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) RegisterEndpointFromFile(ctx context.Context, params RegisterEndpointFromFileRequest) (*RegisterEndpointFromFileResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("RegisterEndpointFromFile")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrRegisterEndpointFromFile, ErrStructValidation, err)
	}

	uri := "/api-definitions/v2/endpoints/files"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRegisterEndpointFromFile, err)
	}

	var result RegisterEndpointFromFileResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRegisterEndpointFromFile, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrRegisterEndpointFromFile, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) ShowEndpoint(ctx context.Context, params ShowEndpointRequest) (*ShowEndpointResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ShowEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrShowEndpoint, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/show", params.APIEndpointID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrShowEndpoint, err)
	}

	var result ShowEndpointResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrShowEndpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrShowEndpoint, a.Error(resp))
	}
	return &result, nil
}

func (a *apidefinitions) HideEndpoint(ctx context.Context, params HideEndpointRequest) (*HideEndpointResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("HideEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrHideEndpoint, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/hide", params.APIEndpointID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrHideEndpoint, err)
	}

	var result HideEndpointResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrHideEndpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrHideEndpoint, a.Error(resp))
	}
	return &result, nil
}

func (a *apidefinitions) DeleteEndpoint(ctx context.Context, params DeleteEndpointRequest) error {
	logger := a.Log(ctx)
	logger.Debug("DeleteEndpoint")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrDeleteEndpoint, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d", params.APIEndpointID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteEndpoint, err)
	}

	resp, err := a.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteEndpoint, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteEndpoint, a.Error(resp))
	}

	return nil
}

func (a *apidefinitions) ListEndpoints(ctx context.Context, params ListEndpointsRequest) (*ListEndpointsResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ListEndpoints")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListEndpoints, ErrStructValidation, err)
	}

	uri, err := url.Parse("/api-definitions/v2/endpoints")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListEndpoints, err)
	}

	q := uri.Query()

	if params.Page <= 0 {
		params.Page = 1
	}
	q.Add("page", strconv.FormatInt(params.Page, 10))

	if params.PageSize <= 0 {
		params.PageSize = 25
	}
	q.Add("pageSize", strconv.FormatInt(params.PageSize, 10))

	if params.Category != "" {
		q.Add("category", params.Category)
	}

	if params.Contains != "" {
		q.Add("contains", params.Contains)
	}

	if params.SortBy != "" {
		q.Add("sortBy", string(params.SortBy))
	}

	if params.SortOrder != "" {
		q.Add("sortOrder", string(params.SortOrder))
	}

	if params.VersionPreference != "" {
		q.Add("versionPreference", string(params.VersionPreference))
	}

	if params.Show != "" {
		q.Add("show", string(params.Show))
	}

	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}

	if params.GroupID != 0 {
		q.Add("groupId", strconv.FormatInt(params.GroupID, 10))
	}

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEndpoints, err)
	}

	var result ListEndpointsResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEndpoints, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEndpoints, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) GetEndpoint(ctx context.Context, params GetEndpointRequest) (*GetEndpointResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetEndpoint")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEndpoint, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d", params.APIEndpointID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEndpoint, err)
	}

	var result GetEndpointResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEndpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEndpoint, a.Error(resp))
	}
	return &result, nil
}

func (a *apidefinitions) ListUserEntitlements(ctx context.Context) (ListUserEntitlementsResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ListUserEntitlements")

	uri := "/api-definitions/v2/endpoints/user-entitlements"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListUserEntitlements, err)
	}

	var result ListUserEntitlementsResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListUserEntitlements, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListUserEntitlements, a.Error(resp))
	}

	return result, nil
}

func (r *restrictionsBool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	switch asString {
	case "0":
		*r = false
	case "1":
		*r = true
	default:
		return fmt.Errorf("boolean unmarshal error: invalid input %s", asString)
	}
	return nil
}

func (r restrictionsBool) MarshalJSON() ([]byte, error) {
	if r {
		return json.Marshal(1)
	}
	return json.Marshal(0)
}
