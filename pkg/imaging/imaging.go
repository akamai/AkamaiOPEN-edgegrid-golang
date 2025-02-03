// Package imaging provides access to the Akamai Image & Video Manager APIs
package imaging

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Imaging is the api interface for Image and Video Manager
	Imaging interface {
		// Policies

		// ListPolicies lists all Policies for the given network and an account
		// See: https://techdocs.akamai.com/ivm/reference/get-policies
		ListPolicies(context.Context, ListPoliciesRequest) (*ListPoliciesResponse, error)

		// GetPolicy gets specific policy by PolicyID
		GetPolicy(context.Context, GetPolicyRequest) (PolicyOutput, error)

		// UpsertPolicy creates or updates the configuration for a policy
		UpsertPolicy(context.Context, UpsertPolicyRequest) (*PolicyResponse, error)

		// DeletePolicy deletes a policy
		DeletePolicy(context.Context, DeletePolicyRequest) (*PolicyResponse, error)

		// GetPolicyHistory retrieves history of changes for a policy
		GetPolicyHistory(context.Context, GetPolicyHistoryRequest) (*GetPolicyHistoryResponse, error)

		// RollbackPolicy reverts a policy to its previous version and deploys it to the network
		RollbackPolicy(ctx context.Context, request RollbackPolicyRequest) (*PolicyResponse, error)

		// PolicySets

		// ListPolicySets lists all PolicySets of specified type for the current account
		ListPolicySets(context.Context, ListPolicySetsRequest) ([]PolicySet, error)

		// GetPolicySet gets specific PolicySet by PolicySetID
		GetPolicySet(context.Context, GetPolicySetRequest) (*PolicySet, error)

		// CreatePolicySet creates configuration for an PolicySet
		CreatePolicySet(context.Context, CreatePolicySetRequest) (*PolicySet, error)

		// UpdatePolicySet creates configuration for an PolicySet
		UpdatePolicySet(context.Context, UpdatePolicySetRequest) (*PolicySet, error)

		// DeletePolicySet deletes configuration for an PolicySet
		DeletePolicySet(context.Context, DeletePolicySetRequest) error
	}

	imaging struct {
		session.Session
	}

	// Option defines an Image and Video Manager option
	Option func(*imaging)

	// ClientFunc is a Image and Video Manager client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Imaging
)

// Client returns a new Image and Video Manager Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Imaging {
	c := &imaging{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
