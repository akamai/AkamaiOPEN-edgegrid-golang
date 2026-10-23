package imaging

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/internal/request"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ListPolicySetsRequest describes the parameters of the ListPolicySets request
	ListPolicySetsRequest struct {
		ContractID string
		Network    Network
	}

	// GetPolicySetRequest describes the parameters of the get GetPolicySet request
	GetPolicySetRequest struct {
		PolicySetID string
		ContractID  string
		Network     Network
	}

	// CreatePolicySet describes the body of the CreatePolicySet request
	CreatePolicySet struct {
		Name          string      `json:"name"`
		Region        Region      `json:"region"`
		Type          MediaType   `json:"type"`
		DefaultPolicy PolicyInput `json:"defaultPolicy,omitempty"`
	}

	// CreatePolicySetRequest describes the parameters of the CreatePolicySet request
	CreatePolicySetRequest struct {
		ContractID string
		CreatePolicySet
	}

	// UpdatePolicySet describes the body of the UpdatePolicySet request
	UpdatePolicySet struct {
		Name   string `json:"name"`
		Region Region `json:"region"`
	}

	// UpdatePolicySetRequest describes the parameters of the UpdatePolicySet request
	UpdatePolicySetRequest struct {
		PolicySetID string
		ContractID  string
		UpdatePolicySet
	}

	// PolicySet is a response returned by CRU operations
	PolicySet struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		Region       Region   `json:"region"`
		Type         string   `json:"type"`
		User         string   `json:"user"`
		Properties   []string `json:"properties"`
		LastModified string   `json:"lastModified"`
	}

	// DeletePolicySetRequest describes the parameters of the delete PolicySet request
	DeletePolicySetRequest struct {
		PolicySetID string
		ContractID  string
	}

	// MediaType of media this Policy Set manages
	MediaType string

	// Network represents the network where policy set is stored
	Network string

	// Region represents the geographic region which media using this Policy Set is optimized for
	Region string
)

const (
	// RegionUS represents US region
	RegionUS Region = "US"
	// RegionEMEA represents EMEA region
	RegionEMEA Region = "EMEA"
	// RegionAsia represents Asia region
	RegionAsia Region = "ASIA"
	// RegionAustralia represents Australia region
	RegionAustralia Region = "AUSTRALIA"
	// RegionJapan represents Japan  region
	RegionJapan Region = "JAPAN"
	// RegionChina represents China region
	RegionChina Region = "CHINA"
)

const (
	// TypeImage represents policy set for Images
	TypeImage MediaType = "IMAGE"
	// TypeVideo represents policy set for Videos
	TypeVideo MediaType = "VIDEO"
)

const (
	// NetworkStaging represents staging network
	NetworkStaging Network = "staging"
	// NetworkProduction represents production network
	NetworkProduction Network = "production"
	// NetworkBoth represent both staging and production network at the same time
	NetworkBoth Network = ""
)

var (
	// ErrListPolicySets is returned when ListPolicySets fails
	ErrListPolicySets = errors.New("list policy sets")
	// ErrGetPolicySet is returned when GetPolicySet fails
	ErrGetPolicySet = errors.New("get policy set")
	// ErrCreatePolicySet is returned when CreatePolicySet fails
	ErrCreatePolicySet = errors.New("create policy set")
	// ErrUpdatePolicySet is returned when UpdatePolicySet fails
	ErrUpdatePolicySet = errors.New("update policy set")
	// ErrDeletePolicySet is returned when DeletePolicySet fails
	ErrDeletePolicySet = errors.New("delete policy set")
)

// Validate validates ListPolicySetsRequest
func (v ListPolicySetsRequest) Validate() error {
	errs := validation.Errors{
		"ContractID": validation.Validate(v.ContractID, validation.Required),
		"Network": validation.Validate(v.Network, validation.In(NetworkStaging, NetworkProduction, NetworkBoth).
			Error(fmt.Sprintf("network has to be '%s', '%s' or empty for both networks at the same time", NetworkStaging, NetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates GetPolicySetRequest
func (v GetPolicySetRequest) Validate() error {
	errs := validation.Errors{
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"Network": validation.Validate(v.Network, validation.In(NetworkStaging, NetworkProduction, NetworkBoth).
			Error(fmt.Sprintf("network has to be '%s', '%s' or empty for both networks at the same time", NetworkStaging, NetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates CreatePolicySetRequest
func (v CreatePolicySetRequest) Validate() error {
	errs := validation.Errors{
		"ContractID": validation.Validate(v.ContractID, validation.Required),
		"Name":       validation.Validate(v.Name, validation.Required),
		"Region": validation.Validate(v.Region, validation.Required, validation.In(RegionUS, RegionEMEA, RegionAsia, RegionAustralia, RegionJapan, RegionChina).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s'", v.Region, RegionUS, RegionEMEA, RegionAsia, RegionAustralia, RegionJapan, RegionChina))),
		"Type": validation.Validate(v.Type, validation.Required, validation.In(TypeImage, TypeVideo).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", v.Type, TypeImage, TypeVideo))),
		"DefaultPolicy": validation.Validate(v.DefaultPolicy),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpdatePolicySetRequest
func (v UpdatePolicySetRequest) Validate() error {
	errs := validation.Errors{
		"ContractID": validation.Validate(v.ContractID, validation.Required),
		"Name":       validation.Validate(v.Name, validation.Required),
		"Region": validation.Validate(v.Region, validation.Required, validation.In(RegionUS, RegionEMEA, RegionAsia, RegionAustralia, RegionJapan, RegionChina).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s'", v.Region, RegionUS, RegionEMEA, RegionAsia, RegionAustralia, RegionJapan, RegionChina))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates DeletePolicySetRequest
func (v DeletePolicySetRequest) Validate() error {
	errs := validation.Errors{
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

func (i *imaging) ListPolicySets(ctx context.Context, params ListPolicySetsRequest) ([]PolicySet, error) {
	logger := i.Log(ctx)
	logger.Debug("ListPolicySets")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrListPolicySets, ErrStructValidation, err)
	}

	builder := request.NewGet(ctx, "/imaging/v2/policysets")
	if params.Network != NetworkBoth {
		builder = request.NewGet(ctx, "/imaging/v2/network/%s/policysets", params.Network)
	}

	req, err := builder.AddHeader("Contract", params.ContractID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrListPolicySets, err)
	}

	var result []PolicySet
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrListPolicySets, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrListPolicySets, i.Error(resp))
	}

	return result, nil
}

func (i *imaging) GetPolicySet(ctx context.Context, params GetPolicySetRequest) (*PolicySet, error) {
	logger := i.Log(ctx)
	logger.Debug("GetPolicySet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrGetPolicySet, ErrStructValidation, err)
	}

	builder := request.NewGet(ctx, "/imaging/v2/policysets/%s", params.PolicySetID)
	if params.Network != NetworkBoth {
		builder = request.NewGet(ctx, "/imaging/v2/network/%s/policysets/%s", params.Network, params.PolicySetID)
	}

	req, err := builder.AddHeader("Contract", params.ContractID).Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrGetPolicySet, err)
	}

	var result PolicySet
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrGetPolicySet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrGetPolicySet, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) CreatePolicySet(ctx context.Context, params CreatePolicySetRequest) (*PolicySet, error) {
	logger := i.Log(ctx)
	logger.Debug("CreatePolicySet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrCreatePolicySet, ErrStructValidation, err)
	}

	req, err := request.NewPost(ctx, "/imaging/v2/policysets").
		AddHeader("Contract", params.ContractID).
		WithBody(params.CreatePolicySet).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrCreatePolicySet, err)
	}

	var result PolicySet
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrCreatePolicySet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %w", ErrCreatePolicySet, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) UpdatePolicySet(ctx context.Context, params UpdatePolicySetRequest) (*PolicySet, error) {
	logger := i.Log(ctx)
	logger.Debug("UpdatePolicySet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w:\n%w", ErrUpdatePolicySet, ErrStructValidation, err)
	}

	req, err := request.NewPut(ctx, "/imaging/v2/policysets/%s", params.PolicySetID).
		AddHeader("Contract", params.ContractID).
		WithBody(params.UpdatePolicySet).
		Build()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %w", ErrUpdatePolicySet, err)
	}

	var result PolicySet

	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %w", ErrUpdatePolicySet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %w", ErrUpdatePolicySet, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) DeletePolicySet(ctx context.Context, params DeletePolicySetRequest) error {
	logger := i.Log(ctx)
	logger.Debug("DeletePolicySet")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%w: %w:\n%w", ErrDeletePolicySet, ErrStructValidation, err)
	}

	req, err := request.NewDelete(ctx, "/imaging/v2/policysets/%s", params.PolicySetID).
		AddHeader("Contract", params.ContractID).
		Build()
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %w", ErrDeletePolicySet, err)
	}

	resp, err := i.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %w", ErrDeletePolicySet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %w", ErrDeletePolicySet, i.Error(resp))
	}

	return nil
}
