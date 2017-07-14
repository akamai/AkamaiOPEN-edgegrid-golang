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

// PrintRules prints a reasonably easy to read tree of all rules and behaviors on a property
func (rules *Rules) PrintRules() error {
	group := NewGroup(NewGroups())
	group.GroupID = rules.GroupID
	group.ContractIDs = []string{rules.ContractID}

	properties, _ := group.GetProperties(nil)
	var property *Property
	for _, property = range properties.Properties.Items {
		if property.PropertyID == rules.PropertyID {
			break
		}
	}

	fmt.Println(property.PropertyName)

	fmt.Println("├── Criteria")
	for _, criteria := range rules.Rules.Criteria {
		fmt.Printf("│   ├── %s\n", criteria.Name)
		i := 0
		for option, value := range *criteria.Options {
			i++
			if i < len(*criteria.Options) {
				fmt.Printf("│   │   ├── %s: %#v\n", option, value)
			} else {
				fmt.Printf("│   │   └── %s: %#v\n", option, value)
			}
		}
	}

	fmt.Println("└── Behaviors")

	prefix := "   │"
	i := 0
	for _, behavior := range rules.Rules.Behaviors {
		i++
		if i < len(rules.Rules.Behaviors) && len(rules.Rules.Children) != 0 {
			fmt.Printf("   ├── Behavior: %s\n", behavior.Name)
		} else {
			fmt.Printf("   └── Behavior: %s\n", behavior.Name)
		}

		j := 0

		for option, value := range *behavior.Options {
			j++
			if i == len(rules.Rules.Behaviors) && len(rules.Rules.Children) == 0 {
				prefix = strings.TrimSuffix(prefix, "│")
			}

			if j < len(*behavior.Options) {
				fmt.Printf("%s   ├── Option: %s: %#v\n", prefix, option, value)
			} else {
				fmt.Printf("%s   └── Option: %s: %#v\n", prefix, option, value)
			}
		}
	}

	if len(rules.Rules.Children) > 0 {
		i := 0
		children := rules.Rules.GetChildren(0, 0)
		for _, child := range children {
			i++
			spacer := strings.TrimSuffix(strings.Repeat(prefix, child.Depth), "│")
			if i < len(children) {
				fmt.Printf("%s├── Section: %s\n", spacer, child.Name)
			} else {
				fmt.Printf("%s└── Section: %s\n", spacer, child.Name)
			}

			spacer = strings.TrimSuffix(strings.Repeat(prefix, child.Depth+1), "│")
			j := 0
			for _, behavior := range child.Behaviors {
				j++
				if j < len(child.Behaviors) {
					fmt.Printf("%s├── Behavior: %s\n", spacer, behavior.Name)
				} else {
					//spacer = strings.TrimSuffix(spacer, "│   ") + "    "
					fmt.Printf("%s└── Behavior: %s\n", spacer, behavior.Name)
				}
				space := strings.TrimSuffix(strings.Repeat(prefix, child.Depth+2), "│")

				fmt.Printf("%s├── Criteria\n", space)
				i := 0
				for _, criteria := range child.Criteria {
					i++
					if i < len(child.Criteria) {
						fmt.Printf("   │%s├── %s\n", space, criteria.Name)
					} else {
						fmt.Printf("   │%s└── %s\n", space, criteria.Name)
					}
					k := 0
					for option, value := range *criteria.Options {
						k++
						if k < len(*criteria.Options) {
							fmt.Printf("   │   │%s├── %s: %#v\n", space, option, value)
						} else {
							fmt.Printf("   │   │%s└── %s: %#v\n", space, option, value)
						}
					}
				}

				k := 0
				for option, value := range *behavior.Options {
					k++
					if k < len(*behavior.Options) {
						fmt.Printf("%s├── Option: %s: %#v\n", space, option, value)
					} else {
						fmt.Printf("%s└── Option: %s: %#v\n", space, option, value)
					}
				}
			}
		}
	}

	return nil
}

// GetRules populates Rules with rule data for a given property
//
// See: Property.GetRules
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletree
// Endpoint: GET /papi/v0/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRules(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v0/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
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
// Endpoint: HEAD /papi/v0/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRulesDigest(property *Property) (string, error) {
	req, err := client.NewRequest(
		Config,
		"HEAD",
		fmt.Sprintf(
			"/papi/v0/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
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
// Endpoint: PUT /papi/v0/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) Save() error {
	rules.Errors = []*RuleErrors{}

	// /papi/v0/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v0/properties/%s/versions/%d/rules/?contractId=%s&groupId=%s",
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
		return fmt.Errorf("there were %d errors. See rules.Errors for details", len(rules.Errors))
	}

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
	if err != nil {
		return err
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
//		"/Performance/JPEG Images/adaptiveImageCompression"",
//		edgegrid.OptionValue{
//			"tier3StandardCompressionValue": 30,
//		},
//	)
//
// This will only change the "tier3StandardCompressionValue" option value.
func (rules *Rules) AddBehaviorOptions(path string, newOptions OptionValue) error {
	behavior, err := rules.FindBehavior(path)
	if err != nil {
		return err
	}

	options := *behavior.Options
	for key, value := range newOptions {
		options[key] = value
	}
	behavior.Options = &options

	return nil
}

// FindBehavior locates a specific behavior by path
//
// See SetBehaviorOptions and AddBehaviorOptions for examples of paths.
func (rules *Rules) FindBehavior(path string) (*Behavior, error) {
	if len(path) <= 1 {
		return nil, fmt.Errorf("Invalid Path: \"%s\"", path)
	}

	sep := "/"
	segments := strings.Split(strings.ToLower(strings.TrimPrefix(path, sep)), sep)

	if len(segments) == 1 {
		for _, behavior := range rules.Rules.Behaviors {
			if strings.ToLower(behavior.Name) == segments[0] {
				return behavior, nil
			}
		}
		return nil, fmt.Errorf("Path not found: \"%s\"", path)
	}

	currentRule := rules.Rules
	i := 0
	for _, segment := range segments {
		i++
		if i < len(segments) {
			for _, rule := range currentRule.GetChildren(0, 1) {
				if strings.ToLower(rule.Name) == segment {
					currentRule = rule
				}
			}
		} else {
			for _, behavior := range currentRule.Behaviors {
				if strings.ToLower(behavior.Name) == segment {
					return behavior, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Path not found: \"%s\"", path)
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
	Comment             string                       `json:"comment,omitempty"`
	CriteriaLocked      bool                         `json:"criteriaLocked,omitempty"`
	CriteriaMustSatisfy RuleCriteriaMustSatisfyValue `json:"criteriaMustSatisfy,omitempty"`
	UUID                string                       `json:"uuid,omitempty"`
	Options             struct {
		IsSecure bool `json:"is_secure,omitempty"`
	} `json:"options,omitempty"`
}

// NewRule creates a new Rule
func NewRule(parent *Rules) *Rule {
	rule := &Rule{parent: parent}
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
func (rule *Rule) AddCriteria(critera *Criteria) {
	rule.Criteria = append(rule.Criteria, critera)
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
	behavior := &Behavior{parent: parent}
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
// Endpoint: GET /papi/v0/properties/{propertyId}/versions/{propertyVersion}/available-criteria{?contractId,groupId}
func (availableCriteria *AvailableCriteria) GetAvailableCriteria(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v0/properties/%s/versions/%d/available-criteria?contractId=%s&groupId=%s",
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
// Endpoint: GET /papi/v0/properties/{propertyId}/versions/{propertyVersion}/available-behaviors{?contractId,groupId}
func (availableBehaviors *AvailableBehaviors) GetAvailableBehaviors(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v0/properties/%s/versions/%d/available-behaviors?contractId=%s&groupId=%s",
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
