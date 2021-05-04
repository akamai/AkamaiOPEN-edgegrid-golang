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
	AttackGroup interface {
		GetAttackGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error)
		GetAttackGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error)
		UpdateAttackGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error)
	}

	AttackGroupConditionException struct {
		AdvancedExceptionsList *AttackGroupAdvancedExceptions `json:"advancedExceptions,omitempty"`
		Exception              *AttackGroupException          `json:"exception,omitempty"`
	}

	AttackGroupAdvancedExceptions struct {
		ConditionOperator                       string                                                      `json:"conditionOperator,omitempty"`
		Conditions                              *AttackGroupConditions                                      `json:"conditions,omitempty"`
		HeaderCookieOrParamValues               *AttackGroupHeaderCookieOrParamValuesAdvanced               `json:"headerCookieOrParamValues,omitempty"`
		SpecificHeaderCookieOrParamNameValue    *AttackGroupSpecificHeaderCookieOrParamNameValAdvanced      `json:"specificHeaderCookieOrParamNameValue,omitempty"`
		SpecificHeaderCookieParamXMLOrJSONNames *AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced `json:"specificHeaderCookieParamXmlOrJsonNames,omitempty"`
	}

	AttackGroupConditions []struct {
		Type          string   `json:"type,omitempty"`
		Extensions    []string `json:"extensions,omitempty"`
		Filenames     []string `json:"filenames,omitempty"`
		Hosts         []string `json:"hosts,omitempty"`
		Ips           []string `json:"ips,omitempty"`
		Methods       []string `json:"methods,omitempty"`
		Paths         []string `json:"paths,omitempty"`
		Header        string   `json:"header,omitempty"`
		CaseSensitive bool     `json:"caseSensitive,omitempty"`
		Name          string   `json:"name,omitempty"`
		NameCase      bool     `json:"nameCase,omitempty"`
		PositiveMatch bool     `json:"positiveMatch"`
		Value         string   `json:"value,omitempty"`
		Wildcard      bool     `json:"wildcard,omitempty"`
		ValueCase     bool     `json:"valueCase,omitempty"`
		ValueWildcard bool     `json:"valueWildcard,omitempty"`
		UseHeaders    bool     `json:"useHeaders,omitempty"`
	}

	AttackGroupAdvancedCriteria []struct {
		Hostnames []string `json:"hostnames,omitempty"`
		Names     []string `json:"names,omitempty"`
		Paths     []string `json:"paths,omitempty"`
		Values    []string `json:"values,omitempty"`
	}

	AttackGroupSpecificHeaderCookieOrParamNameValAdvanced []struct {
		Criteria    *AttackGroupAdvancedCriteria `json:"criteria,omitempty"`
		NamesValues []struct {
			Names  []string `json:"names"`
			Values []string `json:"values"`
		} `json:"namesValues"`
		Selector      string `json:"selector"`
		ValueWildcard bool   `json:"valueWildcard"`
		Wildcard      bool   `json:"wildcard"`
	}

	AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced []struct {
		Criteria *AttackGroupAdvancedCriteria `json:"criteria,omitempty"`
		Names    []string                     `json:"names,omitempty"`
		Selector string                       `json:"selector,omitempty"`
		Wildcard bool                         `json:"wildcard,omitempty"`
	}

	AttackGroupHeaderCookieOrParamValuesAdvanced []struct {
		Criteria      *AttackGroupAdvancedCriteria `json:"criteria,omitempty"`
		ValueWildcard bool                         `json:"valueWildcard"`
		Values        []string                     `json:"values,omitempty"`
	}

	AttackGroupException struct {
		SpecificHeaderCookieParamXMLOrJSONNames *AttackGroupSpecificHeaderCookieParamXMLOrJSONNames `json:"specificHeaderCookieParamXmlOrJsonNames,omitempty"`
	}

	AttackGroupSpecificHeaderCookieParamXMLOrJSONNames []struct {
		Names    []string `json:"names,omitempty"`
		Selector string   `json:"selector,omitempty"`
		Wildcard bool     `json:"wildcard,omitempty"`
	}

	GetAttackGroupsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group,omitempty"`
	}

	GetAttackGroupsResponse struct {
		AttackGroups []struct {
			Group              string                         `json:"group,omitempty"`
			Action             string                         `json:"action,omitempty"`
			ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
		} `json:"attackGroupActions,omitempty"`
	}

	GetAttackGroupRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAttackGroupResponse struct {
		Action             string                         `json:"action,omitempty"`
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}

	UpdateAttackGroupRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		Group          string          `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	UpdateAttackGroupResponse struct {
		Action             string                         `json:"action,omitempty"`
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}
)

// Check Condition Exception is Empty
func (r GetAttackGroupResponse) IsEmptyCodnitionException() bool {
	if r.ConditionException == nil {
		return true
	}
	return false
}

// Validate validates GetAttackGroupConditionExceptionRequest
func (v GetAttackGroupRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates GetAttackGroupConditionExceptionsRequest
func (v GetAttackGroupsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateAttackGroupConditionExceptionRequest
func (v UpdateAttackGroupRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetAttackGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroup")

	var rval GetAttackGroupResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getattackgroup request: %w", err)
	}
	logger.Debugf("BEFORE GetAttackGroup %v", rval)
	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroup  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}
	logger.Debugf("GetAttackGroup %v", rval)
	return &rval, nil

}

func (p *appsec) GetAttackGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupConditionExceptions")

	var rval GetAttackGroupsResponse
	var rvalfiltered GetAttackGroupsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

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

	if params.Group != "" {
		for k, val := range rval.AttackGroups {
			if val.Group == params.Group {
				rvalfiltered.AttackGroups = append(rvalfiltered.AttackGroups, rval.AttackGroups[k])
			}
		}
	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a AttackGroupConditionException.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception

func (p *appsec) UpdateAttackGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAttackGroup")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AttackGroupRequest: %w", err)
	}

	var rval UpdateAttackGroupResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AttackGroup request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
