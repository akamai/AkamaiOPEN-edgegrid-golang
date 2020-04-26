package firewallrules

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestListSubscriptions(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/firewall-rules-manager/v1/subscriptions"))
	mock.
		Get("/firewall-rules-manager/v1/subscriptions").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "subscriptions": [
        {
            "description": "Edge Staging Network",
            "email": "test@akamai.com",
            "serviceId": 7,
            "serviceName": "ESN",
            "signupDate": "2020-04-24"
        },
        {
            "description": "Development Test IPs",
            "email": "test@akamai.com",
            "serviceId": 13,
            "serviceName": "Test IPs",
            "signupDate": "2020-04-24"
        }
    ]
}
                `)

	Init(config)

	response, err := ListSubscriptions()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListSubscriptionsResponse{}, response), true)
	assert.Equal(t, len(response.Subscriptions), 2)
}

func TestListServices(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/firewall-rules-manager/v1/services"))
	mock.
		Get("/firewall-rules-manager/v1/services").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
    [
        {
            "serviceId": 7,
            "serviceName": "ESN",
            "description": "Edge Staging Network"
        },
        {
            "serviceId": 13,
            "serviceName": "Test IPs",
            "description": "Development Test IPs"
        }
    ]
                `)

	Init(config)

	response, err := ListServices()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListServicesResponse{}, response), true)
	assert.Equal(t, len(*response), 2)
}

func TestListCidrBlocks(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/firewall-rules-manager/v1/cidr-blocks"))
	mock.
		Get("/firewall-rules-manager/v1/cidr-blocks").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
[
    {
        "changeDate": null,
        "cidr": "209.170.113.98",
        "cidrId": 303,
        "cidrMask": "/31",
        "creationDate": "2007-09-27",
        "description": "Secure Edge Staging Network",
        "effectiveDate": "2007-10-13",
        "lastAction": "update",
        "maxIp": "209.170.113.99",
        "minIp": "209.170.113.98",
        "port": "80,443",
        "serviceId": 8,
        "serviceName": "SESN"
    },
    {
        "changeDate": null,
        "cidr": "209.170.113.100",
        "cidrId": 304,
        "cidrMask": "/31",
        "creationDate": "2007-09-27",
        "description": "Secure Edge Staging Network",
        "effectiveDate": "2007-10-13",
        "lastAction": "update",
        "maxIp": "209.170.113.101",
        "minIp": "209.170.113.100",
        "port": "80,443",
        "serviceId": 8,
        "serviceName": "SESN"
    }
]

                `)

	Init(config)

	response, err := ListCidrBlocks()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListCidrBlocksResponse{}, response), true)
	assert.Equal(t, len(*response), 2)
}

func TestUpdateSubscriptions(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/firewall-rules-manager/v1/subscriptions"))
	mock.
		Put("/firewall-rules-manager/v1/subscriptions").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "subscriptions": [
        {
            "email": "test@akamai.com",
            "serviceId": 10,
            "serviceName": "NETSTORAGE",
            "signupDate": "2020-04-24"
        }
    ]
}
                `)

	Init(config)

	// Create a new subscription
	var subscription Subscription
	subscription.ServiceID = 10
	subscription.ServiceName = "NETSTORAGE"
	subscription.Email = "test@akamai.com"
	subscription.SignupDate = "2020-04-24"

	// Add it to the list of subscriptions
	var subscriptions = make([]Subscription, 0)
	subscriptions = append(subscriptions, subscription)

	// Wrap it in a request
	var request UpdateSubscriptionsRequest
	request.Subscriptions = subscriptions

	response, err := UpdateSubscriptions(request)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &UpdateSubscriptionsResponse{}, response), true)
	assert.Equal(t, len(response.Subscriptions), 1)
}
