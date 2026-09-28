// Package cloudcertificates provides access to the Akamai Cloud Certificate Manager API.
package cloudcertificates

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

// CloudCertificates defines the interface for Akamai Cloud Certificate Manager operations.
type (
	CloudCertificates interface {
		// CreateLineage creates a new certificate lineage. The API generates a certificate signing request (CSR)
		// for each requested key type, returned as the head generation of the lineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		CreateLineage(ctx context.Context, params CreateLineageRequest) (*CreateLineageResponse, error)

		// GetLineage retrieves a single certificate lineage by its lineageID, optionally expanding its generations.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		GetLineage(ctx context.Context, params GetLineageRequest) (*GetLineageResponse, error)

		// ListLineages returns a paginated list of certificate lineages accessible to the requesting user,
		// filtered and sorted by the given query parameters.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ListLineages(ctx context.Context, params ListLineagesRequest) (*ListLineagesResponse, error)

		// RenameLineage updates the name of a certificate lineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		RenameLineage(ctx context.Context, params RenameLineageRequest) (*RenameLineageResponse, error)

		// DeleteLineage permanently deletes a certificate lineage. The lineage must have no active production or
		// staging generation and no pending activation in progress.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		DeleteLineage(ctx context.Context, params DeleteLineageRequest) error

		// ListLineageActivity returns a paginated, newest-first list of activity events (e.g. creation, renames,
		// certificate uploads, promotions) recorded for a certificate lineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ListLineageActivity(ctx context.Context, params ListLineageActivityRequest) (*ListLineageActivityResponse, error)

		// RenewLineage creates a new head generation for a lineage, generating fresh CSRs for all key types in the
		// lineage's key specs. Any existing head generation is abandoned as part of the renewal.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		RenewLineage(ctx context.Context, params RenewLineageRequest) (*RenewLineageResponse, error)

		// CompleteLineage adds the second algorithm (e.g. ECDSA) to a MULTIPLE_STACK lineage where only one
		// algorithm (e.g. RSA) is currently live on production, creating a new generation with the completed
		// algorithm ready for use.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		CompleteLineage(ctx context.Context, params CompleteLineageRequest) (*CompleteLineageResponse, error)

		// GetGeneration retrieves a single certificate generation within a lineage, including full algorithm
		// instance details (status, expiry dates, and key metadata) for that generation. Works for any generation
		// role: head, currentProduction, previousProduction, or an archived generation.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		GetGeneration(ctx context.Context, params GetGenerationRequest) (*GetGenerationResponse, error)

		// DeleteGeneration permanently deletes a single certificate generation within a lineage. The generation
		// must not be currently deployed on staging or production, and must not be the only generation in the
		// lineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		DeleteGeneration(ctx context.Context, params DeleteGenerationRequest) error

		// ListArchivedGenerations returns the unpaginated list of archived and abandoned generations for a
		// certificate lineage, i.e. every generation except the current head, currentProduction, currentStaging,
		// and previousProduction pointers.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ListArchivedGenerations(ctx context.Context, params ListArchivedGenerationsRequest) (*ListArchivedGenerationsResponse, error)

		// UploadSignedCertificate uploads the signed certificate for one or more key types of a lineage generation.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		UploadSignedCertificate(ctx context.Context, params UploadSignedCertificateRequest) (*UploadSignedCertificateResponse, error)

		// PromoteLineage promotes a lineage's head generation to one or more networks (staging and/or production).
		// GenerationID is mandatory and must match the lineage's head generation. If no networks are given, the
		// API defaults promoting to production only. If the head generation is already active on every
		// requested network, this is a no-op and returns a nil response and nil error.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		PromoteLineage(ctx context.Context, params PromoteLineageRequest) (*PromoteLineageResponse, error)

		// RollbackLineage rolls a lineage's production generation back to its previous production generation.
		// GenerationID is mandatory and must match the lineage's previous production generation. If that
		// generation is already the current production generation, this is a no-op and returns a nil response
		// and nil error.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		RollbackLineage(ctx context.Context, params RollbackLineageRequest) (*RollbackLineageResponse, error)

		// ReplaceStagingLineage replaces the generation currently deployed on staging with the lineage's current
		// production or previous production generation. GenerationID is mandatory and must match one of those
		// two generations. If it already matches the current staging generation, this is a no-op and returns a
		// nil response and nil error.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ReplaceStagingLineage(ctx context.Context, params ReplaceStagingLineageRequest) (*ReplaceStagingLineageResponse, error)

		// GetActivationStatus retrieves the status of a lineage activation request (PROMOTE, ROLLBACK, or
		// REPLACE_STAGING), as returned by PromoteLineage, RollbackLineage, or ReplaceStagingLineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		GetActivationStatus(ctx context.Context, params GetActivationStatusRequest) (*GetActivationStatusResponse, error)

		// ListActivations returns a paginated, newest-first list of activation requests (PROMOTE, ROLLBACK, or
		// REPLACE_STAGING) recorded for a certificate lineage.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ListActivations(ctx context.Context, params ListActivationsRequest) (*ListActivationsResponse, error)

		// ListLineageBindings returns a paginated list of hostname bindings for a certificate lineage, optionally
		// filtered by network and sorted chronologically.
		//
		// See: public API documentation for this endpoint isn't available yet. It's planned for a later release.
		ListLineageBindings(ctx context.Context, params ListLineageBindingsRequest) (*ListLineageBindingsResponse, error)
	}

	cloudcertificates struct {
		session.Session
	}

	// Option is a function that configures the CloudCertificates.
	Option func(*cloudcertificates)
)

// Client creates a new CloudCertificates client.
func Client(sess session.Session, opts ...Option) CloudCertificates {
	c := &cloudcertificates{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
