package purgecache

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// DeleteByURLRequest contains request parameters for deleting cache content by URL or ARL.
	DeleteByURLRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of URLs or ARLs to delete.
		Objects []string `json:"objects"`
	}

	// DeleteByTagRequest contains request parameters for deleting cache content by cache tags.
	DeleteByTagRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of cache tags to delete.
		Objects []string `json:"objects"`
	}

	// DeleteByCPCodeRequest contains request parameters for deleting cache content by CP Codes.
	DeleteByCPCodeRequest struct {
		// Network specifies the network environment: `staging` or `production`.
		// When omitted, the API defaults to production.
		Network PurgeNetwork `json:"-"`

		// Objects is the list of CP codes to delete.
		Objects []int64 `json:"objects"`
	}

	// DeleteResponse contains the response for a cache deletion request.
	DeleteResponse struct {
		// DescribedBy is a URL that describes the API's error response.
		DescribedBy string `json:"describedBy"`

		// Detail contains detailed information about the HTTP status code returned.
		Detail string `json:"detail"`

		// EstimatedSeconds is the estimated number of seconds before the purge is complete.
		EstimatedSeconds int64 `json:"estimatedSeconds"`

		// HTTPStatus is the HTTP code that indicates the status of the delete request.
		HTTPStatus int64 `json:"httpStatus"`

		// PurgeID is the unique identifier for the purge request.
		PurgeID string `json:"purgeId"`

		// SupportID is an identifier to provide Akamai Technical Support if issues arise.
		SupportID string `json:"supportId"`

		// Title describes the response type.
		Title string `json:"title"`

		// RateLimitHeaders contains rate limit information extracted from response headers.
		RateLimitHeaders RateLimitHeaders `json:"-"`
	}
)

// Validate validates DeleteByURLRequest.
func (r DeleteByURLRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one URL or ARL to delete")),
	})
}

// Validate validates DeleteByTagRequest.
func (r DeleteByTagRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one tag to delete")),
	})
}

// Validate validates DeleteByCPCodeRequest.
func (r DeleteByCPCodeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one CP Code to delete")),
	})
}

func (p *purgecache) DeleteByURL(ctx context.Context, params DeleteByURLRequest) (*DeleteResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeleteByURL")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	var err error
	var req *http.Request
	if params.Network == "" {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/url").WithBody(params).Build()
	} else {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/url/%s", params.Network).WithBody(params).Build()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrDelete, err)
	}

	var result DeleteResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrDelete, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	result.RateLimitHeaders = extractRateLimitStatusHeaders(resp, logger)

	return &result, nil
}

func (p *purgecache) DeleteByTag(ctx context.Context, params DeleteByTagRequest) (*DeleteResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeleteByTag")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	var err error
	var req *http.Request
	if params.Network == "" {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/tag").WithBody(params).Build()
	} else {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/tag/%s", params.Network).WithBody(params).Build()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrDelete, err)
	}

	var result DeleteResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrDelete, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	result.RateLimitHeaders = extractRateLimitStatusHeaders(resp, logger)

	return &result, nil
}

func (p *purgecache) DeleteByCPCode(ctx context.Context, params DeleteByCPCodeRequest) (*DeleteResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeleteByCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	var err error
	var req *http.Request
	if params.Network == "" {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/cpcode").WithBody(params).Build()
	} else {
		req, err = request.NewPost(ctx, "/ccu/v3/delete/cpcode/%s", params.Network).WithBody(params).Build()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrDelete, err)
	}

	var result DeleteResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrDelete, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	result.RateLimitHeaders = extractRateLimitStatusHeaders(resp, logger)

	return &result, nil
}
