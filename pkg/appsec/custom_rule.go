package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomRule interface supports creating, retrievinfg, modifying and removing custom rules
	// for a configuration.
	CustomRule interface {
		// GetCustomRules lists custom rules defined in a security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-configs-custom-rules
		GetCustomRules(ctx context.Context, params GetCustomRulesRequest) (*GetCustomRulesResponse, error)

		// GetCustomRule returns the details of a custom rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-config-custom-rule
		GetCustomRule(ctx context.Context, params GetCustomRuleRequest) (*GetCustomRuleResponse, error)

		// CreateCustomRule creates a new custom rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-config-custom-rules
		CreateCustomRule(ctx context.Context, params CreateCustomRuleRequest) (*CreateCustomRuleResponse, error)

		// UpdateCustomRule updates an existing custom rule.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-config-custom-rule
		UpdateCustomRule(ctx context.Context, params UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error)

		// RemoveCustomRule deletes a custom rule as long as it isn't activated.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-config-custom-rule
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
			ID                  int                        `json:"id"`
			Link                string                     `json:"link"`
			Name                string                     `json:"name"`
			Status              string                     `json:"status"`
			Version             int                        `json:"version"`
			EffectiveTimePeriod *CustomRuleEffectivePeriod `json:"effectiveTimePeriod,omitempty"`
			SamplingRate        *int                       `json:"samplingRate,omitempty"`
		} `json:"customRules"`
	}

	// GetCustomRuleRequest is used to retrieve the details of a custom rule.
	GetCustomRuleRequest struct {
		ConfigID int `json:"configid,omitempty"`
		ID       int `json:"id,omitempty"`
	}

	// GetCustomRuleResponse is returned from a call to GetCustomRule.
	GetCustomRuleResponse CustomRuleResponse

	// CustomRuleResponse is returned from calls to GetCustomRule, UpdateCustomRule, and DeleteCustomRule.
	CustomRuleResponse struct {
		ID            int      `json:"-"`
		Name          string   `json:"name"`
		Description   string   `json:"description,omitempty"`
		Version       int      `json:"-"`
		RuleActivated bool     `json:"-"`
		Structured    bool     `json:"-"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Name                  *json.RawMessage `json:"name,omitempty"`
			NameCase              *bool            `json:"nameCase,omitempty"`
			NameWildcard          *bool            `json:"nameWildcard,omitempty"`
			PositiveMatch         bool             `json:"positiveMatch"`
			Type                  string           `json:"type"`
			Value                 *json.RawMessage `json:"value,omitempty"`
			ValueCase             *bool            `json:"valueCase,omitempty"`
			ValueExactMatch       *bool            `json:"valueExactMatch,omitempty"`
			ValueIgnoreSegment    *bool            `json:"valueIgnoreSegment,omitempty"`
			ValueNormalize        *bool            `json:"valueNormalize,omitempty"`
			ValueRecursive        *bool            `json:"valueRecursive,omitempty"`
			ValueWildcard         *bool            `json:"valueWildcard,omitempty"`
			UseXForwardForHeaders *bool            `json:"useXForwardForHeaders,omitempty"`
		} `json:"conditions"`
		EffectiveTimePeriod *CustomRuleEffectivePeriod `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int                        `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage           `json:"loggingOptions,omitempty"`
		Operation           string                     `json:"operation,omitempty"`
	}

	// CustomRuleEffectivePeriod defines the period during which a custom rule is active as well as its current status.
	CustomRuleEffectivePeriod struct {
		EndDate   string `json:"endDate"`
		StartDate string `json:"startDate"`
		Status    string `json:"-"`
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
		RuleActivated bool     `json:"-"`
		Structured    bool     `json:"-"`
		Tag           []string `json:"tag"`
		Conditions    []struct {
			Name                  *json.RawMessage `json:"name,omitempty"`
			NameCase              *bool            `json:"nameCase,omitempty"`
			NameWildcard          *bool            `json:"nameWildcard,omitempty"`
			PositiveMatch         bool             `json:"positiveMatch"`
			Type                  string           `json:"type"`
			Value                 *json.RawMessage `json:"value,omitempty"`
			ValueCase             *bool            `json:"valueCase,omitempty"`
			ValueExactMatch       *bool            `json:"valueExactMatch,omitempty"`
			ValueIgnoreSegment    *bool            `json:"valueIgnoreSegment,omitempty"`
			ValueNormalize        *bool            `json:"valueNormalize,omitempty"`
			ValueRecursive        *bool            `json:"valueRecursive,omitempty"`
			ValueWildcard         *bool            `json:"valueWildcard,omitempty"`
			UseXForwardForHeaders *bool            `json:"useXForwardForHeaders,omitempty"`
		} `json:"conditions"`
		EffectiveTimePeriod *CustomRuleEffectivePeriod `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int                        `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage           `json:"loggingOptions,omitempty"`
		Operation           string                     `json:"operation,omitempty"`
	}

	// UpdateCustomRuleRequest is used to modify an existing custom rule.
	UpdateCustomRuleRequest struct {
		ConfigID       int             `json:"configid,omitempty"`
		ID             int             `json:"id,omitempty"`
		Version        int             `json:"version,omitempty"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateCustomRuleResponse is returned from a call to UpdateCustomRule.
	UpdateCustomRuleResponse GetCustomRuleResponse

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
	logger := p.Log(ctx)
	logger.Debug("GetCustomRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules/%d",
		params.ConfigID,
		params.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRule request: %w", err)
	}

	var result GetCustomRuleResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get custom rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetCustomRules(ctx context.Context, params GetCustomRulesRequest) (*GetCustomRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCustomRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomRules request: %w", err)
	}

	var result GetCustomRulesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get custom rules request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ID != 0 {
		var filteredResult GetCustomRulesResponse
		for _, val := range result.CustomRules {
			if val.ID == params.ID {
				filteredResult.CustomRules = append(filteredResult.CustomRules, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateCustomRule(ctx context.Context, params UpdateCustomRuleRequest) (*UpdateCustomRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateCustomRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules/%d",
		params.ConfigID,
		params.ID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomRule request: %w", err)
	}

	var result UpdateCustomRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update custom rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateCustomRule(ctx context.Context, params CreateCustomRuleRequest) (*CreateCustomRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateCustomRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/custom-rules",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomRule request: %w", err)
	}

	var result CreateCustomRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create custom rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveCustomRule(ctx context.Context, params RemoveCustomRuleRequest) (*RemoveCustomRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveCustomRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result RemoveCustomRuleResponse
	uri := fmt.Sprintf("/appsec/v1/configs/%d/custom-rules/%d", params.ConfigID, params.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveCustomRule request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return nil, fmt.Errorf("remove custom rule request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &result, nil
}
