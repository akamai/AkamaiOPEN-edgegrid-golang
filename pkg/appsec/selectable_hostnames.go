package appsec

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// The SelectableHostnames interface supports retrieving the hostnames that a given configuration version
	// has the ability to protect. Hostnames may show as error hosts when they aren’t currently available. for
	// example, when a contract expires.
	SelectableHostnames interface {
		// GetSelectableHostnames lists the hostnames that a given configuration version has the ability to protect.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-selectable-hostnames
		GetSelectableHostnames(ctx context.Context, params GetSelectableHostnamesRequest) (*GetSelectableHostnamesResponse, error)
	}

	// GetSelectableHostnamesRequest is used to retrieve the selectable hostnames for a configuration.
	GetSelectableHostnamesRequest struct {
		ConfigID   int    `json:"configId"`
		Version    int    `json:"version"`
		ContractID string `json:"-"`
		GroupID    int    `json:"-"`
	}

	// GetSelectableHostnamesResponse is returned from a call to GetSelectableHostnames.
	GetSelectableHostnamesResponse struct {
		AvailableSet []struct {
			ActiveInProduction     bool   `json:"activeInProduction,omitempty"`
			ActiveInStaging        bool   `json:"activeInStaging,omitempty"`
			ArlInclusion           bool   `json:"arlInclusion,omitempty"`
			Hostname               string `json:"hostname,omitempty"`
			ConfigIDInProduction   int    `json:"configIdInProduction,omitempty"`
			ConfigNameInProduction string `json:"configNameInProduction,omitempty"`
		} `json:"availableSet,omitempty"`
		ConfigID                int  `json:"configId,omitempty"`
		ConfigVersion           int  `json:"configVersion,omitempty"`
		ProtectARLInclusionHost bool `json:"protectARLInclusionHost,omitempty"`
	}
)

func (p *appsec) GetSelectableHostnames(ctx context.Context, params GetSelectableHostnamesRequest) (*GetSelectableHostnamesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetSelectableHostnamess")

	var uri string

	if params.ConfigID != 0 {
		uri = fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/selectable-hostnames",
			params.ConfigID,
			params.Version)
	} else {
		uri = fmt.Sprintf(
			"/appsec/v1/contracts/%s/groups/%d/selectable-hostnames",
			params.ContractID,
			params.GroupID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetSelectableHostnames request: %w", err)
	}

	var result GetSelectableHostnamesResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get selectable hostnames request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
