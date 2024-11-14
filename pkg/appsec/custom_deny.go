package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The CustomDeny interface supports creating, retrievinfg, modifying and removing custom deny actions
	// for a configuration.
	CustomDeny interface {
		// GetCustomDenyList returns custom deny actions for a specific security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-custom-deny-actions
		GetCustomDenyList(ctx context.Context, params GetCustomDenyListRequest) (*GetCustomDenyListResponse, error)

		// GetCustomDeny returns the specified custom deny action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-custom-deny-action
		GetCustomDeny(ctx context.Context, params GetCustomDenyRequest) (*GetCustomDenyResponse, error)

		// CreateCustomDeny creates a new custom deny action for a specific configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-custom-deny
		CreateCustomDeny(ctx context.Context, params CreateCustomDenyRequest) (*CreateCustomDenyResponse, error)

		// UpdateCustomDeny updates details for a specific custom deny action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-custom-deny
		UpdateCustomDeny(ctx context.Context, params UpdateCustomDenyRequest) (*UpdateCustomDenyResponse, error)

		// RemoveCustomDeny deletes a custom deny action.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-custom-deny
		RemoveCustomDeny(ctx context.Context, params RemoveCustomDenyRequest) (*RemoveCustomDenyResponse, error)
	}

	customDenyID string

	// GetCustomDenyListRequest is used to retrieve the custom deny actions for a configuration.
	GetCustomDenyListRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		ID       string `json:"id,omitempty"`
	}

	// GetCustomDenyListResponse is returned from a call to GetCustomDenyList.
	GetCustomDenyListResponse struct {
		CustomDenyList []struct {
			Description string       `json:"description,omitempty"`
			Name        string       `json:"name"`
			ID          customDenyID `json:"id"`
			Parameters  []struct {
				DisplayName string `json:"-"`
				Name        string `json:"name"`
				Value       string `json:"value"`
			} `json:"parameters"`
		} `json:"customDenyList"`
	}

	// GetCustomDenyRequest is used to retrieve a specific custom deny action.
	GetCustomDenyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		ID       string `json:"id,omitempty"`
	}

	// GetCustomDenyResponse is returned from a call to GetCustomDeny.
	GetCustomDenyResponse struct {
		Description string       `json:"description,omitempty"`
		Name        string       `json:"name"`
		ID          customDenyID `json:"-"`
		Parameters  []struct {
			DisplayName string `json:"-"`
			Name        string `json:"name"`
			Value       string `json:"value"`
		} `json:"parameters"`
	}

	// CreateCustomDenyRequest is used to create a new custom deny action for a specific configuration.
	CreateCustomDenyRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// CreateCustomDenyResponse is returned from a call to CreateCustomDeny.
	CreateCustomDenyResponse struct {
		Description string       `json:"description,omitempty"`
		Name        string       `json:"name"`
		ID          customDenyID `json:"id"`
		Parameters  []struct {
			DisplayName string `json:"-"`
			Name        string `json:"name"`
			Value       string `json:"value"`
		} `json:"parameters"`
	}

	// UpdateCustomDenyRequest is used to details for a specific custom deny action.
	UpdateCustomDenyRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		ID             string          `json:"id"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

	// UpdateCustomDenyResponse is returned from a call to UpdateCustomDeny.
	UpdateCustomDenyResponse struct {
		Description string       `json:"description,omitempty"`
		Name        string       `json:"name"`
		ID          customDenyID `json:"-"`
		Parameters  []struct {
			DisplayName string `json:"-"`
			Name        string `json:"name"`
			Value       string `json:"value"`
		} `json:"parameters"`
	}

	// RemoveCustomDenyRequest is used to remove an existing custom deny action.
	RemoveCustomDenyRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		ID       string `json:"id,omitempty"`
	}

	// RemoveCustomDenyResponse is returned from a call to RemoveCustomDeny.
	RemoveCustomDenyResponse struct {
		Empty string `json:"-"`
	}
)

// UnmarshalJSON reads a customDenyID struct from its data argument.
func (c *customDenyID) UnmarshalJSON(data []byte) error {
	var nums interface{}
	err := json.Unmarshal(data, &nums)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(nums)
	switch items.Kind() {
	case reflect.String:
		*c = customDenyID(nums.(string))
	case reflect.Int:

		*c = customDenyID(strconv.Itoa(nums.(int)))

	}
	return nil
}

// Validate validates a GetCustomDenyRequest.
func (v GetCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates a GetCustomDenysRequest.
func (v GetCustomDenyListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateCustomDenyRequest.
func (v CreateCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateCustomDenyRequest.
func (v UpdateCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveCustomDenyRequest.
func (v RemoveCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

func (p *appsec) GetCustomDeny(ctx context.Context, params GetCustomDenyRequest) (*GetCustomDenyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCustomDeny")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny/%s",
		params.ConfigID,
		params.Version,
		params.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomDeny request: %w", err)
	}

	var result GetCustomDenyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get custom deny request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) GetCustomDenyList(ctx context.Context, params GetCustomDenyListRequest) (*GetCustomDenyListResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetCustomDenyList")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCustomDenyList request: %w", err)
	}

	var result GetCustomDenyListResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get custom deny list request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ID != "" {
		var filteredResult GetCustomDenyListResponse
		for _, val := range result.CustomDenyList {
			if string(val.ID) == params.ID {
				filteredResult.CustomDenyList = append(filteredResult.CustomDenyList, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}

func (p *appsec) UpdateCustomDeny(ctx context.Context, params UpdateCustomDenyRequest) (*UpdateCustomDenyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateCustomDeny")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny/%s",
		params.ConfigID,
		params.Version,
		params.ID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateCustomDeny request: %w", err)
	}

	var result UpdateCustomDenyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("update custom deny request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateCustomDeny(ctx context.Context, params CreateCustomDenyRequest) (*CreateCustomDenyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateCustomDeny")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateCustomDeny request: %w", err)
	}

	var result CreateCustomDenyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create custom deny request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveCustomDeny(ctx context.Context, params RemoveCustomDenyRequest) (*RemoveCustomDenyResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveCustomDeny")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/custom-deny/%s", params.ConfigID, params.Version, params.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveCustomDeny request: %w", err)
	}

	var result RemoveCustomDenyResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("remove custom deny request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
