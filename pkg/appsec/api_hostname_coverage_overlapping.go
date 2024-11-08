package appsec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ApiHostnameCoverageOverlapping interface supports listing the configuration versions that
	// contain a hostname also included in the given configuration version.
	ApiHostnameCoverageOverlapping interface {
		// GetApiHostnameCoverageOverlapping lists the configuration versions that contain a hostname also included in the current configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-hostname-coverage-overlapping
		GetApiHostnameCoverageOverlapping(ctx context.Context, params GetApiHostnameCoverageOverlappingRequest) (*GetApiHostnameCoverageOverlappingResponse, error)
	}

	// GetApiHostnameCoverageOverlappingRequest is used to retrieve the configuration versions that contain a hostname included in the current configuration version.
	GetApiHostnameCoverageOverlappingRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	// GetApiHostnameCoverageOverlappingResponse is returned from a call to GetApiHostnameCoverageOverlapping.
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

// Validate validates a GetApiHostnameCoverageOverlappingRequest.
func (v GetApiHostnameCoverageOverlappingRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
	}.Filter()
}

func (p *appsec) GetApiHostnameCoverageOverlapping(ctx context.Context, params GetApiHostnameCoverageOverlappingRequest) (*GetApiHostnameCoverageOverlappingResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverageOverlapping")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/hostname-coverage/overlapping?hostname=%s",
		params.ConfigID,
		params.Version,
		params.Hostname,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetApiHostnameCoverageOverlapping request: %w", err)
	}

	var result GetApiHostnameCoverageOverlappingResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get API hostname coverage overlapping request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
