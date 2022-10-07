package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiRequestConstraints interface supports retrieving, modifying, or removing the action
	// taken when any API request constraint is triggered, or when a specific API request constraint
	// is triggered.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#apirequestconstraintsgroup
	ApiRequestConstraints interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapirequestconstraints
		GetApiRequestConstraints(ctx context.Context, params GetApiRequestConstraintsRequest) (*GetApiRequestConstraintsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putapirequestconstraints
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putactionsperapi
		UpdateApiRequestConstraints(ctx context.Context, params UpdateApiRequestConstraintsRequest) (*UpdateApiRequestConstraintsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putapirequestconstraints
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putactionsperapi
		RemoveApiRequestConstraints(ctx context.Context, params RemoveApiRequestConstraintsRequest) (*RemoveApiRequestConstraintsResponse, error)
	}

	// GetApiRequestConstraintsRequest is used to retrieve the list of APIs with their constraints and associated actions.
	GetApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
	}

	// GetApiRequestConstraintsResponse is returned from a call to GetApiRequestConstraints.
	GetApiRequestConstraintsResponse struct {
		APIEndpoints []ApiEndpoint `json:"apiEndpoints,omitempty"`
	}

	// ApiEndpoint describes an API endpoint and its associated action.
	ApiEndpoint struct {
		ID     int    `json:"id"`
		Action string `json:"action"`
	}

	// UpdateApiRequestConstraintsRequest is used to modify the action taken when an API request contraint is triggered.
	UpdateApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
		Action   string `json:"action"`
	}

	// UpdateApiRequestConstraintsResponse is returned from a call to UpdateApiRequestConstraints.
	UpdateApiRequestConstraintsResponse struct {
		Action string `json:"action"`
	}

	// RemoveApiRequestConstraintsRequest is used to remove an API request constraint's action.
	RemoveApiRequestConstraintsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		ApiID    int    `json:"-"`
		Action   string `json:"action"`
	}

	// RemoveApiRequestConstraintsResponse is returned from a call to RemoveApiRequestConstraints.
	RemoveApiRequestConstraintsResponse struct {
		Action string `json:"action"`
	}
)

// Validate validates a GetApiRequestConstraintsRequest.
func (v GetApiRequestConstraintsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateApiRequestConstraintsRequest.
func (v UpdateApiRequestConstraintsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveApiRequestConstraintsRequest.
func (v RemoveApiRequestConstraintsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiRequestConstraints(ctx context.Context, params GetApiRequestConstraintsRequest) (*GetApiRequestConstraintsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiRequestConstraints")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetApiRequestConstraints request: %w", err)
	}

	var result GetApiRequestConstraintsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get API request constraints request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ApiID != 0 {
		var filteredResult GetApiRequestConstraintsResponse
		for _, val := range result.APIEndpoints {
			if val.ID == params.ApiID {
				filteredResult.APIEndpoints = append(filteredResult.APIEndpoints, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateApiRequestConstraints(ctx context.Context, params UpdateApiRequestConstraintsRequest) (*UpdateApiRequestConstraintsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateApiRequestConstraints")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri string
	if params.ApiID != 0 {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints/%d",
			params.ConfigID,
			params.Version,
			params.PolicyID,
			params.ApiID,
		)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateApiRequestConstraints request: %w", err)
	}

	var result UpdateApiRequestConstraintsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update API request constraints request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveApiRequestConstraints(ctx context.Context, params RemoveApiRequestConstraintsRequest) (*RemoveApiRequestConstraintsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveApiRequestConstraints")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri string
	if params.ApiID != 0 {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints/%d",
			params.ConfigID,
			params.Version,
			params.PolicyID,
			params.ApiID,
		)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-request-constraints",
			params.ConfigID,
			params.Version,
			params.PolicyID,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveApiRequestConstraints request: %w", err)
	}

	var result RemoveApiRequestConstraintsResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("remove API request constraints request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
