// Package cloudcertificates provides access to the Akamai Cloud Certificate Manager API.
package cloudcertificates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/errs"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")

	// ErrCreateLineage is returned when creating a lineage fails.
	ErrCreateLineage = errors.New("creating lineage")

	// ErrGetLineage is returned when getting a lineage fails.
	ErrGetLineage = errors.New("getting lineage")

	// ErrListLineages is returned when listing lineages fails.
	ErrListLineages = errors.New("listing lineages")

	// ErrRenameLineage is returned when renaming a lineage fails.
	ErrRenameLineage = errors.New("renaming lineage")

	// ErrDeleteLineage is returned when deleting a lineage fails.
	ErrDeleteLineage = errors.New("deleting lineage")

	// ErrRenewLineage is returned when renewing a lineage fails.
	ErrRenewLineage = errors.New("renewing lineage")

	// ErrCompleteLineage is returned when completing a lineage fails.
	ErrCompleteLineage = errors.New("completing lineage")

	// ErrListLineageActivity is returned when listing lineage activity fails.
	ErrListLineageActivity = errors.New("listing lineage activity")

	// ErrGetGeneration is returned when getting a generation fails.
	ErrGetGeneration = errors.New("getting generation")

	// ErrDeleteGeneration is returned when deleting a generation fails.
	ErrDeleteGeneration = errors.New("deleting generation")

	// ErrListArchivedGenerations is returned when listing archived generations fails.
	ErrListArchivedGenerations = errors.New("listing archived generations")

	// ErrUploadSignedCertificate is returned when uploading a signed certificate fails.
	ErrUploadSignedCertificate = errors.New("uploading signed certificate")

	// ErrPromoteLineage is returned when promoting a lineage fails.
	ErrPromoteLineage = errors.New("promoting lineage")

	// ErrRollbackLineage is returned when rolling back a lineage fails.
	ErrRollbackLineage = errors.New("rolling back lineage")

	// ErrReplaceStagingLineage is returned when replacing a lineage's staging generation fails.
	ErrReplaceStagingLineage = errors.New("replacing lineage staging generation")

	// ErrGetActivationStatus is returned when getting a lineage activation's status fails.
	ErrGetActivationStatus = errors.New("getting activation status")

	// ErrListActivations is returned when listing a lineage's activation requests fails.
	ErrListActivations = errors.New("listing activations")

	// ErrListLineageBindings is returned when listing a lineage's hostname bindings fails.
	ErrListLineageBindings = errors.New("listing lineage bindings")
)

var (
	// ErrLineageAccountNotAllowed represents an error when the account is not permitted to use Certificate Lineage.
	ErrLineageAccountNotAllowed = &Error{Type: "/error-types/lineage-account-not-allowed"}

	// ErrInternalError represents an unexpected server-side error.
	ErrInternalError = &Error{Type: "/error-types/internal-error"}

	// ErrCertAlreadyUploaded represents an error when a signed certificate was already uploaded for an algorithm instance.
	ErrCertAlreadyUploaded = &Error{Type: "/error-types/cert-already-uploaded"}

	// ErrGenerationImmutable represents an error when uploading a signed certificate to a generation that has
	// already been promoted to production (and is therefore immutable) and can no longer be modified.
	ErrGenerationImmutable = &Error{Type: "/error-types/generation-immutable"}

	// ErrLineageUploadValidationWarnings represents an error when the upload is rejected due to warnings that require acknowledgment.
	ErrLineageUploadValidationWarnings = &Error{Type: "/error-types/lineage-upload-validation-warnings"}

	// ErrCertExpiryInvalid represents an error when the uploaded certificate validity dates are invalid.
	ErrCertExpiryInvalid = &Error{Type: "/error-types/cert-expiry-invalid"}

	// ErrCertParseError represents an error when the uploaded certificate is malformed or cannot be parsed.
	ErrCertParseError = &Error{Type: "/error-types/cert-parse-error"}

	// ErrCertCSRMismatch represents an error when the uploaded certificate's public key does not match the
	// stored CSR's public key.
	ErrCertCSRMismatch = &Error{Type: "/error-types/cert-csr-mismatch"}

	// ErrUnknownKeyType represents an error when the given key type is not part of the lineage's key specs.
	ErrUnknownKeyType = &Error{Type: "/error-types/unknown-key-type"}

	// ErrDomainNotValidatedUploadFailed represents an error when the certificate upload is rejected because one
	// or more SANs are not yet Domain Validated.
	ErrDomainNotValidatedUploadFailed = &Error{Type: "/error-types/domain-not-validated-upload-failed"}

	// ErrCertificateNotFound represents an error when the certificate is not found.
	ErrCertificateNotFound = &Error{Type: "/error-types/certificate-not-found"}

	// ErrLineageNotFound represents an error when the certificate lineage is not found.
	ErrLineageNotFound = &Error{Type: "/error-types/lineage-not-found"}

	// ErrLineageNameConflict represents an error when a lineage with the given name already exists for the account.
	ErrLineageNameConflict = &Error{Type: "/error-types/lineage-name-conflict"}

	// ErrGenerationNotFound represents an error when the certificate generation is not found within the lineage.
	ErrGenerationNotFound = &Error{Type: "/error-types/generation-not-found"}

	// ErrGenerationDeployed represents an error when deleting a generation that is currently deployed on staging
	// or production (i.e. is the lineage's currentProduction or currentStaging generation, or otherwise ACTIVE).
	ErrGenerationDeployed = &Error{Type: "/error-types/generation-deployed"}

	// ErrGenerationDeleteConflict represents an error when deleting a generation that is the lineage's
	// previousProduction (rollback candidate) without setting AcknowledgeRollbackCandidateRemoval.
	ErrGenerationDeleteConflict = &Error{Type: "/error-types/generation-delete-conflict"}

	// ErrLastGeneration represents an error when attempting to delete the only remaining generation in a
	// lineage. Delete the lineage itself instead.
	ErrLastGeneration = &Error{Type: "/error-types/last-generation"}

	// ErrLineageHasActiveProduction represents an error when deleting a lineage that still has an active
	// production generation. The production generation must be deactivated first.
	ErrLineageHasActiveProduction = &Error{Type: "/error-types/lineage-has-active-production"}

	// ErrLineageHasActiveStaging represents an error when deleting a lineage that still has an active staging
	// generation. The staging generation must be deactivated first.
	ErrLineageHasActiveStaging = &Error{Type: "/error-types/lineage-has-active-staging"}

	// ErrPendingActivationInProgress represents an error when the lineage already has a pending activation (e.g.
	// a promotion or rollback) in progress that must complete or be waited out before starting another.
	ErrPendingActivationInProgress = &Error{Type: "/error-types/pending-activation-in-progress"}

	// ErrSingleGenerationNoRenew represents an error when renewing a SINGLE_GENERATION lineage that has no
	// production certificate. Only MULTIPLE_GENERATION lineages, or SINGLE_GENERATION lineages already deployed
	// to production, support renewal.
	ErrSingleGenerationNoRenew = &Error{Type: "/error-types/single-generation-no-renew"}

	// ErrHeadDeployedOnNetwork represents an error when renewing a lineage whose head generation is currently
	// deployed on staging or production. The head must be promoted or rolled back before renewal.
	ErrHeadDeployedOnNetwork = &Error{Type: "/error-types/head-deployed-on-network"}

	// ErrRenewBlockedHeadOnStaging represents an error when renewing a lineage that has both a head generation
	// and a staging generation present simultaneously. The head generation must be promoted or abandoned first.
	ErrRenewBlockedHeadOnStaging = &Error{Type: "/error-types/renew-blocked-head-on-staging"}

	// ErrConfirmationRequired represents an error when renewing a lineage whose head generation has a
	// still-valid, unexpired certificate. Set `ConfirmAbandonHead` on the request to explicitly abandon it.
	ErrConfirmationRequired = &Error{Type: "/error-types/confirmation-required"}

	// ErrNotMultipleStack represents an error when completing a lineage that is not MULTIPLE_STACK type. Only
	// MULTIPLE_STACK lineages support the complete operation.
	ErrNotMultipleStack = &Error{Type: "/error-types/not-multiple-stack"}

	// ErrNoCurrentProduction represents an error when completing a lineage that has no current production
	// generation to complete.
	ErrNoCurrentProduction = &Error{Type: "/error-types/no-current-production"}

	// ErrCompletePreconditionFailed represents an error when completing a lineage fails a precondition, e.g. a
	// pending head generation already exists, no unique CSR_READY algorithm was found, or the given key type
	// does not match the CSR_READY algorithm on current production.
	ErrCompletePreconditionFailed = &Error{Type: "/error-types/complete-precondition-failed"}

	// ErrCSRExpired represents an error when completing a lineage whose stored CSR for the given key type has
	// expired.
	ErrCSRExpired = &Error{Type: "/error-types/csr-expired"}

	// ErrCertSANMismatch represents an error when the uploaded certificate's SANs do not match the lineage's
	// SANs.
	ErrCertSANMismatch = &Error{Type: "/error-types/cert-san-mismatch"}

	// ErrCertSubjectMismatch represents an error when the uploaded certificate's non-common-name subject fields
	// differ from the stored CSR's subject fields.
	ErrCertSubjectMismatch = &Error{Type: "/error-types/cert-subject-mismatch"}

	// ErrMediaTypeNotSupported represents an error when the request media type is not supported.
	ErrMediaTypeNotSupported = &Error{Type: "/error-types/media-type-not-supported"}

	// ErrSchemaValidationFailure represents an error when the request body failed JSON schema validation.
	ErrSchemaValidationFailure = &Error{Type: "/error-types/schema-validation-failure"}

	// ErrInvalidCursor represents an error when the pagination `after` cursor is invalid or expired.
	ErrInvalidCursor = &Error{Type: "/error-types/invalid-cursor"}

	// ErrInvalidSortParameter represents an error when the `sort` query parameter references an unsupported value.
	ErrInvalidSortParameter = &Error{Type: "/error-types/invalid-sort-parameter"}

	// ErrLineageNoHeadGeneration represents an error when promoting a lineage that has no head generation to
	// activate.
	ErrLineageNoHeadGeneration = &Error{Type: "/error-types/lineage-no-head-generation"}

	// ErrFirstPromoteRequiresBothNetworks represents an error when a lineage that has never completed its first
	// promotion is promoted to only one of STAGING or PRODUCTION. The first PROMOTE must target both networks.
	ErrFirstPromoteRequiresBothNetworks = &Error{Type: "/error-types/first-promote-requires-both-networks"}

	// ErrSingleGenerationActivationNotSupported represents an error when performing an activation operation
	// (PROMOTE, ROLLBACK, or REPLACE_STAGING) on a SINGLE_GENERATION lineage. Activation operations are only
	// applicable to MULTIPLE_GENERATION lineages.
	ErrSingleGenerationActivationNotSupported = &Error{Type: "/error-types/single-generation-activation-not-supported"}

	// ErrActivationCooldownInEffect represents an error when starting an activation while a recent activation
	// for the same lineage is still propagating to the edge network. The error's `Context["cooldownEndsAt"]`
	// gives the time the cooldown ends.
	ErrActivationCooldownInEffect = &Error{Type: "/error-types/activation-cooldown-in-effect"}

	// ErrReplaceStagingHeadNotOnStaging represents an error when replacing a lineage's staging generation while
	// its head generation is not the lineage's current staging generation. REPLACE_STAGING is only valid when
	// the head is already deployed on staging.
	ErrReplaceStagingHeadNotOnStaging = &Error{Type: "/error-types/replace-staging-head-not-on-staging"}

	// ErrNoCurrentStaging represents an error when replacing a lineage's staging generation while it has no
	// current staging generation at all.
	ErrNoCurrentStaging = &Error{Type: "/error-types/no-current-staging"}

	// ErrLineageBadRequest represents a generic bad-request error with a dynamic, human-readable `reason` (e.g.
	// an unresolvable or invalid source generation for REPLACE_STAGING).
	ErrLineageBadRequest = &Error{Type: "/error-types/lineage-bad-request"}

	// ErrIncompleteCertMaterial represents an error when promoting a generation that has no algorithm instance
	// in READY_FOR_USE status. At least one signed certificate must be uploaded before promoting.
	ErrIncompleteCertMaterial = &Error{Type: "/error-types/incomplete-cert-material"}

	// ErrRollbackUnavailable represents an error when rolling back a lineage that has no previous production
	// generation to roll back to.
	ErrRollbackUnavailable = &Error{Type: "/error-types/rollback-unavailable"}

	// ErrUpstreamActivationError represents an error when an upstream system fails to process an activation
	// operation (PROMOTE, ROLLBACK, or REPLACE_STAGING). This is typically transient; retrying the operation
	// may succeed.
	ErrUpstreamActivationError = &Error{Type: "/error-types/upstream-activation-error"}

	// ErrInvalidField represents a generic invalid-field error, e.g. an unparsable ListActivations `cursor` value.
	ErrInvalidField = &Error{Type: "/error-types/invalid-field"}

	// ErrActivationNotFound represents an error when the requested activation ID does not exist for the lineage.
	ErrActivationNotFound = &Error{Type: "/error-types/activation-not-found"}

	// ErrForbidden represents a generic 403 error when the requesting user lacks the required permission for the
	// operation, e.g. `ACMI_CCM_READ_ONLY` on the lineage.
	ErrForbidden = &Error{Type: "/error-types/forbidden"}
)

type (
	// Error represents an error returned by the Cloud Certificate Manager API.
	Error struct {
		// Type is a URI reference that identifies the error type, e.g. `/error-types/lineage-not-found`. Used to
		// match sentinel errors via `errors.Is`.
		Type string `json:"type"`

		// Title is a short, human-readable summary of the error type, e.g. `Certificate lineage not found.`.
		Title string `json:"title"`

		// Status is the HTTP status code generated by the server for this occurrence of the error.
		Status int `json:"status"`

		// Detail is a human-readable explanation specific to this occurrence of the error, e.g. `Certificate
		// lineage {12345} not found.`.
		Detail string `json:"detail"`

		// Instance is a URI reference that identifies the specific occurrence of the error, typically including a
		// `traceId` query parameter for support/debugging purposes.
		Instance string `json:"instance"`

		// Context provides additional error-type-specific details, e.g. the offending keyType, generationId, or schema validation pointers.
		Context map[string]any `json:"context,omitempty"`
	}
)

// Error parses an error from the CloudCertificates API response.
func (c *cloudcertificates) Error(r *http.Response) error {
	var e Error
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.Status = r.StatusCode
		e.Title = "Failed to read error body"
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		c.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = "Failed to unmarshal error body. CCM API failed. Check details for more information."
		e.Detail = errs.UnescapeContent(string(body))
	}

	e.Status = r.StatusCode

	return &e
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s ", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons.
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}

	ignoreType := t.Type == ""
	ignoreStatus := t.Status == 0
	ignoreTitle := t.Title == ""
	matchType := t.Type == e.Type
	matchStatus := t.Status == e.Status
	matchTitle := t.Title == e.Title

	return (matchType || ignoreType) &&
		(matchStatus || ignoreStatus) &&
		(matchTitle || ignoreTitle)
}
