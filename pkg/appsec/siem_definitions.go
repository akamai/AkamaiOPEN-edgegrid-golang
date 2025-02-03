package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
)

type (
	// The SiemDefinitions interface supports retrieving the available SIEM versions.
	SiemDefinitions interface {
		// GetSiemDefinitions gets available SIEM versions.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-siem-definitions
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

	uri := "/appsec/v1/siem-definitions"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSiemDefinitions request: %w", err)
	}

	var result GetSiemDefinitionsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get siem definitions request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.SiemDefinitionName != "" {
		var filteredResult GetSiemDefinitionsResponse
		for _, val := range result.SiemDefinitions {
			if val.Name == params.SiemDefinitionName {
				filteredResult.SiemDefinitions = append(filteredResult.SiemDefinitions, val)
			}
		}
		return &filteredResult, nil
	}

	return &result, nil
}
