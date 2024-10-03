package cps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetDeploymentScheduleRequest contains parameters for GetDeploymentSchedule
	GetDeploymentScheduleRequest struct {
		ChangeID     int
		EnrollmentID int
	}

	// UpdateDeploymentScheduleRequest contains parameters for UpdateDeploymentSchedule
	UpdateDeploymentScheduleRequest struct {
		ChangeID     int
		EnrollmentID int
		DeploymentSchedule
	}

	// UpdateDeploymentScheduleResponse contains response for UpdateDeploymentSchedule
	UpdateDeploymentScheduleResponse struct {
		Change string `json:"change"`
	}

	// DeploymentSchedule contains the schedule for when you want this change deploy
	DeploymentSchedule struct {
		NotAfter  *string `json:"notAfter,omitempty"`
		NotBefore *string `json:"notBefore,omitempty"`
	}
)

// Validate validates GetDeploymentScheduleRequest
func (c GetDeploymentScheduleRequest) Validate() error {
	return validation.Errors{
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
	}.Filter()
}

// Validate validates UpdateDeploymentScheduleRequest
func (c UpdateDeploymentScheduleRequest) Validate() error {
	return validation.Errors{
		"ChangeID":     validation.Validate(c.ChangeID, validation.Required),
		"EnrollmentID": validation.Validate(c.EnrollmentID, validation.Required),
	}.Filter()
}

var (
	// ErrGetDeploymentSchedule is returned when GetDeploymentSchedule fails
	ErrGetDeploymentSchedule = errors.New("get deployment schedule")
	// ErrUpdateDeploymentSchedule is returned when UpdateDeploymentSchedule fails
	ErrUpdateDeploymentSchedule = errors.New("update deployment schedule")
)

func (c *cps) GetDeploymentSchedule(ctx context.Context, params GetDeploymentScheduleRequest) (*DeploymentSchedule, error) {
	logger := c.Log(ctx)
	logger.Debug("GetDeploymentSchedule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetDeploymentSchedule, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/deployment-schedule", params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetDeploymentSchedule, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.deployment-schedule.v1+json")

	var result DeploymentSchedule
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetDeploymentSchedule, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetDeploymentSchedule, c.Error(resp))
	}

	return &result, nil
}

func (c *cps) UpdateDeploymentSchedule(ctx context.Context, params UpdateDeploymentScheduleRequest) (*UpdateDeploymentScheduleResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateDeploymentSchedule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateDeploymentSchedule, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d/deployment-schedule", params.EnrollmentID, params.ChangeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateDeploymentSchedule, err)
	}
	req.Header.Set("Accept", "application/vnd.akamai.cps.change-id.v1+json")
	req.Header.Set("Content-Type", "application/vnd.akamai.cps.deployment-schedule.v1+json; charset=utf-8")

	var result UpdateDeploymentScheduleResponse
	resp, err := c.Exec(req, &result, params.DeploymentSchedule)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateDeploymentSchedule, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateDeploymentSchedule, c.Error(resp))
	}

	return &result, nil
}
