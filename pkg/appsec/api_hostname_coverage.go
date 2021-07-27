package appsec

import (
	"context"
	"fmt"
	"net/http"
)

// ApiHostnameCoverage represents a collection of ApiHostnameCoverage
//
// See: ApiHostnameCoverage.GetApiHostnameCoverage()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiHostnameCoverage  contains operations available on ApiHostnameCoverage  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapihostnamecoverage
	ApiHostnameCoverage interface {
		GetApiHostnameCoverage(ctx context.Context, params GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error)
	}

	GetApiHostnameCoverageRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	GetApiHostnameCoverageResponse struct {
		HostnameCoverage []struct {
			Configuration  *ConfigurationHostnameCoverage `json:"configuration,omitempty"`
			Status         string                         `json:"status"`
			HasMatchTarget bool                           `json:"hasMatchTarget"`
			Hostname       string                         `json:"hostname"`
			PolicyNames    []string                       `json:"policyNames"`
		} `json:"hostnameCoverage"`
	}

	ConfigurationHostnameCoverage struct {
		ID      int    `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Version int    `json:"version,omitempty"`
	}
)

func (p *appsec) GetApiHostnameCoverage(ctx context.Context, _ GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverage")

	var rval GetApiHostnameCoverageResponse

	uri :=
		"/appsec/v1/hostname-coverage"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapihostnamecoverage request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapihostnamecoverage  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
