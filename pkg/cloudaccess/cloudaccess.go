// Package cloudaccess provides access to the Akamai Cloud Access Manager API
package cloudaccess

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

		// CreateAccessKey creates a new access key
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/post-access-key
		CreateAccessKey(context.Context, CreateAccessKeyRequest) (*CreateAccessKeyResponse, error)

		// GetAccessKey returns details for a specific access key
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key
		GetAccessKey(context.Context, AccessKeyRequest) (*GetAccessKeyResponse, error)

		// ListAccessKeys returns detailed information about all access keys available to the current user account
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-keys
		ListAccessKeys(context.Context, ListAccessKeysRequest) (*ListAccessKeysResponse, error)

		// UpdateAccessKey updates name of an access key
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/put-access-key
		UpdateAccessKey(context.Context, UpdateAccessKeyRequest, AccessKeyRequest) (*UpdateAccessKeyResponse, error)

		// DeleteAccessKey deletes an access key
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/delete-access-key
		DeleteAccessKey(context.Context, AccessKeyRequest) error

		// CreateAccessKeyVersion rotates an access key to a new version.
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/post-access-key-version
		CreateAccessKeyVersion(context.Context, CreateAccessKeyVersionRequest) (*CreateAccessKeyVersionResponse, error)

		// GetAccessKeyVersion returns detailed information for a specific version of an access key.
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version
		GetAccessKeyVersion(context.Context, GetAccessKeyVersionRequest) (*GetAccessKeyVersionResponse, error)

		// ListAccessKeyVersions returns detailed information about all the versions for a specific access key.
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-versions
		ListAccessKeyVersions(context.Context, ListAccessKeyVersionsRequest) (*ListAccessKeyVersionsResponse, error)

		// DeleteAccessKeyVersion deletes a specific version of an access key.
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/delete-access-key-version
		DeleteAccessKeyVersion(context.Context, DeleteAccessKeyVersionRequest) (*DeleteAccessKeyVersionResponse, error)

		// LookupProperties returns information about all the Property Manager properties that use a specific version of an access key.
		// This operation gets the data directly. To avoid any latency problems, use the GetAsyncPropertiesLookupID and PerformAsyncPropertiesLookup
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version-properties
		LookupProperties(ctx context.Context, params LookupPropertiesRequest) (*LookupPropertiesResponse, error)

		// GetAsyncPropertiesLookupID gets the unique identifier used to perform an PerformAsyncPropertiesLookup
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-async-version-property-lookup
		GetAsyncPropertiesLookupID(ctx context.Context, params GetAsyncPropertiesLookupIDRequest) (*GetAsyncPropertiesLookupIDResponse, error)

		// PerformAsyncPropertiesLookup returns in asynchronous way information about all the Property Manager properties that use a specific version of an access key.
		//
		// See: https://techdocs.akamai.com/cloud-access-mgr/reference/get-property-lookup
		PerformAsyncPropertiesLookup(ctx context.Context, params PerformAsyncPropertiesLookupRequest) (*PerformAsyncPropertiesLookupResponse, error)
	}

	cloudaccess struct {
		session.Session
	}

	// Option defines an CloudAccess option
	Option func(*cloudaccess)

	// KeyLink contains hypermedia link for the key
	KeyLink struct {
		AccessKeyUID int64  `json:"accessKeyUid"`
		Link         string `json:"link"`
	}

	// KeyVersion holds details for a version of an access key
	KeyVersion struct {
		AccessKeyUID int64  `json:"accessKeyUid"`
		Link         string `json:"link"`
		Version      int64  `json:"version"`
	}

	// RequestInformation contains information about a request to create an access key
	RequestInformation struct {
		AccessKeyName        string         `json:"accessKeyName"`
		AuthenticationMethod AuthType       `json:"authenticationMethod"`
		ContractID           string         `json:"contractId"`
		GroupID              int64          `json:"groupId"`
		NetworkConfiguration *SecureNetwork `json:"networkConfiguration"`
	}

	// SecureNetwork contains additional information about network
	SecureNetwork struct {
		AdditionalCDN   *CDNType    `json:"additionalCdn,omitempty"`
		SecurityNetwork NetworkType `json:"securityNetwork"`
	}

	// CDNType is a type of additionalCdn
	CDNType string

	// NetworkType is a type of securityNetwork
	NetworkType string

	// AuthType is a type of authentication
	AuthType string

	// ProcessingType is a type of ProcessingStatus
	ProcessingType string
)

const (
	// ChinaCDN represents CDN value of "CHINA_CDN"
	ChinaCDN CDNType = "CHINA_CDN"
	// RussiaCDN represents CDN value of "RUSSIA_CDN"
	RussiaCDN CDNType = "RUSSIA_CDN"

	// NetworkEnhanced represents Network value of "ENHANCED_TLS"
	NetworkEnhanced NetworkType = "ENHANCED_TLS"
	// NetworkStandard represents Network value of "STANDARD_TLS"
	NetworkStandard NetworkType = "STANDARD_TLS"

	// AuthAWS represents Authentication value of "AWS4_HMAC_SHA256"
	AuthAWS AuthType = "AWS4_HMAC_SHA256"
	// AuthGOOG represents Authentication value of "GOOG4_HMAC_SHA256"
	AuthGOOG AuthType = "GOOG4_HMAC_SHA256"

	// ProcessingInProgress represents ProcessingStatus value of 'IN_PROGRESS'
	ProcessingInProgress ProcessingType = "IN_PROGRESS"
	// ProcessingFailed represents ProcessingStatus value of 'FAILED'
	ProcessingFailed ProcessingType = "FAILED"
	// ProcessingDone represents ProcessingStatus value of 'DONE'
	ProcessingDone ProcessingType = "DONE"
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
