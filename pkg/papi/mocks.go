//revive:disable:exported

package papi

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

var _ PAPI = &Mock{}

type (
	// GetGroupsFn is any function having the same signature as GetGroups
	GetGroupsFn func(context.Context) (*GetGroupsResponse, error)

	// GetCPCodesFn is any function having the same signature as GetCPCodes
	GetCPCodesFn func(context.Context, GetCPCodesRequest) (*GetCPCodesResponse, error)

	// CreateCPCodeFn is any function having the same signature as CreateCPCode
	CreateCPCodeFn func(context.Context, CreateCPCodeRequest) (*CreateCPCodeResponse, error)

	// UpdateCPCodeFn is any function having the same signature as UpdateCPCode
	UpdateCPCodeFn func(context.Context, UpdateCPCodeRequest) (*CPCodeDetailResponse, error)

	// GetPropertyFunc is any function having the same signature as GetProperty
	GetPropertyFunc func(context.Context, GetPropertyRequest) (*GetPropertyResponse, error)

	// GetPropertyVersionsFn is any function having the same signature as GetPropertyVersions
	GetPropertyVersionsFn func(context.Context, GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error)

	// GetPropertyVersionHostnamesFn is any function having the same signature as GetPropertyVersionHostnames
	GetPropertyVersionHostnamesFn func(context.Context, GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error)

	// GetRuleTreeFn is any function having the same signature as GetRuleTree
	GetRuleTreeFn = func(context.Context, GetRuleTreeRequest) (*GetRuleTreeResponse, error)

	// UpdateRuleTreeFn is any function having the same signature as UpdateRuleTree
	UpdateRuleTreeFn func(context.Context, UpdateRulesRequest) (*UpdateRulesResponse, error)
)

func (p *Mock) GetGroups(ctx context.Context) (*GetGroupsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetGroupsResponse), args.Error(1)
}

func (p *Mock) GetContracts(ctx context.Context) (*GetContractsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetContractsResponse), args.Error(1)
}

func (p *Mock) CreateActivation(ctx context.Context, r CreateActivationRequest) (*CreateActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateActivationResponse), args.Error(1)
}

func (p *Mock) GetActivations(ctx context.Context, r GetActivationsRequest) (*GetActivationsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationsResponse), args.Error(1)
}

func (p *Mock) GetActivation(ctx context.Context, r GetActivationRequest) (*GetActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetActivationResponse), args.Error(1)
}

func (p *Mock) CancelActivation(ctx context.Context, r CancelActivationRequest) (*CancelActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CancelActivationResponse), args.Error(1)
}

func (p *Mock) GetCPCodes(ctx context.Context, r GetCPCodesRequest) (*GetCPCodesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCPCodesResponse), args.Error(1)
}

func (p *Mock) GetCPCode(ctx context.Context, r GetCPCodeRequest) (*GetCPCodesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetCPCodesResponse), args.Error(1)
}

func (p *Mock) CreateCPCode(ctx context.Context, r CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateCPCodeResponse), args.Error(1)
}

func (p *Mock) UpdateCPCode(ctx context.Context, r UpdateCPCodeRequest) (*CPCodeDetailResponse, error) {
	args := p.Called(ctx, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CPCodeDetailResponse), args.Error(1)
}

func (p *Mock) GetCPCodeDetail(ctx context.Context, r int) (*CPCodeDetailResponse, error) {
	args := p.Called(ctx, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CPCodeDetailResponse), args.Error(1)
}

func (p *Mock) GetProperties(ctx context.Context, r GetPropertiesRequest) (*GetPropertiesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertiesResponse), args.Error(1)
}

func (p *Mock) CreateProperty(ctx context.Context, r CreatePropertyRequest) (*CreatePropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreatePropertyResponse), args.Error(1)
}

func (p *Mock) GetProperty(ctx context.Context, r GetPropertyRequest) (*GetPropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyResponse), args.Error(1)
}

func (p *Mock) RemoveProperty(ctx context.Context, r RemovePropertyRequest) (*RemovePropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RemovePropertyResponse), args.Error(1)
}

func (p *Mock) GetPropertyVersions(ctx context.Context, r GetPropertyVersionsRequest) (*GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyVersionsResponse), args.Error(1)
}

func (p *Mock) GetPropertyVersion(ctx context.Context, r GetPropertyVersionRequest) (*GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyVersionsResponse), args.Error(1)
}

func (p *Mock) CreatePropertyVersion(ctx context.Context, r CreatePropertyVersionRequest) (*CreatePropertyVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreatePropertyVersionResponse), args.Error(1)
}

func (p *Mock) GetLatestVersion(ctx context.Context, r GetLatestVersionRequest) (*GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyVersionsResponse), args.Error(1)
}

func (p *Mock) GetAvailableBehaviors(ctx context.Context, r GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetFeaturesCriteriaResponse), args.Error(1)
}

func (p *Mock) GetAvailableCriteria(ctx context.Context, r GetFeaturesRequest) (*GetFeaturesCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetFeaturesCriteriaResponse), args.Error(1)
}

func (p *Mock) GetEdgeHostnames(ctx context.Context, r GetEdgeHostnamesRequest) (*GetEdgeHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetEdgeHostnamesResponse), args.Error(1)
}

func (p *Mock) GetEdgeHostname(ctx context.Context, r GetEdgeHostnameRequest) (*GetEdgeHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetEdgeHostnamesResponse), args.Error(1)
}

func (p *Mock) CreateEdgeHostname(ctx context.Context, r CreateEdgeHostnameRequest) (*CreateEdgeHostnameResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateEdgeHostnameResponse), args.Error(1)
}

func (p *Mock) GetProducts(ctx context.Context, r GetProductsRequest) (*GetProductsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetProductsResponse), args.Error(1)
}

func (p *Mock) SearchProperties(ctx context.Context, r SearchRequest) (*SearchResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*SearchResponse), args.Error(1)
}

func (p *Mock) GetPropertyVersionHostnames(ctx context.Context, r GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetPropertyVersionHostnamesResponse), args.Error(1)
}

func (p *Mock) UpdatePropertyVersionHostnames(ctx context.Context, r UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error) {
	args := p.Called(ctx, r)

	return args.Get(0).(*UpdatePropertyVersionHostnamesResponse), args.Error(1)
}

func (p *Mock) GetClientSettings(ctx context.Context) (*ClientSettingsBody, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ClientSettingsBody), args.Error(1)
}

func (p *Mock) UpdateClientSettings(ctx context.Context, r ClientSettingsBody) (*ClientSettingsBody, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ClientSettingsBody), args.Error(1)
}

func (p *Mock) GetRuleTree(ctx context.Context, r GetRuleTreeRequest) (*GetRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRuleTreeResponse), args.Error(1)
}

func (p *Mock) UpdateRuleTree(ctx context.Context, r UpdateRulesRequest) (*UpdateRulesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateRulesResponse), args.Error(1)
}

func (p *Mock) GetRuleFormats(ctx context.Context) (*GetRuleFormatsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetRuleFormatsResponse), args.Error(1)
}

func (p *Mock) OnGetGroups(ctx interface{}, impl GetGroupsFn) *mock.Call {
	call := p.On("GetGroups", ctx)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)

		call.Return(impl(callCtx))
	})

	return call
}

func (p *Mock) OnGetCPCodes(impl GetCPCodesFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(callArgs mock.Arguments) {
		ctx := callArgs.Get(0).(context.Context)
		req := callArgs.Get(1).(GetCPCodesRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{mock.Anything, mock.Anything}
	}

	call = p.On("GetCPCodes", args...).Run(runFn)
	return call
}

func (p *Mock) OnCreateCPCode(impl CreateCPCodeFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		req := args.Get(1).(CreateCPCodeRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{mock.Anything, mock.Anything}
	}

	call = p.On("CreateCPCode", args...).Run(runFn)
	return call
}

func (p *Mock) OnUpdateCPCode(impl UpdateCPCodeFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(callArgs mock.Arguments) {
		ctx := callArgs.Get(0).(context.Context)
		req := callArgs.Get(1).(UpdateCPCodeRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{mock.Anything, mock.Anything}
	}

	call = p.On("UpdateCPCode", args...).Run(runFn)
	return call
}

func (p *Mock) OnGetProperty(ctx, req interface{}, impl GetPropertyFunc) *mock.Call {
	call := p.On("GetProperty", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(GetPropertyRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *Mock) OnGetPropertyVersions(ctx, req interface{}, impl GetPropertyVersionsFn) *mock.Call {
	call := p.On("GetPropertyVersions", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(GetPropertyVersionsRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *Mock) OnGetPropertyVersionHostnames(ctx, req interface{}, impl GetPropertyVersionHostnamesFn) *mock.Call {
	call := p.On("GetPropertyVersionHostnames", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(GetPropertyVersionHostnamesRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *Mock) OnGetRuleTree(ctx, req interface{}, impl GetRuleTreeFn) *mock.Call {
	call := p.On("GetRuleTree", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(GetRuleTreeRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *Mock) OnUpdateRuleTree(ctx, req interface{}, impl UpdateRuleTreeFn) *mock.Call {
	call := p.On("UpdateRuleTree", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(UpdateRulesRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *Mock) ListIncludes(ctx context.Context, r ListIncludesRequest) (*ListIncludesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListIncludesResponse), args.Error(1)
}

func (p *Mock) ListIncludeParents(ctx context.Context, r ListIncludeParentsRequest) (*ListIncludeParentsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListIncludeParentsResponse), args.Error(1)
}

func (p *Mock) GetInclude(ctx context.Context, r GetIncludeRequest) (*GetIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetIncludeResponse), args.Error(1)
}

func (p *Mock) CreateInclude(ctx context.Context, r CreateIncludeRequest) (*CreateIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateIncludeResponse), args.Error(1)
}

func (p *Mock) DeleteInclude(ctx context.Context, r DeleteIncludeRequest) (*DeleteIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteIncludeResponse), args.Error(1)
}

func (p *Mock) GetIncludeRuleTree(ctx context.Context, r GetIncludeRuleTreeRequest) (*GetIncludeRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetIncludeRuleTreeResponse), args.Error(1)
}

func (p *Mock) UpdateIncludeRuleTree(ctx context.Context, r UpdateIncludeRuleTreeRequest) (*UpdateIncludeRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateIncludeRuleTreeResponse), args.Error(1)
}

func (p *Mock) ActivateInclude(ctx context.Context, r ActivateIncludeRequest) (*ActivationIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ActivationIncludeResponse), args.Error(1)
}

func (p *Mock) DeactivateInclude(ctx context.Context, r DeactivateIncludeRequest) (*DeactivationIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeactivationIncludeResponse), args.Error(1)
}

func (p *Mock) GetIncludeActivation(ctx context.Context, r GetIncludeActivationRequest) (*GetIncludeActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetIncludeActivationResponse), args.Error(1)
}

func (p *Mock) ListIncludeActivations(ctx context.Context, r ListIncludeActivationsRequest) (*ListIncludeActivationsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListIncludeActivationsResponse), args.Error(1)
}

func (p *Mock) CreateIncludeVersion(ctx context.Context, r CreateIncludeVersionRequest) (*CreateIncludeVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateIncludeVersionResponse), args.Error(1)
}

func (p *Mock) GetIncludeVersion(ctx context.Context, r GetIncludeVersionRequest) (*GetIncludeVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetIncludeVersionResponse), args.Error(1)
}

func (p *Mock) ListIncludeVersions(ctx context.Context, r ListIncludeVersionsRequest) (*ListIncludeVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListIncludeVersionsResponse), args.Error(1)
}

func (p *Mock) ListIncludeVersionAvailableCriteria(ctx context.Context, r ListAvailableCriteriaRequest) (*AvailableCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AvailableCriteriaResponse), args.Error(1)
}

func (p *Mock) ListIncludeVersionAvailableBehaviors(ctx context.Context, r ListAvailableBehaviorsRequest) (*AvailableBehaviorsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AvailableBehaviorsResponse), args.Error(1)
}

func (p *Mock) ListAvailableIncludes(ctx context.Context, r ListAvailableIncludesRequest) (*ListAvailableIncludesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListAvailableIncludesResponse), args.Error(1)
}

func (p *Mock) ListReferencedIncludes(ctx context.Context, r ListReferencedIncludesRequest) (*ListReferencedIncludesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListReferencedIncludesResponse), args.Error(1)
}

func (p *Mock) CancelIncludeActivation(ctx context.Context, r CancelIncludeActivationRequest) (*CancelIncludeActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CancelIncludeActivationResponse), args.Error(1)
}
