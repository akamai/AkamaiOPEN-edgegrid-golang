// Package mtlstruststore provides access to the Akamai mTLS Truststore v2 API.
package mtlstruststore

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

var (
	// ErrStructValidation is returned when given struct validation failed.
	ErrStructValidation = errors.New("struct validation")
)

type (
	// MTLSTruststore is the API interface for mTLS Truststore.
	MTLSTruststore interface {
		// CreateCASet creates a new CA set.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-ca-set
		CreateCASet(ctx context.Context, params CreateCASetRequest) (*CreateCASetResponse, error)

		// GetCASet returns details for a specific CA set.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set
		GetCASet(ctx context.Context, params GetCASetRequest) (*GetCASetResponse, error)

		// ListCASets returns detailed information about CA sets available to the current user account.
		//
		// https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-sets
		ListCASets(ctx context.Context, params ListCASetsRequest) (*ListCASetsResponse, error)

		// DeleteCASet deletes a CA set.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/delete-ca-set
		DeleteCASet(ctx context.Context, params DeleteCASetRequest) error

		// CreateCASetVersion creates a new CA set version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-ca-set-version
		CreateCASetVersion(ctx context.Context, params CreateCASetVersionRequest) (*CreateCASetVersionResponse, error)

		// CloneCASetVersion creates a clone of an existing CA set version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-clone-ca-set-version
		CloneCASetVersion(ctx context.Context, params CloneCASetVersionRequest) (*CloneCASetVersionResponse, error)

		// ListCASetVersions lists all the available CA sets created under the account.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set-versions
		ListCASetVersions(ctx context.Context, params ListCASetVersionsRequest) (*ListCASetVersionsResponse, error)

		// GetCASetVersion returns details of a CA sets version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set-version
		GetCASetVersion(ctx context.Context, params GetCASetVersionRequest) (*GetCASetVersionResponse, error)

		// UpdateCASetVersion updates a CA sets version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/put-ca-set-version
		UpdateCASetVersion(ctx context.Context, params UpdateCASetVersionRequest) (*UpdateCASetVersionResponse, error)

		// DeleteCASetVersion marks a CA set version for future deletion.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/delete-ca-set-version
		DeleteCASetVersion(ctx context.Context, params DeleteCASetVersionRequest) error

		// GetCASetVersionCertificates returns certificates details of a CA sets version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set-version-certificates
		GetCASetVersionCertificates(ctx context.Context, params GetCASetVersionCertificatesRequest) (*GetCASetVersionCertificatesResponse, error)

		// ActivateCASetVersion activates a CA set version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-activate-ca-set-version
		ActivateCASetVersion(ctx context.Context, params ActivateCASetVersionRequest) (*ActivateCASetVersionResponse, error)

		// DeactivateCASetVersion deactivates a CA set version.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-deactivate-ca-set-version
		DeactivateCASetVersion(ctx context.Context, params DeactivateCASetVersionRequest) (*DeactivateCASetVersionResponse, error)

		// GetCASetVersionActivation returns the status of a CA set version activation.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-deployment-request-details
		GetCASetVersionActivation(ctx context.Context, params GetCASetVersionActivationRequest) (*GetCASetVersionActivationResponse, error)

		// ListCASetActivations returns a list of CA set activations for a given CA set.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-deployment-request-for-ca-set
		ListCASetActivations(ctx context.Context, params ListCASetActivationsRequest) (*ListCASetActivationsResponse, error)

		// ListCASetVersionActivations returns a list of CA set version activations.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-version-deployment-request-details
		ListCASetVersionActivations(ctx context.Context, params ListCASetVersionActivationsRequest) (*ListCASetVersionActivationsResponse, error)

		// ListCASetAssociations provides of CA Set associations to a Certificate/Slot in CPS (in Commercial) or to a Hostname in Property Manager (in Defense Edge).
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set-associations
		ListCASetAssociations(ctx context.Context, params ListCASetAssociationsRequest) (*ListCASetAssociationsResponse, error)

		// CloneCASet clones a CA set with provided name and description.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-clone-ca-set
		CloneCASet(ctx context.Context, params CloneCASetRequest) (*CloneCASetResponse, error)

		// GetCASetDeletionStatus fetches the status of delete operation on both the networks.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-deletion-request-for-ca-set
		GetCASetDeletionStatus(ctx context.Context, params GetCASetDeletionStatusRequest) (*GetCASetDeletionStatusResponse, error)

		// ListCASetActivities returns the complete list of activities on a CA set.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/get-ca-set-activities
		ListCASetActivities(ctx context.Context, params ListCASetActivitiesRequest) (*ListCASetActivitiesResponse, error)

		// ValidateCertificates validates provided certificates if they are correct.
		//
		// See: https://techdocs.akamai.com/mtls-edge-truststore/reference/post-validate-certificates
		ValidateCertificates(ctx context.Context, params ValidateCertificatesRequest) (*ValidateCertificatesResponse, error)
	}

	mtlstruststore struct {
		session.Session
	}

	// Option defines an MTLS Truststore option.
	Option func(*mtlstruststore)
)

// Client returns a new mtlstruststore Client instance with the specified controller.
func Client(sess session.Session, opts ...Option) MTLSTruststore {
	c := &mtlstruststore{
		Session: sess,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
