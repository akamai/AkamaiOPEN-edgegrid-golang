package gtm

import (
	"context"
	"fmt"
	"net/http"

	"reflect"
	"strings"
	"unicode"
)

// Domains contains operations available on a Domain resource.
type Domains interface {
	// NullFieldMap retrieves map of null fields.
	NullFieldMap(context.Context, *Domain) (*NullFieldMapStruct, error)
	// GetDomainStatus retrieves current status for the given domain name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-status-current
	GetDomainStatus(context.Context, string) (*ResponseStatus, error)
	// ListDomains retrieves all Domains.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-domains
	ListDomains(context.Context) ([]*DomainItem, error)
	// GetDomain retrieves a Domain with the given domain name.
	//
	// See: https://techdocs.akamai.com/gtm/reference/get-domain
	GetDomain(context.Context, string) (*Domain, error)
	// CreateDomain creates domain.
	//
	// See: https://techdocs.akamai.com/gtm/reference/post-domain
	CreateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
	// DeleteDomain is a method applied to a domain object resulting in removal.
	//
	// See: ** Not Supported by API **
	DeleteDomain(context.Context, *Domain) (*ResponseStatus, error)
	// UpdateDomain is a method applied to a domain object resulting in an update.
	//
	// See: https://techdocs.akamai.com/gtm/reference/put-domain
	UpdateDomain(context.Context, *Domain, map[string]string) (*ResponseStatus, error)
}

// The Domain data structure represents a GTM domain
type Domain struct {
	Name                         string          `json:"name"`
	Type                         string          `json:"type"`
	ASMaps                       []*ASMap        `json:"asMaps,omitempty"`
	Resources                    []*Resource     `json:"resources,omitempty"`
	DefaultUnreachableThreshold  float32         `json:"defaultUnreachableThreshold,omitempty"`
	EmailNotificationList        []string        `json:"emailNotificationList,omitempty"`
	MinPingableRegionFraction    float32         `json:"minPingableRegionFraction,omitempty"`
	DefaultTimeoutPenalty        int             `json:"defaultTimeoutPenalty,omitempty"`
	Datacenters                  []*Datacenter   `json:"datacenters,omitempty"`
	ServermonitorLivenessCount   int             `json:"servermonitorLivenessCount,omitempty"`
	RoundRobinPrefix             string          `json:"roundRobinPrefix,omitempty"`
	ServermonitorLoadCount       int             `json:"servermonitorLoadCount,omitempty"`
	PingInterval                 int             `json:"pingInterval,omitempty"`
	MaxTTL                       int64           `json:"maxTTL,omitempty"`
	LoadImbalancePercentage      float64         `json:"loadImbalancePercentage,omitempty"`
	DefaultHealthMax             float64         `json:"defaultHealthMax,omitempty"`
	LastModified                 string          `json:"lastModified,omitempty"`
	Status                       *ResponseStatus `json:"status,omitempty"`
	MapUpdateInterval            int             `json:"mapUpdateInterval,omitempty"`
	MaxProperties                int             `json:"maxProperties,omitempty"`
	MaxResources                 int             `json:"maxResources,omitempty"`
	DefaultSSLClientPrivateKey   string          `json:"defaultSslClientPrivateKey,omitempty"`
	DefaultErrorPenalty          int             `json:"defaultErrorPenalty,omitempty"`
	Links                        []*Link         `json:"links,omitempty"`
	Properties                   []*Property     `json:"properties,omitempty"`
	MaxTestTimeout               float64         `json:"maxTestTimeout,omitempty"`
	CNameCoalescingEnabled       bool            `json:"cnameCoalescingEnabled"`
	DefaultHealthMultiplier      float64         `json:"defaultHealthMultiplier,omitempty"`
	ServermonitorPool            string          `json:"servermonitorPool,omitempty"`
	LoadFeedback                 bool            `json:"loadFeedback"`
	MinTTL                       int64           `json:"minTTL,omitempty"`
	GeographicMaps               []*GeoMap       `json:"geographicMaps,omitempty"`
	CIDRMaps                     []*CIDRMap      `json:"cidrMaps,omitempty"`
	DefaultMaxUnreachablePenalty int             `json:"defaultMaxUnreachablePenalty"`
	DefaultHealthThreshold       float64         `json:"defaultHealthThreshold,omitempty"`
	LastModifiedBy               string          `json:"lastModifiedBy,omitempty"`
	ModificationComments         string          `json:"modificationComments,omitempty"`
	MinTestInterval              int             `json:"minTestInterval,omitempty"`
	PingPacketSize               int             `json:"pingPacketSize,omitempty"`
	DefaultSSLClientCertificate  string          `json:"defaultSslClientCertificate,omitempty"`
	EndUserMappingEnabled        bool            `json:"endUserMappingEnabled"`
	SignAndServe                 bool            `json:"signAndServe"`
	SignAndServeAlgorithm        *string         `json:"signAndServeAlgorithm"`
}

// DomainsList contains a list of domain items
type DomainsList struct {
	DomainItems []*DomainItem `json:"items"`
}

// DomainItem is a DomainsList item
type DomainItem struct {
	AcgID                 string  `json:"acgId"`
	LastModified          string  `json:"lastModified"`
	Links                 []*Link `json:"links"`
	Name                  string  `json:"name"`
	Status                string  `json:"status"`
	LastModifiedBy        string  `json:"lastModifiedBy"`
	ChangeID              string  `json:"changeId"`
	ActivationState       string  `json:"activationState"`
	ModificationComments  string  `json:"modificationComments"`
	SignAndServe          bool    `json:"signAndServe"`
	SignAndServeAlgorithm string  `json:"signAndServeAlgorithm"`
	DeleteRequestID       string  `json:"deleteRequestId"`
}

// Validate validates Domain
func (d *Domain) Validate() error {
	if len(d.Name) < 1 {
		return fmt.Errorf("Domain is missing Name")
	}
	if len(d.Type) < 1 {
		return fmt.Errorf("Domain is missing Type")
	}

	return nil
}

func (g *gtm) GetDomainStatus(ctx context.Context, domainName string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("GetDomainStatus")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/status/current", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseStatus
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) ListDomains(ctx context.Context) ([]*DomainItem, error) {
	logger := g.Log(ctx)
	logger.Debug("ListDomains")

	getURL := fmt.Sprintf("/config-gtm/v1/domains")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListDomains request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result DomainsList
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("ListDomains request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.DomainItems, nil
}

func (g *gtm) GetDomain(ctx context.Context, domainName string) (*Domain, error) {
	logger := g.Log(ctx)
	logger.Debug("GetDomain")

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result Domain
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return &result, nil
}

// save method; Create or Update
func (d *Domain) save(_ context.Context, g *gtm, queryArgs map[string]string, req *http.Request) (*DomainResponse, error) {
	// set schema version
	setVersionHeader(req, schemaVersion)

	// Look for optional args
	if len(queryArgs) > 0 {
		q := req.URL.Query()
		if val, ok := queryArgs["contractId"]; ok {
			q.Add("contractId", strings.TrimPrefix(val, "ctr_"))
		}
		if val, ok := queryArgs["gid"]; ok {
			q.Add("gid", strings.TrimPrefix(val, "grp_"))
		}
		req.URL.RawQuery = q.Encode()
	}

	var result DomainResponse
	resp, err := g.Exec(req, &result, d)
	if err != nil {
		return nil, fmt.Errorf("domain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, g.Error(resp)
	}

	return &result, nil
}

func (g *gtm) CreateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*DomainResponse, error) {
	logger := g.Log(ctx)
	logger.Debug("CreateDomain")

	if err := domain.Validate(); err != nil {
		return nil, fmt.Errorf("CreateDomain validation failed. %w", err)
	}

	postURL := fmt.Sprintf("/config-gtm/v1/domains/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateDomain request: %w", err)
	}

	return domain.save(ctx, g, queryArgs, req)
}

func (g *gtm) UpdateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("UpdateDomain")

	if err := domain.Validate(); err != nil {
		return nil, fmt.Errorf("UpdateDomain validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateDomain request: %w", err)
	}

	stat, err := domain.save(ctx, g, queryArgs, req)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

func (g *gtm) DeleteDomain(ctx context.Context, domain *Domain) (*ResponseStatus, error) {
	logger := g.Log(ctx)
	logger.Debug("DeleteDomain")

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)

	var result ResponseBody
	resp, err := g.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("DeleteDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	return result.Status, nil
}

// NullPerObjectAttributeStruct represents core and child null object attributes
type NullPerObjectAttributeStruct struct {
	CoreObjectFields  map[string]string
	ChildObjectFields map[string]interface{} // NullObjectAttributeStruct
}

// NullFieldMapStruct returned null Objects structure
type NullFieldMapStruct struct {
	Domain      NullPerObjectAttributeStruct            // entry is domain
	Properties  map[string]NullPerObjectAttributeStruct // entries are properties
	Datacenters map[string]NullPerObjectAttributeStruct // entries are datacenters
	Resources   map[string]NullPerObjectAttributeStruct // entries are resources
	CidrMaps    map[string]NullPerObjectAttributeStruct // entries are cidrmaps
	GeoMaps     map[string]NullPerObjectAttributeStruct // entries are geomaps
	AsMaps      map[string]NullPerObjectAttributeStruct // entries are asmaps
}

// ObjectMap represents ObjectMap datatype
type ObjectMap map[string]interface{}

func (g *gtm) NullFieldMap(ctx context.Context, domain *Domain) (*NullFieldMapStruct, error) {
	logger := g.Log(ctx)
	logger.Debug("NullFieldMap")

	if err := domain.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed. %w", err)
	}

	var nullFieldMap = &NullFieldMapStruct{}
	var domFields = NullPerObjectAttributeStruct{}
	domainMap := make(map[string]string)
	var objMap ObjectMap

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := g.Exec(req, &objMap)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, g.Error(resp)
	}

	for i, d := range objMap {
		objVal := fmt.Sprint(d)
		if fmt.Sprintf("%T", d) == "<nil>" {
			if objVal == "<nil>" {
				domainMap[makeFirstCharUpperCase(i)] = ""
			}
			continue
		}
		list, ok := d.([]interface{})
		if !ok {
			continue
		}

		switch i {
		case "properties":
			nullFieldMap.Properties = processObjectList(list)
		case "datacenters":
			nullFieldMap.Datacenters = processObjectList(list)
		case "resources":
			nullFieldMap.Resources = processObjectList(list)
		case "cidrMaps":
			nullFieldMap.CidrMaps = processObjectList(list)
		case "geographicMaps":
			nullFieldMap.GeoMaps = processObjectList(list)
		case "asMaps":
			nullFieldMap.AsMaps = processObjectList(list)
		}
	}

	domFields.CoreObjectFields = domainMap
	nullFieldMap.Domain = domFields

	return nullFieldMap, nil

}

func makeFirstCharUpperCase(origString string) string {
	a := []rune(origString)
	a[0] = unicode.ToUpper(a[0])
	// hack
	if origString == "cname" {
		a[1] = unicode.ToUpper(a[1])
	}
	return string(a)
}

func processObjectList(objectList []interface{}) map[string]NullPerObjectAttributeStruct {
	nullObjectsList := make(map[string]NullPerObjectAttributeStruct)
	for _, obj := range objectList {
		nullObjectFields := NullPerObjectAttributeStruct{}
		objectName := ""
		objectDCID := ""
		objectMap := make(map[string]string)
		objectChildList := make(map[string]interface{})
		for objF, objD := range obj.(map[string]interface{}) {
			objVal := fmt.Sprint(objD)
			switch fmt.Sprintf("%T", objD) {
			case "<nil>":
				if objVal == "<nil>" {
					objectMap[makeFirstCharUpperCase(objF)] = ""
				}
			case "map[string]interface {}":
				// include null stand alone struct elements in core
				for moName, moValue := range objD.(map[string]interface{}) {
					if fmt.Sprintf("%T", moValue) == "<nil>" {
						objectMap[makeFirstCharUpperCase(moName)] = ""
					}
				}
			case "[]interface {}":
				iSlice := objD.([]interface{})
				if len(iSlice) > 0 && reflect.TypeOf(iSlice[0]).Kind() != reflect.String && reflect.TypeOf(iSlice[0]).Kind() != reflect.Int64 && reflect.TypeOf(iSlice[0]).Kind() != reflect.Float64 && reflect.TypeOf(iSlice[0]).Kind() != reflect.Int32 {
					objectChildList[makeFirstCharUpperCase(objF)] = processObjectList(objD.([]interface{}))
				}
			default:
				if objF == "name" {
					objectName = objVal
				}
				if objF == "datacenterId" {
					objectDCID = objVal
				}
			}
		}
		nullObjectFields.CoreObjectFields = objectMap
		nullObjectFields.ChildObjectFields = objectChildList

		if objectDCID == "" {
			if objectName != "" {
				nullObjectsList[objectName] = nullObjectFields
			} else {
				nullObjectsList["unknown"] = nullObjectFields // TODO: What if mnore than one?
			}
		} else {
			nullObjectsList[objectDCID] = nullObjectFields
		}
	}

	return nullObjectsList
}
