// Package v0 provides access to the Akamai APIDefinitions V0 API
package v0

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type (
	// APIDefinitions is the api definitions api interface
	APIDefinitions interface {
		// RegisterAPI creates the first version of an API endpoint configuration
		RegisterAPI(context.Context, RegisterAPIRequest) (*RegisterAPIResponse, error)
		// GetAPIVersion returns an API version. Use this operation's response object when modifying an endpoint version through Edit a version.
		GetAPIVersion(context.Context, GetAPIVersionRequest) (*GetAPIVersionResponse, error)
		// UpdateAPIVersion updates details about an API version that has never been activated on the staging or production network.
		UpdateAPIVersion(context.Context, UpdateAPIVersionRequest) (*UpdateAPIVersionResponse, error)
		// FromOpenAPIFile map OpenAPI file to API
		FromOpenAPIFile(context.Context, FromOpenAPIFileRequest) (*FromOpenAPIFileResponse, error)
		// ToOpenAPIFile map API to OpenAPI
		ToOpenAPIFile(context.Context, ToOpenAPIFileRequest) (*ToOpenAPIFileResponse, error)
		// GetResourceOperation reads resource operations for a particular endpoint
		GetResourceOperation(context.Context, GetResourceOperationRequest) (*GetResourceOperationResponse, error)
		// UpdateResourceOperation updates resource operations for a particular endpoint
		UpdateResourceOperation(context.Context, UpdateResourceOperationRequest) (*UpdateResourceOperationResponse, error)
		// DeleteResourceOperation deletes resource operations for a particular endpoint
		DeleteResourceOperation(context.Context, DeleteResourceOperationRequest) (*DeleteResourceOperationResponse, error)
	}
	// RegisterAPIRequest contains body for RegisterAPI operation
	RegisterAPIRequest struct {
		APIAttributes
		ContractID string `json:"contractId"`
		GroupID    int64  `json:"groupId"`
	}

	// APIAttributes hold information about API (without contract and group)
	APIAttributes struct {
		Name                      string                                   `json:"name"`
		Hostnames                 []string                                 `json:"hostnames"`
		BasePath                  *string                                  `json:"basePath,omitempty"`
		Tags                      []string                                 `json:"tags,omitempty"`
		Description               *string                                  `json:"description,omitempty"`
		MatchPathSegmentParameter bool                                     `json:"matchPathSegmentParameter,omitempty"`
		MatchCaseSensitive        bool                                     `json:"matchCaseSensitive,omitempty"`
		EnableAPIGateway          bool                                     `json:"enableApiGateway,omitempty"`
		GraphQL                   bool                                     `json:"graphQl,omitempty"`
		SecuritySchemes           *SecuritySchemes                         `json:"securitySchemes,omitempty"`
		Constraints               *Constraints                             `json:"constraints,omitempty"`
		Versioning                *Versioning                              `json:"versioning,omitempty"`
		Resources                 *orderedmap.OrderedMap[string, Resource] `json:"resources,omitempty"`
	}

	// RegisterAPIResponse holds the response from RegisterAPI operation
	RegisterAPIResponse API

	// GetAPIVersionRequest contains parameters for GetAPIVersion method
	GetAPIVersionRequest struct {
		Version int64
		ID      int64
	}

	// GetAPIVersionResponse holds the response from GetAPIVersion operation
	GetAPIVersionResponse API

	// UpdateAPIVersionRequest contains parameters for UpdateEndpointVersion method
	UpdateAPIVersionRequest struct {
		ID      int64
		Version int64
		Body    UpdateAPIVersionRequestBody
	}

	// UpdateAPIVersionRequestBody contains body for UpdateAPIVersion operation
	UpdateAPIVersionRequestBody API

	// UpdateAPIVersionResponse holds the response from UpdateAPIVersion operation
	UpdateAPIVersionResponse API

	// FromOpenAPIFileRequest contains body for FromOpenAPIFile operation
	FromOpenAPIFileRequest struct {
		Content  []byte
		RootFile *string
	}

	// ToOpenAPIFileRequest contains parameters for ToOpenAPIFile method
	ToOpenAPIFileRequest GetAPIVersionRequest

	// ToOpenAPIFileResponse holds the response from ToOpenAPIFile operation
	ToOpenAPIFileResponse string

	apidefinitions struct {
		session.Session
	}

	// Option defines a api definition option
	Option func(*apidefinitions)

	// ClientFunc is an apidefinitions client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) APIDefinitions

	// API holds configuration for an API
	API struct {
		RegisterAPIRequest
		ID            *int64 `json:"id,omitempty"`
		RecordVersion *int64 `json:"recordVersion,omitempty"`
	}

	// SecuritySchemes holds configuration for Security Schemes
	SecuritySchemes struct {
		APIKey *SecurityScheme `json:"apiKey,omitempty"`
	}

	// SecuritySchemeLocation holds location of the SecurityScheme
	SecuritySchemeLocation string

	// SecurityScheme holds configuration for Security Scheme
	SecurityScheme struct {
		In   *SecuritySchemeLocation `json:"in,omitempty"`
		Name *string                 `json:"name,omitempty"`
	}

	// VersioningLocation holds location of the version
	VersioningLocation string

	// Versioning holds configuration for Versioning configuration
	Versioning struct {
		In    *VersioningLocation `json:"in,omitempty"`
		Name  *string             `json:"name,omitempty"`
		Value *string             `json:"value,omitempty"`
	}

	// Constraints holds configuration for Constraints
	Constraints struct {
		EnforceOn   *EnforceOn              `json:"enforceOn,omitempty"`
		RequestBody *ConstraintsRequestBody `json:"requestBody,omitempty"`
	}

	// EnforceOn holds configuration for Constraints enforcement
	EnforceOn struct {
		Request             *bool                `json:"request,omitempty"`
		Response            *bool                `json:"response,omitempty"`
		UndefinedMethods    *UndefinedMethods    `json:"undefinedMethods,omitempty"`
		UndefinedParameters *UndefinedParameters `json:"undefinedParameters,omitempty"`
	}

	// UndefinedMethods hold configuration for undefined method Constraints enforcement
	UndefinedMethods struct {
		Get     bool `json:"get,omitempty"`
		Post    bool `json:"post,omitempty"`
		Put     bool `json:"put,omitempty"`
		Head    bool `json:"head,omitempty"`
		Options bool `json:"options,omitempty"`
		Delete  bool `json:"delete,omitempty"`
		Patch   bool `json:"patch,omitempty"`
	}

	// UndefinedParameters hold configuration for undefined parameters Constraints enforcement
	UndefinedParameters struct {
		RequestCookie  bool `json:"requestCookie,omitempty"`
		RequestHeader  bool `json:"requestHeader,omitempty"`
		RequestQuery   bool `json:"requestQuery,omitempty"`
		RequestBody    bool `json:"requestBody,omitempty"`
		ResponseHeader bool `json:"responseHeader,omitempty"`
		ResponseBody   bool `json:"responseBody,omitempty"`
	}

	// ConsumeType content type the endpoint exchanges
	ConsumeType string

	// ConstraintsRequestBody holds configuration for Constraints
	ConstraintsRequestBody struct {
		ConsumeType     []ConsumeType                     `json:"consumeType,omitempty"`
		MaxBodySize     *int64                            `json:"maxContentLength,omitempty"`
		MaxNestingDepth *int64                            `json:"maxNestingDepth,omitempty"`
		Properties      *ConstraintsRequestBodyProperties `json:"properties,omitempty"`
	}

	// ConstraintsRequestBodyProperties holds configuration for Constraints
	ConstraintsRequestBodyProperties struct {
		MaxStringLength *int64 `json:"maxStringLength,omitempty"`
		MaxIntegerValue *int64 `json:"maxIntegerValue,omitempty"`
		MaxCount        *int64 `json:"maxCount,omitempty"`
		MaxNameLength   *int64 `json:"maxNameLength,omitempty"`
	}

	// Resource holds configuration for an API Resource
	Resource struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Get         *Method `json:"get,omitempty"`
		Post        *Method `json:"post,omitempty"`
		Put         *Method `json:"put,omitempty"`
		Delete      *Method `json:"delete,omitempty"`
		Options     *Method `json:"options,omitempty"`
		Head        *Method `json:"head,omitempty"`
		Patch       *Method `json:"patch,omitempty"`
	}

	// Method holds configuration for an API Method
	Method struct {
		Parameters  []Parameter                              `json:"parameters,omitempty"`
		RequestBody *orderedmap.OrderedMap[string, Property] `json:"requestBody,omitempty"`
		Responses   *Responses                               `json:"responses,omitempty"`
		Constraints *MethodConstraints                       `json:"constraints,omitempty"`
	}

	// MethodConstraints holds configuration for Method Constraints
	MethodConstraints struct {
		EnforceOn *MethodEnforceOn `json:"enforceOn,omitempty"`
	}

	// MethodEnforceOn holds configuration for Method Constraints
	MethodEnforceOn struct {
		UndefinedParameters *UndefinedParameters `json:"undefinedParameters,omitempty"`
	}

	// Responses holds configuration for an API Responses
	Responses struct {
		Headers  []Parameter       `json:"headers,omitempty"`
		Contents []ResponseContent `json:"contents,omitempty"`
	}

	// ResponseContent holds configuration for an API Response
	ResponseContent struct {
		StatusCodes []int64   `json:"statusCodes,omitempty"`
		JSON        *Property `json:"json,omitempty"`
		GraphQL     *Property `json:"graphql,omitempty"`
		XML         *Property `json:"xml,omitempty"`
		URLEncoded  *Property `json:"urlencoded,omitempty"`
		JSONXML     *Property `json:"json/xml,omitempty"`
		Any         *Property `json:"any,omitempty"`
		None        *Property `json:"none,omitempty"`
	}

	// ParameterType type of the parameter
	ParameterType string

	// ParameterLocation location of the parameter
	ParameterLocation string

	// Parameter holds configuration for an API Parameter
	Parameter struct {
		Name        string            `json:"name"`
		Type        ParameterType     `json:"type"`
		In          ParameterLocation `json:"in,omitempty"`
		Required    bool              `json:"required,omitempty"`
		Description *string           `json:"description,omitempty"`
		Minimum     *float32          `json:"minimum,omitempty"`
		Maximum     *float32          `json:"maximum,omitempty"`
		MinLength   *int64            `json:"minLength,omitempty"`
		MaxLength   *int64            `json:"maxLength,omitempty"`
	}

	// PropertyType type of the property
	PropertyType ParameterType

	// MaxBodySize represents MaxBodySize value
	MaxBodySize string

	// XML holds configuration about an XML representation of property
	XML struct {
		Name      *string `json:"name,omitempty"`
		Namespace *string `json:"namespace,omitempty"`
		Prefix    *string `json:"prefix,omitempty"`
		Attribute *bool   `json:"attribute,omitempty"`
		Wrapped   *bool   `json:"wrapped,omitempty"`
	}

	// Property holds configuration for an API Property
	Property struct {
		Name        string       `json:"name"`
		Type        PropertyType `json:"type"`
		Required    bool         `json:"required,omitempty"`
		Description *string      `json:"description,omitempty"`
		Minimum     *float32     `json:"minimum,omitempty"`
		Maximum     *float32     `json:"maximum,omitempty"`
		MinLength   *int64       `json:"minLength,omitempty"`
		MaxLength   *int64       `json:"maxLength,omitempty"`
		MinItems    *int64       `json:"minItems,omitempty"`
		MaxItems    *int64       `json:"maxItems,omitempty"`
		MaxBodySize *MaxBodySize `json:"maxBodySize,omitempty"`
		Properties  []Property   `json:"properties,omitempty"`
		Items       *Items       `json:"items,omitempty"`
		XML         *XML         `json:"xml,omitempty"`
	}

	// Items holds configuration for an Array Property Items
	Items struct {
		Type       PropertyType `json:"type"`
		Required   bool         `json:"required,omitempty"`
		Minimum    *float32     `json:"minimum,omitempty"`
		Maximum    *float32     `json:"maximum,omitempty"`
		MinLength  *int64       `json:"minLength,omitempty"`
		MaxLength  *int64       `json:"maxLength,omitempty"`
		MinItems   *int64       `json:"minItems,omitempty"`
		MaxItems   *int64       `json:"maxItems,omitempty"`
		Properties []Property   `json:"properties,omitempty"`
	}

	// FromOpenAPIFileResponse holds the response for FromOpenAPIFile operation
	FromOpenAPIFileResponse struct {
		Problems []Error       `json:"problems"`
		API      APIAttributes `json:"api"`
	}
)

var (
	// ErrRegisterAPI is returned when GetAPIVersion fails
	ErrRegisterAPI = errors.New("register API")
	// ErrGetAPIVersion is returned when GetEndpointVersion fails
	ErrGetAPIVersion = errors.New("get API version")
	// ErrUpdateAPIVersion is returned when UpdateEndpointVersion fails
	ErrUpdateAPIVersion = errors.New("update API version")
	// ErrFromOpenAPIFile is returned when FromOpenAPIFile fails
	ErrFromOpenAPIFile = errors.New("mapping openapi file")
	// ErrToOpenAPIFile is returned when ErrToOpenAPIFile fails
	ErrToOpenAPIFile = errors.New("to openapi file")
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

const (
	// SecuritySchemeLocationHeader holds value for http 'header' security scheme location
	SecuritySchemeLocationHeader SecuritySchemeLocation = "header"
	// SecuritySchemeLocationQuery holds value for http 'query' security scheme location
	SecuritySchemeLocationQuery SecuritySchemeLocation = "query"
	// SecuritySchemeLocationCookie holds value for http 'cookie' security scheme location
	SecuritySchemeLocationCookie SecuritySchemeLocation = "cookie"

	// ConsumeTypeJSON holds value for consume type json
	ConsumeTypeJSON ConsumeType = "json"
	// ConsumeTypeXML holds value for consume type xml
	ConsumeTypeXML ConsumeType = "xml"
	// ConsumeTypeUrlencoded holds value for consume type urlencoded
	ConsumeTypeUrlencoded ConsumeType = "urlencoded"
	// ConsumeTypeAny holds value for consume type any
	ConsumeTypeAny ConsumeType = "any"

	// ParameterTypeString holds value for string ParameterType
	ParameterTypeString ParameterType = "string"
	// ParameterTypeNumber holds value for number ParameterType
	ParameterTypeNumber ParameterType = "number"
	// ParameterTypeInteger holds value for integer ParameterType
	ParameterTypeInteger ParameterType = "integer"
	// ParameterTypeBoolean holds value for boolean ParameterType
	ParameterTypeBoolean ParameterType = "boolean"

	// ParameterLocationQuery holds value for 'query' parameter location
	ParameterLocationQuery ParameterLocation = "query"
	// ParameterLocationHeader holds value for 'header' parameter location
	ParameterLocationHeader ParameterLocation = "header"
	// ParameterLocationCookie holds value for 'cookie' parameter location
	ParameterLocationCookie ParameterLocation = "cookie"
	// ParameterLocationPath holds value for 'path' parameter location
	ParameterLocationPath ParameterLocation = "path"

	// PropertyTypeObject holds value for boolean PropertyType
	PropertyTypeObject PropertyType = "object"
	// PropertyTypeArray holds value for boolean PropertyType
	PropertyTypeArray PropertyType = "array"

	// MaxBodySizeSize6KB represents MaxBodySize of value "6KB"
	MaxBodySizeSize6KB MaxBodySize = "6KB"
	// MaxBodySizeSize8KB represents MaxBodySize of value "8KB"
	MaxBodySizeSize8KB MaxBodySize = "8KB"
	// MaxBodySizeSize12KB represents MaxBodySize of value "12KB"
	MaxBodySizeSize12KB MaxBodySize = "12KB"
	// MaxBodySizeSize16kB represents MaxBodySize of value "16KB"
	MaxBodySizeSize16kB MaxBodySize = "16KB"

	// VersioningLocationHeader holds value for versioning location 'header'
	VersioningLocationHeader VersioningLocation = "header"
	// VersioningLocationPath holds value for versioning location 'path'
	VersioningLocationPath VersioningLocation = "path"
	// VersioningLocationQuery holds value for versioning location 'query'
	VersioningLocationQuery VersioningLocation = "query"
)

// Client returns a new apidefinitions Client instance with the specified controller
func Client(sess session.Session, opts ...Option) APIDefinitions {
	a := &apidefinitions{
		Session: sess,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Validate validates API
func (r API) Validate() error {
	return r.RegisterAPIRequest.Validate()
}

// Validate validates RegisterAPIRequest
func (r RegisterAPIRequest) Validate() error {

	validationErrors := validation.Errors{
		"Name":            validation.Validate(r.Name, validation.Required),
		"Hostnames":       validation.Validate(r.Hostnames, validation.Required),
		"ContractID":      validation.Validate(r.ContractID, validation.Required),
		"GroupID":         validation.Validate(r.GroupID, validation.Required),
		"Constraints":     validation.Validate(r.Constraints),
		"Resources":       validation.Validate(r.Resources),
		"Versioning":      validation.Validate(r.Versioning),
		"SecuritySchemes": validation.Validate(r.SecuritySchemes),
	}

	if r.Resources != nil {
		for pair := r.Resources.Oldest(); pair != nil; pair = pair.Next() {
			if err := pair.Value.Validate(); err != nil {
				validationErrors[pair.Key] = err
			}
		}
	}
	return edgegriderr.ParseValidationErrors(validationErrors)
}

// Validate validates SecuritySchemes
func (v SecuritySchemes) Validate() error {
	return validation.Errors{
		"APIKey": validation.Validate(v.APIKey),
	}.Filter()
}

// Validate validates SecurityScheme
func (v SecurityScheme) Validate() error {
	return validation.Errors{
		"In": validation.Validate(v.In),
	}.Filter()
}

// Validate validates SecuritySchemeLocation
func (s SecuritySchemeLocation) Validate() error {
	return validation.In(SecuritySchemeLocationCookie, SecuritySchemeLocationHeader, SecuritySchemeLocationQuery).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s'", s, SecuritySchemeLocationCookie, SecuritySchemeLocationHeader, SecuritySchemeLocationQuery)).
		Validate(s)
}

// Validate validates Versioning
func (v Versioning) Validate() error {
	return validation.Errors{
		"In": validation.Validate(v.In),
	}.Filter()
}

// Validate validates ConsumeType
func (s VersioningLocation) Validate() error {
	return validation.In(VersioningLocationHeader, VersioningLocationPath, VersioningLocationQuery).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s'", s, VersioningLocationHeader, VersioningLocationPath, VersioningLocationQuery)).
		Validate(s)
}

// Validate validates Constraints
func (c Constraints) Validate() error {
	return validation.Errors{
		"RequestBody": validation.Validate(c.RequestBody),
		"EnforceOn":   validation.Validate(c.EnforceOn),
	}.Filter()
}

// Validate validates EnforceOn
func (b EnforceOn) Validate() error {
	return validation.Errors{
		"UndefinedMethods":    validation.Validate(b.UndefinedMethods),
		"UndefinedParameters": validation.Validate(b.UndefinedParameters),
	}.Filter()
}

// Validate validates ConstraintsRequestBody
func (c ConstraintsRequestBody) Validate() error {
	return validation.Errors{
		"ConsumeType": validation.Validate(c.ConsumeType),
	}.Filter()
}

// Validate validates ConsumeType
func (s ConsumeType) Validate() error {
	return validation.In(ConsumeTypeJSON, ConsumeTypeXML, ConsumeTypeUrlencoded, ConsumeTypeAny).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s', '%s'", s, ConsumeTypeJSON, ConsumeTypeXML, ConsumeTypeUrlencoded, ConsumeTypeAny)).
		Validate(s)
}

// Validate validates Resource
func (r Resource) Validate() error {
	return validation.Errors{
		"Get":     validation.Validate(r.Get),
		"Post":    validation.Validate(r.Post),
		"Put":     validation.Validate(r.Put),
		"Patch":   validation.Validate(r.Patch),
		"Delete":  validation.Validate(r.Delete),
		"Head":    validation.Validate(r.Head),
		"Options": validation.Validate(r.Options),
	}.Filter()
}

// Validate validates Method
func (m Method) Validate() error {
	validationErrors := validation.Errors{
		"Parameters":  validation.Validate(m.Parameters),
		"RequestBody": validation.Validate(m.RequestBody),
		"Responses":   validation.Validate(m.Responses),
	}

	if m.RequestBody != nil {
		for pair := m.RequestBody.Oldest(); pair != nil; pair = pair.Next() {
			if err := pair.Value.Validate(); err != nil {
				validationErrors[pair.Key] = err
			}
		}
	}
	return validationErrors.Filter()
}

// Validate validates Responses
func (p Responses) Validate() error {
	return validation.Errors{
		"Headers":  validation.Validate(p.Headers),
		"Contents": validation.Validate(p.Contents),
	}.Filter()
}

// Validate validates ResponseContent
func (p ResponseContent) Validate() error {
	return validation.Errors{
		"JSON":       validation.Validate(p.JSON),
		"GraphQL":    validation.Validate(p.GraphQL),
		"XML":        validation.Validate(p.XML),
		"URLEncoded": validation.Validate(p.URLEncoded),
		"JSONXML":    validation.Validate(p.JSONXML),
		"Any":        validation.Validate(p.Any),
		"None":       validation.Validate(p.None),
	}.Filter()
}

// Validate validates Parameter
func (p Parameter) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(p.Name, validation.Required),
		"Type": validation.Validate(p.Type, validation.Required),
		"In":   validation.Validate(p.In),
	}.Filter()
}

// Validate validates ParameterType
func (t ParameterType) Validate() error {
	return validation.In(ParameterTypeNumber, ParameterTypeInteger, ParameterTypeString, ParameterTypeBoolean).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s', '%s'", t, ParameterTypeNumber, ParameterTypeInteger, ParameterTypeString, ParameterTypeBoolean)).
		Validate(t)
}

// Validate validates ParameterLocation
func (t ParameterLocation) Validate() error {
	return validation.In(ParameterLocationCookie, ParameterLocationQuery, ParameterLocationHeader, ParameterLocationPath).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s', '%s'", t, ParameterLocationCookie, ParameterLocationQuery, ParameterLocationHeader, ParameterLocationPath)).
		Validate(t)
}

// Validate validates PropertyType
func (t PropertyType) Validate() error {
	return validation.In(string(ParameterTypeNumber), string(ParameterTypeInteger), string(ParameterTypeString), string(ParameterTypeBoolean), string(PropertyTypeObject), string(PropertyTypeArray)).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s'", t, ParameterTypeNumber, ParameterTypeInteger, ParameterTypeString, ParameterTypeBoolean, PropertyTypeObject, PropertyTypeArray)).
		Validate(string(t))
}

// Validate validates MaxBodySize
func (b MaxBodySize) Validate() error {
	return validation.In(MaxBodySizeSize6KB, MaxBodySizeSize8KB, MaxBodySizeSize12KB, MaxBodySizeSize16kB).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s', '%s' ", b, MaxBodySizeSize6KB, MaxBodySizeSize8KB, MaxBodySizeSize12KB, MaxBodySizeSize16kB)).
		Validate(b)
}

// Validate validates Property
func (p Property) Validate() error {
	return validation.Errors{
		"Name":        validation.Validate(p.Name, validation.Required),
		"Type":        validation.Validate(p.Type, validation.Required),
		"MaxBodySize": validation.Validate(p.MaxBodySize),
		"Items":       validation.Validate(p.Items),
		"Properties":  validation.Validate(p.Properties),
	}.Filter()
}

// Validate validates Items
func (p Items) Validate() error {
	return validation.Errors{
		"Type":       validation.Validate(p.Type),
		"Properties": validation.Validate(p.Properties),
	}.Filter()
}

// Validate validates GetAPIVersionRequest
func (r GetAPIVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ID":      validation.Validate(r.ID, validation.Required),
		"Version": validation.Validate(r.Version, validation.Required),
	})
}

// Validate validates UpdateAPIVersionRequest
func (r UpdateAPIVersionRequest) Validate() interface{} {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ID":      validation.Validate(r.ID, validation.Required),
		"Version": validation.Validate(r.Version, validation.Required),
		"Body":    validation.Validate(r.Body, validation.Required),
	})
}

func (a *apidefinitions) RegisterAPI(ctx context.Context, params RegisterAPIRequest) (*RegisterAPIResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("RegisterAPI")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrRegisterAPI, ErrStructValidation, err)
	}

	uri := "/api-definitions/v0/endpoints"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRegisterAPI, err)
	}

	var result RegisterAPIResponse
	resp, err := a.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRegisterAPI, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrRegisterAPI, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) GetAPIVersion(ctx context.Context, params GetAPIVersionRequest) (*GetAPIVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetAPIVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAPIVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v0/endpoints/%d/versions/%d", params.ID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAPIVersion, err)
	}

	var result GetAPIVersionResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetAPIVersion, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAPIVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) UpdateAPIVersion(ctx context.Context, params UpdateAPIVersionRequest) (*UpdateAPIVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("UpdateAPIVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateAPIVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v0/endpoints/%d/versions/%d", params.ID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateAPIVersion, err)
	}

	var result UpdateAPIVersionResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateAPIVersion, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateAPIVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) FromOpenAPIFile(ctx context.Context, body FromOpenAPIFileRequest) (*FromOpenAPIFileResponse, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("importFile", "api.json")
	if err != nil {
		return nil, fmt.Errorf("%w: io failed: %s", ErrFromOpenAPIFile, err)
	}

	_, err = io.Copy(part, bytes.NewBuffer(body.Content))
	if err != nil {
		return nil, fmt.Errorf("%w: io failed: %s", ErrFromOpenAPIFile, err)
	}

	if body.RootFile != nil {
		part, err = writer.CreateFormField("root")
		if err != nil {
			return nil, fmt.Errorf("%w: io failed: %s", ErrFromOpenAPIFile, err)
		}

		_, err := part.Write([]byte(*body.RootFile))
		if err != nil {
			return nil, fmt.Errorf("%w: io failed: %s", ErrFromOpenAPIFile, err)
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: io failed: %s", ErrFromOpenAPIFile, err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api-definitions/v0/endpoints/openapi", requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrFromOpenAPIFile, err)
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	var result = FromOpenAPIFileResponse{}
	resp, err := a.Exec(r, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrFromOpenAPIFile, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrFromOpenAPIFile, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) ToOpenAPIFile(ctx context.Context, params ToOpenAPIFileRequest) (*ToOpenAPIFileResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ToOpenAPIFile")

	if err := GetAPIVersionRequest(params).Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrToOpenAPIFile, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v0/endpoints/%d/versions/%d/openapi", params.ID, params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrToOpenAPIFile, err)
	}

	resp, err := a.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrToOpenAPIFile, err)
	}
	defer session.CloseResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrToOpenAPIFile, a.Error(resp))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response body: %s", ErrToOpenAPIFile, err)
	}

	return ptr.To(ToOpenAPIFileResponse(data)), nil
}
