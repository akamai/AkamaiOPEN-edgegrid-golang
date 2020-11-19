package appsec

import (
	"context"
	"fmt"
	"net/http"
)

// ConfigurationVersion represents a collection of ConfigurationVersion
//
// See: ConfigurationVersion.GetConfigurationVersion()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ConfigurationVersion  contains operations available on ConfigurationVersion  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationversion
	ConfigurationVersion interface {
		GetConfigurationVersions(ctx context.Context, params GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error)
	}

	GetConfigurationVersionsRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
	}

	GetConfigurationVersionsResponse struct {
		ConfigID           int    `json:"configId"`
		ConfigName         string `json:"configName"`
		LastCreatedVersion int    `json:"lastCreatedVersion"`
		Page               int    `json:"page"`
		PageSize           int    `json:"pageSize"`
		TotalSize          int    `json:"totalSize"`
		VersionList        []struct {
			ConfigID   int `json:"configId"`
			Production struct {
				Status string `json:"status"`
			} `json:"production"`
			Staging struct {
				Status string `json:"status"`
			} `json:"staging"`
			Version int `json:"version"`
			BasedOn int `json:"basedOn,omitempty"`
		} `json:"versionList"`
	}
)

func (p *appsec) GetConfigurationVersions(ctx context.Context, params GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetConfigurationVersions")

	var rval GetConfigurationVersionsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions?page=-1&detail=false",
		params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getconfigurationversions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getconfigurationversions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
