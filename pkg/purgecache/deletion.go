package purgecache

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// DeleteByURLRequest contains request parameters for deleting cache content by URL or ARL.
	DeleteByURLRequest PurgeByURLRequest

	// DeleteByTagRequest contains request parameters for deleting cache content by cache tags.
	DeleteByTagRequest PurgeByTagRequest

	// DeleteByCPCodeRequest contains request parameters for deleting cache content by CP Codes.
	DeleteByCPCodeRequest PurgeByCPCodeRequest

	// DeleteResponse contains the response for a cache deletion request.
	DeleteResponse PurgeResponse
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
	p.Log(ctx).Debug("DeleteByURL")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrDelete, "delete", PurgeTypeURL, params.Network, params)
	return (*DeleteResponse)(result), err
}

func (p *purgecache) DeleteByTag(ctx context.Context, params DeleteByTagRequest) (*DeleteResponse, error) {
	p.Log(ctx).Debug("DeleteByTag")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrDelete, "delete", PurgeTypeTag, params.Network, params)
	return (*DeleteResponse)(result), err
}

func (p *purgecache) DeleteByCPCode(ctx context.Context, params DeleteByCPCodeRequest) (*DeleteResponse, error) {
	p.Log(ctx).Debug("DeleteByCPCode")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w: %w", ErrDelete, ErrStructValidation, err)
	}

	result, err := p.performPurgeRequest(ctx, ErrDelete, "delete", PurgeTypeCPCode, params.Network, params)
	return (*DeleteResponse)(result), err
}
