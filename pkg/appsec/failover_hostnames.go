package appsec

import (
	"context"
	"fmt"
	"net/http"
)

// FailoverHostnames represents a collection of FailoverHostnames
//
// See: FailoverHostnames.GetFailoverHostnames()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// FailoverHostnames  contains operations available on FailoverHostnames  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getfailoverhostnames
	FailoverHostnames interface {
		GetFailoverHostnames(ctx context.Context, params GetFailoverHostnamesRequest) (*GetFailoverHostnamesResponse, error)
	}

	GetFailoverHostnamesRequest struct {
		ConfigID int `json:"-"`
	}

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

	var rval GetFailoverHostnamesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/failover-hostnames",
		params.ConfigID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getfailoverhostnamess request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getfailoverhostnamess request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
