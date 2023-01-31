package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The FailoverHostnames interface supports retrieving the failover hostnames in a configuration.
	FailoverHostnames interface {
		// GetFailoverHostnames returns a list of the failover hostnames in a configuration.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-failover-hostnames
		GetFailoverHostnames(ctx context.Context, params GetFailoverHostnamesRequest) (*GetFailoverHostnamesResponse, error)
	}

	// GetFailoverHostnamesRequest is used to retrieve the failover hostnames for a configuration.
	GetFailoverHostnamesRequest struct {
		ConfigID int `json:"-"`
	}

	// GetFailoverHostnamesResponse is returned from a call to GetFailoverHostnames.
	GetFailoverHostnamesResponse struct {
		ConfigID      int `json:"-"`
		ConfigVersion int `json:"-"`
		HostnameList  []struct {
			Hostname string `json:"hostname"`
		} `json:"hostnameList"`
	}
)

func (p *appsec) GetFailoverHostnames(ctx context.Context, params GetFailoverHostnamesRequest) (*GetFailoverHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetFailoverHostnames")

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/failover-hostnames",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetFailoverHostnames request: %w", err)
	}

	var result GetFailoverHostnamesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get failover hostnames request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
