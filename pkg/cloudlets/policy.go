package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Policy is response returned by GetPolicy or UpdatePolicy
	Policy struct {
		Location         string             `json:"location"`
		PolicyID         int64              `json:"policyId"`
		GroupID          int64              `json:"groupId"`
		Name             string             `json:"name"`
		Description      string             `json:"description"`
		CreatedBy        string             `json:"createdBy"`
		CreateDate       float64            `json:"createDate"`
		LastModifiedBy   string             `json:"lastModifiedBy"`
		LastModifiedDate float64            `json:"lastModifiedDate"`
		Activations      []PolicyActivation `json:"activations"`
		CloudletID       int64              `json:"cloudletId"`
		CloudletCode     string             `json:"cloudletCode"`
		APIVersion       string             `json:"apiVersion"`
		Deleted          bool               `json:"deleted"`
	}

	// PolicyActivation represents a policy activation resource
	PolicyActivation struct {
		APIVersion   string                  `json:"apiVersion"`
		Network      PolicyActivationNetwork `json:"network"`
		PolicyInfo   PolicyInfo              `json:"policyInfo"`
		PropertyInfo PropertyInfo            `json:"propertyInfo"`
	}

	// PolicyInfo represents a policy info resource
	PolicyInfo struct {
		PolicyID       int64                  `json:"policyId"`
		Name           string                 `json:"name"`
		Version        int64                  `json:"version"`
		Status         PolicyActivationStatus `json:"status"`
		StatusDetail   string                 `json:"statusDetail,omitempty"`
		ActivatedBy    string                 `json:"activatedBy"`
		ActivationDate int64                  `json:"activationDate"`
	}

	// PropertyInfo represents a property info resource
	PropertyInfo struct {
		Name           string                 `json:"name"`
		Version        int64                  `json:"version"`
		GroupID        int64                  `json:"groupId"`
		Status         PolicyActivationStatus `json:"status"`
		ActivatedBy    string                 `json:"activatedBy"`
		ActivationDate int64                  `json:"activationDate"`
	}

	// PolicyActivationStatus is an activation status type for policy
	PolicyActivationStatus string

	// GetPolicyRequest describes the body of the get policy request
	GetPolicyRequest struct {
		PolicyID int64
	}

	// CreatePolicyRequest describes the body of the create policy request
	CreatePolicyRequest struct {
		Name         string `json:"name"`
		CloudletID   int64  `json:"cloudletId"`
		Description  string `json:"description,omitempty"`
		PropertyName string `json:"propertyName,omitempty"`
		GroupID      int64  `json:"groupId,omitempty"`
	}

	// UpdatePolicy describes the body of the update policy request
	UpdatePolicy struct {
		Name         string `json:"name,omitempty"`
		Description  string `json:"description,omitempty"`
		PropertyName string `json:"propertyName,omitempty"`
		GroupID      int64  `json:"groupId,omitempty"`
		Deleted      bool   `json:"deleted,omitempty"`
	}

	// ListPoliciesRequest describes the parameters for the list policies request
	ListPoliciesRequest struct {
		CloudletID     *int64
		IncludeDeleted bool
		Offset         int
		PageSize       *int
	}

	// UpdatePolicyRequest describes the parameters for the update policy request
	UpdatePolicyRequest struct {
		UpdatePolicy
		PolicyID int64
	}

	// RemovePolicyRequest describes the body of the remove policy request
	RemovePolicyRequest struct {
		PolicyID int64
	}
)

const (
	// PolicyActivationStatusActive is an activation that is currently active
	PolicyActivationStatusActive PolicyActivationStatus = "active"
	// PolicyActivationStatusDeactivated is an activation that is deactivated
	PolicyActivationStatusDeactivated PolicyActivationStatus = "deactivated"
	// PolicyActivationStatusInactive is an activation that is not active
	PolicyActivationStatusInactive PolicyActivationStatus = "inactive"
	// PolicyActivationStatusPending is status of a pending activation
	PolicyActivationStatusPending PolicyActivationStatus = "pending"
	// PolicyActivationStatusFailed is status of a failed activation
	PolicyActivationStatusFailed PolicyActivationStatus = "failed"
)

var nameRegexp = regexp.MustCompile("^[a-z_A-Z0-9]+$")
var propertyNameRegexp = regexp.MustCompile("^[a-z_A-Z0-9.\\-]+$")

// Validate validates CreatePolicyRequest
func (v CreatePolicyRequest) Validate() error {
	return validation.Errors{
		"Name":         validation.Validate(v.Name, validation.Required, validation.Length(0, 64), validation.Match(nameRegexp)),
		"PropertyName": validation.Validate(v.PropertyName, validation.Match(propertyNameRegexp)),
		"CloudletID":   validation.Validate(v.CloudletID, validation.Min(0), validation.Max(13)),
		"Description":  validation.Validate(v.Description, validation.Length(0, 255)),
		"GroupID":      validation.Validate(v.GroupID),
	}.Filter()
}

// Validate validates UpdatePolicyRequest
func (v UpdatePolicyRequest) Validate() error {
	return validation.Errors{
		"Name":         validation.Validate(v.Name, validation.Length(0, 64), validation.Match(nameRegexp)),
		"Description":  validation.Validate(v.Description, validation.Length(0, 255)),
		"PropertyName": validation.Validate(v.PropertyName, validation.Match(propertyNameRegexp)),
		"GroupID":      validation.Validate(v.GroupID),
		"Deleted":      validation.Validate(v.Deleted),
	}.Filter()
}

var (
	// ErrListPolicies is returned when ListPolicies fails
	ErrListPolicies = errors.New("list policies")
	// ErrGetPolicy is returned when GetPolicy fails
	ErrGetPolicy = errors.New("get policy")
	// ErrCreatePolicy is returned when CreatePolicy fails
	ErrCreatePolicy = errors.New("create policy")
	// ErrRemovePolicy is returned when RemovePolicy fails
	ErrRemovePolicy = errors.New("remove policy")
	// ErrUpdatePolicy is returned when UpdatePolicy fails
	ErrUpdatePolicy = errors.New("update policy")
)

func (c *cloudlets) ListPolicies(ctx context.Context, params ListPoliciesRequest) ([]Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("ListPolicies")

	uri, err := url.Parse("/cloudlets/api/v2/policies")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListPolicies, err)
	}

	q := uri.Query()
	if params.CloudletID != nil {
		q.Add("cloudletId", fmt.Sprintf("%d", *params.CloudletID))
	}
	if params.PageSize != nil {
		q.Add("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}
	q.Add("offset", fmt.Sprintf("%d", params.Offset))
	q.Add("includeDeleted", strconv.FormatBool(params.IncludeDeleted))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicies, err)
	}

	var result []Policy
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicies, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicies, c.Error(resp))
	}

	return result, nil
}

func (c *cloudlets) GetPolicy(ctx context.Context, params GetPolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicy")

	var result Policy

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d", params.PolicyID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetPolicy, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicy, err)
	}

	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicy, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) CreatePolicy(ctx context.Context, params CreatePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("CreatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreatePolicy, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloudlets/api/v2/policies")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreatePolicy, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreatePolicy, err)
	}

	var result Policy

	resp, err := c.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreatePolicy, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) RemovePolicy(ctx context.Context, params RemovePolicyRequest) error {
	logger := c.Log(ctx)
	logger.Debug("RemovePolicy")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d", params.PolicyID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrRemovePolicy, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrRemovePolicy, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrRemovePolicy, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrRemovePolicy, c.Error(resp))
	}

	return nil
}

func (c *cloudlets) UpdatePolicy(ctx context.Context, params UpdatePolicyRequest) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdatePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePolicy, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/cloudlets/api/v2/policies/%d",
		params.PolicyID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdatePolicy, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePolicy, err)
	}

	var result Policy

	resp, err := c.Exec(req, &result, params.UpdatePolicy)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePolicy, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicy, c.Error(resp))
	}

	return &result, nil
}
