package appsec

import (
	"context"
	"fmt"
	"net/http"

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
	ConfigurationClone interface {
		//	GetConfigurationClones(ctx context.Context, params GetConfigurationClonesRequest) (*GetConfigurationClonesResponse, error)
		GetConfigurationClone(ctx context.Context, params GetConfigurationCloneRequest) (*GetConfigurationCloneResponse, error)
		CreateConfigurationClone(ctx context.Context, params CreateConfigurationCloneRequest) (*CreateConfigurationCloneResponse, error)
	}

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

	/*	CreateConfigurationCloneResponse struct {
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
	*/

	CreateConfigurationCloneResponse struct {
		BasedOn    int       `json:"basedOn"`
		ConfigID   int       `json:"configId"`
		CreateDate time.Time `json:"createDate"`
		CreatedBy  string    `json:"createdBy"`
		Production struct {
			Status string `json:"status"`
		} `json:"production"`
		Staging struct {
			Status string `json:"status"`
		} `json:"staging"`
		Version int `json:"version"`
	}

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

	CreateConfigurationCloneRequest struct {
		ConfigID          int  `json:"-"`
		CreateFromVersion int  `json:"createFromVersion"`
		RuleUpdate        bool `json:"ruleUpdate"`
	}
)

// Validate validates GetConfigurationCloneRequest
func (v GetConfigurationCloneRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

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
		return nil, fmt.Errorf("failed to create getconfigurationclone request: %w", err)
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

/// Create will create a new configurationclone.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfigurationclone
func (p *appsec) CreateConfigurationClone(ctx context.Context, params CreateConfigurationCloneRequest) (*CreateConfigurationCloneResponse, error) {
	/*if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}*/

	logger := p.Log(ctx)
	logger.Debug("CreateConfigurationClone")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions",
		params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create configurationclone request: %w", err)
	}

	var rval CreateConfigurationCloneResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create configurationclonerequest failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
