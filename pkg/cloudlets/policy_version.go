package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PolicyVersions is a cloudlets policy versions API interface
	PolicyVersions interface {
		// ListPolicyVersions lists policy versions by policyID
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyversions
		ListPolicyVersions(context.Context, ListPolicyVersionsRequest) ([]PolicyVersion, error)

		// GetPolicyVersion gets policy version by policyID and version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicyversion
		GetPolicyVersion(context.Context, GetPolicyVersionRequest) (*PolicyVersion, error)

		// CreatePolicyVersion creates policy version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postpolicyversions
		CreatePolicyVersion(context.Context, CreatePolicyVersionRequest) (*PolicyVersion, error)

		// DeletePolicyVersion deletes policy version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#deletepolicyversion
		DeletePolicyVersion(context.Context, DeletePolicyVersionRequest) error

		// UpdatePolicyVersion updates policy version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#putpolicyversion
		UpdatePolicyVersion(context.Context, UpdatePolicyVersionRequest) (*PolicyVersion, error)
	}

	// PolicyVersion is response returned by GetPolicyVersion, CreatePolicyVersion or UpdatePolicyVersion
	PolicyVersion struct {
		Location         string          `json:"location"`
		RevisionID       int64           `json:"revisionId"`
		PolicyID         int64           `json:"policyId"`
		Version          int64           `json:"version"`
		Description      string          `json:"description"`
		CreatedBy        string          `json:"createdBy"`
		CreateDate       int             `json:"createDate"`
		LastModifiedBy   string          `json:"lastModifiedBy"`
		LastModifiedDate int             `json:"lastModifiedDate"`
		RulesLocked      bool            `json:"rulesLocked"`
		Activations      []*Activation   `json:"activations"`
		MatchRules       MatchRules      `json:"matchRules"`
		MatchRuleFormat  MatchRuleFormat `json:"matchRuleFormat"`
		Deleted          bool            `json:"deleted,omitempty"`
		Warnings         []Warning       `json:"warnings,omitempty"`
	}

	// ListPolicyVersionsRequest describes the parameters needed to list policy versions
	ListPolicyVersionsRequest struct {
		PolicyID           int64
		IncludeRules       bool
		IncludeDeleted     bool
		IncludeActivations bool
		Offset             int
		PageSize           *int
	}

	// GetPolicyVersionRequest describes the parameters needed to get policy version
	GetPolicyVersionRequest struct {
		PolicyID  int64
		Version   int64
		OmitRules bool
	}

	// CreatePolicyVersionRequest describes the body of the create policy request
	CreatePolicyVersionRequest struct {
		CreatePolicyVersion
		PolicyID int64
	}

	// CreatePolicyVersion describes the body of the create policy request
	CreatePolicyVersion struct {
		Description     string          `json:"description,omitempty"`
		MatchRuleFormat MatchRuleFormat `json:"matchRuleFormat,omitempty"`
		MatchRules      MatchRules      `json:"matchRules"`
	}

	// UpdatePolicyVersion describes the body of the update policy version request
	UpdatePolicyVersion struct {
		Description     string          `json:"description,omitempty"`
		MatchRuleFormat MatchRuleFormat `json:"matchRuleFormat,omitempty"`
		MatchRules      MatchRules      `json:"matchRules"`
		Deleted         bool            `json:"deleted"`
	}

	// UpdatePolicyVersionRequest describes the parameters of the update policy version request
	UpdatePolicyVersionRequest struct {
		UpdatePolicyVersion
		PolicyID int64
		Version  int64
	}

	// DeletePolicyVersionRequest describes the parameters of the delete policy version request
	DeletePolicyVersionRequest struct {
		PolicyID int64
		Version  int64
	}

	// MatchRule is base interface for MarchRuleALB and MatchRuleER
	MatchRule interface {
		// cloudletType is a private method to ensure that only match rules for supported cloudlets can be used
		cloudletType() string
		Validate() error
	}

	// MatchRules is an array of *MarchRuleALB or *MatchRuleER depending on the cloudletId (9 or 0) of the policy
	MatchRules []MatchRule

	// MatchRuleALB represents a match rule resource for create or update resource
	MatchRuleALB struct {
		Name            string             `json:"name,omitempty"`
		Type            MatchRuleType      `json:"type,omitempty"`
		Start           int                `json:"start,omitempty"`
		End             int                `json:"end,omitempty"`
		ID              int64              `json:"id,omitempty"`
		Matches         []MatchCriteriaALB `json:"matches,omitempty"`
		AkaRuleID       string             `json:"akaRuleId,omitempty"`
		MatchURL        string             `json:"matchURL,omitempty"`
		Location        string             `json:"location,omitempty"`
		MatchesAlways   bool               `json:"matchesAlways"`
		ForwardSettings ForwardSettings    `json:"forwardSettings"`
	}

	// ForwardSettings represents forward settings
	ForwardSettings struct {
		OriginID string `json:"originId"`
	}

	// MatchRuleER represents a match rule resource for create or update resource
	MatchRuleER struct {
		Name                     string            `json:"name,omitempty"`
		Type                     MatchRuleType     `json:"type,omitempty"`
		Start                    int               `json:"start,omitempty"`
		End                      int               `json:"end,omitempty"`
		ID                       int64             `json:"id,omitempty"`
		Matches                  []MatchCriteriaER `json:"matches,omitempty"`
		AkaRuleID                string            `json:"akaRuleId,omitempty"`
		UseRelativeURL           string            `json:"useRelativeUrl"`
		StatusCode               int               `json:"statusCode"`
		RedirectURL              string            `json:"redirectURL"`
		MatchURL                 string            `json:"matchURL,omitempty"`
		Location                 string            `json:"location,omitempty"`
		UseIncomingQueryString   bool              `json:"useIncomingQueryString"`
		UseIncomingSchemeAndHost bool              `json:"useIncomingSchemeAndHost"`
	}

	// MatchCriteriaALB represents a match criteria resource for match rule for cloudlet ALB
	// ObjectMatchValue can contain ObjectMatchValueObjectSubtype or ObjectMatchValueRangeOrSimpleSubtype
	MatchCriteriaALB struct {
		MatchType        string        `json:"matchType,omitempty"`
		MatchValue       string        `json:"matchValue,omitempty"`
		MatchOperator    MatchOperator `json:"matchOperator,omitempty"`
		CaseSensitive    bool          `json:"caseSensitive"`
		Negate           bool          `json:"negate"`
		CheckIPs         CheckIPs      `json:"checkIPs,omitempty"`
		ObjectMatchValue interface{}   `json:"objectMatchValue,omitempty"`
	}

	// MatchCriteriaER represents a match criteria resource for match rule for cloudlet ER
	MatchCriteriaER struct {
		MatchType     string        `json:"matchType,omitempty"`
		MatchValue    string        `json:"matchValue,omitempty"`
		MatchOperator MatchOperator `json:"matchOperator,omitempty"`
		CaseSensitive bool          `json:"caseSensitive"`
		Negate        bool          `json:"negate"`
		CheckIPs      CheckIPs      `json:"checkIPs,omitempty"`
	}

	// ObjectMatchValueObjectSubtype represents an object match value resource for match criteria
	ObjectMatchValueObjectSubtype struct {
		Name              string                            `json:"name"`
		Type              ObjectMatchValueObjectTypeSubtype `json:"type"`
		NameCaseSensitive bool                              `json:"nameCaseSensitive"`
		NameHasWildcard   bool                              `json:"nameHasWildcard"`
		Options           *Options                          `json:"options,omitempty"`
	}

	// ObjectMatchValueRangeOrSimpleSubtype represents a match value resource for match criteria
	ObjectMatchValueRangeOrSimpleSubtype struct {
		Type  ObjectMatchValueRangeOrSimpleTypeSubtype `json:"type"`
		Value interface{}                              `json:"value,omitempty"`
	}

	// Options represents an option resource for ObjectMatchValueObjectSubtype
	Options struct {
		Value              interface{} `json:"value,omitempty"`
		ValueHasWildcard   bool        `json:"valueHasWildcard,omitempty"`
		ValueCaseSensitive bool        `json:"valueCaseSensitive,omitempty"`
		ValueEscaped       bool        `json:"valueEscaped,omitempty"`
	}

	//MatchRuleType enum type
	MatchRuleType string
	// MatchRuleFormat enum type
	MatchRuleFormat string
	// MatchOperator enum type
	MatchOperator string
	// CheckIPs enum type
	CheckIPs string
	// ObjectMatchValueRangeOrSimpleTypeSubtype enum type
	ObjectMatchValueRangeOrSimpleTypeSubtype string
	// ObjectMatchValueObjectTypeSubtype enum type
	ObjectMatchValueObjectTypeSubtype string
)

const (
	// MatchRuleTypeALB represents rule type for ALB cloudlets
	MatchRuleTypeALB MatchRuleType = "albMatchRule"
	// MatchRuleTypeER represents rule type for ER cloudlets
	MatchRuleTypeER MatchRuleType = "erMatchRule"
)

const (
	// MatchRuleFormat10 represents default match rule format
	MatchRuleFormat10 MatchRuleFormat = "1.0"
	// MatchRuleFormatDefault represents default match rule format
	MatchRuleFormatDefault = MatchRuleFormat10
)

const (
	// MatchOperatorContains represents contains operator
	MatchOperatorContains MatchOperator = "contains"
	// MatchOperatorExists represents exists operator
	MatchOperatorExists MatchOperator = "exists"
	// MatchOperatorEquals represents equals operator
	MatchOperatorEquals MatchOperator = "equals"
)

const (
	// CheckIPsConnectingIP represents connecting ip option
	CheckIPsConnectingIP CheckIPs = "CONNECTING_IP"
	// CheckIPsXFFHeaders represents xff headers option
	CheckIPsXFFHeaders CheckIPs = "XFF_HEADERS"
	// CheckIPsConnectingIPXFFHeaders represents connecting ip + xff headers option
	CheckIPsConnectingIPXFFHeaders CheckIPs = "CONNECTING_IP XFF_HEADERS"
)

const (
	// ObjectMatchValueRangeOrSimpleTypeSubtypeRange represents range option
	ObjectMatchValueRangeOrSimpleTypeSubtypeRange ObjectMatchValueRangeOrSimpleTypeSubtype = "range"
	// ObjectMatchValueRangeOrSimpleTypeSubtypeSimple represents simple option
	ObjectMatchValueRangeOrSimpleTypeSubtypeSimple ObjectMatchValueRangeOrSimpleTypeSubtype = "simple"
)

const (
	// ObjectMatchValueObjectTypeSubtypeObject represents object option
	ObjectMatchValueObjectTypeSubtypeObject ObjectMatchValueObjectTypeSubtype = "object"
)

// Validate validates ListPolicyVersionsRequest
func (c ListPolicyVersionsRequest) Validate() error {
	return validation.Errors{
		"PolicyID": validation.Validate(c.PolicyID, validation.Required),
		"Offset":   validation.Validate(c.Offset, validation.Min(0)),
	}.Filter()
}

// Validate validates CreatePolicyVersionRequest
func (c CreatePolicyVersionRequest) Validate() error {
	return validation.Errors{
		"Description":     validation.Validate(c.Description, validation.Length(0, 255)),
		"MatchRuleFormat": validation.Validate(c.MatchRuleFormat, validation.In(MatchRuleFormat10)),
		"MatchRules":      validation.Validate(c.MatchRules, validation.Length(0, 5000)),
	}.Filter()
}

// Validate validates MatchRuleALB
func (m MatchRuleALB) Validate() error {
	return validation.Errors{
		"Type":                     validation.Validate(m.Type, validation.In(MatchRuleTypeALB)),
		"Name":                     validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                    validation.Validate(m.Start, validation.Min(0)),
		"End":                      validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                 validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"ForwardSettings.OriginID": validation.Validate(m.ForwardSettings.OriginID, validation.Required, validation.Length(0, 8192)),
		"Location":                 validation.Validate(m.Location, validation.Empty),
		"Matches":                  validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleER
func (m MatchRuleER) Validate() error {
	return validation.Errors{
		"Type":           validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeER)),
		"Name":           validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":          validation.Validate(m.Start, validation.Min(0)),
		"End":            validation.Validate(m.End, validation.Min(0)),
		"MatchURL":       validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"RedirectURL":    validation.Validate(m.RedirectURL, validation.Required, validation.Length(1, 8192)),
		"UseRelativeURL": validation.Validate(m.UseRelativeURL, validation.Required, validation.In("none", "copy_scheme_hostname", "relative_url")),
		"StatusCode":     validation.Validate(m.StatusCode, validation.Required, validation.In(301, 302, 303, 307, 308)),
		"Location":       validation.Validate(m.Location, validation.Empty),
		"Matches":        validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchCriteriaALB
func (m MatchCriteriaALB) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("clientip", "continent", "cookie", "countrycode",
			"deviceCharacteristics", "extension", "header", "hostname", "method", "path", "protocol", "proxy", "query", "regioncode", "range")),
		"MatchValue":       validation.Validate(m.MatchValue, validation.Length(0, 8192), validation.Required.When(m.ObjectMatchValue == nil)),
		"MatchOperator":    validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals)),
		"CheckIPs":         validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders)),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "")),
	}.Filter()
}

// Validate validates MatchCriteriaER
func (m MatchCriteriaER) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query",
			"regex", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy")),
		"MatchValue":    validation.Validate(m.MatchValue, validation.Length(0, 8192)),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals)),
		"CheckIPs":      validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders)),
	}.Filter()
}

// Validate validates ObjectMatchValueRangeOrSimpleSubtype
func (o ObjectMatchValueRangeOrSimpleSubtype) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(o.Type, validation.In(ObjectMatchValueRangeOrSimpleTypeSubtypeRange, ObjectMatchValueRangeOrSimpleTypeSubtypeSimple)),
	}.Filter()
}

// Validate validates ObjectMatchValueObjectSubtype
func (o ObjectMatchValueObjectSubtype) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(o.Name, validation.Required, validation.Length(0, 8192)),
		"Type": validation.Validate(o.Type, validation.Required, validation.In(ObjectMatchValueObjectTypeSubtypeObject)),
	}.Filter()
}

// Validate validates UpdatePolicyVersionRequest
func (o UpdatePolicyVersionRequest) Validate() error {
	return validation.Errors{
		"Description":     validation.Validate(o.Description, validation.Length(0, 255)),
		"MatchRuleFormat": validation.Validate(o.MatchRuleFormat, validation.In(MatchRuleFormat10)),
		"MatchRules":      validation.Validate(o.MatchRules, validation.Length(0, 5000)),
	}.Filter()
}

var (
	// ErrListPolicyVersions is returned when ListPolicyVersions fails
	ErrListPolicyVersions = errors.New("list policy versions")
	// ErrGetPolicyVersion is returned when GetPolicyVersion fails
	ErrGetPolicyVersion = errors.New("get policy versions")
	// ErrCreatePolicyVersion is returned when CreatePolicyVersion fails
	ErrCreatePolicyVersion = errors.New("create policy versions")
	// ErrDeletePolicyVersion is returned when DeletePolicyVersion fails
	ErrDeletePolicyVersion = errors.New("delete policy versions")
	// ErrUpdatePolicyVersion is returned when UpdatePolicyVersion fails
	ErrUpdatePolicyVersion = errors.New("update policy versions")
)

func (m MatchRuleALB) cloudletType() string {
	return "albMatchRule"
}

func (m MatchRuleER) cloudletType() string {
	return "erMatchRule"
}

// matchRuleHandlers contains mapping between name of the type for MatchRule and its implementation
// It makes the UnmarshalJSON more compact and easier to support more cloudlet types
var matchRuleHandlers = map[string]func() MatchRule{
	"albMatchRule": func() MatchRule { return &MatchRuleALB{} },
	"erMatchRule":  func() MatchRule { return &MatchRuleER{} },
}

// UnmarshalJSON helps to un-marshall items of MatchRules array as proper instances of *MatchRuleALB or *MatchRuleER
func (m *MatchRules) UnmarshalJSON(b []byte) error {
	data := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for _, matchRule := range data {
		cloudletType, ok := matchRule["type"]
		if !ok {
			return fmt.Errorf("match rule entry should contain 'type' field")
		}
		cloudletTypeName, ok := cloudletType.(string)
		if !ok {
			return fmt.Errorf("'type' field on match rule entry should be a string")
		}
		byteArr, err := json.Marshal(matchRule)
		if err != nil {
			return err
		}

		matchRuleType, ok := matchRuleHandlers[cloudletTypeName]
		if !ok {
			return fmt.Errorf("unsupported match rule type: %s", cloudletTypeName)
		}
		var dst MatchRule
		dst = matchRuleType()
		err = json.Unmarshal(byteArr, dst)
		if err != nil {
			return err
		}
		*m = append(*m, dst)
	}
	return nil
}

// objectMatchValueHandlers contains mapping between name of the type for ObjectMatchValue and its implementation
// It makes the UnmarshalJSON more compact and easier to support more types
var objectMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObjectSubtype{} },
	"range":  func() interface{} { return &ObjectMatchValueRangeOrSimpleSubtype{} },
	"simple": func() interface{} { return &ObjectMatchValueRangeOrSimpleSubtype{} },
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaALB as proper instance of *ObjectMatchValueObjectSubtype or *ObjectMatchValueRangeOrSimpleSubtype
func (m *MatchCriteriaALB) UnmarshalJSON(b []byte) error {
	// matchCriteriaALB is an alias for MatchCriteriaALB for un-marshalling purposes
	type matchCriteriaALB MatchCriteriaALB

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaALB)(m))
	if err != nil {
		return err
	}
	if m.ObjectMatchValue == nil {
		return nil
	}
	objectMatchValue, ok := m.ObjectMatchValue.(interface{})
	if !ok {
		return fmt.Errorf("object match value should be of type 'interface{}', but was %T", m.ObjectMatchValue)
	}

	objectMatchValueMap, ok := objectMatchValue.(map[string]interface{})
	if !ok {
		return fmt.Errorf("structure of ObjectMatchValue should be the map, but was %T", objectMatchValue)
	}
	objectMatchValueType, ok := objectMatchValueMap["type"]
	if !ok {
		return fmt.Errorf("object should contain 'type' field")
	}
	objectMatchValueTypeName, ok := objectMatchValueType.(string)
	if !ok {
		return fmt.Errorf("'type' should be a string")
	}

	createObjectMatchValue, ok := objectMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("ObjectMatchValue. UnmarshalJSON: unexpected type: %s", objectMatchValueTypeName)
	}
	convertedObjectMatchValue := createObjectMatchValue()
	marshal, err := json.Marshal(objectMatchValue)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, convertedObjectMatchValue)
	if err != nil {
		return err
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

func (c *cloudlets) ListPolicyVersions(ctx context.Context, params ListPolicyVersionsRequest) ([]PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("ListPolicyVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPolicyVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicyVersions, err)
	}

	q := uri.Query()
	q.Add("offset", fmt.Sprintf("%d", params.Offset))
	q.Add("includeRules", strconv.FormatBool(params.IncludeRules))
	q.Add("includeDeleted", strconv.FormatBool(params.IncludeDeleted))
	q.Add("includeActivations", strconv.FormatBool(params.IncludeActivations))
	if params.PageSize != nil {
		q.Add("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicyVersions, err)
	}

	var result []PolicyVersion
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicyVersions, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicyVersions, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) GetPolicyVersion(ctx context.Context, params GetPolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicyVersion")

	var result PolicyVersion

	uri, err := url.Parse(fmt.Sprintf(
		"/cloudlets/api/v2/policies/%d/versions/%d",
		params.PolicyID, params.Version),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPolicyVersion, err)
	}

	q := uri.Query()
	q.Add("omitRules", strconv.FormatBool(params.OmitRules))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyVersion, err)
	}

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreatePolicyVersion(ctx context.Context, params CreatePolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("CreatePolicyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreatePolicyVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreatePolicyVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result, params.CreatePolicyVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePolicyVersion, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) DeletePolicyVersion(ctx context.Context, params DeletePolicyVersionRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeletePolicyVersion")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d", params.PolicyID, params.Version))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeletePolicyVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeletePolicyVersion, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeletePolicyVersion, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeletePolicyVersion, c.Error(resp))
	}

	return nil
}

func (c *cloudlets) UpdatePolicyVersion(ctx context.Context, params UpdatePolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdatePolicyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePolicyVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d", params.PolicyID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdatePolicyVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result, params.UpdatePolicyVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicyVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicyVersion, c.Error(resp))
	}

	return &result, nil
}
