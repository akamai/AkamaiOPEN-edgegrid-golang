package networklists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
)

type (
	// GetNetworkListSubscriptionRequest contains request parameters for GetNetworkListSubscription
	GetNetworkListSubscriptionRequest struct {
		Recipients []string `json:"-"`
		UniqueIDs  []string `json:"-"`
	}

	// GetNetworkListSubscriptionResponse contains response from GetNetworkListSubscription
	GetNetworkListSubscriptionResponse struct {
		Links        NetworkListCreateLink `json:"links"`
		NetworkLists []NetworkListItem     `json:"networkLists"`
	}

	// UpdateNetworkListSubscriptionRequest contains request parameters for UpdateNetworkListSubscription method
	UpdateNetworkListSubscriptionRequest struct {
		Recipients []string `json:"recipients"`
		UniqueIDs  []string `json:"uniqueIds"`
	}

	// UpdateNetworkListSubscriptionResponse contains response from UpdateNetworkListSubscription method
	UpdateNetworkListSubscriptionResponse struct {
		Empty string `json:"-"`
	}

	// RemoveNetworkListSubscriptionResponse contains response from RemoveNetworkListSubscription method
	RemoveNetworkListSubscriptionResponse struct {
		Empty string `json:"-"`
	}

	// RemoveNetworkListSubscriptionRequest contains request parameters for RemoveNetworkListSubscription method
	RemoveNetworkListSubscriptionRequest struct {
		Recipients []string `json:"recipients"`
		UniqueIDs  []string `json:"uniqueIds"`
	}

	// Recipients contains recipients
	Recipients struct {
		Recipients string `json:"notificationRecipients"`
	}
)

func (p *networklists) GetNetworkListSubscription(ctx context.Context, _ GetNetworkListSubscriptionRequest) (*GetNetworkListSubscriptionResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetNetworkListSubscription")

	var rval GetNetworkListSubscriptionResponse

	uri := "/network-list/v2/notifications/subscriptions"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getnetworklistsubscription request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getnetworklistsubscription  request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}

func (p *networklists) UpdateNetworkListSubscription(ctx context.Context, params UpdateNetworkListSubscriptionRequest) (*UpdateNetworkListSubscriptionResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkListSubscription")

	postURL := "/network-list/v2/notifications/subscribe"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create NetworkListSubscriptionrequest: %s", err.Error())
	}

	var rval UpdateNetworkListSubscriptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("remove NetworkListSubscription request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *networklists) RemoveNetworkListSubscription(ctx context.Context, params RemoveNetworkListSubscriptionRequest) (*RemoveNetworkListSubscriptionResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateNetworkListSubscription")

	postURL := "/network-list/v2/notifications/unsubscribe"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create NetworkListSubscriptionrequest: %s", err.Error())
	}

	var rval RemoveNetworkListSubscriptionResponse
	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("remove NetworkListSubscription request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
