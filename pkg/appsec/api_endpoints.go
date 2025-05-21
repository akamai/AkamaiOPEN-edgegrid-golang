package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiEndpoints interface supports retrieving the API endpoints associated with a security policy.
	ApiEndpoints interface {
		// GetApiEndpoints lists the API endpoints associated with a security policy.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-api-endpoints
		GetApiEndpoints(ctx context.Context, params GetApiEndpointsRequest) (*GetApiEndpointsResponse, error)
	}

	// GetApiEndpointsRequest is used to retrieve the endpoints associated with a security policy.
	GetApiEndpointsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Name     string `json:"-"`
		PolicyID string `json:"-"`
		ID       int    `json:"-"`
	}

	// GetApiEndpointsResponse is returned from a call to GetApiEndpoints.
	GetApiEndpointsResponse struct {
		APIEndpoints []struct {
			ID               int      `json:"id"`
			Name             string   `json:"name"`
			BasePath         string   `json:"basePath"`
			APIEndPointHosts []string `json:"apiEndPointHosts"`
			StagingVersion   struct {
				Status        string `json:"status"`
				VersionNumber int    `json:"versionNumber"`
			} `json:"stagingVersion"`
			ProductionVersion struct {
				Status        string `json:"status"`
				VersionNumber int    `json:"versionNumber"`
			} `json:"productionVersion"`
			RequestConstraintsEnabled bool `json:"requestConstraintsEnabled"`
		} `json:"apiEndpoints"`
	}
)

// Validate validates a GetApiEndpointsRequest.
func (v GetApiEndpointsRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiEndpoints(ctx context.Context, params GetApiEndpointsRequest) (*GetApiEndpointsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiEndpoints")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var uri string
	if params.PolicyID != "" {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies/%s/api-endpoints",
			params.ConfigID,
			params.Version,
			params.PolicyID)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/api-endpoints",
			params.ConfigID,
			params.Version,
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetApiEndpoints request: %w", err)
	}

	var result GetApiEndpointsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get API endpoints request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Name != "" {
		var filteredResult GetApiEndpointsResponse
		for _, val := range result.APIEndpoints {
			if val.Name == params.Name {
				filteredResult.APIEndpoints = append(filteredResult.APIEndpoints, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}
