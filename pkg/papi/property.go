package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/spf13/cast"
)

type (
	// PropertyCloneFrom optionally identifies another property instance to clone when making a POST request to create a new property
	PropertyCloneFrom struct {
		CloneFromVersionEtag string `json:"cloneFromVersionEtag"`
		CopyHostnames        bool   `json:"copyHostnames"`
		PropertyID           string `json:"propertyId"`
		Version              int    `json:"version"`
	}

	// Property contains configuration data to apply to edge content.
	Property struct {
		AccountID         string             `json:"accountId"`
		AssetID           string             `json:"assetId"`
		CloneFrom         *PropertyCloneFrom `json:"cloneFrom,omitempty"`
		ContactID         string             `json:"contractId"`
		GroupID           string             `json:"groupId"`
		LatestVersion     int                `json:"latestVersion"`
		Note              string             `json:"note"`
		ProductID         string             `json:"productId"`
		ProductionVersion *int               `json:"productionVersion,omitempty"`
		PropertyID        string             `json:"propertyId"`
		PropertyName      string             `json:"propertyName"`
		RuleFormat        string             `json:"ruleFormat"`
		StagingVersion    *int               `json:"stagingVersion,omitempty"`
	}

	// PropertiesItems is an array of properties
	PropertiesItems struct {
		Items []*Property `json:"items"`
	}

	// GetPropertiesRequest is the argument for GetProperties
	GetPropertiesRequest struct {
		ContractID string `json:"contractId"`
		GroupID    string `json:"groupId"`
	}

	// GetPropertiesResponse is the response for GetProperties
	GetPropertiesResponse struct {
		Properties PropertiesItems `json:"properties"`
	}
)

func (p *papi) GetProperties(ctx context.Context, r GetPropertiesRequest) (*GetPropertiesResponse, error) {
	var rval GetPropertiesResponse

	logger := p.Log(ctx)
	logger.Debug("GetProperties")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/papi/v1/properties?contractId=%s&groupId=%s", r.ContractID, r.GroupID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getproperties request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getproperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &rval, nil
}
