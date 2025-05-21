package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AttackGroup interface supports retrieving and updating attack groups along with their
	// associated actions, conditions, and exceptions.
	AttackGroup interface {
		// GetAttackGroups returns a list of attack groups with their associated actions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-attack-groups-1
		GetAttackGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error)

		// GetAttackGroup returns the action for the attack group.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-policy-attack-group-1
		GetAttackGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error)

		// UpdateAttackGroup updates what action to take when an attack group's rule triggers.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-attack-group-1
		UpdateAttackGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error)
	}

	// AttackGroupConditionException describes an attack group's condition and exception information.
	AttackGroupConditionException struct {
		AdvancedExceptionsList *AttackGroupAdvancedExceptions `json:"advancedExceptions,omitempty"`
		Exception              *AttackGroupException          `json:"exception,omitempty"`
	}

	// AttackGroupAdvancedExceptions describes an attack group's advanced exception information.
	AttackGroupAdvancedExceptions struct {
		ConditionOperator                       string                                                      `json:"conditionOperator,omitempty"`
		Conditions                              *AttackGroupConditions                                      `json:"conditions,omitempty"`
		HeaderCookieOrParamValues               *AttackGroupHeaderCookieOrParamValuesAdvanced               `json:"headerCookieOrParamValues,omitempty"`
		SpecificHeaderCookieOrParamNameValue    *AttackGroupSpecificHeaderCookieOrParamNameValAdvanced      `json:"specificHeaderCookieOrParamNameValue,omitempty"`
		SpecificHeaderCookieParamXMLOrJSONNames *AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced `json:"specificHeaderCookieParamXmlOrJsonNames,omitempty"`
	}

	// AttackGroupConditions describes an attack group's condition information.
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
		ClientLists   []string `json:"clientLists,omitempty"`
	}

	// AttackGroupAdvancedCriteria describes the hostname and path criteria used to limit the scope of an exception.
	AttackGroupAdvancedCriteria []struct {
		Hostnames []string `json:"hostnames,omitempty"`
		Names     []string `json:"names,omitempty"`
		Paths     []string `json:"paths,omitempty"`
		Values    []string `json:"values,omitempty"`
	}

	// AttackGroupSpecificHeaderCookieOrParamNameValAdvanced describes the excepted name-value pairs in a request.
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

	// AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced describes the advanced exception members that allow you to conditionally exclude requests from inspection.
	AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced []struct {
		Criteria *AttackGroupAdvancedCriteria `json:"criteria,omitempty"`
		Names    []string                     `json:"names,omitempty"`
		Selector string                       `json:"selector,omitempty"`
		Wildcard bool                         `json:"wildcard,omitempty"`
	}

	// AttackGroupHeaderCookieOrParamValuesAdvanced describes the list of excepted values in headers, cookies, or query parameters.
	AttackGroupHeaderCookieOrParamValuesAdvanced []struct {
		Criteria      *AttackGroupAdvancedCriteria `json:"criteria,omitempty"`
		ValueWildcard bool                         `json:"valueWildcard"`
		Values        []string                     `json:"values,omitempty"`
	}

	// AttackGroupException is used to describe an exception that can be used to conditionally exclude requests from inspection.
	AttackGroupException struct {
		SpecificHeaderCookieParamXMLOrJSONNames *AttackGroupSpecificHeaderCookieParamXMLOrJSONNames `json:"specificHeaderCookieParamXmlOrJsonNames,omitempty"`
	}

	// AttackGroupSpecificHeaderCookieParamXMLOrJSONNames describes the advanced exception members that can be used to conditionally exclude requests from inspection.
	AttackGroupSpecificHeaderCookieParamXMLOrJSONNames []struct {
		Names    []string `json:"names,omitempty"`
		Selector string   `json:"selector,omitempty"`
		Wildcard bool     `json:"wildcard,omitempty"`
	}

	// GetAttackGroupsRequest is used to retrieve a list of attack groups with their associated actions.
	GetAttackGroupsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group,omitempty"`
	}

	// GetAttackGroupsResponse is returned from a call to GetAttackGroups.
	GetAttackGroupsResponse struct {
		AttackGroups []struct {
			Group              string                         `json:"group,omitempty"`
			Action             string                         `json:"action,omitempty"`
			ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
		} `json:"attackGroupActions,omitempty"`
	}

	// GetAttackGroupRequest is used to retrieve a list of attack groups with their associated actions.
	GetAttackGroupRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	// GetAttackGroupResponse is returned from a call to GetAttackGroup.
	GetAttackGroupResponse struct {
		Action             string                         `json:"action,omitempty"`
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}

	// UpdateAttackGroupRequest is used to modify what action to take when an attack group’s rule triggers.
	UpdateAttackGroupRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		Group          string          `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	// UpdateAttackGroupResponse is returned from a call to UpdateAttackGroup.
	UpdateAttackGroupResponse struct {
		Action             string                         `json:"action,omitempty"`
		ConditionException *AttackGroupConditionException `json:"conditionException,omitempty"`
	}
)

// IsEmptyConditionException checks whether an attack group's ConditionException field is empty.
func (r GetAttackGroupResponse) IsEmptyConditionException() bool {
	return r.ConditionException == nil
}

// Validate validates a GetAttackGroupConditionExceptionRequest.
func (v GetAttackGroupRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a GetAttackGroupConditionExceptionsRequest.
func (v GetAttackGroupsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateAttackGroupConditionExceptionRequest.
func (v UpdateAttackGroupRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetAttackGroup(ctx context.Context, params GetAttackGroupRequest) (*GetAttackGroupResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAttackGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAttackGroup request: %w", err)
	}

	var result GetAttackGroupResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get attack group request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetAttackGroups(ctx context.Context, params GetAttackGroupsRequest) (*GetAttackGroupsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAttackGroups")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAttackGroups request: %w", err)
	}

	var result GetAttackGroupsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get attack groups request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Group != "" {
		var filteredResult GetAttackGroupsResponse
		for k, val := range result.AttackGroups {
			if val.Group == params.Group {
				filteredResult.AttackGroups = append(filteredResult.AttackGroups, result.AttackGroups[k])
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateAttackGroup(ctx context.Context, params UpdateAttackGroupRequest) (*UpdateAttackGroupResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAttackGroup")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAttackGroup request: %w", err)
	}

	var result UpdateAttackGroupResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update attack group request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
