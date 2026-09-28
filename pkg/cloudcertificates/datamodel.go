package cloudcertificates

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// errEmptyAlgorithmKey indicates an Algorithms map contains an empty-string key.
var errEmptyAlgorithmKey = errors.New("must not contain an empty key")

// validateAlgorithmKeysNotEmpty checks that no key in an Algorithms map is an empty string. Key values are
// intentionally not validated against a fixed enum (e.g. RSA/ECDSA) here, since the set of supported
// cryptographic algorithms may be extended by the API without requiring an SDK update.
func validateAlgorithmKeysNotEmpty(value any) error {
	algorithms, ok := value.(map[CryptographicAlgorithm]SignedCertificate)
	if !ok {
		return nil
	}
	for k := range algorithms {
		if k == "" {
			return errEmptyAlgorithmKey
		}
	}
	return nil
}

var (
	_ validation.Validatable = CreateLineageRequest{}
	_ validation.Validatable = CreateLineageRequestBody{}
	_ validation.Validatable = GetLineageRequest{}
	_ validation.Validatable = ListLineagesRequest{}
	_ validation.Validatable = RenameLineageRequest{}
	_ validation.Validatable = DeleteLineageRequest{}
	_ validation.Validatable = RenewLineageRequest{}
	_ validation.Validatable = CompleteLineageRequest{}
	_ validation.Validatable = CompleteLineageRequestBody{}
	_ validation.Validatable = ListLineageActivityRequest{}
	_ validation.Validatable = GetGenerationRequest{}
	_ validation.Validatable = DeleteGenerationRequest{}
	_ validation.Validatable = ListArchivedGenerationsRequest{}
	_ validation.Validatable = UploadSignedCertificateRequest{}
	_ validation.Validatable = UploadSignedCertificateRequestBody{}
	_ validation.Validatable = PromoteLineageRequest{}
	_ validation.Validatable = RollbackLineageRequest{}
	_ validation.Validatable = ReplaceStagingLineageRequest{}
	_ validation.Validatable = GetActivationStatusRequest{}
	_ validation.Validatable = ListActivationsRequest{}
	_ validation.Validatable = ListLineageBindingsRequest{}
	_ validation.Validatable = Subject{}
	_ validation.Validatable = KeySpec{}
	_ validation.Validatable = SignedCertificate{}
	_ validation.Validatable = CryptographicAlgorithm("")
	_ validation.Validatable = KeySize("")
	_ validation.Validatable = SecureNetwork("")
	_ validation.Validatable = LineageType("")
	_ validation.Validatable = GeoClass("")
	_ validation.Validatable = ExpandGenerations("")
	_ validation.Validatable = StackMode("")
	_ validation.Validatable = GenerationStatus("")
	_ validation.Validatable = TargetNetwork("")
	_ validation.Validatable = SortOrder("")
)

var (
	_ GenerationPointer = (*HeadGeneration)(nil)
	_ GenerationPointer = (*ProductionGeneration)(nil)
	_ GenerationPointer = (*StagingGeneration)(nil)
	_ GenerationPointer = (*PreviousProductionGeneration)(nil)
)

const (
	// CertificateStatusReadyForUse indicates the accepted signed certificate for a key type is ready for usage.
	CertificateStatusReadyForUse = "READY_FOR_USE"

	// CertificateStatusCSRReady indicates the CSR generation is complete and available for download.
	CertificateStatusCSRReady = "CSR_READY"

	// CertificateStatusCertUploadProcessing indicates the uploaded certificate is being validated. This is a transient status.
	CertificateStatusCertUploadProcessing = "CERT_UPLOAD_PROCESSING"

	// CertificateStatusAbandoned indicates the algorithm instance was cascaded to abandoned along with its generation.
	CertificateStatusAbandoned = "ABANDONED"

	// CryptographicAlgorithmRSA indicates the `RSA` algorithm.
	CryptographicAlgorithmRSA CryptographicAlgorithm = "RSA"

	// CryptographicAlgorithmECDSA indicates the `ECDSA` algorithm.
	CryptographicAlgorithmECDSA CryptographicAlgorithm = "ECDSA"

	// SecureNetworkEnhancedTLS represents the `ENHANCED_TLS` secure network type.
	SecureNetworkEnhancedTLS SecureNetwork = "ENHANCED_TLS"

	// SecureNetworkStandardTLS represents the `STANDARD_TLS` secure network type.
	SecureNetworkStandardTLS SecureNetwork = "STANDARD_TLS"

	// KeySize2048 represents an RSA key size of 2048 bits.
	KeySize2048 KeySize = "2048"

	// KeySize4096 represents an RSA key size of 4096 bits.
	KeySize4096 KeySize = "4096"

	// KeySizeP256 represents an ECDSA key size of 256 bits.
	KeySizeP256 KeySize = "P-256"

	// KeySizeP384 represents an ECDSA key size of 384 bits.
	KeySizeP384 KeySize = "P-384"

	// GeoClassStandardWorldwide represents the `STANDARD_WORLDWIDE` geographic class.
	GeoClassStandardWorldwide GeoClass = "STANDARD_WORLDWIDE"

	// TargetNetworkStaging represents the `STAGING` deployment network.
	TargetNetworkStaging TargetNetwork = "STAGING"

	// TargetNetworkProduction represents the `PRODUCTION` deployment network.
	TargetNetworkProduction TargetNetwork = "PRODUCTION"

	// ExpandGenerationsHead expands only the head generation of a lineage.
	ExpandGenerationsHead ExpandGenerations = "HEAD"

	// ExpandGenerationsCurrentProduction expands the generation currently deployed to the production network.
	ExpandGenerationsCurrentProduction ExpandGenerations = "CURRENT_PRODUCTION"

	// ExpandGenerationsCurrentStaging expands the generation currently deployed to the staging network.
	ExpandGenerationsCurrentStaging ExpandGenerations = "CURRENT_STAGING"

	// ExpandGenerationsPreviousProduction expands the generation previously deployed to the production network, the rollback candidate.
	ExpandGenerationsPreviousProduction ExpandGenerations = "PREVIOUS_PRODUCTION"

	// LineageTypeMultipleGeneration is the default lineage type, supporting multiple generations over the lineage's lifetime.
	LineageTypeMultipleGeneration LineageType = "MULTIPLE_GENERATION"

	// LineageTypeSingleGeneration restricts the lineage to a single, standalone generation.
	LineageTypeSingleGeneration LineageType = "SINGLE_GENERATION"

	// GenerationStatusCSRReady indicates all algorithm instances have a CSR, awaiting signed certificates.
	GenerationStatusCSRReady GenerationStatus = "CSR_READY"

	// GenerationStatusReadyForUse indicates at least one algorithm instance has been uploaded and validated.
	GenerationStatusReadyForUse GenerationStatus = "READY_FOR_USE"

	// GenerationStatusActive indicates the generation is deployed on `currentProduction` or `currentStaging`.
	GenerationStatusActive GenerationStatus = "ACTIVE"

	// GenerationStatusArchived indicates the generation is terminal, read-only history and has fallen out of the live window.
	GenerationStatusArchived GenerationStatus = "ARCHIVED"

	// GenerationStatusAbandoned indicates the generation was explicitly discarded before ever reaching production.
	GenerationStatusAbandoned GenerationStatus = "ABANDONED"

	// StackModeSingleStack indicates the lineage declares exactly one key algorithm.
	StackModeSingleStack StackMode = "SINGLE_STACK"

	// StackModeMultipleStack indicates the lineage declares both RSA and ECDSA key algorithms.
	StackModeMultipleStack StackMode = "MULTIPLE_STACK"

	// ActivityEventTypeLineageCreated indicates a lineage was created.
	ActivityEventTypeLineageCreated = "LINEAGE_CREATED"

	// ActivityEventTypeCertUploaded indicates a signed certificate was uploaded for a generation.
	ActivityEventTypeCertUploaded = "CERT_UPLOADED"

	// ActivityEventTypeGenerationPromotedToStaging indicates a generation was promoted to the staging network.
	ActivityEventTypeGenerationPromotedToStaging = "GENERATION_PROMOTED_TO_STAGING"

	// ActivityEventTypeGenerationPromotedToProduction indicates a generation was promoted to the production network.
	ActivityEventTypeGenerationPromotedToProduction = "GENERATION_PROMOTED_TO_PRODUCTION"

	// ActivityEventTypeGenerationRolledBack indicates a rollback to the previous production generation.
	ActivityEventTypeGenerationRolledBack = "GENERATION_ROLLED_BACK"

	// ActivityEventTypeStagingReplaced indicates the staging generation was replaced.
	ActivityEventTypeStagingReplaced = "STAGING_REPLACED"

	// ActivityEventTypeGenerationRenewed indicates a new generation was created via renewal.
	ActivityEventTypeGenerationRenewed = "GENERATION_RENEWED"

	// ActivityEventTypeGenerationCompleted indicates a second algorithm was added to a MULTIPLE_STACK lineage.
	ActivityEventTypeGenerationCompleted = "GENERATION_COMPLETED"

	// ActivityEventTypeGenerationAbandoned indicates a generation was abandoned before reaching production.
	ActivityEventTypeGenerationAbandoned = "GENERATION_ABANDONED"

	// ActivityEventTypeGenerationDeleted indicates a generation was deleted.
	ActivityEventTypeGenerationDeleted = "GENERATION_DELETED"

	// ActivityEventTypeLineageRenamed indicates a lineage was renamed.
	ActivityEventTypeLineageRenamed = "LINEAGE_RENAMED"

	// ActivityEventTypeLineageDeleted indicates a lineage was deleted.
	ActivityEventTypeLineageDeleted = "LINEAGE_DELETED"

	// OperationTypePromote indicates a PROMOTE activation operation.
	OperationTypePromote = "PROMOTE"

	// OperationTypeRollback indicates a ROLLBACK activation operation.
	OperationTypeRollback = "ROLLBACK"

	// OperationTypeReplaceStaging indicates a REPLACE_STAGING activation operation.
	OperationTypeReplaceStaging = "REPLACE_STAGING"

	// MaxListLineageActivityPageSize is the maximum page size for listing lineage activities.
	MaxListLineageActivityPageSize = 100

	// MaxListLineagesPageSize is the maximum page size for listing lineages.
	MaxListLineagesPageSize = 100

	// MaxListActivationsPageSize is the maximum page size for listing activations.
	MaxListActivationsPageSize = 100

	// MaxListLineageBindingsPageSize is the maximum page size for listing lineage hostname bindings.
	MaxListLineageBindingsPageSize = 100

	// DefaultListLineagesPageSize is the page size the API uses for listing lineages when none is specified.
	DefaultListLineagesPageSize = 10

	// DefaultListLineageActivityPageSize is the page size the API uses for listing lineage activity when none is specified.
	DefaultListLineageActivityPageSize = 20

	// DefaultListActivationsPageSize is the page size the API uses for listing activations when none is specified.
	DefaultListActivationsPageSize = 20

	// DefaultListLineageBindingsPageSize is the page size the API uses for listing lineage hostname bindings when none is specified.
	DefaultListLineageBindingsPageSize = 50

	// SortOrderAscending sorts results in ascending order, oldest first.
	SortOrderAscending SortOrder = "ASC"

	// SortOrderDescending sorts results in descending order, newest first.
	SortOrderDescending SortOrder = "DESC"
)

type (
	// CreateLineageRequest represents the parameters for creating a new certificate lineage.
	CreateLineageRequest struct {
		// Body of the create lineage request containing the lineage details.
		Body CreateLineageRequestBody
	}

	// CreateLineageRequestBody contains the details for creating a certificate lineage.
	CreateLineageRequestBody struct {
		// ContractID under which this lineage will be created.
		ContractID string `json:"contractId"`

		// GroupID is a unique identifier for the group.
		GroupID int64 `json:"groupId"`

		// GeoClass is the geographic class of the certificate. If empty, the API assigns a default value.
		GeoClass GeoClass `json:"geoClass,omitempty"`

		// LineageName is the name of the lineage, up to 270 characters. If empty, the API generates one.
		LineageName string `json:"lineageName,omitempty"`

		// LineageType is the type of the lineage. If empty, the API assigns [LineageTypeMultipleGeneration] as the default.
		LineageType LineageType `json:"lineageType,omitempty"`

		// SecureNetwork is the secure network type to use for the certificate. Valid values are [SecureNetworkEnhancedTLS] and [SecureNetworkStandardTLS].
		SecureNetwork SecureNetwork `json:"secureNetwork"`

		// KeySpecs is the list of key specifications for which CSRs are generated.
		KeySpecs []KeySpec `json:"keySpecs"`

		// SANs is the list of Subject Alternative Names (SANs) for the certificate.
		SANs []string `json:"sans"`

		// Subject fields as defined in X.509 certificates (RFC 5280). If nil, no subject fields are sent.
		Subject *Subject `json:"subject,omitempty"`
	}

	// CreateLineageResponse contains the response for creating a certificate lineage.
	CreateLineageResponse Lineage

	// GetLineageRequest represents the parameters for retrieving a certificate lineage.
	GetLineageRequest struct {
		// LineageID is the unique identifier of the lineage on which to perform the desired operation.
		LineageID int64

		// ExpandGenerations optionally expands generations of the lineage in the response, e.g. [ExpandGenerationsHead],
		// [ExpandGenerationsCurrentProduction], [ExpandGenerationsCurrentStaging], [ExpandGenerationsPreviousProduction].
		// Multiple values are combined into a single comma-separated query parameter.
		ExpandGenerations []ExpandGenerations
	}

	// GetLineageResponse contains the response for retrieving a certificate lineage.
	GetLineageResponse Lineage

	// ListLineagesRequest represents the parameters for listing certificate lineages.
	ListLineagesRequest struct {
		// ContractID filters lineages by Akamai contract ID.
		ContractID string

		// LineageName performs a case-insensitive substring match on the lineage name.
		LineageName string

		// LineageIDs filters lineages by one or more lineage IDs. Multiple IDs are sent as a comma-separated query parameter.
		LineageIDs []int64

		// SecureNetwork filters lineages by secure network type, e.g. [SecureNetworkEnhancedTLS] or [SecureNetworkStandardTLS].
		SecureNetwork SecureNetwork

		// StackMode filters lineages by stack mode, e.g. [StackModeSingleStack] or [StackModeMultipleStack].
		StackMode StackMode

		// LineageType filters lineages by type, e.g. [LineageTypeMultipleGeneration] or [LineageTypeSingleGeneration].
		LineageType LineageType

		// Domain performs a case-insensitive exact or wildcard match against the lineage's SANs and subject common
		// name, e.g. `domain=example.com` also matches `*.example.com`.
		Domain string

		// GenerationStatus filters lineages where at least one active pointer generation (head, currentProduction,
		// previousProduction, or currentStaging; [GenerationStatusArchived] generations excluded) has one of the given statuses.
		GenerationStatus []GenerationStatus

		// ExpiringInDays filters lineages where at least one [CertificateStatusReadyForUse] signed certificate across active pointer
		// algorithm instances expires within N days. A value of 0 or less returns lineages with already-expired
		// certificates. Nil omits this filter.
		ExpiringInDays *int

		// KeyType filters lineages that declare the given algorithm in their key specs, e.g. [CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA].
		KeyType CryptographicAlgorithm

		// Issuer performs a case-insensitive substring match on the signed certificate issuer across [CertificateStatusReadyForUse]
		// active pointer algorithm instances.
		Issuer string

		// ExpandGenerations optionally returns full generation objects instead of identifier stubs for the given
		// pointer roles, e.g. [ExpandGenerationsHead], [ExpandGenerationsCurrentProduction], [ExpandGenerationsCurrentStaging],
		// [ExpandGenerationsPreviousProduction].
		ExpandGenerations []ExpandGenerations

		// PageSize is the maximum number of lineages to return per page. Defaults to [DefaultListLineagesPageSize] if 0. Maximum is [MaxListLineagesPageSize].
		PageSize int

		// After is an opaque keyset pagination cursor obtained from a previous response's [ListLineagesResponse.NextCursor] field.
		After string

		// Sort orders results by one field prefixed with `+` (ascending) or `-` (descending), e.g. `-modifiedDate`.
		// Supported fields: `modifiedDate`, `createdDate`, `lineageName`, `expirationDate`. Defaults to `-modifiedDate`
		// if empty.
		Sort string
	}

	// ListLineagesResponse contains a page of certificate lineages.
	ListLineagesResponse struct {
		// Lineages is the page of certificate lineages matching the request filters.
		Lineages []Lineage `json:"lineages"`

		// NextCursor is an opaque keyset pagination cursor to pass as [ListLineagesRequest.After] to retrieve the next page,
		// or nil on the last page.
		NextCursor *string `json:"nextCursor"`

		// TotalCount is the total number of lineages matching the request filters, across all pages.
		TotalCount int64 `json:"totalCount"`
	}

	// RenameLineageRequest represents the parameters for renaming a certificate lineage.
	RenameLineageRequest struct {
		// LineageID is the unique identifier of the lineage to rename.
		LineageID int64

		// LineageName is the new name for the lineage.
		LineageName string
	}

	// RenameLineageResponse contains the response for renaming a certificate lineage.
	RenameLineageResponse Lineage

	// lineagePatchOperation represents a single JSON Patch (RFC 6902) operation. It is used internally to build
	// the wire request body for RenameLineage; the only currently supported operation is replacing
	// `/lineageName`.
	lineagePatchOperation struct {
		// Op is the patch operation. The only supported value is `replace`.
		Op string `json:"op"`

		// Path is the JSON pointer to the field being patched. The only supported value is `/lineageName`.
		Path string `json:"path"`

		// Value is the new lineage name.
		Value string `json:"value"`
	}

	// DeleteLineageRequest represents the parameters for deleting a certificate lineage. This permanently
	// deletes the lineage and all its generations, activations, and activity history. The lineage must have
	// no active production or staging generation and no pending activation in progress, or the API rejects
	// the request.
	DeleteLineageRequest struct {
		// LineageID is the unique identifier of the lineage to delete.
		LineageID int64
	}

	// RenewLineageRequest represents the parameters for renewing a certificate lineage: creating a new head
	// generation with fresh CSRs for all key types in the lineage's key specs.
	RenewLineageRequest struct {
		// LineageID is the unique identifier of the lineage to renew.
		LineageID int64

		// ConfirmAbandonHead confirms abandoning the existing head generation when it has a still-valid,
		// unexpired certificate. Not required when there is no head, the head has no valid certificate, or its
		// CSRs have expired.
		ConfirmAbandonHead bool
	}

	// RenewLineageResponse contains the response for renewing a certificate lineage.
	RenewLineageResponse Lineage

	// CompleteLineageRequest represents the parameters for completing a certificate lineage: adding the second
	// algorithm (e.g. [CryptographicAlgorithmECDSA]) to a [StackModeMultipleStack] lineage where only one algorithm (e.g. [CryptographicAlgorithmRSA])
	// is currently live on production.
	CompleteLineageRequest struct {
		// LineageID is the unique identifier of the lineage to complete.
		LineageID int64

		// AcknowledgeWarnings acknowledges warnings and retries the complete operation when the response contains warnings.
		AcknowledgeWarnings bool

		// Body of the complete lineage request containing the signed certificate for the missing algorithm.
		Body CompleteLineageRequestBody
	}

	// CompleteLineageRequestBody contains the signed certificate for the algorithm being completed, keyed by
	// cryptographic algorithm. Exactly one entry is expected, matching the [GenerationStatusCSRReady] algorithm on the lineage's
	// current production generation.
	CompleteLineageRequestBody struct {
		// SourceGenerationID optionally identifies the generation the completed algorithm is added to. If given, it
		// must match the lineage's current production generation.
		SourceGenerationID *int64 `json:"sourceGenerationId,omitempty"`

		// Algorithms maps a cryptographic algorithm ([CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA]) to its signed certificate.
		Algorithms map[CryptographicAlgorithm]SignedCertificate `json:"algorithms"`

		// AutoActivate optionally triggers an immediate PROMOTE to the given networks, in order, once the complete
		// operation succeeds (e.g. [TargetNetworkStaging] then [TargetNetworkProduction]). The complete operation is
		// always committed first; if activation fails, the complete result is preserved and the activation error is
		// returned alongside the response.
		AutoActivate []TargetNetwork `json:"autoActivate,omitempty"`
	}

	// CompleteLineageResponse contains the response for completing a certificate lineage.
	CompleteLineageResponse Lineage

	// ListLineageActivityRequest represents the parameters for listing the activity events of a certificate lineage.
	ListLineageActivityRequest struct {
		// LineageID is the unique identifier of the lineage to list activity for.
		LineageID int64

		// PageSize is the maximum number of events to return per page. Defaults to [DefaultListLineageActivityPageSize] if 0. Maximum is
		// [MaxListLineageActivityPageSize].
		PageSize int

		// Cursor is an opaque keyset pagination cursor obtained from a previous response's [ListLineageActivityResponse.NextCursor] field.
		Cursor string
	}

	// ListLineageActivityResponse contains a page of certificate lineage activity events, newest first.
	ListLineageActivityResponse struct {
		// Events is the page of activity events matching the request, ordered newest first.
		Events []LineageActivityEvent `json:"events"`

		// NextCursor is an opaque keyset pagination cursor to pass as [ListLineageActivityRequest.Cursor] to retrieve the next page,
		// or nil on the last page.
		NextCursor *string `json:"nextCursor"`

		// TotalCount is the total number of activity events for the lineage, across all pages.
		TotalCount int64 `json:"totalCount"`
	}

	// LineageActivityEvent represents a single activity event recorded for a certificate lineage.
	LineageActivityEvent struct {
		// ActivityID is the unique identifier of the activity event.
		ActivityID int64 `json:"activityId"`

		// EventType is the type of the activity event, e.g. [ActivityEventTypeLineageCreated], [ActivityEventTypeCertUploaded],
		// [ActivityEventTypeLineageRenamed].
		EventType string `json:"eventType"`

		// LineageID is the unique identifier of the lineage the event belongs to.
		LineageID int64 `json:"lineageId"`

		// GenerationID is the unique identifier of the generation the event relates to, or nil for lineage-level events.
		GenerationID *int64 `json:"generationId"`

		// CreatedBy is the user who triggered the event.
		CreatedBy string `json:"createdBy"`

		// Network is the target network the event relates to, e.g. [TargetNetworkStaging] or [TargetNetworkProduction],
		// or nil if not applicable.
		Network *string `json:"network"`

		// Outcome is the outcome of the event, e.g. `ALL_SUCCESS`, `PARTIAL_SUCCESS`, or `ALL_FAILED`, or nil if
		// not applicable.
		Outcome *string `json:"outcome"`

		// CreatedTime is the time the event was recorded.
		CreatedTime time.Time `json:"createdTime"`
	}

	// GetGenerationRequest represents the parameters for retrieving a single certificate generation within a
	// lineage.
	GetGenerationRequest struct {
		// LineageID is the unique identifier of the lineage the generation belongs to.
		LineageID int64

		// GenerationID is the unique identifier of the generation to retrieve.
		GenerationID int64
	}

	// GetGenerationResponse contains the full detail of a single certificate generation within a lineage,
	// including complete algorithm instance metadata.
	//
	// This endpoint returns unprefixed generationId/generationStatus fields, unlike the embedded generation
	// pointers on [Lineage]/[GetLineageResponse] (HeadGeneration, ProductionGeneration, StagingGeneration,
	// PreviousProductionGeneration), which use role-specific prefixes such as headGenerationId.
	//
	GetGenerationResponse struct {
		Generation

		// GenerationID is the unique identifier of the generation.
		GenerationID int64 `json:"generationId"`

		// GenerationStatus is the status of the generation.
		GenerationStatus string `json:"generationStatus"`
	}

	// DeleteGenerationRequest represents the parameters for deleting a single certificate generation within a
	// lineage. The generation must not be currently deployed (i.e. must not be the lineage's currentProduction
	// or currentStaging generation, nor otherwise ACTIVE), and must not be the only generation in the lineage.
	DeleteGenerationRequest struct {
		// LineageID is the unique identifier of the lineage the generation belongs to.
		LineageID int64

		// GenerationID is the unique identifier of the generation to delete.
		GenerationID int64

		// AcknowledgeRollbackCandidateRemoval confirms deleting the generation even though it is the lineage's
		// previousProduction (rollback candidate). Not required when the generation is not the rollback candidate.
		AcknowledgeRollbackCandidateRemoval bool
	}

	// ListArchivedGenerationsRequest represents the parameters for listing the archived and abandoned generations
	// of a certificate lineage (i.e. every generation except the current [ExpandGenerationsHead], [ExpandGenerationsCurrentProduction]
	// [ExpandGenerationsCurrentStaging], and [ExpandGenerationsPreviousProduction] pointers).
	ListArchivedGenerationsRequest struct {
		// LineageID is the unique identifier of the lineage to list archived generations for.
		LineageID int64

		// IncludeAlgorithms controls whether full algorithm instance detail (CSR/signed certificate PEM, issuer,
		// validity dates, etc.) is returned for each generation's Algorithms. Defaults to false, in which case
		// each [Algorithm] only has identifier-only fields populated ([Algorithm.AlgorithmInstanceID], [Algorithm.CertificateStatus],
		// [Algorithm.KeyType], [Algorithm.SignedCertificateNotValidAfterDate], and [Algorithm.SignedCertificateSerialNumber]).
		IncludeAlgorithms bool
	}

	// ListArchivedGenerationsResponse contains the unpaginated list of archived and abandoned generations for a
	// certificate lineage.
	ListArchivedGenerationsResponse struct {
		// Items is the list of archived and abandoned generations for the lineage.
		Items []ArchivedGeneration `json:"items"`
	}

	// ArchivedGeneration represents a single archived or abandoned certificate generation within a lineage.
	ArchivedGeneration struct {
		Generation

		// GenerationID is the unique identifier of the generation.
		GenerationID int64 `json:"generationId"`

		// GenerationStatus is the status of the generation.
		GenerationStatus string `json:"generationStatus"`
	}

	// UploadSignedCertificateRequest represents the parameters for uploading a signed certificate to a lineage generation.
	UploadSignedCertificateRequest struct {
		// LineageID is the unique identifier of the lineage.
		LineageID int64

		// GenerationID is the unique identifier of the generation within the lineage.
		GenerationID int64

		// AcknowledgeWarnings acknowledges warnings and retries the certificate upload when the response contains warnings.
		AcknowledgeWarnings bool

		// Body of the upload signed certificate request.
		Body UploadSignedCertificateRequestBody
	}

	// UploadSignedCertificateRequestBody contains the signed certificates keyed by cryptographic algorithm.
	UploadSignedCertificateRequestBody struct {
		// Algorithms maps a cryptographic algorithm ([CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA]) to its signed certificate.
		Algorithms map[CryptographicAlgorithm]SignedCertificate `json:"algorithms"`

		// AutoActivate optionally triggers an immediate PROMOTE to the given networks, in order, once the upload
		// succeeds (e.g. [TargetNetworkStaging] then [TargetNetworkProduction]). The upload is always committed
		// first; if activation fails, the upload result is preserved and the activation error is returned
		// alongside the response.
		AutoActivate []TargetNetwork `json:"autoActivate,omitempty"`
	}

	// SignedCertificate contains the signed certificate material for a single key type.
	SignedCertificate struct {
		// SignedCertificatePEM is the PEM-encoded signed certificate.
		SignedCertificatePEM string `json:"signedCertificatePem"`

		// TrustChainPEM is the optional PEM-encoded trust chain for the signed certificate.
		TrustChainPEM string `json:"trustChainPem,omitempty"`
	}

	// UploadSignedCertificateResponse contains the response for uploading a signed certificate.
	UploadSignedCertificateResponse struct {
		// Algorithms is the list of per key type certificate details after the upload.
		Algorithms []Algorithm `json:"algorithms"`

		// FirstPromotedToProductionAt is the time the generation was first promoted to the production network, or
		// nil if it has never been promoted to production. Only becomes non-nil within this response when
		// [UploadSignedCertificateRequestBody.AutoActivate] included [TargetNetworkProduction] and the promotion
		// succeeded as part of this request.
		FirstPromotedToProductionAt *time.Time `json:"firstPromotedToProductionAt"`

		// GenerationID is the unique identifier of the generation.
		GenerationID int64 `json:"generationId"`

		// GenerationModifiedBy is the user who last modified the generation.
		GenerationModifiedBy string `json:"generationModifiedBy"`

		// GenerationModifiedTime is the time the generation was last modified.
		GenerationModifiedTime time.Time `json:"generationModifiedTime"`

		// GenerationStatus is the status of the generation.
		GenerationStatus string `json:"generationStatus"`

		// LineageID is the unique identifier of the lineage.
		LineageID int64 `json:"lineageId"`

		// ValidationResults contains any errors, warnings, or notices about the upload. Nil if there are none.
		ValidationResults *UploadValidationResults `json:"validationResults"`
	}

	// UploadValidationResults contains validation errors, warnings, and notices returned by [UploadSignedCertificate].
	// Unlike [ValidationResults] (used by lineage and activation responses), this endpoint's response can also
	// carry Errors and Notices in addition to Warnings.
	UploadValidationResults struct {
		// Errors is a list of validation errors.
		Errors []ValidationResultItem `json:"errors"`

		// Notices is a list of validation notices.
		Notices []ValidationResultItem `json:"notices"`

		// Warnings is a list of validation warnings.
		Warnings []ValidationResultItem `json:"warnings"`
	}

	// PromoteLineageRequest represents the parameters for promoting a lineage's head generation to one or more networks.
	PromoteLineageRequest struct {
		// LineageID is the unique identifier of the lineage to promote.
		LineageID int64

		// GenerationID is the unique identifier of the generation to promote. It must be the lineage's head
		// generation; the API rejects the request otherwise.
		GenerationID int64

		// Target networks for head generation promotion. Allowed values are [TargetNetworkStaging]
		// and [TargetNetworkProduction]. If left empty, defaults to both networks for the initial
		// promotion and [TargetNetworkProduction] for subsequent promotions.
		Networks []TargetNetwork
	}

	// activationWireRequestBody is the actual JSON body sent to the API for PromoteLineage, RollbackLineage, and
	// ReplaceStagingLineage: the `POST /{id}/activations` endpoint shares this single request body shape across
	// all three operations, discriminated by `OperationType`. `Networks` is only meaningful for PROMOTE; the
	// Promote/Rollback/ReplaceStaging requests deliberately don't expose the `operationType` discriminator (or,
	// for Rollback/ReplaceStaging, `networks`) to callers, since each method only ever performs its own operation.
	activationWireRequestBody struct {
		OperationType string          `json:"operationType"`
		GenerationID  int64           `json:"generationId"`
		Networks      []TargetNetwork `json:"networks,omitempty"`
	}

	// lineageActivationResponse contains the common response structure for lineage activation operations
	// (promote, rollback, replace staging).
	lineageActivationResponse struct {
		// Items is the list of activation request responses.
		Items []GetActivationStatusResponse `json:"items"`

		// ValidationResults contains any warnings about the activation operation. Nil if there are none.
		ValidationResults *ValidationResults `json:"validationResults"`
	}

	// PromoteLineageResponse contains the response for promoting a lineage's head generation.
	//
	// Reuses [GetActivationStatusResponse] for each item.
	//
	// PromoteLineage returns a nil *PromoteLineageResponse (and nil error) when the head generation is already
	// active on every requested network: the API responds 204 No Content in that case, rather than 200 with an
	// activation page.
	PromoteLineageResponse lineageActivationResponse

	// RollbackLineageRequest represents the parameters for rolling back a lineage's production generation to its
	// previous production generation.
	RollbackLineageRequest struct {
		// LineageID is the unique identifier of the lineage to roll back.
		LineageID int64

		// GenerationID is the unique identifier of the generation to roll back to. It must be the lineage's
		// previous production generation; the API rejects the request otherwise.
		GenerationID int64
	}

	// RollbackLineageResponse contains the response for rolling back a lineage's production generation.
	//
	// Reuses [GetActivationStatusResponse] for each item.
	//
	// RollbackLineage returns a nil *RollbackLineageResponse (and nil error) when the previous production
	// generation is already the current production generation: the API responds 204 No Content in that case,
	// rather than 200 with an activation page.
	RollbackLineageResponse lineageActivationResponse

	// ReplaceStagingLineageRequest represents the parameters for replacing the generation currently deployed on
	// staging with the lineage's current production or previous production generation.
	ReplaceStagingLineageRequest struct {
		// LineageID is the unique identifier of the lineage.
		LineageID int64

		// GenerationID is the unique identifier of the generation to use as the source for the replace-staging
		// operation. It must be the lineage's current production or previous production generation; the API
		// rejects the request otherwise.
		GenerationID int64
	}

	// ReplaceStagingLineageResponse contains the response for replacing a lineage's current staging generation.
	//
	// Reuses [GetActivationStatusResponse] for each item.
	//
	// ReplaceStagingLineage returns a nil *ReplaceStagingLineageResponse (and nil error) when [ReplaceStagingLineageRequest.GenerationID]
	// is already the current staging generation: the API responds 204 No Content in that case, rather than 200
	// with an activation page.
	ReplaceStagingLineageResponse lineageActivationResponse

	// GetActivationStatusRequest represents the parameters for getting the status of a lineage activation
	// request (as returned by PromoteLineage, RollbackLineage, or ReplaceStagingLineage).
	GetActivationStatusRequest struct {
		// LineageID is the unique identifier of the lineage.
		LineageID int64

		// ActivationID is the unique identifier of the activation request.
		ActivationID int64
	}

	// GetActivationStatusResponse contains the status of a lineage activation request.
	GetActivationStatusResponse struct {
		// ActivationID is the unique identifier of the activation request.
		ActivationID int64 `json:"activationId"`

		// LineageID is the unique identifier of the lineage the activation belongs to.
		LineageID int64 `json:"lineageId"`

		// ActivationType is the type of the activation operation: [OperationTypePromote], [OperationTypeRollback], or
		// [OperationTypeReplaceStaging].
		ActivationType string `json:"activationType"`

		// ActivationStatus is the status of the activation request: `PENDING`, `IN_PROGRESS`, `COMPLETE`,
		// `PARTIAL_SUCCESS`, `FAILED`, `ABORTED`, `PRE_EMPTED`, or `DELAYED`.
		ActivationStatus string `json:"activationStatus"`

		// GenerationID is the unique identifier of the generation being activated.
		GenerationID int64 `json:"generationId"`

		// ActivationCreatedTime is the time the activation request was created.
		ActivationCreatedTime time.Time `json:"activationCreatedTime"`

		// ActivationModifiedTime is the time the activation request was last updated.
		ActivationModifiedTime time.Time `json:"activationModifiedTime"`

		// TargetEnvironment is the network the generation is being activated to: [TargetNetworkStaging] or [TargetNetworkProduction].
		TargetEnvironment string `json:"targetEnvironment"`

		// TotalHostnameCount is the total number of hostnames being deployed as part of this activation, or nil
		// if not yet known.
		TotalHostnameCount *int64 `json:"totalHostnameCount"`

		// InProgressHostnameCount is the number of hostnames still in progress for this activation, or nil if
		// not yet known.
		InProgressHostnameCount *int64 `json:"inProgressHostnameCount"`

		// PreEmptedBy optionally identifies the activation request that pre-empted (superseded) this one, or
		// nil if this activation was not pre-empted.
		PreEmptedBy *int64 `json:"preEmptedBy"`

		// CreatedBy is the user who created the activation request.
		CreatedBy string `json:"createdBy"`

		// ModifiedBy is the user who last modified the activation request.
		ModifiedBy string `json:"modifiedBy"`

		// ErrorTypes optionally contains error type information when the activation failed, empty otherwise.
		ErrorTypes string `json:"errorTypes"`
	}

	// ListActivationsRequest represents the parameters for listing the activation requests ([OperationTypePromote],
	// [OperationTypeRollback], or [OperationTypeReplaceStaging]) recorded for a certificate lineage.
	ListActivationsRequest struct {
		// LineageID is the unique identifier of the lineage to list activations for.
		LineageID int64

		// PageSize is the maximum number of activations to return per page. Defaults to [DefaultListActivationsPageSize] if 0. Maximum is [MaxListActivationsPageSize].
		PageSize int

		// Cursor is an opaque keyset pagination cursor obtained from a previous response's [ListActivationsResponse.NextCursor] field.
		Cursor string
	}

	// ListActivationsResponse contains a page of certificate lineage activation requests, newest first.
	ListActivationsResponse struct {
		// Items is the page of activation requests matching the request, ordered newest first.
		//
		// Reuses GetActivationStatusResponse for each item, since the server returns the exact same underlying
		// model for both endpoints.
		Items []GetActivationStatusResponse `json:"items"`

		// NextCursor is an opaque keyset pagination cursor to pass as [ListActivationsRequest.Cursor] to retrieve the next page,
		// or nil on the last page.
		NextCursor *string `json:"nextCursor"`

		// TotalCount is the total number of activation requests for the lineage, across all pages.
		TotalCount int64 `json:"totalCount"`

		// ValidationResults contains any warnings aggregated across this page's activations. Nil if there are none.
		ValidationResults *ValidationResults `json:"validationResults"`
	}

	// ListLineageBindingsRequest represents the parameters for listing the hostname bindings of a certificate lineage.
	ListLineageBindingsRequest struct {
		// LineageID is the unique identifier of the lineage to list hostname bindings for.
		LineageID int64

		// Network filters bindings by network, e.g., [TargetNetworkStaging] or [TargetNetworkProduction].
		// If empty, bindings on both networks are returned.
		Network TargetNetwork

		// PageSize is the maximum number of bindings to return per page. Defaults to [DefaultListLineageBindingsPageSize] when set to 0.
		// Maximum is [MaxListLineageBindingsPageSize].
		PageSize int

		// After is an opaque keyset pagination cursor obtained from a previous response's [ListLineageBindingsResponse.NextCursor] field.
		// It remains stable as long as Network and Sort are unchanged between requests.
		After string

		// Sort orders results by an internal chronological binding identifier.
		// Valid values are [SortOrderAscending] (oldest first) and [SortOrderDescending] (newest first). Defaults to
		// [SortOrderAscending] if empty. Must not change between pages when using After.
		Sort SortOrder
	}

	// ListLineageBindingsResponse contains a page of hostname bindings for a certificate lineage.
	ListLineageBindingsResponse struct {
		// Bindings is the page of hostname bindings matching the request.
		Bindings []LineageBinding `json:"bindings"`

		// NextCursor is an opaque keyset pagination cursor to pass as After to retrieve the next page, or nil on
		// the last page.
		NextCursor *string `json:"nextCursor"`

		// TotalCount is the total number of bindings matching the request filters, across all pages.
		TotalCount int64 `json:"totalCount"`
	}

	// LineageBinding represents a single hostname bound to a certificate lineage.
	LineageBinding struct {
		// Active indicates whether the binding is currently active.
		Active bool `json:"active"`

		// Hostname is the bound hostname.
		Hostname string `json:"hostname"`

		// Networks is the list of networks the hostname is bound on, [TargetNetworkStaging], [TargetNetworkProduction], or both.
		Networks []string `json:"network"`
	}

	// Lineage represents a certificate lineage and its metadata.
	Lineage struct {
		// AccountID is the account identifier associated with the ContractID.
		AccountID string `json:"accountId"`

		// ContractID is the contract identifier.
		ContractID string `json:"contractId"`

		// GeoClass is the geographic class of the certificate.
		GeoClass string `json:"geoClass"`

		// GroupID is the unique identifier for the group.
		GroupID int64 `json:"groupId"`

		// Head is the head generation of the lineage, or nil if the lineage currently has no head (e.g. once
		// fully promoted to production with no in-progress renewal or completion).
		Head *HeadGeneration `json:"head"`

		// CurrentProduction is the generation currently deployed to the production network, or nil if none.
		// Populated only when requested via [ExpandGenerationsCurrentProduction].
		CurrentProduction *ProductionGeneration `json:"currentProduction"`

		// PreviousProduction is the generation previously deployed to production, the rollback candidate, or nil if none.
		// Populated only when requested via [ExpandGenerationsPreviousProduction].
		PreviousProduction *PreviousProductionGeneration `json:"previousProduction"`

		// CurrentStaging is the generation currently deployed to the staging network, or nil if none.
		// Populated only when requested via [ExpandGenerationsCurrentStaging].
		CurrentStaging *StagingGeneration `json:"currentStaging"`

		// KeySpecs is the list of key specifications for the lineage.
		KeySpecs []KeySpecResponse `json:"keySpecs"`

		// LineageCreatedBy is the user who created the lineage.
		LineageCreatedBy string `json:"lineageCreatedBy"`

		// LineageCreatedTime is the time the lineage was created.
		LineageCreatedTime time.Time `json:"lineageCreatedTime"`

		// LineageID is the unique identifier of the lineage.
		LineageID int64 `json:"lineageId"`

		// LineageModifiedBy is the user who last modified the lineage.
		LineageModifiedBy string `json:"lineageModifiedBy"`

		// LineageModifiedTime is the time the lineage was last modified.
		LineageModifiedTime time.Time `json:"lineageModifiedTime"`

		// LineageName is the name of the lineage.
		LineageName string `json:"lineageName"`

		// LineageType is the type of the lineage.
		LineageType string `json:"lineageType"`

		// SANs is the list of Subject Alternative Names (SANs) for the certificate.
		SANs []string `json:"sans"`

		// SecureNetwork is the secure network type of the certificate.
		SecureNetwork string `json:"secureNetwork"`

		// StackMode is the stack mode of the lineage.
		StackMode string `json:"stackMode"`

		// Subject contains the subject details of the certificate.
		Subject Subject `json:"subject"`

		// ValidationResults contains any warnings related to the lineage.
		ValidationResults *ValidationResults `json:"validationResults"`
	}

	// Generation contains the fields common to every certificate generation returned by the API, regardless of
	// which pointer role (head, currentProduction, currentStaging, previousProduction) it is returned under. It
	// is embedded by [HeadGeneration], [ProductionGeneration], [StagingGeneration], [PreviousProductionGeneration],
	// [ArchivedGeneration], and [GetGenerationResponse].
	Generation struct {
		// Algorithms is the list of per key type certificate details for the generation, empty when not expanded.
		Algorithms []Algorithm `json:"algorithms"`

		// FirstPromotedToProductionTime is the time the generation was first promoted to the production network,
		// or nil if it has never been promoted to production (or the generation was not expanded).
		FirstPromotedToProductionTime *time.Time `json:"firstPromotedToProductionTime"`

		// GenerationCreatedBy is the user who created the generation, or nil if not expanded.
		GenerationCreatedBy *string `json:"generationCreatedBy"`

		// GenerationCreatedTime is the time the generation was created, or nil if not expanded.
		GenerationCreatedTime *time.Time `json:"generationCreatedTime"`

		// GenerationModifiedBy is the user who last modified the generation, or nil if not expanded.
		GenerationModifiedBy *string `json:"generationModifiedBy"`

		// GenerationModifiedTime is the time the generation was last modified, or nil if not expanded or if the
		// generation has never been modified since creation.
		GenerationModifiedTime *time.Time `json:"generationModifiedTime"`
	}

	// HeadGeneration represents the head generation of a lineage: the generation currently being prepared (CSR
	// generation, certificate upload) ahead of activation.
	HeadGeneration struct {
		Generation

		// HeadGenerationID is the unique identifier of the head generation.
		HeadGenerationID int64 `json:"headGenerationId"`

		// HeadGenerationStatus is the status of the head generation.
		HeadGenerationStatus string `json:"headGenerationStatus"`
	}

	// ProductionGeneration represents the generation currently deployed to the production network.
	ProductionGeneration struct {
		Generation

		// ProductionGenerationID is the unique identifier of the generation currently deployed to production.
		ProductionGenerationID int64 `json:"productionGenerationId"`

		// ProductionGenerationStatus is the status of the generation currently deployed to production.
		ProductionGenerationStatus string `json:"productionGenerationStatus"`
	}

	// StagingGeneration represents the generation currently deployed to the staging network.
	StagingGeneration struct {
		Generation

		// StagingGenerationID is the unique identifier of the generation currently deployed to staging.
		StagingGenerationID int64 `json:"stagingGenerationId"`

		// StagingGenerationStatus is the status of the generation currently deployed to staging.
		StagingGenerationStatus string `json:"stagingGenerationStatus"`
	}

	// PreviousProductionGeneration represents the generation previously deployed to production, the rollback
	// candidate.
	PreviousProductionGeneration struct {
		Generation

		// PreviousProductionGenerationID is the unique identifier of the generation previously deployed to production.
		PreviousProductionGenerationID int64 `json:"previousProductionGenerationId"`

		// PreviousProductionGenerationStatus is the status of the generation previously deployed to production.
		PreviousProductionGenerationStatus string `json:"previousProductionGenerationStatus"`
	}

	// GenerationPointer is implemented by HeadGeneration, ProductionGeneration, StagingGeneration, and
	// PreviousProductionGeneration. Each role uses its own prefixed field names on the wire (e.g.
	// headGenerationId vs productionGenerationId); this interface lets callers read a pointer's ID, status, and
	// common generation detail the same way regardless of role.
	GenerationPointer interface {
		// ID returns the generation's unique identifier.
		ID() int64

		// Status returns the generation's status.
		Status() string

		// Common returns the fields common to every generation, embedded via Generation.
		Common() Generation
	}

	// Algorithm represents the certificate material and status for a single key type within a generation.
	Algorithm struct {
		// AlgorithmInstanceCreatedBy is the user who created the algorithm instance.
		AlgorithmInstanceCreatedBy string `json:"algorithmInstanceCreatedBy"`

		// AlgorithmInstanceCreatedTime is the time the algorithm instance was created.
		AlgorithmInstanceCreatedTime *time.Time `json:"algorithmInstanceCreatedTime"`

		// AlgorithmInstanceID is the unique identifier of the algorithm instance.
		AlgorithmInstanceID int64 `json:"algorithmInstanceId"`

		// AlgorithmInstanceModifiedBy is the user who last modified the algorithm instance.
		AlgorithmInstanceModifiedBy *string `json:"algorithmInstanceModifiedBy"`

		// AlgorithmInstanceModifiedTime is the time the algorithm instance was last modified.
		AlgorithmInstanceModifiedTime *time.Time `json:"algorithmInstanceModifiedTime"`

		// CertificateStatus is the status of the certificate for this key type.
		CertificateStatus string `json:"certificateStatus"`

		// CSRExpirationDate is the date when the CSR expires.
		CSRExpirationDate *time.Time `json:"csrExpirationDate"`

		// CSRPEM is the PEM-encoded certificate signing request generated for this key type.
		CSRPEM string `json:"csrPem"`

		// KeyType is the key type of the algorithm, either [CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA].
		KeyType string `json:"keyType"`

		// SignedCertificateIssuer is the issuer field of the signed certificate.
		SignedCertificateIssuer *string `json:"signedCertificateIssuer"`

		// SignedCertificateNotValidAfterDate is the expiration date of the signed certificate.
		SignedCertificateNotValidAfterDate *time.Time `json:"signedCertificateNotValidAfterDate"`

		// SignedCertificateNotValidBeforeDate is the date before which the signed certificate is not valid.
		SignedCertificateNotValidBeforeDate *time.Time `json:"signedCertificateNotValidBeforeDate"`

		// SignedCertificatePEM is the PEM-encoded signed certificate uploaded for this key type.
		SignedCertificatePEM *string `json:"signedCertificatePem"`

		// SignedCertificateSerialNumber is the signed certificate serial number in hex format.
		SignedCertificateSerialNumber *string `json:"signedCertificateSerialNumber"`

		// SignedCertificateSHA256Fingerprint is the SHA-256 fingerprint of the signed certificate.
		SignedCertificateSHA256Fingerprint *string `json:"signedCertificateSha256Fingerprint"`

		// TrustChainPEM is the PEM-encoded trust chain uploaded alongside the signed certificate, or nil if none
		// was uploaded.
		TrustChainPEM *string `json:"trustChainPem"`
	}

	// KeySpec represents a key specification for a certificate.
	KeySpec struct {
		// KeyType is the key type. Valid values are [CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA].
		KeyType CryptographicAlgorithm `json:"keyType"`

		// KeySize is the key size. Valid values for [CryptographicAlgorithmRSA]: [KeySize2048], [KeySize4096].
		// Valid values for [CryptographicAlgorithmECDSA]: [KeySizeP256], [KeySizeP384].
		KeySize KeySize `json:"keySize"`
	}

	// KeySpecResponse represents a key specification for a certificate as returned by the API.
	KeySpecResponse struct {
		// KeyType is the key type, e.g. [CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA].
		KeyType string `json:"keyType"`

		// KeySize is the key size, e.g. [KeySize2048], [KeySize4096], [KeySizeP256], or [KeySizeP384].
		KeySize string `json:"keySize"`
	}

	// Subject contains the subject details for a certificate.
	Subject struct {
		// CommonName is the fully qualified domain name (FQDN) or other name associated with the subject.
		CommonName string `json:"commonName,omitempty"`

		// Organization is the legal name of the organization.
		Organization string `json:"organization,omitempty"`

		// OrganizationalUnit is the organizational unit associated with the subject.
		OrganizationalUnit string `json:"organizationalUnit,omitempty"`

		// Country is the two-letter ISO 3166 country code.
		Country string `json:"country,omitempty"`

		// State is the full name of the state or province.
		State string `json:"state,omitempty"`

		// Locality is the city or locality name.
		Locality string `json:"locality,omitempty"`
	}

	// ValidationResults contains validation warnings returned by the API.
	ValidationResults struct {
		// Warnings is a list of validation warnings.
		Warnings []ValidationResultItem `json:"warnings"`
	}

	// ValidationResultItem represents a validation result item (error, warning, or notice).
	ValidationResultItem struct {
		// Context provides additional context about the item.
		Context map[string]any `json:"context"`

		// Detail provides a human-readable description of the item.
		Detail string `json:"detail"`

		// Instance is a URI reference that identifies the specific occurrence of the item.
		Instance string `json:"instance"`

		// Status is the HTTP status code associated with the item.
		Status int `json:"status"`

		// Title is a human-readable summary of the item.
		Title string `json:"title"`

		// Type is a URI reference that identifies the item type.
		Type string `json:"type"`
	}

	// CryptographicAlgorithm represents the cryptographic algorithm type: [CryptographicAlgorithmRSA] or [CryptographicAlgorithmECDSA].
	CryptographicAlgorithm string

	// SecureNetwork represents the type of secure network, e.g. [SecureNetworkEnhancedTLS].
	SecureNetwork string

	// KeySize represents the size of the key: [KeySize2048] or [KeySize4096] for [CryptographicAlgorithmRSA],
	// [KeySizeP256] or [KeySizeP384] for [CryptographicAlgorithmECDSA].
	KeySize string

	// GeoClass represents the geographic class of a certificate. [GeoClassStandardWorldwide] is the only value
	// confirmed by API documentation and live responses to date.
	GeoClass string

	// ExpandGenerations represents which generations of a lineage to expand in a response: [ExpandGenerationsHead], [ExpandGenerationsCurrentProduction],
	// [ExpandGenerationsCurrentStaging], or [ExpandGenerationsPreviousProduction].
	ExpandGenerations string

	// TargetNetwork represents a deployment network target for lineage activation: [TargetNetworkStaging] or [TargetNetworkProduction].
	TargetNetwork string

	// LineageType represents the type of certificate lineage: [LineageTypeMultipleGeneration] or [LineageTypeSingleGeneration].
	LineageType string

	// StackMode represents the stack mode of a lineage used to filter lineages: [StackModeSingleStack] or [StackModeMultipleStack].
	StackMode string

	// GenerationStatus represents the status of a certificate generation used to filter lineages: [GenerationStatusCSRReady],
	// [GenerationStatusReadyForUse], [GenerationStatusActive], [GenerationStatusArchived], or [GenerationStatusAbandoned].
	GenerationStatus string

	// SortOrder represents the sort direction for a keyset-paginated list: [SortOrderAscending] or [SortOrderDescending].
	SortOrder string
)

// ID returns the head generation's unique identifier.
func (h *HeadGeneration) ID() int64 { return h.HeadGenerationID }

// Status returns the head generation's status.
func (h *HeadGeneration) Status() string { return h.HeadGenerationStatus }

// Common returns the head generation's common fields.
func (h *HeadGeneration) Common() Generation { return h.Generation }

// ID returns the production generation's unique identifier.
func (p *ProductionGeneration) ID() int64 { return p.ProductionGenerationID }

// Status returns the production generation's status.
func (p *ProductionGeneration) Status() string { return p.ProductionGenerationStatus }

// Common returns the production generation's common fields.
func (p *ProductionGeneration) Common() Generation { return p.Generation }

// ID returns the staging generation's unique identifier.
func (s *StagingGeneration) ID() int64 { return s.StagingGenerationID }

// Status returns the staging generation's status.
func (s *StagingGeneration) Status() string { return s.StagingGenerationStatus }

// Common returns the staging generation's common fields.
func (s *StagingGeneration) Common() Generation { return s.Generation }

// ID returns the previous production generation's unique identifier.
func (p *PreviousProductionGeneration) ID() int64 { return p.PreviousProductionGenerationID }

// Status returns the previous production generation's status.
func (p *PreviousProductionGeneration) Status() string { return p.PreviousProductionGenerationStatus }

// Common returns the previous production generation's common fields.
func (p *PreviousProductionGeneration) Common() Generation { return p.Generation }

// Validate validates CreateLineageRequest.
func (r CreateLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Body": validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates CreateLineageRequestBody.
func (r CreateLineageRequestBody) Validate() error {
	return validation.Errors{
		"ContractID":    validation.Validate(r.ContractID, validation.Required),
		"GroupID":       validation.Validate(r.GroupID, validation.Required),
		"GeoClass":      validation.Validate(r.GeoClass),
		"LineageName":   validation.Validate(r.LineageName, validation.Length(0, 270)),
		"LineageType":   validation.Validate(r.LineageType),
		"SecureNetwork": validation.Validate(r.SecureNetwork, validation.Required),
		"KeySpecs":      validation.Validate(r.KeySpecs, validation.Required),
		"SANs":          validation.Validate(r.SANs, validation.Required),
		"Subject":       validation.Validate(r.Subject),
	}.Filter()
}

// Validate validates GetLineageRequest.
func (r GetLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":         validation.Validate(r.LineageID, validation.Required),
		"ExpandGenerations": validation.Validate(r.ExpandGenerations),
	})
}

// Validate validates ListLineagesRequest.
func (r ListLineagesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageIDs": validation.Validate(r.LineageIDs, validation.When(r.LineageIDs != nil,
			validation.Required,
			validation.Each(validation.Min(1)),
		)),
		"SecureNetwork":     validation.Validate(r.SecureNetwork),
		"StackMode":         validation.Validate(r.StackMode),
		"LineageType":       validation.Validate(r.LineageType),
		"GenerationStatus":  validation.Validate(r.GenerationStatus),
		"KeyType":           validation.Validate(r.KeyType),
		"ExpandGenerations": validation.Validate(r.ExpandGenerations),
		"PageSize":          validation.Validate(r.PageSize, validation.Min(0), validation.Max(MaxListLineagesPageSize)),
	})
}

// Validate validates GetGenerationRequest.
func (r GetGenerationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
	})
}

// Validate validates DeleteGenerationRequest.
func (r DeleteGenerationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
	})
}

// Validate validates RenameLineageRequest.
func (r RenameLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":   validation.Validate(r.LineageID, validation.Required),
		"LineageName": validation.Validate(r.LineageName, validation.Required, validation.Length(1, 270)),
	})
}

// Validate validates DeleteLineageRequest.
func (r DeleteLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
	})
}

// Validate validates RenewLineageRequest.
func (r RenewLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
	})
}

// Validate validates ListArchivedGenerationsRequest.
func (r ListArchivedGenerationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
	})
}

// Validate validates CompleteLineageRequest.
func (r CompleteLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
		"Body":      validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates CompleteLineageRequestBody.
func (r CompleteLineageRequestBody) Validate() error {
	return validation.Errors{
		"Algorithms":   validation.Validate(r.Algorithms, validation.Required, validation.Length(1, 1), validation.By(validateAlgorithmKeysNotEmpty)),
		"AutoActivate": validation.Validate(r.AutoActivate),
	}.Filter()
}

// Validate validates ListLineageActivityRequest.
func (r ListLineageActivityRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
		"PageSize":  validation.Validate(r.PageSize, validation.Min(0), validation.Max(MaxListLineageActivityPageSize)),
	})
}

// Validate validates ListActivationsRequest.
func (r ListActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required),
		"PageSize":  validation.Validate(r.PageSize, validation.Min(0), validation.Max(MaxListActivationsPageSize)),
	})
}

// Validate validates ListLineageBindingsRequest.
func (r ListLineageBindingsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID": validation.Validate(r.LineageID, validation.Required, validation.Min(1)),
		"Network":   validation.Validate(r.Network),
		"PageSize":  validation.Validate(r.PageSize, validation.Min(0), validation.Max(MaxListLineageBindingsPageSize)),
		"Sort":      validation.Validate(r.Sort),
	})
}

// Validate validates UploadSignedCertificateRequest.
func (r UploadSignedCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
		"Body":         validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates UploadSignedCertificateRequestBody.
func (r UploadSignedCertificateRequestBody) Validate() error {
	return validation.Errors{
		"Algorithms":   validation.Validate(r.Algorithms, validation.Required, validation.By(validateAlgorithmKeysNotEmpty)),
		"AutoActivate": validation.Validate(r.AutoActivate),
	}.Filter()
}

// Validate validates PromoteLineageRequest.
func (r PromoteLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
		"Networks":     validation.Validate(r.Networks),
	})
}

// Validate validates RollbackLineageRequest.
func (r RollbackLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
	})
}

// Validate validates ReplaceStagingLineageRequest.
func (r ReplaceStagingLineageRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"GenerationID": validation.Validate(r.GenerationID, validation.Required),
	})
}

// Validate validates GetActivationStatusRequest.
func (r GetActivationStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"LineageID":    validation.Validate(r.LineageID, validation.Required),
		"ActivationID": validation.Validate(r.ActivationID, validation.Required),
	})
}

// Validate validates KeySpec.
func (k KeySpec) Validate() error {
	return validation.Errors{
		"KeyType": validation.Validate(k.KeyType, validation.Required),
		"KeySize": validation.Validate(k.KeySize, validation.Required),
	}.Filter()
}

// Validate validates Subject.
func (s Subject) Validate() error {
	return validation.Errors{
		"CommonName":   validation.Validate(s.CommonName, validation.Length(1, 64), validation.Match(regexp.MustCompile(`\S`))),
		"Organization": validation.Validate(s.Organization, validation.Length(1, 64), validation.Match(regexp.MustCompile(`\S`))),
		"Country":      validation.Validate(s.Country, validation.Length(2, 2), validation.Match(regexp.MustCompile(`\S`))),
		"State":        validation.Validate(s.State, validation.Length(1, 128), validation.Match(regexp.MustCompile(`\S`))),
		"Locality":     validation.Validate(s.Locality, validation.Length(1, 128), validation.Match(regexp.MustCompile(`\S`))),
	}.Filter()
}

// Validate validates SignedCertificate.
func (s SignedCertificate) Validate() error {
	return validation.Errors{
		"SignedCertificatePEM": validation.Validate(s.SignedCertificatePEM, validation.Required),
	}.Filter()
}

// Validate validates CryptographicAlgorithm.
func (c CryptographicAlgorithm) Validate() error {
	return validation.In(CryptographicAlgorithmRSA, CryptographicAlgorithmECDSA).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'", c, CryptographicAlgorithmRSA, CryptographicAlgorithmECDSA)).
		Validate(c)
}

// Validate validates KeySize.
func (k KeySize) Validate() error {
	return validation.In(KeySize2048, KeySize4096, KeySizeP256, KeySizeP384).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', or '%s'",
			k, KeySize2048, KeySize4096, KeySizeP256, KeySizeP384)).
		Validate(k)
}

// Validate validates SecureNetwork.
func (s SecureNetwork) Validate() error {
	return validation.In(SecureNetworkEnhancedTLS, SecureNetworkStandardTLS).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'",
			s, SecureNetworkEnhancedTLS, SecureNetworkStandardTLS)).
		Validate(s)
}

// Validate validates GeoClass.
func (g GeoClass) Validate() error {
	return validation.In(GeoClassStandardWorldwide).
		Error(fmt.Sprintf("value '%s' is invalid. Must be '%s'", g, GeoClassStandardWorldwide)).
		Validate(g)
}

// Validate validates ExpandGenerations.
func (e ExpandGenerations) Validate() error {
	return validation.In(ExpandGenerationsHead, ExpandGenerationsCurrentProduction, ExpandGenerationsCurrentStaging, ExpandGenerationsPreviousProduction).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', or '%s'",
			e, ExpandGenerationsHead, ExpandGenerationsCurrentProduction, ExpandGenerationsCurrentStaging, ExpandGenerationsPreviousProduction)).
		Validate(e)
}

// Validate validates LineageType.
func (l LineageType) Validate() error {
	return validation.In(LineageTypeMultipleGeneration, LineageTypeSingleGeneration).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'",
			l, LineageTypeMultipleGeneration, LineageTypeSingleGeneration)).
		Validate(l)
}

// Validate validates StackMode.
func (s StackMode) Validate() error {
	return validation.In(StackModeSingleStack, StackModeMultipleStack).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'",
			s, StackModeSingleStack, StackModeMultipleStack)).
		Validate(s)
}

// Validate validates GenerationStatus.
func (g GenerationStatus) Validate() error {
	return validation.In(GenerationStatusCSRReady, GenerationStatusReadyForUse, GenerationStatusActive, GenerationStatusArchived, GenerationStatusAbandoned).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', or '%s'",
			g, GenerationStatusCSRReady, GenerationStatusReadyForUse, GenerationStatusActive, GenerationStatusArchived, GenerationStatusAbandoned)).
		Validate(g)
}

// Validate validates TargetNetwork.
func (n TargetNetwork) Validate() error {
	return validation.In(TargetNetworkStaging, TargetNetworkProduction).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'", n, TargetNetworkStaging, TargetNetworkProduction)).
		Validate(n)
}

// Validate validates SortOrder.
func (s SortOrder) Validate() error {
	return validation.In(SortOrderAscending, SortOrderDescending).
		Error(fmt.Sprintf("value '%s' is invalid. Must be either '%s' or '%s'", s, SortOrderAscending, SortOrderDescending)).
		Validate(s)
}
