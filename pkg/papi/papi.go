// Package papi provides access to the Akamai Property APIs
package papi

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type (
	// PAPI is the papi api interface
	PAPI interface {
		// GetGroups provides a read-only list of groups, which may contain properties.
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getgroups
		GetGroups(context.Context) (*GetGroupsResponse, error)

		// GetContract provides a read-only list of contract names and identifiers
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcontracts
		GetContracts(context.Context) (*GetContractResponse, error)

		// CreateActivation creates a new activation or deactivation request
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postpropertyactivations
		CreateActivation(ctx context.Context, req CreateActivationRequest) (*CreateActivationResponse, error)

		// GetActivation gets details about an activation
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyactivation
		GetActivation(ctx context.Context, req *GetActivationRequest) (*GetActivationResponse, error)
	}

	papi struct {
		session.Session
	}
)

var (
	// UsePrefixes is set to tell the PAPI api to accept and return prefixed object idenfiers
	UsePrefixes = true
)

// New returns a new papi New instance with the specified controller
func New(sess session.Session) PAPI {
	return &papi{
		sess,
	}
}
