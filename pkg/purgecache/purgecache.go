// Package purgecache provides access to the Purge Cache API.
package purgecache

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

type (
	// PurgeCache is the interface for the Purge Cache API.
	PurgeCache interface {
		// RateLimitStatus checks the rate and object limit statuses for a specific purge type.
		//
		// See: https://techdocs.akamai.com/purge-cache/reference/post-rate-limit-status
		RateLimitStatus(ctx context.Context, params RateLimitStatusRequest) (*RateLimitStatusResponse, error)
	}

	purgecache struct {
		session.Session
	}

	// Option is a function that configures the Purge Cache client.
	Option func(*purgecache)
)

// Client creates a new PurgeCache client.
func Client(sess session.Session, opts ...Option) PurgeCache {
	c := &purgecache{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
