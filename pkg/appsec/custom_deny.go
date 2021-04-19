package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CustomDeny represents a collection of CustomDeny
//
// See: CustomDeny.GetCustomDeny()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// CustomDeny  contains operations available on CustomDeny  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomdeny
	CustomDeny interface {
		GetCustomDenyList(ctx context.Context, params GetCustomDenyListRequest) (*GetCustomDenyListResponse, error)
		GetCustomDeny(ctx context.Context, params GetCustomDenyRequest) (*GetCustomDenyResponse, error)
		CreateCustomDeny(ctx context.Context, params CreateCustomDenyRequest) (*CreateCustomDenyResponse, error)
		UpdateCustomDeny(ctx context.Context, params UpdateCustomDenyRequest) (*UpdateCustomDenyResponse, error)
		RemoveCustomDeny(ctx context.Context, params RemoveCustomDenyRequest) (*RemoveCustomDenyResponse, error)
	}

	customDenyID string

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

	GetCustomDenyListRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		ID       string `json:"id,omitempty"`
	}

	GetCustomDenyRequest struct {
		ConfigID int    `json:"configId"`
		Version  int    `json:"version"`
		ID       string `json:"id,omitempty"`
	}

	CreateCustomDenyRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

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

	UpdateCustomDenyRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		ID             string          `json:"id"`
		JsonPayloadRaw json.RawMessage `json:"-"`
	}

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

	RemoveCustomDenyRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		ID       string `json:"id,omitempty"`
	}

	RemoveCustomDenyResponse struct {
		Empty string `json:"-"`
	}
)

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

// Validate validates GetCustomDenyRequest
func (v GetCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates GetCustomDenysRequest
func (v GetCustomDenyListRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates CreateCustomDenyRequest
func (v CreateCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateCustomDenyRequest
func (v UpdateCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

// Validate validates RemoveCustomDenyRequest
func (v RemoveCustomDenyRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"ID":       validation.Validate(v.ID, validation.Required),
	}.Filter()
}

func (p *appsec) GetCustomDeny(ctx context.Context, params GetCustomDenyRequest) (*GetCustomDenyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomDeny")

	var rval GetCustomDenyResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny/%s",
		params.ConfigID,
		params.Version,
		params.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcustomdeny request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) GetCustomDenyList(ctx context.Context, params GetCustomDenyListRequest) (*GetCustomDenyListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetCustomDenyList")

	var rval GetCustomDenyListResponse
	var rvalfiltered GetCustomDenyListResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getcustomdenylist request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getcustomdenylist request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.ID != "" {
		for _, val := range rval.CustomDenyList {
			if string(val.ID) == params.ID {
				rvalfiltered.CustomDenyList = append(rvalfiltered.CustomDenyList, val)
			}
		}

	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}

// Update will update a CustomDeny.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putcustomdeny

func (p *appsec) UpdateCustomDeny(ctx context.Context, params UpdateCustomDenyRequest) (*UpdateCustomDenyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateCustomDeny")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny/%s",
		params.ConfigID,
		params.Version,
		params.ID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create CustomDenyrequest: %w", err)
	}

	var rval UpdateCustomDenyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create CustomDeny request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Create will create a new customdeny.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postcustomdeny
func (p *appsec) CreateCustomDeny(ctx context.Context, params CreateCustomDenyRequest) (*CreateCustomDenyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateCustomDeny")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create customdeny request: %w", err)
	}

	var rval CreateCustomDenyResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &rval, params.JsonPayloadRaw)
	if err != nil {
		return nil, fmt.Errorf("create customdenyrequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a CustomDeny
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletecustomdeny

func (p *appsec) RemoveCustomDeny(ctx context.Context, params RemoveCustomDenyRequest) (*RemoveCustomDenyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveCustomDenyResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveCustomDeny")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/custom-deny/%s",
		params.ConfigID,
		params.Version,
		params.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delcustomdeny request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delcustomdeny request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
