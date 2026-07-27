package purgecache

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// InvalidateByURLRequest contains request parameters for invalidating cache content by URL or ARL.
	InvalidateByURLRequest PurgeByURLRequest

	// InvalidateByTagRequest contains request parameters for invalidating cache content by cache tags.
	InvalidateByTagRequest PurgeByTagRequest

	// InvalidateByCPCodeRequest contains request parameters for invalidating cache content by CP Codes.
	InvalidateByCPCodeRequest PurgeByCPCodeRequest

	// InvalidateResponse contains the response for a cache invalidation request.
	InvalidateResponse PurgeResponse
)

// Validate validates InvalidateByURLRequest.
func (r InvalidateByURLRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one URL or ARL to invalidate")),
	})
}

// Validate validates InvalidateByTagRequest.
func (r InvalidateByTagRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one tag to invalidate")),
	})
}

// Validate validates InvalidateByCPCodeRequest.
func (r InvalidateByCPCodeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Network": validation.Validate(r.Network),
		"Objects": validation.Validate(r.Objects, validation.Required.Error("must contain at least one CP Code to invalidate")),
	})
}

func (p *purgecache) InvalidateByURL(ctx context.Context, params InvalidateByURLRequest) (*InvalidateResponse, error) {
	p.Log(ctx).Debug("InvalidateByURL")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrInvalidate, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrInvalidate, "invalidate", PurgeTypeURL, params.Network, params)
	return (*InvalidateResponse)(result), err
}

func (p *purgecache) InvalidateByTag(ctx context.Context, params InvalidateByTagRequest) (*InvalidateResponse, error) {
	p.Log(ctx).Debug("InvalidateByTag")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrInvalidate, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrInvalidate, "invalidate", PurgeTypeTag, params.Network, params)
	return (*InvalidateResponse)(result), err
}

func (p *purgecache) InvalidateByCPCode(ctx context.Context, params InvalidateByCPCodeRequest) (*InvalidateResponse, error) {
	p.Log(ctx).Debug("InvalidateByCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrInvalidate, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrInvalidate, "invalidate", PurgeTypeCPCode, params.Network, params)
	return (*InvalidateResponse)(result), err
}
