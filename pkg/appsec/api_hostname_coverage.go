package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The ApiHostnameCoverage interface supports retrieving hostnames with their current protections,
	// activation statuses, and other summary information.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#hostnamecoverage
	ApiHostnameCoverage interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethostnamecoverage
		GetApiHostnameCoverage(ctx context.Context, params GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error)
	}

	// GetApiHostnameCoverageRequest is used to call GetApiHostnameCoverage.
	GetApiHostnameCoverageRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	// GetApiHostnameCoverageResponse is returned from a call to GetApiHostnameCoverage.
	GetApiHostnameCoverageResponse struct {
		HostnameCoverage []struct {
			Configuration  *ConfigurationHostnameCoverage `json:"configuration,omitempty"`
			Status         string                         `json:"status"`
			HasMatchTarget bool                           `json:"hasMatchTarget"`
			Hostname       string                         `json:"hostname"`
			PolicyNames    []string                       `json:"policyNames"`
		} `json:"hostnameCoverage"`
	}

	// ConfigurationHostnameCoverage describes a specific configuration version.
	ConfigurationHostnameCoverage struct {
		ID      int    `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Version int    `json:"version,omitempty"`
	}
)

func (p *appsec) GetApiHostnameCoverage(ctx context.Context, _ GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverage")

	var result GetApiHostnameCoverageResponse

	uri := "/appsec/v1/hostname-coverage"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetApiHostnameCoverage request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetApiHostnameCoverage request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}
