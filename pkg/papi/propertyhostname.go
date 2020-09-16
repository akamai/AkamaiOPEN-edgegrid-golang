package papi

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
	"net/http"
)

type (
	GetPropertyHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupId           string
		ValidateHostnames bool
	}

	GetPropertyHostnamesResponse struct {
		AccountID       string        `json:"accountId"`
		ContractID      string        `json:"contractId"`
		GroupID         string        `json:"groupId"`
		PropertyID      string        `json:"propertyId"`
		PropertyVersion int           `json:"propertyVersion"`
		Etag            string        `json:"etag"`
		Hostnames       HostnameItems `json:"hostnames"`
	}

	HostnameItems struct {
		Items []HostnameItem `json:"items"`
	}

	HostnameItem struct {
		CnameType      string `json:"cnameType"`
		EdgeHostnameID string `json:"edgeHostnameId"`
		CnameFrom      string `json:"cnameFrom"`
		CnameTo        string `json:"cnameTo"`
	}
)

// Validate validates GetPropertyHostnamesRequest
func (ph GetPropertyHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ph.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ph.PropertyVersion, validation.Required),
	}.Filter()
}

func (p *papi) GetPropertyHostnames(ctx context.Context, params GetPropertyHostnamesRequest) (*GetPropertyHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/hostnames?contractId=%s&groupId=%s&validateHostnames=%v",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupId,
		params.ValidateHostnames)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the GetPropertyHostnames request: %v", err.Error())
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var hostnames GetPropertyHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("GetPropertyHostnames request failed: %v", err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w, %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &hostnames, nil
}
