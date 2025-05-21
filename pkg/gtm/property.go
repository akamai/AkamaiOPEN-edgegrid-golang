package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// TrafficTarget struct contains information about where to direct data center traffic
	TrafficTarget struct {
		DatacenterID int      `json:"datacenterId"`
		Enabled      bool     `json:"enabled"`
		Weight       float64  `json:"weight,omitempty"`
		Servers      []string `json:"servers,omitempty"`
		Name         string   `json:"name,omitempty"`
		HandoutCName string   `json:"handoutCName,omitempty"`
		Precedence   *int     `json:"precedence,omitempty"`
	}

	// HTTPHeader struct contains HTTP headers to send if the testObjectProtocol is http or https
	HTTPHeader struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// LivenessTest contains configuration of liveness tests to determine whether your servers respond to requests
	LivenessTest struct {
		Name                          string       `json:"name"`
		ErrorPenalty                  float64      `json:"errorPenalty,omitempty"`
		PeerCertificateVerification   bool         `json:"peerCertificateVerification"`
		TestInterval                  int          `json:"testInterval,omitempty"`
		TestObject                    string       `json:"testObject,omitempty"`
		Links                         []Link       `json:"links,omitempty"`
		RequestString                 string       `json:"requestString,omitempty"`
		ResponseString                string       `json:"responseString,omitempty"`
		HTTPError3xx                  bool         `json:"httpError3xx"`
		HTTPError4xx                  bool         `json:"httpError4xx"`
		HTTPError5xx                  bool         `json:"httpError5xx"`
		HTTPMethod                    *string      `json:"httpMethod"`
		HTTPRequestBody               *string      `json:"httpRequestBody"`
		Disabled                      bool         `json:"disabled"`
		TestObjectProtocol            string       `json:"testObjectProtocol,omitempty"`
		TestObjectPassword            string       `json:"testObjectPassword,omitempty"`
		TestObjectPort                int          `json:"testObjectPort,omitempty"`
		SSLClientPrivateKey           string       `json:"sslClientPrivateKey,omitempty"`
		SSLClientCertificate          string       `json:"sslClientCertificate,omitempty"`
		Pre2023SecurityPosture        bool         `json:"pre2023SecurityPosture"`
		DisableNonstandardPortWarning bool         `json:"disableNonstandardPortWarning"`
		HTTPHeaders                   []HTTPHeader `json:"httpHeaders,omitempty"`
		TestObjectUsername            string       `json:"testObjectUsername,omitempty"`
		TestTimeout                   float32      `json:"testTimeout,omitempty"`
		TimeoutPenalty                float64      `json:"timeoutPenalty,omitempty"`
		AnswersRequired               bool         `json:"answersRequired"`
		ResourceType                  string       `json:"resourceType,omitempty"`
		RecursionRequested            bool         `json:"recursionRequested"`
		AlternateCACertificates       []string     `json:"alternateCACertificates"`
	}

	// StaticRRSet contains static recordset
	StaticRRSet struct {
		Type  string   `json:"type"`
		TTL   int      `json:"ttl"`
		Rdata []string `json:"rdata"`
	}

	// Property represents a GTM property
	Property struct {
		Name                      string          `json:"name"`
		Type                      string          `json:"type"`
		IPv6                      bool            `json:"ipv6"`
		ScoreAggregationType      string          `json:"scoreAggregationType"`
		StickinessBonusPercentage int             `json:"stickinessBonusPercentage,omitempty"`
		StickinessBonusConstant   int             `json:"stickinessBonusConstant,omitempty"`
		HealthThreshold           float64         `json:"healthThreshold,omitempty"`
		UseComputedTargets        bool            `json:"useComputedTargets"`
		BackupIP                  string          `json:"backupIp,omitempty"`
		BalanceByDownloadScore    bool            `json:"balanceByDownloadScore"`
		StaticTTL                 int             `json:"staticTTL,omitempty"`
		StaticRRSets              []StaticRRSet   `json:"staticRRSets,omitempty"`
		LastModified              string          `json:"lastModified"`
		UnreachableThreshold      float64         `json:"unreachableThreshold,omitempty"`
		MinLiveFraction           float64         `json:"minLiveFraction,omitempty"`
		HealthMultiplier          float64         `json:"healthMultiplier,omitempty"`
		DynamicTTL                int             `json:"dynamicTTL,omitempty"`
		MaxUnreachablePenalty     int             `json:"maxUnreachablePenalty,omitempty"`
		MapName                   string          `json:"mapName,omitempty"`
		HandoutLimit              int             `json:"handoutLimit"`
		HandoutMode               string          `json:"handoutMode"`
		FailoverDelay             int             `json:"failoverDelay,omitempty"`
		BackupCName               string          `json:"backupCName,omitempty"`
		FailbackDelay             int             `json:"failbackDelay,omitempty"`
		LoadImbalancePercentage   float64         `json:"loadImbalancePercentage,omitempty"`
		HealthMax                 float64         `json:"healthMax,omitempty"`
		GhostDemandReporting      bool            `json:"ghostDemandReporting"`
		Comments                  string          `json:"comments,omitempty"`
		CName                     string          `json:"cname,omitempty"`
		WeightedHashBitsForIPv4   int             `json:"weightedHashBitsForIPv4,omitempty"`
		WeightedHashBitsForIPv6   int             `json:"weightedHashBitsForIPv6,omitempty"`
		TrafficTargets            []TrafficTarget `json:"trafficTargets,omitempty"`
		Links                     []Link          `json:"links,omitempty"`
		LivenessTests             []LivenessTest  `json:"livenessTests,omitempty"`
	}

	// PropertyRequest contains request parameters
	PropertyRequest struct {
		Property   *Property
		DomainName string
	}

	// PropertyList contains a list of property items
	PropertyList struct {
		PropertyItems []Property `json:"items"`
	}
	// GetPropertyRequest contains request parameters for GetProperty
	GetPropertyRequest struct {
		DomainName   string
		PropertyName string
	}

	// GetPropertyResponse contains the response data from GetProperty operation
	GetPropertyResponse Property

	// ListPropertiesRequest contains request parameters for ListProperties
	ListPropertiesRequest struct {
		DomainName string
	}

	// CreatePropertyRequest contains request parameters for CreateProperty
	CreatePropertyRequest PropertyRequest

	// CreatePropertyResponse contains the response data from CreateProperty operation
	CreatePropertyResponse struct {
		Resource *Property       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// UpdatePropertyRequest contains request parameters for UpdatePropertyResponse
	UpdatePropertyRequest PropertyRequest

	// UpdatePropertyResponse contains the response data from UpdatePropertyResponse operation
	UpdatePropertyResponse struct {
		Resource *Property       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}

	// DeletePropertyRequest contains request parameters for DeleteProperty
	DeletePropertyRequest struct {
		DomainName   string
		PropertyName string
	}

	// DeletePropertyResponse contains the response data from DeleteProperty operation
	DeletePropertyResponse struct {
		Resource *Property       `json:"resource"`
		Status   *ResponseStatus `json:"status"`
	}
)

var (
	// ErrGetProperty is returned when GetProperty fails.
	ErrGetProperty = errors.New("get property")
	// ErrListProperties is returned when ListProperties fails.
	ErrListProperties = errors.New("list properties")
	// ErrCreateProperty is returned when CreateProperty fails.
	ErrCreateProperty = errors.New("create Property")
	// ErrUpdateProperty is returned when UpdateProperty fails
	ErrUpdateProperty = errors.New("update Property")
	// ErrDeleteProperty is returned when DeleteProperty fails
	ErrDeleteProperty = errors.New("delete Property")
)

// Validate validates GetPropertyRequest
func (r GetPropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
		"PropertyName": validation.Validate(r.PropertyName, validation.Required),
	})
}

// Validate validates ListPropertiesRequest
func (r ListPropertiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
	})
}

// Validate validates CreatePropertyRequest
func (r CreatePropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Property":   validation.Validate(r.Property, validation.Required),
	})
}

// Validate validates UpdatePropertyRequest
func (r UpdatePropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName": validation.Validate(r.DomainName, validation.Required),
		"Property":   validation.Validate(r.Property, validation.Required),
	})
}

// Validate validates DeletePropertyRequest
func (r DeletePropertyRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"DomainName":   validation.Validate(r.DomainName, validation.Required),
		"PropertyName": validation.Validate(r.PropertyName, validation.Required),
	})
}

// Validate validates Property
func (p *Property) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Name":                  validation.Validate(p.Name, validation.Required),
		"Type":                  validation.Validate(p.Type, validation.Required),
		"ScoreAggregationTypes": validation.Validate(p.ScoreAggregationType, validation.Required),
		"HandoutMode":           validation.Validate(p.HandoutMode, validation.Required),
		"TrafficTargets":        validation.Validate(p.TrafficTargets, validation.When(p.Type == "ranked-failover", validation.By(validateRankedFailoverTrafficTargets))),
	})
}

// validateRankedFailoverTrafficTargets validates traffic targets when property type is 'ranked-failover'
func validateRankedFailoverTrafficTargets(value interface{}) error {
	tt := value.([]TrafficTarget)
	if len(tt) == 0 {
		return fmt.Errorf("no traffic targets are enabled")
	}
	precedenceCounter := map[int]int{}
	minPrecedence := 256
	for _, t := range tt {
		if t.Precedence == nil {
			precedenceCounter[0]++
			minPrecedence = 0
		} else {
			if *t.Precedence > 255 || *t.Precedence < 0 {
				return fmt.Errorf("'Precedence' value has to be between 0 and 255")
			}
			precedenceCounter[*t.Precedence]++
			if *t.Precedence < minPrecedence {
				minPrecedence = *t.Precedence
			}
		}
	}
	if precedenceCounter[minPrecedence] > 1 {
		return fmt.Errorf("property cannot have multiple primary traffic targets (targets with lowest precedence)")
	}

	return nil
}

func (g *gtm) ListProperties(ctx context.Context, params ListPropertiesRequest) ([]Property, error) {
	logger := g.Log(ctx)
	logger.Debug("ListProperties")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListProperties, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties", params.DomainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListProperties request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result PropertyList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListProperties request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.PropertyItems, nil
}

func (g *gtm) GetProperty(ctx context.Context, params GetPropertyRequest) (*GetPropertyResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("GetProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProperty, ErrStructValidation, err)
	}

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", params.DomainName, params.PropertyName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetProperty request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result GetPropertyResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetProperty request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateProperty(ctx context.Context, params CreatePropertyRequest) (*CreatePropertyResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateProperty, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", params.DomainName, params.Property.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result CreatePropertyResponse
	resp, err := g.Exec(req, &result, params.Property)
	if err != nil {
		return nil, fmt.Errorf("property request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) UpdateProperty(ctx context.Context, params UpdatePropertyRequest) (*UpdatePropertyResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateProperty, ErrStructValidation, err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", params.DomainName, params.Property.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result UpdatePropertyResponse
	resp, err := g.Exec(req, &result, params.Property)
	if err != nil {
		return nil, fmt.Errorf("property request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteProperty(ctx context.Context, params DeletePropertyRequest) (*DeletePropertyResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteProperty")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteProperty, ErrStructValidation, err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", params.DomainName, params.PropertyName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DeletePropertyResponse
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteProperty request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}
