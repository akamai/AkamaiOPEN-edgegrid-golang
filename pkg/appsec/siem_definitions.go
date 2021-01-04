package appsec

import (
	"context"
	"fmt"
	"net/http"
)

// SiemDefinitions represents a collection of SiemDefinitions
//
// See: SiemDefinitions.GetSiemDefinitions()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// SiemDefinitions  contains operations available on SiemDefinitions  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemdefinitions
	SiemDefinitions interface {
		GetSiemDefinitions(ctx context.Context, params GetSiemDefinitionsRequest) (*GetSiemDefinitionsResponse, error)
	}

	GetSiemDefinitionsResponse struct {
		SiemDefinitions []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"siemDefinitions"`
	}

	GetSiemDefinitionsRequest struct {
		ID                 int    `json:"id"`
		SiemDefinitionName string `json:"name"`
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
		return nil, fmt.Errorf("failed to create getsiemdefinitions request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getsiemdefinitions  request failed: %w", err)
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
