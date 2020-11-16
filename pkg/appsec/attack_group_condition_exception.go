package appsec

import (
	"context"
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
		Group    string `json:"group"`
	}

	GetAttackGroupConditionExceptionsResponse struct {
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

	GetAttackGroupConditionExceptionResponse struct {
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

	UpdateAttackGroupConditionExceptionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
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

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroupconditionexception  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

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
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
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
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AttackGroupConditionException request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
