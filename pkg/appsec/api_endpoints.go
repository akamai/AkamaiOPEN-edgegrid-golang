package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiEndpoints interface supports retrieving the API endpoints associated with a security policy.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#apiendpoint
	ApiEndpoints interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapiendpoints
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
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetApiEndpoints")

	var rval GetApiEndpointsResponse
	var rvalfiltered GetApiEndpointsResponse

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

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetApiEndpoints request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.Name != "" {
		for _, val := range rval.APIEndpoints {
			if val.Name == params.Name {
				rvalfiltered.APIEndpoints = append(rvalfiltered.APIEndpoints, val)
			}
		}

	} else {
		rvalfiltered = rval
	}
	return &rvalfiltered, nil

}
