// Package accountprotection provides access to the Akamai Application Security Account Protection APIs
package accountprotection

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// AccountProtection is the account protection api interface
	//
	// See: https://techdocs.akamai.com/account-protector/reference/api
	AccountProtection interface {
		// ListProtectedOperations Get List of protected operations
		//
		// See: https://techdocs.akamai.com/account-protector/reference/get-account-protection
		ListProtectedOperations(ctx context.Context, params ListProtectedOperationsRequest) (*ListProtectedOperationsResponse, error)

		// GetProtectedOperationByID Get protected protected operation by operationId
		//
		// See: https://techdocs.akamai.com/account-protector/reference/get-account-protection-op
		GetProtectedOperationByID(ctx context.Context, params GetProtectedOperationByIDRequest) (*ListProtectedOperationsResponse, error)

		// CreateProtectedOperations Create a list of protected operations
		//
		// See: https://techdocs.akamai.com/account-protector/reference/post-account-protection
		CreateProtectedOperations(ctx context.Context, params CreateProtectedOperationsRequest) (*ListProtectedOperationsResponse, error)

		// UpdateProtectedOperation Update a protected operation
		//
		// See: https://techdocs.akamai.com/account-protector/reference/put-account-protection-op
		UpdateProtectedOperation(ctx context.Context, params UpdateProtectedOperationRequest) (map[string]any, error)

		// RemoveProtectedOperation Delete a protected operation
		//
		// See: https://techdocs.akamai.com/account-protector/reference/delete-account-protection-op
		RemoveProtectedOperation(ctx context.Context, params RemoveProtectedOperationRequest) error

		// GetGeneralSettings Get general settings for account protection for a given security policy
		//
		// See: https://techdocs.akamai.com/account-protector/reference/get-account-protection-settings
		GetGeneralSettings(ctx context.Context, params GetGeneralSettingsRequest) (map[string]any, error)

		// UpsertGeneralSettings Update or create general settings for account protection for a given security policy
		//
		// See: https://techdocs.akamai.com/account-protector/reference/put-account-protection-settings
		UpsertGeneralSettings(ctx context.Context, params UpsertGeneralSettingsRequest) (map[string]any, error)

		// GetUserRiskResponseStrategy Get User Risk Response Strategy for a given security configuration
		//
		// See: https://techdocs.akamai.com/account-protector/reference/get-user-risk-response-strategy
		GetUserRiskResponseStrategy(ctx context.Context, params GetUserRiskResponseStrategyRequest) (map[string]any, error)

		// UpsertUserRiskResponseStrategy Update or create User Risk Response Strategy for a given security configuration
		//
		// See: https://techdocs.akamai.com/account-protector/reference/put-user-risk-response-strategy
		UpsertUserRiskResponseStrategy(ctx context.Context, params UpsertUserRiskResponseStrategyRequest) (map[string]any, error)

		// GetUserAllowListID Get User Allow List ID for a given security configuration
		//
		// See: https://techdocs.akamai.com/account-protector/reference/get-user-allow-list
		GetUserAllowListID(ctx context.Context, params GetUserAllowListIDRequest) (map[string]any, error)

		// UpsertUserAllowListID Update User Allow List ID for a given security configuration
		//
		// See: https://techdocs.akamai.com/account-protector/reference/put-get-user-allow-list
		UpsertUserAllowListID(ctx context.Context, params UpsertUserAllowListIDRequest) (map[string]any, error)

		// DeleteUserAllowListID Delete User Allow List ID for a given security configuration
		//
		// See: https://techdocs.akamai.com/account-protector/reference/delete-get-user-allow-list
		DeleteUserAllowListID(ctx context.Context, params DeleteUserAllowListIDRequest) error
	}

	accountProtection struct {
		session.Session
	}

	// Option defines a AccountProtection option
	Option func(*accountProtection)

	// ClientFunc is a AccountProtection client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) AccountProtection
)

// Client returns a new AccountProtection Client instance with the specified controller
func Client(sess session.Session, opts ...Option) AccountProtection {
	p := &accountProtection{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
