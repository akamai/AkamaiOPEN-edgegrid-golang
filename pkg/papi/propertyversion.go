package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPropertyVersionsRequest contains path and query params used for listing property versions
	GetPropertyVersionsRequest struct {
		PropertyID string
		ContractID string
		GroupID    string
		Limit      int
		Offset     int
	}

	// GetPropertyVersionsResponse contains GET response returned while fetching property versions or specific version
	GetPropertyVersionsResponse struct {
		PropertyID   string               `json:"propertyId"`
		PropertyName string               `json:"propertyName"`
		AccountID    string               `json:"accountId"`
		ContractID   string               `json:"contractId"`
		GroupID      string               `json:"groupId"`
		AssetID      string               `json:"assetId"`
		Versions     PropertyVersionItems `json:"versions"`
		Version      PropertyVersionGetItem
	}

	// PropertyVersionItems contains collection of property version details
	PropertyVersionItems struct {
		Items []PropertyVersionGetItem `json:"items"`
	}

	// PropertyVersionGetItem contains detailed information about specific property version returned in GET
	PropertyVersionGetItem struct {
		Etag             string        `json:"etag"`
		Note             string        `json:"note"`
		ProductID        string        `json:"productId"`
		ProductionStatus VersionStatus `json:"productionStatus"`
		PropertyVersion  int           `json:"propertyVersion"`
		RuleFormat       string        `json:"ruleFormat"`
		StagingStatus    VersionStatus `json:"stagingStatus"`
		UpdatedByUser    string        `json:"updatedByUser"`
		UpdatedDate      string        `json:"updatedDate"`
	}

	// GetPropertyVersionRequest contains path and query params used for fetching specific property version
	GetPropertyVersionRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		GroupID         string
	}

	// CreatePropertyVersionRequest contains path and query params, as well as request body required to execute POST /versions request
	CreatePropertyVersionRequest struct {
		PropertyID string
		ContractID string
		GroupID    string
		Version    PropertyVersionCreate
	}

	// PropertyVersionCreate contains request body used in POST /versions request
	PropertyVersionCreate struct {
		CreateFromVersion     int    `json:"createFromVersion"`
		CreateFromVersionEtag string `json:"createFromVersionEtag,omitempty"`
	}

	// CreatePropertyVersionResponse contains a link returned after creating new property version and version number of this version
	CreatePropertyVersionResponse struct {
		VersionLink     string `json:"versionLink"`
		PropertyVersion int
	}

	// GetLatestVersionRequest contains path and query params required to fetch latest property version
	GetLatestVersionRequest struct {
		PropertyID  string
		ActivatedOn string
		ContractID  string
		GroupID     string
	}

	// GetAvailableItemsRequest contains path and query params required to fetch available behaviors or criteria for a property
	GetAvailableItemsRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		GroupID         string
	}

	// GetAvailableBehaviorsRequest contains path and query params required to fetch available behaviors for a property
	GetAvailableBehaviorsRequest GetAvailableItemsRequest

	// GetAvailableCriteriaRequest contains path and query params required to fetch available criteria for a property
	GetAvailableCriteriaRequest GetAvailableItemsRequest

	// GetBehaviorsResponse represents a response object returned by GetAvailableBehaviors
	GetBehaviorsResponse AvailableBehaviorsResponse

	// GetCriteriaResponse represents a response object returned by GetAvailableCriteria
	GetCriteriaResponse AvailableCriteriaResponse

	// VersionStatus represents ProductionVersion and StagingVersion of a Version struct
	VersionStatus string

	// ListAvailableIncludesRequest contains path and query params required to fetch list of available includes
	ListAvailableIncludesRequest ListAvailableReferencedIncludesRequest

	// ListReferencedIncludesRequest contains path and query params required to fetch  list of referenced includes
	ListReferencedIncludesRequest ListAvailableReferencedIncludesRequest

	//ListAvailableReferencedIncludesRequest common request struct for ListReferencedIncludesRequest and ListAvailableIncludesRequest
	ListAvailableReferencedIncludesRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		GroupID         string
	}

	// ListAvailableIncludesResponse contains response received when fetching list of available includes
	ListAvailableIncludesResponse struct {
		AvailableIncludes []ExternalIncludeData
	}

	// ListReferencedIncludesResponse contains response received when fetching list of referenced includes
	// The response from the API is a map, but we convert it to the array for better usability.
	ListReferencedIncludesResponse struct {
		Includes IncludeItems `json:"includes"`
	}

	// ExternalIncludeData contains data for a specific include from AvailableIncludes
	ExternalIncludeData struct {
		IncludeID   string      `json:"id"`
		IncludeName string      `json:"name"`
		IncludeType IncludeType `json:"includeType"`
		FileName    string      `json:"fileName"`
		ProductName string      `json:"productName"`
		RuleFormat  string      `json:"ruleFormat"`
	}
)

const (
	// VersionStatusActive const
	VersionStatusActive VersionStatus = "ACTIVE"
	// VersionStatusInactive const
	VersionStatusInactive VersionStatus = "INACTIVE"
	// VersionStatusPending const
	VersionStatusPending VersionStatus = "PENDING"
	// VersionStatusDeactivated const
	VersionStatusDeactivated VersionStatus = "DEACTIVATED"
	// VersionProduction const
	VersionProduction = "PRODUCTION"
	// VersionStaging const
	VersionStaging = "STAGING"
)

// Validate validates GetPropertyVersionsRequest
func (v GetPropertyVersionsRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(v.PropertyID, validation.Required),
	}.Filter()
}

// Validate validates GetPropertyVersionRequest
func (v GetPropertyVersionRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates CreatePropertyVersionRequest
func (v CreatePropertyVersionRequest) Validate() error {
	errs := validation.Errors{
		"PropertyID": validation.Validate(v.PropertyID, validation.Required),
		"Version":    validation.Validate(v.Version),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates PropertyVersionCreate
func (v PropertyVersionCreate) Validate() error {
	return validation.Errors{
		"CreateFromVersion": validation.Validate(v.CreateFromVersion, validation.Required),
	}.Filter()
}

// Validate validates GetLatestVersionRequest
func (v GetLatestVersionRequest) Validate() error {
	return validation.Errors{
		"PropertyID":  validation.Validate(v.PropertyID, validation.Required),
		"ActivatedOn": validation.Validate(v.ActivatedOn, validation.In(VersionProduction, VersionStaging)),
	}.Filter()
}

// Validate validates GetAvailableBehaviorsRequest
func (v GetAvailableBehaviorsRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates GetAvailableCriteriaRequest
func (v GetAvailableCriteriaRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates ListAvailableIncludesRequest
func (v ListAvailableIncludesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates ListReferencedIncludesRequest
func (v ListReferencedIncludesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
		"GroupID":         validation.Validate(v.GroupID, validation.Required),
		"ContractID":      validation.Validate(v.ContractID, validation.Required),
	}.Filter()
}

var (
	// ErrGetPropertyVersions represents error when fetching property versions fails
	ErrGetPropertyVersions = errors.New("fetching property versions")
	// ErrGetPropertyVersion represents error when fetching property version fails
	ErrGetPropertyVersion = errors.New("fetching property version")
	// ErrGetLatestVersion represents error when fetching latest property version fails
	ErrGetLatestVersion = errors.New("fetching latest property version")
	// ErrCreatePropertyVersion represents error when creating property version fails
	ErrCreatePropertyVersion = errors.New("creating property version")
	// ErrGetAvailableBehaviors represents error when fetching available behaviors fails
	ErrGetAvailableBehaviors = errors.New("fetching available behaviors")
	// ErrGetAvailableCriteria represents error when fetching available criteria fails
	ErrGetAvailableCriteria = errors.New("fetching available criteria")
	// ErrListAvailableIncludes represents error when fetching available includes
	ErrListAvailableIncludes = errors.New("fetching available includes")
	// ErrListReferencedIncludes represents error when fetching referenced includes
	ErrListReferencedIncludes = errors.New("fetching referenced includes")
)

func (p *papi) GetPropertyVersions(ctx context.Context, params GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersions, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyVersions")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ContractID,
		params.GroupID,
	)
	if params.Limit != 0 {
		getURL += fmt.Sprintf("&limit=%d", params.Limit)
	}
	if params.Offset != 0 {
		getURL += fmt.Sprintf("&offset=%d", params.Offset)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyVersions, err)
	}

	var versions GetPropertyVersionsResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyVersions, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyVersions, p.Error(resp))
	}

	return &versions, nil
}

func (p *papi) GetLatestVersion(ctx context.Context, params GetLatestVersionRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetLatestVersion, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetLatestVersion")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/latest?contractId=%s&groupId=%s",
		params.PropertyID,
		params.ContractID,
		params.GroupID,
	)
	if params.ActivatedOn != "" {
		getURL += fmt.Sprintf("&activatedOn=%s", params.ActivatedOn)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetLatestVersion, err)
	}

	var version GetPropertyVersionsResponse
	resp, err := p.Exec(req, &version)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetLatestVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetLatestVersion, p.Error(resp))
	}
	if len(version.Versions.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: latest version for PropertyID: %s", ErrGetLatestVersion, ErrNotFound, params.PropertyID)
	}
	version.Version = version.Versions.Items[0]
	return &version, nil
}

func (p *papi) GetPropertyVersion(ctx context.Context, params GetPropertyVersionRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersion, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyVersion")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d?contractId=%s&groupId=%s",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyVersion, err)
	}

	var versions GetPropertyVersionsResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyVersion, p.Error(resp))
	}
	if len(versions.Versions.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: Version %d for PropertyID: %s", ErrGetPropertyVersion, ErrNotFound, params.PropertyVersion, params.PropertyID)
	}
	versions.Version = versions.Versions.Items[0]
	return &versions, nil
}

func (p *papi) CreatePropertyVersion(ctx context.Context, request CreatePropertyVersionRequest) (*CreatePropertyVersionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreatePropertyVersion, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreatePropertyVersion")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions?contractId=%s&groupId=%s",
		request.PropertyID,
		request.ContractID,
		request.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePropertyVersion, err)
	}

	var version CreatePropertyVersionResponse
	resp, err := p.Exec(req, &version, request.Version)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePropertyVersion, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePropertyVersion, p.Error(resp))
	}
	propertyVersion, err := ResponseLinkParse(version.VersionLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreatePropertyVersion, ErrInvalidResponseLink, err)
	}
	versionNumber, err := strconv.Atoi(propertyVersion)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s: %s", ErrCreatePropertyVersion, ErrInvalidResponseLink, "version should be a number", propertyVersion)
	}
	version.PropertyVersion = versionNumber
	return &version, nil
}

func (p *papi) GetAvailableBehaviors(ctx context.Context, params GetAvailableBehaviorsRequest) (*GetBehaviorsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAvailableBehaviors, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetAvailableBehaviors")

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/versions/%d/available-behaviors", params.PropertyID, params.PropertyVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse uri: %s", ErrGetAvailableBehaviors, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAvailableBehaviors, err)
	}

	var behaviors GetBehaviorsResponse
	resp, err := p.Exec(req, &behaviors)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetAvailableBehaviors, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAvailableBehaviors, p.Error(resp))
	}

	return &behaviors, nil
}

func (p *papi) GetAvailableCriteria(ctx context.Context, params GetAvailableCriteriaRequest) (*GetCriteriaResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetAvailableCriteria, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetAvailableCriteria")

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/versions/%d/available-criteria", params.PropertyID, params.PropertyVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse uri: %s", ErrGetAvailableCriteria, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetAvailableCriteria, err)
	}

	var criteria GetCriteriaResponse
	resp, err := p.Exec(req, &criteria)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetAvailableCriteria, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetAvailableCriteria, p.Error(resp))
	}

	return &criteria, nil
}

func (p *papi) ListAvailableIncludes(ctx context.Context, params ListAvailableIncludesRequest) (*ListAvailableIncludesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListAvailableIncludes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListAvailableIncludes, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/versions/%d/external-resources", params.PropertyID, params.PropertyVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListAvailableIncludes, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListAvailableIncludes, err)
	}

	var result ListAvailableIncludesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListAvailableIncludes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListAvailableIncludes, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListReferencedIncludes(ctx context.Context, params ListReferencedIncludesRequest) (*ListReferencedIncludesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListReferencedIncludes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListReferencedIncludes, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/properties/%s/versions/%d/includes", params.PropertyID, params.PropertyVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListReferencedIncludes, err)
	}

	q := uri.Query()
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListReferencedIncludes, err)
	}

	var result ListReferencedIncludesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListReferencedIncludes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListReferencedIncludes, p.Error(resp))
	}

	return &result, nil
}

// UnmarshalJSON reads a ListAvailableIncludesResponse struct from its data argument and transform map of includes into array for better usability.
func (r *ListAvailableIncludesResponse) UnmarshalJSON(data []byte) error {
	var response struct {
		ExternalResources struct {
			ExternalIncludes map[string]ExternalIncludeData `json:"include"`
		} `json:"externalResources"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	r.AvailableIncludes = make([]ExternalIncludeData, 0, len(response.ExternalResources.ExternalIncludes))
	for _, include := range response.ExternalResources.ExternalIncludes {
		r.AvailableIncludes = append(r.AvailableIncludes, include)
	}

	return nil
}
