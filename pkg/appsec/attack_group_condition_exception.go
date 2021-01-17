package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AttackGroupConditionException represents a collection of AttackGroupConditionException
//
// See: AttackGroupConditionException.GetAttackGroupConditionException()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AttackGroupConditionException  contains operations available on AttackGroupConditionException  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getattackgroupconditionexception
	AttackGroupConditionException interface {
		GetAttackGroupConditionExceptions(ctx context.Context, params GetAttackGroupConditionExceptionsRequest) (*GetAttackGroupConditionExceptionsResponse, error)
		GetAttackGroupConditionException(ctx context.Context, params GetAttackGroupConditionExceptionRequest) (*GetAttackGroupConditionExceptionResponse, error)
		UpdateAttackGroupConditionException(ctx context.Context, params UpdateAttackGroupConditionExceptionRequest) (*UpdateAttackGroupConditionExceptionResponse, error)
		RemoveAttackGroupConditionException(ctx context.Context, params RemoveAttackGroupConditionExceptionRequest) (*RemoveAttackGroupConditionExceptionResponse, error)
	}

	GetAttackGroupConditionExceptionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAttackGroupConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group,omitempty"`
	}

	GetAttackGroupConditionExceptionsResponse struct {
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
			SpecificHeaderCookieOrParamPrefix struct {
				Prefix   string `json:"prefix,omitempty"`
				Selector string `json:"selector,omitempty"`
			} `json:"specificHeaderCookieOrParamPrefix,omitempty"`
		} `json:"exception,omitempty"`
	}

	GetAttackGroupConditionExceptionResponse struct {
		Conditions []interface{} `json:"conditions,omitempty"`
		Exception  struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues,omitempty"`
			SpecificHeaderCookieOrParamNames []struct {
				Names    []string `json:"names,omitempty"`
				Selector string   `json:"selector,omitempty"`
			} `json:"specificHeaderCookieOrParamNames,omitempty"`
			SpecificHeaderCookieOrParamPrefix struct {
				Prefix   string `json:"prefix,omitempty"`
				Selector string `json:"selector,omitempty"`
			} `json:"specificHeaderCookieOrParamPrefix,omitempty"`
		} `json:"exception,omitempty"`
	}

	UpdateAttackGroupConditionExceptionRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		Group          string          `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	UpdateAttackGroupConditionExceptionResponse struct {
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
			SpecificHeaderCookieOrParamPrefix struct {
				Prefix   string `json:"prefix"`
				Selector string `json:"selector"`
			} `json:"specificHeaderCookieOrParamPrefix"`
		} `json:"exception"`
	}

	RemoveAttackGroupConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"-"`
		Empty    string `json:"empty"`
	}

	RemoveAttackGroupConditionExceptionResponse struct {
		Conditions []interface{} `json:"conditions"`
		Exception  struct {
			HeaderCookieOrParamValues        []string `json:"headerCookieOrParamValues"`
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
)

// Validate validates GetAttackGroupConditionExceptionRequest
func (v GetAttackGroupConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetAttackGroupConditionExceptionsRequest
func (v GetAttackGroupConditionExceptionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateAttackGroupConditionExceptionRequest
func (v UpdateAttackGroupConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateRuleConditionExceptionRequest
func (v RemoveAttackGroupConditionExceptionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
	}.Filter()
}

func (p *appsec) GetAttackGroupConditionException(ctx context.Context, params GetAttackGroupConditionExceptionRequest) (*GetAttackGroupConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupConditionException")

	var rval GetAttackGroupConditionExceptionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getattackgroupconditionexception request: %w", err)
	}
	logger.Debugf("BEFORE GetAttackGroupConditionException %v", rval)
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroupconditionexception  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}
	logger.Debugf("GetAttackGroupConditionException %v", rval)
	return &rval, nil

}

func (p *appsec) GetAttackGroupConditionExceptions(ctx context.Context, params GetAttackGroupConditionExceptionsRequest) (*GetAttackGroupConditionExceptionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupConditionExceptions")

	var rval GetAttackGroupConditionExceptionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getattackgroupconditionexceptions request: %w", err)
	}
	logger.Debugf("BEFORE GetAttackGroupConditionException %v", rval)
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroupconditionexceptions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a AttackGroupConditionException.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception

func (p *appsec) UpdateAttackGroupConditionException(ctx context.Context, params UpdateAttackGroupConditionExceptionRequest) (*UpdateAttackGroupConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAttackGroupConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AttackGroupConditionExceptionrequest: %w", err)
	}

	var rval UpdateAttackGroupConditionExceptionResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create AttackGroupConditionException request failed: %w", err)
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

func (p *appsec) RemoveAttackGroupConditionException(ctx context.Context, params RemoveAttackGroupConditionExceptionRequest) (*RemoveAttackGroupConditionExceptionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("RemoveAttackGroupConditionException")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s/condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remove AttackGroupConditionExceptionrequest: %w", err)
	}

	var rval RemoveAttackGroupConditionExceptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create RemoveAttackGroupConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
