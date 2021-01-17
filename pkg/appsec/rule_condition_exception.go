package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RuleConditionException represents a collection of RuleConditionException
//
// See: RuleConditionException.GetRuleConditionException()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// RuleConditionException  contains operations available on RuleConditionException  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleconditionexception
	RuleConditionException interface {
		GetRuleConditionExceptions(ctx context.Context, params GetRuleConditionExceptionsRequest) (*GetRuleConditionExceptionsResponse, error)
		GetRuleConditionException(ctx context.Context, params GetRuleConditionExceptionRequest) (*GetRuleConditionExceptionResponse, error)
		UpdateRuleConditionException(ctx context.Context, params UpdateRuleConditionExceptionRequest) (*UpdateRuleConditionExceptionResponse, error)
		RemoveRuleConditionException(ctx context.Context, params RemoveRuleConditionExceptionRequest) (*RemoveRuleConditionExceptionResponse, error)
	}

	GetRuleConditionExceptionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRuleConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	GetRuleConditionExceptionsResponse struct {
		Conditions []struct {
			Type          string   `json:"type,omitempty"`
			Filenames     []string `json:"filenames,omitempty"`
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			Methods       []string `json:"methods,omitempty"`
		} `json:"conditions,omitempty"`
		Exception struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues,omitempty"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names,omitempty"`
				Selector string   `json:"selector,omitempty"`
			} `json:"specificHeaderCookieOrParamNames,omitempty"`
		} `json:"exception,omitempty"`
	}

	GetRuleConditionExceptionResponse struct {
		Conditions []struct {
			Type          string   `json:"type,omitempty"`
			Filenames     []string `json:"filenames,omitempty"`
			PositiveMatch bool     `json:"positiveMatch,omitempty"`
			Methods       []string `json:"methods,omitempty"`
		} `json:"conditions,omitempty"`
		Exception struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues,omitempty"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names,omitempty"`
				Selector string   `json:"selector,omitempty"`
			} `json:"specificHeaderCookieOrParamNames,omitempty"`
		} `json:"exception,omitempty"`
	}

	UpdateRuleConditionExceptionRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		RuleID         int             `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	UpdateRuleConditionExceptionResponse struct {
		Conditions []struct {
			Type          string   `json:"type"`
			Filenames     []string `json:"filenames,omitempty"`
			PositiveMatch bool     `json:"positiveMatch"`
			Methods       []string `json:"methods,omitempty"`
		} `json:"conditions"`
		Exception struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names"`
				Selector string   `json:"selector"`
			} `json:"specificHeaderCookieOrParamNames"`
		} `json:"exception"`
	}

	RemoveRuleConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Empty    string `json:"empty"`
	}

	RemoveRuleConditionExceptionResponse struct {
		Conditions []struct {
			Type          string   `json:"type"`
			Filenames     []string `json:"filenames,omitempty"`
			PositiveMatch bool     `json:"positiveMatch"`
			Methods       []string `json:"methods,omitempty"`
		} `json:"conditions"`
		Exception struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names"`
				Selector string   `json:"selector"`
			} `json:"specificHeaderCookieOrParamNames"`
		} `json:"exception"`
	}
)

// Validate validates GetRuleConditionExceptionRequest
func (v GetRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetRuleConditionExceptionsRequest
func (v GetRuleConditionExceptionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleConditionExceptionRequest
func (v UpdateRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleConditionExceptionRequest
func (v RemoveRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetRuleConditionException(ctx context.Context, params GetRuleConditionExceptionRequest) (*GetRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleConditionException")

	var rval GetRuleConditionExceptionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getruleconditionexception request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getruleconditionexception  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetRuleConditionExceptions(ctx context.Context, params GetRuleConditionExceptionsRequest) (*GetRuleConditionExceptionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleConditionExceptions")

	var rval GetRuleConditionExceptionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getruleconditionexceptions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getruleconditionexceptions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a RuleConditionException.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleconditionexception

func (p *appsec) UpdateRuleConditionException(ctx context.Context, params UpdateRuleConditionExceptionRequest) (*UpdateRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRuleConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create RuleConditionExceptionrequest: %w", err)
	}

	var rval UpdateRuleConditionExceptionResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create RuleConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Remove will remove a RuleConditionException.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleconditionexception

func (p *appsec) RemoveRuleConditionException(ctx context.Context, params RemoveRuleConditionExceptionRequest) (*RemoveRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveRuleConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create RuleConditionExceptionrequest: %w", err)
	}

	var rval RemoveRuleConditionExceptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RuleConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
