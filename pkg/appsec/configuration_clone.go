package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ConfigurationClone interface supports cloning an existing configuration and retrieving a configuration version.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#configurationclone
	ConfigurationClone interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationversion
		// Note: this method is DEPRECATED and will be removed in a future release.
		GetConfigurationClone(ctx context.Context, params GetConfigurationCloneRequest) (*GetConfigurationCloneResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsummarylistofconfigurationversions
		CreateConfigurationClone(ctx context.Context, params CreateConfigurationCloneRequest) (*CreateConfigurationCloneResponse, error)
	}

	// GetConfigurationCloneRequest is used to retrieve information about an existing security configuration.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetConfigurationCloneRequest struct {
		ConfigID     int       `json:"configId"`
		ConfigName   string    `json:"configName"`
		Version      int       `json:"version"`
		VersionNotes string    `json:"versionNotes"`
		CreateDate   time.Time `json:"createDate"`
		CreatedBy    string    `json:"createdBy"`
		BasedOn      int       `json:"basedOn"`
		Production   struct {
			Status string    `json:"status"`
			Time   time.Time `json:"time"`
		} `json:"production"`
		Staging struct {
			Status string `json:"status"`
		} `json:"staging"`
	}

	// GetConfigurationCloneResponse is returned from a call to GetConfigurationClone.
	// Note: this struct is DEPRECATED and will be removed in a future release.
	GetConfigurationCloneResponse struct {
		ConfigID     int       `json:"configId"`
		ConfigName   string    `json:"configName"`
		Version      int       `json:"version"`
		VersionNotes string    `json:"versionNotes"`
		CreateDate   time.Time `json:"createDate"`
		CreatedBy    string    `json:"createdBy"`
		BasedOn      int       `json:"basedOn"`
		Production   struct {
			Status string    `json:"status"`
			Time   time.Time `json:"time"`
		} `json:"production"`
		Staging struct {
			Status string `json:"status"`
		} `json:"staging"`
	}

	// CreateConfigurationCloneRequest is used to clone an existing security configuration.
	CreateConfigurationCloneRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ContractID  string   `json:"contractId"`
		GroupID     int      `json:"groupId"`
		Hostnames   []string `json:"hostnames"`
		CreateFrom  struct {
			ConfigID int `json:"configId"`
			Version  int `json:"version"`
		} `json:"createFrom"`
	}

	// CreateConfigurationCloneResponse is returned from a call to CreateConfigurationClone.
	CreateConfigurationCloneResponse struct {
		ConfigID    int    `json:"configId"`
		Version     int    `json:"version"`
		Description string `json:"description"`
		Name        string `json:"name"`
	}
)

// Validate validates a GetConfigurationCloneRequest.
// Note: this method is DEPRECATED and will be removed in a future release.
func (v GetConfigurationCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateConfigurationCloneRequest.
func (v CreateConfigurationCloneRequest) Validate() error {
	return validation.Errors{
		"CreateFromConfigID": validation.Validate(v.CreateFrom.ConfigID, validation.Required),
	}.Filter()
}

// Note: this method is DEPRECATED and will be removed in a future release.
func (p *appsec) GetConfigurationClone(ctx context.Context, params GetConfigurationCloneRequest) (*GetConfigurationCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetConfigurationClone")

	var rval GetConfigurationCloneResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfigurationClone request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetConfigurationClone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *appsec) CreateConfigurationClone(ctx context.Context, params CreateConfigurationCloneRequest) (*CreateConfigurationCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateConfigurationClone")

	uri := "/appsec/v1/configs/"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateConfigurationClone request: %w", err)
	}

	var rval CreateConfigurationCloneResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("CreateConfigurationClone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
