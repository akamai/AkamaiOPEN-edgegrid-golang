package appsec

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ConfigurationClone represents a collection of ConfigurationClone
//
// See: ConfigurationClone.GetConfigurationClone()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ConfigurationClone  contains operations available on ConfigurationClone  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationclone
	ConfigurationVersionClone interface {
		GetConfigurationVersionClone(ctx context.Context, params GetConfigurationVersionCloneRequest) (*GetConfigurationVersionCloneResponse, error)
		CreateConfigurationVersionClone(ctx context.Context, params CreateConfigurationVersionCloneRequest) (*CreateConfigurationVersionCloneResponse, error)
		RemoveConfigurationVersionClone(ctx context.Context, params RemoveConfigurationVersionCloneRequest) (*RemoveConfigurationVersionCloneResponse, error)
	}

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

	CreateConfigurationVersionCloneRequest struct {
		ConfigID          int  `json:"-"`
		CreateFromVersion int  `json:"createFromVersion"`
		RuleUpdate        bool `json:"ruleUpdate"`
	}

	RemoveConfigurationVersionCloneRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	RemoveConfigurationVersionCloneResponse struct {
		Empty string `json:"-"`
	}
)

// Validate validates GetConfigurationCloneRequest
func (v GetConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates GetConfigurationCloneRequest
func (v CreateConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.CreateFromVersion, validation.Required),
	}.Filter()
}

// Validate validates GetConfigurationCloneRequest
func (v RemoveConfigurationVersionCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetConfigurationVersionClone(ctx context.Context, params GetConfigurationVersionCloneRequest) (*GetConfigurationVersionCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetConfigurationVersionClone")

	var rval GetConfigurationVersionCloneResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getconfigurationversionclone request: %w", err)
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

/// Create will create a new configurationVersionclone.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfigurationclone
func (p *appsec) CreateConfigurationVersionClone(ctx context.Context, params CreateConfigurationVersionCloneRequest) (*CreateConfigurationVersionCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateConfigurationVersionClone")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions",
		params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create configurationVersionclone request: %w", err)
	}

	var rval CreateConfigurationVersionCloneResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create configurationVersionclonerequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Delete will delete a ConfigurationVersion
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteconfiguration

func (p *appsec) RemoveConfigurationVersionClone(ctx context.Context, params RemoveConfigurationVersionCloneRequest) (*RemoveConfigurationVersionCloneResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var rval RemoveConfigurationVersionCloneResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveConfiguration")

	uri, err := url.Parse(fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d",
		params.ConfigID,
		params.Version,
	),
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delconfigurationversion request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("delconfigurationversion request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
