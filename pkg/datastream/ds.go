// Package datastream provides access to the Akamai DataStream APIs
//
// See: https://techdocs.akamai.com/datastream2/reference/api
package datastream

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// DS is the ds api interface
	DS interface {
		// Activation

		// ActivateStream activates stream with given ID.
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/post-stream-activate
		ActivateStream(context.Context, ActivateStreamRequest) (*DetailedStreamVersion, error)

		// DeactivateStream deactivates stream with given ID.
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/post-stream-deactivate
		DeactivateStream(context.Context, DeactivateStreamRequest) (*DetailedStreamVersion, error)

		// GetActivationHistory returns a history of activation status changes for all versions of a stream.
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/get-stream-activation-history
		GetActivationHistory(context.Context, GetActivationHistoryRequest) ([]ActivationHistoryEntry, error)

		// Properties

		// GetProperties returns properties that are active on the production and staging network for a specific product type that are available within a group
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/get-properties-cdn
		GetProperties(context.Context, GetPropertiesRequest) (*PropertiesDetails, error)

		// GetDatasetFields returns groups of data set fields available in the template.
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/get-dataset-fields
		GetDatasetFields(context.Context, GetDatasetFieldsRequest) (*DataSets, error)

		// Application Security Configurations

		// GetAppSecConfigs returns the AppSec configurations available for streams for a given group and contract.
		//
		// See: https://techdocs.akamai.com/datastream2/reference/get-configs-appsec
		GetAppSecConfigs(context.Context, GetAppSecConfigsRequest) ([]AppSecConfigDetails, error)

		// Stream

		// CreateStream creates a stream
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/post-stream-cdn
		CreateStream(context.Context, CreateStreamRequest) (*DetailedStreamVersion, error)

		// GetStream gets stream details
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/get-stream
		GetStream(context.Context, GetStreamRequest) (*DetailedStreamVersion, error)

		// UpdateStream updates a stream
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/put-stream-cdn
		UpdateStream(context.Context, UpdateStreamRequest) (*DetailedStreamVersion, error)

		// DeleteStream deletes a stream
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/delete-stream
		DeleteStream(context.Context, DeleteStreamRequest) error

		// ListStreams retrieves list of streams
		//
		// See: https://techdocs.akamai.com/datastream2/v3/reference/get-streams
		ListStreams(context.Context, ListStreamsRequest) ([]StreamDetails, error)
	}

	ds struct {
		session.Session
	}

	// Option defines a DS option
	Option func(*ds)

	// ClientFunc is a ds client new method, this can be used for mocking
	ClientFunc func(sess session.Session, ops ...Option) DS
)

// Client returns a new ds Client instance with the specified controller
func Client(sess session.Session, opts ...Option) DS {
	c := &ds{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DelimiterTypePtr returns the address of the DelimiterType
func DelimiterTypePtr(d DelimiterType) *DelimiterType {
	return &d
}
