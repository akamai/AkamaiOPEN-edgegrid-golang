package appsec

import (
	"context"
	"fmt"
	"net/http"
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
	}

	GetConfigurationsRequest struct {
		ConfigID int `json:"configId"`
	}

	GetConfigurationsResponse struct {
		Configurations []struct {
			Description         string   `json:"description,omitempty"`
			FileType            string   `json:"fileType"`
			ID                  int      `json:"id"`
			LatestVersion       int      `json:"latestVersion"`
			Name                string   `json:"name,omitempty"`
			StagingVersion      int      `json:"stagingVersion,omitempty"`
			TargetProduct       string   `json:"targetProduct"`
			ProductionHostnames []string `json:"productionHostnames,omitempty"`
			ProductionVersion   int      `json:"productionVersion,omitempty"`
		} `json:"configurations"`
	}
)

func (p *appsec) GetConfigurations(ctx context.Context, params GetConfigurationsRequest) (*GetConfigurationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetConfigurations")

	var rval GetConfigurationsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs",
	)

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
