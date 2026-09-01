package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ConfigurationVersion interface supports retrieving the versions of a configuration.
	ConfigurationVersion interface {
		// GetConfigurationVersions lists available versions for the specified security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-config-versions
		GetConfigurationVersions(ctx context.Context, params GetConfigurationVersionsRequest) (*GetConfigurationVersionsResponse, error)

		// GetConfigurationVersion returns information about a specific configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-version-number
		GetConfigurationVersion(ctx context.Context, params GetConfigurationVersionRequest) (*GetConfigurationVersionResponse, error)
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
			ConfigID   int               `json:"configId,omitempty"`
			Production EnvironmentStatus `json:"production,omitempty"`
			Staging    EnvironmentStatus `json:"staging,omitempty"`
			Version    int               `json:"version,omitempty"`
			BasedOn    int               `json:"basedOn,omitempty"`
		} `json:"versionList,omitempty"`
	}

	// GetConfigurationVersionRequest is used to retrieve information about a specific configuration version.
	GetConfigurationVersionRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

	// GetConfigurationVersionResponse is returned from a call to GetConfigurationVersion.
	GetConfigurationVersionResponse struct {
		BasedOn      int               `json:"basedOn,omitempty"`
		ConfigID     int               `json:"configId,omitempty"`
		ConfigName   string            `json:"configName,omitempty"`
		CreateDate   time.Time         `json:"createDate,omitempty"`
		CreatedBy    string            `json:"createdBy,omitempty"`
		Production   EnvironmentStatus `json:"production,omitempty"`
		Staging      EnvironmentStatus `json:"staging,omitempty"`
		Version      int               `json:"version,omitempty"`
		VersionNotes string            `json:"versionNotes,omitempty"`
	}

	// EnvironmentStatus represents the activation status for a configuration network (production or staging).
	EnvironmentStatus struct {
		Status string    `json:"status,omitempty"`
		Time   time.Time `json:"time,omitempty"`
	}
)

// Validate validates a GetConfigurationVersionRequest.
func (v GetConfigurationVersionRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

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

func (p *appsec) GetConfigurationVersion(ctx context.Context, params GetConfigurationVersionRequest) (*GetConfigurationVersionResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetConfigurationVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfigurationVersion request: %w", err)
	}

	var result GetConfigurationVersionResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get configuration version request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
