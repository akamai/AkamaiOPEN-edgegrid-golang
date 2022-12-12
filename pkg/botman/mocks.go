//revive:disable:exported

package botman

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ BotMan = &Mock{}

func (p *Mock) GetAkamaiBotCategoryList(ctx context.Context, params GetAkamaiBotCategoryListRequest) (*GetAkamaiBotCategoryListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAkamaiBotCategoryListResponse), nil
}

func (p *Mock) GetAkamaiBotCategoryActionList(ctx context.Context, params GetAkamaiBotCategoryActionListRequest) (*GetAkamaiBotCategoryActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAkamaiBotCategoryActionListResponse), nil
}
func (p *Mock) GetAkamaiBotCategoryAction(ctx context.Context, params GetAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateAkamaiBotCategoryAction(ctx context.Context, params UpdateAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]interface{}), nil
}

func (p *Mock) GetAkamaiDefinedBotList(ctx context.Context, params GetAkamaiDefinedBotListRequest) (*GetAkamaiDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetAkamaiDefinedBotListResponse), nil
}
func (p *Mock) GetBotAnalyticsCookie(ctx context.Context, params GetBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateBotAnalyticsCookie(ctx context.Context, params UpdateBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetBotAnalyticsCookieValues(ctx context.Context) (map[string]interface{}, error) {
	args := p.Called(ctx)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetBotCategoryException(ctx context.Context, params GetBotCategoryExceptionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateBotCategoryException(ctx context.Context, params UpdateBotCategoryExceptionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetBotDetectionActionList(ctx context.Context, params GetBotDetectionActionListRequest) (*GetBotDetectionActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBotDetectionActionListResponse), nil
}
func (p *Mock) GetBotDetectionAction(ctx context.Context, params GetBotDetectionActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateBotDetectionAction(ctx context.Context, params UpdateBotDetectionActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetBotEndpointCoverageReport(ctx context.Context, params GetBotEndpointCoverageReportRequest) (*GetBotEndpointCoverageReportResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBotEndpointCoverageReportResponse), nil
}
func (p *Mock) GetBotManagementSetting(ctx context.Context, params GetBotManagementSettingRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateBotManagementSetting(ctx context.Context, params UpdateBotManagementSettingRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetChallengeActionList(ctx context.Context, params GetChallengeActionListRequest) (*GetChallengeActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetChallengeActionListResponse), nil
}
func (p *Mock) GetChallengeAction(ctx context.Context, params GetChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateChallengeAction(ctx context.Context, params CreateChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateChallengeAction(ctx context.Context, params UpdateChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveChallengeAction(ctx context.Context, params RemoveChallengeActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) UpdateGoogleReCaptchaSecretKey(ctx context.Context, params UpdateGoogleReCaptchaSecretKeyRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetChallengeInterceptionRules(ctx context.Context, params GetChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateChallengeInterceptionRules(ctx context.Context, params UpdateChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetClientSideSecurity(ctx context.Context, params GetClientSideSecurityRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateClientSideSecurity(ctx context.Context, params UpdateClientSideSecurityRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetConditionalActionList(ctx context.Context, params GetConditionalActionListRequest) (*GetConditionalActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConditionalActionListResponse), nil
}
func (p *Mock) GetConditionalAction(ctx context.Context, params GetConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateConditionalAction(ctx context.Context, params CreateConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateConditionalAction(ctx context.Context, params UpdateConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveConditionalAction(ctx context.Context, params RemoveConditionalActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetCustomBotCategoryList(ctx context.Context, params GetCustomBotCategoryListRequest) (*GetCustomBotCategoryListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomBotCategoryListResponse), nil
}
func (p *Mock) GetCustomBotCategory(ctx context.Context, params GetCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateCustomBotCategory(ctx context.Context, params CreateCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateCustomBotCategory(ctx context.Context, params UpdateCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveCustomBotCategory(ctx context.Context, params RemoveCustomBotCategoryRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetCustomBotCategoryActionList(ctx context.Context, params GetCustomBotCategoryActionListRequest) (*GetCustomBotCategoryActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomBotCategoryActionListResponse), nil
}
func (p *Mock) GetCustomBotCategoryAction(ctx context.Context, params GetCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateCustomBotCategoryAction(ctx context.Context, params UpdateCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetCustomBotCategorySequence(ctx context.Context, params GetCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CustomBotCategorySequenceResponse), nil
}
func (p *Mock) UpdateCustomBotCategorySequence(ctx context.Context, params UpdateCustomBotCategorySequenceRequest) (*CustomBotCategorySequenceResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CustomBotCategorySequenceResponse), nil
}
func (p *Mock) GetCustomClientList(ctx context.Context, params GetCustomClientListRequest) (*GetCustomClientListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomClientListResponse), nil
}
func (p *Mock) GetCustomClient(ctx context.Context, params GetCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateCustomClient(ctx context.Context, params CreateCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateCustomClient(ctx context.Context, params UpdateCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveCustomClient(ctx context.Context, params RemoveCustomClientRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetCustomDefinedBotList(ctx context.Context, params GetCustomDefinedBotListRequest) (*GetCustomDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomDefinedBotListResponse), nil
}
func (p *Mock) GetCustomDefinedBot(ctx context.Context, params GetCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateCustomDefinedBot(ctx context.Context, params CreateCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateCustomDefinedBot(ctx context.Context, params UpdateCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveCustomDefinedBot(ctx context.Context, params RemoveCustomDefinedBotRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetCustomDenyActionList(ctx context.Context, params GetCustomDenyActionListRequest) (*GetCustomDenyActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomDenyActionListResponse), nil
}
func (p *Mock) GetCustomDenyAction(ctx context.Context, params GetCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateCustomDenyAction(ctx context.Context, params CreateCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateCustomDenyAction(ctx context.Context, params UpdateCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveCustomDenyAction(ctx context.Context, params RemoveCustomDenyActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}

func (p *Mock) GetRecategorizedAkamaiDefinedBotList(ctx context.Context, params GetRecategorizedAkamaiDefinedBotListRequest) (*GetRecategorizedAkamaiDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRecategorizedAkamaiDefinedBotListResponse), nil
}
func (p *Mock) GetRecategorizedAkamaiDefinedBot(ctx context.Context, params GetRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *Mock) CreateRecategorizedAkamaiDefinedBot(ctx context.Context, params CreateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *Mock) UpdateRecategorizedAkamaiDefinedBot(ctx context.Context, params UpdateRecategorizedAkamaiDefinedBotRequest) (*RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *Mock) RemoveRecategorizedAkamaiDefinedBot(ctx context.Context, params RemoveRecategorizedAkamaiDefinedBotRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetResponseActionList(ctx context.Context, params GetResponseActionListRequest) (*GetResponseActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetResponseActionListResponse), nil
}

func (p *Mock) GetTransactionalEndpointList(ctx context.Context, params GetTransactionalEndpointListRequest) (*GetTransactionalEndpointListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetTransactionalEndpointListResponse), nil
}
func (p *Mock) GetTransactionalEndpoint(ctx context.Context, params GetTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateTransactionalEndpoint(ctx context.Context, params CreateTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateTransactionalEndpoint(ctx context.Context, params UpdateTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveTransactionalEndpoint(ctx context.Context, params RemoveTransactionalEndpointRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *Mock) GetTransactionalEndpointProtection(ctx context.Context, params GetTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateTransactionalEndpointProtection(ctx context.Context, params UpdateTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}

func (p *Mock) GetJavascriptInjection(ctx context.Context, params GetJavascriptInjectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateJavascriptInjection(ctx context.Context, params UpdateJavascriptInjectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) GetServeAlternateActionList(ctx context.Context, params GetServeAlternateActionListRequest) (*GetServeAlternateActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetServeAlternateActionListResponse), nil
}
func (p *Mock) GetServeAlternateAction(ctx context.Context, params GetServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) CreateServeAlternateAction(ctx context.Context, params CreateServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) UpdateServeAlternateAction(ctx context.Context, params UpdateServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *Mock) RemoveServeAlternateAction(ctx context.Context, params RemoveServeAlternateActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}

func (p *Mock) GetBotDetectionList(ctx context.Context, params GetBotDetectionListRequest) (*GetBotDetectionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBotDetectionListResponse), nil
}
