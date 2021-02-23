package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// VersionNotes represents a collection of VersionNotes
//
// See: VersionNotes.GetVersionNotes()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// VersionNotes  contains operations available on VersionNotes  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getversionnotes
	VersionNotes interface {
		GetVersionNotes(ctx context.Context, params GetVersionNotesRequest) (*GetVersionNotesResponse, error)
		UpdateVersionNotes(ctx context.Context, params UpdateVersionNotesRequest) (*UpdateVersionNotesResponse, error)
	}

	GetVersionNotesRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	GetVersionNotesResponse struct {
		Notes string `json:"notes"`
	}

	UpdateVersionNotesRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`

		Notes string `json:"notes"`
	}

	UpdateVersionNotesResponse struct {
		Notes string `json:"notes"`
	}
)

// Validate validates GetVersionNotesRequest
func (v GetVersionNotesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateVersionNotesRequest
func (v UpdateVersionNotesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetVersionNotes(ctx context.Context, params GetVersionNotesRequest) (*GetVersionNotesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetVersionNotes")

	var rval GetVersionNotesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/version-notes",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getversionnotes request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getversionnotes  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

// Update will update a VersionNotes.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putversionnotes

func (p *appsec) UpdateVersionNotes(ctx context.Context, params UpdateVersionNotesRequest) (*UpdateVersionNotesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateVersionNotes")

	putURL := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/version-notes",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create VersionNotesrequest: %w", err)
	}

	var rval UpdateVersionNotesResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create VersionNotes request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
