// Package purgecache provides access to the Purge Cache API.
package purgecache

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

type (
	// PurgeCache is the interface for the Purge Cache API.
	PurgeCache interface {
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
