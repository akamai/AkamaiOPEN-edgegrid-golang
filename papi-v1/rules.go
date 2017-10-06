package papi

import (
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Rules is a collection of property rules
type Rules struct {
	client.Resource
	AccountID       string        `json:"accountId"`
	ContractID      string        `json:"contractId"`
	GroupID         string        `json:"groupId"`
	PropertyID      string        `json:"propertyId"`
	PropertyVersion int           `json:"propertyVersion"`
	Etag            string        `json:"etag"`
	RuleFormat      string        `json:"ruleFormat"`
	Rule            *Rule         `json:"rules"`
	Errors          []*RuleErrors `json:"errors,omitempty"`
}

// NewRules creates a new Rules
func NewRules() *Rules {
	rules := &Rules{}
	rules.Rule = NewRule()
	rules.Rule.Name = "default"
	rules.Init()

	return rules
}

// PreMarshalJSON is called before JSON marshaling
//
// See: jsonhooks-v1/json.Marshal()
func (rules *Rules) PreMarshalJSON() error {
	rules.Errors = nil
	return nil
}

// GetRules populates Rules with rule data for a given property
//
// See: Property.GetRules
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletree
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRules(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
			property.PropertyID,
			property.LatestVersion,
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	return nil
}

// GetRulesDigest fetches the Etag for a rule tree
//
// See: Property.GetRulesDigest()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletreesdigest
// Endpoint: HEAD /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRulesDigest(property *Property) (string, error) {
	req, err := client.NewRequest(
		Config,
		"HEAD",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
			property.PropertyID,
			property.LatestVersion,
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		nil,
	)
	if err != nil {
		return "", err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return "", err
	}

	if client.IsError(res) {
		return "", client.NewAPIError(res)
	}

	return res.Header.Get("Etag"), nil
}

// Save creates/updates a rule tree for a property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#updatearuletree
// Endpoint: PUT /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) Save() error {
	rules.Errors = []*RuleErrors{}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules/?contractId=%s&groupId=%s",
			rules.PropertyID,
			rules.PropertyVersion,
			rules.ContractID,
			rules.GroupID,
		),
		rules,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	if len(rules.Errors) != 0 {
		return ErrorMap[ErrInvalidRules]
	}

	return nil
}

// AddChildRule adds a rule as a child of the top level default rule
func (rules *Rules) AddChildRule(rule *Rule) error {
	rules.Rule.Children = append(rules.Rule.Children, rule)
	return nil
}

// FindRule locates a specific rule by path
func (rules *Rules) FindRule(path string) (*Rule, error) {
	if path == "" {
		return rules.Rule, nil
	}

	sep := "/"
	segments := strings.Split(path, sep)

	currentRule := rules.Rule
	for _, segment := range segments {
		found := false
		for _, rule := range currentRule.Children {
			if strings.ToLower(rule.Name) == segment {
				currentRule = rule
				found = true
			}
		}
		if found != true {
			return nil, ErrorMap[ErrRuleNotFound]
		}
	}

	return currentRule, nil
}

// Freeze pins a properties rule set to a specific rule set version
func (rules *Rules) Freeze(format string) error {
	rules.Errors = []*RuleErrors{}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules/?contractId=%s&groupId=%s",
			rules.PropertyID,
			rules.PropertyVersion,
			rules.ContractID,
			rules.GroupID,
		),
		rules,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", fmt.Sprintf("application/vnd.akamai.papirules.%s+json", format))

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	if len(rules.Errors) != 0 {
		return ErrorMap[ErrInvalidRules]
	}

	return nil
}

// Rule represents a property rule resource
type Rule struct {
	client.Resource
	Depth               int                          `json:"-"`
	Name                string                       `json:"name"`
	Criteria            []*Criteria                  `json:"criteria,omitempty"`
	Behaviors           []*Behavior                  `json:"behaviors,omitempty"`
	Children            []*Rule                      `json:"children,omitempty"`
	Comments            string                       `json:"comments,omitempty"`
	CriteriaLocked      bool                         `json:"criteriaLocked,omitempty"`
	CriteriaMustSatisfy RuleCriteriaMustSatisfyValue `json:"criteriaMustSatisfy,omitempty"`
	UUID                string                       `json:"uuid,omitempty"`
	Options             struct {
		IsSecure bool `json:"is_secure,omitempty"`
	} `json:"options,omitempty"`
}

// NewRule creates a new Rule
func NewRule() *Rule {
	rule := &Rule{}
	rule.Init()

	return rule
}

// AddChildRule appends a child rule
func (rule *Rule) AddChildRule(child *Rule) {
	for k, v := range rule.Children {
		if v.Name == child.Name {
			rule.Children[k] = child
			return
		}
	}
	rule.Children = append(rule.Children, child)
}

// AddCriteria appends a rule criteria
func (rule *Rule) AddCriteria(criteria *Criteria) {
	for k, v := range rule.Criteria {
		if v.Name == criteria.Name {
			rule.Criteria[k] = criteria
			return
		}
	}
	rule.Criteria = append(rule.Criteria, criteria)
}

// AddBehavior appends a rule behavior
func (rule *Rule) AddBehavior(behavior *Behavior) {
	for k, v := range rule.Behaviors {
		if v.Name == behavior.Name {
			rule.Behaviors[k] = behavior
			return
		}
	}
	rule.Behaviors = append(rule.Behaviors, behavior)
}

// Criteria represents a rule criteria resource
type Criteria struct {
	client.Resource
	Name    string       `json:"name"`
	Options *OptionValue `json:"options"`
}

// NewCriteria creates a new Criteria
func NewCriteria() *Criteria {
	criteria := &Criteria{}
	criteria.Init()

	return criteria
}

// Behavior represents a rule behavior resource
type Behavior struct {
	client.Resource
	Name    string       `json:"name"`
	Options *OptionValue `json:"options"`
}

// NewBehavior creates a new Behavior
func NewBehavior() *Behavior {
	behavior := &Behavior{}
	behavior.Init()

	return behavior
}

// OptionValue represents a generic option value
//
// OptionValue is a map with string keys, and any
// type of value. You can nest OptionValues as necessary
// to create more complex values.
type OptionValue map[string]interface{}

// RuleErrors represents an validate error returned for a rule
type RuleErrors struct {
	client.Resource
	Type         string `json:"type"`
	Title        string `json:"title"`
	Detail       string `json:"detail"`
	Instance     string `json:"instance"`
	BehaviorName string `json:"behaviorName"`
}

// NewRuleErrors creates a new RuleErrors
func NewRuleErrors() *RuleErrors {
	ruleErrors := &RuleErrors{}
	ruleErrors.Init()

	return ruleErrors
}

type RuleCriteriaMustSatisfyValue string

const (
	RuleCriteriaMustSatisfyAll RuleCriteriaMustSatisfyValue = "any"
	RuleCriteriaMustSatisfyAny RuleCriteriaMustSatisfyValue = "all"
)
