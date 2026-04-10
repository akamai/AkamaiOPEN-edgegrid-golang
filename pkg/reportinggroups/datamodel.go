package reportinggroups

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	_ validation.Validatable = CreateReportingGroupRequest{}
	_ validation.Validatable = GetReportingGroupsRequest{}
	_ validation.Validatable = UpdateReportingGroupRequest{}
	_ validation.Validatable = DeleteReportingGroupRequest{}
	_ validation.Validatable = ListProductsRequest{}
	_ validation.Validatable = GetReportingGroupsWaterMarkLimitsRequest{}
	_ validation.Validatable = GetCPCodesWaterMarkLimitsRequest{}
	_ validation.Validatable = GetCPCodeRequest{}
	_ validation.Validatable = UpdateCPCodeRequest{}
	_ validation.Validatable = CPCodeContract{}
	_ validation.Validatable = Product{}
	_ validation.Validatable = CPCodeTimeZone{}
)

type (
	// CreateReportingGroupRequest is the request body for creating a reporting group.
	CreateReportingGroupRequest struct {
		// A group that controls access to specific CP codes.
		AccessGroup AccessGroupModel `json:"accessGroup"`

		// A collection of contracts and CP codes assigned to the reporting group.
		Contracts []ContractCreateModel `json:"contracts"`

		// The descriptive label for the reporting group.
		ReportingGroupName string `json:"reportingGroupName"`
	}

	// AccessGroupModel is the model for the access group in the reporting group.
	AccessGroupModel struct {
		// Identifies the contract assigned to the access control group. Get and store this ID from your account using the Contracts API.
		ContractID string `json:"contractId"`

		// Identifies the access control group. Reporting groups may belong to many groups. If that happens, this member's value is null.
		GroupID *int64 `json:"groupId"`
	}

	// ContractCreateModel is the model for the contract in the reporting group.
	ContractCreateModel struct {
		// Identifies the contract assigned to the reporting group.
		ContractID string `json:"contractId"`

		// A collection of CP codes assigned to the reporting group.
		CpCodes []CpCodeCreateModel `json:"cpcodes"`
	}

	// CpCodeCreateModel is the model for the CP code in the reporting group.
	CpCodeCreateModel struct {
		// Identifies a CP code.
		CpCodeID int64 `json:"cpcodeId"`
	}

	// ContractModel is the response model for the contract in the reporting group.
	ContractModel struct {
		// Identifies the contract assigned to the reporting group.
		ContractID string `json:"contractId"`

		// A collection of CP codes assigned to the reporting group.
		CpCodes []CpCodeModel `json:"cpcodes"`
	}

	// CpCodeModel is the resposne model for the CP code in the reporting group.
	CpCodeModel struct {
		// Identifies a CP code.
		CpCodeID int64 `json:"cpcodeId"`

		// The descriptive label for the CP code.
		CpCodeName string `json:"cpcodeName"`
	}

	// CreateReportingGroupResponse is the response body for creating a reporting group.
	CreateReportingGroupResponse struct {
		// ReportingGroupItem contains detailed information about the created reporting group.
		ReportingGroupItem

		// ResourceLimits contains information about reporting groups limits.
		ResourceLimits ResourceLimitsMetadata
	}

	// ResourceLimitsMetadata is the model for reporting groups resource limits.
	ResourceLimitsMetadata struct {
		// Total number of reporting groups allowed.
		ReportingGroupsLimitTotal *int64

		// Number of remaining reporting groups that can be created.
		ReportingGroupsLimitRemaining *int64
	}

	// GetReportingGroupsRequest is the request body for getting reporting groups.
	GetReportingGroupsRequest struct {
		// The identifier for the reporting group.
		ReportingGroupID int64 `json:"reportingGroupId,omitempty"`
	}

	// GetReportingGroupResponse is the response body for getting a reporting group.
	GetReportingGroupResponse ReportingGroupItem

	// UpdateReportingGroupRequest is the request body for updating a reporting group.
	UpdateReportingGroupRequest struct {
		// The identifier for the reporting group.
		ReportingGroupID int64 `json:"reportingGroupId"`

		// The descriptive label for the reporting group.
		ReportingGroupName string `json:"reportingGroupName,omitempty"`

		// A collection of contracts and CP codes assigned to the reporting group.
		Contracts []ContractModel `json:"contracts,omitempty"`
	}

	// UpdateReportingGroupResponse is the response body for updating a reporting group.
	UpdateReportingGroupResponse ReportingGroupItem

	// DeleteReportingGroupRequest is the request body for deleting a reporting group.
	DeleteReportingGroupRequest struct {
		// The identifier for the reporting group.
		ReportingGroupID int64 `json:"reportingGroupId"`
	}

	// ListReportingGroupsRequest is the request for listing reporting groups.
	ListReportingGroupsRequest struct {
		// Identifies the contract to filter data by.
		ContractID string `json:"contractId,omitempty"`

		// Identifies the access group to filter data by.
		GroupID int64 `json:"groupId,omitempty"`

		// The name of the reporting group to filter data by.
		ReportingGroupName string `json:"reportingGroupName,omitempty"`

		// Identifies the CP code to filter data by.
		CpCodeID string `json:"cpcodeId,omitempty"`
	}

	// ListReportingGroupsResponse is the response body for listing reporting groups.
	ListReportingGroupsResponse struct {
		// A set of reporting groups available for your contract.
		Groups []ReportingGroupItem `json:"groups"`
	}

	// ReportingGroupItem is the model for a reporting group item in the list response.
	ReportingGroupItem struct {
		// A group that controls access to specific CP codes.
		AccessGroup AccessGroupModel `json:"accessGroup"`

		// A collection of contracts and CP codes assigned to the reporting group.
		Contracts []ContractModel `json:"contracts"`

		// Identifies the reporting group.
		ReportingGroupID int64 `json:"reportingGroupId"`

		// The descriptive label for the reporting group.
		ReportingGroupName string `json:"reportingGroupName"`
	}

	// ListProductsRequest is the request for listing products within a reporting group.
	ListProductsRequest struct {
		// The identifier for the reporting group.
		ReportingGroupID int64 `json:"reportingGroupId"`
	}

	// ListProductsResponse is the response body for listing products within a reporting group.
	ListProductsResponse struct {
		// A collection of products and services assigned to the reporting group.
		Products []Product `json:"products"`
	}

	// Product is the model for a product or service assigned to a reporting group.
	Product struct {
		// Identifies a product or service.
		ProductID string `json:"productId"`

		// The descriptive label for a product or service.
		ProductName string `json:"productName"`
	}

	// CPCodeContract contains contract data for a CP code.
	CPCodeContract struct {
		// ContractID identifies the contract assigned to the CP code.
		ContractID string `json:"contractId"`

		// Status is the status of the contract assigned to the CP code.
		Status string `json:"status,omitempty"`
	}

	// CPCodeTimeZone contains time zone data for a CP code.
	CPCodeTimeZone struct {
		// TimeZoneID identifies the time zone.
		TimeZoneID string `json:"timezoneId"`

		// TimeZoneValue is the GMT time zone value.
		TimeZoneValue string `json:"timezoneValue"`
	}

	// ListCPCodesRequest contains the optional query parameters for listing CP codes.
	ListCPCodesRequest struct {
		// ContractID identifies the contract to filter data by.
		ContractID string

		// GroupID identifies the access group to filter data by.
		GroupID string

		// ProductID identifies the product or service to filter data by.
		ProductID string

		// CPCodeName is the name of the CP code to filter data by.
		CPCodeName string
	}

	// ListCPCodesResponse is the response for listing CP codes.
	ListCPCodesResponse struct {
		// CPCodes is a collection of CP codes details available for your contract.
		CPCodes []CPCodeDetails `json:"cpcodes"`
	}

	// CPCodeDetails contains detailed information about a single CP code in the list response.
	CPCodeDetails struct {
		// AccessGroup is the group that controls access to specific CP codes.
		AccessGroup AccessGroupModel `json:"accessGroup"`

		// AccountID identifies an account assigned to the CP code.
		AccountID string `json:"accountId"`

		// Contracts provides detailed information about the contracts assigned to the CP code.
		Contracts []CPCodeContract `json:"contracts"`

		// CPCodeID identifies the CP code.
		CPCodeID int64 `json:"cpcodeId"`

		// CPCodeName is the descriptive label for the CP code.
		CPCodeName string `json:"cpcodeName"`

		// DefaultTimeZone is the default GMT time zone assigned to the CP code.
		DefaultTimeZone string `json:"defaultTimezone"`

		// OverrideTimeZone is the GMT time zone that overrides the default time zone.
		OverrideTimeZone CPCodeTimeZone `json:"overrideTimezone"`

		// Products is a collection of products and services assigned to the CP code.
		Products []Product `json:"products"`

		// Purgeable indicates whether you can purge the content cached by the CP code.
		Purgeable bool `json:"purgeable"`

		// Type indicates whether the CP code supports serving traffic from 'China', or 'Regular' traffic elsewhere.
		Type string `json:"type"`
	}

	// GetReportingGroupsWaterMarkLimitsRequest is the request body for getting reporting groups water-mark limits.
	GetReportingGroupsWaterMarkLimitsRequest struct {
		// ContractID is the identifier for the contract for which you want to check water-mark limits.
		ContractID string
	}

	// GetReportingGroupsWaterMarkLimitsResponse is the response body for getting reporting groups water-mark limits.
	GetReportingGroupsWaterMarkLimitsResponse struct {
		// CurrentCapacity is the current number of reporting groups.
		CurrentCapacity int `json:"currentCapacity"`

		// Limit is the number of allowed reporting groups.
		Limit int `json:"limit"`

		// LimitType identifies whether the limit is global or applies to an account or a contract.
		LimitType string `json:"limitType"`
	}

	// GetCPCodesWaterMarkLimitsRequest contains the parameters for getting CP codes water-mark limits.
	GetCPCodesWaterMarkLimitsRequest struct {
		// ContractID is the identifier for the contract for which you want to check water-mark limits.
		ContractID string
	}

	// GetCPCodeRequest contains the parameters for fetching a single CP code's details.
	GetCPCodeRequest struct {
		// CPCodeID identifies the CP code.
		CPCodeID int64
	}

	// GetCPCodeResponse is an alias for CPCodeDetails, returned by GetCPCode.
	GetCPCodeResponse = CPCodeDetails

	// UpdateCPCodeRequest contains the parameters for updating a specific CP code.
	UpdateCPCodeRequest struct {
		// CPCodeID identifies the CP code. Used as path parameter and included in the request body.
		CPCodeID int64 `json:"cpcodeId"`

		// CPCodeName is the descriptive label for the CP code.
		CPCodeName string `json:"cpcodeName"`

		// Purgeable indicates whether you can purge the content cached by the CP code.
		Purgeable *bool `json:"purgeable,omitempty"`

		// OverrideTimeZone is the GMT time zone that overrides the default time zone.
		OverrideTimeZone *CPCodeTimeZone `json:"overrideTimezone,omitempty"`

		// Contracts provides contract status information (required by API, pass through from GET response).
		Contracts []CPCodeContract `json:"contracts"`

		// Products is a collection of products and services assigned to the CP code (required by API, pass through from GET response).
		Products []Product `json:"products"`
	}

	// UpdateCPCodeResponse is an alias for CPCodeDetails, returned by UpdateCPCode.
	UpdateCPCodeResponse = CPCodeDetails

	// GetCPCodesWaterMarkLimitsResponse is the response for getting CP codes water-mark limits.
	GetCPCodesWaterMarkLimitsResponse struct {
		// CurrentCapacity is the current number of CP codes.
		CurrentCapacity int `json:"currentCapacity"`

		// Limit is the number of allowed CP codes.
		Limit int `json:"limit"`

		// LimitType identifies whether the limit is global or applies to an account or a contract.
		LimitType string `json:"limitType"`
	}
)

// Validate validates CreateReportingGroupRequest.
func (r CreateReportingGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ReportingGroupName": validation.Validate(r.ReportingGroupName, validation.Required),
		"AccessGroup":        validation.Validate(r.AccessGroup.ContractID, validation.Required),
		"Contracts":          validation.Validate(r.Contracts, validation.Required, validation.Length(1, 1)),
	})
}

// Validate validates GetReportingGroupsRequest.
func (r GetReportingGroupsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ReportingGroupID": validation.Validate(r.ReportingGroupID, validation.Required),
	})
}

// Validate validates UpdateReportingGroupRequest.
func (r UpdateReportingGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ReportingGroupID":   validation.Validate(r.ReportingGroupID, validation.Required),
		"ReportingGroupName": validation.Validate(r.ReportingGroupName, validation.Required),
		"Contracts":          validation.Validate(r.Contracts, validation.Required, validation.Length(1, 1)),
	})
}

// Validate validates DeleteReportingGroupRequest.
func (r DeleteReportingGroupRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ReportingGroupID": validation.Validate(r.ReportingGroupID, validation.Required),
	})
}

// Validate validates ListProductsRequest.
func (r ListProductsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ReportingGroupID": validation.Validate(r.ReportingGroupID, validation.Required),
	})
}

// Validate validates GetReportingGroupsWaterMarkLimitsRequest.
func (r GetReportingGroupsWaterMarkLimitsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID": validation.Validate(r.ContractID, validation.Required),
	})
}

// Validate validates GetCPCodesWaterMarkLimitsRequest.
func (r GetCPCodesWaterMarkLimitsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ContractID": validation.Validate(r.ContractID, validation.Required),
	})
}

// Validate validates GetCPCodeRequest.
func (r GetCPCodeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CPCodeID": validation.Validate(r.CPCodeID, validation.Required),
	})
}

// Validate validates UpdateCPCodeRequest.
func (r UpdateCPCodeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CPCodeID":         validation.Validate(r.CPCodeID, validation.Required),
		"CPCodeName":       validation.Validate(r.CPCodeName, validation.Required),
		"Contracts":        validation.Validate(r.Contracts, validation.Required),
		"Products":         validation.Validate(r.Products, validation.Required),
		"OverrideTimeZone": validation.Validate(r.OverrideTimeZone),
	})
}

// Validate validates CPCodeContract.
func (c CPCodeContract) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(c.ContractID, validation.Required),
	}.Filter()
}

// Validate validates Product.
func (p Product) Validate() error {
	return validation.Errors{
		"ProductID": validation.Validate(p.ProductID, validation.Required),
	}.Filter()
}

// Validate validates CPCodeTimeZone.
func (tz CPCodeTimeZone) Validate() error {
	return validation.Errors{
		"TimeZoneID": validation.Validate(tz.TimeZoneID, validation.Required),
	}.Filter()
}
