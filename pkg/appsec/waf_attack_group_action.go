package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// WAFAttackGroupAction represents a collection of WAFAttackGroupAction
//
// See: WAFAttackGroupAction.GetWAFAttackGroupAction()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// WAFAttackGroupAction  contains operations available on WAFAttackGroupAction  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getwafattackgroupaction
	WAFAttackGroupAction interface {
		GetWAFAttackGroupActions(ctx context.Context, params GetWAFAttackGroupActionsRequest) (*GetWAFAttackGroupActionsResponse, error)
		GetWAFAttackGroupAction(ctx context.Context, params GetWAFAttackGroupActionRequest) (*GetWAFAttackGroupActionResponse, error)
		UpdateWAFAttackGroupAction(ctx context.Context, params UpdateWAFAttackGroupActionRequest) (*UpdateWAFAttackGroupActionResponse, error)
	}

	GetWAFAttackGroupActionsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetWAFAttackGroupActionsResponse struct {
		AttackGroupActions []struct {
			Action string `json:"action"`
			Group  string `json:"group"`
		} `json:"attackGroupActions"`
	}

	GetWAFAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Group    string `json:"group"`
	}

	GetWAFAttackGroupActionResponse struct {
		Action string `json:"action"`
	}

	CreateWAFAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	CreateWAFAttackGroupActionResponse struct {
		Action string `json:"action"`
	}

	UpdateWAFAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	UpdateWAFAttackGroupActionResponse struct {
		Action string `json:"action"`
	}

	RemoveWAFAttackGroupActionRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		Action   string `json:"action"`
		Group    string `json:"group"`
	}

	RemoveWAFAttackGroupActionResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetWAFAttackGroupActionRequest
func (v GetWAFAttackGroupActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
	}.Filter()
}

// Validate validates GetWAFAttackGroupActionsRequest
func (v GetWAFAttackGroupActionsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateWAFAttackGroupActionRequest
func (v UpdateWAFAttackGroupActionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"Group":    validation.Validate(v.Group, validation.Required),
	}.Filter()
}

func (p *appsec) GetWAFAttackGroupAction(ctx context.Context, params GetWAFAttackGroupActionRequest) (*GetWAFAttackGroupActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFAttackGroupAction")

	var rval GetWAFAttackGroupActionResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafattackgroupaction request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafattackgroupaction  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetWAFAttackGroupActions(ctx context.Context, params GetWAFAttackGroupActionsRequest) (*GetWAFAttackGroupActionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetWAFAttackGroupActions")

	var rval GetWAFAttackGroupActionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getwafattackgroupactions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getwafattackgroupactions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a WAFAttackGroupAction.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putwafattackgroupaction

func (p *appsec) UpdateWAFAttackGroupAction(ctx context.Context, params UpdateWAFAttackGroupActionRequest) (*UpdateWAFAttackGroupActionResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateWAFAttackGroupAction")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/attack-groups/%s",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.Group,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create WAFAttackGroupActionrequest: %w", err)
	}

	var rval UpdateWAFAttackGroupActionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create WAFAttackGroupAction request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
