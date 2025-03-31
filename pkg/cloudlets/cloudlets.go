// Package cloudlets provides access to the Akamai Cloudlets APIs
package cloudlets

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
	// Cloudlets is the api interface for cloudlets
	Cloudlets interface {
		//LoadBalancers

		// ListOrigins lists all origins of specified type for the current account.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-origins
		ListOrigins(context.Context, ListOriginsRequest) ([]OriginResponse, error)

		// GetOrigin gets specific origin by originID.
		// This operation is only available for the APPLICATION_LOAD_BALANCER origin type.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-origin
		GetOrigin(context.Context, GetOriginRequest) (*Origin, error)

		// CreateOrigin creates configuration for an origin.
		// This operation is only available for the APPLICATION_LOAD_BALANCER origin type.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-origin
		CreateOrigin(context.Context, CreateOriginRequest) (*Origin, error)

		// UpdateOrigin creates configuration for an origin.
		// This operation is only available for the APPLICATION_LOAD_BALANCER origin type.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/put-origin
		UpdateOrigin(context.Context, UpdateOriginRequest) (*Origin, error)

		// LoadBalancerVersions

		// CreateLoadBalancerVersion creates load balancer version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-origin-versions
		CreateLoadBalancerVersion(context.Context, CreateLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// GetLoadBalancerVersion gets specific load balancer version by originID and version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-origin-version
		GetLoadBalancerVersion(context.Context, GetLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// UpdateLoadBalancerVersion updates specific load balancer version by originID and version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/put-origin-version
		UpdateLoadBalancerVersion(context.Context, UpdateLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// ListLoadBalancerVersions lists all versions of Origin with type APPLICATION_LOAD_BALANCER.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-origin-versions
		ListLoadBalancerVersions(context.Context, ListLoadBalancerVersionsRequest) ([]LoadBalancerVersion, error)

		// LoadBalancerActivations

		// ListLoadBalancerActivations fetches activations with the most recent listed first.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-origin-activations
		ListLoadBalancerActivations(context.Context, ListLoadBalancerActivationsRequest) ([]LoadBalancerActivation, error)

		// ActivateLoadBalancerVersion activates the load balancing version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-origin-activations
		ActivateLoadBalancerVersion(context.Context, ActivateLoadBalancerVersionRequest) (*LoadBalancerActivation, error)

		// Policies

		// ListPolicies lists policies.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policies
		ListPolicies(context.Context, ListPoliciesRequest) ([]Policy, error)

		// GetPolicy gets policy by policyID.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy
		GetPolicy(context.Context, GetPolicyRequest) (*Policy, error)

		// CreatePolicy creates policy.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-policy
		CreatePolicy(context.Context, CreatePolicyRequest) (*Policy, error)

		// RemovePolicy removes policy.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/delete-policy
		RemovePolicy(context.Context, RemovePolicyRequest) error

		// UpdatePolicy updates policy.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/put-policy
		UpdatePolicy(context.Context, UpdatePolicyRequest) (*Policy, error)

		// PolicyProperties

		// GetPolicyProperties gets all the associated properties by the policyID.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy-properties
		GetPolicyProperties(context.Context, GetPolicyPropertiesRequest) (map[string]PolicyProperty, error)

		// DeletePolicyProperty removes a property from a policy activation associated_properties list.
		DeletePolicyProperty(context.Context, DeletePolicyPropertyRequest) error

		// PolicyVersions

		// ListPolicyVersions lists policy versions by policyID.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy-versions
		ListPolicyVersions(context.Context, ListPolicyVersionsRequest) ([]PolicyVersion, error)

		// GetPolicyVersion gets policy version by policyID and version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy-version
		GetPolicyVersion(context.Context, GetPolicyVersionRequest) (*PolicyVersion, error)

		// CreatePolicyVersion creates policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-policy-versions
		CreatePolicyVersion(context.Context, CreatePolicyVersionRequest) (*PolicyVersion, error)

		// DeletePolicyVersion deletes policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/delete-policy-version
		DeletePolicyVersion(context.Context, DeletePolicyVersionRequest) error

		// UpdatePolicyVersion updates policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/put-policy-version
		UpdatePolicyVersion(context.Context, UpdatePolicyVersionRequest) (*PolicyVersion, error)

		// PolicyVersionActivations

		// ListPolicyActivations returns the complete activation history for the selected policy in reverse chronological order.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy-activations
		ListPolicyActivations(context.Context, ListPolicyActivationsRequest) ([]PolicyActivation, error)

		// ActivatePolicyVersion activates the selected cloudlet policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-policy-version-activations
		ActivatePolicyVersion(context.Context, ActivatePolicyVersionRequest) ([]PolicyActivation, error)

		// GetPolicyVersionRule returns information for a specific rule within a policy version.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/get-policy-version-rule
		GetPolicyVersionRule(context.Context, GetPolicyVersionRuleRequest) (MatchRule, error)

		// CreatePolicyVersionRule adds a new rule to an existing policy version. You can only add one rule at a time.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/post-policy-version-rules
		CreatePolicyVersionRule(context.Context, CreatePolicyVersionRuleRequest) (MatchRule, error)

		// UpdatePolicyVersionRule updates attributes of an existing rule within a policy version. When updating a rule, set disabled to true to prevent the rule from being evaluated against incoming requests.
		//
		// See: https://techdocs.akamai.com/cloudlets/v2/reference/put-policy-version-rule
		UpdatePolicyVersionRule(context.Context, UpdatePolicyVersionRuleRequest) (MatchRule, error)
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
