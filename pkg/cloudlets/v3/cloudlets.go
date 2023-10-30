// Package v3 provides access to the Akamai Cloudlets V3 APIs
package v3

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// Cloudlets is the api interface for cloudlets
	Cloudlets interface {
		// ListCloudlets returns details of Cloudlets that you can create a shared policy for, including a Cloudlet name and Cloudlet type
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/get-cloudlets
		ListCloudlets(context.Context) ([]ListCloudletsItem, error)
		// ListSharedPolicies returns shared policies that are available within your group
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/get-policies
		ListSharedPolicies(context.Context, ListSharedPoliciesRequest) (*ListSharedPoliciesResponse, error)
		// CreateSharedPolicy creates a shared policy for a specific Cloudlet type
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/post-policy
		CreateSharedPolicy(context.Context, CreateSharedPolicyRequest) (*Policy, error)
		// DeleteSharedPolicy deletes an existing Cloudlets policy
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/delete-policy
		DeleteSharedPolicy(context.Context, DeleteSharedPolicyRequest) error
		// GetSharedPolicy returns information about a shared policy, including its activation status on the staging and production networks
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/get-policy
		GetSharedPolicy(context.Context, GetSharedPolicyRequest) (*Policy, error)
		// UpdateSharedPolicy updates an existing policy
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/put-policy
		UpdateSharedPolicy(context.Context, UpdateSharedPolicyRequest) (*Policy, error)
		// ClonePolicy clones the staging, production, and last modified versions of a non-shared (API v2) or shared policy into a new shared policy
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/post-policy-clone
		ClonePolicy(context.Context, ClonePolicyRequest) (*Policy, error)
		// ListActivePolicyProperties returns all active properties that are assigned to the policy
		//
		// See: https://techdocs.akamai.com/cloudlets/reference/get-policy-properties
		ListActivePolicyProperties(context.Context, ListActivePolicyPropertiesRequest) (*PolicyProperty, error)
		// ListPolicyVersions lists policy versions by policyID.
		//
		// See: https://techdocs.akamai.com/cloudlets/v3/reference/get-policy-versions
		ListPolicyVersions(context.Context, ListPolicyVersionsRequest) (*ListPolicyVersions, error)

		// GetPolicyVersion gets policy version by policyID and version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v3/reference/get-policy-version
		GetPolicyVersion(context.Context, GetPolicyVersionRequest) (*PolicyVersion, error)

		// CreatePolicyVersion creates policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v3/reference/post-policy-version
		CreatePolicyVersion(context.Context, CreatePolicyVersionRequest) (*PolicyVersion, error)

		// DeletePolicyVersion deletes policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v3/reference/delete-policy-version
		DeletePolicyVersion(context.Context, DeletePolicyVersionRequest) error

		// UpdatePolicyVersion updates policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v3/reference/put-policy-version
		UpdatePolicyVersion(context.Context, UpdatePolicyVersionRequest) (*PolicyVersion, error)
	}

	cloudlets struct {
		session.Session
	}

	// Option defines a Cloudlets option
	Option func(*cloudlets)

	// ClientFunc is a Cloudlets client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) Cloudlets
)

// Client returns a new cloudlets Client instance with the specified controller
func Client(sess session.Session, opts ...Option) Cloudlets {
	c := &cloudlets{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
