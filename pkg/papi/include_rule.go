package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// IncludeRules contains operations available on IncludeRule resource
	IncludeRules interface {
		// GetIncludeRuleTree gets the entire rule tree for an include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/get-include-version-rules
		GetIncludeRuleTree(context.Context, GetIncludeRuleTreeRequest) (*GetIncludeRuleTreeResponse, error)

		// UpdateIncludeRuleTree updates the rule tree for an include version
		//
		// See: https://techdocs.akamai.com/property-mgr/reference/patch-include-version-rules
		UpdateIncludeRuleTree(context.Context, UpdateIncludeRuleTreeRequest) (*UpdateIncludeRuleTreeResponse, error)
	}

	// GetIncludeRuleTreeRequest contains path and query params necessary to perform GetIncludeRuleTree
	GetIncludeRuleTreeRequest struct {
		ContractID     string
		GroupID        string
		IncludeID      string
		IncludeVersion int
		RuleFormat     string
		ValidateMode   string
		ValidateRules  bool
	}

	// GetIncludeRuleTreeResponse contains data returned by performing GetIncludeRuleTree request
	GetIncludeRuleTreeResponse struct {
		Response
		Comments       string      `json:"comments,omitempty"`
		Etag           string      `json:"etag"`
		IncludeID      string      `json:"includeId"`
		IncludeName    string      `json:"includeName"`
		IncludeType    IncludeType `json:"includeType"`
		IncludeVersion int         `json:"includeVersion"`
		RuleFormat     string      `json:"ruleFormat"`
		Rules          Rules       `json:"rules"`
	}

	// UpdateIncludeRuleTreeRequest contains path and query params, as well as request body necessary to perform UpdateIncludeRuleTree
	UpdateIncludeRuleTreeRequest struct {
		ContractID     string
		DryRun         bool
		GroupID        string
		IncludeID      string
		IncludeVersion int
		Rules          RulesUpdate
		ValidateMode   string
		ValidateRules  bool
	}

	// UpdateIncludeRuleTreeResponse contains data returned by performing UpdateIncludeRuleTree request
	UpdateIncludeRuleTreeResponse struct {
		Response
		ResponseHeaders UpdateIncludeResponseHeaders
		Comments        string      `json:"comments,omitempty"`
		Etag            string      `json:"etag"`
		IncludeID       string      `json:"includeId"`
		IncludeName     string      `json:"includeName"`
		IncludeType     IncludeType `json:"includeType"`
		IncludeVersion  int         `json:"includeVersion"`
		RuleFormat      string      `json:"ruleFormat"`
		Rules           Rules       `json:"rules"`
	}

	// UpdateIncludeResponseHeaders contains information received in response headers when making a UpdateIncludeRuleTree request
	UpdateIncludeResponseHeaders struct {
		ElementsPerPropertyRemaining      string
		ElementsPerPropertyTotal          string
		MaxNestedRulesPerIncludeRemaining string
		MaxNestedRulesPerIncludeTotal     string
	}
)

// Validate validates GetIncludeRuleTreeRequest struct
func (i GetIncludeRuleTreeRequest) Validate() error {
	errs := validation.Errors{
		"ContractID":     validation.Validate(i.ContractID, validation.Required),
		"GroupID":        validation.Validate(i.GroupID, validation.Required),
		"IncludeID":      validation.Validate(i.IncludeID, validation.Required),
		"IncludeVersion": validation.Validate(i.IncludeVersion, validation.Required),
		"RuleFormat":     validation.Validate(i.RuleFormat, validation.Match(validRuleFormat)),
		"ValidateMode":   validation.Validate(i.ValidateMode, validation.In(RuleValidateModeFast, RuleValidateModeFull)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpdateIncludeRuleTreeRequest struct
func (i UpdateIncludeRuleTreeRequest) Validate() error {
	errs := validation.Errors{
		"ContractID":     validation.Validate(i.ContractID, validation.Required),
		"GroupID":        validation.Validate(i.GroupID, validation.Required),
		"IncludeID":      validation.Validate(i.IncludeID, validation.Required),
		"IncludeVersion": validation.Validate(i.IncludeVersion, validation.Required),
		"Rules":          validation.Validate(i.Rules),
		"ValidateMode":   validation.Validate(i.ValidateMode, validation.In(RuleValidateModeFast, RuleValidateModeFull)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

var (
	// ErrGetIncludeRuleTree represents error when fetching rule tree fails
	ErrGetIncludeRuleTree = errors.New("fetching include rule tree")
	// ErrUpdateIncludeRuleTree represents error when updating rule tree fails
	ErrUpdateIncludeRuleTree = errors.New("updating include rule tree")
)

func (p *papi) GetIncludeRuleTree(ctx context.Context, params GetIncludeRuleTreeRequest) (*GetIncludeRuleTreeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIncludeRuleTree")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetIncludeRuleTree, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/versions/%d/rules", params.IncludeID, params.IncludeVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetIncludeRuleTree, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	if params.ValidateMode != "" {
		q.Add("validateMode", params.ValidateMode)
	}
	if !params.ValidateRules {
		q.Add("validateRules", strconv.FormatBool(params.ValidateRules))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetIncludeRuleTree, err)
	}

	if params.RuleFormat != "" {
		req.Header.Set("Accept", fmt.Sprintf("application/vnd.akamai.papirules.%s+json", params.RuleFormat))
	}

	var result GetIncludeRuleTreeResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetIncludeRuleTree, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetIncludeRuleTree, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) UpdateIncludeRuleTree(ctx context.Context, params UpdateIncludeRuleTreeRequest) (*UpdateIncludeRuleTreeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateIncludeRuleTree")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateIncludeRuleTree, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/versions/%d/rules", params.IncludeID, params.IncludeVersion))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateIncludeRuleTree, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	if params.ValidateMode != "" {
		q.Add("validateMode", params.ValidateMode)
	}
	if !params.ValidateRules {
		q.Add("validateRules", strconv.FormatBool(params.ValidateRules))
	}
	if params.DryRun {
		q.Add("dryRun", strconv.FormatBool(params.DryRun))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateIncludeRuleTree, err)
	}

	var result UpdateIncludeRuleTreeResponse
	resp, err := p.Exec(req, &result, params.Rules)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateIncludeRuleTree, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateIncludeRuleTree, p.Error(resp))
	}

	result.ResponseHeaders.ElementsPerPropertyRemaining = resp.Header.Get("x-limit-elements-per-property-remaining")
	result.ResponseHeaders.ElementsPerPropertyTotal = resp.Header.Get("x-limit-elements-per-property-limit")
	result.ResponseHeaders.MaxNestedRulesPerIncludeRemaining = resp.Header.Get("x-limit-max-nested-rules-per-include-remaining")
	result.ResponseHeaders.MaxNestedRulesPerIncludeTotal = resp.Header.Get("x-limit-max-nested-rules-per-include-limit")

	return &result, nil
}
