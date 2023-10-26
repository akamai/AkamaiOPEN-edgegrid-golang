package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListSharedPoliciesRequest contains request parameters for ListSharedPolicies
	ListSharedPoliciesRequest struct {
		Page int
		Size int
	}

	// CreateSharedPolicyRequest contains request parameters for CreateSharedPolicy
	CreateSharedPolicyRequest struct {
		CloudletType CloudletType `json:"cloudletType"`
		Description  *string      `json:"description,omitempty"`
		GroupID      int64        `json:"groupId"`
		Name         string       `json:"name"`
		PolicyType   PolicyType   `json:"policyType,omitempty"`
	}

	// DeleteSharedPolicyRequest contains request parameters for DeleteSharedPolicy
	DeleteSharedPolicyRequest struct {
		PolicyID int64
	}

	// GetSharedPolicyRequest contains request parameters for GetSharedPolicy
	GetSharedPolicyRequest struct {
		PolicyID int64
	}

	// UpdateSharedPolicyRequest contains request parameters for UpdateSharedPolicy
	UpdateSharedPolicyRequest struct {
		PolicyID   int64
		BodyParams UpdateSharedPolicyBodyParams
	}

	// ClonePolicyRequest contains request parameters for ClonePolicy
	ClonePolicyRequest struct {
		PolicyID   int64
		BodyParams ClonePolicyBodyParams
	}

	// ClonePolicyBodyParams contains request body parameters used in ClonePolicy operation
	ClonePolicyBodyParams struct {
		AdditionalVersions []int64 `json:"additionalVersions,omitempty"`
		GroupID            int64   `json:"groupId,omitempty"`
		NewName            string  `json:"newName"`
	}

	// UpdateSharedPolicyBodyParams contains request body parameters used in UpdateSharedPolicy operation
	UpdateSharedPolicyBodyParams struct {
		GroupID     int64   `json:"groupId"`
		Description *string `json:"description,omitempty"`
	}

	// PolicyType represents the type of the policy
	PolicyType string

	// CloudletType represents the type of the cloudlet
	CloudletType string

	// ListSharedPoliciesResponse contains the response data from ListSharedPolicies operation
	ListSharedPoliciesResponse struct {
		Content []Policy `json:"content"`
		Links   []Link   `json:"links"`
		Page    Page     `json:"page"`
	}

	// Policy contains information about shared policy
	Policy struct {
		CloudletType       CloudletType       `json:"cloudletType"`
		CreatedBy          string             `json:"createdBy"`
		CreatedDate        string             `json:"createdDate"`
		CurrentActivations CurrentActivations `json:"currentActivations"`
		Description        *string            `json:"description"`
		GroupID            int64              `json:"groupId"`
		ID                 int64              `json:"id"`
		Links              []Link             `json:"links"`
		ModifiedBy         string             `json:"modifiedBy"`
		ModifiedDate       string             `json:"modifiedDate"`
		Name               string             `json:"name"`
		PolicyType         PolicyType         `json:"policyType"`
	}

	// CurrentActivations contains information about the active policy version that's currently in use and the status of the most recent activation
	// or deactivation operation on the policy's versions for the production and staging networks
	CurrentActivations struct {
		Production Activation `json:"production"`
		Staging    Activation `json:"staging"`
	}

	// Activation contains information about effective and latest activations
	Activation struct {
		Effective ActivationInfo `json:"effective"`
		Latest    ActivationInfo `json:"latest"`
	}

	// ActivationInfo contains information about activation
	ActivationInfo struct {
		CreatedBy            string `json:"createdBy"`
		CreatedDate          string `json:"createdDate"`
		FinishDate           string `json:"finishDate"`
		ID                   int64  `json:"id"`
		Links                []Link `json:"links"`
		Network              string `json:"network"`
		Operation            string `json:"operation"`
		PolicyID             int64  `json:"policyId"`
		PolicyVersion        int64  `json:"policyVersion"`
		PolicyVersionDeleted bool   `json:"policyVersionDeleted"`
		Status               string `json:"status"`
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
	// CloudletTypeVWR represents cloudlet of type VWR
	CloudletTypeVWR = CloudletType("VWR")
)

var (
	// ErrListSharedPolicies is returned when ListSharedPolicies fails
	ErrListSharedPolicies = errors.New("list shared policies")
	// ErrCreateSharedPolicy is returned when CreateSharedPolicy fails
	ErrCreateSharedPolicy = errors.New("create shared policy")
	// ErrDeleteSharedPolicy is returned when DeleteSharedPolicy fails
	ErrDeleteSharedPolicy = errors.New("delete shared policy")
	// ErrGetSharedPolicy is returned when GetSharedPolicy fails
	ErrGetSharedPolicy = errors.New("get shared policy")
	// ErrUpdateSharedPolicy is returned when UpdateSharedPolicy fails
	ErrUpdateSharedPolicy = errors.New("update shared policy")
	// ErrClonePolicy is returned when ClonePolicy fails
	ErrClonePolicy = errors.New("clone policy")
)

// Validate validates ListSharedPoliciesRequest
func (r ListSharedPoliciesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Page": validation.Validate(r.Page, validation.Min(0)),
		"Size": validation.Validate(r.Size, validation.Min(10)),
	})
}

// Validate validates CreateSharedPolicyRequest
func (r CreateSharedPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CloudletType": validation.Validate(r.CloudletType, validation.Required, validation.In(CloudletTypeAP, CloudletTypeAS, CloudletTypeCD, CloudletTypeER, CloudletTypeFR, CloudletTypeIG, CloudletTypeVWR).
			Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s', '%s'", r.CloudletType, CloudletTypeAP, CloudletTypeAS, CloudletTypeCD, CloudletTypeER, CloudletTypeFR, CloudletTypeIG, CloudletTypeVWR))),
		"Name": validation.Validate(r.Name, validation.Required, validation.Length(0, 64), validation.Match(regexp.MustCompile("^[a-z_A-Z0-9]+$")).
			Error(fmt.Sprintf("value '%s' is invalid. Must be of format: ^[a-z_A-Z0-9]+$", r.Name))),
		"GroupID":     validation.Validate(r.GroupID, validation.Required),
		"Description": validation.Validate(r.Description, validation.Length(0, 255)),
		"PolicyType":  validation.Validate(r.PolicyType, validation.In(PolicyTypeShared).Error(fmt.Sprintf("value '%s' is invalid. Must be '%s'", r.PolicyType, PolicyTypeShared))),
	})
}

// Validate validates DeleteSharedPolicyRequest
func (r DeleteSharedPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
	})
}

// Validate validates GetSharedPolicyRequest
func (r GetSharedPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID": validation.Validate(r.PolicyID, validation.Required),
	})
}

// Validate validates UpdateSharedPolicyRequest
func (r UpdateSharedPolicyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"PolicyID":   validation.Validate(r.PolicyID, validation.Required),
		"BodyParams": validation.Validate(r.BodyParams, validation.Required),
	})
}

// Validate validates UpdateSharedPolicyBodyParams
func (b UpdateSharedPolicyBodyParams) Validate() error {
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

func (c *cloudlets) ListSharedPolicies(ctx context.Context, params ListSharedPoliciesRequest) (*ListSharedPoliciesResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListSharedPolicies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListSharedPolicies, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloudlets/v3/policies")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListSharedPolicies, err)
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
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListSharedPolicies, err)
	}

	var result ListSharedPoliciesResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListSharedPolicies, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListSharedPolicies, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreateSharedPolicy(ctx context.Context, params CreateSharedPolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateSharedPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateSharedPolicy, ErrStructValidation, err)
	}

	uri := "/cloudlets/v3/policies"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateSharedPolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateSharedPolicy, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateSharedPolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) DeleteSharedPolicy(ctx context.Context, params DeleteSharedPolicyRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteSharedPolicy")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteSharedPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteSharedPolicy, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteSharedPolicy, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteSharedPolicy, c.Error(resp))
	}

	return nil
}

func (c *cloudlets) GetSharedPolicy(ctx context.Context, params GetSharedPolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("GetSharedPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetSharedPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetSharedPolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetSharedPolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetSharedPolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) UpdateSharedPolicy(ctx context.Context, params UpdateSharedPolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateSharedPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateSharedPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateSharedPolicy, err)
	}

	var result Policy
	resp, err := c.Exec(req, &result, params.BodyParams)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateSharedPolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateSharedPolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) ClonePolicy(ctx context.Context, params ClonePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("ClonePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrClonePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloudlets/v3/policies/%d", params.PolicyID)

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
