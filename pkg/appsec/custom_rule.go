package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CustomRule represents a collection of CustomRule
//
// See: CustomRule.GetCustomRule()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// CustomRule  contains operations available on CustomRule  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomrule
	CustomRule interface {
		GetCustomRules(ctx context.Context, params GetCustomRulesRequest) (*GetCustomRulesResponse, error)
		GetCustomRule(ctx context.Context, params GetCustomRuleRequest) (*GetCustomRuleResponse, error)
		CreateCustomRule(ctx context.Context, params CreateCustomRuleRequest) (*CreateCustomRuleResponse, error)
		UpdateCustomRule(ctx context.Context, params UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error)
		RemoveCustomRule(ctx context.Context, params RemoveCustomRuleRequest) (*RemoveCustomRuleResponse, error)
	}

	CustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	CustomRulesResponse struct {
		CustomRules []struct {
			ID      int    `json:"id"`
			Link    string `json:"link"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			Version int    `json:"version"`
		} `json:"customRules"`
	}

	GetCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	GetCustomRulesResponse struct {
		CustomRules []struct {
			ID      int    `json:"id"`
			Link    string `json:"link"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			Version int    `json:"version"`
		} `json:"customRules"`
	}

	GetCustomRulesRequest struct {
		ConfigID int `json:"configid,omitempty"`
	}

	GetCustomRuleRequest struct {
		ConfigID      int      `json:"configid,omitempty"`
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	CreateCustomRuleRequest struct {
		ConfigID      int      `json:"configid,omitempty"`
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	CreateCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	UpdateCustomRuleRequest struct {
		ConfigID      int      `json:"configid,omitempty"`
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	UpdateCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	RemoveCustomRuleRequest struct {
		ConfigID      int      `json:"configid,omitempty"`
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}

	RemoveCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Version       int      `json:"version,omitempty"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type          string   `json:"type"`
			PositiveMatch bool     `json:"positiveMatch"`
			Value         []string `json:"value,omitempty"`
			ValueWildcard bool     `json:"valueWildcard,omitempty"`
			ValueCase     bool     `json:"valueCase,omitempty"`
			NameWildcard  bool     `json:"nameWildcard,omitempty"`
			Name          []string `json:"name,omitempty"`
			NameCase      bool     `json:"nameCase,omitempty"`
		} `json:"conditions"`
	}
)

// Validate validates GetCustomRuleRequest
func (v GetCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates GetCustomRulesRequest
func (v GetCustomRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates CreateCustomRuleRequest
func (v CreateCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates UpdateCustomRuleRequest
func (v UpdateCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates RemoveCustomRuleRequest
func (v RemoveCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

func (p *appsec) GetCustomRule(ctx context.Context, params GetCustomRuleRequest) (*GetCustomRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomRule")

	var rval GetCustomRuleResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules/%d",
		params.ConfigID,
		params.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcustomrule request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil

}

func (p *appsec) GetCustomRules(ctx context.Context, params GetCustomRulesRequest) (*GetCustomRulesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomRules")

	var rval GetCustomRulesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcustomrules request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getcustomrules request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil

}

// Update will update a CustomRule.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putcustomrule

func (p *appsec) UpdateCustomRule(ctx context.Context, params UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateCustomRule")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules/%d",
		params.ConfigID,
		params.ID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create CustomRulerequest: %w", err)
	}

	var rval UpdateCustomRuleResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create CustomRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil
}

// Create will create a new customrule.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postcustomrule
func (p *appsec) CreateCustomRule(ctx context.Context, params CreateCustomRuleRequest) (*CreateCustomRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateCustomRule")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create customrule request: %w", err)
	}

	var rval CreateCustomRuleResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create customrulerequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil

}

// Delete will delete a CustomRule
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletecustomrule

func (p *appsec) RemoveCustomRule(ctx context.Context, params RemoveCustomRuleRequest) (*RemoveCustomRuleResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveCustomRuleResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveCustomRule")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules/%d",
		params.ConfigID,
		params.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delcustomrule request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delcustomrule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil
}
