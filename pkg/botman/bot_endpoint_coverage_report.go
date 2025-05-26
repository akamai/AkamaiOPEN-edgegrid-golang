package botman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The BotEndpointCoverageReport interface supports retrieving bot endpoint coverage report for an account or a specific configuration
	BotEndpointCoverageReport interface {
		// GetBotEndpointCoverageReport https://techdocs.akamai.com/bot-manager/reference/get-bot-endpoint-coverage-report
		GetBotEndpointCoverageReport(ctx context.Context, params GetBotEndpointCoverageReportRequest) (*GetBotEndpointCoverageReportResponse, error)
	}

	// GetBotEndpointCoverageReportRequest is used to retrieve the akamai bot category actions for a policy.
	GetBotEndpointCoverageReportRequest struct {
		ConfigID    int64
		Version     int64
		OperationID string
	}

	// GetBotEndpointCoverageReportResponse is returned from a call to GetBotEndpointCoverageReport.
	GetBotEndpointCoverageReportResponse struct {
		Operations []map[string]interface{} `json:"operations"`
	}
)

// Validate validates GetBotEndpointCoverageReportRequest
func (v GetBotEndpointCoverageReportRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.When(v.Version != 0, validation.Required)),
		"Version":  validation.Validate(v.Version, validation.When(v.ConfigID != 0, validation.Required)),
	}.Filter()
}

func (b *botman) GetBotEndpointCoverageReport(ctx context.Context, params GetBotEndpointCoverageReportRequest) (*GetBotEndpointCoverageReportResponse, error) {
	logger := b.Log(ctx)
	logger.Debug("GetBotEndpointCoverageReport")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri string
	if params.ConfigID != 0 && params.Version != 0 {
		uri = fmt.Sprintf("/appsec/v1/configs/%d/versions/%d/bot-endpoint-coverage-report", params.ConfigID, params.Version)
	} else {
		uri = "/appsec/v1/bot-endpoint-coverage-report"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBotEndpointCoverageReport request: %w", err)
	}

	var result GetBotEndpointCoverageReportResponse
	resp, err := b.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetBotEndpointCoverageReport request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, b.Error(resp)
	}

	var filteredResult GetBotEndpointCoverageReportResponse
	if params.OperationID != "" {
		for _, val := range result.Operations {
			if val["operationId"].(string) == params.OperationID {
				filteredResult.Operations = append(filteredResult.Operations, val)
			}
		}
	} else {
		filteredResult = result
	}
	return &filteredResult, nil
}
