// Package reportinggroups provides access to the Reporting Groups API.
package reportinggroups

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

type (
	// ReportingGroups is the interface for the Reporting Groups API.
	ReportingGroups interface {
	}

	reportinggroups struct {
		session.Session
	}

	// Option is a function that configures the Reporting Groups client.
	Option func(*reportinggroups)
)

// Client creates a new ReportingGroups client.
func Client(sess session.Session, opts ...Option) ReportingGroups {
	c := &reportinggroups{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
