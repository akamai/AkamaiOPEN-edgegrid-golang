// Package netstorage provides access to the Akamai Net Storage APIs
package netstorage

import (
	"context"
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrNotFound is returned when requested resource was not found
	ErrNotFound = errors.New("resource not found")
)

type (
	// NetStorage is the Net Storage api interface
	NetStorage interface {
		// Groups

		// ListStorageGroups Get a list of all of the storage groups in your NetStorage instance.
		//
		// See https://techdocs.akamai.com/netstorage/reference/get-storage-groups
		ListStorageGroups(ctx context.Context, params ListStorageGroupsRequest) (*ListStorageGroupsResponse, error)

		// GetStorageGroup Get a storage group in your NetStorage instance.
		//
		// See https://techdocs.akamai.com/netstorage/reference/get-storage-group
		GetStorageGroup(ctx context.Context, params GetStorageGroupRequest) (*GetStorageGroupResponse, error)
	}

	netstorage struct {
		session.Session
	}

	// Option defines a NetStorage option
	Option func(*netstorage)

	// ClientFunc is a netstorage client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) NetStorage
)

// Client returns a new netstorage Client instance with the specified controller
func Client(sess session.Session, opts ...Option) NetStorage {
	p := &netstorage{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Exec overrides the session.Exec to add netstorage options
func (p *netstorage) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	return p.Session.Exec(r, out, in...)
}
