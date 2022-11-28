package cloudlets

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

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

	// MatchRuleALB represents an Application Load Balancer (ALB) match rule resource for create or update resource
	MatchRuleALB struct {
		Name            string             `json:"name,omitempty"`
		Type            MatchRuleType      `json:"type,omitempty"`
		Start           int64              `json:"start,omitempty"`
		End             int64              `json:"end,omitempty"`
		ID              int64              `json:"id,omitempty"`
		Matches         []MatchCriteriaALB `json:"matches,omitempty"`
		MatchURL        string             `json:"matchURL,omitempty"`
		MatchesAlways   bool               `json:"matchesAlways"`
		ForwardSettings ForwardSettingsALB `json:"forwardSettings"`
		Disabled        bool               `json:"disabled,omitempty"`
	}

	// ForwardSettingsALB represents forward settings for an Application Load Balancer (ALB)
	ForwardSettingsALB struct {
		OriginID string `json:"originId"`
	}

	// MatchRuleAP represents an API Prioritization (AP) match rule resource for create or update
	MatchRuleAP struct {
		Name               string            `json:"name,omitempty"`
		Type               MatchRuleType     `json:"type,omitempty"`
		Start              int64             `json:"start,omitempty"`
		End                int64             `json:"end,omitempty"`
		ID                 int64             `json:"id,omitempty"`
		Matches            []MatchCriteriaAP `json:"matches,omitempty"`
		MatchURL           string            `json:"matchURL,omitempty"`
		PassThroughPercent *float64          `json:"passThroughPercent"`
		Disabled           bool              `json:"disabled,omitempty"`
	}

	// MatchRuleAS represents an Application Segmentation (AS) match rule resource for create or update resource
	MatchRuleAS struct {
		Name            string            `json:"name,omitempty"`
		Type            MatchRuleType     `json:"type,omitempty"`
		Start           int64             `json:"start,omitempty"`
		End             int64             `json:"end,omitempty"`
		ID              int64             `json:"id,omitempty"`
		Matches         []MatchCriteriaAS `json:"matches,omitempty"`
		MatchURL        string            `json:"matchURL,omitempty"`
		ForwardSettings ForwardSettingsAS `json:"forwardSettings"`
		Disabled        bool              `json:"disabled,omitempty"`
	}

	// ForwardSettingsAS represents forward settings for an Application Segmentation (AS)
	ForwardSettingsAS struct {
		PathAndQS              string `json:"pathAndQS,omitempty"`
		UseIncomingQueryString bool   `json:"useIncomingQueryString,omitempty"`
		OriginID               string `json:"originId,omitempty"`
	}

	// MatchRulePR represents a Phased Release (PR aka CD) match rule resource for create or update resource
	MatchRulePR struct {
		Name            string            `json:"name,omitempty"`
		Type            MatchRuleType     `json:"type,omitempty"`
		Start           int64             `json:"start,omitempty"`
		End             int64             `json:"end,omitempty"`
		ID              int64             `json:"id,omitempty"`
		Matches         []MatchCriteriaPR `json:"matches,omitempty"`
		MatchURL        string            `json:"matchURL,omitempty"`
		ForwardSettings ForwardSettingsPR `json:"forwardSettings"`
		Disabled        bool              `json:"disabled,omitempty"`
		MatchesAlways   bool              `json:"matchesAlways,omitempty"`
	}

	// ForwardSettingsPR represents forward settings for a Phased Release (PR aka CD)
	ForwardSettingsPR struct {
		OriginID string `json:"originId"`
		Percent  int    `json:"percent"`
	}

	// MatchRuleER represents an Edge Redirector (ER) match rule resource for create or update resource
	MatchRuleER struct {
		Name                     string            `json:"name,omitempty"`
		Type                     MatchRuleType     `json:"type,omitempty"`
		Start                    int64             `json:"start,omitempty"`
		End                      int64             `json:"end,omitempty"`
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

	// MatchRuleFR represents a Forward Rewrite (FR) match rule resource for create or update resource
	MatchRuleFR struct {
		Name            string            `json:"name,omitempty"`
		Type            MatchRuleType     `json:"type,omitempty"`
		Start           int64             `json:"start,omitempty"`
		End             int64             `json:"end,omitempty"`
		ID              int64             `json:"id,omitempty"`
		Matches         []MatchCriteriaFR `json:"matches,omitempty"`
		MatchURL        string            `json:"matchURL,omitempty"`
		ForwardSettings ForwardSettingsFR `json:"forwardSettings"`
		Disabled        bool              `json:"disabled,omitempty"`
	}

	// ForwardSettingsFR represents forward settings for a Forward Rewrite (FR)
	ForwardSettingsFR struct {
		PathAndQS              string `json:"pathAndQS,omitempty"`
		UseIncomingQueryString bool   `json:"useIncomingQueryString,omitempty"`
		OriginID               string `json:"originId,omitempty"`
	}

	// MatchRuleRC represents a Request Control (RC aka IG) match rule resource for create or update resource
	MatchRuleRC struct {
		Name          string            `json:"name,omitempty"`
		Type          MatchRuleType     `json:"type,omitempty"`
		Start         int64             `json:"start,omitempty"`
		End           int64             `json:"end,omitempty"`
		ID            int64             `json:"id,omitempty"`
		Matches       []MatchCriteriaRC `json:"matches,omitempty"`
		MatchesAlways bool              `json:"matchesAlways,omitempty"`
		AllowDeny     AllowDeny         `json:"allowDeny"`
		Disabled      bool              `json:"disabled,omitempty"`
	}

	// MatchRuleVP represents a Visitor Prioritization (VP) match rule resource for create or update resource
	MatchRuleVP struct {
		Name               string            `json:"name,omitempty"`
		Type               MatchRuleType     `json:"type,omitempty"`
		Start              int64             `json:"start,omitempty"`
		End                int64             `json:"end,omitempty"`
		ID                 int64             `json:"id,omitempty"`
		Matches            []MatchCriteriaVP `json:"matches,omitempty"`
		MatchURL           string            `json:"matchURL,omitempty"`
		PassThroughPercent *float64          `json:"passThroughPercent"`
		Disabled           bool              `json:"disabled,omitempty"`
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

	// MatchCriteriaALB represents a match criteria resource for match rule for cloudlet Application Load Balancer (ALB)
	// ObjectMatchValue can contain ObjectMatchValueObject, ObjectMatchValueSimple or ObjectMatchValueRange
	MatchCriteriaALB MatchCriteria

	// MatchCriteriaAP represents a match criteria resource for match rule for cloudlet API Prioritization (AP)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaAP MatchCriteria

	// MatchCriteriaAS represents a match criteria resource for match rule for cloudlet Application Segmentation (AS)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple or ObjectMatchValueRange
	MatchCriteriaAS MatchCriteria

	// MatchCriteriaPR represents a match criteria resource for match rule for cloudlet Phased Release (PR aka CD)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaPR MatchCriteria

	// MatchCriteriaER represents a match criteria resource for match rule for cloudlet Edge Redirector (ER)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaER MatchCriteria

	// MatchCriteriaFR represents a match criteria resource for match rule for cloudlet Forward Rewrite (FR)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaFR MatchCriteria

	// MatchCriteriaRC represents a match criteria resource for match rule for cloudlet Request Control (RC aka IG)
	// ObjectMatchValue can contain ObjectMatchValueObject or ObjectMatchValueSimple
	MatchCriteriaRC MatchCriteria

	// MatchCriteriaVP represents a match criteria resource for match rule for cloudlet Visitor Prioritization (VP)
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
	// AllowDeny enum type
	AllowDeny string
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
	// MatchRuleTypeALB represents rule type for Application Load Balancer (ALB) cloudlets
	MatchRuleTypeALB MatchRuleType = "albMatchRule"
	// MatchRuleTypeAP represents rule type for API Prioritization (AP) cloudlets
	MatchRuleTypeAP MatchRuleType = "apMatchRule"
	// MatchRuleTypeAS represents rule type for Application Segmentation (AS) cloudlets
	MatchRuleTypeAS MatchRuleType = "asMatchRule"
	// MatchRuleTypePR represents rule type for Phased Release (PR aka CD) cloudlets
	MatchRuleTypePR MatchRuleType = "cdMatchRule"
	// MatchRuleTypeER represents rule type for Edge Redirector (ER) cloudlets
	MatchRuleTypeER MatchRuleType = "erMatchRule"
	// MatchRuleTypeFR represents rule type for Forward Rewrite (FR) cloudlets
	MatchRuleTypeFR MatchRuleType = "frMatchRule"
	// MatchRuleTypeRC represents rule type for Request Control (RC aka IG) cloudlets
	MatchRuleTypeRC MatchRuleType = "igMatchRule"
	// MatchRuleTypeVP represents rule type for Visitor Prioritization (VP) cloudlets
	MatchRuleTypeVP MatchRuleType = "vpMatchRule"
)

const (
	// MatchRuleFormat10 represents default match rule format
	MatchRuleFormat10 MatchRuleFormat = "1.0"
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
	// Allow represents allow option
	Allow AllowDeny = "allow"
	// Deny represents deny option
	Deny AllowDeny = "deny"
	// DenyBranded represents denybranded option
	DenyBranded AllowDeny = "denybranded"
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

var (
	// ErrUnmarshallMatchCriteriaALB is returned when unmarshalling of MatchCriteriaALB fails
	ErrUnmarshallMatchCriteriaALB = errors.New("unmarshalling MatchCriteriaALB")
	// ErrUnmarshallMatchCriteriaAP is returned when unmarshalling of MatchCriteriaAP fails
	ErrUnmarshallMatchCriteriaAP = errors.New("unmarshalling MatchCriteriaAP")
	// ErrUnmarshallMatchCriteriaAS is returned when unmarshalling of MatchCriteriaAS fails
	ErrUnmarshallMatchCriteriaAS = errors.New("unmarshalling MatchCriteriaAS")
	// ErrUnmarshallMatchCriteriaPR is returned when unmarshalling of MatchCriteriaPR fails
	ErrUnmarshallMatchCriteriaPR = errors.New("unmarshalling MatchCriteriaPR")
	// ErrUnmarshallMatchCriteriaER is returned when unmarshalling of MatchCriteriaER fails
	ErrUnmarshallMatchCriteriaER = errors.New("unmarshalling MatchCriteriaER")
	// ErrUnmarshallMatchCriteriaFR is returned when unmarshalling of MatchCriteriaFR fails
	ErrUnmarshallMatchCriteriaFR = errors.New("unmarshalling MatchCriteriaFR")
	// ErrUnmarshallMatchCriteriaRC is returned when unmarshalling of MatchCriteriaRC fails
	ErrUnmarshallMatchCriteriaRC = errors.New("unmarshalling MatchCriteriaRC")
	// ErrUnmarshallMatchCriteriaVP is returned when unmarshalling of MatchCriteriaVP fails
	ErrUnmarshallMatchCriteriaVP = errors.New("unmarshalling MatchCriteriaVP")
	// ErrUnmarshallMatchRules is returned when unmarshalling of MatchRules fails
	ErrUnmarshallMatchRules = errors.New("unmarshalling MatchRules")
)

// matchRuleHandlers contains mapping between name of the type for MatchRule and its implementation
// It makes the UnmarshalJSON more compact and easier to support more cloudlet types
var matchRuleHandlers = map[string]func() MatchRule{
	"albMatchRule": func() MatchRule { return &MatchRuleALB{} },
	"apMatchRule":  func() MatchRule { return &MatchRuleAP{} },
	"asMatchRule":  func() MatchRule { return &MatchRuleAS{} },
	"cdMatchRule":  func() MatchRule { return &MatchRulePR{} },
	"erMatchRule":  func() MatchRule { return &MatchRuleER{} },
	"frMatchRule":  func() MatchRule { return &MatchRuleFR{} },
	"igMatchRule":  func() MatchRule { return &MatchRuleRC{} },
	"vpMatchRule":  func() MatchRule { return &MatchRuleVP{} },
}

// objectOrRangeOrSimpleMatchValueHandlers contains mapping between name of the type for ObjectMatchValue and its implementation
// It makes the UnmarshalJSON more compact and easier to support more types
var objectOrRangeOrSimpleMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObject{} },
	"range":  func() interface{} { return &ObjectMatchValueRange{} },
	"simple": func() interface{} { return &ObjectMatchValueSimple{} },
}

// simpleObjectMatchValueHandlers contains mapping between name of the types (simple or object) for ObjectMatchValue and their implementations
// It makes the UnmarshalJSON more compact and easier to support more types
var simpleObjectMatchValueHandlers = map[string]func() interface{}{
	"object": func() interface{} { return &ObjectMatchValueObject{} },
	"simple": func() interface{} { return &ObjectMatchValueSimple{} },
}

// Validate validates MatchRules
func (m MatchRules) Validate() error {
	type matchRules MatchRules

	errs := validation.Errors{
		"MatchRules": validation.Validate(matchRules(m), validation.Length(0, 5000)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates MatchRuleALB
func (m MatchRuleALB) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeALB).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'albMatchRule'", (&m).Type))),
		"Name":                     validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                    validation.Validate(m.Start, validation.Min(0)),
		"End":                      validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                 validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"ForwardSettings.OriginID": validation.Validate(m.ForwardSettings.OriginID, validation.Required, validation.Length(0, 8192)),
		"Matches":                  validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleAP
func (m MatchRuleAP) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeAP).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'apMatchRule'", (&m).Type))),
		"Name":               validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":              validation.Validate(m.Start, validation.Min(0)),
		"End":                validation.Validate(m.End, validation.Min(0)),
		"MatchURL":           validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"PassThroughPercent": validation.Validate(m.PassThroughPercent, validation.By(passThroughPercentValidation)),
		"Matches":            validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleAS
func (m MatchRuleAS) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeAS).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'asMatchRule'", (&m).Type))),
		"Name":                      validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                     validation.Validate(m.Start, validation.Min(0)),
		"End":                       validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                  validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"Matches":                   validation.Validate(m.Matches),
		"ForwardSettings":           validation.Validate(m.ForwardSettings, validation.Required),
		"ForwardSettings.PathAndQS": validation.Validate(m.ForwardSettings.PathAndQS, validation.Length(1, 8192)),
		"ForwardSettings.OriginID":  validation.Validate(m.ForwardSettings.OriginID, validation.Length(0, 8192)),
	}.Filter()
}

// Validate validates MatchRulePR
func (m MatchRulePR) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypePR).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'cdMatchRule'", (&m).Type))),
		"Name":                     validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                    validation.Validate(m.Start, validation.Min(0)),
		"End":                      validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                 validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"ForwardSettings":          validation.Validate(m.ForwardSettings, validation.Required),
		"ForwardSettings.OriginID": validation.Validate(m.ForwardSettings.OriginID, validation.Required, validation.Length(0, 8192)),
		"ForwardSettings.Percent":  validation.Validate(m.ForwardSettings.Percent, validation.Required, validation.Min(1), validation.Max(100)),
		"Matches":                  validation.Validate(m.Matches),
	}.Filter()
}

// Validate validates MatchRuleER
func (m MatchRuleER) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeER).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'erMatchRule'", (&m).Type))),
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

// Validate validates MatchRuleFR
func (m MatchRuleFR) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeFR).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'frMatchRule'", (&m).Type))),
		"Name":                      validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":                     validation.Validate(m.Start, validation.Min(0)),
		"End":                       validation.Validate(m.End, validation.Min(0)),
		"MatchURL":                  validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"Matches":                   validation.Validate(m.Matches),
		"ForwardSettings":           validation.Validate(m.ForwardSettings, validation.Required),
		"ForwardSettings.PathAndQS": validation.Validate(m.ForwardSettings.PathAndQS, validation.Length(1, 8192)),
		"ForwardSettings.OriginID":  validation.Validate(m.ForwardSettings.OriginID, validation.Length(0, 8192)),
	}.Filter()
}

// Validate validates MatchRuleRC
func (m MatchRuleRC) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeRC).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'igMatchRule'", (&m).Type))),
		"Name":    validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":   validation.Validate(m.Start, validation.Min(0)),
		"End":     validation.Validate(m.End, validation.Min(0)),
		"Matches": validation.Validate(m.Matches, validation.When(m.MatchesAlways, validation.Empty.Error("must be blank when 'matchesAlways' is set"))),
		"AllowDeny": validation.Validate(m.AllowDeny, validation.Required, validation.In(Allow, Deny, DenyBranded).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s' or '%s'", (&m).AllowDeny, Allow, Deny, DenyBranded),
		)),
	}.Filter()
}

// Validate validates MatchRuleVP
func (m MatchRuleVP) Validate() error {
	return validation.Errors{
		"Type": validation.Validate(m.Type, validation.Required, validation.In(MatchRuleTypeVP).Error(
			fmt.Sprintf("value '%s' is invalid. Must be: 'vpMatchRule'", (&m).Type))),
		"Name":               validation.Validate(m.Name, validation.Length(0, 8192)),
		"Start":              validation.Validate(m.Start, validation.Min(0)),
		"End":                validation.Validate(m.End, validation.Min(0)),
		"MatchURL":           validation.Validate(m.MatchURL, validation.Length(0, 8192)),
		"PassThroughPercent": validation.Validate(m.PassThroughPercent, validation.By(passThroughPercentValidation)),
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
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrRangeOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaAP
func (m MatchCriteriaAP) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In(
			"header", "hostname", "path", "extension", "query", "cookie", "deviceCharacteristics", "clientip",
			"continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaAS
func (m MatchCriteriaAS) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query", "range",
			"regex", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'range', "+
				"'regex', 'cookie', 'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrRangeOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaPR
func (m MatchCriteriaPR) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension",
			"query", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol",
			"method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaER
func (m MatchCriteriaER) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query",
			"regex", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'regex', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy' or '' (empty)", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaFR
func (m MatchCriteriaFR) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.Required, validation.In("header", "hostname", "path", "extension", "query", "regex",
			"cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'regex', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaRC
func (m MatchCriteriaRC) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.Required, validation.In("header", "hostname", "path", "extension", "query", "cookie",
			"deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'cookie', 'deviceCharacteristics', "+
				"'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
	}.Filter()
}

// Validate validates MatchCriteriaVP
func (m MatchCriteriaVP) Validate() error {
	return validation.Errors{
		"MatchType": validation.Validate(m.MatchType, validation.In("header", "hostname", "path", "extension", "query",
			"cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'header', 'hostname', 'path', 'extension', 'query', 'cookie', "+
				"'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'", (&m).MatchType))),
		"MatchValue": validation.Validate(m.MatchValue, validation.Length(1, 8192), validation.Required.When(m.ObjectMatchValue == nil).Error("cannot be blank when ObjectMatchValue is blank"),
			validation.Empty.When(m.ObjectMatchValue != nil).Error("must be blank when ObjectMatchValue is set")),
		"MatchOperator": validation.Validate(m.MatchOperator, validation.In(MatchOperatorContains, MatchOperatorExists, MatchOperatorEquals).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'contains', 'exists', 'equals' or '' (empty)", (&m).MatchOperator))),
		"CheckIPs": validation.Validate(m.CheckIPs, validation.In(CheckIPsConnectingIP, CheckIPsXFFHeaders, CheckIPsConnectingIPXFFHeaders).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS' or '' (empty)", (&m).CheckIPs))),
		"ObjectMatchValue": validation.Validate(m.ObjectMatchValue, validation.Required.When(m.MatchValue == "").Error("cannot be blank when MatchValue is blank"),
			validation.Empty.When(m.MatchValue != "").Error("must be blank when MatchValue is set"), validation.By(objectMatchValueSimpleOrObjectValidation)),
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

func passThroughPercentValidation(value interface{}) error {
	v, ok := value.(*float64)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be *float64", value)
	}
	if v == nil {
		return fmt.Errorf("cannot be blank")
	}
	if *v < -1 {
		return fmt.Errorf("must be no less than -1")
	}
	if *v > 100 {
		return fmt.Errorf("must be no greater than 100")
	}
	return nil
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

func (m MatchRuleALB) cloudletType() string {
	return "albMatchRule"
}

func (m MatchRuleAP) cloudletType() string {
	return "apMatchRule"
}

func (m MatchRuleAS) cloudletType() string {
	return "asMatchRule"
}

func (m MatchRulePR) cloudletType() string {
	return "cdMatchRule"
}

func (m MatchRuleER) cloudletType() string {
	return "erMatchRule"
}

func (m MatchRuleFR) cloudletType() string {
	return "frMatchRule"
}

func (m MatchRuleRC) cloudletType() string {
	return "igMatchRule"
}

func (m MatchRuleVP) cloudletType() string {
	return "vpMatchRule"
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

	createObjectMatchValue, ok := objectOrRangeOrSimpleMatchValueHandlers[objectMatchValueTypeName]
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

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaAP as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaAP) UnmarshalJSON(b []byte) error {
	// matchCriteriaER is an alias for MatchCriteriaER for un-marshalling purposes
	type matchCriteriaAP MatchCriteriaAP

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaAP)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAP, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAP, err)
	}

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaAP, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAP, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaAS as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple or *ObjectMatchValueRange
func (m *MatchCriteriaAS) UnmarshalJSON(b []byte) error {
	// matchCriteriaAS is an alias for MatchCriteriaAS for un-marshalling purposes
	type matchCriteriaAS MatchCriteriaAS

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaAS)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAS, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAS, err)
	}

	createObjectMatchValue, ok := objectOrRangeOrSimpleMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaAS, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaAS, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaPR as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaPR) UnmarshalJSON(b []byte) error {
	// matchCriteriaPR is an alias for MatchCriteriaPR for un-marshalling purposes
	type matchCriteriaPR MatchCriteriaPR

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaPR)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaPR, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaPR, err)
	}

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaPR, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaPR, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
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

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
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

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaFR as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaFR) UnmarshalJSON(b []byte) error {
	// matchCriteriaFR is an alias for MatchCriteriaFR for un-marshalling purposes
	type matchCriteriaFR MatchCriteriaFR

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaFR)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaFR, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaFR, err)
	}

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaFR, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaFR, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaRC as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaRC) UnmarshalJSON(b []byte) error {
	// matchCriteriaRC is an alias for MatchCriteriaRC for un-marshalling purposes
	type matchCriteriaRC MatchCriteriaRC

	// populate common attributes using default json unmarshaler using aliased type
	err := json.Unmarshal(b, (*matchCriteriaRC)(m))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaRC, err)
	}
	if m.ObjectMatchValue == nil {
		return nil
	}

	objectMatchValueTypeName, err := getObjectMatchValueType(m.ObjectMatchValue)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaRC, err)
	}

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
	if !ok {
		return fmt.Errorf("%w: objectMatchValue has unexpected type: '%s'", ErrUnmarshallMatchCriteriaRC, objectMatchValueTypeName)
	}
	convertedObjectMatchValue, err := convertObjectMatchValue(m.ObjectMatchValue, createObjectMatchValue())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallMatchCriteriaRC, err)
	}
	m.ObjectMatchValue = convertedObjectMatchValue

	return nil
}

// UnmarshalJSON helps to un-marshall field ObjectMatchValue of MatchCriteriaVP as proper instance of *ObjectMatchValueObject or *ObjectMatchValueSimple
func (m *MatchCriteriaVP) UnmarshalJSON(b []byte) error {
	// matchCriteriaVP is an alias for MatchCriteriaVP for un-marshalling purposes
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

	createObjectMatchValue, ok := simpleObjectMatchValueHandlers[objectMatchValueTypeName]
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
