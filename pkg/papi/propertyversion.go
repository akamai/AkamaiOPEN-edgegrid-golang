package papi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PropertyVersions contains operations available on PropertyVersions resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#propertyversionsgroup
	PropertyVersions interface {
		// GetPropertyVersions fetches available property versions
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyversions
		GetPropertyVersions(context.Context, GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error)

		// GetPropertyVersion fetches specific property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyversion
		GetPropertyVersion(context.Context, GetPropertyVersionRequest) (*GetPropertyVersionsResponse, error)

		// CreatePropertyVersion creates a new property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postpropertyversions
		CreatePropertyVersion(context.Context, CreatePropertyVersionRequest) (*CreatePropertyVersionResponse, error)

		// GetLatestVersion fetches latest property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getlatestversion
		GetLatestVersion(context.Context, GetLatestVersionRequest) (*GetPropertyVersionsResponse, error)

		// GetAvailableBehaviors fetches a list of behaviors applied to property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getavailablebehaviors
		GetAvailableBehaviors(context.Context, GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error)

		// GetAvailableCriteria fetches a list of criteria applied to property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getavailablecriteria
		GetAvailableCriteria(context.Context, GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error)
	}

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
	}

	// PropertyVersionItems contains collection of property version details
	PropertyVersionItems struct {
		Items []PropertyVersionGetItem `json:"items"`
	}

	// PropertyVersionGetItem contains detailed information about specific property version returned in GET
	PropertyVersionGetItem struct {
		Etag             string `json:"etag"`
		Note             string `json:"note"`
		ProductID        string `json:"productId"`
		ProductionStatus string `json:"productionStatus"`
		PropertyVersion  int    `json:"propertyVersion"`
		RuleFormat       string `json:"ruleFormat"`
		StagingStatus    string `json:"stagingStatus"`
		UpdatedByUser    string `json:"updatedByUser"`
		UpdatedDate      string `json:"updatedDate"`
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
		CreateFromVersionEtag string `json:"createFromVersionEtag"`
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

	// GetFeaturesRequest contains path and query params required to fetch both available behaviors and available criteria for a property
	GetFeaturesRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		GroupID         string
	}

	// AvailableFeature represents details of a single feature (behavior or criteria available for selected property version
	AvailableFeature struct {
		Name       string `json:"name"`
		SchemaLink string `json:"schemaLink"`
	}

	// GetFeaturesCriteriaResponse contains response received when fetching both available behaviors and available criteria for a property
	GetFeaturesCriteriaResponse struct {
		ContractID         string                `json:"contractId"`
		GroupID            string                `json:"groupId"`
		ProductID          string                `json:"productId"`
		RuleFormat         string                `json:"ruleFormat"`
		AvailableBehaviors AvailableFeatureItems `json:"availableBehaviors"`
	}

	// AvailableFeatureItems contains a slice of AvailableFeature items
	AvailableFeatureItems struct {
		Items []AvailableFeature `json:"items"`
	}
)

const (
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
	return validation.Errors{
		"PropertyID": validation.Validate(v.PropertyID, validation.Required),
		"Version":    validation.Validate(v.Version),
	}.Filter()
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

// Validate validates GetFeaturesRequest
func (v GetFeaturesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(v.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(v.PropertyVersion, validation.Required),
	}.Filter()
}

// GetPropertyVersions returns list of property versions for give propertyID, contractID and groupID
func (p *papi) GetPropertyVersions(ctx context.Context, params GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
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
		return nil, fmt.Errorf("failed to create getpropertyversions request: %w", err)
	}

	var versions GetPropertyVersionsResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("getpropertyversions request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &versions, nil
}

// GetLatestVersion returns either the latest property version overall, or the latest ACTIVE version on production or staging network
func (p *papi) GetLatestVersion(ctx context.Context, params GetLatestVersionRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
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
		return nil, fmt.Errorf("failed to create getlatestversion request: %w", err)
	}

	var version GetPropertyVersionsResponse
	resp, err := p.Exec(req, &version)
	if err != nil {
		return nil, fmt.Errorf("getlatestversion request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}
	return &version, nil
}

// GetPropertyVersion returns property version with provided version number
func (p *papi) GetPropertyVersion(ctx context.Context, params GetPropertyVersionRequest) (*GetPropertyVersionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
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
		return nil, fmt.Errorf("failed to create getpropertyversion request: %w", err)
	}

	var versions GetPropertyVersionsResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("getpropertyversion request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &versions, nil
}

// CreatePropertyVersion creates a new property version and returns location and number for the new version
func (p *papi) CreatePropertyVersion(ctx context.Context, request CreatePropertyVersionRequest) (*CreatePropertyVersionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
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
		return nil, fmt.Errorf("failed to create createpropertyversion request: %w", err)
	}

	var version CreatePropertyVersionResponse
	resp, err := p.Exec(req, &version)
	if err != nil {
		return nil, fmt.Errorf("createpropertyversion request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, session.NewAPIError(resp, logger)
	}
	propertyVersion, err := tools.FetchIDFromLocation(version.VersionLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrInvalidLocation, err.Error())
	}
	versionNumber, err := strconv.Atoi(propertyVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s", tools.ErrInvalidLocation, "version should be a number", propertyVersion)
	}
	version.PropertyVersion = versionNumber
	return &version, nil
}

// GetAvailableBehaviors lists available behaviors for given property version
func (p *papi) GetAvailableBehaviors(ctx context.Context, params GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAvailableBehaviors")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/available-behaviors?contractId=%s&groupId=%s",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getavailablebehaviors request: %w", err)
	}

	var versions GetFeaturesCriteriaResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("getavailablebehaviors request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &versions, nil
}

// GetAvailableCriteria lists available criteria for given property version
func (p *papi) GetAvailableCriteria(ctx context.Context, params GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAvailableCriteria")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/available-criteria?contractId=%s&groupId=%s",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getavailablecriteria request: %w", err)
	}

	var versions GetFeaturesCriteriaResponse
	resp, err := p.Exec(req, &versions)
	if err != nil {
		return nil, fmt.Errorf("getavailablecriteria request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &versions, nil
}
