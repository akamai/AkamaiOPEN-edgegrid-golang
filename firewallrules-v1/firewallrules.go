package firewallrules

import (
	"bytes"
	"encoding/json"
	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type ListServicesResponse []struct {
	ServiceID   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
}

type Subscription struct {
	ServiceID   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Email       string `json:"email"`
	SignupDate  string `json:"signupDate"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type ListCidrBlocksResponse []struct {
	CidrID        int    `json:"cidrId"`
	ServiceID     int    `json:"serviceId"`
	ServiceName   string `json:"serviceName"`
	Cidr          string `json:"cidr"`
	CidrMask      string `json:"cidrMask"`
	Port          string `json:"port"`
	CreationDate  string `json:"creationDate"`
	EffectiveDate string `json:"effectiveDate"`
	ChangeDate    string `json:"changeDate"`
	MinIP         string `json:"minIp"`
	MaxIP         string `json:"maxIp"`
	LastAction    string `json:"lastAction"`
}

type UpdateSubscriptionsRequest struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type UpdateSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

func ListSubscriptions() (*ListSubscriptionsResponse, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		"/firewall-rules-manager/v1/subscriptions",
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

	var response ListSubscriptionsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil

}
func ListServices() (*ListServicesResponse, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		"/firewall-rules-manager/v1/services",
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

	var response ListServicesResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil

}

func UpdateSubscriptions(updatesubscriptionsrequest UpdateSubscriptionsRequest) (*UpdateSubscriptionsResponse, error) {

	s, err := json.Marshal(updatesubscriptionsrequest)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PUT",
		"/firewall-rules-manager/v1/subscriptions",
		bytes.NewReader(s),
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

	var response UpdateSubscriptionsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil

}
func ListCidrBlocks() (*ListCidrBlocksResponse, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		"/firewall-rules-manager/v1/cidr-blocks",
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

	var response ListCidrBlocksResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil

}
