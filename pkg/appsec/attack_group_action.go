package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// AttackGroupAction represents a collection of AttackGroupAction
//
// See: AttackGroupAction.GetAttackGroupAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// AttackGroupAction  contains operations available on AttackGroupAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getattackgroupaction
	AttackGroupAction interface {
		GetAttackGroupActions(ctx context.Context, params GetAttackGroupActionsRequest) (*GetAttackGroupActionsResponse, error)
		GetAttackGroupAction(ctx context.Context, params GetAttackGroupActionRequest) (*GetAttackGroupActionResponse, error)
		UpdateAttackGroupAction(ctx context.Context, params UpdateAttackGroupActionRequest) (*UpdateAttackGroupActionResponse, error)
	}

	GetAttackGroupActionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAttackGroupActionsResponse struct {
		AttackGroupActions []struct {
			Action string `json:"action,omitempty"`
			Group  string `json:"group,omitempty"`
		} `json:"attackGroupActions,omitempty"`
	}

	GetAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetAttackGroupActionResponse struct {
		Action string `json:"action,omitempty"`
	}

	CreateAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	CreateAttackGroupActionResponse struct {
		Action string `json:"action"`
	}

	UpdateAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	UpdateAttackGroupActionResponse struct {
		Action string `json:"action"`
	}

	RemoveAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	RemoveAttackGroupActionResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetAttackGroupActionRequest
func (v GetAttackGroupActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
	}.Filter()
}

// Validate validates GetAttackGroupActionsRequest
func (v GetAttackGroupActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateAttackGroupActionRequest
func (v UpdateAttackGroupActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
	}.Filter()
}

func (p *appsec) GetAttackGroupAction(ctx context.Context, params GetAttackGroupActionRequest) (*GetAttackGroupActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupAction")

	var rval GetAttackGroupActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getattackgroupaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroupaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetAttackGroupActions(ctx context.Context, params GetAttackGroupActionsRequest) (*GetAttackGroupActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetAttackGroupActions")

	var rval GetAttackGroupActionsResponse
	var rvalfiltered GetAttackGroupActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getattackgroupactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getattackgroupactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Group != "" {
		for k, val := range rval.AttackGroupActions {
			if val.Group == params.Group {
				rvalfiltered.AttackGroupActions = append(rvalfiltered.AttackGroupActions, rval.AttackGroupActions[k])
			}
		}
	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}

// Update will update a AttackGroupAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupaction

func (p *appsec) UpdateAttackGroupAction(ctx context.Context, params UpdateAttackGroupActionRequest) (*UpdateAttackGroupActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateAttackGroupAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create AttackGroupActionrequest: %w", err)
	}

	var rval UpdateAttackGroupActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create AttackGroupAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
