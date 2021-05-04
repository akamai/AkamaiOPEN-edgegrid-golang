package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// MatchTargetSequence represents a collection of MatchTargetSequence
//
// See: MatchTargetSequence.GetMatchTargetSequence()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// MatchTargetSequence  contains operations available on MatchTargetSequence  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmatchtargetsequence
	MatchTargetSequence interface {
		GetMatchTargetSequence(ctx context.Context, params GetMatchTargetSequenceRequest) (*GetMatchTargetSequenceResponse, error)
		UpdateMatchTargetSequence(ctx context.Context, params UpdateMatchTargetSequenceRequest) (*UpdateMatchTargetSequenceResponse, error)
	}

	// GetMatchTargetSequence is the argument for GetProperties
	GetMatchTargetSequenceRequest struct {
		ConfigID      int    `json:"configId"`
		ConfigVersion int    `json:"configVersion"`
		Type          string `json:"type"`
	}

	GetMatchTargetSequenceResponse struct {
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// UpdateMatchTargetRequest is the argument for GetProperties
	UpdateMatchTargetSequenceRequest struct {
		ConfigID       int               `json:"-"`
		ConfigVersion  int               `json:"-"`
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// UpdateMatchTargetResponse ...
	UpdateMatchTargetSequenceResponse struct {
		TargetSequence []MatchTargetItem `json:"targetSequence"`
		Type           string            `json:"type"`
	}

	// BypassNetworkList ...
	MatchTargetItem struct {
		Sequence int `json:"sequence"`
		TargetID int `json:"targetId"`
	}
)

// Validate validates GetMatchTargetSequenceRequest
func (v GetMatchTargetSequenceRequest) Validate() error {
	return validation.Errors{
		"ConfigID":      validation.Validate(v.ConfigID, validation.Required),
		"ConfigVersion": validation.Validate(v.ConfigVersion, validation.Required),
		"Type":          validation.Validate(v.Type, validation.Required),
	}.Filter()
}

// Validate validates UpdateMatchTargetSequenceRequest
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
		return nil, fmt.Errorf("failed to create getmatchtargetsequence request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a MatchTargetSequence.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmatchtargetsequence

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
		return nil, fmt.Errorf("failed to create create MatchTargetSequencerequest: %w", err)
	}

	var rval UpdateMatchTargetSequenceResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create MatchTargetSequence request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
