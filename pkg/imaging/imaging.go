// Package imaging provides access to the Akamai Image & Video Manager APIs
package imaging

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Imaging is the api interface for Image and Video Manager
	Imaging interface {
		// Policies

		// ListPolicies lists all Policies for the given network and an account.
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policies
		ListPolicies(context.Context, ListPoliciesRequest) (*ListPoliciesResponse, error)

		// GetPolicy gets specific policy by PolicyID.
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policy
		GetPolicy(context.Context, GetPolicyRequest) (PolicyOutput, error)

		// UpsertPolicy creates or updates the configuration for a policy.
		//
		// See: https://techdocs.akamai.com/ivm/reference/put-policy
		UpsertPolicy(context.Context, UpsertPolicyRequest) (*PolicyResponse, error)

		// DeletePolicy deletes a policy.
		//
		// See: https://techdocs.akamai.com/ivm/reference/delete-policy
		DeletePolicy(context.Context, DeletePolicyRequest) (*PolicyResponse, error)

		// GetPolicyHistory retrieves history of changes for a policy.
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policy-history
		GetPolicyHistory(context.Context, GetPolicyHistoryRequest) (*GetPolicyHistoryResponse, error)

		// RollbackPolicy reverts a policy to its previous version and deploys it to the network.
		//
		// See: https://techdocs.akamai.com/ivm/reference/put-rollback-policy
		RollbackPolicy(ctx context.Context, request RollbackPolicyRequest) (*PolicyResponse, error)

		// PolicySets

		// ListPolicySets lists all PolicySets of specified type for the current account.
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policysets
		ListPolicySets(context.Context, ListPolicySetsRequest) ([]PolicySet, error)

		// GetPolicySet gets specific PolicySet by PolicySetID.
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policyset
		GetPolicySet(context.Context, GetPolicySetRequest) (*PolicySet, error)

		// CreatePolicySet creates configuration for an PolicySet.
		//
		// See: https://techdocs.akamai.com/ivm/reference/post-policyset
		CreatePolicySet(context.Context, CreatePolicySetRequest) (*PolicySet, error)

		// UpdatePolicySet creates configuration for an PolicySet.
		//
		// See: https://techdocs.akamai.com/ivm/reference/put-policyset
		UpdatePolicySet(context.Context, UpdatePolicySetRequest) (*PolicySet, error)

		// DeletePolicySet deletes configuration for an PolicySet.
		//
		// See: https://techdocs.akamai.com/ivm/reference/delete-policyset
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
