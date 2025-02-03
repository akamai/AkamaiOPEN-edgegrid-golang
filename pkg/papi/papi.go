// Package papi provides access to the Akamai Property APIs
package papi

import (
	"context"
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/spf13/cast"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrNotFound is returned when requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrSBDNotEnabled indicates that secure-by-default is not enabled on the given account
	ErrSBDNotEnabled = errors.New("secure-by-default is not enabled")

	// ErrDefaultCertLimitReached indicates that the limit for DEFAULT certificates has been reached
	ErrDefaultCertLimitReached = errors.New("the limit for DEFAULT certificates has been reached")

	// ErrMissingComplianceRecord is returned when compliance record is required and is not provided
	ErrMissingComplianceRecord = errors.New("compliance record must be specified")
)

type (
	// PAPI is the papi api interface
	PAPI interface {
		// Activations

		// CreateActivation creates a new activation or deactivation request
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-property-activations
		CreateActivation(context.Context, CreateActivationRequest) (*CreateActivationResponse, error)

		// GetActivations returns a list of the property activations
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-activations
		GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error)

		// GetActivation gets details about an activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-activation
		GetActivation(context.Context, GetActivationRequest) (*GetActivationResponse, error)

		// CancelActivation allows for canceling an activation while it is still PENDING
		//
		// https://techdocs.akamai.com/property-mgr/reference/delete-property-activation
		CancelActivation(context.Context, CancelActivationRequest) (*CancelActivationResponse, error)

		// ClientSettings

		// GetClientSettings returns client's settings.
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-client-settings
		GetClientSettings(context.Context) (*ClientSettingsBody, error)

		// UpdateClientSettings updates client's settings.
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/put-client-settings
		UpdateClientSettings(context.Context, ClientSettingsBody) (*ClientSettingsBody, error)

		// Contracts

		// GetContracts provides a read-only list of contract names and identifiers
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-contracts
		GetContracts(context.Context) (*GetContractsResponse, error)

		// CPCodes

		// GetCPCodes lists all available CP codes
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-cpcodes
		GetCPCodes(context.Context, GetCPCodesRequest) (*GetCPCodesResponse, error)

		// GetCPCode gets the CP code with provided ID
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-cpcode
		GetCPCode(context.Context, GetCPCodeRequest) (*GetCPCodesResponse, error)

		// GetCPCodeDetail lists detailed information about a specific CP code
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/get-cpcode
		GetCPCodeDetail(context.Context, int) (*CPCodeDetailResponse, error)

		// CreateCPCode creates a new CP code
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-cpcodes
		CreateCPCode(context.Context, CreateCPCodeRequest) (*CreateCPCodeResponse, error)

		// UpdateCPCode modifies a specific CP code. You should only modify a CP code's name, time zone, and purgeable member
		//
		// See: https://techdocs.akamai.com/cp-codes/reference/put-cpcode
		UpdateCPCode(context.Context, UpdateCPCodeRequest) (*CPCodeDetailResponse, error)

		// EdgeHostnames

		// GetEdgeHostnames fetches a list of edge hostnames
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-edgehostnames
		GetEdgeHostnames(context.Context, GetEdgeHostnamesRequest) (*GetEdgeHostnamesResponse, error)

		// GetEdgeHostname fetches edge hostname with given ID
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-edgehostname
		GetEdgeHostname(context.Context, GetEdgeHostnameRequest) (*GetEdgeHostnamesResponse, error)

		// CreateEdgeHostname creates a new edge hostname
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-edgehostnames
		CreateEdgeHostname(context.Context, CreateEdgeHostnameRequest) (*CreateEdgeHostnameResponse, error)

		// Groups

		// GetGroups provides a read-only list of groups, which may contain properties.
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-groups
		GetGroups(context.Context) (*GetGroupsResponse, error)

		// Includes

		// ListIncludes lists Includes available for the current contract and group
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-includes
		ListIncludes(context.Context, ListIncludesRequest) (*ListIncludesResponse, error)

		// ListIncludeParents lists parents of a specific Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-parents
		ListIncludeParents(context.Context, ListIncludeParentsRequest) (*ListIncludeParentsResponse, error)

		// GetInclude gets information about a specific Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include
		GetInclude(context.Context, GetIncludeRequest) (*GetIncludeResponse, error)

		// CreateInclude creates a new Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-includes
		CreateInclude(context.Context, CreateIncludeRequest) (*CreateIncludeResponse, error)

		// DeleteInclude deletes an Include
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/delete-include
		DeleteInclude(context.Context, DeleteIncludeRequest) (*DeleteIncludeResponse, error)

		// IncludeRules

		// GetIncludeRuleTree gets the entire rule tree for an include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-version-rules
		GetIncludeRuleTree(context.Context, GetIncludeRuleTreeRequest) (*GetIncludeRuleTreeResponse, error)

		// UpdateIncludeRuleTree updates the rule tree for an include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/patch-include-version-rules
		UpdateIncludeRuleTree(context.Context, UpdateIncludeRuleTreeRequest) (*UpdateIncludeRuleTreeResponse, error)

		// IncludeActivations

		// ActivateInclude creates a new include activation, which deactivates any current activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-include-activation
		ActivateInclude(context.Context, ActivateIncludeRequest) (*ActivationIncludeResponse, error)

		// DeactivateInclude deactivates the include activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-include-activation
		DeactivateInclude(context.Context, DeactivateIncludeRequest) (*DeactivationIncludeResponse, error)

		// CancelIncludeActivation cancels specified include activation, if it is still in `PENDING` state
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/delete-include-activation
		CancelIncludeActivation(context.Context, CancelIncludeActivationRequest) (*CancelIncludeActivationResponse, error)

		// GetIncludeActivation gets details about an activation
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-activation
		GetIncludeActivation(context.Context, GetIncludeActivationRequest) (*GetIncludeActivationResponse, error)

		// ListIncludeActivations lists all activations for all versions of the include, on both production and staging networks
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-activations
		ListIncludeActivations(context.Context, ListIncludeActivationsRequest) (*ListIncludeActivationsResponse, error)

		// IncludeVersions

		// CreateIncludeVersion creates a new include version based on any previous version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-include-versions
		CreateIncludeVersion(context.Context, CreateIncludeVersionRequest) (*CreateIncludeVersionResponse, error)

		// GetIncludeVersion polls the state of a specific include version, for example to check its activation status
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-version
		GetIncludeVersion(context.Context, GetIncludeVersionRequest) (*GetIncludeVersionResponse, error)

		// ListIncludeVersions lists the include versions, with results limited to the 500 most recent versions
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-versions
		ListIncludeVersions(context.Context, ListIncludeVersionsRequest) (*ListIncludeVersionsResponse, error)

		// ListIncludeVersionAvailableCriteria lists available criteria for the include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-available-criteria
		ListIncludeVersionAvailableCriteria(context.Context, ListAvailableCriteriaRequest) (*AvailableCriteriaResponse, error)

		// ListIncludeVersionAvailableBehaviors lists available behaviors for the include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-available-behaviors
		ListIncludeVersionAvailableBehaviors(context.Context, ListAvailableBehaviorsRequest) (*AvailableBehaviorsResponse, error)

		// Products

		// GetProducts lists all available Products
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-products
		GetProducts(context.Context, GetProductsRequest) (*GetProductsResponse, error)

		// Properties

		// GetProperties lists properties available for the current contract and group
		//
		// https://techdocs.akamai.com/property-mgr/reference/get-properties
		GetProperties(ctx context.Context, r GetPropertiesRequest) (*GetPropertiesResponse, error)

		// CreateProperty creates a new property from scratch or bases one on another property's rule tree and optionally its set of assigned hostnames
		//
		// https://techdocs.akamai.com/property-mgr/reference/post-properties
		CreateProperty(ctx context.Context, params CreatePropertyRequest) (*CreatePropertyResponse, error)

		// GetProperty gets a specific property
		//
		// https://techdocs.akamai.com/property-mgr/reference/get-property
		GetProperty(ctx context.Context, params GetPropertyRequest) (*GetPropertyResponse, error)

		// RemoveProperty removes a specific property, which you can only do if none of its versions are currently active
		//
		// https://techdocs.akamai.com/property-mgr/reference/delete-property
		RemoveProperty(ctx context.Context, params RemovePropertyRequest) (*RemovePropertyResponse, error)

		// MapPropertyNameToID returns (PAPI) property ID for given property name
		// Mainly to be used to map (IAM) Property ID to (PAPI) Property ID
		// To get property name for the mapping, please use iam.MapPropertyIDToName
		MapPropertyNameToID(context.Context, MapPropertyNameToIDRequest) (*string, error)

		// PropertyRules

		// GetRuleTree gets the entire rule tree for a property version.
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version-rules
		GetRuleTree(context.Context, GetRuleTreeRequest) (*GetRuleTreeResponse, error)

		// UpdateRuleTree updates the rule tree for a property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/put-property-version-rules
		UpdateRuleTree(context.Context, UpdateRulesRequest) (*UpdateRulesResponse, error)

		// PropertyVersionHostnames

		// GetPropertyVersionHostnames lists all the hostnames assigned to a property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version-hostnames
		GetPropertyVersionHostnames(context.Context, GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error)

		// UpdatePropertyVersionHostnames modifies the set of hostnames for a property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/patch-property-version-hostnames
		UpdatePropertyVersionHostnames(context.Context, UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error)

		// PropertyVersions

		// GetPropertyVersions fetches available property versions
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-versions
		GetPropertyVersions(context.Context, GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error)

		// GetPropertyVersion fetches specific property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version
		GetPropertyVersion(context.Context, GetPropertyVersionRequest) (*GetPropertyVersionsResponse, error)

		// CreatePropertyVersion creates a new property version and returns location and number for the new version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-property-versions
		CreatePropertyVersion(context.Context, CreatePropertyVersionRequest) (*CreatePropertyVersionResponse, error)

		// GetLatestVersion fetches the latest property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-latest-property-version
		GetLatestVersion(context.Context, GetLatestVersionRequest) (*GetPropertyVersionsResponse, error)

		// GetAvailableBehaviors fetches a list of available behaviors for given property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-available-behaviors
		GetAvailableBehaviors(context.Context, GetAvailableBehaviorsRequest) (*GetBehaviorsResponse, error)

		// GetAvailableCriteria fetches a list of available criteria for given property version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-available-criteria
		GetAvailableCriteria(context.Context, GetAvailableCriteriaRequest) (*GetCriteriaResponse, error)

		// ListAvailableIncludes lists external resources that can be applied within a property version's rules
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version-external-resources
		ListAvailableIncludes(context.Context, ListAvailableIncludesRequest) (*ListAvailableIncludesResponse, error)

		// ListReferencedIncludes lists referenced includes for parent property
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-property-version-includes
		ListReferencedIncludes(context.Context, ListReferencedIncludesRequest) (*ListReferencedIncludesResponse, error)

		// RuleFormats

		// GetRuleFormats provides a list of rule formats
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-rule-formats
		GetRuleFormats(context.Context) (*GetRuleFormatsResponse, error)

		// Search

		// SearchProperties searches properties by name, or by the hostname or edge hostname for which it’s currently active
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/post-search-find-by-value
		SearchProperties(context.Context, SearchRequest) (*SearchResponse, error)
	}

	papi struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a PAPI option
	Option func(*papi)

	// ClientFunc is a papi client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) PAPI

	// Response is a base PAPI Response type
	Response struct {
		AccountID  string   `json:"accountId,omitempty"`
		ContractID string   `json:"contractId,omitempty"`
		GroupID    string   `json:"groupId,omitempty"`
		Etag       string   `json:"etag,omitempty"`
		Errors     []*Error `json:"errors,omitempty"`
		Warnings   []*Error `json:"warnings,omitempty"`
	}
)

// Client returns a new papi Client instance with the specified controller
func Client(sess session.Session, opts ...Option) PAPI {
	p := &papi{
		Session:     sess,
		usePrefixes: true,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithUsePrefixes sets the `PAPI-Use-Prefixes` header on requests
// See: https://techdocs.akamai.com/property-mgr/reference/id-prefixes
func WithUsePrefixes(usePrefixes bool) Option {
	return func(p *papi) {
		p.usePrefixes = usePrefixes
	}
}

// Exec overrides the session.Exec to add papi options
func (p *papi) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	// explicitly add the PAPI-Use-Prefixes header
	r.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	return p.Session.Exec(r, out, in...)
}
