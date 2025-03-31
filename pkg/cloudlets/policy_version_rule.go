package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetPolicyVersionRuleRequest describes the parameters needed to get policy version rule
	GetPolicyVersionRuleRequest struct {
		AkaRuleID string
		Version   int64
		PolicyID  int64
	}

	// policyMatchRule is wrapper for MatchRule interface to allow unmarshalling of the response in GetPolicyVersionRule, CreatePolicyVersionRule, UpdatePolicyVersionRule
	policyMatchRule struct {
		MatchRule
	}

	// CreatePolicyVersionRuleRequest describes the parameters needed to create policy version rule
	CreatePolicyVersionRuleRequest struct {
		Version  int64
		PolicyID int64
		Index    int64
		MatchRule
	}

	// UpdatePolicyVersionRuleRequest describes the parameters needed to update policy version rule
	UpdatePolicyVersionRuleRequest struct {
		AkaRuleID string
		Version   int64
		PolicyID  int64
		MatchRule
	}
)

// Validate validates GetPolicyVersionRuleRequest
func (r GetPolicyVersionRuleRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Version":   validation.Validate(r.Version, validation.Required, validation.Min(1)),
		"AkaRuleID": validation.Validate(r.AkaRuleID, validation.Required),
		"PolicyID":  validation.Validate(r.PolicyID, validation.Required),
	})
}

// Validate validates CreatePolicyVersionRuleRequest
func (r CreatePolicyVersionRuleRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Version":  validation.Validate(r.Version, validation.Required, validation.Min(1)),
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
		"Index":    validation.Validate(r.Index, validation.Min(0)),
	})
}

// Validate validates UpdatePolicyVersionRuleRequest
func (r UpdatePolicyVersionRuleRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Version":   validation.Validate(r.Version, validation.Required, validation.Min(1)),
		"AkaRuleID": validation.Validate(r.AkaRuleID, validation.Required),
		"PolicyID":  validation.Validate(r.PolicyID, validation.Required),
	})
}

var (
	// ErrGetPolicyVersionRule is returned when GetPolicyVersionRule fails
	ErrGetPolicyVersionRule = errors.New("get policy version rule")
	// ErrCreatePolicyVersionRule is returned when CreatePolicyVersionRule fails
	ErrCreatePolicyVersionRule = errors.New("create policy version rule")
	// ErrUpdatePolicyVersionRule is returned when UpdatePolicyVersionRule fails
	ErrUpdatePolicyVersionRule = errors.New("update policy version rule")
)

func (c *cloudlets) GetPolicyVersionRule(ctx context.Context, params GetPolicyVersionRuleRequest) (MatchRule, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicyVersionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetPolicyVersionRule, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d/rules/%s",
		params.PolicyID, params.Version, params.AkaRuleID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPolicyVersionRule, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyVersionRule, err)
	}

	var res policyMatchRule
	resp, err := c.Exec(req, &res)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyVersionRule, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyVersionRule, c.Error(resp))
	}

	return res.MatchRule, nil
}

func (c *cloudlets) CreatePolicyVersionRule(ctx context.Context, params CreatePolicyVersionRuleRequest) (MatchRule, error) {
	logger := c.Log(ctx)
	logger.Debug("CreatePolicyVersionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreatePolicyVersionRule, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d/rules",
		params.PolicyID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreatePolicyVersionRule, err)
	}

	q := uri.Query()
	if params.Index > 0 {
		q.Set("index", strconv.FormatInt(params.Index, 10))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicyVersionRule, err)
	}

	var res policyMatchRule
	resp, err := c.Exec(req, &res, params.MatchRule)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePolicyVersionRule, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicyVersionRule, c.Error(resp))
	}

	return res.MatchRule, nil
}

func (c *cloudlets) UpdatePolicyVersionRule(ctx context.Context, params UpdatePolicyVersionRuleRequest) (MatchRule, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdatePolicyVersionRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdatePolicyVersionRule, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d/rules/%s",
		params.PolicyID, params.Version, params.AkaRuleID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdatePolicyVersionRule, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicyVersionRule, err)
	}

	var res policyMatchRule
	resp, err := c.Exec(req, &res, params.MatchRule)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicyVersionRule, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicyVersionRule, c.Error(resp))
	}

	return res.MatchRule, nil
}
