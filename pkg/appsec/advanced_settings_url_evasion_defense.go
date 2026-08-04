package appsec

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// AdvancedSettingsURLEvasionDefense supports retrieving or modifying the URL Evasion Defense settings
	AdvancedSettingsURLEvasionDefense interface {
		// GetAdvancedSettingsURLEvasionDefense retrieves the URL Evasion Defense settings
		GetAdvancedSettingsURLEvasionDefense(ctx context.Context, params GetAdvancedSettingsURLEvasionDefenseRequest) (*GetAdvancedSettingsURLEvasionDefenseResponse, error)

		// UpdateAdvancedSettingsURLEvasionDefense modifies the URL Evasion Defense settings
		UpdateAdvancedSettingsURLEvasionDefense(ctx context.Context, params UpdateAdvancedSettingsURLEvasionDefenseRequest) (*UpdateAdvancedSettingsURLEvasionDefenseResponse, error)
	}

	// GetAdvancedSettingsURLEvasionDefenseRequest is used to retrieve the URL Evasion Defense settings
	GetAdvancedSettingsURLEvasionDefenseRequest struct {
		ConfigID int64
		Version  int64
	}

	// AdvancedSettingsURLEvasionDefenseStatus is a URL Evasion Defense status
	AdvancedSettingsURLEvasionDefenseStatus string
	// AdvancedSettingsURLEvasionDefenseConditionOperator is a URL Evasion Defense conditions operator
	AdvancedSettingsURLEvasionDefenseConditionOperator string

	// GetAdvancedSettingsURLEvasionDefenseResponse returns the URL Evasion Defense settings
	GetAdvancedSettingsURLEvasionDefenseResponse struct {
		Status      string                                          `json:"status"`
		LockVersion int64                                           `json:"lockVersion"`
		BypassLists []string                                        `json:"bypassLists"`
		Rules       []AdvancedSettingsURLEvasionDefenseRuleResponse `json:"rules"`
	}
	// UpdateAdvancedSettingsURLEvasionDefenseRequest is used to update the URL Evasion Defense settings
	UpdateAdvancedSettingsURLEvasionDefenseRequest struct {
		ConfigID                      int64
		Version                       int64
		DisableLegacyEvasivePathMatch bool
		Body                          UpdateAdvancedSettingsURLEvasionDefenseRequestBody
	}

	// UpdateAdvancedSettingsURLEvasionDefenseRequestBody is the JSON payload for updating URL Evasion Defense settings.
	// Rules cannot be shared with GetAdvancedSettingsURLEvasionDefenseResponse because the Action field
	// differs between request (string, marshalled always) and response (*string, omitted when nil)
	UpdateAdvancedSettingsURLEvasionDefenseRequestBody struct {
		Status      AdvancedSettingsURLEvasionDefenseStatus        `json:"status"`
		LockVersion int64                                          `json:"lockVersion"`
		BypassLists []string                                       `json:"bypassLists,omitempty"`
		Rules       []AdvancedSettingsURLEvasionDefenseRuleRequest `json:"rules,omitempty"`
	}

	// UpdateAdvancedSettingsURLEvasionDefenseResponse returns the result of updating the URL Evasion Defense settings
	UpdateAdvancedSettingsURLEvasionDefenseResponse GetAdvancedSettingsURLEvasionDefenseResponse

	// AdvancedSettingsURLEvasionDefenseRuleRequest defines a URL Evasion Defense rule payload
	// Action cannot be shared with RuleResponse because it is a plain string in requests
	// but a *string in responses (omitted when nil)
	AdvancedSettingsURLEvasionDefenseRuleRequest struct {
		RuleID            int64                                               `json:"ruleId"`
		ConditionOperator *AdvancedSettingsURLEvasionDefenseConditionOperator `json:"conditionOperator,omitempty"`
		Conditions        []AdvancedSettingsURLEvasionDefenseRuleCondition    `json:"conditions,omitempty"`
		Action            string                                              `json:"action"`
	}

	// AdvancedSettingsURLEvasionDefenseRuleResponse defines a URL Evasion Defense rule response
	AdvancedSettingsURLEvasionDefenseRuleResponse struct {
		RuleID            int64                                            `json:"ruleId"`
		ConditionOperator *string                                          `json:"conditionOperator"`
		Conditions        []AdvancedSettingsURLEvasionDefenseRuleCondition `json:"conditions"`
		Action            *string                                          `json:"action"`
		Name              string                                           `json:"name"`
		Description       string                                           `json:"description"`
		DefaultAction     string                                           `json:"defaultAction"`
	}
	// AdvancedSettingsURLEvasionDefenseRuleCondition defines a URL Evasion Defense rule condition
	AdvancedSettingsURLEvasionDefenseRuleCondition struct {
		Type          string   `json:"type,omitempty"`
		Extensions    []string `json:"extensions,omitempty"`
		Filenames     []string `json:"filenames,omitempty"`
		Hosts         []string `json:"hosts,omitempty"`
		IPs           []string `json:"ips,omitempty"`
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
)

const (
	// AdvancedSettingsURLEvasionDefenseStatusEnabled enables URL Evasion Defense
	AdvancedSettingsURLEvasionDefenseStatusEnabled AdvancedSettingsURLEvasionDefenseStatus = "enabled"
	// AdvancedSettingsURLEvasionDefenseStatusDisabled disables URL Evasion Defense
	AdvancedSettingsURLEvasionDefenseStatusDisabled AdvancedSettingsURLEvasionDefenseStatus = "disabled"

	// AdvancedSettingsURLEvasionDefenseConditionOperatorAnd requires all rule conditions to match
	AdvancedSettingsURLEvasionDefenseConditionOperatorAnd AdvancedSettingsURLEvasionDefenseConditionOperator = "AND"
	// AdvancedSettingsURLEvasionDefenseConditionOperatorOr requires at least one rule condition to match
	AdvancedSettingsURLEvasionDefenseConditionOperatorOr AdvancedSettingsURLEvasionDefenseConditionOperator = "OR"
)

// Validate validates GetAdvancedSettingsURLEvasionDefenseRequest
func (v GetAdvancedSettingsURLEvasionDefenseRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	})
}

// Validate validates AdvancedSettingsURLEvasionDefenseRuleRequest
func (v AdvancedSettingsURLEvasionDefenseRuleRequest) Validate() error {
	errors := validation.Errors{
		"RuleID":            validation.Validate(v.RuleID, validation.Required),
		"Action":            validation.Validate(v.Action, validation.Required),
		"ConditionOperator": validation.Validate(v.ConditionOperator),
		"Conditions":        validation.Validate(v.Conditions),
	}

	return errors.Filter()
}

// Validate validates UpdateAdvancedSettingsURLEvasionDefenseRequest
func (v UpdateAdvancedSettingsURLEvasionDefenseRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"Body":     validation.Validate(v.Body),
	})
}

// Validate validates UpdateAdvancedSettingsURLEvasionDefenseRequestBody
func (v UpdateAdvancedSettingsURLEvasionDefenseRequestBody) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Status": validation.Validate(v.Status),
		"Rules":  validation.Validate(v.Rules),
	})
}

// Validate validates AdvancedSettingsURLEvasionDefenseRuleCondition
func (c AdvancedSettingsURLEvasionDefenseRuleCondition) Validate() error {
	allowedFieldsByType := map[string][]string{
		"extensionMatch":     {"Extensions"},
		"filenameMatch":      {"Filenames"},
		"hostMatch":          {"Hosts"},
		"pathMatch":          {"Paths"},
		"requestMethodMatch": {"Methods"},
		"ipMatch":            {"Ips", "UseHeaders"},
		"clientListMatch":    {"ClientLists", "UseHeaders"},
		"requestHeaderMatch": {"Header", "Value", "ValueCase", "ValueWildcard"},
		"uriQueryMatch":      {"Name", "Value", "NameCaseSensitive", "CaseSensitive", "Wildcard"},
	}

	typeIn := func(condType string, allowedTypes ...string) validation.RuleFunc {
		return func(_ interface{}) error {
			if _, knownType := allowedFieldsByType[condType]; !knownType {
				return nil
			}
			if slices.Contains(allowedTypes, condType) {
				return nil
			}
			allowedFields := append(slices.Clone(allowedFieldsByType[condType]), "Type", "PositiveMatch")
			return fmt.Errorf("field not supported for type '%s'; allowed fields: %s", condType, fmt.Sprintf("%v", allowedFields))
		}
	}

	errors := validation.Errors{
		"Type": validation.Validate(c.Type, validation.Required,
			validation.In("clientListMatch", "extensionMatch", "filenameMatch", "hostMatch", "ipMatch", "pathMatch", "requestHeaderMatch", "requestMethodMatch", "uriQueryMatch").
				Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'clientListMatch', 'extensionMatch', 'filenameMatch', 'hostMatch', 'ipMatch', 'pathMatch', 'requestHeaderMatch', 'requestMethodMatch', 'uriQueryMatch'", c.Type)),
		),
		"Extensions":        validation.Validate(c.Extensions, validation.When(len(c.Extensions) > 0, validation.By(typeIn(c.Type, "extensionMatch")))),
		"Filenames":         validation.Validate(c.Filenames, validation.When(len(c.Filenames) > 0, validation.By(typeIn(c.Type, "filenameMatch")))),
		"Hosts":             validation.Validate(c.Hosts, validation.When(len(c.Hosts) > 0, validation.By(typeIn(c.Type, "hostMatch")))),
		"Ips":               validation.Validate(c.IPs, validation.When(len(c.IPs) > 0, validation.By(typeIn(c.Type, "ipMatch")))),
		"Methods":           validation.Validate(c.Methods, validation.When(len(c.Methods) > 0, validation.By(typeIn(c.Type, "requestMethodMatch")))),
		"Paths":             validation.Validate(c.Paths, validation.When(len(c.Paths) > 0, validation.By(typeIn(c.Type, "pathMatch")))),
		"ClientLists":       validation.Validate(c.ClientLists, validation.When(len(c.ClientLists) > 0, validation.By(typeIn(c.Type, "clientListMatch")))),
		"Header":            validation.Validate(c.Header, validation.When(c.Header != "", validation.By(typeIn(c.Type, "requestHeaderMatch")))),
		"Name":              validation.Validate(c.Name, validation.When(c.Name != "", validation.By(typeIn(c.Type, "uriQueryMatch")))),
		"Value":             validation.Validate(c.Value, validation.When(c.Value != "", validation.By(typeIn(c.Type, "requestHeaderMatch", "uriQueryMatch")))),
		"CaseSensitive":     validation.Validate(c.CaseSensitive, validation.When(c.CaseSensitive, validation.By(typeIn(c.Type, "uriQueryMatch")))),
		"NameCaseSensitive": validation.Validate(c.NameCase, validation.When(c.NameCase, validation.By(typeIn(c.Type, "uriQueryMatch")))),
		"Wildcard":          validation.Validate(c.Wildcard, validation.When(c.Wildcard, validation.By(typeIn(c.Type, "uriQueryMatch")))),
		"ValueCase":         validation.Validate(c.ValueCase, validation.When(c.ValueCase, validation.By(typeIn(c.Type, "requestHeaderMatch")))),
		"ValueWildcard":     validation.Validate(c.ValueWildcard, validation.When(c.ValueWildcard, validation.By(typeIn(c.Type, "requestHeaderMatch")))),
		"UseHeaders":        validation.Validate(c.UseHeaders, validation.When(c.UseHeaders, validation.By(typeIn(c.Type, "ipMatch", "clientListMatch")))),
	}

	return errors.Filter()
}

// Validate validates AdvancedSettingsURLEvasionDefenseStatus
func (s AdvancedSettingsURLEvasionDefenseStatus) Validate() error {
	return validation.In(AdvancedSettingsURLEvasionDefenseStatusEnabled, AdvancedSettingsURLEvasionDefenseStatusDisabled).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s'", s, AdvancedSettingsURLEvasionDefenseStatusEnabled, AdvancedSettingsURLEvasionDefenseStatusDisabled)).
		Validate(s)
}

// Validate validates AdvancedSettingsURLEvasionDefenseConditionOperator
func (o AdvancedSettingsURLEvasionDefenseConditionOperator) Validate() error {
	return validation.In(AdvancedSettingsURLEvasionDefenseConditionOperatorAnd, AdvancedSettingsURLEvasionDefenseConditionOperatorOr).
		Error(fmt.Sprintf("value '%v' is invalid. Must be one of: '%s', '%s'", o, AdvancedSettingsURLEvasionDefenseConditionOperatorAnd, AdvancedSettingsURLEvasionDefenseConditionOperatorOr)).
		Validate(o)
}

func (p *appsec) GetAdvancedSettingsURLEvasionDefense(ctx context.Context, params GetAdvancedSettingsURLEvasionDefenseRequest) (*GetAdvancedSettingsURLEvasionDefenseResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsURLEvasionDefense")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewGet(ctx, "/appsec/v1/configs/%d/versions/%d/advanced-settings/url-evasion-defense", params.ConfigID, params.Version).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsURLEvasionDefense request: %w", err)
	}

	var result GetAdvancedSettingsURLEvasionDefenseResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings URL Evasion Defense request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAdvancedSettingsURLEvasionDefense(ctx context.Context, params UpdateAdvancedSettingsURLEvasionDefenseRequest) (*UpdateAdvancedSettingsURLEvasionDefenseResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsURLEvasionDefense")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	req, err := request.NewPut(ctx, "/appsec/v1/configs/%d/versions/%d/advanced-settings/url-evasion-defense", params.ConfigID, params.Version).
		AddQueryParam("disableLegacyEvasivePathMatch", strconv.FormatBool(params.DisableLegacyEvasivePathMatch)).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsURLEvasionDefense request: %w", err)
	}

	var result UpdateAdvancedSettingsURLEvasionDefenseResponse
	resp, err := p.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings URL Evasion Defense request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
