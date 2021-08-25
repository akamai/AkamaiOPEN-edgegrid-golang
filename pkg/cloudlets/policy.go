package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Policies is a cloudlets policies API interface
	Policies interface {
		// GetPolicy gets policy by policyID
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getpolicy
		GetPolicy(context.Context, int64) (*Policy, error)

		// CreatePolicy creates policy
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postpolicies
		CreatePolicy(context.Context, CreatePolicyRequest) (*Policy, error)

		// RemovePolicy removes policy
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#deletepolicy
		RemovePolicy(context.Context, int64) error

		// UpdatePolicy updates policy
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#putpolicy
		UpdatePolicy(context.Context, UpdatePolicyRequest) (*Policy, error)
	}

	// Policy is response returned by GetPolicy or UpdatePolicy
	Policy struct {
		Location         string       `json:"location"`
		PolicyID         int64        `json:"policyId"`
		GroupID          int64        `json:"groupId"`
		Name             string       `json:"name"`
		Description      string       `json:"description"`
		CreatedBy        string       `json:"createdBy"`
		CreateDate       float64      `json:"createDate"`
		LastModifiedBy   string       `json:"lastModifiedBy"`
		LastModifiedDate float64      `json:"lastModifiedDate"`
		Activations      []Activation `json:"activations"`
		CloudletId       int64        `json:"cloudletId"`
		CloudletCode     string       `json:"cloudletCode"`
		APIVersion       string       `json:"apiVersion"`
		Deleted          bool         `json:"deleted"`
	}

	// Activation represents a policy activation resource
	Activation struct {
		APIVersion   string       `json:"apiVersion"`
		Network      string       `json:"network"`
		PolicyInfo   PolicyInfo   `json:"policyInfo"`
		PropertyInfo PropertyInfo `json:"propertyInfo"`
	}

	// PolicyInfo represents a policy info resource
	PolicyInfo struct {
		PolicyID       int64  `json:"policyId"`
		Name           string `json:"name"`
		Version        int64  `json:"version"`
		Status         Status `json:"status"`
		StatusDetail   string `json:"statusDetail,omitempty"`
		ActivatedBy    string `json:"activatedBy"`
		ActivationDate int64  `json:"activationDate"`
	}

	// PropertyInfo represents a property info resource
	PropertyInfo struct {
		Name           string `json:"name"`
		Version        int64  `json:"version"`
		GroupID        int64  `json:"groupId"`
		Status         Status `json:"status"`
		ActivatedBy    string `json:"activatedBy"`
		ActivationDate int64  `json:"activationDate"`
	}

	// Status is an activation status value
	Status string

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

	// UpdatePolicyRequest describes the parameters for the update policy request
	UpdatePolicyRequest struct {
		UpdatePolicy
		PolicyID int64
	}
)

const (
	// StatusActive represents active value
	StatusActive Status = "active"
	// StatusInactive represents inactive value
	StatusInactive Status = "inactive"
	// StatusPending represents pending value
	StatusPending Status = "pending"
	// StatusFailed represents failed value
	StatusFailed Status = "failed"
	// StatusDeactivated represents deactivated value
	StatusDeactivated Status = "deactivated"
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
	// ErrGetPolicy is returned when GetPolicy fails
	ErrGetPolicy = errors.New("get policy")
	// ErrCreatePolicy is returned when CreatePolicy fails
	ErrCreatePolicy = errors.New("create policy")
	// ErrRemovePolicy is returned when RemovePolicy fails
	ErrRemovePolicy = errors.New("remove policy")
	// ErrUpdatePolicy is returned when UpdatePolicy fails
	ErrUpdatePolicy = errors.New("update policy")
)

func (c *cloudlets) GetPolicy(ctx context.Context, policyID int64) (*Policy, error) {
	logger := c.Log(ctx)
	logger.Debug("GetPolicy")

	var result Policy

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d", policyID))
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

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies"))
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

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreatePolicy, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) RemovePolicy(ctx context.Context, policyID int64) error {
	logger := c.Log(ctx)
	logger.Debug("RemovePolicy")

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/policies/%d", policyID))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePolicy, c.Error(resp))
	}

	return &result, nil
}
