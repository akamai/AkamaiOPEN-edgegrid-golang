package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EvalRuleConditionException represents a collection of EvalRuleConditionException
//
// See: EvalRuleConditionException.GetEvalRuleConditionException()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// EvalRuleConditionException  contains operations available on EvalRuleConditionException  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalruleconditionexception
	EvalRuleConditionException interface {
		GetEvalRuleConditionExceptions(ctx context.Context, params GetEvalRuleConditionExceptionsRequest) (*GetEvalRuleConditionExceptionsResponse, error)
		GetEvalRuleConditionException(ctx context.Context, params GetEvalRuleConditionExceptionRequest) (*GetEvalRuleConditionExceptionResponse, error)
		UpdateEvalRuleConditionException(ctx context.Context, params UpdateEvalRuleConditionExceptionRequest) (*UpdateEvalRuleConditionExceptionResponse, error)
		RemoveEvalRuleConditionException(ctx context.Context, params RemoveEvalRuleConditionExceptionRequest) (*RemoveEvalRuleConditionExceptionResponse, error)
	}

	GetEvalRuleConditionExceptionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
	}

	GetEvalRuleConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"ruleId"`
	}

	GetEvalRuleConditionExceptionsResponse struct {
		Conditions []struct {
			Type          string   `json:"type"`
			Extensions    []string `json:"extensions,omitempty"`
			PositiveMatch bool     `json:"positiveMatch"`
			Filenames     []string `json:"filenames,omitempty"`
			Hosts         []string `json:"hosts,omitempty"`
			Ips           []string `json:"ips,omitempty"`
			UseHeaders    bool     `json:"useHeaders,omitempty"`
			CaseSensitive bool     `json:"caseSensitive,omitempty"`
			Name          string   `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
			Value         string   `json:"value,omitempty"`
			Wildcard      bool     `json:"wildcard,omitempty"`
			Header        string   `json:"header,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			Methods       []string `json:"methods,omitempty"`
			Paths         []string `json:"paths,omitempty"`
		} `json:"conditions"`
		Exception struct {
			HeaderCookieOrParamValues            []string `json:"headerCookieOrParamValues"`
			SpecificHeaderCookieOrParamNameValue struct {
				Name     string `json:"name"`
				Selector string `json:"selector"`
				Value    string `json:"value"`
			} `json:"specificHeaderCookieOrParamNameValue"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names"`
				Selector string   `json:"selector"`
			} `json:"specificHeaderCookieOrParamNames"`
			SpecificHeaderCookieOrParamPrefix struct {
				Prefix   string `json:"prefix"`
				Selector string `json:"selector"`
			} `json:"specificHeaderCookieOrParamPrefix"`
		} `json:"exception"`
	}

	GetEvalRuleConditionExceptionResponse struct {
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

	UpdateEvalRuleConditionExceptionRequest struct {
		ConfigID   int    `json:"-"`
		Version    int    `json:"-"`
		PolicyID   string `json:"-"`
		RuleID     int    `json:"-"`
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

	UpdateEvalRuleConditionExceptionResponse struct {
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

	RemoveEvalRuleConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
		Empty    string `json:"empty"`
	}

	RemoveEvalRuleConditionExceptionResponse struct {
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

// Validate validates GetEvalRuleConditionExceptionRequest
func (v GetEvalRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetEvalRuleConditionExceptionsRequest
func (v GetEvalRuleConditionExceptionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateEvalRuleConditionExceptionRequest
func (v UpdateEvalRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleConditionExceptionRequest
func (v RemoveEvalRuleConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalRuleConditionException(ctx context.Context, params GetEvalRuleConditionExceptionRequest) (*GetEvalRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRuleConditionException")

	var rval GetEvalRuleConditionExceptionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleconditionexception request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleconditionexception  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetEvalRuleConditionExceptions(ctx context.Context, params GetEvalRuleConditionExceptionsRequest) (*GetEvalRuleConditionExceptionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetEvalRuleConditionExceptions")

	var rval GetEvalRuleConditionExceptionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getevalruleconditionexceptions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getevalruleconditionexceptions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a EvalRuleConditionException.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalruleconditionexception

func (p *appsec) UpdateEvalRuleConditionException(ctx context.Context, params UpdateEvalRuleConditionExceptionRequest) (*UpdateEvalRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateEvalRuleConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create EvalRuleConditionExceptionrequest: %w", err)
	}

	var rval UpdateEvalRuleConditionExceptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create EvalRuleConditionException request failed: %w", err)
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

func (p *appsec) RemoveEvalRuleConditionException(ctx context.Context, params RemoveEvalRuleConditionExceptionRequest) (*RemoveEvalRuleConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveRuleConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remove EvalRuleConditionExceptionrequest: %w", err)
	}

	var rval RemoveEvalRuleConditionExceptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RemoveRuleConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
