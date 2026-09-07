package apidefinitions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListEndpointVersionsRequest contains parameters for ListEndpointVersions
	ListEndpointVersionsRequest struct {
		APIEndpointID int64
		Page          int64
		PageSize      int64
		SortBy        ListEndpointVersionSortType
		SortOrder     SortOrderType
		Show          Visibility
	}

	// ListEndpointVersionSortType represents the type of the sorting is based on
	ListEndpointVersionSortType string

	// SortOrderType represents the type of the sort order, either ascending or descending
	SortOrderType string

	// Visibility represents the visibility of the endpoint, either 'ALL', 'ONLY_HIDDEN' or 'ONLY_VISIBLE'
	Visibility string

	// EndpointVersionRequest contains parameters for DeleteEndpointVersion, CloneEndpointVersion and GetEndpointVersion methods
	EndpointVersionRequest struct {
		VersionNumber int64
		APIEndpointID int64
	}

	// DeleteEndpointVersionRequest contains parameters for DeleteEndpointVersion method
	DeleteEndpointVersionRequest EndpointVersionRequest

	// CloneEndpointVersionRequest contains parameters for CloneEndpointVersion method
	CloneEndpointVersionRequest EndpointVersionRequest

	// GetEndpointVersionRequest contains parameters for GetEndpointVersion method
	GetEndpointVersionRequest EndpointVersionRequest

	// UpdateEndpointVersionRequest contains parameters for UpdateEndpointVersion method
	UpdateEndpointVersionRequest struct {
		VersionNumber int64
		APIEndpointID int64
		Body          UpdateEndpointVersionRequestBody
	}

	// ListEndpointVersionsResponse holds response data for ListEndpointVersions
	ListEndpointVersionsResponse struct {
		TotalSize       int64        `json:"totalSize"`
		Page            int64        `json:"page"`
		PageSize        int64        `json:"pageSize"`
		APIEndpointID   int64        `json:"apiEndPointId"`
		APIEndpointName string       `json:"apiEndPointName"`
		APIVersions     []APIVersion `json:"apiVersions"`
	}

	// GetEndpointVersionResponse holds response parameter for GetEndpointVersion method
	GetEndpointVersionResponse EndpointVersionResponse

	// CloneEndpointVersionResponse holds response parameter for CloneEndpointVersion method
	CloneEndpointVersionResponse EndpointVersionResponse

	// UpdateEndpointVersionResponse holds response parameter for UpdateEndpointVersion method
	UpdateEndpointVersionResponse EndpointVersionResponse

	// APIVersion represents version object returned by ListEndpointVersions
	APIVersion struct {
		CreateDate           string            `json:"createDate"`
		CreatedBy            string            `json:"createdBy"`
		UpdateDate           string            `json:"updateDate"`
		UpdatedBy            string            `json:"updatedBy"`
		APIEndpointVersionID int64             `json:"apiEndPointVersionId"`
		BasePath             string            `json:"basePath"`
		VersionNumber        int64             `json:"versionNumber"`
		Description          *string           `json:"description"`
		BasedOn              *int64            `json:"basedOn"`
		StagingStatus        *ActivationStatus `json:"stagingStatus"`
		ProductionStatus     *ActivationStatus `json:"productionStatus"`
		StagingDate          *string           `json:"stagingDate"`
		ProductionDate       *string           `json:"productionDate"`
		IsVersionLocked      bool              `json:"isVersionLocked"`
		Hidden               bool              `json:"hidden"`
		AvailableActions     []string          `json:"availableActions"`
		CloningStatus        *string           `json:"cloningStatus"`
		LockVersion          int64             `json:"lockVersion"`
	}

	// EndpointVersionResponse holds response body parameters used in GetEndpointVersion, UpdateEndpointVersion, CloneEndpointVersion
	EndpointVersionResponse struct {
		SecurityScheme             *SecurityScheme             `json:"securityScheme"`
		AkamaiSecurityRestrictions *AkamaiSecurityRestrictions `json:"akamaiSecurityRestrictions"`
		ContractID                 string                      `json:"contractId"`
		GroupID                    int64                       `json:"groupId"`
		APIEndpointID              int64                       `json:"apiEndPointId"`
		APIEndpointVersion         *int64                      `json:"apiEndPointVersion"`
		VersionNumber              int64                       `json:"versionNumber"`
		APIEndpointName            string                      `json:"apiEndPointName"`
		Description                *string                     `json:"description"`
		BasePath                   string                      `json:"basePath"`
		ClonedFromVersion          *int64                      `json:"clonedFromVersion"`
		APIEndpointLocked          bool                        `json:"apiEndPointLocked"`
		APIEndpointScheme          *string                     `json:"apiEndPointScheme"`
		ConsumeType                *string                     `json:"consumeType"`
		APIEndpointHosts           []string                    `json:"apiEndPointHosts"`
		APICategoryIDs             []int64                     `json:"apiCategoryIds"`
		LockVersion                int64                       `json:"lockVersion"`
		UpdatedBy                  string                      `json:"updatedBy"`
		CreatedBy                  string                      `json:"createdBy"`
		CreateDate                 string                      `json:"createDate"`
		UpdateDate                 string                      `json:"updateDate"`
		PositiveConstrainsEnabled  bool                        `json:"positiveConstrainsEnabled"`
		CaseSensitive              *bool                       `json:"caseSensitive"`
		MatchPathSegmentParam      bool                        `json:"matchPathSegmentParam"`
		Source                     *Source                     `json:"source"`
		StagingVersion             *VersionState               `json:"stagingVersion"`
		ProductionVersion          *VersionState               `json:"productionVersion"`
		ProductionStatus           *string                     `json:"productionStatus"`
		StagingStatus              *string                     `json:"stagingStatus"`
		ProtectedByAPIKey          bool                        `json:"protectedByApiKey"`
		IsGraphQL                  bool                        `json:"isGraphQL"`
		AvailableActions           []string                    `json:"availableActions"`
		VersionHidden              bool                        `json:"versionHidden"`
		EndpointHidden             bool                        `json:"endpointHidden"`
		APISource                  *string                     `json:"apiSource"`
		APISourceDetails           []APISourceDiff             `json:"apiSourceDetails"`
		CloningStatus              *string                     `json:"cloningStatus"`
		APIGatewayEnabled          *bool                       `json:"apiGatewayEnabled"`
		GraphQL                    bool                        `json:"graphQL"`
		DiscoveredPIIIDs           []int64                     `json:"discoveredPiiIds"`
		APIVersionInfo             *APIVersionInfo             `json:"apiVersionInfo"`
		APIResources               []APIResource               `json:"apiResources"`
		Locked                     bool                        `json:"locked"`
	}

	// UpdateEndpointVersionRequestBody holds request body parameters of UpdateEndpointVersionRequest
	UpdateEndpointVersionRequestBody struct {
		SecurityScheme             *SecurityScheme             `json:"securityScheme"`
		AkamaiSecurityRestrictions *AkamaiSecurityRestrictions `json:"akamaiSecurityRestrictions,omitempty"`
		ContractID                 string                      `json:"contractId"`
		GroupID                    int64                       `json:"groupId"`
		APIEndpointID              int64                       `json:"apiEndPointId"`
		APIEndpointVersion         *int64                      `json:"apiEndPointVersion"`
		VersionNumber              int64                       `json:"versionNumber"`
		APIEndpointName            string                      `json:"apiEndPointName"`
		Description                *string                     `json:"description"`
		BasePath                   string                      `json:"basePath,omitempty"`
		APIEndpointScheme          *APIEndpointScheme          `json:"apiEndPointScheme"`
		ConsumeType                *ConsumeType                `json:"consumeType,omitempty"`
		APIEndpointHosts           []string                    `json:"apiEndPointHosts"`
		APICategoryIDs             []int64                     `json:"apiCategoryIds"`
		LockVersion                int64                       `json:"lockVersion"`
		CaseSensitive              *bool                       `json:"caseSensitive,omitempty"`
		MatchPathSegmentParam      bool                        `json:"matchPathSegmentParam"`
		IsGraphQL                  bool                        `json:"isGraphQL"`
		APIGatewayEnabled          *bool                       `json:"apiGatewayEnabled,omitempty"`
		GraphQL                    bool                        `json:"graphQL"`
		APIVersionInfo             *APIVersionInfo             `json:"apiVersionInfo,omitempty"`
		APIResources               []APIResource               `json:"apiResources,omitempty"`
	}
)

const (
	// DescriptionSort sorts by 'Description'
	DescriptionSort ListEndpointVersionSortType = "description"
	// VersionNumberSort sorts by 'VersionNumber'
	VersionNumberSort ListEndpointVersionSortType = "versionNumber"
	// UpdateDateSort sorts by 'UpdateDate'
	UpdateDateSort ListEndpointVersionSortType = "updateDate"
	// UpdatedBySort sorts by 'UpdatedBy'
	UpdatedBySort ListEndpointVersionSortType = "updatedBy"
	// BasedOnSort sorts by 'BasedOn'
	BasedOnSort ListEndpointVersionSortType = "basedOn"
	// StagingStatusSort sorts by 'StagingStatus'
	StagingStatusSort ListEndpointVersionSortType = "stagingStatus"
	// ProductionStatusSort sorts by 'ProductionStatus'
	ProductionStatusSort ListEndpointVersionSortType = "productionStatus"

	// AscSortOrder sets ordering to ascending
	AscSortOrder SortOrderType = "asc"
	// DescSortOrder sets ordering to descending
	DescSortOrder SortOrderType = "desc"

	// AllVisibility lists all endpoint versions
	AllVisibility Visibility = "ALL"
	// OnlyHiddenVisibility lists all hidden endpoint versions
	OnlyHiddenVisibility Visibility = "ONLY_HIDDEN"
	// OnlyVisibleVisibility lists all visible endpoint versions
	OnlyVisibleVisibility Visibility = "ONLY_VISIBLE"
)

var (
	// ErrGetEndpointVersion is returned when GetEndpointVersion fails
	ErrGetEndpointVersion = errors.New("get endpoint version")
	// ErrDeleteEndpointVersion is returned when DeleteEndpointVersion fails
	ErrDeleteEndpointVersion = errors.New("delete endpoint version")
	// ErrListEndpointVersions is return when ListEndpointVersions fails
	ErrListEndpointVersions = errors.New("list endpoint versions")
	// ErrCloneEndpointVersion is returned when CloneEndpointVersion fails
	ErrCloneEndpointVersion = errors.New("clone endpoint version")
	// ErrUpdateEndpointVersion is returned when UpdateEndpointVersion fails
	ErrUpdateEndpointVersion = errors.New("update endpoint version")
)

// Validate validates ListEndpointVersionSortType
func (s ListEndpointVersionSortType) Validate() error {
	return validation.In(DescriptionSort, VersionNumberSort, UpdateDateSort, UpdatedBySort, BasedOnSort, StagingStatusSort, ProductionStatusSort).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s', '%s'.", s, DescriptionSort, VersionNumberSort, UpdateDateSort, UpdatedBySort, BasedOnSort, StagingStatusSort, ProductionStatusSort)).
		Validate(s)
}

// Validate validates SortOrderType
func (s SortOrderType) Validate() error {
	return validation.In(AscSortOrder, DescSortOrder).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'.", s, AscSortOrder, DescSortOrder)).
		Validate(s)
}

// Validate validates Visibility
func (v Visibility) Validate() error {
	return validation.In(AllVisibility, OnlyHiddenVisibility, OnlyVisibleVisibility).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s'.", v, AllVisibility, OnlyHiddenVisibility, OnlyVisibleVisibility)).
		Validate(v)
}

// Validate validates ListEndpointVersionsRequest
func (r ListEndpointVersionsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
		"SortBy":        validation.Validate(r.SortBy),
		"SortOrder":     validation.Validate(r.SortOrder),
		"Show":          validation.Validate(r.Show),
	})
}

// Validate validates EndpointVersionRequest
func (r EndpointVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(r.APIEndpointID, validation.Required),
		"VersionNumber": validation.Validate(r.VersionNumber, validation.Required),
	})
}

// Validate validates GetEndpointVersionRequest
func (r GetEndpointVersionRequest) Validate() error {
	return EndpointVersionRequest(r).Validate()
}

// Validate validates CloneEndpointVersionRequest
func (r CloneEndpointVersionRequest) Validate() error {
	return EndpointVersionRequest(r).Validate()
}

// Validate validates DeleteEndpointVersionRequest
func (r DeleteEndpointVersionRequest) Validate() error {
	return EndpointVersionRequest(r).Validate()
}

// Validate validates UpdateEndpointVersionRequest
func (u UpdateEndpointVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"APIEndpointID": validation.Validate(u.APIEndpointID, validation.Required),
		"VersionNumber": validation.Validate(u.VersionNumber, validation.Required),
		"Body":          validation.Validate(u.Body, validation.Required),
	})
}

// Validate validates UpdateEndpointVersionRequestBody
func (e UpdateEndpointVersionRequestBody) Validate() error {
	return validation.Errors{
		"APIResources":               validation.Validate(e.APIResources),
		"APIVersionInfo":             validation.Validate(e.APIVersionInfo),
		"APIEndpointID":              validation.Validate(e.APIEndpointID, validation.Required),
		"APIEndpointName":            validation.Validate(e.APIEndpointName, validation.Required),
		"ConsumeType":                validation.Validate(e.ConsumeType),
		"ContractID":                 validation.Validate(e.ContractID, validation.Required),
		"GroupID":                    validation.Validate(e.GroupID, validation.Required),
		"APIEndpointHosts":           validation.Validate(e.APIEndpointHosts, validation.Required, validation.Length(1, 0)),
		"SecurityScheme":             validation.Validate(e.SecurityScheme),
		"AkamaiSecurityRestrictions": validation.Validate(e.AkamaiSecurityRestrictions),
		"APIEndpointScheme":          validation.Validate(e.APIEndpointScheme),
	}.Filter()
}

// Validate validates ConsumeType
func (c ConsumeType) Validate() error {
	return validation.In(ConsumeTypeJSON, ConsumeTypeXML, ConsumeTypeJSONXML, ConsumeTypeAny, ConsumeTypeUrlencoded, ConsumeTypeJSONUrlencoded, ConsumeTypeXMLUrlencoded, ConsumeTypeJSONXMLUrlencoded, ConsumeTypeNone).
		Error(fmt.Sprintf("value '%v' in invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'", c, ConsumeTypeJSON, ConsumeTypeXML, ConsumeTypeJSONXML, ConsumeTypeAny, ConsumeTypeUrlencoded, ConsumeTypeJSONUrlencoded, ConsumeTypeXMLUrlencoded, ConsumeTypeJSONXMLUrlencoded, ConsumeTypeNone)).
		Validate(c)
}

// Validate validates APIEndpointScheme
func (s APIEndpointScheme) Validate() error {
	return validation.In(APIEndpointSchemeHTTP, APIEndpointSchemeHTTPS, APIEndpointSchemeHTTPHTTPS).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s', '%s'", s, APIEndpointSchemeHTTP, APIEndpointSchemeHTTPS, APIEndpointSchemeHTTPHTTPS)).
		Validate(s)
}

// Validate validates APIEndpointScheme
func (s APISource) Validate() error {
	return validation.In(APISourceUser, APISourceAPIDiscovery).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s'", s, APISourceUser, APISourceAPIDiscovery)).
		Validate(s)
}

// Validate validates SourceType
func (s SourceType) Validate() error {
	return validation.In(SourceTypeSwagger, SourceTypeRaml).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", s, SourceTypeSwagger, SourceTypeRaml)).
		Validate(s)
}

// Validate validates Source
func (s Source) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(s.Type),
	}.Filter()
}

func (a *apidefinitions) ListEndpointVersions(ctx context.Context, params ListEndpointVersionsRequest) (*ListEndpointVersionsResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("ListEndpointVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListEndpointVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions", params.APIEndpointID))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrListEndpointVersions, err)
	}

	query := uri.Query()
	if params.Page != 0 {
		query.Add("page", strconv.FormatInt(params.Page, 10))
	}
	if params.PageSize != 0 {
		query.Add("pageSize", strconv.FormatInt(params.PageSize, 10))
	}
	if params.SortBy != "" {
		query.Add("sortBy", string(params.SortBy))
	}
	if params.SortOrder != "" {
		query.Add("sortOrder", string(params.SortOrder))
	}
	if params.Show != "" {
		query.Add("show", string(params.Show))
	}
	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListEndpointVersions, err)
	}

	var result ListEndpointVersionsResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListEndpointVersions, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListEndpointVersions, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) GetEndpointVersion(ctx context.Context, params GetEndpointVersionRequest) (*GetEndpointVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("GetEndpointVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEndpointVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d/resources-detail", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEndpointVersion, err)
	}

	var result GetEndpointVersionResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEndpointVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEndpointVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) CloneEndpointVersion(ctx context.Context, params CloneEndpointVersionRequest) (*CloneEndpointVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("CloneEndpointVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCloneEndpointVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d/cloneVersion", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCloneEndpointVersion, err)
	}

	var result CloneEndpointVersionResponse
	resp, err := a.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCloneEndpointVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCloneEndpointVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) UpdateEndpointVersion(ctx context.Context, params UpdateEndpointVersionRequest) (*UpdateEndpointVersionResponse, error) {
	logger := a.Log(ctx)
	logger.Debug("UpdateEndpointVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateEndpointVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateEndpointVersion, err)
	}

	var result UpdateEndpointVersionResponse
	resp, err := a.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEndpointVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateEndpointVersion, a.Error(resp))
	}

	return &result, nil
}

func (a *apidefinitions) DeleteEndpointVersion(ctx context.Context, params DeleteEndpointVersionRequest) error {
	logger := a.Log(ctx)
	logger.Debug("DeleteEndpointVersion")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteEndpointVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/api-definitions/v2/endpoints/%d/versions/%d", params.APIEndpointID, params.VersionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteEndpointVersion, err)
	}

	resp, err := a.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteEndpointVersion, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteEndpointVersion, a.Error(resp))
	}

	return nil
}
