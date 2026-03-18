// Package reportinggroups provides access to the Reporting Groups API.
package reportinggroups

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

type (
	// ReportingGroups is the interface for the Reporting Groups API.
	ReportingGroups interface {
		// CreateReportingGroup creates a new reporting group.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/post-reporting-groups
		CreateReportingGroup(ctx context.Context, params CreateReportingGroupRequest) (*CreateReportingGroupResponse, error)

		// GetReportingGroup retrieves a reporting group by its ID.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/get-reporting-group
		GetReportingGroup(ctx context.Context, params GetReportingGroupsRequest) (*GetReportingGroupResponse, error)

		// ListReportingGroups lists detailed information about reporting groups available for your account and contract.
		// Optionally filtered by ContractID, GroupID, ReportingGroupName, or CPCodeID.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/get-reporting-groups
		ListReportingGroups(ctx context.Context, params ListReportingGroupsRequest) (*ListReportingGroupsResponse, error)

		// UpdateReportingGroup updates an existing reporting group.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/put-reporting-group
		UpdateReportingGroup(ctx context.Context, params UpdateReportingGroupRequest) (*UpdateReportingGroupResponse, error)

		// DeleteReportingGroup deletes a reporting group by its ID.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/delete-reporting-group
		DeleteReportingGroup(ctx context.Context, params DeleteReportingGroupRequest) error

		// ListProducts lists products and services assigned to a specific reporting group.
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/get-reporting-group-products
		ListProducts(ctx context.Context, params ListProductsRequest) (*ListProductsResponse, error)
	}

	reportinggroups struct {
		session.Session
	}

	// Option is a function that configures the Reporting Groups client.
	Option func(*reportinggroups)
)

// Client creates a new ReportingGroups client.
func Client(sess session.Session, opts ...Option) ReportingGroups {
	c := &reportinggroups{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
