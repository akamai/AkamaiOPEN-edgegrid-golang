package gtm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegriderr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Properties contains operations available on a Property resource.
type Properties interface {
	// ListProperties retrieves all Properties for the provided domainName.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-properties
	ListProperties(context.Context, string) ([]*Property, error)
	// GetProperty retrieves a Property with the given domain and property names.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-property
	GetProperty(context.Context, string, string) (*Property, error)
	// CreateProperty creates property.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-property
	CreateProperty(context.Context, *Property, string) (*PropertyResponse, error)
	// DeleteProperty is a method applied to a property object resulting in removal.
	//
	// See: https://techdocs.akamai.com/gtm/reference/delete-property
	DeleteProperty(context.Context, *Property, string) (*ResponseStatus, error)
	// UpdateProperty is a method applied to a property object resulting in an update.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-property
	UpdateProperty(context.Context, *Property, string) (*ResponseStatus, error)
}

// TrafficTarget struct contains information about where to direct data center traffic
type TrafficTarget struct {
	DatacenterID int      `json:"datacenterId"`
	Enabled      bool     `json:"enabled"`
	Weight       float64  `json:"weight,omitempty"`
	Servers      []string `json:"servers,omitempty"`
	Name         string   `json:"name,omitempty"`
	HandoutCName string   `json:"handoutCName,omitempty"`
	Precedence   *int     `json:"precedence,omitempty"`
}

// HTTPHeader struct contains HTTP headers to send if the testObjectProtocol is http or https
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// LivenessTest contains configuration of liveness tests to determine whether your servers respond to requests
type LivenessTest struct {
	Name                          string        `json:"name"`
	ErrorPenalty                  float64       `json:"errorPenalty,omitempty"`
	PeerCertificateVerification   bool          `json:"peerCertificateVerification"`
	TestInterval                  int           `json:"testInterval,omitempty"`
	TestObject                    string        `json:"testObject,omitempty"`
	Links                         []*Link       `json:"links,omitempty"`
	RequestString                 string        `json:"requestString,omitempty"`
	ResponseString                string        `json:"responseString,omitempty"`
	HTTPError3xx                  bool          `json:"httpError3xx"`
	HTTPError4xx                  bool          `json:"httpError4xx"`
	HTTPError5xx                  bool          `json:"httpError5xx"`
	HTTPMethod                    *string       `json:"httpMethod"`
	HTTPRequestBody               *string       `json:"httpRequestBody"`
	Disabled                      bool          `json:"disabled"`
	TestObjectProtocol            string        `json:"testObjectProtocol,omitempty"`
	TestObjectPassword            string        `json:"testObjectPassword,omitempty"`
	TestObjectPort                int           `json:"testObjectPort,omitempty"`
	SSLClientPrivateKey           string        `json:"sslClientPrivateKey,omitempty"`
	SSLClientCertificate          string        `json:"sslClientCertificate,omitempty"`
	Pre2023SecurityPosture        bool          `json:"pre2023SecurityPosture"`
	DisableNonstandardPortWarning bool          `json:"disableNonstandardPortWarning"`
	HTTPHeaders                   []*HTTPHeader `json:"httpHeaders,omitempty"`
	TestObjectUsername            string        `json:"testObjectUsername,omitempty"`
	TestTimeout                   float32       `json:"testTimeout,omitempty"`
	TimeoutPenalty                float64       `json:"timeoutPenalty,omitempty"`
	AnswersRequired               bool          `json:"answersRequired"`
	ResourceType                  string        `json:"resourceType,omitempty"`
	RecursionRequested            bool          `json:"recursionRequested"`
	AlternateCACertificates       []string      `json:"alternateCACertificates"`
}

// StaticRRSet contains static recordset
type StaticRRSet struct {
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Rdata []string `json:"rdata"`
}

// Property represents a GTM property
type Property struct {
	Name                      string           `json:"name"`
	Type                      string           `json:"type"`
	IPv6                      bool             `json:"ipv6"`
	ScoreAggregationType      string           `json:"scoreAggregationType"`
	StickinessBonusPercentage int              `json:"stickinessBonusPercentage,omitempty"`
	StickinessBonusConstant   int              `json:"stickinessBonusConstant,omitempty"`
	HealthThreshold           float64          `json:"healthThreshold,omitempty"`
	UseComputedTargets        bool             `json:"useComputedTargets"`
	BackupIP                  string           `json:"backupIp,omitempty"`
	BalanceByDownloadScore    bool             `json:"balanceByDownloadScore"`
	StaticTTL                 int              `json:"staticTTL,omitempty"`
	StaticRRSets              []*StaticRRSet   `json:"staticRRSets,omitempty"`
	LastModified              string           `json:"lastModified"`
	UnreachableThreshold      float64          `json:"unreachableThreshold,omitempty"`
	MinLiveFraction           float64          `json:"minLiveFraction,omitempty"`
	HealthMultiplier          float64          `json:"healthMultiplier,omitempty"`
	DynamicTTL                int              `json:"dynamicTTL,omitempty"`
	MaxUnreachablePenalty     int              `json:"maxUnreachablePenalty,omitempty"`
	MapName                   string           `json:"mapName,omitempty"`
	HandoutLimit              int              `json:"handoutLimit"`
	HandoutMode               string           `json:"handoutMode"`
	FailoverDelay             int              `json:"failoverDelay,omitempty"`
	BackupCName               string           `json:"backupCName,omitempty"`
	FailbackDelay             int              `json:"failbackDelay,omitempty"`
	LoadImbalancePercentage   float64          `json:"loadImbalancePercentage,omitempty"`
	HealthMax                 float64          `json:"healthMax,omitempty"`
	GhostDemandReporting      bool             `json:"ghostDemandReporting"`
	Comments                  string           `json:"comments,omitempty"`
	CName                     string           `json:"cname,omitempty"`
	WeightedHashBitsForIPv4   int              `json:"weightedHashBitsForIPv4,omitempty"`
	WeightedHashBitsForIPv6   int              `json:"weightedHashBitsForIPv6,omitempty"`
	TrafficTargets            []*TrafficTarget `json:"trafficTargets,omitempty"`
	Links                     []*Link          `json:"links,omitempty"`
	LivenessTests             []*LivenessTest  `json:"livenessTests,omitempty"`
}

// PropertyList contains a list of property items
type PropertyList struct {
	PropertyItems []*Property `json:"items"`
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
	tt := value.([]*TrafficTarget)
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

func (g *gtm) ListProperties(ctx context.Context, domainName string) ([]*Property, error) {
	logger := g.Log(ctx)
	logger.Debug("ListProperties")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties", domainName)
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

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.PropertyItems, nil
}

func (g *gtm) GetProperty(ctx context.Context, propertyName, domainName string) (*Property, error) {
	logger := g.Log(ctx)
	logger.Debug("GetProperty")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, propertyName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetProperty request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result Property
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetProperty request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateProperty(ctx context.Context, property *Property, domainName string) (*PropertyResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateProperty")

	return property.save(ctx, g, domainName)
}

func (g *gtm) UpdateProperty(ctx context.Context, property *Property, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateProperty")

	stat, err := property.save(ctx, g, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save Property updates method
func (p *Property) save(ctx context.Context, g *gtm, domainName string) (*PropertyResponse, error) {

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("property validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, p.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result PropertyResponse
	resp, err := g.Exec(req, &result, p)
	if err != nil {
		return nil, fmt.Errorf("property request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) DeleteProperty(ctx context.Context, property *Property, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteProperty")

	if err := property.Validate(); err != nil {
		return nil, fmt.Errorf("DeleteProperty validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, property.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteProperty request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}
