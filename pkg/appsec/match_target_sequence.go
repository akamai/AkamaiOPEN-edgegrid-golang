package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The MatchTargetSequence interface supports querying and modifying the order of match targets.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#matchtargetorder
	MatchTargetSequence interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmatchtargets
		GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsequence
		UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error)
	}

	// GetMatchTargetSequenceRequest is used to retrieve the sequence of match targets for a configuration.
	GetMatchTargetSequenceRequest struct {
		ConfigID      int    `json:"configId"`
		ConfigVersion int    `json:"configVersion"`
		Type          string `json:"type"`
	}

	// GetMatchTargetSequenceResponse is returned from a call to GetMatchTargetSequence.
	GetMatchTargetSequenceResponse struct {
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// UpdateMatchTargetSequenceRequest UpdateMatchTargetSequenceRequest is used to modify an existing match target sequence.
	UpdateMatchTargetSequenceRequest struct {
		ConfigID       int               `json:"-"`
		ConfigVersion  int               `json:"-"`
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// UpdateMatchTargetSequenceResponse is returned from a call to UpdateMatchTargetSequence.
	UpdateMatchTargetSequenceResponse struct {
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// MatchTargetItem describes a match target and its sequence number.
	MatchTargetItem struct {
		Sequence int `json:"sequence"`
		TargetID int `json:"targetId"`
	}
)

// Validate validates a GetMatchTargetSequenceRequest.
func (v GetMatchTargetSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Type":          validation.Validate(v.Type, validation.Required),
	}.Filter()
}

// Validate validates an UpdateMatchTargetSequenceRequest.
func (v UpdateMatchTargetSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Type":          validation.Validate(v.ConfigVersion, validation.Required),
	}.Filter()
}

func (p *appsec) GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetMatchTargetSequence")

	var rval GetMatchTargetSequenceResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets/sequence?type=%s",
		params.ConfigID,
		params.ConfigVersion,
		params.Type,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetMatchTargetSequence request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetMatchTargetSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateMatchTargetSequence")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets/sequence",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateMatchTargetSequence request: %w", err)
	}

	var rval UpdateMatchTargetSequenceResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateMatchTargetSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
