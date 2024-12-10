package netstorage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// StorageGroup represents a NetStorage storage group resource
	StorageGroup struct {
		ContractID          string                        `json:"contractId,omitempty"`
		StorageGroupID      int                           `json:"storageGroupId,omitempty"`
		StorageGroupName    string                        `json:"storageGroupName,omitempty"`
		StorageGroupType    StorageGroupType              `json:"storageGroupType,omitempty"`
		StorageGroupPurpose StorageGroupPurpose           `json:"storageGroupPurpose,omitempty"`
		DomainPrefix        string                        `json:"domainPrefix,omitempty"`
		AsperaEnabled       bool                          `json:"asperaEnabled,omitempty"`
		PciEnabled          bool                          `json:"pciEnabled,omitempty"`
		EstimatedUsageGB    float64                       `json:"estimatedUsageGB,omitempty"`
		AllowEdit           bool                          `json:"allowEdit,omitempty"`
		ProvisionStatus     StorageGroupProvisionStatus   `json:"provisionStatus,omitempty"`
		PropagationStatus   StorageGroupPropagationStatus `json:"propagationStatus,omitempty"`
		LastModifiedBy      string                        `json:"lastModifiedBy,omitempty"`
		LastModifiedDate    string                        `json:"lastModifiedDate,omitempty"`
		CPCodes             []*StorageGroupCPCode         `json:"cpcodes,omitempty"`
	}

	// ListStorageGroupsResponse is the get activation response
	ListStorageGroupsResponse struct {
		Items []*StorageGroup `json:"items"`
	}

	// ListStorageGroupsRequest is the list storage groups request
	ListStorageGroupsRequest struct {
		CPCodeID            string
		StorageGroupPurpose string
	}

	// GetStorageGroupRequest is the get storage group request
	GetStorageGroupRequest struct {
		StorageGroupID int
	}

	// GetStorageGroupResponse is the get storage group response
	GetStorageGroupResponse StorageGroup

	// StorageGroupType is a storage group type value
	StorageGroupType string

	// StorageGroupPurpose is a storage group purpose value
	StorageGroupPurpose string

	// StorageGroupProvisionStatus is a storage group provision status value
	StorageGroupProvisionStatus string

	PropagationStatus string

	RequestURICaseConversion string

	QueryStringConversionMode string

	QueryStringConversion struct {
		QueryStringConversionMode QueryStringConversionMode `json:"queryStringConversionMode"`
	}

	Encoding string

	// The type of encoding to use when displaying paths in XML-aware contexts.
	EncodingConfig struct {
		Encoding        Encoding `json:"encoding"`
		EnforceEncoding bool     `json:"enforceEncoding"`
	}

	DirListingSearchOn404 string

	DirListing struct {
		IndexFileName string                `json:"indexFileName"`
		MaxListSize   int                   `json:"maxListSize"`
		SearchOn404   DirListingSearchOn404 `json:"searchOn404"`
	}

	PathCheckAndConversion string

	// StorageGroupPropagationStatus is a storage group propagation status value
	StorageGroupPropagationStatus struct {
		Status PropagationStatus `json:"status"`
	}

	StorageGroupCPCode struct {
		CPCodeID                 int                      `json:"cpcodeId"`
		DownloadSecurity         string                   `json:"downloadSecurity"`
		UseSSL                   bool                     `json:"useSsl"`
		ServeFromZip             bool                     `json:"serveFromZip"`
		SendHash                 bool                     `json:"sendHash"`
		QuickDelete              bool                     `json:"quickDelete"`
		NumberOfFiles            int                      `json:"numberOfFiles"`
		NumberOfBytes            int                      `json:"numberOfBytes"`
		LastChangesPropagated    bool                     `json:"lastChangesPropagated"`
		RequestURICaseConversion RequestURICaseConversion `json:"requestUriCaseConversion"`
		QueryStringConversion    QueryStringConversion    `json:"queryStringConversion"`
		PathCheckAndConversion   PathCheckAndConversion   `json:"pathCheckAndConversion"`
		EncodingConfig           EncodingConfig           `json:"encodingConfig"`
		DirListing               DirListing               `json:"dirListing"`
		LastModifiedBy           string                   `json:"lastModifiedBy"`
		LastModifiedDate         string                   `json:"lastModifiedDate"`
	}
)

const (
	// StorageGroupTypeObjectStore is used for ...
	StorageGroupTypeObjectStore StorageGroupType = "OBJECTSTORE"

	// A universal NetStorage storage group.
	StorageGroupPurposeNetStorage StorageGroupPurpose = "NETSTORAGE"

	// Provisioned for universal streaming over edge servers.
	StorageGroupPurposeEdgestream StorageGroupPurpose = "EDGESTREAM"

	// Provisioned for iPhone (HLS) 2.1 streaming over edge servers.
	StorageGroupPurposeEdgestreamIphone StorageGroupPurpose = "EDGESTREAM_IPHONE"

	// Provisioned for Adaptive Media Delivery over edge servers.
	StorageGroupPurposeAdaptiveEdge StorageGroupPurpose = "ADAPTIVEEDGE"

	// Provisioned for Ad Insertion use.
	StorageGroupPurposeAdInsertion StorageGroupPurpose = "AD_INSERTION"

	// For use with Media Services on Demand: Content Preparation.
	StorageGroupPurposeContentPreparation StorageGroupPurpose = "CONTENT_PREPARATION"

	// Provisioned as an Origin for Media Services Live.
	StorageGroupPurposeMSLOrigin StorageGroupPurpose = "MSL_ORIGIN"

	//
	StorageGroupPurposeFEO StorageGroupPurpose = "FEO"

	// Indicates the group hasn't been requested for provisioning.
	StorageGroupProvisionStatusNotProvisioned StorageGroupProvisionStatus = "NOT_PROVISIONED"

	// Indicates the group is ready for use.
	StorageGroupProvisionStatusProvisioned StorageGroupProvisionStatus = "PROVISIONED"

	// Indicates deprovisioning has been requested, but the group is still accessible.
	StorageGroupProvisionStatusMarkedForDeprovisioning StorageGroupProvisionStatus = "MARKED_FOR_DEPROVISIONING"

	// Indicates deprovisioning has completed, and the group is no longer available.
	StorageGroupProvisionStatusDeprovisioned StorageGroupProvisionStatus = "DEPROVISIONED"

	PropagationStatusPending PropagationStatus = "PENDING"

	PropagationStatusActive PropagationStatus = "ACTIVE"

	RequestURICaseConversionNoConversion RequestURICaseConversion = "NO_CONVERSION"

	RequestURICaseConversionConvertToUpper RequestURICaseConversion = "CONVERT_TO_UPPER"

	RequestURICaseConversionConvertToLower RequestURICaseConversion = "CONVERT_TO_LOWER"

	// Indicates case requirements are applied to support the Stream OS product.
	RequestURICaseConversionStreamOS RequestURICaseConversion = "STREAM_OS"

	// All query strings are stripped and ignored.
	QueryStringConversionModeStripAllIncomingQuery QueryStringConversionMode = "STRIP_ALL_INCOMING_QUERY"

	// The key is reviewed and stripped down.
	QueryStringConversionModeApplyACSQueryConversion QueryStringConversionMode = "APPLY_ACS_QUERY_CONVERSION"

	// The string is left as is.
	QueryStringConversionModeLeaveIncomingQueryAsIs QueryStringConversionMode = "LEAVE_INCOMING_QUERY_AS_IS"

	// Stops additional searches, and returns a 404 error.
	DirListingSearchOn404DoNotSearch DirListingSearchOn404 = "DO_NOT_SEARCH"

	// Looks for an explicit directory matching the path specified in the request
	DirListingSearchOn404LookForExplicitDirOnly DirListingSearchOn404 = "LOOK_FOR_EXPLICIT_DIR_ONLY"

	// Looks for both an explicit and implicit directory that may match a path defined in the request.
	DirListingSearchOn404LookForImplicitAndExplicitDir DirListingSearchOn404 = "LOOK_FOR_IMPLICIT_AND_EXPLICIT_DIR"

	PathCheckAndConversionDisallowAllImproperPaths PathCheckAndConversion = "DISALLOW_ALL_IMPROPER_PATHS"

	PathCheckAndConversionDisallowImproperPathsOnUploadOnly PathCheckAndConversion = "DISALLOW_IMPROPER_PATHS_ON_UPLOAD_ONLY"

	PathCheckAndConversionTranslateToCanonical PathCheckAndConversion = "TRANSLATE_TO_CANONICAL"

	PathCheckAndConversionDoNotCheckPaths PathCheckAndConversion = "DO_NOT_CHECK_PATHS"

	// 8-bit, single-byte coded graphic character sets
	EncodingISO88591 Encoding = "ISO_8859_1"

	// Variable length, 8-bit code units via UTF-8 character encoding.
	EncodingUTF8 Encoding = "UTF_8"
)

// Validate validates GetStorageGroupRequest
func (v GetStorageGroupRequest) Validate() error {
	return validation.Errors{
		"StorageGroupID": validation.Validate(v.StorageGroupID, validation.Required),
	}.Filter()
}

// Validate validates ListStorageGroupsRequest
func (v ListStorageGroupsRequest) Validate() error {
	return validation.Errors{}.Filter()
}

var (
	// ErrListStorageGroups represents error when fetching storage groups fails
	ErrListStorageGroups = errors.New("fetching storage groups")
	// ErrGetStorageGroup represents error when fetching storage group fails
	ErrGetStorageGroup = errors.New("fetching storage group")
)

func (p *netstorage) ListStorageGroups(ctx context.Context, params ListStorageGroupsRequest) (*ListStorageGroupsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListStorageGroups, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("ListStorageGroups")

	uri, err := url.Parse(
		"/storage/v1/storage-groups",
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListStorageGroups, err)
	}
	q := uri.Query()
	if params.CPCodeID != "" {
		q.Add("cpcodeId", params.CPCodeID)
	}
	if params.StorageGroupPurpose != "" {
		q.Add("storageGroupPurpose", params.StorageGroupPurpose)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListStorageGroups, err)
	}

	var rval ListStorageGroupsResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListStorageGroups, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListStorageGroups, p.Error(resp))
	}

	return &rval, nil
}

func (p *netstorage) GetStorageGroup(ctx context.Context, params GetStorageGroupRequest) (*GetStorageGroupResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetStorageGroup, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetStorageGroup")

	uri := fmt.Sprintf(
		"/storage/v1/storage-groups/%d",
		params.StorageGroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetStorageGroup, err)
	}

	var rval GetStorageGroupResponse

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetStorageGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetStorageGroup, p.Error(resp))
	}

	return &rval, nil
}
