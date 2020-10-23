package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PenaltyBox represents a collection of PenaltyBox
//
// See: PenaltyBox.GetPenaltyBox()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// PenaltyBox  contains operations available on PenaltyBox  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpenaltybox
	PenaltyBox interface {
		GetPenaltyBoxes(ctx context.Context, params GetPenaltyBoxesRequest) (*GetPenaltyBoxesResponse, error)
		GetPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error)
		UpdatePenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error)
	}

	GetPenaltyBoxesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	GetPenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	GetPenaltyBoxesResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	GetPenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	UpdatePenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	UpdatePenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}
)

// Validate validates GetPenaltyBoxRequest
func (v GetPenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetPenaltyBoxsRequest
func (v GetPenaltyBoxesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdatePenaltyBoxRequest
func (v UpdatePenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetPenaltyBox")

	var rval GetPenaltyBoxResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getpenaltybox request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getpenaltybox  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetPenaltyBoxes(ctx context.Context, params GetPenaltyBoxesRequest) (*GetPenaltyBoxesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetPenaltyBoxs")

	var rval GetPenaltyBoxesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getpenaltyboxs request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getpenaltyboxs request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a PenaltyBox.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpenaltybox

func (p *appsec) UpdatePenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdatePenaltyBox")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/penalty-box",
		params.ConfigID,
		params.Version,
		params.PolicyID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create PenaltyBoxrequest: %w", err)
	}

	var rval UpdatePenaltyBoxResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create PenaltyBox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
