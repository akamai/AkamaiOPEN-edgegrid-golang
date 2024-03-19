package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Configurations is a CloudWrapper configurations API interface
	Configurations interface {
		// GetConfiguration gets a specific Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-configuration
		GetConfiguration(context.Context, GetConfigurationRequest) (*Configuration, error)
		// ListConfigurations lists all Cloud Wrapper configurations on your contract
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/get-configurations
		ListConfigurations(context.Context) (*ListConfigurationsResponse, error)
		// CreateConfiguration creates a Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration
		CreateConfiguration(context.Context, CreateConfigurationRequest) (*Configuration, error)
		// UpdateConfiguration updates a saved or inactive configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/put-configuration
		UpdateConfiguration(context.Context, UpdateConfigurationRequest) (*Configuration, error)
		// DeleteConfiguration deletes configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/delete-configuration
		DeleteConfiguration(context.Context, DeleteConfigurationRequest) error
		// ActivateConfiguration activates a Cloud Wrapper configuration
		//
		// See: https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration-activations
		ActivateConfiguration(context.Context, ActivateConfigurationRequest) error
	}

	// GetConfigurationRequest holds parameters for GetConfiguration
	GetConfigurationRequest struct {
		ConfigID int64
	}

	// CreateConfigurationRequest holds parameters for CreateConfiguration
	CreateConfigurationRequest struct {
		Activate bool
		Body     CreateConfigurationBody
	}

	// CreateConfigurationBody holds request body parameters for CreateConfiguration
	CreateConfigurationBody struct {
		CapacityAlertsThreshold *int                `json:"capacityAlertsThreshold,omitempty"`
		Comments                string              `json:"comments"`
		ContractID              string              `json:"contractId"`
		Locations               []ConfigLocationReq `json:"locations"`
		MultiCDNSettings        *MultiCDNSettings   `json:"multiCdnSettings,omitempty"`
		ConfigName              string              `json:"configName"`
		NotificationEmails      []string            `json:"notificationEmails,omitempty"`
		PropertyIDs             []string            `json:"propertyIds"`
		RetainIdleObjects       bool                `json:"retainIdleObjects,omitempty"`
	}

	// UpdateConfigurationRequest holds parameters for UpdateConfiguration
	UpdateConfigurationRequest struct {
		ConfigID int64
		Activate bool
		Body     UpdateConfigurationBody
	}

	// UpdateConfigurationBody holds request body parameters for UpdateConfiguration
	UpdateConfigurationBody struct {
		CapacityAlertsThreshold *int                `json:"capacityAlertsThreshold,omitempty"`
		Comments                string              `json:"comments"`
		Locations               []ConfigLocationReq `json:"locations"`
		MultiCDNSettings        *MultiCDNSettings   `json:"multiCdnSettings,omitempty"`
		NotificationEmails      []string            `json:"notificationEmails,omitempty"`
		PropertyIDs             []string            `json:"propertyIds"`
		RetainIdleObjects       bool                `json:"retainIdleObjects,omitempty"`
	}

	// DeleteConfigurationRequest holds parameters for DeleteConfiguration
	DeleteConfigurationRequest struct {
		ConfigID int64
	}

	// ActivateConfigurationRequest holds parameters for ActivateConfiguration
	ActivateConfigurationRequest struct {
		ConfigurationIDs []int `json:"configurationIds"`
	}

	// Configuration represents CloudWrapper configuration
	Configuration struct {
		CapacityAlertsThreshold *int                 `json:"capacityAlertsThreshold"`
		Comments                string               `json:"comments"`
		ContractID              string               `json:"contractId"`
		ConfigID                int64                `json:"configId"`
		Locations               []ConfigLocationResp `json:"locations"`
		MultiCDNSettings        *MultiCDNSettings    `json:"multiCdnSettings"`
		Status                  StatusType           `json:"status"`
		ConfigName              string               `json:"configName"`
		LastUpdatedBy           string               `json:"lastUpdatedBy"`
		LastUpdatedDate         string               `json:"lastUpdatedDate"`
		LastActivatedBy         *string              `json:"lastActivatedBy"`
		LastActivatedDate       *string              `json:"lastActivatedDate"`
		NotificationEmails      []string             `json:"notificationEmails"`
		PropertyIDs             []string             `json:"propertyIds"`
		RetainIdleObjects       bool                 `json:"retainIdleObjects"`
	}

	// ListConfigurationsResponse contains response from ListConfigurations
	ListConfigurationsResponse struct {
		Configurations []Configuration `json:"configurations"`
	}

	// ConfigLocationReq represents location to be configured for the configuration
	ConfigLocationReq struct {
		Comments      string   `json:"comments"`
		TrafficTypeID int      `json:"trafficTypeId"`
		Capacity      Capacity `json:"capacity"`
	}

	// ConfigLocationResp represents location to be configured for the configuration
	ConfigLocationResp struct {
		Comments      string   `json:"comments"`
		TrafficTypeID int      `json:"trafficTypeId"`
		Capacity      Capacity `json:"capacity"`
		MapName       string   `json:"mapName"`
	}

	// MultiCDNSettings represents details about Multi CDN Settings
	MultiCDNSettings struct {
		BOCC             *BOCC        `json:"bocc"`
		CDNs             []CDN        `json:"cdns"`
		DataStreams      *DataStreams `json:"dataStreams"`
		EnableSoftAlerts bool         `json:"enableSoftAlerts,omitempty"`
		Origins          []Origin     `json:"origins"`
	}

	// BOCC represents diagnostic data beacon details
	BOCC struct {
		ConditionalSamplingFrequency SamplingFrequency `json:"conditionalSamplingFrequency,omitempty"`
		Enabled                      bool              `json:"enabled"`
		ForwardType                  ForwardType       `json:"forwardType,omitempty"`
		RequestType                  RequestType       `json:"requestType,omitempty"`
		SamplingFrequency            SamplingFrequency `json:"samplingFrequency,omitempty"`
	}

	// CDN represents a CDN added for the configuration
	CDN struct {
		CDNAuthKeys []CDNAuthKey `json:"cdnAuthKeys,omitempty"`
		CDNCode     string       `json:"cdnCode"`
		Enabled     bool         `json:"enabled"`
		HTTPSOnly   bool         `json:"httpsOnly,omitempty"`
		IPACLCIDRs  []string     `json:"ipAclCidrs,omitempty"`
	}

	// CDNAuthKey represents auth key configured for the CDN
	CDNAuthKey struct {
		AuthKeyName string `json:"authKeyName"`
		ExpiryDate  string `json:"expiryDate,omitempty"`
		HeaderName  string `json:"headerName,omitempty"`
		Secret      string `json:"secret,omitempty"`
	}

	// DataStreams represents data streams details
	DataStreams struct {
		DataStreamIDs []int64 `json:"dataStreamIds,omitempty"`
		Enabled       bool    `json:"enabled"`
		SamplingRate  *int    `json:"samplingRate,omitempty"`
	}

	// Origin represents origin corresponding to the properties selected in the configuration
	Origin struct {
		Hostname   string `json:"hostname"`
		OriginID   string `json:"originId"`
		PropertyID int    `json:"propertyId"`
	}

	// SamplingFrequency is a type of sampling frequency. Either 'ZERO' or 'ONE_TENTH'
	SamplingFrequency string

	// ForwardType is a type of forward
	ForwardType string

	// RequestType is a type of request
	RequestType string

	// StatusType is a type of status
	StatusType string
)

const (
	// SamplingFrequencyZero represents SamplingFrequency value of 'ZERO'
	SamplingFrequencyZero SamplingFrequency = "ZERO"
	// SamplingFrequencyOneTenth represents SamplingFrequency value of 'ONE_TENTH'
	SamplingFrequencyOneTenth SamplingFrequency = "ONE_TENTH"
	// ForwardTypeOriginOnly represents ForwardType value of 'ORIGIN_ONLY'
	ForwardTypeOriginOnly ForwardType = "ORIGIN_ONLY"
	// ForwardTypeMidgressOnly represents ForwardType value of 'MIDGRESS_ONLY'
	ForwardTypeMidgressOnly ForwardType = "MIDGRESS_ONLY"
	// ForwardTypeOriginAndMidgress represents ForwardType value of 'ORIGIN_AND_MIDGRESS'
	ForwardTypeOriginAndMidgress ForwardType = "ORIGIN_AND_MIDGRESS"
	// RequestTypeEdgeOnly represents RequestType value of 'EDGE_ONLY'
	RequestTypeEdgeOnly RequestType = "EDGE_ONLY"
	// RequestTypeEdgeAndMidgress represents RequestType value of 'EDGE_AND_MIDGRESS'
	RequestTypeEdgeAndMidgress RequestType = "EDGE_AND_MIDGRESS"

	// StatusActive represents Status value of 'ACTIVE'
	StatusActive StatusType = "ACTIVE"
	// StatusSaved represents Status value of 'SAVED'
	StatusSaved StatusType = "SAVED"
	// StatusInProgress represents Status value of 'IN_PROGRESS'
	StatusInProgress StatusType = "IN_PROGRESS"
	// StatusDeleteInProgress represents Status value of 'DELETE_IN_PROGRESS'
	StatusDeleteInProgress StatusType = "DELETE_IN_PROGRESS"
	// StatusFailed represents Status value of 'FAILED'
	StatusFailed StatusType = "FAILED"
)

// Validate validates GetConfigurationRequest
func (r GetConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(r.ConfigID, validation.Required),
	})
}

// Validate validates CreateConfigurationRequest
func (r CreateConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Body": validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates CreateConfigurationBody
func (b CreateConfigurationBody) Validate() error {
	return validation.Errors{
		"Comments":                validation.Validate(b.Comments, validation.Required),
		"Locations":               validation.Validate(b.Locations, validation.Required),
		"ConfigName":              validation.Validate(b.ConfigName, validation.Required),
		"ContractID":              validation.Validate(b.ContractID, validation.Required),
		"PropertyIDs":             validation.Validate(b.PropertyIDs, validation.Required),
		"MultiCDNSettings":        validation.Validate(b.MultiCDNSettings),
		"CapacityAlertsThreshold": validation.Validate(b.CapacityAlertsThreshold, validation.Min(50), validation.Max(100).Error(fmt.Sprintf("value '%d' is invalid. Must be between 50 and 100", b.CapacityAlertsThreshold))),
	}.Filter()
}

// Validate validates UpdateConfigurationRequest
func (r UpdateConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(r.ConfigID, validation.Required),
		"Body":     validation.Validate(r.Body, validation.Required),
	})
}

// Validate validates UpdateConfigurationBody
func (b UpdateConfigurationBody) Validate() error {
	return validation.Errors{
		"Comments":                validation.Validate(b.Comments, validation.Required),
		"Locations":               validation.Validate(b.Locations, validation.Required),
		"PropertyIDs":             validation.Validate(b.PropertyIDs, validation.Required),
		"MultiCDNSettings":        validation.Validate(b.MultiCDNSettings),
		"CapacityAlertsThreshold": validation.Validate(b.CapacityAlertsThreshold, validation.Min(50), validation.Max(100).Error(fmt.Sprintf("value '%d' is invalid. Must be between 50 and 100", b.CapacityAlertsThreshold))),
	}.Filter()
}

// Validate validates DeleteConfigurationRequest
func (r DeleteConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigID": validation.Validate(r.ConfigID, validation.Required),
	})
}

// Validate validates ActivateConfigurationRequest
func (r ActivateConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"ConfigurationIDs": validation.Validate(r.ConfigurationIDs, validation.Required),
	})
}

// Validate validates ConfigurationLocation
func (c ConfigLocationReq) Validate() error {
	return validation.Errors{
		"Comments":      validation.Validate(c.Comments, validation.Required),
		"Capacity":      validation.Validate(c.Capacity, validation.Required),
		"TrafficTypeID": validation.Validate(c.TrafficTypeID, validation.Required),
	}.Filter()
}

// Validate validates Capacity
func (c Capacity) Validate() error {
	return validation.Errors{
		"Unit":  validation.Validate(c.Unit, validation.Required, validation.In(UnitGB, UnitTB).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", c.Unit, UnitGB, UnitTB))),
		"Value": validation.Validate(c.Value, validation.Required, validation.Min(1), validation.Max(int64(10000000000))),
	}.Filter()
}

// Validate validates MultiCDNSettings
func (m MultiCDNSettings) Validate() error {
	return validation.Errors{
		"BOCC":        validation.Validate(m.BOCC, validation.Required),
		"CDNs":        validation.Validate(m.CDNs, validation.By(validateCDNs)),
		"DataStreams": validation.Validate(m.DataStreams, validation.Required),
		"Origins":     validation.Validate(m.Origins, validation.Required),
	}.Filter()
}

// Validate validates BOCC
func (b BOCC) Validate() error {
	return validation.Errors{
		"Enabled":                      validation.Validate(b.Enabled, validation.NotNil),
		"ConditionalSamplingFrequency": validation.Validate(b.ConditionalSamplingFrequency, validation.Required.When(b.Enabled), validation.In(SamplingFrequencyZero, SamplingFrequencyOneTenth).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", b.ConditionalSamplingFrequency, SamplingFrequencyZero, SamplingFrequencyOneTenth))),
		"ForwardType":                  validation.Validate(b.ForwardType, validation.Required.When(b.Enabled), validation.In(ForwardTypeOriginOnly, ForwardTypeMidgressOnly, ForwardTypeOriginAndMidgress).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s'", b.ForwardType, ForwardTypeOriginOnly, ForwardTypeMidgressOnly, ForwardTypeOriginAndMidgress))),
		"RequestType":                  validation.Validate(b.RequestType, validation.Required.When(b.Enabled), validation.In(RequestTypeEdgeOnly, RequestTypeEdgeAndMidgress).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", b.RequestType, RequestTypeEdgeOnly, RequestTypeEdgeAndMidgress))),
		"SamplingFrequency":            validation.Validate(b.SamplingFrequency, validation.Required.When(b.Enabled), validation.In(SamplingFrequencyZero, SamplingFrequencyOneTenth).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s'", b.RequestType, SamplingFrequencyZero, SamplingFrequencyOneTenth))),
	}.Filter()
}

// Validate validates CDN
func (c CDN) Validate() error {
	return validation.Errors{
		"CDNAuthKeys": validation.Validate(c.CDNAuthKeys),
		"Enabled":     validation.Validate(c.Enabled, validation.NotNil),
		"CDNCode":     validation.Validate(c.CDNCode, validation.Required),
	}.Filter()
}

// Validate validates CDNAuthKey
func (c CDNAuthKey) Validate() error {
	return validation.Errors{
		"AuthKeyName": validation.Validate(c.AuthKeyName, validation.Required),
		"Secret":      validation.Validate(c.Secret, validation.Length(24, 24)),
	}.Filter()
}

// Validate validates DataStreams
func (d DataStreams) Validate() error {
	return validation.Errors{
		"Enabled":      validation.Validate(d.Enabled, validation.NotNil),
		"SamplingRate": validation.Validate(d.SamplingRate, validation.When(d.SamplingRate != nil, validation.Min(1), validation.Max(100).Error(fmt.Sprintf("value '%d' is invalid. Must be between 1 and 100", d.SamplingRate)))),
	}.Filter()
}

// Validate validates Origin
func (o Origin) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(o.PropertyID, validation.Required),
		"Hostname":   validation.Validate(o.Hostname, validation.Required),
		"OriginID":   validation.Validate(o.OriginID, validation.Required),
	}.Filter()
}

// validateCDNs validates CDNs by checking if at least one is enabled and one of authKeys or IP ACLs is specified
func validateCDNs(value interface{}) error {
	v, ok := value.([]CDN)
	if !ok {
		return fmt.Errorf("type %T is invalid. Must be []CDN", value)
	}
	if v == nil {
		return fmt.Errorf("cannot be blank")
	}
	var isEnabled bool
	for _, cdn := range v {
		if cdn.Enabled {
			isEnabled = true
		}
		if cdn.CDNAuthKeys == nil && cdn.IPACLCIDRs == nil {
			return fmt.Errorf("at least one authentication method is required for CDN. Either IP ACL or header authentication must be enabled")
		}
	}
	if !isEnabled {
		return fmt.Errorf("at least one of CDNs must be enabled")
	}

	return nil
}

var (
	// ErrGetConfiguration is returned when GetConfiguration fails
	ErrGetConfiguration = errors.New("get configuration")
	// ErrListConfigurations is returned when ListConfigurations fails
	ErrListConfigurations = errors.New("list configurations")
	// ErrCreateConfiguration is returned when CreateConfiguration fails
	ErrCreateConfiguration = errors.New("create configuration")
	// ErrUpdateConfiguration is returned when UpdateConfiguration fails
	ErrUpdateConfiguration = errors.New("update configuration")
	// ErrDeleteConfiguration is returned when DeleteConfiguration fails
	ErrDeleteConfiguration = errors.New("delete configuration")
	// ErrActivateConfiguration is returned when ActivateConfiguration fails
	ErrActivateConfiguration = errors.New("activate configuration")
)

func (c *cloudwrapper) GetConfiguration(ctx context.Context, params GetConfigurationRequest) (*Configuration, error) {
	logger := c.Log(ctx)
	logger.Debug("GetConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetConfiguration, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloud-wrapper/v1/configurations/%d", params.ConfigID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetConfiguration, err)
	}

	var result Configuration
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetConfiguration, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetConfiguration, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) ListConfigurations(ctx context.Context) (*ListConfigurationsResponse, error) {
	logger := c.Log(ctx)
	logger.Debug("ListConfigurations")

	uri := "/cloud-wrapper/v1/configurations"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListConfigurations, err)
	}

	var result ListConfigurationsResponse
	resp, err := c.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListConfigurations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListConfigurations, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) CreateConfiguration(ctx context.Context, params CreateConfigurationRequest) (*Configuration, error) {
	logger := c.Log(ctx)
	logger.Debug("CreateConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateConfiguration, ErrStructValidation, err)
	}

	uri, err := url.Parse("/cloud-wrapper/v1/configurations")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCreateConfiguration, err)
	}

	q := uri.Query()
	q.Add("activate", strconv.FormatBool(params.Activate))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateConfiguration, err)
	}

	var result Configuration
	resp, err := c.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateConfiguration, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateConfiguration, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) UpdateConfiguration(ctx context.Context, params UpdateConfigurationRequest) (*Configuration, error) {
	logger := c.Log(ctx)
	logger.Debug("UpdateConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateConfiguration, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/cloud-wrapper/v1/configurations/%d", params.ConfigID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateConfiguration, err)
	}

	q := uri.Query()
	q.Add("activate", strconv.FormatBool(params.Activate))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateConfiguration, err)
	}

	var result Configuration
	resp, err := c.Exec(req, &result, params.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateConfiguration, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateConfiguration, c.Error(resp))
	}

	return &result, nil
}

func (c *cloudwrapper) DeleteConfiguration(ctx context.Context, params DeleteConfigurationRequest) error {
	logger := c.Log(ctx)
	logger.Debug("DeleteConfiguration")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteConfiguration, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/cloud-wrapper/v1/configurations/%d", params.ConfigID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteConfiguration, err)
	}

	resp, err := c.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteConfiguration, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%s: %w", ErrDeleteConfiguration, c.Error(resp))
	}

	return nil
}

func (c *cloudwrapper) ActivateConfiguration(ctx context.Context, params ActivateConfigurationRequest) error {
	logger := c.Log(ctx)
	logger.Debug("ActivateConfiguration")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrActivateConfiguration, ErrStructValidation, err)
	}

	uri := "/cloud-wrapper/v1/configurations/activate"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrActivateConfiguration, err)
	}

	resp, err := c.Exec(req, nil, params)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrActivateConfiguration, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrActivateConfiguration, c.Error(resp))
	}

	return nil
}
