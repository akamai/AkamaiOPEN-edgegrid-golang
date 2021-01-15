package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ApiRequestConstraints represents a collection of ApiRequestConstraints
//
// See: ApiRequestConstraints.GetApiRequestConstraints()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiRequestConstraints  contains operations available on ApiRequestConstraints  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapirequestconstraints
	ApiRequestConstraints interface {
		GetApiRequestConstraints(ctx context.Context, params GetApiRequestConstraintsRequest) (*GetApiRequestConstraintsResponse, error)
		UpdateApiRequestConstraints(ctx context.Context, params UpdateApiRequestConstraintsRequest) (*UpdateApiRequestConstraintsResponse, error)
	}

	GetApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
	}

	GetApiRequestConstraintsResponse struct {
		APIEndpoints []struct {
			ID     int    `json:"id"`
			Action string `json:"action"`
		} `json:"apiEndpoints"`
	}

	UpdateApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
		Action   string `json:"action"`
	}

	UpdateApiRequestConstraintsResponse struct {
		Action string `json:"action"`
	}

	RemoveApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
		Action   string `json:"action"`
	}

	RemoveApiRequestConstraintsResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates GetApiRequestConstraintsRequest
func (v GetApiRequestConstraintsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates UpdateApiRequestConstraintsRequest
func (v UpdateApiRequestConstraintsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiRequestConstraints(ctx context.Context, params GetApiRequestConstraintsRequest) (*GetApiRequestConstraintsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetApiRequestConstraints")

	var rval GetApiRequestConstraintsResponse
	var rvalfiltered GetApiRequestConstraintsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapirequestconstraints request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapirequestconstraints  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ApiID != 0 {
		for _, val := range rval.APIEndpoints {
			if val.ID == params.ApiID {
				rvalfiltered.APIEndpoints = append(rvalfiltered.APIEndpoints, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}

// Update will update a ApiRequestConstraints.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putapirequestconstraints

func (p *appsec) UpdateApiRequestConstraints(ctx context.Context, params UpdateApiRequestConstraintsRequest) (*UpdateApiRequestConstraintsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateApiRequestConstraints")

	var putURL string
	if params.ApiID != 0 {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints/%d",
			params.ConfigID,
			params.Version,
			params.PolicyID,
			params.ApiID,
		)
	} else {
		putURL = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create ApiRequestConstraintsrequest: %w", err)
	}

	var rval UpdateApiRequestConstraintsResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create ApiRequestConstraints request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
