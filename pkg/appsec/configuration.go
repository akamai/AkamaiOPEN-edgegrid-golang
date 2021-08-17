package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The Configuration interface supports creating, retrieving, updating and deleting security configurations.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#configuration
	Configuration interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurations
		GetConfigurations(ctx context.Context, params GetConfigurationsRequest) (*GetConfigurationsResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurations
		GetConfiguration(ctx context.Context, params GetConfigurationRequest) (*GetConfigurationResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfigurations
		CreateConfiguration(ctx context.Context, params CreateConfigurationRequest) (*CreateConfigurationResponse, error)

		UpdateConfiguration(ctx context.Context, params UpdateConfigurationRequest) (*UpdateConfigurationResponse, error)

		RemoveConfiguration(ctx context.Context, params RemoveConfigurationRequest) (*RemoveConfigurationResponse, error)
	}

	// GetConfigurationsRequest is used to list the available security configurations.
	GetConfigurationsRequest struct {
		ConfigID int    `json:"configId"`
		Name     string `json:"-"`
	}

	// GetConfigurationsResponse is returned from a call to GetConfigurations.
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

	// GetConfigurationRequest GetConfigurationRequest is used to retrieve information about a specific configuration.
	GetConfigurationRequest struct {
		ConfigID int `json:"configId"`
	}

	// GetConfigurationResponse is returned from a call to GetConfiguration.
	GetConfigurationResponse struct {
		Description         string   `json:"description,omitempty"`
		FileType            string   `json:"fileType,omitempty"`
		ID                  int      `json:"id,omitempty"`
		LatestVersion       int      `json:"latestVersion,omitempty"`
		Name                string   `json:"name,omitempty"`
		StagingVersion      int      `json:"stagingVersion,omitempty"`
		TargetProduct       string   `json:"targetProduct,omitempty"`
		ProductionHostnames []string `json:"productionHostnames,omitempty"`
		ProductionVersion   int      `json:"productionVersion,omitempty"`
	}

	// CreateConfigurationRequest is used to create a new WAP or KSD security configuration.
	CreateConfigurationRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ContractID  string   `json:"contractId"`
		GroupID     int      `json:"groupId"`
		Hostnames   []string `json:"hostnames"`
	}

	// CreateConfigurationResponse is returned from a call to CreateConfiguration.
	CreateConfigurationResponse struct {
		ConfigID    int    `json:"configId"`
		Version     int    `json:"version"`
		Description string `json:"description"`
		Name        string `json:"name"`
	}

	// UpdateConfigurationRequest is used tdo modify the name or description of an existing security configuration.
	UpdateConfigurationRequest struct {
		ConfigID    int    `json:"-"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// UpdateConfigurationResponse  is returned from a call to UpdateConfiguration.
	UpdateConfigurationResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// RemoveConfigurationRequest is used to remove an existing security configuration.
	RemoveConfigurationRequest struct {
		ConfigID int `json:"configId"`
	}

	// RemoveConfigurationResponse is returned from a call to RemoveConfiguration.
	RemoveConfigurationResponse struct {
		Empty int `json:"-"`
	}
)

// Validate validates a GetConfigurationRequest.
func (v GetConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates a GetConfigurationsRequest.
func (v GetConfigurationsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateConfigurationRequest.
func (v UpdateConfigurationRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
	}.Filter()
}

// Validate validates a RemoveConfigurationRequest.
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

	var getConfigurationResponse GetConfigurationResponse

	configid := params.ConfigID
	uri := fmt.Sprintf("/appsec/v1/configs/%d", configid)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfiguration request: %w", err)
	}

	resp, err := p.Exec(req, &getConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("GetConfiguration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &getConfigurationResponse, nil

}

func (p *appsec) GetConfigurations(ctx context.Context, _ GetConfigurationsRequest) (*GetConfigurationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetConfigurations")

	var rval GetConfigurationsResponse

	uri := "/appsec/v1/configs"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfigurations request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetConfigurations request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

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
		return nil, fmt.Errorf("failed to create UpdateConfiguration request: %w", err)
	}

	var rval UpdateConfigurationResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateConfiguration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *appsec) CreateConfiguration(ctx context.Context, params CreateConfigurationRequest) (*CreateConfigurationResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateConfiguration")

	uri :=
		"/appsec/v1/configs"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateConfiguration request: %w", err)
	}

	var rval CreateConfigurationResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("CreateConfiguration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) RemoveConfiguration(ctx context.Context, params RemoveConfigurationRequest) (*RemoveConfigurationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveConfigurationResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveConfiguration")

	uri, err := url.Parse(fmt.Sprintf("/appsec/v1/configs/%d", params.ConfigID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveConfiguration request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("RemoveConfiguration request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
