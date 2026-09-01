package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The AdvancedSettingsAsePenaltyBox interface supports retrieving or modifying the ASE Penalty Box setting.
	AdvancedSettingsAsePenaltyBox interface {
		// GetAdvancedSettingsAsePenaltyBox retrieves the ASE Penalty Box setting.
		//
		// See: TBD
		GetAdvancedSettingsAsePenaltyBox(ctx context.Context, params GetAdvancedSettingsAsePenaltyBoxRequest) (*GetAdvancedSettingsAsePenaltyBoxResponse, error)

		// UpdateAdvancedSettingsAsePenaltyBox modifies the ASE Penalty Box setting.
		//
		// See: TBD
		UpdateAdvancedSettingsAsePenaltyBox(ctx context.Context, params UpdateAdvancedSettingsAsePenaltyBoxRequest) (*UpdateAdvancedSettingsAsePenaltyBoxResponse, error)

		// RemoveAdvancedSettingsAsePenaltyBox removes the ASE Penalty Box setting.
		//
		RemoveAdvancedSettingsAsePenaltyBox(ctx context.Context, params RemoveAdvancedSettingsAsePenaltyBoxRequest) (*RemoveAdvancedSettingsAsePenaltyBoxResponse, error)
	}

	// GetAdvancedSettingsAsePenaltyBoxRequest is used to retrieve the AsePenaltyBox setting.
	GetAdvancedSettingsAsePenaltyBoxRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// GetAdvancedSettingsAsePenaltyBoxResponse returns the AsePenaltyBox setting.
	GetAdvancedSettingsAsePenaltyBoxResponse struct {
		RequestCount            int                      `json:"requestCount,omitempty"`
		BlockDuration           int                      `json:"blockDuration,omitempty"`
		ClientIdentifiers       []string                 `json:"clientIdentifiers,omitempty"`
		AkamaiManagedExclusions *AkamaiManagedExclusions `json:"akamaiManagedExclusions,omitempty"`
		QualificationExclusions *QualificationExclusions `json:"qualificationExclusions,omitempty"`
	}

	// UpdateAdvancedSettingsAsePenaltyBoxRequest is used to update the AsePenaltyBox setting.
	UpdateAdvancedSettingsAsePenaltyBoxRequest struct {
		ConfigID                int                      `json:"-"`
		Version                 int                      `json:"-"`
		BlockDuration           int                      `json:"blockDuration,omitempty"`
		QualificationExclusions *QualificationExclusions `json:"qualificationExclusions,omitempty"`
	}

	// UpdateAdvancedSettingsAsePenaltyBoxResponse returns the result of updating the AsePenaltyBox setting.
	UpdateAdvancedSettingsAsePenaltyBoxResponse struct {
		RequestCount            int                      `json:"requestCount,omitempty"`
		BlockDuration           int                      `json:"blockDuration,omitempty"`
		ClientIdentifiers       []string                 `json:"clientIdentifiers,omitempty"`
		AkamaiManagedExclusions *AkamaiManagedExclusions `json:"akamaiManagedExclusions,omitempty"`
		QualificationExclusions *QualificationExclusions `json:"qualificationExclusions,omitempty"`
	}

	// RemoveAdvancedSettingsAsePenaltyBoxRequest is used to clear the AsePenaltyBox setting.
	RemoveAdvancedSettingsAsePenaltyBoxRequest struct {
		ConfigID int `json:"-"`
		Version  int `json:"-"`
	}

	// RemoveAdvancedSettingsAsePenaltyBoxResponse returns the result of clearing the AsePenaltyBox setting.
	RemoveAdvancedSettingsAsePenaltyBoxResponse struct {
		RequestCount            int                      `json:"requestCount,omitempty"`
		BlockDuration           int                      `json:"blockDuration,omitempty"`
		ClientIdentifiers       []string                 `json:"clientIdentifiers,omitempty"`
		AkamaiManagedExclusions *AkamaiManagedExclusions `json:"akamaiManagedExclusions,omitempty"`
		QualificationExclusions *QualificationExclusions `json:"qualificationExclusions,omitempty"`
	}

	// QualificationExclusions returns the attack groups and rules that are excluded from the penalty box.
	QualificationExclusions struct {
		AttackGroups []string `json:"attackGroups,omitempty"`
		Rules        []int    `json:"rules,omitempty"`
	}

	// AkamaiManagedExclusions returns the rules managed by Akamai that are excluded from the penalty box.
	AkamaiManagedExclusions struct {
		Rules       []int  `json:"rules,omitempty"`
		LastUpdated string `json:"lastUpdated,omitempty"`
	}
)

// Validate validates GetAdvancedSettingsAsePenaltyBoxRequest
func (v GetAdvancedSettingsAsePenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsAsePenaltyBoxRequest
func (v UpdateAdvancedSettingsAsePenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

// Validate validates UpdateAdvancedSettingsAsePenaltyBoxRequest
func (v RemoveAdvancedSettingsAsePenaltyBoxRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetAdvancedSettingsAsePenaltyBox(ctx context.Context, params GetAdvancedSettingsAsePenaltyBoxRequest) (*GetAdvancedSettingsAsePenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetAdvancedSettingsAsePenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri = fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/ase-penalty-box",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAdvancedSettingsAsePenaltyBox request: %w", err)
	}

	var result GetAdvancedSettingsAsePenaltyBoxResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get advanced settings ASE Penalty Box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) UpdateAdvancedSettingsAsePenaltyBox(ctx context.Context, params UpdateAdvancedSettingsAsePenaltyBoxRequest) (*UpdateAdvancedSettingsAsePenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateAdvancedSettingsAsePenaltyBox")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri = fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/advanced-settings/ase-penalty-box",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateAdvancedSettingsAsePenaltyBox request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result UpdateAdvancedSettingsAsePenaltyBoxResponse
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("update advanced settings ASE Penalty Box request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}

func (p *appsec) RemoveAdvancedSettingsAsePenaltyBox(ctx context.Context, params RemoveAdvancedSettingsAsePenaltyBoxRequest) (*RemoveAdvancedSettingsAsePenaltyBoxResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("RemoveAdvancedSettingsAsePenaltyBox")

	request := UpdateAdvancedSettingsAsePenaltyBoxRequest{
		ConfigID:      params.ConfigID,
		Version:       params.Version,
		BlockDuration: 10,
	}

	resp, err := p.UpdateAdvancedSettingsAsePenaltyBox(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("remove advanced settings ASE Penalty Box request failed: %w", err)
	}

	response := RemoveAdvancedSettingsAsePenaltyBoxResponse{
		RequestCount:            resp.RequestCount,
		BlockDuration:           resp.BlockDuration,
		ClientIdentifiers:       resp.ClientIdentifiers,
		AkamaiManagedExclusions: resp.AkamaiManagedExclusions,
		QualificationExclusions: resp.QualificationExclusions,
	}
	return &response, nil
}
