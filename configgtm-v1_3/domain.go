package configgtm

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"net/http"
	"strings"
)

//
// Support gtm domains thru Edgegrid
// Based on 1.3 Schema
//

// The Domain data structure represents a GTM domain
type Domain struct {
	Name                         string          `json:"name"`
	Type                         string          `json:"type"`
	AsMaps                       []*AsMap        `json:"asMaps"`
	Resources                    []*Resource     `json:"resources"`
	DefaultUnreachableThreshold  float32         `json:"defaultUnreachableThreshold"`
	EmailNotificationList        []string        `json:"emailNotificationList"`
	MinPingableRegionFraction    float32         `json:"minPingableRegionFraction"`
	DefaultTimeoutPenalty        int             `json:"defaultTimeoutPenalty"`
	Datacenters                  []*Datacenter   `json:"datacenters"`
	ServermonitorLivenessCount   int             `json:"servermonitorLivenessCount"`
	RoundRobinPrefix             string          `json:"roundRobinPrefix"`
	ServermonitorLoadCount       int             `json:"servermonitorLoadCount"`
	PingInterval                 int             `json:"pingInterval"`
	MaxTTL                       int64           `json:"maxTTL"`
	LoadImbalancePercentage      float64         `json:"loadImbalancePercentage"`
	DefaultHealthMax             int             `json:"defaultHealthMax"`
	LastModified                 string          `json:"lastModified"`
	Status                       *ResponseStatus `json:"status"`
	MapUpdateInterval            int             `json:"mapUpdateInterval"`
	MaxProperties                int             `json:"maxProperties"`
	MaxResources                 int             `json:"maxResources"`
	DefaultSslClientPrivateKey   string          `json:"defaultSslClientPrivateKey"`
	DefaultErrorPenalty          int             `json:"defaultErrorPenalty"`
	Links                        []*Link         `json:"links"`
	Properties                   []*Property     `json:"properties"`
	MaxTestTimeout               float64         `json:"maxTestTimeout"`
	CnameCoalescingEnabled       bool            `json:"cnameCoalescingEnabled"`
	DefaultHealthMultiplier      int             `json:"defaultHealthMultiplier"`
	ServermonitorPool            string          `json:"servermonitorPool"`
	LoadFeedback                 bool            `json:"loadFeedback"`
	MinTTL                       int64           `json:"minTTL"`
	GeographicMaps               []*GeoMap       `json:"geographicMaps"`
	CidrMaps                     []*CidrMap      `json:"cidrMaps"`
	DefaultMaxUnreachablePenalty int             `json:"defaultMaxUnreachablePenalty"`
	DefaultHealthThreshold       int             `json:"defaultHealthThreshold"`
	LastModifiedBy               string          `json:"lastModifiedBy"`
	ModificationComments         string          `json:"modificationComments"`
	MinTestInterval              int             `json:"minTestInterval"`
	PingPacketSize               int             `json:"pingPacketSize"`
	DefaultSslClientCertificate  string          `json:"defaultSslClientCertificate"`
}

type DomainsList struct {
	DomainItems []*DomainItem `json:"items"`
}

// DomainItem is a DomainsList item
type DomainItem struct {
	AcgId        string  `json:"acgId"`
	LastModified string  `json:"lastModified"`
	Links        []*Link `json:"links"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
}

// NewDomain is a utility function that creates a new Domain object.
func NewDomain(domainName, domainType string) *Domain {
	domain := &Domain{}
	domain.Name = domainName
	domain.Type = domainType
	return domain
}

// GetStatus retrieves current status for the given domainname.
func GetDomainStatus(domainName string) (*ResponseStatus, error) {
	stat := &ResponseStatus{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/status/current", domainName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Domain", name: domainName}
	} else {
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}

		return stat, nil
	}
}

// ListDomains retrieves all Domains.
func ListDomains() ([]*DomainItem, error) {
	domains := &DomainsList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/",
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Domain"}
	} else {
		err = client.BodyJSON(res, domains)
		if err != nil {
			return nil, err
		}

		return domains.DomainItems, nil
	}
}

// GetDomain retrieves a Domain with the given domainname.
func GetDomain(domainName string) (*Domain, error) {
	domain := NewDomain(domainName, "basic")
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domainName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Domain", name: domainName}
	} else {
		err = client.BodyJSON(res, domain)
		if err != nil {
			return nil, err
		}

		return domain, nil
	}
}

// Save method; Create or Update
func (domain *Domain) save(queryArgs map[string]string, req *http.Request) (*DomainResponse, error) {

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

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domain.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domain.Name, apiErrorMessage: err.Detail, err: err}
	}

	// TODO: What validation can we do? E.g. if not equivalent there was a concurrent change...
	responseBody := &DomainResponse{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil

}

// Create is a method applied to a domain object resulting in creation.
func (domain *Domain) Create(queryArgs map[string]string) (*DomainResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf("/config-gtm/v1/domains/"),
		domain,
	)
	if err != nil {
		return nil, err
	}

	return domain.save(queryArgs, req)

}

// Update is a method applied to a domain object resulting in an update.
func (domain *Domain) Update(queryArgs map[string]string) (*ResponseStatus, error) {

	// Any validation to do?
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name),
		domain,
	)
	if err != nil {
		return nil, err
	}

	stat, err := domain.save(queryArgs, req)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Delete is a method applied to a domain object resulting in removal.
func (domain *Domain) Delete() (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domain.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domain.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil

}
