package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

//
// Handle Operations on gtm geomaps
// Based on 1.3 schema
//

type GeoAssignment struct {
        DatacenterBase                                                    
        Countries                []string          `json:"countries,omitempty"`
}

// GeoMap  represents a GTM GeoMap   
type GeoMap struct {
        DefaultDatacenter        *DatacenterBase        `json:"defaultDatacenter"`
        Assignments              []*GeoAssignment       `json:"assignments,omitempty"`
        Name                     string                 `json:"name"`
        Links                    []*Link                `json:"links,omitempty"`
}

type GeoMapList struct {
        GeoMapItems             []*GeoMap             `json:"items"`
}

// NewGeoMap creates a new GeoMap object
func NewGeoMap(geomapname string) *GeoMap {
	geomap := &GeoMap{Name: geomapname}
        return geomap
}

// ListGeoMap retreieves all GeoMaps
func ListGeoMaps(domainname string) ([]*GeoMap, error) {
        geos := &GeoMapList{}
        req, err := client.NewRequest(
                Config,
                "GET",
                "/config-gtm/v1/domains/"+domainname+"/geographic-maps",
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
                return nil, &CommonError{entityName: "geoMap"}
        } 
        err = client.BodyJSON(res, geos)
        if err != nil {
                return nil, err
        }

        return geos.GeoMapItems, nil
         
}

// GetGeoMap retrieves a GeoMap with the given name.                     
func GetGeoMap(geomapname, domainname string) (*GeoMap, error) {
	geo := NewGeoMap(geomapname)

	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/"+domainname+"/geographic-maps/"+geomapname,
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
		return nil, &CommonError{entityName: "GeographicMap", name: geomapname}
	} else {
		err = client.BodyJSON(res, geo)
		if err != nil {
			return nil, err
		}

		return geo, nil
	}
}

// Instantiate new Assignment struct
func (geo *GeoMap) NewAssignment(dcid int, nickname string) *GeoAssignment {
        geoAssign := &GeoAssignment{}
        geoAssign.DatacenterId = dcid
        geoAssign.Nickname = nickname

        return geoAssign

}       

// Instantiate new Default Datacenter Struct
func (geo *GeoMap) NewDefaultDatacenter(dcid int) *DatacenterBase {
        return &DatacenterBase{DatacenterId: dcid}
}       

// Create GeoMap in provided domain                        
func (geo *GeoMap) Create(domainname string) (*GeoMapResponse, error) {

        // Use common code. Any specific validation needed?

        return geo.save(domainname)

}

// Update GeoMap in given domain
func (geo *GeoMap) Update(domainname string) (*ResponseStatus, error) {

        // common code
   
        stat, err := geo.save(domainname)
        if err != nil {
                return nil, err
        }
        return stat.Status, err

}

// Save GeoMap in given domain. Common path for Create and Update.
func (geo *GeoMap) save(domainname string) (*GeoMapResponse, error) {

        req, err := client.NewJSONRequest(
                Config,
                "PUT",
                "/config-gtm/v1/domains/"+domainname+"/geographic-maps/"+geo.Name,
                geo,
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
                        entityName:       "geographicMap",
                        name:             geo.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "geographicMap", name: geo.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &GeoMapResponse{}
        // Unmarshall whole response body for updated entity and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody, nil
}

// Delete GeoMap method
func (geo *GeoMap) Delete(domainname string) (*ResponseStatus, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
                "/config-gtm/v1/domains/"+domainname+"/geographic-maps/"+geo.Name,
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
                        entityName:       "geographicMap",
                        name:             geo.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "geographicMap", name:geo.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &ResponseBody{}
        // Unmarshall whole response body in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil
}

