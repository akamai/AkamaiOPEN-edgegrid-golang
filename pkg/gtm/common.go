package gtm

import (
	"fmt"
	"net/http"
)

//
// Common data types and methods
// Based on 1.3 schemas
//

// Append url args to req
func appendReqArgs(req *http.Request, queryArgs map[string]string) {

	// Look for optional args
	if len(queryArgs) > 0 {
		q := req.URL.Query()
		for argName, argVal := range queryArgs {
			q.Add(argName, argVal)
		}
		req.URL.RawQuery = q.Encode()
	}

}

// default schema version
// TODO: retrieve from environment or elsewhere in Service Init
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
func (p *gtm) NewDefaultDatacenter(dcid int) *DatacenterBase {
	return &DatacenterBase{DatacenterId: dcid}
}

// ResponseStatus is returned on Create, Update or Delete operations for all entity types
type ResponseStatus struct {
	ChangeId              string  `json:"changeId,omitempty"`
	Links                 *[]Link `json:"links,omitempty"`
	Message               string  `json:"message,omitempty"`
	PassingValidation     bool    `json:"passingValidation,omitempty"`
	PropagationStatus     string  `json:"propagationStatus,omitempty"`
	PropagationStatusDate string  `json:"propagationStatusDate,omitempty"`
}

// NewResponseStatus returns a new ResponseStatus struct
func NewResponseStatus() *ResponseStatus {

	return &ResponseStatus{}

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

// CidrMapResponse contains a response after creating or updating CidrMap
type CidrMapResponse struct {
	Resource *CidrMap        `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// GeoMapResponse contains a response after creating or updating GeoMap
type GeoMapResponse struct {
	Resource *GeoMap         `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// AsMapResponse contains a response after creating or updating AsMap
type AsMapResponse struct {
	Resource *AsMap          `json:"resource"`
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
	DatacenterId int    `json:"datacenterId"`
}

// NewDatacenterBase returns a new DatacenterBase structure
func NewDatacenterBase() *DatacenterBase {
	return &DatacenterBase{}
}
