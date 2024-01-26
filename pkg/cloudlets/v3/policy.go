package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPoliciesRequest contains request parameters for ListPolicies
	ListPoliciesRequest struct {
		Page int
		Size int
	}

	// CreatePolicyRequest contains request parameters for CreatePolicy
	CreatePolicyRequest struct {
		CloudletType CloudletType `json:"cloudletType"`
		Description  *string      `json:"description,omitempty"`
		GroupID      int64        `json:"groupId"`
		Name         string       `json:"name"`
		PolicyType   PolicyType   `json:"policyType,omitempty"`
	}

	// DeletePolicyRequest contains request parameters for DeletePolicy
	DeletePolicyRequest struct {
		PolicyID int64
	}

	// GetPolicyRequest contains request parameters for GetPolicy
	GetPolicyRequest struct {
		PolicyID int64
	}

	// UpdatePolicyRequest contains request parameters for UpdatePolicy
	UpdatePolicyRequest struct {
		PolicyID   int64
		BodyParams UpdatePolicyBodyParams
	}

	// ClonePolicyRequest contains request parameters for ClonePolicy
	ClonePolicyRequest struct {
		PolicyID   int64
		BodyParams ClonePolicyBodyParams
	}

	// ClonePolicyBodyParams contains request body parameters used in ClonePolicy operation
	// GroupID is required only when cloning v2
	ClonePolicyBodyParams struct {
		AdditionalVersions []int64 `json:"additionalVersions,omitempty"`
		GroupID            int64   `json:"groupId,omitempty"`
		NewName            string  `json:"newName"`
	}

	// UpdatePolicyBodyParams contains request body parameters used in UpdatePolicy operation
	UpdatePolicyBodyParams struct {
		GroupID     int64   `json:"groupId"`
		Description *string `json:"description,omitempty"`
	}

	// PolicyType represents the type of the policy
	PolicyType string

	// CloudletType represents the type of the cloudlet
	CloudletType string

	// ListPoliciesResponse contains the response data from ListPolicies operation
	ListPoliciesResponse struct {
		Content []Policy `json:"content"`
		Links   []Link   `json:"links"`
		Page    Page     `json:"page"`
	}

	// Policy contains information about shared policy
	Policy struct {
		CloudletType       CloudletType       `json:"cloudletType"`
		CreatedBy          string             `json:"createdBy"`
		CreatedDate        time.Time          `json:"createdDate"`
		CurrentActivations CurrentActivations `json:"currentActivations"`
		Description        *string            `json:"description"`
		GroupID            int64              `json:"groupId"`
		ID                 int64              `json:"id"`
		Links              []Link             `json:"links"`
		ModifiedBy         string             `json:"modifiedBy"`
		ModifiedDate       *time.Time         `json:"modifiedDate,omitempty"`
		Name               string             `json:"name"`
		PolicyType         PolicyType         `json:"policyType"`
	}

	// CurrentActivations contains information about the active policy version that's currently in use and the status of the most recent activation
	// or deactivation operation on the policy's versions for the production and staging networks
	CurrentActivations struct {
		Production ActivationInfo `json:"production"`
		Staging    ActivationInfo `json:"staging"`
	}

	// ActivationInfo contains information about effective and latest activations
	ActivationInfo struct {
		Effective *PolicyActivation `json:"effective"`
		Latest    *PolicyActivation `json:"latest"`
	}
)

const (
	// PolicyTypeShared represents policy of type SHARED
	PolicyTypeShared = PolicyType("SHARED")
	// CloudletTypeAP represents cloudlet of type AP
	CloudletTypeAP = CloudletType("AP")
	// CloudletTypeAS represents cloudlet of type AS
	CloudletTypeAS = CloudletType("AS")
	// CloudletTypeCD represents cloudlet of type CD
	CloudletTypeCD = CloudletType("CD")
	// CloudletTypeER represents cloudlet of type ER
	CloudletTypeER = CloudletType("ER")
	// CloudletTypeFR represents cloudlet of type FR
	CloudletTypeFR = CloudletType("FR")
	// CloudletTypeIG represents cloudlet of type IG
	CloudletTypeIG = CloudletType("IG")
)

var (
	// ErrListPolicies is returned when ListPolicies fails
	ErrListPolicies = errors.New("list shared policies")
	// ErrCreatePolicy is returned when CreatePolicy fails
	ErrCreatePolicy = errors.New("create shared policy")
	// ErrDeletePolicy is returned when DeletePolicy fails
	ErrDeletePolicy = errors.New("delete shared policy")
	// ErrGetPolicy is returned when GetPolicy fails
	ErrGetPolicy = errors.New("get shared policy")
	// ErrUpdatePolicy is returned when UpdatePolicy fails
	ErrUpdatePolicy = errors.New("update shared policy")
	// ErrClonePolicy is returned when ClonePolicy fails
	ErrClonePolicy = errors.New("clone policy")
)

// Validate validates ListPoliciesRequest
func (r ListPoliciesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Page": validation.Validate(r.Page, validation.Min(0)),
		"Size": validation.Validate(r.Size, validation.Min(10)),
	})
}

// Validate validates CreatePolicyRequest
func (r CreatePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CloudletType": validation.Validate(r.CloudletType, validation.Required, validation.In(CloudletTypeAP, CloudletTypeAS, CloudletTypeCD, CloudletTypeER, CloudletTypeFR, CloudletTypeIG).
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s'", r.CloudletType, CloudletTypeAP, CloudletTypeAS, CloudletTypeCD, CloudletTypeER, CloudletTypeFR, CloudletTypeIG))),
		"Name": validation.Validate(r.Name, validation.Required, validation.Length(0, 64), validation.Match(regexp.MustCompile("^[a-z_A-Z0-9]+$")).
			Error(fmt.Sprintf("value '%s' is invalid. Must be of format: ^[a-z_A-Z0-9]+$", r.Name))),
		"GroupID":     validation.Validate(r.GroupID, validation.Required),
		"Description": validation.Validate(r.Description, validation.Length(0, 255)),
		"PolicyType":  validation.Validate(r.PolicyType, validation.In(PolicyTypeShared).Error(fmt.Sprintf("value '%s' is invalid. Must be '%s'", r.PolicyType, PolicyTypeShared))),
	})
}

// Validate validates DeletePolicyRequest
func (r DeletePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
	})
}

// Validate validates GetPolicyRequest
func (r GetPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
	})
}

// Validate validates UpdatePolicyRequest
func (r UpdatePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":   validation.Validate(r.PolicyID, validation.Required),
		"BodyParams": validation.Validate(r.BodyParams, validation.Required),
	})
}

// Validate validates UpdatePolicyBodyParams
func (b UpdatePolicyBodyParams) Validate() error {
	return validation.Errors{
		"GroupID":     validation.Validate(b.GroupID, validation.Required),
		"Description": validation.Validate(b.Description, validation.Length(0, 255)),
	}.Filter()
}

// Validate validates ClonePolicyRequest
func (r ClonePolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":   validation.Validate(r.PolicyID, validation.Required),
		"BodyParams": validation.Validate(r.BodyParams, validation.Required),
	})
}

// Validate validates ClonePolicyBodyParams
func (b ClonePolicyBodyParams) Validate() error {
	return validation.Errors{
		"NewName": validation.Validate(b.NewName, validation.Required, validation.Length(0, 64), validation.Match(regexp.MustCompile("^[a-z_A-Z0-9]+$")).
			Error(fmt.Sprintf("value '%s' is invalid. Must be of format: ^[a-z_A-Z0-9]+$", b.NewName))),
	}.Filter()
}

func (c *cloudlets) ListPolicies(ctx context.Context, params ListPoliciesRequest) (*ListPoliciesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListPolicies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListPolicies, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloudlets/v3/policies")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicies, err)
	}

	q := uri.Query()
	if params.Size != 0 {
		q.Add("size", strconv.Itoa(params.Size))
	}
	if params.Page != 0 {
		q.Add("page", strconv.Itoa(params.Page))
	}

	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicies, err)
	}

	var result ListPoliciesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicies, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicies, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreatePolicy(ctx context.Context, params CreatePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("CreatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreatePolicy, ErrStructValidation, err)
	}

	uri := "/cloudlets/v3/policies"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePolicy, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) DeletePolicy(ctx context.Context, params DeletePolicyRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeletePolicy")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeletePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeletePolicy, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeletePolicy, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeletePolicy, c.Error(resp))
	}

	return nil
}

func (c *cloudlets) GetPolicy(ctx context.Context, params GetPolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicy, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPolicy, ErrPolicyNotFound, c.Error(resp))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) UpdatePolicy(ctx context.Context, params UpdatePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result, params.BodyParams)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) ClonePolicy(ctx context.Context, params ClonePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("ClonePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrClonePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d/clone", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrClonePolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result, params.BodyParams)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrClonePolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrClonePolicy, c.Error(resp))
	}

	return &result, nil
}
