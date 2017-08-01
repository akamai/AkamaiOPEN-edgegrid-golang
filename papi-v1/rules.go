package papi

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/xeipuuv/gojsonschema"
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
	Rules           *Rule         `json:"rules"`
	Errors          []*RuleErrors `json:"errors,omitempty"`
}

// NewRules creates a new Rules
func NewRules() *Rules {
	rules := &Rules{}
	rules.Rules = NewRule(rules)
	rules.Init()

	return rules
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (rules *Rules) PostUnmarshalJSON() error {
	rules.Init()

	for key := range rules.Rules.Behaviors {
		rules.Rules.Behaviors[key].parent = rules.Rules
		if len(rules.Rules.Children) > 0 {
			for _, v := range rules.Rules.GetChildren(0, 0) {
				for _, j := range v.Behaviors {
					j.parent = rules.Rules
				}
			}
		}
	}

	for key := range rules.Rules.Criteria {
		rules.Rules.Criteria[key].parent = rules.Rules
	}

	return nil
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

// GetAllRules returns a flattened rule tree for easy iteration
//
// Each Rule has a Depth property that makes it possible to identify
// it's original placement in the tree.
func (rules *Rules) GetAllRules() []*Rule {
	var flatRules []*Rule
	flatRules = append(flatRules, rules.Rules)
	flatRules = append(flatRules, rules.Rules.GetChildren(0, 0)...)

	return flatRules
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

// AddChildBehavior adds a behavior as a child of the given rule path
//
// If the behavior already exists, it's options are merged with the existing
// options.
func (rules *Rules) AddChildBehavior(path string, behavior *Behavior) error {
	_, err := rules.FindBehavior(path + "/" + behavior.Name)
	if err == nil {
		path = path + "/" + behavior.Name

		return rules.AddBehaviorOptions(path, *behavior.Options)
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Behaviors = append(parent.Behaviors, behavior)
	return nil
}

// SetChildBehavior adds a behavior as a child of the given rule path
//
// If the behavior already exists it is replaced with the given behavior
func (rules *Rules) SetChildBehavior(path string, behavior *Behavior) error {
	existingBehavior, err := rules.FindBehavior(path + "/" + behavior.Name)
	if err == nil {
		*existingBehavior = *behavior
		return nil
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Behaviors = append(parent.Behaviors, behavior)
	return nil
}

// SetBehaviorOptions sets the options on a given behavior path
//
// path is a / delimited path from the root of the rule set to the behavior.
// All existing options are overwritten. To add/replace options see AddBehaviorOptions
// instead.
//
// For example, to set the CP Code, the behavior exists at the root, and is called "cpCode",
// making the path, "/cpCode":
//
//	rules.SetBehaviorOptions(
//		"/cpCode",
//		edgegrid.OptionValue{
//			"value": edgegrid.OptionValue{
//				"id": cpcode,
//			},
//		},
//	)
//
// The adaptiveImageCompression behavior defaults to being under the Performance -> JPEG Images,
// making the path "/Performance/JPEG Images/adaptiveImageCompression":
//
//	rules.SetBehaviorOptions(
//		"/Performance/JPEG Images/adaptiveImageCompression"",
//		edgegrid.OptionValue{
//			"tier3StandardCompressionValue": 30,
//		},
//	)
//
// However, this would replace all other options for that behavior, meaning the example
// would fail to validate without the other required options.
func (rules *Rules) SetBehaviorOptions(path string, newOptions OptionValue) error {
	behavior, err := rules.FindBehavior(path)
	if err != nil && err != ErrorMap[ErrBehaviorNotFound] {
		return err
	}

	// Create the missing behavior
	if err != nil && err == ErrorMap[ErrBehaviorNotFound] {
		parent, err := rules.FindParentRule(path)
		if err != nil {
			return err
		}

		sep := "/"
		segments := strings.Split(path, sep)
		behavior = NewBehavior(parent)
		behavior.Name = segments[len(segments)-1]
		rules.AddChildBehavior(strings.Join(segments[:len(segments)-1], sep), behavior)
	}

	behavior.Options = &newOptions
	return nil
}

// AddBehaviorOptions adds/replaces options on a given behavior path
//
// path is a / delimited path from the root of the rule set to the behavior.
// Individual existing options are overwritten. To replace all options, see SetBehaviorOptions
// instead.
//
// For example, to change the CP Code, the behavior exists at the root, and is called "cpCode",
// making the path, "/cpCode":
//
//	rules.AddBehaviorOptions(
//		"/cpCode",
//		edgegrid.OptionValue{
//			"value": edgegrid.OptionValue{
//				"id": cpcode,
//			},
//		},
//	)
//
// The adaptiveImageCompression behavior defaults to being under the Performance -> JPEG Images,
// making the path "/Performance/JPEG Images/adaptiveImageCompression":
//
//	rules.AddBehaviorOptions(
//		"/Performance/JPEG Images/adaptiveImageCompression",
//		edgegrid.OptionValue{
//			"tier3StandardCompressionValue": 30,
//		},
//	)
//
// This will only change the "tier3StandardCompressionValue" option value.
func (rules *Rules) AddBehaviorOptions(path string, newOptions OptionValue) error {
	behavior, err := rules.FindBehavior(path)
	if err != nil && err != ErrorMap[ErrBehaviorNotFound] {
		return err
	}

	// Create the missing behavior
	if err != nil && err == ErrorMap[ErrBehaviorNotFound] {
		parent, err := rules.FindParentRule(path)
		if err != nil {
			return err
		}

		sep := "/"
		segments := strings.Split(path, sep)
		behavior = NewBehavior(parent)
		behavior.Name = segments[len(segments)-1]
		rules.AddChildBehavior(strings.Join(segments[:len(segments)-1], sep), behavior)
	}

	options := *behavior.Options
	for key, value := range newOptions {
		options[key] = value
	}
	behavior.Options = &options

	return nil
}

// AddChildRule adds a rule as a child of the rule specified by the given path
//
// If the rule already exists, criteria, behaviors, and child rules are added to
// the existing rule.
func (rules *Rules) AddChildRule(path string, rule *Rule) error {
	existingRule, err := rules.FindRule(path + "/" + rule.Name)
	if err == nil {
		path = path + "/" + rule.Name

		existingRule.Name = rule.Name

		for _, criteria := range rule.Criteria {
			err := rules.AddCriteriaOptions(path+"/"+criteria.Name, *criteria.Options)
			if err != nil {
				return err
			}
		}

		for _, behavior := range rule.Behaviors {
			err := rules.AddBehaviorOptions(path+"/"+behavior.Name, *behavior.Options)
			if err != nil {
				return err
			}
		}

		for _, childRule := range rule.Children {
			err := rules.AddChildRule(path, childRule)
			if err != nil {
				return err
			}
		}

		return nil
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Children = append(parent.Children, rule)

	return nil
}

// SetChildRule adds a rule as a child of the rule specified by the given path
//
// If the rule already exists, it is replaced by the given rule.
func (rules *Rules) SetChildRule(path string, rule *Rule) error {
	exists, err := rules.FindRule(path + "/" + rule.Name)
	if err != nil {
		*exists = *rule
		return nil
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Children = append(parent.Children, rule)

	return nil
}

// SetCriteriaOptions sets the options on a given criteria path
//
// path is a / delimited path from the root of the rule set to the criteria.
// All existing options are overwritten. To add/replace options see AddCriteriaOptions
// instead.
func (rules *Rules) SetCriteriaOptions(path string, newOptions OptionValue) error {
	criteria, err := rules.FindCriteria(path)
	if err != nil && err != ErrorMap[ErrCriteriaNotFound] {
		return err
	}

	// Create the missing criteria
	if err != nil && err == ErrorMap[ErrCriteriaNotFound] {
		parent, err := rules.FindParentRule(path)
		if err != nil {
			return err
		}

		sep := "/"
		segments := strings.Split(path, sep)
		criteria = NewCriteria(parent)
		criteria.Name = segments[len(segments)-1]
		rules.AddChildCriteria(strings.Join(segments[:len(segments)-1], sep), criteria)
	}

	criteria.Options = &newOptions
	return nil
}

// AddCriteriaOptions adds/replaces options on a given Criteria path
//
// path is a / delimited path from the root of the rule set to the criteria.
// Individual existing options are overwritten. To replace all options, see SetCriteriaOptions
// instead.
func (rules *Rules) AddCriteriaOptions(path string, newOptions OptionValue) error {
	criteria, err := rules.FindCriteria(path)
	if err != nil && err != ErrorMap[ErrCriteriaNotFound] {
		return err
	}

	// Create the missing criteria
	if err != nil && err == ErrorMap[ErrCriteriaNotFound] {
		parent, err := rules.FindParentRule(path)
		if err != nil {
			return err
		}

		sep := "/"
		segments := strings.Split(path, sep)
		criteria = NewCriteria(parent)
		criteria.Name = segments[len(segments)-1]
		rules.AddChildCriteria(strings.Join(segments[:len(segments)-1], sep), criteria)
	}

	options := *criteria.Options
	for key, value := range newOptions {
		options[key] = value
	}
	criteria.Options = &options

	return nil
}

// AddChildCriteria adds criteria as a child of the rule specified by the given path
//
// If the criteria already exists it's options are added to the existing criteria.
func (rules *Rules) AddChildCriteria(path string, criteria *Criteria) error {
	_, err := rules.FindCriteria(path + "/" + criteria.Name)
	if err == nil {
		path = path + "/" + criteria.Name

		return rules.AddCriteriaOptions(path, *criteria.Options)
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Criteria = append(parent.Criteria, criteria)
	return nil
}

// SetChildCriteria adds criteria as a child of the rule specified by the given path
//
// If the criteria already exists it will replace the existing criteria.
func (rules *Rules) SetChildCriteria(path string, criteria *Criteria) error {
	existingCriteria, err := rules.FindCriteria(path + "/" + criteria.Name)
	if err == nil {
		*existingCriteria = *criteria
		return nil
	}

	parent, err := rules.FindRule(path)
	if err != nil {
		return err
	}

	parent.Criteria = append(parent.Criteria, criteria)
	return nil
}

// Find the parent rule for a given rule, criteria, or behavior path
func (rules *Rules) FindParentRule(path string) (*Rule, error) {
	sep := "/"
	segments := strings.Split(strings.ToLower(strings.TrimPrefix(path, sep)), sep)
	parentPath := strings.Join(segments[0:len(segments)-1], sep)

	return rules.FindRule(parentPath)
}

// FindBehavior locates a specific behavior by path
//
// See SetBehaviorOptions and AddBehaviorOptions for examples of paths.
func (rules *Rules) FindBehavior(path string) (*Behavior, error) {
	if len(path) <= 1 {
		return nil, ErrorMap[ErrInvalidPath]
	}

	rule, err := rules.FindParentRule(path)
	if err != nil {
		return nil, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	behaviorName := strings.ToLower(segments[len(segments)-1])
	for _, behavior := range rule.Behaviors {
		if strings.ToLower(behavior.Name) == behaviorName {
			return behavior, nil
		}
	}

	return nil, ErrorMap[ErrBehaviorNotFound]
}

// FindCriteria locates a specific Critieria by path
//
// See SetCriteriaOptions and AddCriteriaOptions for examples of paths.
func (rules *Rules) FindCriteria(path string) (*Criteria, error) {
	if len(path) <= 1 {
		return nil, ErrorMap[ErrInvalidPath]
	}

	rule, err := rules.FindParentRule(path)
	if err != nil {
		return nil, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	criteriaName := strings.ToLower(segments[len(segments)-1])
	for _, criteria := range rule.Criteria {
		if strings.ToLower(criteria.Name) == criteriaName {
			return criteria, nil
		}
	}

	return nil, ErrorMap[ErrCriteriaNotFound]
}

// FindRule locates a specific rule by path
//
// See SetBehaviorOptions and AddBehaviorOptions for examples of paths.
func (rules *Rules) FindRule(path string) (*Rule, error) {
	path = rules.fixupPath(path)

	if path == "" {
		return rules.Rules, nil
	}

	sep := "/"
	segments := strings.Split(path, sep)

	currentRule := rules.Rules
	for _, segment := range segments {
		found := false
		for _, rule := range currentRule.GetChildren(0, 1) {
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

func (rules *Rules) fixupPath(path string) string {
	sep := "/"

	path = strings.Replace(path, sep+sep, sep, -1)
	path = strings.TrimSuffix(path, sep)

	if path == "" || path == "/default" {
		return ""
	}

	if strings.HasPrefix(path, sep+"default"+sep) {
		path = "/" + strings.TrimPrefix(path, "/default/")
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.ToLower(path)

	return path
}

// Rule represents a property rule resource
type Rule struct {
	client.Resource
	parent              *Rules
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
func NewRule(parent *Rules) *Rule {
	rule := &Rule{parent: parent, Children: []*Rule{}}
	rule.Init()

	return rule
}

// GetChildren recurses a Rule tree and retrieves all child rules
func (rule *Rule) GetChildren(depth int, limit int) []*Rule {
	depth++

	if limit != 0 && depth > limit {
		return nil
	}

	var children []*Rule
	if len(rule.Children) > 0 {
		for _, v := range rule.Children {
			v.Depth = depth
			children = append(children, v)
			children = append(children, v.GetChildren(depth, limit)...)
		}
	}

	return children
}

// AddChildRule appends a child rule
func (rule *Rule) AddChildRule(child *Rule) {
	rule.Children = append(rule.Children, child)
}

// AddCriteria appends a rule criteria
func (rule *Rule) AddCriteria(criteria *Criteria) {
	rule.Criteria = append(rule.Criteria, criteria)
}

// AddBehavior appends a rule behavior
func (rule *Rule) AddBehavior(behavior *Behavior) {
	rule.Behaviors = append(rule.Behaviors, behavior)
}

// Criteria represents a rule criteria resource
type Criteria struct {
	client.Resource
	parent  *Rule
	Name    string       `json:"name"`
	Options *OptionValue `json:"options"`
}

// NewCriteria creates a new Criteria
func NewCriteria(parent *Rule) *Criteria {
	criteria := &Criteria{parent: parent}
	criteria.Init()

	return criteria
}

// Behavior represents a rule behavior resource
type Behavior struct {
	client.Resource
	parent  *Rule
	Name    string       `json:"name"`
	Options *OptionValue `json:"options"`
}

// NewBehavior creates a new Behavior
func NewBehavior(parent *Rule) *Behavior {
	behavior := &Behavior{parent: parent, Options: &OptionValue{}}
	behavior.Init()

	return behavior
}

// OptionValue represents a generic option value
//
// OptionValue is a map with string keys, and any
// type of value. You can nest OptionValues as necessary
// to create more complex values.
type OptionValue client.JSONBody

// AvailableCriteria represents a collection of available rule criteria
type AvailableCriteria struct {
	client.Resource
	ContractID        string `json:"contractId"`
	GroupID           string `json:"groupId"`
	ProductID         string `json:"productId"`
	RuleFormat        string `json:"ruleFormat"`
	AvailableCriteria struct {
		Items []struct {
			Name       string `json:"name"`
			SchemaLink string `json:"schemaLink"`
		} `json:"items"`
	} `json:"availableCriteria"`
}

// NewAvailableCriteria creates a new AvailableCriteria
func NewAvailableCriteria() *AvailableCriteria {
	availableCriteria := &AvailableCriteria{}
	availableCriteria.Init()

	return availableCriteria
}

// GetAvailableCriteria retrieves criteria available for a given property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listavailablecriteria
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/available-criteria{?contractId,groupId}
func (availableCriteria *AvailableCriteria) GetAvailableCriteria(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/available-criteria?contractId=%s&groupId=%s",
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

	if err = client.BodyJSON(res, availableCriteria); err != nil {
		return err
	}

	return nil
}

// AvailableBehaviors represents a collection of available rule behaviors
type AvailableBehaviors struct {
	client.Resource
	ContractID string `json:"contractId"`
	GroupID    string `json:"groupId"`
	ProductID  string `json:"productId"`
	RuleFormat string `json:"ruleFormat"`
	Behaviors  struct {
		Items []AvailableBehavior `json:"items"`
	} `json:"behaviors"`
}

// NewAvailableBehaviors creates a new AvailableBehaviors
func NewAvailableBehaviors() *AvailableBehaviors {
	availableBehaviors := &AvailableBehaviors{}
	availableBehaviors.Init()

	return availableBehaviors
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (availableBehaviors *AvailableBehaviors) PostUnmarshalJSON() error {
	availableBehaviors.Init()

	for key := range availableBehaviors.Behaviors.Items {
		availableBehaviors.Behaviors.Items[key].parent = availableBehaviors
	}

	availableBehaviors.Complete <- true

	return nil
}

// GetAvailableBehaviors retrieves available behaviors for a given property
//
// See: Property.GetAvailableBehaviors
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listavailablebehaviors
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/available-behaviors{?contractId,groupId}
func (availableBehaviors *AvailableBehaviors) GetAvailableBehaviors(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/available-behaviors?contractId=%s&groupId=%s",
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

	if err = client.BodyJSON(res, availableBehaviors); err != nil {
		return err
	}

	return nil
}

// AvailableBehavior represents an available behavior resource
type AvailableBehavior struct {
	client.Resource
	parent     *AvailableBehaviors
	Name       string `json:"name"`
	SchemaLink string `json:"schemaLink"`
}

// NewAvailableBehavior creates a new AvailableBehavior
func NewAvailableBehavior(parent *AvailableBehaviors) *AvailableBehavior {
	availableBehavior := &AvailableBehavior{parent: parent}
	availableBehavior.Init()

	return availableBehavior
}

// GetSchema retrieves the JSON schema for an available behavior
func (behavior *AvailableBehavior) GetSchema() (*gojsonschema.Schema, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		behavior.SchemaLink,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	schemaBytes, _ := ioutil.ReadAll(res.Body)
	schemaBody := string(schemaBytes)
	loader := gojsonschema.NewStringLoader(schemaBody)
	schema, err := gojsonschema.NewSchema(loader)

	return schema, err
}

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
