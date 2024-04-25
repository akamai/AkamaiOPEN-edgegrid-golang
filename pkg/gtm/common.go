package gtm

import (
	"fmt"
	"net/http"
)

// default schema version
var schemaVersion = "1.6"

// internal method to set version. passed in as string
func setVersionHeader(req *http.Request, version string) {

	req.Header.Set("Accept", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))

	if req.Method != "GET" {
		req.Header.Set("Content-Type", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))
	}

	return

}

// ResponseStatus is returned on Create, Update or Delete operations for all entity types
type ResponseStatus struct {
	ChangeID              string `json:"changeId,omitempty"`
	Links                 []Link `json:"links,omitempty"`
	Message               string `json:"message,omitempty"`
	PassingValidation     bool   `json:"passingValidation,omitempty"`
	PropagationStatus     string `json:"propagationStatus,omitempty"`
	PropagationStatusDate string `json:"propagationStatusDate,omitempty"`
}

// DatacenterResponse contains a response after creating or updating Datacenter
type DatacenterResponse struct {
	Status   *ResponseStatus `json:"status"`
	Resource *Datacenter     `json:"resource"`
}

// ResourceResponse contains a response after creating or updating Resource
type ResourceResponse struct {
	Resource *Resource       `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// Link is Probably THE most common type
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

// LoadObject contains information about the load reporting interface
type LoadObject struct {
	LoadObject     string   `json:"loadObject,omitempty"`
	LoadObjectPort int      `json:"loadObjectPort,omitempty"`
	LoadServers    []string `json:"loadServers,omitempty"`
}

// DatacenterBase is a placeholder for default Datacenter
type DatacenterBase struct {
	Nickname     string `json:"nickname,omitempty"`
	DatacenterID int    `json:"datacenterId"`
}
