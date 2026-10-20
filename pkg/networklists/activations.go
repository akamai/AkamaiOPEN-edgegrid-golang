package networklists

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// GetActivationsRequest contains request parameters for getting activation status
	GetActivationsRequest struct {
		UniqueID     string `json:"-"`
		Action       string `json:"-"`
		Network      string `json:"network"`
		ActivationID int    `json:"activationId"`
	}

	// GetActivationRequest contains request parameters for getting activation details
	GetActivationRequest struct {
		ActivationID int `json:"activationId"`
	}

	// GetActivationsResponse contains response with activation status
	GetActivationsResponse struct {
		ActivationID       int             `json:"activationId"`
		ActivationComments string          `json:"activationComments"`
		ActivationStatus   string          `json:"activationStatus"`
		SyncPoint          int             `json:"syncPoint"`
		UniqueID           string          `json:"uniqueId"`
		Fast               bool            `json:"fast"`
		DispatchCount      int             `json:"dispatchCount"`
		Links              ActivationLinks `json:"links"`
	}

	// GetActivationResponse contains response with activation details
	GetActivationResponse struct {
		ActivationID     int                   `json:"activationId"`
		CreateDate       time.Time             `json:"createDate"`
		CreatedBy        string                `json:"createdBy"`
		Environment      string                `json:"environment"`
		Fast             bool                  `json:"fast"`
		ActivationStatus string                `json:"status"`
		NetworkList      ActivationNetworkList `json:"networkList"`
	}

	// CreateActivationsRequest contains request parameters for creating new activation
	CreateActivationsRequest struct {
		UniqueID               string   `json:"-"`
		Action                 string   `json:"-"`
		Network                string   `json:"network"`
		Comments               string   `json:"comments"`
		NotificationRecipients []string `json:"notificationRecipients"`
	}

	// CreateActivationsResponse contains response after creating new activation
	CreateActivationsResponse struct {
		ActivationID       int             `json:"activationId"`
		ActivationComments string          `json:"activationComments"`
		ActivationStatus   string          `json:"activationStatus"`
		SyncPoint          int             `json:"syncPoint"`
		UniqueID           string          `json:"uniqueId"`
		Fast               bool            `json:"fast"`
		DispatchCount      int             `json:"dispatchCount"`
		Links              ActivationLinks `json:"links"`
	}

	// RemoveActivationsRequest contains request parameters of Activation to deactivate
	RemoveActivationsRequest struct {
		UniqueID               string   `json:"-"`
		ActivationID           int      `json:"-"`
		Action                 string   `json:"action"`
		Network                string   `json:"network"`
		Comments               string   `json:"comments"`
		NotificationRecipients []string `json:"notificationRecipients"`
	}

	// RemoveActivationsResponse contains response of Activation deactivation
	RemoveActivationsResponse struct {
		ActivationID       int             `json:"activationId"`
		ActivationComments string          `json:"activationComments"`
		ActivationStatus   string          `json:"activationStatus"`
		SyncPoint          int             `json:"syncPoint"`
		UniqueID           string          `json:"uniqueId"`
		Fast               bool            `json:"fast"`
		DispatchCount      int             `json:"dispatchCount"`
		Links              ActivationLinks `json:"links"`
	}

	// ActivationLinks contains hypermedia links for activation responses.
	ActivationLinks struct {
		AppendItems        ActivationActionLink `json:"appendItems"`
		Retrieve           ActivationLink       `json:"retrieve"`
		StatusInProduction ActivationLink       `json:"statusInProduction"`
		StatusInStaging    ActivationLink       `json:"statusInStaging"`
		SyncPointHistory   ActivationLink       `json:"syncPointHistory"`
		Update             ActivationActionLink `json:"update"`
		ActivationDetails  ActivationLink       `json:"activationDetails"`
	}

	// ActivationNetworkList describes network list details in an activation response.
	ActivationNetworkList struct {
		ActivationComments string                     `json:"activationComments"`
		ActivationStatus   string                     `json:"activationStatus"`
		Links              ActivationNetworkListLinks `json:"links"`
		SyncPoint          int                        `json:"syncPoint"`
		UniqueID           string                     `json:"uniqueId"`
	}

	// ActivationNetworkListLinks contains hypermedia links for network list details.
	ActivationNetworkListLinks struct {
		AppendItems        ActivationActionLink `json:"appendItems"`
		Retrieve           ActivationLink       `json:"retrieve"`
		StatusInProduction ActivationLink       `json:"statusInProduction"`
		StatusInStaging    ActivationLink       `json:"statusInStaging"`
		SyncPointHistory   ActivationLink       `json:"syncPointHistory"`
		Update             ActivationActionLink `json:"update"`
	}

	// ActivationActionLink describes a hypermedia link that includes an HTTP method.
	ActivationActionLink struct {
		Href   string `json:"href"`
		Method string `json:"method"`
	}

	// ActivationLink describes a standard hypermedia reference.
	ActivationLink struct {
		Href string `json:"href"`
	}

	// ActivationValue is used to create an "enum" of possible Activation.ActivationType values
	ActivationValue string

	// NetworkValue is used to create an "enum" of possible Activation.Network values
	NetworkValue string

	// StatusValue is used to create an "enum" of possible Activation.Status values
	StatusValue string
)

const (
	// ActivationTypeActivate Activation.ActivationType value ACTIVATE
	ActivationTypeActivate ActivationValue = "ACTIVATE"
	// ActivationTypeDeactivate Activation.ActivationType value DEACTIVATE
	ActivationTypeDeactivate ActivationValue = "DEACTIVATE"

	// NetworkProduction Activation.Network value PRODUCTION
	NetworkProduction NetworkValue = "PRODUCTION"
	// NetworkStaging Activation.Network value STAGING
	NetworkStaging NetworkValue = "STAGING"

	// StatusActive Activation.Status value ACTIVE
	StatusActive StatusValue = "ACTIVATED"
	// StatusInactive Activation.Status value INACTIVE
	StatusInactive StatusValue = "INACTIVE"
	// StatusPending Activation.Status value RECEIVED
	StatusPending StatusValue = "RECEIVED"
	// StatusAborted Activation.Status value ABORTED
	StatusAborted StatusValue = "ABORTED"
	// StatusFailed Activation.Status value FAILED
	StatusFailed StatusValue = "FAILED"
	// StatusDeactivated Activation.Status value DEACTIVATED
	StatusDeactivated StatusValue = "DEACTIVATED"
	// StatusPendingDeactivation Activation.Status value PENDING_DEACTIVATION
	StatusPendingDeactivation StatusValue = "PENDING_DEACTIVATION"
	// StatusNew Activation.Status value NEW
	StatusNew StatusValue = "NEW"
)

// Validate validates GetActivationsRequest
func (v GetActivationsRequest) Validate() error {
	return validation.Errors{
		"UniqueID": validation.Validate(v.UniqueID, validation.Required),
	}.Filter()
}

// Validate validates GetActivationRequest
func (v GetActivationRequest) Validate() error {
	return validation.Errors{
		"ActivationID": validation.Validate(v.ActivationID, validation.Required),
	}.Filter()
}

func (p *networklists) GetActivations(ctx context.Context, params GetActivationsRequest) (*GetActivationsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetActivations")

	var rval GetActivationsResponse

	uri := fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/status",
		params.UniqueID,
		params.Network,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getactivations request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getactivations request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *networklists) GetActivation(ctx context.Context, params GetActivationRequest) (*GetActivationResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	logger := p.Log(ctx)
	logger.Debug("GetActivation")

	var rval GetActivationResponse

	uri := fmt.Sprintf("/network-list/v2/activations/%d",
		params.ActivationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getactivation request: %s", err.Error())
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getactivation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}

func (p *networklists) CreateActivations(ctx context.Context, params CreateActivationsRequest) (*CreateActivationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateActivations")

	uri := fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/activate",
		params.UniqueID,
		params.Network,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create create activation request: %s", err.Error())
	}

	var rval CreateActivationsResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("create activation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	var rvalget CreateActivationsResponse

	uriget := fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/status",
		params.UniqueID,
		params.Network,
	)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, uriget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation request: %s", err.Error())
	}

	resp, err = p.Exec(req, &rvalget)
	if err != nil {
		return nil, fmt.Errorf("get activation request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rvalget, nil

}

func (p *networklists) RemoveActivations(ctx context.Context, params RemoveActivationsRequest) (*RemoveActivationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("RemoveActivations")

	uri := fmt.Sprintf("/network-list/v2/network-lists/%s/environments/%s/deactivate",
		params.UniqueID,
		params.Network,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create remove activation request: %s", err.Error())
	}

	var rval RemoveActivationsResponse

	resp, err := p.Exec(req, &rval, params)
	if err != nil {
		return nil, fmt.Errorf("remove activation request failed: %s", err.Error())
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil
}
