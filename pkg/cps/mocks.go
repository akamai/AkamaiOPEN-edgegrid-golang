//revive:disable:exported

package cps

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ CPS = &Mock{}

func (m *Mock) ListEnrollments(ctx context.Context, r ListEnrollmentsRequest) (*ListEnrollmentsResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListEnrollmentsResponse), args.Error(1)
}

func (m *Mock) GetEnrollment(ctx context.Context, r GetEnrollmentRequest) (*Enrollment, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Enrollment), args.Error(1)
}

func (m *Mock) CreateEnrollment(ctx context.Context, r CreateEnrollmentRequest) (*CreateEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateEnrollmentResponse), args.Error(1)
}

func (m *Mock) UpdateEnrollment(ctx context.Context, r UpdateEnrollmentRequest) (*UpdateEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateEnrollmentResponse), args.Error(1)
}

func (m *Mock) RemoveEnrollment(ctx context.Context, r RemoveEnrollmentRequest) (*RemoveEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RemoveEnrollmentResponse), args.Error(1)
}

func (m *Mock) GetChangeStatus(ctx context.Context, r GetChangeStatusRequest) (*Change, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Change), args.Error(1)
}

func (m *Mock) CancelChange(ctx context.Context, r CancelChangeRequest) (*CancelChangeResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CancelChangeResponse), args.Error(1)
}

func (m *Mock) UpdateChange(ctx context.Context, r UpdateChangeRequest) (*UpdateChangeResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateChangeResponse), args.Error(1)
}

func (m *Mock) GetChangeLetsEncryptChallenges(ctx context.Context, r GetChangeRequest) (*DVArray, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DVArray), args.Error(1)
}

func (m *Mock) GetChangePreVerificationWarnings(ctx context.Context, r GetChangeRequest) (*PreVerificationWarnings, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*PreVerificationWarnings), args.Error(1)
}

func (m *Mock) AcknowledgeDVChallenges(ctx context.Context, r AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}

func (m *Mock) AcknowledgePreVerificationWarnings(ctx context.Context, r AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}

func (m *Mock) GetChangeManagementInfo(ctx context.Context, r GetChangeRequest) (*ChangeManagementInfoResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ChangeManagementInfoResponse), args.Error(1)
}

func (m *Mock) GetChangeDeploymentInfo(ctx context.Context, r GetChangeRequest) (*ChangeDeploymentInfoResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ChangeDeploymentInfoResponse), args.Error(1)
}

func (m *Mock) ListDeployments(ctx context.Context, r ListDeploymentsRequest) (*ListDeploymentsResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListDeploymentsResponse), args.Error(1)
}

func (m *Mock) GetProductionDeployment(ctx context.Context, r GetDeploymentRequest) (*GetProductionDeploymentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetProductionDeploymentResponse), args.Error(1)
}

func (m *Mock) GetStagingDeployment(ctx context.Context, r GetDeploymentRequest) (*GetStagingDeploymentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetStagingDeploymentResponse), args.Error(1)
}

func (m *Mock) GetDeploymentSchedule(ctx context.Context, r GetDeploymentScheduleRequest) (*DeploymentSchedule, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeploymentSchedule), args.Error(1)
}

func (m *Mock) UpdateDeploymentSchedule(ctx context.Context, r UpdateDeploymentScheduleRequest) (*UpdateDeploymentScheduleResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateDeploymentScheduleResponse), args.Error(1)
}

func (m *Mock) GetDVHistory(ctx context.Context, r GetDVHistoryRequest) (*GetDVHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetDVHistoryResponse), args.Error(1)
}

func (m *Mock) GetCertificateHistory(ctx context.Context, r GetCertificateHistoryRequest) (*GetCertificateHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCertificateHistoryResponse), args.Error(1)
}

func (m *Mock) GetChangeHistory(ctx context.Context, r GetChangeHistoryRequest) (*GetChangeHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetChangeHistoryResponse), args.Error(1)
}

func (m *Mock) GetChangePostVerificationWarnings(ctx context.Context, r GetChangeRequest) (*PostVerificationWarnings, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*PostVerificationWarnings), args.Error(1)
}

func (m *Mock) GetChangeThirdPartyCSR(ctx context.Context, r GetChangeRequest) (*ThirdPartyCSRResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ThirdPartyCSRResponse), args.Error(1)
}

func (m *Mock) AcknowledgeChangeManagement(ctx context.Context, r AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}

func (m *Mock) AcknowledgePostVerificationWarnings(ctx context.Context, r AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}

func (m *Mock) UploadThirdPartyCertAndTrustChain(ctx context.Context, r UploadThirdPartyCertAndTrustChainRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}
