package gtm

import (
	"fmt"
	"net/http"
)

// default schema version
var schemaVersion = "1.4"

// internal method to set version. passed in as string
func setVersionHeader(req *http.Request, version string) {

	req.Header.Set("Accept", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))

	if req.Method != "GET" {
		req.Header.Set("Content-Type", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))
	}

	return

}

// NewDefaultDatacenter instantiates new Default Datacenter Struct
func (g *gtm) NewDefaultDatacenter(dcID int) *DatacenterBase {
	return &DatacenterBase{DatacenterID: dcID}
}

// ResponseStatus is returned on Create, Update or Delete operations for all entity types
type ResponseStatus struct {
	ChangeID              string  `json:"changeId,omitempty"`
	Links                 *[]Link `json:"links,omitempty"`
	Message               string  `json:"message,omitempty"`
	PassingValidation     bool    `json:"passingValidation,omitempty"`
	PropagationStatus     string  `json:"propagationStatus,omitempty"`
	PropagationStatusDate string  `json:"propagationStatusDate,omitempty"`
}

// ResponseBody is a generic response struct
type ResponseBody struct {
	Resource interface{}     `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// DomainResponse contains a response after creating or updating Domain
type DomainResponse struct {
	Resource *Domain         `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// DatacenterResponse contains a response after creating or updating Datacenter
type DatacenterResponse struct {
	Status   *ResponseStatus `json:"status"`
	Resource *Datacenter     `json:"resource"`
}

// PropertyResponse contains a response after creating or updating Property
type PropertyResponse struct {
	Resource *Property       `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// ResourceResponse contains a response after creating or updating Resource
type ResourceResponse struct {
	Resource *Resource       `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// CIDRMapResponse contains a response after creating or updating CIDRMap
type CIDRMapResponse struct {
	Resource *CIDRMap        `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// GeoMapResponse contains a response after creating or updating GeoMap
type GeoMapResponse struct {
	Resource *GeoMap         `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// ASMapResponse contains a response after creating or updating ASMap
type ASMapResponse struct {
	Resource *ASMap          `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// Link is Probably THE most common type
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// LoadObject contains information about the load reporting interface
type LoadObject struct {
	LoadObject     string   `json:"loadObject,omitempty"`
	LoadObjectPort int      `json:"loadObjectPort,omitempty"`
	LoadServers    []string `json:"loadServers,omitempty"`
}

// NewLoadObject returns a new LoadObject structure
func NewLoadObject() *LoadObject {
	return &LoadObject{}
}

// DatacenterBase is a placeholder for default Datacenter
type DatacenterBase struct {
	Nickname     string `json:"nickname,omitempty"`
	DatacenterID int    `json:"datacenterId"`
}
