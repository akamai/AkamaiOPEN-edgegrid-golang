package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// LoadBalancerVersions is a cloudlets LoadBalancer version API interface
	LoadBalancerVersions interface {
		// CreateLoadBalancerVersion creates load balancer version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#postloadbalancingversions
		CreateLoadBalancerVersion(context.Context, CreateLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// GetLoadBalancerVersion gets specific load balancer version by originID and version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getloadbalancingconfigversion
		GetLoadBalancerVersion(context.Context, GetLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// UpdateLoadBalancerVersion updates specific load balancer version by originID and version
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#putloadbalancingconfigversion
		UpdateLoadBalancerVersion(context.Context, UpdateLoadBalancerVersionRequest) (*LoadBalancerVersion, error)

		// ListLoadBalancerVersions lists all versions of Origin with type APPLICATION_LOAD_BALANCER
		//
		// See: https://developer.akamai.com/api/web_performance/cloudlets/v2.html#getloadbalancingversions
		ListLoadBalancerVersions(context.Context, ListLoadBalancerVersionsRequest) ([]LoadBalancerVersion, error)
	}

	// DataCenter represents the dataCenter field of load balancer version
	DataCenter struct {
		City                          string   `json:"city,omitempty"`
		CloudServerHostHeaderOverride bool     `json:"cloudServerHostHeaderOverride,omitempty"`
		CloudService                  bool     `json:"cloudService"`
		Continent                     string   `json:"continent"`
		Country                       string   `json:"country"`
		Hostname                      string   `json:"hostname,omitempty"`
		Latitude                      *float64 `json:"latitude"`
		LivenessHosts                 []string `json:"livenessHosts,omitempty"`
		Longitude                     *float64 `json:"longitude"`
		OriginID                      string   `json:"originId"`
		Percent                       *float64 `json:"percent"`
		StateOrProvince               *string  `json:"stateOrProvince,omitempty"`
	}

	// LivenessSettings represents the livenessSettings field of load balancer version
	LivenessSettings struct {
		HostHeader                  string            `json:"hostHeader,omitempty"`
		AdditionalHeaders           map[string]string `json:"additionalHeaders,omitempty"`
		Interval                    int               `json:"interval,omitempty"`
		Path                        string            `json:"path,omitempty"`
		PeerCertificateVerification bool              `json:"peerCertificateVerification,omitempty"`
		Port                        int               `json:"port"`
		Protocol                    string            `json:"protocol"`
		RequestString               string            `json:"requestString,omitempty"`
		ResponseString              string            `json:"responseString,omitempty"`
		Status3xxFailure            bool              `json:"status3xxFailure,omitempty"`
		Status4xxFailure            bool              `json:"status4xxFailure,omitempty"`
		Status5xxFailure            bool              `json:"status5xxFailure,omitempty"`
		Timeout                     float64           `json:"timeout,omitempty"`
	}

	// BalancingType is a type for BalancingType field
	BalancingType string

	// LoadBalancerVersion describes the body of the create and update load balancer version request
	LoadBalancerVersion struct {
		BalancingType    BalancingType     `json:"balancingType,omitempty"`
		CreatedBy        string            `json:"createdBy,omitempty"`
		CreatedDate      string            `json:"createdDate,omitempty"`
		DataCenters      []DataCenter      `json:"dataCenters,omitempty"`
		Deleted          bool              `json:"deleted"`
		Description      string            `json:"description,omitempty"`
		Immutable        bool              `json:"immutable"`
		LastModifiedBy   string            `json:"lastModifiedBy,omitempty"`
		LastModifiedDate string            `json:"lastModifiedDate,omitempty"`
		LivenessSettings *LivenessSettings `json:"livenessSettings,omitempty"`
		OriginID         string            `json:"originID,omitempty"`
		Version          int64             `json:"version,omitempty"`
		Warnings         []Warning         `json:"warnings,omitempty"`
	}

	// CreateLoadBalancerVersionRequest describes the parameters needed to create load balancer version
	CreateLoadBalancerVersionRequest struct {
		OriginID            string
		LoadBalancerVersion LoadBalancerVersion
	}

	// GetLoadBalancerVersionRequest describes the parameters needed to get load balancer version
	GetLoadBalancerVersionRequest struct {
		OriginID       string
		Version        int64
		ShouldValidate bool
	}

	// UpdateLoadBalancerVersionRequest describes the parameters needed to update load balancer version
	UpdateLoadBalancerVersionRequest struct {
		OriginID            string
		ShouldValidate      bool
		Version             int64
		LoadBalancerVersion LoadBalancerVersion
	}

	// ListLoadBalancerVersionsRequest describes the parameters needed to list load balancer versions
	ListLoadBalancerVersionsRequest struct {
		OriginID string
	}
)

const (
	// BalancingTypeWeighted represents weighted balancing type for load balancer version
	BalancingTypeWeighted BalancingType = "WEIGHTED"
	// BalancingTypePerformance represents performance balancing type for load balancer version
	BalancingTypePerformance BalancingType = "PERFORMANCE"
)

var (
	// ErrCreateLoadBalancerVersion is returned when CreateLoadBalancerVersion fails
	ErrCreateLoadBalancerVersion = errors.New("create origin version")
	// ErrGetLoadBalancerVersion is returned when GetLoadBalancerVersion fails
	ErrGetLoadBalancerVersion = errors.New("get origin version")
	// ErrUpdateLoadBalancerVersion is returned when UpdateLoadBalancerVersion fails
	ErrUpdateLoadBalancerVersion = errors.New("update origin version")
	// ErrListLoadBalancerVersions is returned when ListLoadBalancerVersions fails
	ErrListLoadBalancerVersions = errors.New("list origin versions")
)

// Validate validates DataCenter
func (v DataCenter) Validate() error {
	return validation.Errors{
		"Continent": validation.Validate(v.Continent, validation.Required, validation.In("AF", "AS", "EU", "NA", "OC", "OT", "SA").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'AF', 'AS', 'EU', 'NA', 'OC', 'OT' or 'SA'", (&v).Continent))),
		"Country":   validation.Validate(v.Country, validation.Required, validation.Length(2, 2)),
		"Hostname":  validation.Validate(v.Hostname, validation.Length(0, 256)),
		"Latitude":  validation.Validate(v.Latitude, validation.NotNil, validation.Min(-180.0), validation.Max(180.0)),
		"Longitude": validation.Validate(v.Longitude, validation.NotNil, validation.Min(-180.0), validation.Max(180.0)),
		"OriginID":  validation.Validate(v.OriginID, validation.Required, validation.Length(1, 128)),
		"Percent":   validation.Validate(v.Percent, validation.NotNil, validation.Min(0.0), validation.Max(100.0)),
	}.Filter()
}

// generateHostHeaderRules generates case insensitive validation rules for host header
// its required because schema requires that host header value contains at least 1 char
// but it doesnt put such requirement on other headers, so these two cases need to be considered separately
func generateHostHeaderRules(headers map[string]string) []*validation.KeyRules {
	var hostRules []*validation.KeyRules

	for k := range headers {
		if strings.ToLower(k) == "host" {
			hostRules = append(hostRules, validation.Key(k, validation.Length(1, 256)))
		}
	}
	return hostRules
}

// Validate validates LivenessSettings
func (v LivenessSettings) Validate() error {
	return validation.Errors{
		"HostHeader":        validation.Validate(v.HostHeader, validation.Length(1, 256)),
		"AdditionalHeaders": validation.Validate(v.AdditionalHeaders, validation.Map(generateHostHeaderRules(v.AdditionalHeaders)...).AllowExtraKeys()),
		"Interval":          validation.Validate(v.Interval, validation.Min(10), validation.Max(3600)),
		"Path": validation.Validate(v.Path,
			validation.When(v.Protocol == "HTTP" || v.Protocol == "HTTPS", validation.Required, validation.Length(1, 256)),
		),
		"Port": validation.Validate(v.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		"Protocol": validation.Validate(v.Protocol, validation.Required, validation.In("HTTP", "HTTPS", "TCP", "TCPS").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'HTTP', 'HTTPS', 'TCP' or 'TCPS'", (&v).Protocol))),
		"RequestString": validation.Validate(v.RequestString,
			validation.When(v.Protocol == "TCP" || v.Protocol == "TCPS", validation.Required),
		),
		"ResponseString": validation.Validate(v.ResponseString,
			validation.When(v.Protocol == "TCP" || v.Protocol == "TCPS", validation.Required),
		),
		"Timeout": validation.Validate(v.Timeout, validation.Min(0.001), validation.Max(60.0)),
	}.Filter()
}

// Validate validates Warning
func (v Warning) Validate() error {
	return validation.Errors{
		"Detail":      validation.Validate(v.Detail, validation.Required),
		"JSONPointer": validation.Validate(v.JSONPointer, validation.Length(0, 128)),
		"Title":       validation.Validate(v.Title, validation.Required),
		"Type":        validation.Validate(v.Type, validation.Required),
	}
}

// Validate validates LoadBalancerVersion
func (v LoadBalancerVersion) Validate() error {
	return validation.Errors{
		"BalancingType": validation.Validate(v.BalancingType, validation.In(BalancingTypeWeighted, BalancingTypePerformance).Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'WEIGHTED', 'PERFORMANCE' or '' (empty)", (&v).BalancingType))),
		"CreatedDate":      validation.Validate(v.CreatedDate, validation.Date(time.RFC3339)),
		"DataCenters":      validation.Validate(v.DataCenters, validation.Length(1, 199)),
		"LastModifiedDate": validation.Validate(v.LastModifiedDate, validation.Date(time.RFC3339)),
		"LivenessSettings": validation.Validate(v.LivenessSettings),
		"OriginID":         validation.Validate(v.OriginID, validation.Length(2, 62)),
		"Version":          validation.Validate(v.Version, validation.Min(0)),
		"Warnings":         validation.Validate(v.Warnings),
	}.Filter()
}

// Validate validates CreateLoadBalancerVersionRequest
func (v CreateLoadBalancerVersionRequest) Validate() error {
	errs := validation.Errors{
		"OriginID":            validation.Validate(v.OriginID, validation.Length(2, 62)),
		"LoadBalancerVersion": validation.Validate(v.LoadBalancerVersion),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates GetLoadBalancerVersionRequest
func (v GetLoadBalancerVersionRequest) Validate() error {
	errs := validation.Errors{
		"OriginID": validation.Validate(v.OriginID, validation.Length(2, 62)),
		"Version":  validation.Validate(v.Version, validation.Min(0)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpdateLoadBalancerVersionRequest
func (v UpdateLoadBalancerVersionRequest) Validate() error {
	errs := validation.Errors{
		"OriginID":            validation.Validate(v.OriginID, validation.Length(2, 62)),
		"Version":             validation.Validate(v.Version, validation.Min(0)),
		"LoadBalancerVersion": validation.Validate(v.LoadBalancerVersion),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates ListLoadBalancerVersionsRequest
func (v ListLoadBalancerVersionsRequest) Validate() error {
	errs := validation.Errors{
		"OriginID": validation.Validate(v.OriginID, validation.Required, validation.Length(2, 62)),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

func (c *cloudlets) CreateLoadBalancerVersion(ctx context.Context, params CreateLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrCreateLoadBalancerVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/versions", params.OriginID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateLoadBalancerVersion, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateLoadBalancerVersion, err)
	}

	var result LoadBalancerVersion
	resp, err := c.Exec(req, &result, params.LoadBalancerVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateLoadBalancerVersion, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) GetLoadBalancerVersion(ctx context.Context, params GetLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("GetLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetLoadBalancerVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/versions/%d", params.OriginID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetLoadBalancerVersion, err)
	}

	if params.ShouldValidate {
		q := uri.Query()
		q.Add("validate", "true")
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetLoadBalancerVersion, err)
	}

	var result LoadBalancerVersion
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetLoadBalancerVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) UpdateLoadBalancerVersion(ctx context.Context, params UpdateLoadBalancerVersionRequest) (*LoadBalancerVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateLoadBalancerVersion")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpdateLoadBalancerVersion, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/versions/%d", params.OriginID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateLoadBalancerVersion, err)
	}

	if params.ShouldValidate {
		q := uri.Query()
		q.Add("validate", "true")
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateLoadBalancerVersion, err)
	}

	var result LoadBalancerVersion
	resp, err := c.Exec(req, &result, params.LoadBalancerVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateLoadBalancerVersion, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateLoadBalancerVersion, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudlets) ListLoadBalancerVersions(ctx context.Context, params ListLoadBalancerVersionsRequest) ([]LoadBalancerVersion, error) {
	logger := c.Log(ctx)
	logger.Debug("ListLoadBalancerVersions")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListLoadBalancerVersions, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloudlets/api/v2/origins/%s/versions?includeModel=true", params.OriginID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListLoadBalancerVersions, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListLoadBalancerVersions, err)
	}

	var result []LoadBalancerVersion
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListLoadBalancerVersions, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListLoadBalancerVersions, c.Error(resp))
	}

	return result, nil
}
