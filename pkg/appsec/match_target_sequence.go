package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The MatchTargetSequence interface supports querying and modifying the order of match targets.
	MatchTargetSequence interface {
		// GetMatchTargetSequence returns match targets defined in the specified security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-match-targets
		GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error)

		// UpdateMatchTargetSequence updates the sequence of Match Targets in a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-match-targets-sequence
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
	logger := p.Log(ctx)
	logger.Debug("GetMatchTargetSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

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

	var result GetMatchTargetSequenceResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get match target sequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateMatchTargetSequence")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/match-targets/sequence",
		params.ConfigID,
		params.ConfigVersion,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateMatchTargetSequence request: %w", err)
	}

	var result UpdateMatchTargetSequenceResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update match target sequence request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
