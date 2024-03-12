// Package cloudaccess provides access to the Akamai Cloud Access Manager API
package cloudaccess

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
	// CloudAccess is the API interface for Cloud Access Manager
	CloudAccess interface {
		// GetAccessKeyStatus gets the current status and other details for a request to create a new access key
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-create-request
		GetAccessKeyStatus(context.Context, GetAccessKeyStatusRequest) (*GetAccessKeyStatusResponse, error)

		// GetAccessKeyVersionStatus gets the current status and other details for a request to create a new access key version
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version-create-request
		GetAccessKeyVersionStatus(context.Context, GetAccessKeyVersionStatusRequest) (*GetAccessKeyVersionStatusResponse, error)
	}

	cloudaccess struct {
		session.Session
	}

	// Option defines an CloudAccess option
	Option func(*cloudaccess)
)

// Client returns a new cloudaccess Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CloudAccess {
	c := &cloudaccess{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
