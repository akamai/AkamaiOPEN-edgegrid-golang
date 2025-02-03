package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The ConfigurationVersion interface supports retrieving the versions of a configuration.
	ConfigurationVersion interface {
		// GetConfigurationVersions lists available versions for the specified security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-config-versions
		GetConfigurationVersions(ctx context.Context, params GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error)
	}

	// GetConfigurationVersionsRequest is used to retrieve the versions of a security configuration.
	GetConfigurationVersionsRequest struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
	}

	// GetConfigurationVersionsResponse is returned from a call to GetConfigurationVersions.
	GetConfigurationVersionsResponse struct {
		ConfigID           int    `json:"configId,omitempty"`
		ConfigName         string `json:"configName,omitempty"`
		LastCreatedVersion int    `json:"lastCreatedVersion,omitempty"`
		Page               int    `json:"page,omitempty"`
		PageSize           int    `json:"pageSize,omitempty"`
		TotalSize          int    `json:"totalSize,omitempty"`
		VersionList        []struct {
			ConfigID   int `json:"configId,omitempty"`
			Production struct {
				Status string `json:"status,omitempty"`
			} `json:"production,omitempty"`
			Staging struct {
				Status string `json:"status,omitempty"`
			} `json:"staging,omitempty"`
			Version int `json:"version,omitempty"`
			BasedOn int `json:"basedOn,omitempty"`
		} `json:"versionList,omitempty"`
	}
)

func (p *appsec) GetConfigurationVersions(ctx context.Context, params GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetConfigurationVersions")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions?page=-1&detail=false",
		params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfigurationVersions request: %w", err)
	}

	var result GetConfigurationVersionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get configuration versions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
