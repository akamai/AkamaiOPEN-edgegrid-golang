package hapi

import (
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ChangeRequestStatusPending selects pending change requests.
	ChangeRequestStatusPending ChangeRequestStatus = "PENDING"
	// ChangeRequestStatusSucceeded represents a successful change request.
	ChangeRequestStatusSucceeded = "SUCCEEDED"
	// ChangeRequestStatusFailed represents a failed change request.
	ChangeRequestStatusFailed = "FAILED"
	// ChangeRequestStatusIgnored represents an ignored change request.
	ChangeRequestStatusIgnored = "IGNORED"

	// ChangeRequestOptimizationStatusProcessing represents optimization in progress.
	ChangeRequestOptimizationStatusProcessing = "PROCESSING"
	// ChangeRequestOptimizationStatusProcessed represents completed optimization.
	ChangeRequestOptimizationStatusProcessed = "PROCESSED"
	// ChangeRequestOptimizationStatusRejected represents rejected optimization.
	ChangeRequestOptimizationStatusRejected = "REJECTED"
	// ChangeRequestOptimizationStatusRetryFailed represents a failed optimization retry.
	ChangeRequestOptimizationStatusRetryFailed = "RETRY_FAILED"
	// ChangeRequestOptimizationStatusIgnored represents ignored optimization.
	ChangeRequestOptimizationStatusIgnored = "IGNORED"
)

type (
	// DeleteEdgeHostnameRequest represents the parameters for deleting an edge hostname.
	DeleteEdgeHostnameRequest struct {
		// DNSZone is the Akamai-specific portion of the edge hostname.
		DNSZone string

		// RecordName is the record name of the edge hostname.
		RecordName string

		// StatusUpdateEmail contains the email addresses to notify about the updates of the change.
		StatusUpdateEmail []string

		// Comments contains comments about the requested deletion.
		Comments string
	}

	// DeleteEdgeHostnameResponse contains details of a deletion request.
	// It is separate from ChangeRequest because the delete response schema does not include optimizationStatus.
	DeleteEdgeHostnameResponse struct {
		// Action is the action represented by the request.
		Action string `json:"action"`

		// ChangeID is the unique identifier of the deletion request.
		ChangeID int64 `json:"changeId"`

		// Comments contains the comments added when the deletion request was created.
		Comments string `json:"comments"`

		// Status is the status of the deletion request. Can be "PENDING", "SUCCEEDED", "FAILED", or "IGNORED".
		Status string `json:"status"`

		// StatusMessage is the API-generated status description of the deletion request.
		StatusMessage string `json:"statusMessage"`

		// StatusUpdateDate is the ISO 8601 formatted time of the most recent status update.
		StatusUpdateDate time.Time `json:"statusUpdateDate"`

		// StatusUpdateEmail is the comma-separated list of recipient email addresses for status updates on the deletion request.
		StatusUpdateEmail *string `json:"statusUpdateEmail"`

		// SubmitDate is the ISO 8601 formatted time the deletion request was submitted.
		SubmitDate time.Time `json:"submitDate"`

		// Submitter is the username of the person who submitted the deletion request.
		Submitter string `json:"submitter"`

		// SubmitterEmail is the email address of the deletion request's submitter.
		SubmitterEmail string `json:"submitterEmail"`

		// EdgeHostnames contains the edge hostnames affected by the deletion request.
		EdgeHostnames []EdgeHostname `json:"edgeHostnames"`
	}

	// UpdateEdgeHostnameRequest represents the parameters for updating an edge hostname.
	UpdateEdgeHostnameRequest struct {
		// DNSZone is the Akamai-specific portion of the edge hostname.
		DNSZone string

		// RecordName is the record name of the edge hostname.
		RecordName string

		// StatusUpdateEmail contains the email addresses to notify about the updates of the change.
		StatusUpdateEmail []string

		// Comments contains additional comments about the requested update.
		Comments string

		// Body contains the JSON Patch operations to apply to the edge hostname.
		Body []UpdateEdgeHostnameRequestBody
	}

	// UpdateEdgeHostnameRequestBody specifies a JSON Patch operation for an edge hostname.
	UpdateEdgeHostnameRequestBody struct {
		// Op is the JSON Patch operator to apply.
		Op string `json:"op"`

		// Path is the JSON Pointer expression, starting with a slash character, that locates the item within the object to change.
		Path string `json:"path"`

		// Value provides the replacement text.
		Value string `json:"value"`
	}

	// UpdateEdgeHostnameResponse contains details of an update request.
	UpdateEdgeHostnameResponse ChangeRequest

	// EdgeHostname contains the configuration of an edge hostname affected by a change request.
	EdgeHostname struct {
		// EdgeHostnameID is the unique identifier of the edge hostname.
		EdgeHostnameID int64 `json:"edgeHostnameId"`

		// RecordName is the origin domain portion of the edge hostname.
		RecordName string `json:"recordName"`

		// DNSZone is the Akamai-specific portion of the edge hostname.
		DNSZone string `json:"dnsZone"`

		// SecurityType indicates whether "STANDARD-TLS" or "ENHANCED-TLS" applies.
		SecurityType string `json:"securityType"`

		// UseDefaultTTL indicates whether the product's default TTL is used.
		UseDefaultTTL bool `json:"useDefaultTtl"`

		// UseDefaultMap indicates whether the product's default map is used.
		UseDefaultMap bool `json:"useDefaultMap"`

		// TTL is the number of seconds to keep the edge hostname assigned to a map or target.
		TTL int64 `json:"ttl"`

		// Map correlates the edge hostname with the IP address it resolves to.
		Map string `json:"map"`

		// SlotNumber is the edge hostname's slot number.
		SlotNumber *int64 `json:"slotNumber"`

		// IPVersionBehavior specifies which versions of the IP protocol to use. Either "IPV4" or "IPV6_IPV4_DUALSTACK".
		IPVersionBehavior string `json:"ipVersionBehavior"`

		// Comments contains the comments provided when the edge hostname was created.
		Comments string `json:"comments"`

		// ChinaCDN contains the edge hostname's China CDN settings.
		ChinaCDN ChinaCDN `json:"chinaCdn"`

		// CustomTarget is the user-specified target for the edge hostname.
		CustomTarget *string `json:"customTarget"`

		// IsEdgeIPBindingEnabled indicates whether fixed-IP delivery is enabled.
		IsEdgeIPBindingEnabled *bool `json:"isEdgeIPBindingEnabled"`

		// MapAlias is an alternate name for a map.
		MapAlias *string `json:"mapAlias"`

		// ProductID identifies the product associated with the edge hostname.
		ProductID string `json:"productId"`

		// SerialNumber is the edge hostname's certificate serial number.
		SerialNumber int64 `json:"serialNumber"`

		// UseCases contains the preset scenarios used to optimize traffic.
		UseCases []UseCase `json:"useCases"`

		// HTTPSServiceBinding specifies the protocols advertised by the HTTPS service binding DNS record.
		HTTPSServiceBinding *string `json:"httpsServiceBinding"`

		// DeliveryPolicy controls the geographical regions in which edge traffic is kept.
		DeliveryPolicy *DeliveryPolicy `json:"deliveryPolicy"`

		// DNSSerial is the serial number for the edge hostname.
		DNSSerial *int64 `json:"dnsSerial"`
	}

	// ChinaCDN contains settings for serving edge hostname traffic in mainland China.
	ChinaCDN struct {
		// IsChinaCDN indicates whether the edge hostname can serve traffic in China.
		IsChinaCDN bool `json:"isChinaCdn"`

		// CustomChinaCDNMap is the hostname configured with the China partner.
		CustomChinaCDNMap *string `json:"customChinaCdnMap"`

		// ChinaPartner is the name of the China partner.
		ChinaPartner *string `json:"chinaPartner"`

		// GTMDomain is the Global Traffic Management domain that serves traffic in China.
		GTMDomain *string `json:"gtmDomain"`
	}

	// DeliveryPolicy specifies the delivery policy assigned to an edge hostname.
	DeliveryPolicy struct {
		// PolicyType is the type of delivery policy.
		PolicyType string `json:"policyType"`

		// PolicyAttributes contains the delivery policy details.
		PolicyAttributes []DeliveryPolicyAttribute `json:"policyAttributes"`
	}

	// DeliveryPolicyAttribute specifies a delivery policy attribute.
	DeliveryPolicyAttribute struct {
		// AttributeType is the type of policy attribute.
		AttributeType *string `json:"attributeType"`

		// Code is the code of the geographical region in which edge traffic is kept.
		Code string `json:"code"`

		// Name is the name of the policy attribute.
		Name *string `json:"name"`
	}

	// UseCase specifies a preset scenario for optimizing edge hostname traffic.
	UseCase struct {
		// Type identifies the Akamai network over which traffic deploys.
		Type string `json:"type"`

		// Option specifies an available option for the use case's traffic type.
		Option string `json:"option"`

		// UseCase uniquely identifies the preset scenario.
		UseCase string `json:"useCase"`
	}

	// GetEdgeHostnameRequest represents the parameters for retrieving an edge hostname.
	GetEdgeHostnameRequest struct {
		// EdgeHostnameID is the unique identifier of the edge hostname.
		EdgeHostnameID int64

		// IncludeDeliveryPolicy includes the delivery policy when set. False by default.
		IncludeDeliveryPolicy *bool
	}

	// GetEdgeHostnameResponse contains an edge hostname's configuration.
	GetEdgeHostnameResponse EdgeHostname

	// GetCertificateRequest identifies the edge hostname whose certificate is to be retrieved.
	GetCertificateRequest struct {
		// DNSZone is the Akamai-specific portion of the edge hostname.
		DNSZone string

		// RecordName is the origin domain portion of the edge hostname.
		RecordName string
	}

	// GetCertificateResponse contains the certificate associated with an edge hostname.
	GetCertificateResponse struct {
		// AvailableDomains contains the certificate's domains, including the common name.
		AvailableDomains []string `json:"availableDomains"`

		// CertificateID is the unique identifier of the certificate.
		CertificateID string `json:"certificateId"`

		// CertificateType specifies the type of certificate. Can be "SAN", "SINGLE", "THIRD_PARTY", "WILDCARD", or "WILDCARD_SAN".
		CertificateType string `json:"certificateType"`

		// CommonName is the hostname for which the certificate is valid.
		CommonName string `json:"commonName"`

		// ExpirationDate is the ISO 8601 formatted time at which the certificate expires.
		ExpirationDate time.Time `json:"expirationDate"`

		// SerialNumber is the certificate's serial number.
		SerialNumber string `json:"serialNumber"`

		// SlotNumber is the edge hostname's slot number.
		SlotNumber int64 `json:"slotNumber"`

		// Status is the certificate's deployment or expiration status.
		Status string `json:"status"`

		// ValidationType specifies the level or source of certificate validation. Can be "DOMAIN_VALIDATION", "ORGANIZATIONAL_VALIDATION", "THIRD_PARTY", or "EXTENDED_VALIDATION".
		ValidationType string `json:"validationType"`
	}

	// GetChangeRequestRequest identifies a change request and the optional details to retrieve.
	GetChangeRequestRequest struct {
		// ChangeID is the unique identifier of the change request.
		ChangeID int64

		// IncludeDeliveryPolicy includes the delivery policy when set. False by default.
		IncludeDeliveryPolicy *bool
	}

	// GetChangeRequestResponse contains the status and details of an edge hostname change request.
	GetChangeRequestResponse ChangeRequest

	// ChangeRequest contains the status and details of an edge hostname change request returned by update, get, and list operations.
	ChangeRequest struct {
		// Action is the action represented by the change request.
		Action string `json:"action"`

		// ChangeID is the unique identifier of the change request.
		ChangeID int64 `json:"changeId"`

		// Comments contains the comments added when the change request was created.
		Comments string `json:"comments"`

		// EdgeHostnames contains the edge hostnames affected by the change request.
		EdgeHostnames []EdgeHostname `json:"edgeHostnames"`

		// Status is the status of the change request. Can be "PENDING", "SUCCEEDED", "FAILED", or "IGNORED".
		Status string `json:"status"`

		// StatusMessage is the API-generated status description of the change request.
		StatusMessage string `json:"statusMessage"`

		// StatusUpdateEmail is the comma-separated list of recipient email addresses for status updates on the change request.
		StatusUpdateEmail *string `json:"statusUpdateEmail"`

		// StatusUpdateDate is the ISO 8601 formatted time of the most recent status update.
		StatusUpdateDate time.Time `json:"statusUpdateDate"`

		// SubmitDate is the ISO 8601 formatted time the change request was submitted.
		SubmitDate time.Time `json:"submitDate"`

		// Submitter is the username of the person who submitted the change request.
		Submitter string `json:"submitter"`

		// SubmitterEmail is the email address of the change request's submitter.
		SubmitterEmail string `json:"submitterEmail"`

		// OptimizationStatus is the change request's optimization status, either "PROCESSING", "PROCESSED", "REJECTED", "RETRY_FAILED", or "IGNORED".
		// It can be null.
		OptimizationStatus *string `json:"optimizationStatus"`
	}

	// ChangeRequestStatus is the status filter for listing change requests.
	ChangeRequestStatus string

	// ListChangeRequestsRequest represents the parameters for listing change requests.
	ListChangeRequestsRequest struct {
		// Status filters change requests. Only [ChangeRequestStatusPending] is supported.
		Status ChangeRequestStatus

		// IncludeDeliveryPolicy includes the delivery policy when set. False by default.
		IncludeDeliveryPolicy *bool
	}

	// ListChangeRequestsResponse represents the response for listing change requests.
	ListChangeRequestsResponse struct {
		// ChangeRequests contains the change requests matching the request criteria.
		ChangeRequests []ChangeRequest `json:"changeRequests"`
	}

	// ListChangeRequestsPerEdgeHostnameRequest represents the parameters for listing change requests for a specific edge hostname.
	ListChangeRequestsPerEdgeHostnameRequest struct {
		// DNSZone is the Akamai-specific portion of the edge hostname.
		DNSZone string

		// RecordName is the origin domain portion of the edge hostname.
		RecordName string

		// IncludeDeliveryPolicy includes the delivery policy when set. False by default.
		IncludeDeliveryPolicy *bool
	}

	// ListChangeRequestsPerEdgeHostnameResponse represents the response for listing change requests for a specific edge hostname.
	ListChangeRequestsPerEdgeHostnameResponse ListChangeRequestsResponse
)

// Validate validates DeleteEdgeHostnameRequest.
func (r DeleteEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	})
}

// Validate validates UpdateEdgeHostnameRequest.
func (r UpdateEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
		"Body":       validation.Validate(r.Body),
	})
}

// Validate validates UpdateEdgeHostnameRequestBody.
func (b UpdateEdgeHostnameRequestBody) Validate() error {
	return validation.Errors{
		"Path": validation.Validate(b.Path, validation.Required, validation.In("/ttl", "/ipVersionBehavior", "/httpsServiceBinding").Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '/ttl', '/ipVersionBehavior' or '/httpsServiceBinding'", b.Path))),
		"Op": validation.Validate(b.Op, validation.Required,
			validation.When(b.Path == "/ttl" || b.Path == "/ipVersionBehavior", validation.In("replace").Error(fmt.Sprintf("value '%s' is invalid. Must use 'replace' for '%s' Path", b.Op, b.Path))),
			validation.When(b.Path == "/httpsServiceBinding", validation.In("replace", "add", "remove").Error(fmt.Sprintf("value '%s' is invalid. Must use 'replace', 'add' or 'remove' for '%s' Path", b.Op, b.Path))),
		),
		"Value": validation.Validate(b.Value,
			validation.When(b.Path == "/ttl" || b.Path == "/ipVersionBehavior", validation.Required),
			validation.When(b.Path == "/httpsServiceBinding" && (b.Op == "replace" || b.Op == "add"), validation.Required, validation.In("H2", "H3", "H2_AND_H3").Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'H2', 'H3' or 'H2_AND_H3'", b.Value)))),
	}.Filter()
}

// Validate validates GetEdgeHostnameRequest.
func (r GetEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"EdgeHostnameID": validation.Validate(r.EdgeHostnameID, validation.Required),
	})
}

// Validate validates GetCertificateRequest.
func (r GetCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	})
}

// Validate validates GetChangeRequestRequest.
func (r GetChangeRequestRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ChangeID": validation.Validate(r.ChangeID, validation.Required),
	})
}

// Validate validates ListChangeRequestsRequest.
func (r ListChangeRequestsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Status": validation.Validate(r.Status, validation.In(ChangeRequestStatusPending)),
	})
}

// Validate validates ListChangeRequestsPerEdgeHostnameRequest.
func (r ListChangeRequestsPerEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	})
}
