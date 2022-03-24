package hapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// EdgeHostnames contains operations available on Edge Hostname resource
	// See: https://developer.akamai.com/api/core_features/edge_hostnames/v1.html#edgehostname
	EdgeHostnames interface {
		// DeleteEdgeHostname allows deleting a specific edge hostname.
		// You must have an Admin or Technical role in order to delete an edge hostname.
		// You can delete any hostname that’s not currently part of an active Property Manager configuration.
		//
		// See: https://developer.akamai.com/api/core_features/edge_hostnames/v1.html#deleteedgehostnamebyname
		DeleteEdgeHostname(context.Context, DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error)

		// GetEdgeHostname gets a specific edge hostname's details including its product ID, IP version behavior,
		// and China CDN or Edge IP Binding status.
		//
		// See: https://techdocs.akamai.com/edge-hostnames/reference/get-edgehostnameid
		GetEdgeHostname(context.Context, int) (*GetEdgeHostnameResponse, error)
	}

	// DeleteEdgeHostnameRequest is used to delete edge hostname
	DeleteEdgeHostnameRequest struct {
		DNSZone           string
		RecordName        string
		StatusUpdateEmail []string
		Comments          string
	}

	// DeleteEdgeHostnameResponse is a response from deleting edge hostname
	DeleteEdgeHostnameResponse struct {
		Action            string         `json:"action"`
		ChangeID          int            `json:"changeId"`
		Comments          string         `json:"comments"`
		Status            string         `json:"status"`
		StatusMessage     string         `json:"statusMessage"`
		StatusUpdateDate  string         `json:"statusUpdateDate"`
		StatusUpdateEmail string         `json:"statusUpdateEmail"`
		SubmitDate        string         `json:"submitDate"`
		Submitter         string         `json:"submitter"`
		SubmitterEmail    string         `json:"submitterEmail"`
		EdgeHostnames     []EdgeHostname `json:"edgeHostnames"`
	}

	// EdgeHostname represents edge hostname part of DeleteEdgeHostnameResponse
	EdgeHostname struct {
		EdgeHostnameID    int      `json:"edgeHostnameId"`
		RecordName        string   `json:"recordName"`
		DNSZone           string   `json:"dnsZone"`
		SecurityType      string   `json:"securityType"`
		UseDefaultTTL     bool     `json:"useDefaultTtl"`
		UseDefaultMap     bool     `json:"useDefaultMap"`
		TTL               int      `json:"ttl"`
		Map               string   `json:"map"`
		SlotNumber        int      `json:"slotNumber"`
		IPVersionBehavior string   `json:"ipVersionBehavior"`
		Comments          string   `json:"comments"`
		ChinaCDN          ChinaCDN `json:"chinaCdn"`
	}

	// ChinaCDN represents China CDN settings of EdgeHostname
	ChinaCDN struct {
		IsChinaCDN        bool   `json:"isChinaCdn"`
		CustomChinaCDNMap string `json:"customChinaCdnMap,omitempty"`
	}

	// GetEdgeHostnameResponse represents edge hostname
	GetEdgeHostnameResponse struct {
		EdgeHostnameID         int      `json:"edgeHostnameId"`
		RecordName             string   `json:"recordName"`
		DNSZone                string   `json:"dnsZone"`
		SecurityType           string   `json:"securityType"`
		UseDefaultTTL          bool     `json:"useDefaultTtl"`
		UseDefaultMap          bool     `json:"useDefaultMap"`
		IPVersionBehavior      string   `json:"ipVersionBehavior"`
		ProductID              string   `json:"productId"`
		TTL                    int      `json:"ttl"`
		Map                    string   `json:"map,omitempty"`
		SlotNumber             int      `json:"slotNumber,omitempty"`
		Comments               string   `json:"comments"`
		SerialNumber           int      `json:"serialNumber,omitempty"`
		CustomTarget           string   `json:"customTarget,omitempty"`
		ChinaCdn               ChinaCDN `json:"chinaCdn,omitempty"`
		IsEdgeIPBindingEnabled bool     `json:"isEdgeIPBindingEnabled,omitempty"`
	}
)

// Validate validates DeleteEdgeHostnameRequest
func (r DeleteEdgeHostnameRequest) Validate() error {
	return validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	}.Filter()
}

var (
	// ErrDeleteEdgeHostname represents error when deleting edge hostname fails
	ErrDeleteEdgeHostname = errors.New("delete edge hostname")
	// ErrGetEdgeHostname represents error when getting edge hostname fails
	ErrGetEdgeHostname = errors.New("get edge hostname")
)

func (h *hapi) DeleteEdgeHostname(ctx context.Context, params DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteEdgeHostname, ErrStructValidation, err)
	}

	logger := h.Log(ctx)
	logger.Debug("DeleteEdgeHostname")

	uri := fmt.Sprintf(
		"/hapi/v1/dns-zones/%s/edge-hostnames/%s",
		params.DNSZone,
		params.RecordName,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteEdgeHostname, err)
	}

	q := req.URL.Query()
	if len(params.StatusUpdateEmail) > 0 {
		emails := strings.Join(params.StatusUpdateEmail, ",")
		q.Add("statusUpdateEmail", emails)
	}
	if params.Comments != "" {
		q.Add("comments", params.Comments)
	}
	req.URL.RawQuery = q.Encode()

	var rval DeleteEdgeHostnameResponse

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteEdgeHostname, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrDeleteEdgeHostname, h.Error(resp))
	}

	return &rval, nil
}

func (h *hapi) GetEdgeHostname(ctx context.Context, edgeHostnameID int) (*GetEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetEdgeHostname")

	uri := fmt.Sprintf("/hapi/v1/edge-hostnames/%d", edgeHostnameID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeHostname, err)
	}

	var rval GetEdgeHostnameResponse

	resp, err := h.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeHostname, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeHostname, h.Error(resp))
	}

	return &rval, nil
}
