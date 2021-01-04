package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ApiEndpoints represents a collection of ApiEndpoints
//
// See: ApiEndpoints.GetApiEndpoints()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiEndpoints  contains operations available on ApiEndpoints  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapiendpoints
	ApiEndpoints interface {
		//	GetApiEndpointss(ctx context.Context, params GetApiEndpointssRequest) (*GetApiEndpointssResponse, error)
		GetApiEndpoints(ctx context.Context, params GetApiEndpointsRequest) (*GetApiEndpointsResponse, error)
		//	UpdateApiEndpoints(ctx context.Context, params UpdateApiEndpointsRequest) (*UpdateApiEndpointsResponse, error)
	}

	GetApiEndpointsRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Name     string `json:"-"`
	}

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

// Validate validates GetApiEndpointsRequest
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

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/api-endpoints",
		params.ConfigID,
		params.Version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapiendpoints request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapiendpoints  request failed: %w", err)
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
