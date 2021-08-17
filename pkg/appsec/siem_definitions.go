package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The SiemDefinitions interface supports retrieving the available SIEM versions.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#siem
	SiemDefinitions interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemversions
		GetSiemDefinitions(ctx context.Context, params GetSiemDefinitionsRequest) (*GetSiemDefinitionsResponse, error)
	}

	// GetSiemDefinitionsRequest is used to retrieve the available SIEM versions.
	GetSiemDefinitionsRequest struct {
		ID                 int    `json:"id"`
		SiemDefinitionName string `json:"name"`
	}

	// GetSiemDefinitionsResponse is returned from a call to GetSiemDefinitions.
	GetSiemDefinitionsResponse struct {
		SiemDefinitions []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"siemDefinitions"`
	}
)

func (p *appsec) GetSiemDefinitions(ctx context.Context, params GetSiemDefinitionsRequest) (*GetSiemDefinitionsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetSiemDefinitions")

	var rval GetSiemDefinitionsResponse
	var rvalfiltered GetSiemDefinitionsResponse

	uri := "/appsec/v1/siem-definitions"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSiemDefinitions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetSiemDefinitions request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.SiemDefinitionName != "" {
		for _, val := range rval.SiemDefinitions {
			if val.Name == params.SiemDefinitionName {
				rvalfiltered.SiemDefinitions = append(rvalfiltered.SiemDefinitions, val)
			}
		}

	} else {
		rvalfiltered = rval
	}

	return &rvalfiltered, nil

}
