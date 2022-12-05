//revive:disable:exported

package appsec

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ APPSEC = &Mock{}

func (m *Mock) UpdateWAPSelectedHostnames(ctx context.Context, req UpdateWAPSelectedHostnamesRequest) (*UpdateWAPSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateWAPSelectedHostnamesResponse), args.Error(1)
}

func (m *Mock) UpdateWAFProtection(ctx context.Context, req UpdateWAFProtectionRequest) (*UpdateWAFProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateWAFProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateWAFMode(ctx context.Context, req UpdateWAFModeRequest) (*UpdateWAFModeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateWAFModeResponse), args.Error(1)
}

func (m *Mock) UpdateVersionNotes(ctx context.Context, req UpdateVersionNotesRequest) (*UpdateVersionNotesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateVersionNotesResponse), args.Error(1)
}

func (m *Mock) GetRuleRecommendations(ctx context.Context, params GetRuleRecommendationsRequest) (*GetRuleRecommendationsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRuleRecommendationsResponse), args.Error(1)
}

func (m *Mock) UpdateThreatIntel(ctx context.Context, req UpdateThreatIntelRequest) (*UpdateThreatIntelResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateThreatIntelResponse), args.Error(1)
}

func (m *Mock) UpdateSlowPostProtectionSetting(ctx context.Context, req UpdateSlowPostProtectionSettingRequest) (*UpdateSlowPostProtectionSettingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSlowPostProtectionSettingResponse), args.Error(1)
}

func (m *Mock) UpdateSlowPostProtection(ctx context.Context, req UpdateSlowPostProtectionRequest) (*UpdateSlowPostProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSlowPostProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateSiemSettings(ctx context.Context, req UpdateSiemSettingsRequest) (*UpdateSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSiemSettingsResponse), args.Error(1)
}

func (m *Mock) UpdateSelectedHostnames(ctx context.Context, req UpdateSelectedHostnamesRequest) (*UpdateSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSelectedHostnamesResponse), args.Error(1)
}

func (m *Mock) UpdateSelectedHostname(ctx context.Context, req UpdateSelectedHostnameRequest) (*UpdateSelectedHostnameResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSelectedHostnameResponse), args.Error(1)
}

func (m *Mock) UpdateSecurityPolicy(ctx context.Context, req UpdateSecurityPolicyRequest) (*UpdateSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateSecurityPolicyResponse), args.Error(1)
}

func (m *Mock) UpdateRuleUpgrade(ctx context.Context, req UpdateRuleUpgradeRequest) (*UpdateRuleUpgradeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateRuleUpgradeResponse), args.Error(1)
}

func (m *Mock) UpdateRuleConditionException(ctx context.Context, req UpdateConditionExceptionRequest) (*UpdateConditionExceptionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateConditionExceptionResponse), args.Error(1)
}

func (m *Mock) UpdateRule(ctx context.Context, req UpdateRuleRequest) (*UpdateRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateRuleResponse), args.Error(1)
}

func (m *Mock) UpdateReputationProtection(ctx context.Context, req UpdateReputationProtectionRequest) (*UpdateReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateReputationProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateReputationProfileAction(ctx context.Context, req UpdateReputationProfileActionRequest) (*UpdateReputationProfileActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateReputationProfileActionResponse), args.Error(1)
}

func (m *Mock) UpdateReputationProfile(ctx context.Context, req UpdateReputationProfileRequest) (*UpdateReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateReputationProfileResponse), args.Error(1)
}

func (m *Mock) UpdateReputationAnalysis(ctx context.Context, req UpdateReputationAnalysisRequest) (*UpdateReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateReputationAnalysisResponse), args.Error(1)
}

func (m *Mock) UpdateRateProtection(ctx context.Context, req UpdateRateProtectionRequest) (*UpdateRateProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateRateProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateRatePolicyAction(ctx context.Context, req UpdateRatePolicyActionRequest) (*UpdateRatePolicyActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateRatePolicyActionResponse), args.Error(1)
}

func (m *Mock) UpdateRatePolicy(ctx context.Context, req UpdateRatePolicyRequest) (*UpdateRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateRatePolicyResponse), args.Error(1)
}

func (m *Mock) UpdatePolicyProtections(ctx context.Context, req UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyProtectionsResponse), args.Error(1)
}

func (m *Mock) UpdatePenaltyBox(ctx context.Context, req UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdatePenaltyBoxResponse), args.Error(1)
}

func (m *Mock) UpdateNetworkLayerProtection(ctx context.Context, req UpdateNetworkLayerProtectionRequest) (*UpdateNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateNetworkLayerProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateMatchTargetSequence(ctx context.Context, req UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateMatchTargetSequenceResponse), args.Error(1)
}

func (m *Mock) UpdateMatchTarget(ctx context.Context, req UpdateMatchTargetRequest) (*UpdateMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateMatchTargetResponse), args.Error(1)
}

func (m *Mock) UpdateIPGeoProtection(ctx context.Context, req UpdateIPGeoProtectionRequest) (*UpdateIPGeoProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateIPGeoProtectionResponse), args.Error(1)
}

func (m *Mock) UpdateIPGeo(ctx context.Context, req UpdateIPGeoRequest) (*UpdateIPGeoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateIPGeoResponse), args.Error(1)
}

func (m *Mock) UpdateEvalRule(ctx context.Context, req UpdateEvalRuleRequest) (*UpdateEvalRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateEvalRuleResponse), args.Error(1)
}

func (m *Mock) UpdateEvalProtectHost(ctx context.Context, req UpdateEvalProtectHostRequest) (*UpdateEvalProtectHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateEvalProtectHostResponse), args.Error(1)
}

func (m *Mock) UpdateEvalHost(ctx context.Context, req UpdateEvalHostRequest) (*UpdateEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateEvalHostResponse), args.Error(1)
}

func (m *Mock) UpdateEvalGroup(ctx context.Context, req UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAttackGroupResponse), args.Error(1)
}

func (m *Mock) UpdateEval(ctx context.Context, req UpdateEvalRequest) (*UpdateEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateEvalResponse), args.Error(1)
}

func (m *Mock) UpdateCustomRuleAction(ctx context.Context, req UpdateCustomRuleActionRequest) (*UpdateCustomRuleActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateCustomRuleActionResponse), args.Error(1)
}

func (m *Mock) UpdateCustomRule(ctx context.Context, req UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateCustomRuleResponse), args.Error(1)
}

func (m *Mock) UpdateCustomDeny(ctx context.Context, req UpdateCustomDenyRequest) (*UpdateCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateCustomDenyResponse), args.Error(1)
}

func (m *Mock) UpdateConfiguration(ctx context.Context, req UpdateConfigurationRequest) (*UpdateConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateConfigurationResponse), args.Error(1)
}

func (m *Mock) UpdateBypassNetworkLists(ctx context.Context, req UpdateBypassNetworkListsRequest) (*UpdateBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) UpdateAttackGroup(ctx context.Context, req UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAttackGroupResponse), args.Error(1)
}

func (m *Mock) UpdateApiRequestConstraints(ctx context.Context, req UpdateApiRequestConstraintsRequest) (*UpdateApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateApiRequestConstraintsResponse), args.Error(1)
}

func (m *Mock) UpdateAdvancedSettingsPrefetch(ctx context.Context, req UpdateAdvancedSettingsPrefetchRequest) (*UpdateAdvancedSettingsPrefetchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAdvancedSettingsPrefetchResponse), args.Error(1)
}

func (m *Mock) UpdateAdvancedSettingsPragma(ctx context.Context, req UpdateAdvancedSettingsPragmaRequest) (*UpdateAdvancedSettingsPragmaResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAdvancedSettingsPragmaResponse), args.Error(1)
}

func (m *Mock) UpdateAdvancedSettingsLogging(ctx context.Context, req UpdateAdvancedSettingsLoggingRequest) (*UpdateAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *Mock) UpdateAdvancedSettingsEvasivePathMatch(ctx context.Context, req UpdateAdvancedSettingsEvasivePathMatchRequest) (*UpdateAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *Mock) UpdateAPIConstraintsProtection(ctx context.Context, req UpdateAPIConstraintsProtectionRequest) (*UpdateAPIConstraintsProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateAPIConstraintsProtectionResponse), args.Error(1)
}

func (m *Mock) RemoveSiemSettings(ctx context.Context, req RemoveSiemSettingsRequest) (*RemoveSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveSiemSettingsResponse), args.Error(1)
}

func (m *Mock) RemoveSecurityPolicy(ctx context.Context, req RemoveSecurityPolicyRequest) (*RemoveSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveSecurityPolicyResponse), args.Error(1)
}

func (m *Mock) RemoveReputationProtection(ctx context.Context, req RemoveReputationProtectionRequest) (*RemoveReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveReputationProtectionResponse), args.Error(1)
}

func (m *Mock) RemoveReputationProfile(ctx context.Context, req RemoveReputationProfileRequest) (*RemoveReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveReputationProfileResponse), args.Error(1)
}

func (m *Mock) RemoveReputationAnalysis(ctx context.Context, req RemoveReputationAnalysisRequest) (*RemoveReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveReputationAnalysisResponse), args.Error(1)
}

func (m *Mock) RemoveRatePolicy(ctx context.Context, req RemoveRatePolicyRequest) (*RemoveRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveRatePolicyResponse), args.Error(1)
}

func (m *Mock) RemovePolicyProtections(ctx context.Context, req UpdatePolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyProtectionsResponse), args.Error(1)
}

func (m *Mock) RemoveNetworkLayerProtection(ctx context.Context, req RemoveNetworkLayerProtectionRequest) (*RemoveNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveNetworkLayerProtectionResponse), args.Error(1)
}

func (m *Mock) RemoveMatchTarget(ctx context.Context, req RemoveMatchTargetRequest) (*RemoveMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveMatchTargetResponse), args.Error(1)
}

func (m *Mock) RemoveEval(ctx context.Context, req RemoveEvalRequest) (*RemoveEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveEvalResponse), args.Error(1)
}

func (m *Mock) RemoveEvalHost(ctx context.Context, req RemoveEvalHostRequest) (*RemoveEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveEvalHostResponse), args.Error(1)
}

func (m *Mock) RemoveCustomRule(ctx context.Context, req RemoveCustomRuleRequest) (*RemoveCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveCustomRuleResponse), args.Error(1)
}

func (m *Mock) RemoveCustomDeny(ctx context.Context, req RemoveCustomDenyRequest) (*RemoveCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveCustomDenyResponse), args.Error(1)
}

func (m *Mock) RemoveConfigurationVersionClone(ctx context.Context, req RemoveConfigurationVersionCloneRequest) (*RemoveConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveConfigurationVersionCloneResponse), args.Error(1)
}

func (m *Mock) RemoveConfiguration(ctx context.Context, req RemoveConfigurationRequest) (*RemoveConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveConfigurationResponse), args.Error(1)
}

func (m *Mock) RemoveBypassNetworkLists(ctx context.Context, req RemoveBypassNetworkListsRequest) (*RemoveBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) RemoveApiRequestConstraints(ctx context.Context, req RemoveApiRequestConstraintsRequest) (*RemoveApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveApiRequestConstraintsResponse), args.Error(1)
}

func (m *Mock) RemoveAdvancedSettingsLogging(ctx context.Context, req RemoveAdvancedSettingsLoggingRequest) (*RemoveAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *Mock) RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, req RemoveAdvancedSettingsEvasivePathMatchRequest) (*RemoveAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *Mock) RemoveActivations(ctx context.Context, req RemoveActivationsRequest) (*RemoveActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveActivationsResponse), args.Error(1)
}

func (m *Mock) GetWAPSelectedHostnames(ctx context.Context, req GetWAPSelectedHostnamesRequest) (*GetWAPSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAPSelectedHostnamesResponse), args.Error(1)
}

func (m *Mock) GetWAFProtections(ctx context.Context, req GetWAFProtectionsRequest) (*GetWAFProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAFProtectionsResponse), args.Error(1)
}

func (m *Mock) GetWAFProtection(ctx context.Context, req GetWAFProtectionRequest) (*GetWAFProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAFProtectionResponse), args.Error(1)
}

func (m *Mock) GetWAFModes(ctx context.Context, req GetWAFModesRequest) (*GetWAFModesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAFModesResponse), args.Error(1)
}

func (m *Mock) GetWAFMode(ctx context.Context, req GetWAFModeRequest) (*GetWAFModeResponse, error) {

	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAFModeResponse), args.Error(1)
}

func (m *Mock) GetVersionNotes(ctx context.Context, req GetVersionNotesRequest) (*GetVersionNotesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetVersionNotesResponse), args.Error(1)
}

func (m *Mock) GetTuningRecommendations(ctx context.Context, req GetTuningRecommendationsRequest) (*GetTuningRecommendationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetTuningRecommendationsResponse), args.Error(1)
}

func (m *Mock) GetThreatIntel(ctx context.Context, req GetThreatIntelRequest) (*GetThreatIntelResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetThreatIntelResponse), args.Error(1)
}

func (m *Mock) GetSlowPostProtections(ctx context.Context, req GetSlowPostProtectionsRequest) (*GetSlowPostProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSlowPostProtectionsResponse), args.Error(1)
}

func (m *Mock) GetSlowPostProtectionSettings(ctx context.Context, req GetSlowPostProtectionSettingsRequest) (*GetSlowPostProtectionSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSlowPostProtectionSettingsResponse), args.Error(1)
}

func (m *Mock) GetSlowPostProtectionSetting(ctx context.Context, req GetSlowPostProtectionSettingRequest) (*GetSlowPostProtectionSettingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSlowPostProtectionSettingResponse), args.Error(1)
}

func (m *Mock) GetSlowPostProtection(ctx context.Context, req GetSlowPostProtectionRequest) (*GetSlowPostProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSlowPostProtectionResponse), args.Error(1)
}

func (m *Mock) GetSiemSettings(ctx context.Context, req GetSiemSettingsRequest) (*GetSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSiemSettingsResponse), args.Error(1)
}

func (m *Mock) GetSiemDefinitions(ctx context.Context, req GetSiemDefinitionsRequest) (*GetSiemDefinitionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSiemDefinitionsResponse), args.Error(1)
}

func (m *Mock) GetSelectedHostnames(ctx context.Context, req GetSelectedHostnamesRequest) (*GetSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSelectedHostnamesResponse), args.Error(1)
}

func (m *Mock) GetSelectedHostname(ctx context.Context, req GetSelectedHostnameRequest) (*GetSelectedHostnameResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSelectedHostnameResponse), args.Error(1)
}

func (m *Mock) GetSelectableHostnames(ctx context.Context, req GetSelectableHostnamesRequest) (*GetSelectableHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSelectableHostnamesResponse), args.Error(1)
}

func (m *Mock) GetSecurityPolicyClones(ctx context.Context, req GetSecurityPolicyClonesRequest) (*GetSecurityPolicyClonesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSecurityPolicyClonesResponse), args.Error(1)
}

func (m *Mock) GetSecurityPolicyClone(ctx context.Context, req GetSecurityPolicyCloneRequest) (*GetSecurityPolicyCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSecurityPolicyCloneResponse), args.Error(1)
}

func (m *Mock) GetSecurityPolicy(ctx context.Context, req GetSecurityPolicyRequest) (*GetSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSecurityPolicyResponse), args.Error(1)
}

func (m *Mock) GetRuleUpgrade(ctx context.Context, req GetRuleUpgradeRequest) (*GetRuleUpgradeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRuleUpgradeResponse), args.Error(1)
}

func (m *Mock) GetRules(ctx context.Context, req GetRulesRequest) (*GetRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRulesResponse), args.Error(1)
}

func (m *Mock) GetRule(ctx context.Context, req GetRuleRequest) (*GetRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRuleResponse), args.Error(1)
}

func (m *Mock) GetReputationProtections(ctx context.Context, req GetReputationProtectionsRequest) (*GetReputationProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProtectionsResponse), args.Error(1)
}

func (m *Mock) GetReputationProtection(ctx context.Context, req GetReputationProtectionRequest) (*GetReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProtectionResponse), args.Error(1)
}

func (m *Mock) GetReputationProfiles(ctx context.Context, req GetReputationProfilesRequest) (*GetReputationProfilesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProfilesResponse), args.Error(1)
}

func (m *Mock) GetReputationProfileActions(ctx context.Context, req GetReputationProfileActionsRequest) (*GetReputationProfileActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProfileActionsResponse), args.Error(1)
}

func (m *Mock) GetReputationProfileAction(ctx context.Context, req GetReputationProfileActionRequest) (*GetReputationProfileActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProfileActionResponse), args.Error(1)
}

func (m *Mock) GetReputationProfile(ctx context.Context, req GetReputationProfileRequest) (*GetReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationProfileResponse), args.Error(1)
}

func (m *Mock) GetReputationAnalysis(ctx context.Context, req GetReputationAnalysisRequest) (*GetReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetReputationAnalysisResponse), args.Error(1)
}

func (m *Mock) GetRateProtections(ctx context.Context, req GetRateProtectionsRequest) (*GetRateProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRateProtectionsResponse), args.Error(1)
}

func (m *Mock) GetRateProtection(ctx context.Context, req GetRateProtectionRequest) (*GetRateProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRateProtectionResponse), args.Error(1)
}

func (m *Mock) GetRatePolicyActions(ctx context.Context, req GetRatePolicyActionsRequest) (*GetRatePolicyActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRatePolicyActionsResponse), args.Error(1)
}

func (m *Mock) GetRatePolicyAction(ctx context.Context, req GetRatePolicyActionRequest) (*GetRatePolicyActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRatePolicyActionResponse), args.Error(1)
}

func (m *Mock) GetRatePolicy(ctx context.Context, req GetRatePolicyRequest) (*GetRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRatePolicyResponse), args.Error(1)
}

func (m *Mock) GetRatePolicies(ctx context.Context, req GetRatePoliciesRequest) (*GetRatePoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetRatePoliciesResponse), args.Error(1)
}

func (m *Mock) GetSecurityPolicies(ctx context.Context, req GetSecurityPoliciesRequest) (*GetSecurityPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetSecurityPoliciesResponse), args.Error(1)
}

func (m *Mock) GetPolicyProtections(ctx context.Context, req GetPolicyProtectionsRequest) (*PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PolicyProtectionsResponse), args.Error(1)
}

func (m *Mock) GetPenaltyBoxes(ctx context.Context, req GetPenaltyBoxesRequest) (*GetPenaltyBoxesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPenaltyBoxesResponse), args.Error(1)
}

func (m *Mock) GetPenaltyBox(ctx context.Context, req GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetPenaltyBoxResponse), args.Error(1)
}

func (m *Mock) GetEvalPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPenaltyBoxResponse), args.Error(1)
}

func (m *Mock) UpdateEvalPenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdatePenaltyBoxResponse), args.Error(1)
}

func (m *Mock) GetNetworkLayerProtections(ctx context.Context, req GetNetworkLayerProtectionsRequest) (*GetNetworkLayerProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNetworkLayerProtectionsResponse), args.Error(1)
}

func (m *Mock) GetNetworkLayerProtection(ctx context.Context, req GetNetworkLayerProtectionRequest) (*GetNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetNetworkLayerProtectionResponse), args.Error(1)
}

func (m *Mock) GetMatchTargets(ctx context.Context, req GetMatchTargetsRequest) (*GetMatchTargetsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMatchTargetsResponse), args.Error(1)
}

func (m *Mock) GetMatchTargetSequence(ctx context.Context, req GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMatchTargetSequenceResponse), args.Error(1)
}

func (m *Mock) GetMatchTarget(ctx context.Context, req GetMatchTargetRequest) (*GetMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetMatchTargetResponse), args.Error(1)
}

func (m *Mock) GetIPGeoProtection(ctx context.Context, req GetIPGeoProtectionRequest) (*GetIPGeoProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetIPGeoProtectionResponse), args.Error(1)
}

func (m *Mock) GetIPGeoProtections(ctx context.Context, req GetIPGeoProtectionsRequest) (*GetIPGeoProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetIPGeoProtectionsResponse), args.Error(1)
}

func (m *Mock) GetIPGeo(ctx context.Context, req GetIPGeoRequest) (*GetIPGeoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetIPGeoResponse), args.Error(1)
}

func (m *Mock) GetFailoverHostnames(ctx context.Context, req GetFailoverHostnamesRequest) (*GetFailoverHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetFailoverHostnamesResponse), args.Error(1)
}

func (m *Mock) GetExportConfiguration(ctx context.Context, req GetExportConfigurationRequest) (*GetExportConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetExportConfigurationResponse), args.Error(1)
}

func (m *Mock) GetExportConfigurations(ctx context.Context, req GetExportConfigurationsRequest) (*GetExportConfigurationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetExportConfigurationsResponse), args.Error(1)
}

func (m *Mock) GetEvals(ctx context.Context, req GetEvalsRequest) (*GetEvalsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalsResponse), args.Error(1)
}

func (m *Mock) GetEvalRules(ctx context.Context, req GetEvalRulesRequest) (*GetEvalRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalRulesResponse), args.Error(1)
}

func (m *Mock) GetEvalRule(ctx context.Context, req GetEvalRuleRequest) (*GetEvalRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalRuleResponse), args.Error(1)
}

func (m *Mock) GetEvalProtectHosts(ctx context.Context, req GetEvalProtectHostsRequest) (*GetEvalProtectHostsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalProtectHostsResponse), args.Error(1)
}

func (m *Mock) GetEvalProtectHost(ctx context.Context, req GetEvalProtectHostRequest) (*GetEvalProtectHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalProtectHostResponse), args.Error(1)
}

func (m *Mock) GetEvalHosts(ctx context.Context, req GetEvalHostsRequest) (*GetEvalHostsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalHostsResponse), args.Error(1)
}

func (m *Mock) GetEvalHost(ctx context.Context, req GetEvalHostRequest) (*GetEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalHostResponse), args.Error(1)
}

func (m *Mock) GetEvalGroups(ctx context.Context, req GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAttackGroupsResponse), args.Error(1)
}

func (m *Mock) GetEvalGroup(ctx context.Context, req GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAttackGroupResponse), args.Error(1)
}

func (m *Mock) GetEval(ctx context.Context, req GetEvalRequest) (*GetEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetEvalResponse), args.Error(1)
}

func (m *Mock) GetCustomRules(ctx context.Context, req GetCustomRulesRequest) (*GetCustomRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomRulesResponse), args.Error(1)
}

func (m *Mock) GetCustomRuleActions(ctx context.Context, req GetCustomRuleActionsRequest) (*GetCustomRuleActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomRuleActionsResponse), args.Error(1)
}

func (m *Mock) GetCustomRuleAction(ctx context.Context, req GetCustomRuleActionRequest) (*GetCustomRuleActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomRuleActionResponse), args.Error(1)
}

func (m *Mock) GetCustomRule(ctx context.Context, req GetCustomRuleRequest) (*GetCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomRuleResponse), args.Error(1)
}

func (m *Mock) GetCustomDenyList(ctx context.Context, req GetCustomDenyListRequest) (*GetCustomDenyListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomDenyListResponse), args.Error(1)
}

func (m *Mock) GetCustomDeny(ctx context.Context, req GetCustomDenyRequest) (*GetCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetCustomDenyResponse), args.Error(1)
}

func (m *Mock) GetContractsGroups(ctx context.Context, req GetContractsGroupsRequest) (*GetContractsGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetContractsGroupsResponse), args.Error(1)
}

func (m *Mock) GetConfigurations(ctx context.Context, req GetConfigurationsRequest) (*GetConfigurationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConfigurationsResponse), args.Error(1)
}

func (m *Mock) GetConfigurationVersions(ctx context.Context, req GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConfigurationVersionsResponse), args.Error(1)
}

func (m *Mock) GetConfigurationVersionClone(ctx context.Context, req GetConfigurationVersionCloneRequest) (*GetConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConfigurationVersionCloneResponse), args.Error(1)
}

func (m *Mock) GetConfigurationClone(ctx context.Context, req GetConfigurationCloneRequest) (*GetConfigurationCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConfigurationCloneResponse), args.Error(1)
}

func (m *Mock) GetBypassNetworkLists(ctx context.Context, req GetBypassNetworkListsRequest) (*GetBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) GetAttackGroups(ctx context.Context, req GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAttackGroupsResponse), args.Error(1)
}

func (m *Mock) GetAttackGroupRecommendations(ctx context.Context, req GetAttackGroupRecommendationsRequest) (*GetAttackGroupRecommendationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAttackGroupRecommendationsResponse), args.Error(1)
}

func (m *Mock) GetAttackGroup(ctx context.Context, req GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAttackGroupResponse), args.Error(1)
}

func (m *Mock) GetApiRequestConstraints(ctx context.Context, req GetApiRequestConstraintsRequest) (*GetApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetApiRequestConstraintsResponse), args.Error(1)
}

func (m *Mock) GetApiHostnameCoverageOverlapping(ctx context.Context, req GetApiHostnameCoverageOverlappingRequest) (*GetApiHostnameCoverageOverlappingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetApiHostnameCoverageOverlappingResponse), args.Error(1)
}

func (m *Mock) GetApiHostnameCoverageMatchTargets(ctx context.Context, req GetApiHostnameCoverageMatchTargetsRequest) (*GetApiHostnameCoverageMatchTargetsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetApiHostnameCoverageMatchTargetsResponse), args.Error(1)
}

func (m *Mock) GetApiHostnameCoverage(ctx context.Context, req GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetApiHostnameCoverageResponse), args.Error(1)
}

func (m *Mock) GetApiEndpoints(ctx context.Context, req GetApiEndpointsRequest) (*GetApiEndpointsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetApiEndpointsResponse), args.Error(1)
}

func (m *Mock) GetAdvancedSettingsPragma(ctx context.Context, req GetAdvancedSettingsPragmaRequest) (*GetAdvancedSettingsPragmaResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAdvancedSettingsPragmaResponse), args.Error(1)
}

func (m *Mock) GetAdvancedSettingsPrefetch(ctx context.Context, req GetAdvancedSettingsPrefetchRequest) (*GetAdvancedSettingsPrefetchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAdvancedSettingsPrefetchResponse), args.Error(1)
}

func (m *Mock) GetAdvancedSettingsLogging(ctx context.Context, req GetAdvancedSettingsLoggingRequest) (*GetAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *Mock) GetAdvancedSettingsEvasivePathMatch(ctx context.Context, req GetAdvancedSettingsEvasivePathMatchRequest) (*GetAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *Mock) GetActivations(ctx context.Context, req GetActivationsRequest) (*GetActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetActivationsResponse), args.Error(1)
}

func (m *Mock) GetAPIConstraintsProtection(ctx context.Context, req GetAPIConstraintsProtectionRequest) (*GetAPIConstraintsProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetAPIConstraintsProtectionResponse), args.Error(1)
}

func (m *Mock) CreateSecurityPolicyClone(ctx context.Context, req CreateSecurityPolicyCloneRequest) (*CreateSecurityPolicyCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateSecurityPolicyCloneResponse), args.Error(1)
}

func (m *Mock) CreateSecurityPolicy(ctx context.Context, req CreateSecurityPolicyRequest) (*CreateSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateSecurityPolicyResponse), args.Error(1)
}

func (m *Mock) CreateReputationProfile(ctx context.Context, req CreateReputationProfileRequest) (*CreateReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateReputationProfileResponse), args.Error(1)
}

func (m *Mock) CreateRatePolicy(ctx context.Context, req CreateRatePolicyRequest) (*CreateRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateRatePolicyResponse), args.Error(1)
}

func (m *Mock) CreateMatchTarget(ctx context.Context, req CreateMatchTargetRequest) (*CreateMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateMatchTargetResponse), args.Error(1)
}

func (m *Mock) CreateCustomRule(ctx context.Context, req CreateCustomRuleRequest) (*CreateCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateCustomRuleResponse), args.Error(1)
}

func (m *Mock) CreateCustomDeny(ctx context.Context, req CreateCustomDenyRequest) (*CreateCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateCustomDenyResponse), args.Error(1)
}

func (m *Mock) CreateConfiguration(ctx context.Context, req CreateConfigurationRequest) (*CreateConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateConfigurationResponse), args.Error(1)
}

func (m *Mock) CreateConfigurationVersionClone(ctx context.Context, req CreateConfigurationVersionCloneRequest) (*CreateConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateConfigurationVersionCloneResponse), args.Error(1)
}

func (m *Mock) CreateConfigurationClone(ctx context.Context, req CreateConfigurationCloneRequest) (*CreateConfigurationCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateConfigurationCloneResponse), args.Error(1)
}

func (m *Mock) CreateActivations(ctx context.Context, req CreateActivationsRequest, _ bool) (*CreateActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateActivationsResponse), args.Error(1)
}

func (m *Mock) GetConfiguration(ctx context.Context, req GetConfigurationRequest) (*GetConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetConfigurationResponse), args.Error(1)
}

func (m *Mock) GetWAPBypassNetworkLists(ctx context.Context, req GetWAPBypassNetworkListsRequest) (*GetWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) RemoveWAPBypassNetworkLists(ctx context.Context, req RemoveWAPBypassNetworkListsRequest) (*RemoveWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RemoveWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) UpdateWAPBypassNetworkLists(ctx context.Context, req UpdateWAPBypassNetworkListsRequest) (*UpdateWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpdateWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *Mock) GetActivationHistory(ctx context.Context, req GetActivationHistoryRequest) (*GetActivationHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetActivationHistoryResponse), args.Error(1)
}

func (m *Mock) GetMalwareProtection(ctx context.Context, params GetMalwareProtectionRequest) (*GetMalwareProtectionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetMalwareProtectionResponse), args.Error(1)
}

func (m *Mock) GetMalwareProtections(ctx context.Context, params GetMalwareProtectionsRequest) (*GetMalwareProtectionsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetMalwareProtectionsResponse), args.Error(1)
}

func (m *Mock) UpdateMalwareProtection(ctx context.Context, params UpdateMalwareProtectionRequest) (*UpdateMalwareProtectionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateMalwareProtectionResponse), args.Error(1)
}

func (m *Mock) GetMalwareContentTypes(ctx context.Context, params GetMalwareContentTypesRequest) (*GetMalwareContentTypesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetMalwareContentTypesResponse), args.Error(1)
}

func (m *Mock) CreateMalwarePolicy(ctx context.Context, params CreateMalwarePolicyRequest) (*MalwarePolicyResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*MalwarePolicyResponse), args.Error(1)
}

func (m *Mock) GetMalwarePolicy(ctx context.Context, params GetMalwarePolicyRequest) (*MalwarePolicyResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*MalwarePolicyResponse), args.Error(1)
}

func (m *Mock) GetMalwarePolicies(ctx context.Context, params GetMalwarePoliciesRequest) (*MalwarePoliciesResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*MalwarePoliciesResponse), args.Error(1)
}

func (m *Mock) UpdateMalwarePolicy(ctx context.Context, params UpdateMalwarePolicyRequest) (*MalwarePolicyResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*MalwarePolicyResponse), args.Error(1)
}

func (m *Mock) RemoveMalwarePolicy(ctx context.Context, params RemoveMalwarePolicyRequest) error {
	args := m.Called(ctx, params)

	return args.Error(0)
}

func (m *Mock) GetMalwarePolicyActions(ctx context.Context, params GetMalwarePolicyActionsRequest) (*GetMalwarePolicyActionsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetMalwarePolicyActionsResponse), args.Error(1)
}

func (m *Mock) UpdateMalwarePolicyAction(ctx context.Context, params UpdateMalwarePolicyActionRequest) (*UpdateMalwarePolicyActionResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateMalwarePolicyActionResponse), args.Error(1)
}

func (m *Mock) UpdateMalwarePolicyActions(ctx context.Context, params UpdateMalwarePolicyActionsRequest) (*UpdateMalwarePolicyActionsResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateMalwarePolicyActionsResponse), args.Error(1)
}
