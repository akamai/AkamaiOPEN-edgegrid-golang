package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
 
        "fmt"
        "net/http/httputil"
)

//
// Handle Operations on gtm asmaps
// Based on 1.3 schema
//

type AsAssignment struct {
        DatacenterBase                                                    
        AsNumbers                []int64           `json:"asNumbers"`
}

// AsMap  represents a GTM AsMap   
type AsMap struct {
        DefaultDatacenter        *DatacenterBase   `json:"defaultDatacenter"`
        Assignments              []*AsAssignment   `json:"assignments,omitempty"`
        Name                     string            `json:"name,omitempty"`
        Links                    []*Link           `json:"links,omitempty"`
}

// NewAsMap creates a new asMap
func NewAsMap(asmapname string) *AsMap {
	asmap := &AsMap{Name: asmapname}
        return asmap
}

// GetAsMap retrieves a asMap with the given name.                     
func GetAsMap(asmapname, domainname string) (*AsMap, error) {
	as := NewAsMap(asmapname)
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/"+domainname+"/as-maps/"+asmapname,
		nil,
	)
	if err != nil {
		return nil, err
	}

        SetHeader(req)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &CommonError{entityName: "asMap", name: asmapname}
	} else {
		err = client.BodyJSON(res, as)
		if err != nil {
			return nil, err
		}

		return as, nil
	}
}

// Instantiate new Assignment struct
func (as *AsMap) NewAssignment(dcid int, nickname string) *AsAssignment {
        asAssign := &AsAssignment{}
        asAssign.DatacenterId = dcid 
        asAssign.Nickname =nickname

        return asAssign

}

// Instantiate new Default Datacenter Struct
func (as *AsMap) NewDefaultDatacenter(dcid int) *DatacenterBase {
        return &DatacenterBase{DatacenterId: dcid}
}

// Create asMap in provided domain                        
func (as *AsMap) Create(domainname string) (*AsMapResponse, error) {

        // Use common code. Any specific validation needed?

        return as.save(domainname)

}

// Update AsMap in given domain
func (as *AsMap) Update(domainname string) (*ResponseStatus, error) {

        // common code
   
        stat, err := as.save(domainname)
        if err != nil {
                return nil, err
        }
        return stat.Status, err

}

// Save AsMap in given domain. Common path for Create and Update.
func (as *AsMap) save(domainname string) (*AsMapResponse, error) {

        req, err := client.NewJSONRequest(
                Config,
                "PUT",
                "/config-gtm/v1/domains/"+domainname+"/as-maps/"+as.Name,
                as,
        )
        if err != nil {
                return nil, err
        }

        SetHeader(req)

        b, err := httputil.DumpRequestOut(req, true)
        if err == nil {
                fmt.Println(string(b))
        }

        res, err := client.Do(Config, req)

        // Network error
        if err != nil {
                return nil, &CommonError{
                        entityName:       "asMap",
                        name:             as.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "asMap", name: as.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &AsMapResponse{}
        // Unmarshall whole response body for updated entity and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody, nil
}

// Delete AsMap method
func (as *AsMap) Delete(domainname string) (*ResponseStatus, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
                "/config-gtm/v1/domains/"+domainname+"/as-maps/"+as.Name,
                nil,
        )
        if err != nil {
                return nil, err
        }

        SetHeader(req)

        res, err := client.Do(Config, req)
        if err != nil {
                return nil, err
        }

        res, err = client.Do(Config, req)

        // Network error
        if err != nil {
                return nil, &CommonError{
                        entityName:       "asMap",
                        name:             as.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "asMap", name:as.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &ResponseBody{}
        // Unmarshall whole response body in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil
}

// Add/Update AsMap element
func (as *AsMap) AddElement(element interface{}) error {

        // Do we need?

	return nil
}

// Remove AsMap element
func (as *AsMap) RemoveElement(element interface{}) error {

        // What does this mean?

	return nil
}

// Retrieve a specific element
func (as *AsMap) GetElement(element interface{}) interface{} {

        // useful?
        return nil

}
