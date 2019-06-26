package configgtm

import (
	"github.com/edglynes/AkamaiOPEN-edgegrid-golang/client-v1"

        "fmt"
)

//
// Handle Operations on gtm resources
// Based on 1.3 schema
//

type ResourceInstance struct {
        DatacenterId             int               `json:"datacenterId"`
        UseDefaultLoadObject     bool              `json:"seDefaultLoadObject,omitempty"`
        LoadObject
}

// Resource represents a GTM resource   
type Resource struct {
        Type                              string                 `json:"type"`
        HostHeader                        string                 `json:"hostHeader,omitempty"`
        LeastSquaresDecay                 int                    `json:"leastSquaresDecay,omitempty"`
        Description                       string                 `json:"description,omitempty"`
        LeaderString                      string                 `json:"leaderString,omitempty"`
        ConstrainedProperty               string                 `json:"constrainedProperty,omitempty"`
        ResourceInstances                 []*ResourceInstance    `json:"resourceInstances,omitempty"`
        AggregationType                   string                 `json:"aggregationType"`
        Links                             []*Link                `json:"links,omitempty"`
        LoadImbalancePercent              float64                `json:"loadImbalancePercent,omitempty"`
        UpperBound                        int                    `json:"upperBound,omitempty"`
        Name                              string                 `json:"name"`
        MaxUMultiplicativeIncrement       float64                `json:"maxUMultiplicativeIncrement,omitempty"`
        DecayRate                         float64                `json:"decayRate,omitempty"`
}

type ResourceList struct {
        ResourceItems            []*Resource       `json:"items"`
}

// NewResourceInstance instantiates a new ResourceInstance. 
func (rsrc *Resource) NewResourceInstance(dcid int) *ResourceInstance {
 
        return &ResourceInstance{DatacenterId: dcid}  

}

// NewResource creates a new Resource object.
func NewResource(resourcename string) *Resource {
	resource := &Resource{Name: resourcename}
        return resource
}

// ListResources retreieves all Resources in the specified domain.
func ListResources(domainname string) ([]*Resource, error) {
        rsrcs := &ResourceList{}
        req, err := client.NewRequest(
                Config,
                "GET",
                "/config-gtm/v1/domains/"+domainname+"/resources",
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
                return nil, &CommonError{entityName: "Resources"}
        } 
        err = client.BodyJSON(res, rsrcs)
        if err != nil {
                return nil, err
        }

        return rsrcs.ResourceItems, nil
         
}

// GetResource retrieves a Resource with the given name in the specified domain.
func GetResource(resourcename, domainname string) (*Resource, error) {
	rsc := NewResource(resourcename)
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/"+domainname+"/resources/"+resourcename,
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
		return nil, &CommonError{entityName: "Resource", name: resourcename}
	} else {
		err = client.BodyJSON(res, rsc)
		if err != nil {
			return nil, err
		}

		return rsc, nil
	}
}


// Create the resource identified by the receiver argument in the specified domain.                        
func (rsrc *Resource) Create(domainname string) (*ResourceResponse, error) {

        // Use common code. Any specific validation needed?

        return rsrc.save(domainname)

}

// Update the resourceidentified in the receiver argument in the specified domain.
func (rsrc *Resource) Update(domainname string) (*ResponseStatus, error) {

        // common code
   
        stat, err := rsrc.save(domainname)
        if err != nil {
                return nil, err
        }
        return stat.Status, err

}

// Save Resource in given domain. Common path for Create and Update.
func (rsrc *Resource) save(domainname string) (*ResourceResponse, error) {

        fmt.Println("Creating request!!!!")
        req, err := client.NewJSONRequest(
                Config,
                "PUT",
                "/config-gtm/v1/domains/"+domainname+"/resources/"+rsrc.Name,
                rsrc,
        )
        fmt.Println("BACK!!!!!!!!!!!!!!!!")
      
        if err != nil {
                return nil, err
        }

        fmt.Println("Setting header!!!")

        SetHeader(req)

        printHttpRequest(req, true)

        res, err := client.Do(Config, req)

        printHttpResponse(res, true)

        // Network error
        if err != nil {
                return nil, &CommonError{
                        entityName:       "Resource",
                        name:             rsrc.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "Resource", name: rsrc.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &ResourceResponse{}
        // Unmarshall whole response body for updated entity and in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody, nil

}

// Delete the resource identified in the receiver argument from the specified domain. 
func (rsrc *Resource) Delete(domainname string) (*ResponseStatus, error) {

        req, err := client.NewRequest(
                Config,
                "DELETE",
                "/config-gtm/v1/domains/"+domainname+"/resources/"+rsrc.Name,
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
                        entityName:       "Resource",
                        name:             rsrc.Name,
                        httpErrorMessage: err.Error(),
                        err:              err,
                }
        }

        // API error
        if client.IsError(res) {
                err := client.NewAPIError(res)
                return nil, &CommonError{entityName: "Resource", name:rsrc.Name, apiErrorMessage: err.Detail, err: err}
        }

        responseBody := &ResponseBody{}
        // Unmarshall whole response body in case want status
        err = client.BodyJSON(res, responseBody)
        if err != nil {
                return nil, err
        }

        return responseBody.Status, nil

}

