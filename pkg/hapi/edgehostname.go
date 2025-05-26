package hapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
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

	// UpdateEdgeHostnameRequest is a request used to update edge hostname
	UpdateEdgeHostnameRequest struct {
		DNSZone           string
		RecordName        string
		StatusUpdateEmail []string
		Comments          string
		Body              []UpdateEdgeHostnameRequestBody
	}

	// UpdateEdgeHostnameRequestBody is a request's body used to update edge hostname
	UpdateEdgeHostnameRequestBody struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	// UpdateEdgeHostnameResponse is a response from deleting edge hostname
	UpdateEdgeHostnameResponse struct {
		Action            string         `json:"action,omitempty"`
		ChangeID          int            `json:"changeId,omitempty"`
		Comments          string         `json:"comments,omitempty"`
		Status            string         `json:"status,omitempty"`
		StatusMessage     string         `json:"statusMessage,omitempty"`
		StatusUpdateDate  string         `json:"statusUpdateDate,omitempty"`
		StatusUpdateEmail string         `json:"statusUpdateEmail,omitempty"`
		SubmitDate        string         `json:"submitDate,omitempty"`
		Submitter         string         `json:"submitter,omitempty"`
		SubmitterEmail    string         `json:"submitterEmail,omitempty"`
		EdgeHostnames     []EdgeHostname `json:"edgeHostnames,omitempty"`
	}

	// EdgeHostname represents edge hostname part of DeleteEdgeHostnameResponse and UpdateEdgeHostnameResponse
	EdgeHostname struct {
		EdgeHostnameID         int       `json:"edgeHostnameId,omitempty"`
		RecordName             string    `json:"recordName"`
		DNSZone                string    `json:"dnsZone"`
		SecurityType           string    `json:"securityType"`
		UseDefaultTTL          bool      `json:"useDefaultTtl"`
		UseDefaultMap          bool      `json:"useDefaultMap"`
		TTL                    int       `json:"ttl"`
		Map                    string    `json:"map,omitempty"`
		SlotNumber             int       `json:"slotNumber,omitempty"`
		IPVersionBehavior      string    `json:"ipVersionBehavior,omitempty"`
		Comments               string    `json:"comments,omitempty"`
		ChinaCDN               ChinaCDN  `json:"chinaCdn,omitempty"`
		CustomTarget           string    `json:"customTarget,omitempty"`
		IsEdgeIPBindingEnabled bool      `json:"isEdgeIPBindingEnabled,omitempty"`
		MapAlias               string    `json:"mapAlias,omitempty"`
		ProductId              string    `json:"productId,omitempty"`
		SerialNumber           int       `json:"serialNumber,omitempty"`
		UseCases               []UseCase `json:"useCases,omitempty"`
	}

	// ChinaCDN represents China CDN settings of EdgeHostname
	ChinaCDN struct {
		IsChinaCDN        bool   `json:"isChinaCdn,omitempty"`
		CustomChinaCDNMap string `json:"customChinaCdnMap,omitempty"`
	}

	// UseCase represents useCase attribute in EdgeHostname
	UseCase struct {
		Type    string `json:"type,omitempty"`
		Option  string `json:"option"`
		UseCase string `json:"useCase"`
	}

	// GetEdgeHostnameResponse represents edge hostname
	GetEdgeHostnameResponse struct {
		EdgeHostnameID         int       `json:"edgeHostnameId"`
		RecordName             string    `json:"recordName"`
		DNSZone                string    `json:"dnsZone"`
		SecurityType           string    `json:"securityType"`
		UseDefaultTTL          bool      `json:"useDefaultTtl"`
		UseDefaultMap          bool      `json:"useDefaultMap"`
		IPVersionBehavior      string    `json:"ipVersionBehavior"`
		ProductID              string    `json:"productId"`
		TTL                    int       `json:"ttl"`
		Map                    string    `json:"map,omitempty"`
		SlotNumber             int       `json:"slotNumber,omitempty"`
		Comments               string    `json:"comments"`
		SerialNumber           int       `json:"serialNumber,omitempty"`
		CustomTarget           string    `json:"customTarget,omitempty"`
		ChinaCdn               ChinaCDN  `json:"chinaCdn,omitempty"`
		IsEdgeIPBindingEnabled bool      `json:"isEdgeIPBindingEnabled,omitempty"`
		MapAlias               string    `json:"mapAlias"`
		UseCases               []UseCase `json:"useCases"`
	}

	// GetCertificateRequest is used to get certificate associated with edge hostname
	GetCertificateRequest struct {
		DNSZone    string
		RecordName string
	}

	// GetCertificateResponse represents edge hostname certificate
	GetCertificateResponse struct {
		AvailableDomains []string  `json:"availableDomains"`
		CertificateID    string    `json:"certificateId"`
		CertificateType  string    `json:"certificateType"`
		CommonName       string    `json:"commonName"`
		ExpirationDate   time.Time `json:"expirationDate"`
		SerialNumber     string    `json:"serialNumber"`
		SlotNumber       int       `json:"slotNumber"`
		Status           string    `json:"status"`
		ValidationType   string    `json:"validationType"`
	}
)

// Validate validates DeleteEdgeHostnameRequest
func (r DeleteEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	})
}

// Validate validates DeleteEdgeHostnameRequest
func (r UpdateEdgeHostnameRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
		"Body":       validation.Validate(r.Body),
	})
}

// Validate validates UpdateEdgeHostnameRequestBody
func (b UpdateEdgeHostnameRequestBody) Validate() error {
	return validation.Errors{
		"Path":  validation.Validate(b.Path, validation.Required, validation.In("/ttl", "/ipVersionBehavior").Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '/ttl' or '/ipVersionBehavior'", b.Path))),
		"Op":    validation.Validate(b.Op, validation.Required, validation.In("replace").Error(fmt.Sprintf("value '%s' is invalid. Must use 'replace'", b.Op))),
		"Value": validation.Validate(b.Value, validation.Required),
	}.Filter()
}

// Validate validates GetCertificateRequest
func (r GetCertificateRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DNSZone":    validation.Validate(r.DNSZone, validation.Required),
		"RecordName": validation.Validate(r.RecordName, validation.Required),
	})
}

var (
	// ErrDeleteEdgeHostname represents error when deleting edge hostname fails
	ErrDeleteEdgeHostname = errors.New("delete edge hostname")
	// ErrGetEdgeHostname represents error when getting edge hostname fails
	ErrGetEdgeHostname = errors.New("get edge hostname")
	// ErrUpdateEdgeHostname represents error when updating edge hostname fails
	ErrUpdateEdgeHostname = errors.New("update edge hostname")
	// ErrGetCertificate represents error when getting edge hostname certificate fails
	ErrGetCertificate = errors.New("get edge hostname certificate")
	// ErrNotFound represents error when getting edge hostname fails
	ErrNotFound = errors.New("not found")
)

func (h *hapi) DeleteEdgeHostname(ctx context.Context, params DeleteEdgeHostnameRequest) (*DeleteEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("DeleteEdgeHostname")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteEdgeHostname, ErrStructValidation, err)
	}

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
	defer session.CloseResponseBody(resp)

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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeHostname, h.Error(resp))
	}

	return &rval, nil
}

func (h *hapi) UpdateEdgeHostname(ctx context.Context, request UpdateEdgeHostnameRequest) (*UpdateEdgeHostnameResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("UpdateEdgeHostname")

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrUpdateEdgeHostname, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", request.DNSZone, request.RecordName)

	body, err := buildBody(request.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if len(request.StatusUpdateEmail) > 0 {
		emails := strings.Join(request.StatusUpdateEmail, ",")
		q.Add("statusUpdateEmail", emails)
	}
	if request.Comments != "" {
		q.Add("comments", request.Comments)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json-patch+json")

	var result UpdateEdgeHostnameResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateEdgeHostname, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%w: %s", ErrUpdateEdgeHostname, h.Error(resp))
	}

	return &result, nil
}

func (h *hapi) GetCertificate(ctx context.Context, params GetCertificateRequest) (*GetCertificateResponse, error) {
	logger := h.Log(ctx)
	logger.Debug("GetCertificate")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCertificate, ErrStructValidation, err)
	}

	uri := fmt.Sprintf(
		"/hapi/v1/dns-zones/%s/edge-hostnames/%s/certificate",
		params.DNSZone,
		params.RecordName,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCertificate, err)
	}

	var result GetCertificateResponse

	resp, err := h.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCertificate, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %s: %w", ErrGetCertificate, ErrNotFound, h.Error(resp))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCertificate, h.Error(resp))
	}

	return &result, nil
}

func buildBody(body []UpdateEdgeHostnameRequestBody) (io.Reader, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(reqBody), nil
}
