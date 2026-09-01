package reportinggroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (r *reportinggroups) ListCPCodes(ctx context.Context, params ListCPCodesRequest) (*ListCPCodesResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("ListCPCodes")

	req, err := request.NewGet(ctx, "/cprg/v1/cpcodes").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("groupId", params.GroupID, params.GroupID != "").
		AddQueryParamIf("productId", params.ProductID, params.ProductID != "").
		AddQueryParamIf("cpcodeName", params.CPCodeName, params.CPCodeName != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListCPCodes, err)
	}

	var result ListCPCodesResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListCPCodes, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListCPCodes, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) GetCPCodesWaterMarkLimits(ctx context.Context, params GetCPCodesWaterMarkLimitsRequest) (*GetCPCodesWaterMarkLimitsResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("GetCPCodesWaterMarkLimits")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetCPCodesWaterMarkLimits, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cprg/v1/cpcodes/contracts/%s/watermark-limits", params.ContractID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetCPCodesWaterMarkLimits, err)
	}

	var result GetCPCodesWaterMarkLimitsResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetCPCodesWaterMarkLimits, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetCPCodesWaterMarkLimits, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) GetCPCode(ctx context.Context, params GetCPCodeRequest) (*GetCPCodeResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("GetCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetCPCode, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cprg/v1/cpcodes/%d", params.CPCodeID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetCPCode, err)
	}

	var result GetCPCodeResponse
	resp, err := r.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetCPCode, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetCPCode, r.Error(resp))
	}

	return &result, nil
}

func (r *reportinggroups) UpdateCPCode(ctx context.Context, params UpdateCPCodeRequest) (*UpdateCPCodeResponse, error) {
	logger := r.Log(ctx)
	logger.Debug("UpdateCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUpdateCPCode, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cprg/v1/cpcodes/%d", params.CPCodeID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrUpdateCPCode, err)
	}

	var result UpdateCPCodeResponse
	resp, err := r.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrUpdateCPCode, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrUpdateCPCode, r.Error(resp))
	}

	return &result, nil
}
