package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The PenaltyBox interface supports retrieving or modifying the penalty box settings for
	// a specified security policy
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#penaltybox
	PenaltyBox interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpenaltybox
		// Note: this method is DEPRECATED and will be removed in a future release.
		GetPenaltyBoxes(ctx context.Context, params GetPenaltyBoxesRequest) (*GetPenaltyBoxesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getpenaltybox
		GetPenaltyBox(ctx context.Context, params GetPenaltyBoxRequest) (*GetPenaltyBoxResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpenaltybox
		UpdatePenaltyBox(ctx context.Context, params UpdatePenaltyBoxRequest) (*UpdatePenaltyBoxResponse, error)
	}

	// GetPenaltyBoxesRequest is used to retrieve the penalty box settings.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetPenaltyBoxesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Current  string `json:"current"`
		Mode     string `json:"mode"`
		Eval     string `json:"eval"`
	}

	// GetPenaltyBoxesResponse is returned from a call to GetPenaltyBoxes.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetPenaltyBoxesResponse struct {
		Action               string `json:"action,omitempty"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

	// GetPenaltyBoxRequest is used to retrieve the penalty box settings.
	GetPenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// GetPenaltyBoxResponse is returned from a call to GetPenaltyBox.
	GetPenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

	// UpdatePenaltyBoxRequest is used to modify the penalty box settings.
	UpdatePenaltyBoxRequest struct {
		ConfigID             int    `json:"-"`
		Version              int    `json:"-"`
		PolicyID             string `json:"-"`
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// UpdatePenaltyBoxResponse is returned from a call to UpdatePenaltyBox.
	UpdatePenaltyBoxResponse struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}
)

// Validate validates a GetPenaltyBoxRequest.
func (v GetPenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetPenaltyBoxesRequest.
// Note: this method is DEPRECATED and will be removed in a future release.
func (v GetPenaltyBoxesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdatePenaltyBoxRequest.
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
		return nil, fmt.Errorf("failed to create GetPenaltyBox request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetPenaltyBox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Note: this method is DEPRECATED and will be removed in a future release.
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
		return nil, fmt.Errorf("failed to create GetPenaltyBoxes request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetPenaltyBoxes request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdatePenaltyBox: %w", err)
	}

	var rval UpdatePenaltyBoxResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdatePenaltyBox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
