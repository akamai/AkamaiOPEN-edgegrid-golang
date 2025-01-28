// Package cps provides access to the Akamai CPS APIs
package cps

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CPS is the cps api interface
	CPS interface {
		// ChangeManagementInfo

		// GetChangeManagementInfo gets information about acknowledgement status,
		// and may include warnings about potential conflicts that may occur if you proceed with acknowledgement
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeManagementInfo(ctx context.Context, params GetChangeRequest) (*ChangeManagementInfoResponse, error)

		// GetChangeDeploymentInfo gets deployment currently deployed to the staging network
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeDeploymentInfo(ctx context.Context, params GetChangeRequest) (*ChangeDeploymentInfoResponse, error)

		// AcknowledgeChangeManagement sends acknowledgement request to CPS to proceed deploying the certificate to the production network
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgeChangeManagement(context.Context, AcknowledgementRequest) error

		// GetChangeStatus fetches change status for given enrollment and change ID
		//
		// See: https://techdocs.akamai.com/cps/reference/get-enrollment-change
		GetChangeStatus(context.Context, GetChangeStatusRequest) (*Change, error)

		// ChangeOperations

		// CancelChange cancels a pending change
		//
		// See: https://techdocs.akamai.com/cps/reference/delete-enrollment-change
		CancelChange(context.Context, CancelChangeRequest) (*CancelChangeResponse, error)

		// Deployments

		// ListDeployments fetches deployments for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployments
		ListDeployments(context.Context, ListDeploymentsRequest) (*ListDeploymentsResponse, error)

		// GetProductionDeployment fetches production deployment for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployments-production
		GetProductionDeployment(context.Context, GetDeploymentRequest) (*GetProductionDeploymentResponse, error)

		// GetStagingDeployment fetches staging deployment for given enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-deployment-staging
		GetStagingDeployment(context.Context, GetDeploymentRequest) (*GetStagingDeploymentResponse, error)

		// DeploymentSchedules

		// GetDeploymentSchedule fetches the current deployment schedule settings describing when a change deploys to the network
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-deployment-schedule
		GetDeploymentSchedule(context.Context, GetDeploymentScheduleRequest) (*DeploymentSchedule, error)

		// UpdateDeploymentSchedule updates the current deployment schedule
		//
		// See: https://techdocs.akamai.com/cps/reference/put-change-deployment-schedule
		UpdateDeploymentSchedule(context.Context, UpdateDeploymentScheduleRequest) (*UpdateDeploymentScheduleResponse, error)

		// DVChallenges

		// GetChangeLetsEncryptChallenges gets detailed information about Domain Validation challenges
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeLetsEncryptChallenges(context.Context, GetChangeRequest) (*DVArray, error)

		// AcknowledgeDVChallenges sends acknowledgement request to CPS informing that the validation is completed
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgeDVChallenges(context.Context, AcknowledgementRequest) error

		// Enrollments

		// ListEnrollments fetches all enrollments with given contractId
		//
		// See https://techdocs.akamai.com/cps/reference/get-enrollments
		ListEnrollments(context.Context, ListEnrollmentsRequest) (*ListEnrollmentsResponse, error)

		// GetEnrollment fetches enrollment object with given ID
		//
		// See: https://techdocs.akamai.com/cps/reference/get-enrollment
		GetEnrollment(context.Context, GetEnrollmentRequest) (*GetEnrollmentResponse, error)

		// CreateEnrollment creates a new enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/post-enrollment
		CreateEnrollment(context.Context, CreateEnrollmentRequest) (*CreateEnrollmentResponse, error)

		// UpdateEnrollment updates a single enrollment entry with given ID
		//
		// See: https://techdocs.akamai.com/cps/reference/put-enrollment
		UpdateEnrollment(context.Context, UpdateEnrollmentRequest) (*UpdateEnrollmentResponse, error)

		// RemoveEnrollment removes an enrollment with given ID
		//
		// See: https://techdocs.akamai.com/cps/reference/delete-enrollment
		RemoveEnrollment(context.Context, RemoveEnrollmentRequest) (*RemoveEnrollmentResponse, error)

		// History

		// GetDVHistory is a domain name validation history for the enrollment
		//
		// See: https://techdocs.akamai.com/cps/reference/get-dv-history
		GetDVHistory(context.Context, GetDVHistoryRequest) (*GetDVHistoryResponse, error)

		// GetCertificateHistory views the certificate history.
		//
		// See: https://techdocs.akamai.com/cps/reference/get-history-certificates
		GetCertificateHistory(context.Context, GetCertificateHistoryRequest) (*GetCertificateHistoryResponse, error)

		// GetChangeHistory views the change history for enrollment.
		//
		// See: https://techdocs.akamai.com/cps/reference/get-history-changes
		GetChangeHistory(context.Context, GetChangeHistoryRequest) (*GetChangeHistoryResponse, error)

		// PostVerification

		// GetChangePostVerificationWarnings gets information about post verification warnings
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangePostVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PostVerificationWarnings, error)
		// AcknowledgePostVerificationWarnings sends acknowledgement request to CPS informing that the warnings should be ignored
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgePostVerificationWarnings(context.Context, AcknowledgementRequest) error

		// PreVerification

		// GetChangePreVerificationWarnings gets detailed information about Domain Validation challenges
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangePreVerificationWarnings(ctx context.Context, params GetChangeRequest) (*PreVerificationWarnings, error)

		// AcknowledgePreVerificationWarnings sends acknowledgement request to CPS informing that the warnings should be ignored
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		AcknowledgePreVerificationWarnings(context.Context, AcknowledgementRequest) error

		// ThirdPartyCSR

		// GetChangeThirdPartyCSR gets certificate signing request
		//
		// See: https://techdocs.akamai.com/cps/reference/get-change-allowed-input-param
		GetChangeThirdPartyCSR(ctx context.Context, params GetChangeRequest) (*ThirdPartyCSRResponse, error)

		// UploadThirdPartyCertAndTrustChain uploads signed certificate and trust chain to cps
		//
		// See: https://techdocs.akamai.com/cps/reference/post-change-allowed-input-param
		UploadThirdPartyCertAndTrustChain(context.Context, UploadThirdPartyCertAndTrustChainRequest) error
	}

	cps struct {
		session.Session
	}

	// Option defines a CPS option
	Option func(*cps)

	// ClientFunc is a cps client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) CPS
)

// Client returns a new cps Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CPS {
	c := &cps{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
