package appsec

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ApiHostnameCoverageOverlapping represents a collection of ApiHostnameCoverageOverlapping
//
// See: ApiHostnameCoverageOverlapping.GetApiHostnameCoverageOverlapping()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiHostnameCoverageOverlapping  contains operations available on ApiHostnameCoverageOverlapping  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapihostnamecoverageoverlapping
	ApiHostnameCoverageOverlapping interface {
		GetApiHostnameCoverageOverlapping(ctx context.Context, params GetApiHostnameCoverageOverlappingRequest) (*GetApiHostnameCoverageOverlappingResponse, error)
	}

	GetApiHostnameCoverageOverlappingRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	GetApiHostnameCoverageOverlappingResponse struct {
		OverLappingList []struct {
			ConfigID      int      `json:"configId"`
			ConfigName    string   `json:"configName"`
			ConfigVersion int      `json:"configVersion"`
			ContractID    string   `json:"contractId"`
			ContractName  string   `json:"contractName"`
			VersionTags   []string `json:"versionTags"`
		} `json:"overLappingList"`
	}
)

// Validate validates GetApiHostnameCoverageOverlappingRequest
func (v GetApiHostnameCoverageOverlappingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiHostnameCoverageOverlapping(ctx context.Context, params GetApiHostnameCoverageOverlappingRequest) (*GetApiHostnameCoverageOverlappingResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverageOverlapping")

	var rval GetApiHostnameCoverageOverlappingResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/hostname-coverage/match-targets?hostname=www.example.com",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapihostnamecoverageoverlapping request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapihostnamecoverageoverlapping  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
