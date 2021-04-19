package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Configuration represents a collection of Configuration
//
// See: Configuration.GetConfiguration()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// Configuration  contains operations available on Configuration  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfiguration
	Configuration interface {
		GetConfigurations(ctx context.Context, params GetConfigurationsRequest) (*GetConfigurationsResponse, error)
		GetConfiguration(ctx context.Context, params GetConfigurationRequest) (*GetConfigurationResponse, error)
		CreateConfiguration(ctx context.Context, params CreateConfigurationRequest) (*CreateConfigurationResponse, error)
		UpdateConfiguration(ctx context.Context, params UpdateConfigurationRequest) (*UpdateConfigurationResponse, error)
		RemoveConfiguration(ctx context.Context, params RemoveConfigurationRequest) (*RemoveConfigurationResponse, error)
	}

	GetConfigurationsRequest struct {
		ConfigID int    `json:"configId"`
		Name     string `json:"-"`
	}

	GetConfigurationsResponse struct {
		Configurations []struct {
			Description         string   `json:"description,omitempty"`
			FileType            string   `json:"fileType,omitempty"`
			ID                  int      `json:"id,omitempty"`
			LatestVersion       int      `json:"latestVersion,omitempty"`
			Name                string   `json:"name,omitempty"`
			StagingVersion      int      `json:"stagingVersion,omitempty"`
			TargetProduct       string   `json:"targetProduct,omitempty"`
			ProductionHostnames []string `json:"productionHostnames,omitempty"`
			ProductionVersion   int      `json:"productionVersion,omitempty"`
		} `json:"configurations,omitempty"`
	}

	GetConfigurationRequest struct {
		ConfigID int `json:"configId"`
	}

	GetConfigurationResponse struct {
		Configurations []struct {
			Description         string   `json:"description,omitempty"`
			FileType            string   `json:"fileType,omitempty"`
			ID                  int      `json:"id,omitempty"`
			LatestVersion       int      `json:"latestVersion,omitempty"`
			Name                string   `json:"name,omitempty"`
			StagingVersion      int      `json:"stagingVersion,omitempty"`
			TargetProduct       string   `json:"targetProduct,omitempty"`
			ProductionHostnames []string `json:"productionHostnames,omitempty"`
			ProductionVersion   int      `json:"productionVersion,omitempty"`
		} `json:"configurations,omitempty"`
	}

	CreateConfigurationRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ContractID  string   `json:"contractId"`
		GroupID     int      `json:"groupId"`
		Hostnames   []string `json:"hostnames"`
	}

	CreateConfigurationResponse struct {
		ConfigID    int    `json:"configId"`
		Version     int    `json:"version"`
		Description string `json:"description"`
		Name        string `json:"name"`
	}

	UpdateConfigurationRequest struct {
		ConfigID    int    `json:"-"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UpdateConfigurationResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	RemoveConfigurationRequest struct {
		ConfigID int `json:"configId"`
	}

	RemoveConfigurationResponse struct {
		Empty int `json:"-"`
	}
)

// Validate validates GetConfigurationRequest
func (v GetConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates GetConfigurationsRequest
func (v GetConfigurationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates UpdateConfigurationRequest
func (v UpdateConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates RemoveConfigurationRequest
func (v RemoveConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

func (p *appsec) GetConfiguration(ctx context.Context, params GetConfigurationRequest) (*GetConfigurationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetConfiguration")

	var rval GetConfigurationResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getconfiguration request: %w", err)
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

func (p *appsec) GetConfigurations(ctx context.Context, params GetConfigurationsRequest) (*GetConfigurationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetConfigurations")

	var rval GetConfigurationsResponse

	uri := "/appsec/v1/configs"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getconfigurations request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getconfigurations request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a Configuration.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putconfiguration

func (p *appsec) UpdateConfiguration(ctx context.Context, params UpdateConfigurationRequest) (*UpdateConfigurationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateConfiguration")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create Configurationrequest: %w", err)
	}

	var rval UpdateConfigurationResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create Configuration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

// Create will create a new configuration.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfiguration
func (p *appsec) CreateConfiguration(ctx context.Context, params CreateConfigurationRequest) (*CreateConfigurationResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateConfiguration")

	uri :=
		"/appsec/v1/configs"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create configuration request: %w", err)
	}

	var rval CreateConfigurationResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create configurationrequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a Configuration
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteconfiguration

func (p *appsec) RemoveConfiguration(ctx context.Context, params RemoveConfigurationRequest) (*RemoveConfigurationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveConfigurationResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveConfiguration")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d",
		params.ConfigID,
	),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delconfiguration request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delconfiguration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
