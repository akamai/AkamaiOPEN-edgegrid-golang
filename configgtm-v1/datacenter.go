package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
        "strconv"

)

//
// Handle Operations on gtm datacenters
// Based on 1.3 schema
//

// Datacenter represents a GTM datacenter
type Datacenter struct {
        City                          string                 `json:"city,omitempty"`
        CloneOf                       int                    `json:"cloneOf,omitempty"`
        CloudServerTargeting          bool                   `json:"cloudServerTargeting,omitempty"`
        Continent                     string                 `json:"continent,omitempty"`
        Country                       string                 `json:"country,omitempty"`
        DefaultLoadObject             *LoadObject            `json:"defaultLoadObject,omitempty"`
        Latitude                      float64                `json:"latitude,omitempty"`
        Links                         []*Link                 `json:"links,omitempty"`
        Longitude                     float64                `json:"longitude,omitempty"`
        Nickname                      string                 `json:"nickname,omitempty"`
        PingInterval                  int                    `json:"pingInterval,omitempty"`
        PingPacketSize                int                    `json:"pingPacketSize,omitempty"`
        DatacenterId                  int                    `json:"datacenterId,omitempty"`
        ScorePenalty                  int                    `json:"scorePenalty,omitempty"`
        ServermonitorLivenessCount    int                    `json:"servermonitorLivenessCount,omitempty"`
        ServermonitorLoadCount        int                    `json:"servermonitorLoadCount,omitempty"`
        ServermonitorPool             string                 `json:"servermonitorPool,omitempty"`
        SstateOrProvince              string                 `json:"stateOrProvince,omitempty"`
        Virtual                       bool                   `json:"virtual,omitempty"`
}

type DatacenterList struct {
        DatacenterItems          []*Datacenter     `json:"items"`
}

func NewDatacenterResponse() *DatacenterResponse {
        dcResp := &DatacenterResponse{}
        return dcResp
}

// NewDatacenter creates a new Datacenter object
func NewDatacenter() *Datacenter {
	dc := &Datacenter{}
        return dc
}

// ListDatacenters retreieves all Datacenters
func ListDatacenters(domainname string) ([]*Datacenter, error) {
        dcs := &DatacenterList{}
        req, err := client.NewRequest(
                Config,
                "GET",
                "/config-gtm/v1/domains/"+domainname+"/datacenters",
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
                return nil, &CommonError{entityName: "Datacenter"}
        } else {
                err = client.BodyJSON(res, dcs)
                if err != nil {
                        return nil, err
                }

                return dcs.DatacenterItems, nil
        }
}

// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
func GetDatacenter(dcId int, domainname string) (*Datacenter, error) {
	dc := NewDatacenter()
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/"+domainname+"/datacenters/"+strconv.Itoa(dcId),
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

        printHttpRequest(req, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &CommonError{entityName: "Datacenter", name: strconv.Itoa(dcId)}
	} else {
		err = client.BodyJSON(res, dc)
		if err != nil {
			return nil, err
		}

		return dc, nil
	}
}


// Create datacenter in provided domain and Datacenter object
func (dc *Datacenter) Create(domainname string) (*DatacenterResponse, error) {

        req, err := client.NewJSONRequest(
                Config,
                "POST",
                "/config-gtm/v1/domains/"+domainname+"/datacenters",
                dc,
        )
        if err != nil {
                return nil, err
        }

        SetHeader(req)

        printHttpRequest(req, true)

        res, err := client.Do(Config, req)

        printHttpResponse(res, true)

        // Network
        if err != nil {
                return nil, &CommonError{
                        entityName:       "Domain",
                        name:             domainname,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "Domain", name: domainname, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := NewDatacenterResponse()
        // Unmarshall whole response body for updated DC and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody, nil

}

func (dc *Datacenter) Update(domainname string) (*ResponseStatus, error) {

        req, err := client.NewJSONRequest(
                Config,
                "PUT",
                "/config-gtm/v1/domains/"+domainname+"/datacenters/"+strconv.Itoa(dc.DatacenterId),
                dc,
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
                        entityName:       "Datacenter",
                        name:             strconv.Itoa(dc.DatacenterId),
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "Datacenter", name: string(dc.DatacenterId), apiErrorMessage: err.Detail, err: err}
        }

        responseBody := NewDatacenterResponse()
        // Unmarshall whole response body for updated entity and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil
}

// Delete Datacenter method
func (dc *Datacenter) Delete(domainname string) (*ResponseStatus, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
                "/config-gtm/v1/domains/"+domainname+"/datacenters/"+strconv.Itoa(dc.DatacenterId),
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
                        entityName:       "Datacenter",
                        name:             strconv.Itoa(dc.DatacenterId),
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "Datacenter", name: string(dc.DatacenterId), apiErrorMessage: err.Detail, err: err}
        }

        responseBody := NewDatacenterResponse()
        // Unmarshall whole response body in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil
}

// Add/Update Datacenter element
func (dc *Datacenter) AddElement(element interface{}) error {

        // Do we need?

	return nil
}

// Remove Datacenter element
func (dc *Datacenter) RemoveElement(element interface{}) error {

        // What does this mean?

	return nil
}

// Retrieve a specific element
func (dc *Datacenter) GetElement(element interface{}) interface{} {

        // useful?
        return nil

}