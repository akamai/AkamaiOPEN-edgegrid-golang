package configgtm

//
// Common data types
// Based on 1.3 schemas
//

// response Status is returned on Create, Update or Delete operations for all entity types
type ResponseStatus struct {
        ChangeId              string          `json:"changeId"`
        Links                 *[]Link         `json:"links"`
        Message               string          `json:"message"`
        PassingValidation     bool            `json:"passingValidation"`
        PropagationStatus     string          `json:"propagationStatus"`
        PropagationStatusDate string          `json:"propagationStatusDate"`
}

func NewResponseStatus() *ResponseStatus {

        return &ResponseStatus{}

}

// Generic response structs
type ResponseBody struct {
        Resource       interface{}            `json:"resource"`
        Status         *ResponseStatus        `json:"status"` 
}

// Response structs by Entity Type
type DomainResponse struct {
        Resource       *Domain                `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

type DatacenterResponse struct {
        Status         *ResponseStatus        `json:"status"`
        Resource       *Datacenter            `json:"resource"`
}

type PropertyResponse struct {
        Resource       *Property              `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

type ResourceResponse struct {
        Resource       *Resource              `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

type CidrMapResponse struct {
        Resource       *CidrMap               `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

type GeoMapResponse struct {
        Resource       *GeoMap                `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

type AsMapResponse struct {
        Resource       *AsMap                 `json:"resource"`
        Status         *ResponseStatus        `json:"status"`
}

// Probably THE most common type
type Link struct {
        Rel            string                 `json:"rel"`
        Href           string                 `json:"href"`
}

//
type LoadObject struct {
        LoadObject               string            `json:"loadObject, omitempty"`
        LoadObjectPort           int               `json:"loadObjectPort, omitempty"`
        LoadServers              []string          `json:"loadServers, omitempty"`
}

func NewLoadObject() *LoadObject {
        return &LoadObject{}
}

type DatacenterBase struct {
        Nickname                 string            `json:"nickname"`
        DatacenterId             int               `json:"datacenterId"`
}

func NewDatacenterBase() *DatacenterBase {
        return &DatacenterBase{}
}

// util method
func logGtmSingleLineOutput(logline string) error {

        return nil

}

func logGtmMultiLineOutput(loglines string) error {

        return nil

}

