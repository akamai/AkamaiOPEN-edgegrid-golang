package cloudlets

import (
	"encoding/json"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
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
		MatchURL        string             `json:"matchURL,omitempty"`
		MatchesAlways   bool               `json:"matchesAlways"`
		ForwardSettings ForwardSettings    `json:"forwardSettings"`
		Disabled        bool               `json:"disabled,omitempty"`
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
		UseRelativeURL           string            `json:"useRelativeUrl,omitempty"`
		StatusCode               int               `json:"statusCode"`
		RedirectURL              string            `json:"redirectURL"`
		MatchURL                 string            `json:"matchURL,omitempty"`
		UseIncomingQueryString   bool              `json:"useIncomingQueryString"`
		UseIncomingSchemeAndHost bool              `json:"useIncomingSchemeAndHost"`
		Disabled                 bool              `json:"disabled,omitempty"`
	}

	// MatchRuleVP represents a match rule resource for create or update resource
	MatchRuleVP struct {
		Name                   string            `json:"name,omitempty"`
		Type                   MatchRuleType     `json:"type,omitempty"`
		Start                  int               `json:"start,omitempty"`
		End                    int               `json:"end,omitempty"`
		ID                     int64             `json:"id,omitempty"`
		Matches                []MatchCriteriaVP `json:"matches,omitempty"`
		MatchURL               string            `json:"matchURL,omitempty"`
		UseIncomingQueryString bool              `json:"useIncomingQueryString,omitempty"`
		PassThroughPercent     float64           `json:"passThroughPercent"`
		Disabled               bool              `json:"disabled,omitempty"`
	}

	// MatchCriteria represents a match criteria resource for match rule for cloudlet
	MatchCriteria struct {
		MatchType        string        `json:"matchType,omitempty"`
		MatchValue       string        `json:"matchValue,omitempty"`
		MatchOperator    MatchOperator `json:"matchOperator,omitempty"`
		CaseSensitive    bool          `json:"caseSensitive"`
		Negate           bool          `json:"negate"`
		CheckIPs         CheckIPs      `json:"checkIPs,omitempty"`
		ObjectMatchValue interface{}   `json:"objectMatchValue,omitempty"`
	}

	// MatchCriteriaALB represents a match criteria resource for match rule for cloudlet ALB
	// ObjectMatchValue can contain ObjectMatchValueObject, ObjectMatchValueSimple or ObjectMatchValueRange
	MatchCriteriaALB MatchCriteria

	// MatchCriteriaER represents a match criteria resource for match rule for cloudlet ER
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaER MatchCriteria

	// MatchCriteriaVP represents a match criteria resource for match rule for cloudlet VP
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaVP MatchCriteria

	// ObjectMatchValueObject represents an object match value resource for match criteria of type object
	ObjectMatchValueObject struct {
		Name              string                     `json:"name"`
		Type              ObjectMatchValueObjectType `json:"type"`
		NameCaseSensitive bool                       `json:"nameCaseSensitive"`
		NameHasWildcard   bool                       `json:"nameHasWildcard"`
		Options           *Options                   `json:"options,omitempty"`
	}

	// ObjectMatchValueSimple represents an object match value resource for match criteria of type simple
	ObjectMatchValueSimple struct {
		Type  ObjectMatchValueSimpleType `json:"type"`
		Value []string                   `json:"value,omitempty"`
	}

	// ObjectMatchValueRange represents an object match value resource for match criteria of type range
	ObjectMatchValueRange struct {
		Type  ObjectMatchValueRangeType `json:"type"`
		Value []int64                   `json:"value,omitempty"`
	}

	// Options represents an option resource for ObjectMatchValueObject
	Options struct {
		Value              []string `json:"value,omitempty"`
		ValueHasWildcard   bool     `json:"valueHasWildcard,omitempty"`
		ValueCaseSensitive bool     `json:"valueCaseSensitive,omitempty"`
		ValueEscaped       bool     `json:"valueEscaped,omitempty"`
	}

	//MatchRuleType enum type
	MatchRuleType string
	// MatchRuleFormat enum type
	MatchRuleFormat string
	// MatchOperator enum type
	MatchOperator string
	// CheckIPs enum type
	CheckIPs string
	// ObjectMatchValueRangeType enum type
	ObjectMatchValueRangeType string
	// ObjectMatchValueSimpleType enum type
	ObjectMatchValueSimpleType string
	// ObjectMatchValueObjectType enum type
	ObjectMatchValueObjectType string
)

const (
	// MatchRuleTypeALB represents rule type for ALB cloudlets
	MatchRuleTypeALB MatchRuleType = "albMatchRule"
	// MatchRuleTypeER represents rule type for ER cloudlets
	MatchRuleTypeER MatchRuleType = "erMatchRule"
	// MatchRuleTypeVP represents rule type for VP cloudlets
	MatchRuleTypeVP MatchRuleType = "vpMatchRule"
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
	// Range represents range option
	Range ObjectMatchValueRangeType = "range"
	// Simple represents simple option
	Simple ObjectMatchValueSimpleType = "simple"
	// Object represents object option
	Object ObjectMatchValueObjectType = "object"
)

// Validate validates MatchRuleALB
func (m MatchRuleALB) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.In(MatchRuleTypeALB).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'albMatchRule' or '' (empty)", (&m).Type))),
		"Name":                     validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                    validation.Validate(m.Start, validation.Min(0)),
		"End":                      validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                 validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"ForwardSettings.OriginID": validation.Validate(m.ForwardSettings.OriginID, validation.Required, validation.Length(0, 8192)),
		"Matches":                  validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleER
func (m MatchRuleER) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeER).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'erMatchRule' or '' (empty)", (&m).Type))),
		"Name":        validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":       validation.Validate(m.Start, validation.Min(0)),
		"End":         validation.Validate(m.End, validation.Min(0)),
		"MatchURL":    validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"RedirectURL": validation.Validate(m.RedirectURL, validation.Required, validation.Length(1, 8192)),
		"UseRelativeURL": validation.Validate(m.UseRelativeURL, validation.In("none", "copy_scheme_hostname", "relative_url").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'none', 'copy_scheme_hostname', 'relative_url' or '' (empty)", (&m).UseRelativeURL))),
		"StatusCode": validation.Validate(m.StatusCode, validation.Required, validation.In(301, 302, 303, 307, 308).Error(
			fmt.Sprintf("value '%d' is invalid. Must be one of: 301, 302, 303, 307 or 308", (&m).StatusCode))),
		"Matches": validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleVP
func (m MatchRuleVP) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeVP).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'vpMatchRule' or '' (empty)", (&m).Type))),
		"Name":               validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":              validation.Validate(m.Start, validation.Min(0)),
		"End":                validation.Validate(m.End, validation.Min(0)),
		"MatchURL":           validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"PassThroughPercent": validation.Validate(m.PassThroughPercent, validation.Required, validation.Min(-1.0), validation.Max(100.0)),
		"Matches":            validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchCriteriaALB
func (m MatchCriteriaALB) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("clientip", "continent", "cookie", "countrycode",
			"deviceCharacteristics", "extension", "header", "hostname", "method", "path", "protocol", "proxy", "query", "regioncode", "range").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'clientip', 'continent', 'cookie', 'countrycode', 'deviceCharacteristics', "+
				"'extension', 'header', 'hostname', 'method', 'path', 'protocol', 'proxy', 'query', 'regioncode', 'range' or '' (empty)", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(0, 8192), validation.Required.When(m.ObjectMatchValue == nil)),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == ""), validation.By(objectMatchValueSimpleOrRangeOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaER
func (m MatchCriteriaER) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query",
			"regex", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'regex', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy' or '' (empty)", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(0, 8192)),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == ""), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaVP
func (m MatchCriteriaVP) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query",
			"cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy' or '' (empty)", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(0, 8192)),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == ""), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

func objectMatchValueSimpleOrObjectValidation(value interface{}) error {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case *ObjectMatchValueObject, *ObjectMatchValueSimple:
		return nil
	default:
		return fmt.Errorf("type %T is invalid. Must be one of: 'simple' or 'object'", value)
	}
}

func objectMatchValueSimpleOrRangeOrObjectValidation(value interface{}) error {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case *ObjectMatchValueObject, *ObjectMatchValueSimple, *ObjectMatchValueRange:
		return nil
	default:
		return fmt.Errorf("type %T is invalid. Must be one of: 'simple', 'range' or 'object'", value)
	}
}

// Validate validates ObjectMatchValueRange
func (o ObjectMatchValueRange) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(o.Type, validation.In(Range).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'range'", (&o).Type))),
	}.Filter()
}

// Validate validates ObjectMatchValueSimple
func (o ObjectMatchValueSimple) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(o.Type, validation.In(Simple).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'simple'", (&o).Type))),
	}.Filter()
}

// Validate validates ObjectMatchValueObject
func (o ObjectMatchValueObject) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(o.Name, validation.Required, validation.Length(0, 8192)),
		"Type": validation.Validate(o.Type, validation.Required, validation.In(Object).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'object'", (&o).Type))),
	}.Filter()
}

var (
	// ErrUnmarshallMatchCriteriaALB is returned when unmarshalling of MatchCriteriaALB fails
	ErrUnmarshallMatchCriteriaALB = errors.New("unmarshalling MatchCriteriaALB")
	// ErrUnmarshallMatchCriteriaER is returned when unmarshalling of MatchCriteriaER fails
	ErrUnmarshallMatchCriteriaER = errors.New("unmarshalling MatchCriteriaER")
	// ErrUnmarshallMatchCriteriaVP is returned when unmarshalling of MatchCriteriaVP fails
	ErrUnmarshallMatchCriteriaVP = errors.New("unmarshalling MatchCriteriaVP")
	// ErrUnmarshallMatchRules is returned when unmarshalling of MatchRules fails
	ErrUnmarshallMatchRules = errors.New("unmarshalling MatchRules")
)

func (m MatchRuleALB) cloudletType() string {
	return "albMatchRule"
}

func (m MatchRuleER) cloudletType() string {
	return "erMatchRule"
}

func (m MatchRuleVP) cloudletType() string {
	return "vpMatchRule"
}

// matchRuleHandlers contains mapping between name of the type for MatchRule and its implementation
// It makes the UnmarshalJSON more compact and easier to support more cloudlet types
var matchRuleHandlers = map[string]func() MatchRule{
	"albMatchRule": func() MatchRule { return &MatchRuleALB{} },
	"erMatchRule":  func() MatchRule { return &MatchRuleER{} },
	"vpMatchRule":  func() MatchRule { return &MatchRuleVP{} },
}

// UnmarshalJSON helps to un-marshall items of MatchRules array as proper instances of *MatchRuleALB or *MatchRuleER
func (m *MatchRules) UnmarshalJSON(b []byte) error {
	data := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchRules, err)
	}
	for _, matchRule := range data {
		cloudletType, ok := matchRule["type"]
		if !ok {
			return fmt.Errorf("%w: match rule entry should contain 'type' field", ErrUnmarshallMatchRules)
		}
		cloudletTypeName, ok := cloudletType.(string)
		if !ok {
			return fmt.Errorf("%w: 'type' field on match rule entry should be a string", ErrUnmarshallMatchRules)
		}
		byteArr, err := json.Marshal(matchRule)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshallMatchRules, err)
		}

		matchRuleType, ok := matchRuleHandlers[cloudletTypeName]
		if !ok {
			return fmt.Errorf("%w: unsupported match rule type: %s", ErrUnmarshallMatchRules, cloudletTypeName)
		}
		dst := matchRuleType()
		err = json.Unmarshal(byteArr, dst)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshallMatchRules, err)
		}
		*m = append(*m, dst)
	}
	return nil
}

// objectALBMatchValueHandlers contains mapping between name of the type for ObjectMatchValue and its implementation
// It makes the UnmarshalJSON more compact and easier to support more types
var objectALBMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObject{} },
	"range":  func() interface{} { return &ObjectMatchValueRange{} },
	"simple": func() interface{} { return &ObjectMatchValueSimple{} },
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaALB as proper instance of *ObjectMatchValueObject, *ObjectMatchValueSimple or *ObjectMatchValueRange
func (m *MatchCriteriaALB) UnmarshalJSON(b []byte) error {
	// matchCriteriaALB is an alias for MatchCriteriaALB for un-marshalling purposes
	type matchCriteriaALB MatchCriteriaALB

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaALB)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaALB, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaALB, err)
	}

	createObjectMatchValue, ok := objectALBMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaALB, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaALB, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// objectERMatchValueHandlers contains mapping between name of the type for ObjectMatchValue and its implementation
// It makes the UnmarshalJSON more compact and easier to support more types
var objectERMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObject{} },
	"simple": func() interface{} { return &ObjectMatchValueSimple{} },
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaER as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaER) UnmarshalJSON(b []byte) error {
	// matchCriteriaER is an alias for MatchCriteriaER for un-marshalling purposes
	type matchCriteriaER MatchCriteriaER

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaER)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaER, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaER, err)
	}

	createObjectMatchValue, ok := objectERMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaER, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaER, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// objectVPMatchValueHandlers contains mapping between name of the type for ObjectMatchValue and its implementation
// It makes the UnmarshalJSON more compact and easier to support more types
var objectVPMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObject{} },
	"simple": func() interface{} { return &ObjectMatchValueSimple{} },
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaER as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaVP) UnmarshalJSON(b []byte) error {
	// matchCriteriaER is an alias for MatchCriteriaER for un-marshalling purposes
	type matchCriteriaVP MatchCriteriaVP

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaVP)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaVP, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaVP, err)
	}

	createObjectMatchValue, ok := objectVPMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaVP, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaVP, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

func getObjectMatchValueType(omv interface{}) (string, error) {
	objectMatchValueMap, ok := omv.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("structure of objectMatchValue should be 'map', but was '%T'", omv)
	}
	objectMatchValueType, ok := objectMatchValueMap["type"]
	if !ok {
		return "", fmt.Errorf("objectMatchValue should contain 'type' field")
	}
	objectMatchValueTypeName, ok := objectMatchValueType.(string)
	if !ok {
		return "", fmt.Errorf("'type' should be a string")
	}
	return objectMatchValueTypeName, nil
}

func convertObjectMatchValue(in, out interface{}) (interface{}, error) {
	marshal, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	err = json.Unmarshal(marshal, out)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return out, nil
}
