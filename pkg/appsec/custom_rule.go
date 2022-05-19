package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomRule interface supports creating, retrievinfg, modifying and removing custom rules
	// for a configuration.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#customrule
	CustomRule interface {
		//https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomrules
		GetCustomRules(ctx context.Context, params GetCustomRulesRequest) (*GetCustomRulesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getruleid
		GetCustomRule(ctx context.Context, params GetCustomRuleRequest) (*GetCustomRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postcustomrules
		CreateCustomRule(ctx context.Context, params CreateCustomRuleRequest) (*CreateCustomRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleid
		UpdateCustomRule(ctx context.Context, params UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteruleid
		RemoveCustomRule(ctx context.Context, params RemoveCustomRuleRequest) (*RemoveCustomRuleResponse, error)
	}

	// CustomRuleConditionsValue is a slice of strings used to indicate condition values in custom rule conditions.
	CustomRuleConditionsValue []string

	// CustomRuleConditionsName is a slice of strings used to indicate condition names in custom rule conditions.
	CustomRuleConditionsName []string

	// GetCustomRulesRequest is used to retrieve the custom rules for a configuration.
	GetCustomRulesRequest struct {
		ConfigID int `json:"configid,omitempty"`
		ID       int `json:"-"`
	}

	// GetCustomRulesResponse is returned from a call to GetCustomRules.
	GetCustomRulesResponse struct {
		CustomRules []struct {
			ID                  int              `json:"id"`
			Link                string           `json:"link"`
			Name                string           `json:"name"`
			Status              string           `json:"status"`
			Version             int              `json:"version"`
			EffectiveTimePeriod *json.RawMessage `json:"effectiveTimePeriod,omitempty"`
			SamplingRate        int              `json:"samplingRate,omitempty"`
		} `json:"customRules"`
	}

	// GetCustomRuleRequest is used to retrieve the details of a custom rule.
	GetCustomRuleRequest struct {
		ConfigID int `json:"configid,omitempty"`
		ID       int `json:"id,omitempty"`
	}

	// GetCustomRuleResponse is returned from a call to GetCustomRule.
	GetCustomRuleResponse struct {
		ID            int      `json:"-"`
		Name          string   `json:"name"`
		Description   string   `json:"description,omitempty"`
		Version       int      `json:"-"`
		RuleActivated bool     `json:"-"`
		Structured    bool     `json:"-"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type           string           `json:"type"`
			PositiveMatch  bool             `json:"positiveMatch"`
			Value          *json.RawMessage `json:"value,omitempty"`
			ValueWildcard  *json.RawMessage `json:"valueWildcard,omitempty"`
			ValueCase      *json.RawMessage `json:"valueCase,omitempty"`
			NameWildcard   *json.RawMessage `json:"nameWildcard,omitempty"`
			Name           *json.RawMessage `json:"name,omitempty"`
			NameCase       *json.RawMessage `json:"nameCase,omitempty"`
			ValueNormalize *json.RawMessage `json:"valueNormalize,omitempty"`
		} `json:"conditions"`
		EffectiveTimePeriod *json.RawMessage `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int              `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage `json:"loggingOptions,omitempty"`
		Operation           string           `json:"operation,omitempty"`
	}

	// CreateCustomRuleRequest is used to create a custom rule.
	CreateCustomRuleRequest struct {
		ConfigID       int             `json:"configid,omitempty"`
		Version        int             `json:"version,omitempty"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// CreateCustomRuleResponse is returned from a call to CreateCustomRule.
	CreateCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description,omitempty"`
		Version       int      `json:"-"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type           string                    `json:"type"`
			PositiveMatch  bool                      `json:"positiveMatch"`
			Value          CustomRuleConditionsValue `json:"value,omitempty"`
			ValueWildcard  *json.RawMessage          `json:"valueWildcard,omitempty"`
			ValueCase      *json.RawMessage          `json:"valueCase,omitempty"`
			NameWildcard   *json.RawMessage          `json:"nameWildcard,omitempty"`
			Name           CustomRuleConditionsName  `json:"name,omitempty"`
			NameCase       *json.RawMessage          `json:"nameCase,omitempty"`
			ValueNormalize *json.RawMessage          `json:"valueNormalize,omitempty"`
		} `json:"conditions"`
		EffectiveTimePeriod *json.RawMessage `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int              `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage `json:"loggingOptions,omitempty"`
		Operation           string           `json:"operation,omitempty"`
	}

	// UpdateCustomRuleRequest is used to modify an existing custom rule.
	UpdateCustomRuleRequest struct {
		ConfigID       int             `json:"configid,omitempty"`
		ID             int             `json:"id,omitempty"`
		Version        int             `json:"version,omitempty"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateCustomRuleResponse is returned from a call to UpdateCustomRule.
	UpdateCustomRuleResponse struct {
		ID            int      `json:"id,omitempty"`
		Name          string   `json:"name"`
		Description   string   `json:"description,omitempty"`
		Version       int      `json:"-"`
		RuleActivated bool     `json:"ruleActivated"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Type           string           `json:"type"`
			PositiveMatch  bool             `json:"positiveMatch"`
			Value          *json.RawMessage `json:"value,omitempty"`
			ValueWildcard  *json.RawMessage `json:"valueWildcard,omitempty"`
			ValueCase      *json.RawMessage `json:"valueCase,omitempty"`
			NameWildcard   *json.RawMessage `json:"nameWildcard,omitempty"`
			Name           *json.RawMessage `json:"name,omitempty"`
			NameCase       *json.RawMessage `json:"nameCase,omitempty"`
			ValueNormalize *json.RawMessage `json:"valueNormalize,omitempty"`
		} `json:"conditions"`
		EffectiveTimePeriod *json.RawMessage `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int              `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage `json:"loggingOptions,omitempty"`
		Operation           string           `json:"operation,omitempty"`
	}

	// RemoveCustomRuleRequest is used to remove a custom rule.
	RemoveCustomRuleRequest struct {
		ConfigID int `json:"configid,omitempty"`
		ID       int `json:"id,omitempty"`
	}

	// RemoveCustomRuleResponse is returned from a call to RemoveCustomRule.
	RemoveCustomRuleResponse UpdateCustomRuleResponse
)

// UnmarshalJSON reads a CustomRuleConditionsValue from its data argument.
func (c *CustomRuleConditionsValue) UnmarshalJSON(data []byte) error {
	var nums interface{}
	err := json.Unmarshal(data, &nums)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(nums)
	switch items.Kind() {
	case reflect.String:
		*c = append(*c, items.String())

	case reflect.Slice:
		*c = make(CustomRuleConditionsValue, 0, items.Len())
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				*c = append(*c, item.String())
			case reflect.Interface:
				*c = append(*c, item.Interface().(string))
			}
		}
	}
	return nil
}

// UnmarshalJSON reads a CustomRuleConditionsName from its data argument.
func (c *CustomRuleConditionsName) UnmarshalJSON(data []byte) error {
	var nums interface{}
	err := json.Unmarshal(data, &nums)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(nums)
	switch items.Kind() {
	case reflect.String:
		*c = append(*c, items.String())

	case reflect.Slice:
		*c = make(CustomRuleConditionsName, 0, items.Len())
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				*c = append(*c, item.String())
			case reflect.Interface:
				*c = append(*c, item.Interface().(string))
			}
		}
	}
	return nil
}

// Validate validates a GetCustomRuleRequest.
func (v GetCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomRulesRequest.
func (v GetCustomRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomRuleRequest.
func (v CreateCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomRuleRequest.
func (v UpdateCustomRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomRuleRequest.
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
		return nil, fmt.Errorf("failed to create GetCustomRule request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetCustomRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
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
	var rvalfiltered GetCustomRulesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRules request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetCustomRules request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ID != 0 {
		for _, val := range rval.CustomRules {
			if val.ID == params.ID {
				rvalfiltered.CustomRules = append(rvalfiltered.CustomRules, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateCustomRule request: %w", err)
	}

	var rval UpdateCustomRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

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
		return nil, fmt.Errorf("failed to create CreateCustomRule request: %w", err)
	}

	var rval CreateCustomRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("CreateCustomRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveCustomRule request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("RemoveCustomRule request failed: %w", err)
	}
	logger.Debugf("RemoveCustomRule RESP CODE %v", resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
