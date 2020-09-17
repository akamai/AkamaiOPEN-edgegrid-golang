package papi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
)

type (
	GetHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupId           string
		ValidateHostnames bool
	}

	GetHostnamesResponse struct {
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

	CreateHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
	}

	CreateHostnamesResponse struct {
		AccountID       string        `json:"accountId"`
		ContractID      string        `json:"contractId"`
		GroupID         string        `json:"groupId"`
		PropertyID      string        `json:"propertyId"`
		PropertyVersion int           `json:"propertyVersion"`
		Etag            string        `json:"etag"`
		Hostnames       HostnameItems `json:"hostnames"`
	}
)

// Validate validates GetHostnamesRequest
func (ph GetHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ph.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ph.PropertyVersion, validation.Required),
	}.Filter()
}

func (p *papi) GetHostnames(ctx context.Context, params GetHostnamesRequest) (*GetHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/hostnames?contractId=%s&groupId=%s&validateHostnames=%v",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupId,
		params.ValidateHostnames)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the GetHostnames request: %v", err.Error())
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	var hostnames GetHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("GetHostnames request failed: %v", err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w, %s", session.ErrNotFound, getURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &hostnames, nil
}

// Validate validates CreateHostnamesRequest
func (ch CreateHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ch.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ch.PropertyVersion, validation.Required),
	}.Filter()
}

func (p *papi) CreateHostnames(ctx context.Context, params CreateHostnamesRequest) (*CreateHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("CreateHostnames")

	putURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%v/hostnames?contractId=%s&groupId=%s&validateHostnames=%v",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
		params.ValidateHostnames,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create createhostnames request: %w", err)
	}

	req.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))
	var hostnames CreateHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("createhostnames request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", session.ErrNotFound, putURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, session.NewAPIError(resp, logger)
	}

	return &hostnames, nil
}
