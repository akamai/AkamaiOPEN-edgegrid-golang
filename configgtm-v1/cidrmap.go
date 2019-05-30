package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

//
// Handle Operations on gtm cidrmaps
// Based on 1.3 schema
//

type CidrAssignment struct {
        DatacenterBase                                                    
        Blocks                   []string          `json:"blocks"`
}

// CidrMap  represents a GTM cidrMap   
type CidrMap struct {
        DefaultDatacenter        *DatacenterBase        `json:"defaultDatacenter"`
        Assignments              []*CidrAssignment      `json:"assignments,omitempty"`
        Name                     string                 `json:"name"`
        Links                    []*Link                `json:"links, omitempty"`
}

type CidrMapList struct {
        CidrMapItems             []*CidrMap             `json:"items"`
}


// NewCidrMap creates a new CidrMap object
func NewCidrMap(cidrmapname string) *CidrMap {
	cidrmap := &CidrMap{Name: cidrmapname}
        return cidrmap
}

// ListCidrMap retreieves all CidrMaps
func ListCidrMaps(domainname string) ([]*CidrMap, error) {
        cidrs := &CidrMapList{}
        req, err := client.NewRequest(
                Config,
                "GET",
                "/config-gtm/v1/domains/"+domainname+"/cidr-maps",
                nil,
        )
        if err != nil {
                return nil, err
        }       

        SetHeader(req)

        printHttpRequest(req, true)

        res, err := client.Do(Config, req)
        if err != nil {
                return nil, err
        }

        printHttpResponse(res, true)

        if client.IsError(res) && res.StatusCode != 404 {
                return nil, client.NewAPIError(res)
        } else if res.StatusCode == 404 {
                return nil, &CommonError{entityName: "cidrMap"}
        } 
        err = client.BodyJSON(res, cidrs)
        if err != nil {
                return nil, err
        }

        return cidrs.CidrMapItems, nil
         
}

// GetCidrMap retrieves a CidrMap with the given name.                     
func GetCidrMap(cidrmapname, domainname string) (*CidrMap, error) {
	cidr := NewCidrMap(cidrmapname)
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/"+domainname+"/cidr-maps/"+cidrmapname,
		nil,
	)
	if err != nil {
		return nil, err
	}

        SetHeader(req)

        printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

        printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &CommonError{entityName: "cidrMap", name: cidrmapname}
	} else {
		err = client.BodyJSON(res, cidr)
		if err != nil {
			return nil, err
		}

		return cidr, nil
	}
}

// Instantiate new Assignment struct
func (cidr *CidrMap) NewAssignment(dcid int, nickname string) *CidrAssignment {
        cidrAssign := &CidrAssignment{}
        cidrAssign.DatacenterId = dcid
        cidrAssign.Nickname = nickname

        return cidrAssign 
}

// Instantiate new Default Datacenter Struct
func (cidr *CidrMap) NewDefaultDatacenter(dcid int) *DatacenterBase {
        return &DatacenterBase{DatacenterId: dcid}
}

// Create CidrMap in provided domain                        
func (cidr *CidrMap) Create(domainname string) (*CidrMapResponse, error) {

        // Use common code. Any specific validation needed?

        return cidr.save(domainname)

}

// Update CidrMap in given domain
func (cidr *CidrMap) Update(domainname string) (*ResponseStatus, error) {

        // common code
   
        stat, err := cidr.save(domainname)
        if err != nil {
                return nil, err
        } 
        return stat.Status, err

}

// Save CidrMap in given domain. Common path for Create and Update.
func (cidr *CidrMap) save(domainname string) (*CidrMapResponse, error) {

        req, err := client.NewJSONRequest(
                Config,
                "PUT",
                "/config-gtm/v1/domains/"+domainname+"/cidr-maps/"+cidr.Name,
                cidr,
        )
        if err != nil {
                return nil, err
        }

        SetHeader(req)

        printHttpRequest(req, true)

        res, err := client.Do(Config, req)

        printHttpResponse(res, true)

        // Network error
        if err != nil {
                return nil, &CommonError{
                        entityName:       "cidrMap",
                        name:             cidr.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "cidrMap", name: cidr.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &CidrMapResponse{}
        // Unmarshall whole response body for updated entity and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody, nil
}

// Delete CidrMap method
func (cidr *CidrMap) Delete(domainname string) (*ResponseStatus, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
                "/config-gtm/v1/domains/"+domainname+"/cidr-maps/"+cidr.Name,
                nil,
        )
        if err != nil {
                return nil, err
        }

        SetHeader(req)

        printHttpRequest(req, true)

        res, err := client.Do(Config, req)
        if err != nil {
                return nil, err
        }

        printHttpResponse(res, true)

        // Network error
        if err != nil {
                return nil, &CommonError{
                        entityName:       "cidrMap",
                        name:             cidr.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "cidrMap", name:cidr.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &ResponseBody{}
        // Unmarshall whole response body in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil

}

// Add/Update CidrMap element
func (cidr *CidrMap) AddElement(element interface{}) error {

        // Do we need?

	return nil
}

// Remove CidrMap element
func (cidr *CidrMap) RemoveElement(element interface{}) error {

        // What does this mean?

	return nil
}

// Retrieve a specific element
func (cidr *CidrMap) GetElement(element interface{}) interface{} {

        // useful?
        return nil

}
