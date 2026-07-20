package cloudcertificates

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
)

// The API's UploadSignedCertificate endpoint declares a fixed `produces` media type on its server-side mapping,
// rejecting the session's default `Accept: application/json` with a 406 media-type-not-acceptable error. This
// vendor media type must be set explicitly as the Accept header.
const mediaTypeCertificateLineageUploadV3 = "application/prs.akamai.cps-ccm-api.third-party.certificate-lineage-create.v3+json; charset=UTF-8"

func (c *cloudcertificates) RenewLineage(ctx context.Context, params RenewLineageRequest) (*RenewLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("RenewLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrRenewLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages/%d/generations", params.LineageID).
		AddQueryParamIf("confirmAbandonHead", strconv.FormatBool(params.ConfirmAbandonHead), params.ConfirmAbandonHead).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrRenewLineage, err)
	}

	var result RenewLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrRenewLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %w", ErrRenewLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) CompleteLineage(ctx context.Context, params CompleteLineageRequest) (*CompleteLineageResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("CompleteLineage")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrCompleteLineage, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/ccm/v2/lineages/%d/generations/complete", params.LineageID).
		AddQueryParamIf("acknowledgeWarnings", strconv.FormatBool(params.AcknowledgeWarnings), params.AcknowledgeWarnings).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrCompleteLineage, err)
	}

	var result CompleteLineageResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrCompleteLineage, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %w", ErrCompleteLineage, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) GetGeneration(ctx context.Context, params GetGenerationRequest) (*GetGenerationResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("GetGeneration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrGetGeneration, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/generations/%d", params.LineageID, params.GenerationID).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetGeneration, err)
	}

	var result GetGenerationResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrGetGeneration, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetGeneration, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) DeleteGeneration(ctx context.Context, params DeleteGenerationRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteGeneration")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w: %w", ErrDeleteGeneration, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/ccm/v2/lineages/%d/generations/%d", params.LineageID, params.GenerationID).
		AddQueryParamIf("acknowledgeRollbackCandidateRemoval", strconv.FormatBool(params.AcknowledgeRollbackCandidateRemoval), params.AcknowledgeRollbackCandidateRemoval).
		Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %w", ErrDeleteGeneration, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request execution failed: %w", ErrDeleteGeneration, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %w", ErrDeleteGeneration, c.Error(resp))
	}

	return nil
}

func (c *cloudcertificates) ListArchivedGenerations(ctx context.Context, params ListArchivedGenerationsRequest) (*ListArchivedGenerationsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListArchivedGenerations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrListArchivedGenerations, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/ccm/v2/lineages/%d/history", params.LineageID).
		AddQueryParamIf("includeAlgorithms", strconv.FormatBool(params.IncludeAlgorithms), params.IncludeAlgorithms).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListArchivedGenerations, err)
	}

	var result ListArchivedGenerationsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrListArchivedGenerations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListArchivedGenerations, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudcertificates) UploadSignedCertificate(ctx context.Context, params UploadSignedCertificateRequest) (*UploadSignedCertificateResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("UploadSignedCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrUploadSignedCertificate, ErrStructValidation, err)
	}

	req, err := request.NewPatch(ctx, "/ccm/v2/lineages/%d/generations/%d", params.LineageID, params.GenerationID).
		AddQueryParamIf("acknowledgeWarnings", strconv.FormatBool(params.AcknowledgeWarnings), params.AcknowledgeWarnings).
		AddHeader("Accept", mediaTypeCertificateLineageUploadV3).
		WithBody(params.Body).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrUploadSignedCertificate, err)
	}

	var result UploadSignedCertificateResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request execution failed: %w", ErrUploadSignedCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrUploadSignedCertificate, c.Error(resp))
	}

	return &result, nil
}
