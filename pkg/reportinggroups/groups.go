package reportinggroups

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
)

func (r *reportinggroups) CreateReportingGroup(ctx context.Context, params CreateReportingGroupRequest) (*CreateReportingGroupResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("CreateReportingGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrCreateReportingGroup, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/cprg/v1/reporting-groups").Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrCreateReportingGroup, err)
	}

	var result CreateReportingGroupResponse
	resp, err := r.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrCreateReportingGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %w", ErrCreateReportingGroup, r.Error(resp))
	}

	result.ResourceLimits = extractReportingGroupLimitHeaders(resp, logger)

	return &result, nil
}

func (r *reportinggroups) GetReportingGroup(ctx context.Context, params GetReportingGroupsRequest) (*GetReportingGroupResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("GetReportingGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetReportingGroup, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cprg/v1/reporting-groups/%d", params.ReportingGroupID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetReportingGroup, err)
	}

	var result GetReportingGroupResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetReportingGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetReportingGroup, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) ListReportingGroups(ctx context.Context, params ListReportingGroupsRequest) (*ListReportingGroupsResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("ListReportingGroups")

	req, err := request.NewGet(ctx, "/cprg/v1/reporting-groups").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", strconv.FormatInt(params.GroupID, 10), params.GroupID != 0).
		AddQueryParamIf("reportingGroupName", params.ReportingGroupName, params.ReportingGroupName != "").
		AddQueryParamIf("cpcodeId", params.CpCodeID, params.CpCodeID != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListReportingGroups, err)
	}

	var result ListReportingGroupsResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListReportingGroups, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListReportingGroups, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) UpdateReportingGroup(ctx context.Context, params UpdateReportingGroupRequest) (*UpdateReportingGroupResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("UpdateReportingGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUpdateReportingGroup, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cprg/v1/reporting-groups/%d", params.ReportingGroupID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrUpdateReportingGroup, err)
	}

	var result UpdateReportingGroupResponse
	resp, err := r.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrUpdateReportingGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrUpdateReportingGroup, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) DeleteReportingGroup(ctx context.Context, params DeleteReportingGroupRequest) error {
	logger := r.Log(ctx)
	logger.Debug("DeleteReportingGroup")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDeleteReportingGroup, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/cprg/v1/reporting-groups/%d", params.ReportingGroupID).Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %w", ErrDeleteReportingGroup, err)
	}

	resp, err := r.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %w", ErrDeleteReportingGroup, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %w", ErrDeleteReportingGroup, r.Error(resp))
	}

	return nil
}

func (r *reportinggroups) ListProducts(ctx context.Context, params ListProductsRequest) (*ListProductsResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("ListProducts")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListProducts, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cprg/v1/reporting-groups/%d/products", params.ReportingGroupID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListProducts, err)
	}

	var result ListProductsResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListProducts, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListProducts, r.Error(resp))
	}

	return &result, nil
}

func extractReportingGroupLimitHeaders(resp *http.Response, logger log.Interface) ResourceLimitsMetadata {
	var reportingGroupsLimits ResourceLimitsMetadata

	limitTotal, err := strconv.ParseInt(resp.Header.Get("X-Limit-Max-Reporting-Groups-Limit"), 10, 64)
	if err == nil {
		reportingGroupsLimits.ReportingGroupsLimitTotal = &limitTotal
	} else {
		logger.Warnf("Failed to parse X-Limit-Max-Reporting-Groups-Limit header: %v", err)
	}

	limitRemaining, err := strconv.ParseInt(resp.Header.Get("X-Limit-Max-Reporting-Groups-Remaining"), 10, 64)
	if err == nil {
		reportingGroupsLimits.ReportingGroupsLimitRemaining = &limitRemaining
	} else {
		logger.Warnf("Failed to parse X-Limit-Max-Reporting-Groups-Remaining header: %v", err)
	}

	return reportingGroupsLimits
}

func (r *reportinggroups) GetReportingGroupsWaterMarkLimits(ctx context.Context, params GetReportingGroupsWaterMarkLimitsRequest) (*GetReportingGroupsWaterMarkLimitsResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("GetReportingGroupsWaterMarkLimits")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetReportingGroupsWaterMarkLimits, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cprg/v1/reporting-groups/contracts/%s/watermark-limits", params.ContractID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetReportingGroupsWaterMarkLimits, err)
	}

	var result GetReportingGroupsWaterMarkLimitsResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetReportingGroupsWaterMarkLimits, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetReportingGroupsWaterMarkLimits, r.Error(resp))
	}

	return &result, nil
}
