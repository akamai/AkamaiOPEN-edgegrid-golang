// Package cloudwrapper provides access to the Akamai Cloud Wrapper API
package cloudwrapper

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CloudWrapper is the api interface for Cloud Wrapper
	CloudWrapper interface {
		// Capacities

		// ListCapacities fetches capacities available for a given contractId.
		// If no contract id is provided, lists all available capacity locations
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-capacity-inventory
		ListCapacities(context.Context, ListCapacitiesRequest) (*ListCapacitiesResponse, error)

		// Configurations

		// GetConfiguration gets a specific Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-configuration
		GetConfiguration(context.Context, GetConfigurationRequest) (*Configuration, error)

		// ListConfigurations lists all Cloud Wrapper configurations on your contract
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-configurations
		ListConfigurations(context.Context) (*ListConfigurationsResponse, error)

		// CreateConfiguration creates a Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration
		CreateConfiguration(context.Context, CreateConfigurationRequest) (*Configuration, error)

		// UpdateConfiguration updates a saved or inactive configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/put-configuration
		UpdateConfiguration(context.Context, UpdateConfigurationRequest) (*Configuration, error)

		// DeleteConfiguration deletes configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/delete-configuration
		DeleteConfiguration(context.Context, DeleteConfigurationRequest) error

		// ActivateConfiguration activates a Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration-activations
		ActivateConfiguration(context.Context, ActivateConfigurationRequest) error

		// Locations

		// ListLocations returns a list of locations available to distribute Cloud Wrapper capacity
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-locations
		ListLocations(context.Context) (*ListLocationResponse, error)

		// MultiCDN

		// ListAuthKeys lists the cdnAuthKeys for a specified contractId and cdnCode
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-auth-keys
		ListAuthKeys(context.Context, ListAuthKeysRequest) (*ListAuthKeysResponse, error)

		// ListCDNProviders lists CDN providers
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-providers
		ListCDNProviders(context.Context) (*ListCDNProvidersResponse, error)

		// Properties

		// ListProperties lists unused properties
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-properties
		ListProperties(context.Context, ListPropertiesRequest) (*ListPropertiesResponse, error)

		// ListOrigins lists property origins
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-origins
		ListOrigins(context.Context, ListOriginsRequest) (*ListOriginsResponse, error)
	}

	cloudwrapper struct {
		session.Session
	}

	// Option defines an CloudWrapper option
	Option func(*cloudwrapper)
)

// Client returns a new cloudwrapper Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CloudWrapper {
	c := &cloudwrapper{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
