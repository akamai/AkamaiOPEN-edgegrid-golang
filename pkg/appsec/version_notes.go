package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The VersionNotes interface supports retrieving and modifying the version notes for a configuration and version.
	VersionNotes interface {
		// GetVersionNotes gets the most recent version notes for a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-version-notes
		GetVersionNotes(ctx context.Context, params GetVersionNotesRequest) (*GetVersionNotesResponse, error)

		// UpdateVersionNotes updates the most recent version notes for a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/put-version-notes
		UpdateVersionNotes(ctx context.Context, params UpdateVersionNotesRequest) (*UpdateVersionNotesResponse, error)
	}

	// GetVersionNotesRequest is used to retrieve the version notes for a configuration version.
	GetVersionNotesRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetVersionNotesResponse is returned from a call to GetVersionNotes.
	GetVersionNotesResponse struct {
		Notes string `json:"notes"`
	}

	// UpdateVersionNotesRequest is used to modify the version notes for a configuration version.
	UpdateVersionNotesRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`

		Notes string `json:"notes"`
	}

	// UpdateVersionNotesResponse is returned from a call to UpdateVersionNotes.
	UpdateVersionNotesResponse struct {
		Notes string `json:"notes"`
	}
)

// Validate validates a GetVersionNotesRequest.
func (v GetVersionNotesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates an UpdateVersionNotesRequest.
func (v UpdateVersionNotesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetVersionNotes(ctx context.Context, params GetVersionNotesRequest) (*GetVersionNotesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetVersionNotes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/version-notes",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetVersionNotes request: %w", err)
	}

	var result GetVersionNotesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get version notes request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateVersionNotes(ctx context.Context, params UpdateVersionNotesRequest) (*UpdateVersionNotesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateVersionNotes")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/version-notes",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateVersionNotes request: %w", err)
	}

	var result UpdateVersionNotesResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update version notes request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}
