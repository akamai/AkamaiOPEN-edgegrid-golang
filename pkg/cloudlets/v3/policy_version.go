package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPolicyVersions is response returned by ListPolicyVersions
	ListPolicyVersions struct {
		PolicyVersions []ListPolicyVersionsItem `json:"content"`
		Links          []Link                   `json:"links"`
		Page           Page                     `json:"page"`
	}

	// ListPolicyVersionsItem is a content struct of ListPolicyVersion response
	ListPolicyVersionsItem struct {
		CreatedBy     string     `json:"createdBy"`
		CreatedDate   time.Time  `json:"createdDate"`
		Description   *string    `json:"description"`
		ID            int64      `json:"id"`
		Immutable     bool       `json:"immutable"`
		Links         []Link     `json:"links"`
		ModifiedBy    string     `json:"modifiedBy"`
		ModifiedDate  *time.Time `json:"modifiedDate,omitempty"`
		PolicyID      int64      `json:"policyId"`
		PolicyVersion int64      `json:"version"`
	}

	// PolicyVersion is response returned by GetPolicyVersion, CreatePolicyVersion or UpdatePolicyVersion
	PolicyVersion struct {
		CreatedBy          string              `json:"createdBy"`
		CreatedDate        time.Time           `json:"createdDate"`
		Description        *string             `json:"description"`
		ID                 int64               `json:"id"`
		Immutable          bool                `json:"immutable"`
		MatchRules         MatchRules          `json:"matchRules"`
		MatchRulesWarnings []MatchRulesWarning `json:"matchRulesWarnings"`
		ModifiedBy         string              `json:"modifiedBy"`
		ModifiedDate       *time.Time          `json:"modifiedDate,omitempty"`
		PolicyID           int64               `json:"policyId"`
		PolicyVersion      int64               `json:"version"`
	}

	// MatchRulesWarning describes the warnings struct
	MatchRulesWarning struct {
		Detail      string `json:"detail"`
		JSONPointer string `json:"jsonPointer,omitempty"`
		Title       string `json:"title"`
		Type        string `json:"type"`
	}

	// ListPolicyVersionsRequest describes the parameters needed to list policy versions
	ListPolicyVersionsRequest struct {
		PolicyID int64
		Page     int
		Size     int
	}

	// GetPolicyVersionRequest describes the parameters needed to get policy version
	GetPolicyVersionRequest struct {
		PolicyID      int64
		PolicyVersion int64
	}

	// CreatePolicyVersionRequest describes the body of the create policy request
	CreatePolicyVersionRequest struct {
		CreatePolicyVersion
		PolicyID int64
	}

	// CreatePolicyVersion describes the body of the create policy request
	CreatePolicyVersion struct {
		Description *string    `json:"description,omitempty"`
		MatchRules  MatchRules `json:"matchRules"`
	}

	// UpdatePolicyVersion describes the body of the update policy version request
	UpdatePolicyVersion struct {
		Description *string    `json:"description,omitempty"`
		MatchRules  MatchRules `json:"matchRules"`
	}

	// DeletePolicyVersionRequest describes the parameters of the delete policy version request
	DeletePolicyVersionRequest struct {
		PolicyID      int64
		PolicyVersion int64
	}

	// UpdatePolicyVersionRequest describes the parameters of the update policy version request
	UpdatePolicyVersionRequest struct {
		UpdatePolicyVersion
		PolicyID      int64
		PolicyVersion int64
	}
)

// Validate validates ListPolicyVersionsRequest
func (c ListPolicyVersionsRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID": validation.Validate(c.PolicyID, validation.Required),
		"Size":     validation.Validate(c.Size, validation.Min(10)),
		"Page":     validation.Validate(c.Page, validation.Min(0)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates CreatePolicyVersionRequest
func (c CreatePolicyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":    validation.Validate(c.PolicyID, validation.Required),
		"Description": validation.Validate(c.Description, validation.Length(0, 255)),
		"MatchRules":  validation.Validate(c.MatchRules, validation.Length(0, 5000)),
	})
}

// Validate validates UpdatePolicyVersionRequest
func (c UpdatePolicyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":      validation.Validate(c.PolicyID, validation.Required),
		"PolicyVersion": validation.Validate(c.PolicyVersion, validation.Required),
		"Description":   validation.Validate(c.Description, validation.Length(0, 255)),
		"MatchRules":    validation.Validate(c.MatchRules, validation.Length(0, 5000)),
	})
}

// Validate validates DeletePolicyVersionRequest
func (c DeletePolicyVersionRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":      validation.Validate(c.PolicyID, validation.Required),
		"PolicyVersion": validation.Validate(c.PolicyVersion, validation.Required),
	})
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

func (c *cloudlets) ListPolicyVersions(ctx context.Context, params ListPolicyVersionsRequest) (*ListPolicyVersions, error) {
	logger := c.Log(ctx)
	logger.Debug("ListPolicyVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListPolicyVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/v3/policies/%d/versions", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicyVersions, err)
	}

	q := uri.Query()
	q.Add("page", fmt.Sprintf("%d", params.Page))
	if params.Size != 0 {
		q.Add("size", fmt.Sprintf("%d", params.Size))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicyVersions, err)
	}

	var result *ListPolicyVersions
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

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/versions/%d", params.PolicyID, params.PolicyVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyVersion, err)
	}

	var result PolicyVersion

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

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/versions", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result, params.CreatePolicyVersion)
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

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w:\n%s", ErrDeletePolicyVersion, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/versions/%d", params.PolicyID, params.PolicyVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
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

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/versions/%d", params.PolicyID, params.PolicyVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicyVersion, err)
	}

	var result PolicyVersion

	resp, err := c.Exec(req, &result, params.UpdatePolicyVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicyVersion, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicyVersion, c.Error(resp))
	}

	return &result, nil
}
