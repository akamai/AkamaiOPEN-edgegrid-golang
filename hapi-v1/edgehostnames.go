package hapi

import (
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Certificate struct {
	CertificateID    string    `json:"certificateId"`
	CommonName       string    `json:"commonName"`
	SerialNumber     string    `json:"serialNumber"`
	SlotNumber       int       `json:"slotNumber"`
	ExpirationDate   time.Time `json:"expirationDate"`
	CertificateType  string    `json:"certificateType"`
	ValidationType   string    `json:"validationType"`
	Status           string    `json:"status"`
	AvailableDomains []string  `json:"availableDomains"`
}

type EdgeHostname struct {
	EdgeHostnameID    int    `json:"edgeHostnameId"`
	RecordName        string `json:"recordName"`
	DNSZone           string `json:"dnsZone"`
	SecurityType      string `json:"securityType"`
	UseDefaultTTL     bool   `json:"useDefaultTtl"`
	UseDefaultMap     bool   `json:"useDefaultMap"`
	IPVersionBehavior string `json:"ipVersionBehavior"`
	ProductID         string `json:"productId"`
	TTL               int    `json:"ttl"`
	Map               string `json:"map,omitempty"`
	SlotNumber        int    `json:"slotNumber,omitempty"`
	Comments          string `json:"comments"`
	SerialNumber      int    `json:"serialNumber,omitempty"`
	CustomTarget      string `json:"customTarget,omitempty"`
	ChinaCdn          struct {
		IsChinaCdn bool `json:"isChinaCdn"`
	} `json:"chinaCdn,omitempty"`
	IsEdgeIPBindingEnabled bool `json:"isEdgeIPBindingEnabled,omitempty"`
} 

type ListEdgeHostnamesResponse struct {
	EdgeHostnames []EdgeHostname `json:"edgeHostnames"`
}

type ChangeRequest struct {
	Action            string    `json:"action"`
	ChangeID          int       `json:"changeId"`
	Comments          string    `json:"comments"`
	Status            string    `json:"status"`
	StatusMessage     string    `json:"statusMessage"`
	StatusUpdateDate  time.Time `json:"statusUpdateDate"`
	StatusUpdateEmail string    `json:"statusUpdateEmail"`
	SubmitDate        time.Time `json:"submitDate"`
	Submitter         string    `json:"submitter"`
	SubmitterEmail    string    `json:"submitterEmail"`
	EdgeHostnames     []EdgeHostname `json:"edgeHostnames"`
}

type ListChangeRequestsResponse struct {
	ChangeRequests []ChangeRequest `json:"changeRequests"`
}

type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type Products struct {
	ProductDisplayNames []struct {
		ProductID          string `json:"productId"`
		ProductDisplayName string `json:"productDisplayName"`
	} `json:"productDisplayNames"`
}

func ListEdgeHostnames() (*ListEdgeHostnamesResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/edge-hostnames",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListEdgeHostnamesResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func GetEdgeHostname(recordName string, dnsZone string) (*EdgeHostname, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response EdgeHostname
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func GetCertificate(recordName string, dnsZone string) (*Certificate, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s/certificate", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Certificate
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func GetEdgeHostnameById(id string) (*EdgeHostname, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/edge-hostnames/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response EdgeHostname
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func DeleteEdgeHostname(recordName string, dnsZone string) (*ChangeRequest, error) {
	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func PatchEdgeHostname(recordName string, dnsZone string, patch []Patch) (*ChangeRequest, error) {

        r, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PATCH",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		bytes.NewReader(r),
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func ListChangeRequestsByHostname(recordName string, dnsZone string) (*ListChangeRequestsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s/change-requests", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListChangeRequestsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func ListChangeRequests() (*ListChangeRequestsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/change-requests",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListChangeRequestsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func GetChangeRequest(id string) (*ChangeRequest, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/change-requests/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 

func ListProducts() (*Products, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/products/display-names",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Products
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
} 
