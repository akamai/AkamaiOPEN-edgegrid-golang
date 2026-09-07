package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PolicyVersion is response returned by GetPolicyVersion, CreatePolicyVersion or UpdatePolicyVersion
	PolicyVersion struct {
		Location         string             `json:"location"`
		RevisionID       int64              `json:"revisionId"`
		PolicyID         int64              `json:"policyId"`
		Version          int64              `json:"version"`
		Description      string             `json:"description"`
		CreatedBy        string             `json:"createdBy"`
		CreateDate       int64              `json:"createDate"`
		LastModifiedBy   string             `json:"lastModifiedBy"`
		LastModifiedDate int64              `json:"lastModifiedDate"`
		RulesLocked      bool               `json:"rulesLocked"`
		Activations      []PolicyActivation `json:"activations"`
		MatchRules       MatchRules         `json:"matchRules"`
		MatchRuleFormat  MatchRuleFormat    `json:"matchRuleFormat"`
		Deleted          bool               `json:"deleted,omitempty"`
		Warnings         []Warning          `json:"warnings,omitempty"`
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

	// DeletePolicyVersionRequest describes the parameters of the delete policy version request
	DeletePolicyVersionRequest struct {
		PolicyID int64
		Version  int64
	}

	// UpdatePolicyVersionRequest describes the parameters of the update policy version request
	UpdatePolicyVersionRequest struct {
		UpdatePolicyVersion
		PolicyID int64
		Version  int64
	}
)

// Validate validates ListPolicyVersionsRequest
func (c ListPolicyVersionsRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID": validation.Validate(c.PolicyID, validation.Required),
		"Offset":   validation.Validate(c.Offset, validation.Min(0)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates CreatePolicyVersionRequest
func (c CreatePolicyVersionRequest) Validate() error {
	errs := validation.Errors{
		"Description": validation.Validate(c.Description, validation.Length(0, 255)),
		"MatchRuleFormat": validation.Validate(c.MatchRuleFormat, validation.In(MatchRuleFormat10).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '1.0' or '' (empty)", (&c).MatchRuleFormat))),
		"MatchRules": validation.Validate(c.MatchRules, validation.Length(0, 5000)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpdatePolicyVersionRequest
func (o UpdatePolicyVersionRequest) Validate() error {
	errs := validation.Errors{
		"Description": validation.Validate(o.Description, validation.Length(0, 255)),
		"MatchRuleFormat": validation.Validate(o.MatchRuleFormat, validation.In(MatchRuleFormat10).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '1.0' or '' (empty)", (&o).MatchRuleFormat))),
		"MatchRules": validation.Validate(o.MatchRules, validation.Length(0, 5000)),
	}
	return edgegriderr.ParseValidationErrors(errs)
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

func (c *cloudlets) ListPolicyVersions(ctx context.Context, params ListPolicyVersionsRequest) ([]PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("ListPolicyVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListPolicyVersions, ErrStructValidation, err)
	}

	req, err := request.NewGet(ctx, "/cloudlets/api/v2/policies/%d/versions", params.PolicyID).
		AddQueryParam("offset", fmt.Sprintf("%d", params.Offset)).
		AddQueryParam("includeRules", strconv.FormatBool(params.IncludeRules)).
		AddQueryParam("includeDeleted", strconv.FormatBool(params.IncludeDeleted)).
		AddQueryParam("includeActivations", strconv.FormatBool(params.IncludeActivations)).
		AddQueryParamFunc("pageSize", func() string {
			return fmt.Sprintf("%d", *params.PageSize)
		}, params.PageSize != nil).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicyVersions, err)
	}

	var result []PolicyVersion
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicyVersions, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicyVersions, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) GetPolicyVersion(ctx context.Context, params GetPolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicyVersion")

	var result PolicyVersion

	req, err := request.NewGet(ctx, "/cloudlets/api/v2/policies/%d/versions/%d", params.PolicyID, params.Version).
		AddQueryParam("omitRules", strconv.FormatBool(params.OmitRules)).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyVersion, err)
	}

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreatePolicyVersion(ctx context.Context, params CreatePolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("CreatePolicyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreatePolicyVersion, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/cloudlets/api/v2/policies/%d/versions", params.PolicyID).
		WithBody(params.CreatePolicyVersion).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePolicyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicyVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) DeletePolicyVersion(ctx context.Context, params DeletePolicyVersionRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeletePolicyVersion")

	req, err := request.NewDelete(ctx, "/cloudlets/api/v2/policies/%d/versions/%d", params.PolicyID, params.Version).Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeletePolicyVersion, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeletePolicyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeletePolicyVersion, c.Error(resp))
	}

	return nil
}

func (c *cloudlets) UpdatePolicyVersion(ctx context.Context, params UpdatePolicyVersionRequest) (*PolicyVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdatePolicyVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdatePolicyVersion, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/cloudlets/api/v2/policies/%d/versions/%d", params.PolicyID, params.Version).
		WithBody(params.UpdatePolicyVersion).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicyVersion, c.Error(resp))
	}

	return &result, nil
}
