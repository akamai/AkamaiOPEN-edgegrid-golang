package appsec

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ConfigurationVersionClone interface supports creating, retrieving, and removing
	// clones of a configuration version.
	ConfigurationVersionClone interface {
		// GetConfigurationVersionClone returns basic detail about a configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-version-number
		GetConfigurationVersionClone(ctx context.Context, params GetConfigurationVersionCloneRequest) (*GetConfigurationVersionCloneResponse, error)

		// CreateConfigurationVersionClone  creates a new version of the specified security configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/post-config-versions
		CreateConfigurationVersionClone(ctx context.Context, params CreateConfigurationVersionCloneRequest) (*CreateConfigurationVersionCloneResponse, error)

		// RemoveConfigurationVersionClone deletes the specified configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/delete-version-number
		RemoveConfigurationVersionClone(ctx context.Context, params RemoveConfigurationVersionCloneRequest) (*RemoveConfigurationVersionCloneResponse, error)
	}

	// GetConfigurationVersionCloneRequest is used to retrieve information about an existing configuration version.
	GetConfigurationVersionCloneRequest struct {
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

	// GetConfigurationVersionCloneResponse is returned from a call to GetConfigurationVersionClone.
	GetConfigurationVersionCloneResponse struct {
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

	// CreateConfigurationVersionCloneRequest is used to clone an existing configuration version.
	CreateConfigurationVersionCloneRequest struct {
		ConfigID          int  `json:"-"`
		CreateFromVersion int  `json:"createFromVersion"`
		RuleUpdate        bool `json:"ruleUpdate"`
	}

	// CreateConfigurationVersionCloneResponse is returned from a call to CreateConfigurationVersionClone.
	CreateConfigurationVersionCloneResponse struct {
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

	// RemoveConfigurationVersionCloneRequest is used to remove an existing configuration version.
	RemoveConfigurationVersionCloneRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// RemoveConfigurationVersionCloneResponse is returned from a call to RemoveConfigurationVersionClone.
	RemoveConfigurationVersionCloneResponse struct {
		Empty string `json:"-"`
	}
)

// Validate validates a GetConfigurationCloneRequest.
func (v GetConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates a CreateConfigurationCloneRequest.
func (v CreateConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.CreateFromVersion, validation.Required),
	}.Filter()
}

// Validate validates a RemoveConfigurationCloneRequest.
func (v RemoveConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetConfigurationVersionClone(ctx context.Context, params GetConfigurationVersionCloneRequest) (*GetConfigurationVersionCloneResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetConfigurationVersionClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetConfigurationVersionClone request: %w", err)
	}

	var result GetConfigurationVersionCloneResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get configuration version clone request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) CreateConfigurationVersionClone(ctx context.Context, params CreateConfigurationVersionCloneRequest) (*CreateConfigurationVersionCloneResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("CreateConfigurationVersionClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions", params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateConfigurationVersionClone request: %w", err)
	}

	var result CreateConfigurationVersionCloneResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("create configuration version clone request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveConfigurationVersionClone(ctx context.Context, params RemoveConfigurationVersionCloneRequest) (*RemoveConfigurationVersionCloneResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveConfigurationVersionClone")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf("/appsec/v1/configs/%d/versions/%d", params.ConfigID, params.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoveConfigurationVersionClone request: %w", err)
	}

	var result RemoveConfigurationVersionCloneResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("remove configuration version clone request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
