package reportinggroups

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
