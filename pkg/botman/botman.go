// Package botman provides access to the Akamai Application Security Botman APIs
package botman

import (
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// BotMan is the botman api interface
	BotMan interface {
		AkamaiBotCategory
		AkamaiBotCategoryAction
		AkamaiDefinedBot
		BotAnalyticsCookie
		BotAnalyticsCookieValues
		BotCategoryException
		BotDetection
		BotDetectionAction
		BotEndpointCoverageReport
		BotManagementSetting
		ChallengeAction
		ChallengeInterceptionRules
		ClientSideSecurity
		ConditionalAction
		CustomBotCategory
		CustomBotCategoryAction
		CustomBotCategorySequence
		CustomClient
		CustomDefinedBot
		CustomDenyAction
		JavascriptInjection
		RecategorizedAkamaiDefinedBot
		ResponseAction
		ServeAlternateAction
		TransactionalEndpoint
		TransactionalEndpointProtection
	}

	botman struct {
		session.Session
	}

	// Option defines a BotMan option
	Option func(*botman)

	// ClientFunc is a botman client new method, this can be used for mocking
	ClientFunc func(sess session.Session, opts ...Option) BotMan
)

// Client returns a new botman Client instance with the specified controller
func Client(sess session.Session, opts ...Option) BotMan {
	p := &botman{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}
