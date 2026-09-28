package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/texts"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

func (c *cloudcertificates) CreateLineage(ctx context.Context, params CreateLineageRequest) (*CreateLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrCreateLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages").
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrCreateLineage, err)
	}

	var result CreateLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrCreateLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %w", ErrCreateLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) GetLineage(ctx context.Context, params GetLineageRequest) (*GetLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetLineage, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d", params.LineageID).
		AddQueryParamsFunc("expandGenerations", func() []string {
			return texts.ToStrings(params.ExpandGenerations)
		}, len(params.ExpandGenerations) > 0).
		UseCommaSeparatedQuery().
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetLineage, err)
	}

	var result GetLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrGetLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) ListLineages(ctx context.Context, params ListLineagesRequest) (*ListLineagesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLineages")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListLineages, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages").
		AddQueryParamIf("contractId", params.ContractID, params.ContractID != "").
		AddQueryParamIf("lineageName", params.LineageName, params.LineageName != "").
		AddQueryParamsFunc("lineageId", func() []string {
			lineageIDs := make([]string, len(params.LineageIDs))
			for i, lineageID := range params.LineageIDs {
				lineageIDs[i] = strconv.FormatInt(lineageID, 10)
			}
			return lineageIDs
		}, len(params.LineageIDs) > 0).
		AddQueryParamIf("secureNetwork", string(params.SecureNetwork), params.SecureNetwork != "").
		AddQueryParamIf("stackMode", string(params.StackMode), params.StackMode != "").
		AddQueryParamIf("lineageType", string(params.LineageType), params.LineageType != "").
		AddQueryParamIf("domain", params.Domain, params.Domain != "").
		AddQueryParamsFunc("generationStatus", func() []string {
			return texts.ToStrings(params.GenerationStatus)
		}, len(params.GenerationStatus) > 0).
		AddQueryParamFunc("expiringInDays", func() string {
			return strconv.Itoa(*params.ExpiringInDays)
		}, params.ExpiringInDays != nil).
		AddQueryParamIf("keyType", string(params.KeyType), params.KeyType != "").
		AddQueryParamIf("issuer", params.Issuer, params.Issuer != "").
		AddQueryParamsFunc("expandGenerations", func() []string {
			return texts.ToStrings(params.ExpandGenerations)
		}, len(params.ExpandGenerations) > 0).
		AddQueryParamIf("pageSize", strconv.Itoa(params.PageSize), params.PageSize > 0).
		AddQueryParamIf("after", params.After, params.After != "").
		AddQueryParamIf("sort", params.Sort, params.Sort != "").
		UseCommaSeparatedQuery().
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListLineages, err)
	}

	var result ListLineagesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListLineages, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListLineages, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) RenameLineage(ctx context.Context, params RenameLineageRequest) (*RenameLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("RenameLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrRenameLineage, ErrStructValidation, err)
	}

	body := []lineagePatchOperation{
		{Op: "replace", Path: "/lineageName", Value: params.LineageName},
	}

	req, err := request.NewPatch(ctx, "/ccm/v2/lineages/%d", params.LineageID).
		WithBody(body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrRenameLineage, err)
	}

	var result RenameLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrRenameLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrRenameLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) DeleteLineage(ctx context.Context, params DeleteLineageRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteLineage")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDeleteLineage, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/ccm/v2/lineages/%d", params.LineageID).
		Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %w", ErrDeleteLineage, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request execution failed: %w", ErrDeleteLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %w", ErrDeleteLineage, c.Error(resp))
	}

	return nil
}

func (c *cloudcertificates) ListLineageActivity(ctx context.Context, params ListLineageActivityRequest) (*ListLineageActivityResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLineageActivity")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListLineageActivity, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/activity", params.LineageID).
		AddQueryParamIf("pageSize", strconv.Itoa(params.PageSize), params.PageSize > 0).
		AddQueryParamIf("after", params.Cursor, params.Cursor != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListLineageActivity, err)
	}

	var result ListLineageActivityResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListLineageActivity, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListLineageActivity, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) PromoteLineage(ctx context.Context, params PromoteLineageRequest) (*PromoteLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("PromoteLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrPromoteLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages/%d/activations", params.LineageID).
		WithBody(activationWireRequestBody{
			OperationType: OperationTypePromote,
			GenerationID:  params.GenerationID,
			Networks:      params.Networks,
		}).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrPromoteLineage, err)
	}

	var result PromoteLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrPromoteLineage, err)
	}
	defer session.CloseResponseBody(resp)

	// The generation is already active on every requested network: the API returns 204 No Content (no body,
	// nothing was activated) instead of 200 with an activation list.
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %w", ErrPromoteLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) RollbackLineage(ctx context.Context, params RollbackLineageRequest) (*RollbackLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("RollbackLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrRollbackLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages/%d/activations", params.LineageID).
		WithBody(activationWireRequestBody{
			OperationType: OperationTypeRollback,
			GenerationID:  params.GenerationID,
		}).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrRollbackLineage, err)
	}

	var result RollbackLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrRollbackLineage, err)
	}
	defer session.CloseResponseBody(resp)

	// The previous production generation is already active: the API returns 204 No Content (no body, nothing
	// was activated) instead of 200 with an activation page.
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %w", ErrRollbackLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) ReplaceStagingLineage(ctx context.Context, params ReplaceStagingLineageRequest) (*ReplaceStagingLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ReplaceStagingLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrReplaceStagingLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages/%d/activations", params.LineageID).
		WithBody(activationWireRequestBody{
			OperationType: OperationTypeReplaceStaging,
			GenerationID:  params.GenerationID,
		}).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrReplaceStagingLineage, err)
	}

	var result ReplaceStagingLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrReplaceStagingLineage, err)
	}
	defer session.CloseResponseBody(resp)

	// GenerationID is already the current staging generation: the API returns 204 No Content (no body, nothing
	// was activated) instead of 200 with an activation page.
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %w", ErrReplaceStagingLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) GetActivationStatus(ctx context.Context, params GetActivationStatusRequest) (*GetActivationStatusResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetActivationStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetActivationStatus, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/activations/%d", params.LineageID, params.ActivationID).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetActivationStatus, err)
	}

	var result GetActivationStatusResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrGetActivationStatus, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetActivationStatus, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) ListActivations(ctx context.Context, params ListActivationsRequest) (*ListActivationsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListActivations, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/activations", params.LineageID).
		AddQueryParamIf("pageSize", strconv.Itoa(params.PageSize), params.PageSize > 0).
		AddQueryParamIf("cursor", params.Cursor, params.Cursor != "").
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListActivations, err)
	}

	var result ListActivationsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListActivations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListActivations, c.Error(resp))
	}

	return &result, nil
}
