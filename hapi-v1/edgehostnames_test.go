package hapi

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)


func TestListEdgehostnames (t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/edge-hostnames").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeHostnames": [
        {
            "edgeHostnameId": 1,
            "recordName": "www.example.com",
            "dnsZone": "edgekey.net",
            "securityType": "ENHANCED-TLS",
            "useDefaultTtl": true,
            "useDefaultMap": true,
            "ipVersionBehavior": "IPV6_IPV4_DUALSTACK",
            "productId": "DSD",
            "ttl": 21600,
            "map": "e;dscb.akamaiedge.net",
            "slotNumber": 11838,
            "comments": "Created for example.com"
        },
        {
            "edgeHostnameId": 2,
            "recordName": "www.examples.com",
            "dnsZone": "edgesuite.net",
            "securityType": "STANDARD-TLS",
            "useDefaultTtl": true,
            "useDefaultMap": true,
            "ipVersionBehavior": "IPV4",
            "productId": "DSA",
            "ttl": 21600,
            "map": "a;g.akamai.net",
            "serialNumber": 140,
            "comments": "Created for site"
        },
        {
            "edgeHostnameId": 3,
            "recordName": "www.example-china.com",
            "dnsZone": "edgesuite.net",
            "useDefaultTtl": true,
            "useDefaultMap": true,
            "securityType": "STANDARD-TLS",
            "ipVersionBehavior": "IPV4",
            "productId": "DSA",
            "ttl": 1800,
            "map": "a;g.akamai.net",
            "serialNumber": 2322,
            "customTarget": "www.example-china.com.edgesuite.net.globalredir.akadns.net",
            "comments": "Created for China CDN",
            "chinaCdn": {
                "isChinaCdn": true
            }
        },
        {
            "edgeHostnameId": 4,
            "recordName": "www.example-eip.com",
            "dnsZone": "edgekey.net",
            "useDefaultTtl": true,
            "useDefaultMap": true,
            "securityType": "ENHANCED-TLS",
            "ipVersionBehavior": "IPV4",
            "productId": "DSA",
            "ttl": 21600,
            "isEdgeIPBindingEnabled": true,
            "slotNumber": 2322,
            "customTarget": "www.example-eip.com.eip.akadns.net",
            "comments": "Created for EIP"
        }
    ]

}
                `)


	Init(config)

	response, err := ListEdgeHostnames()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListEdgeHostnamesResponse{}, response), true)
	assert.Equal(t, 4, len(response.EdgeHostnames))
}

func TestGetEdgeHostname(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.examples.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeHostnameId": 1,
    "recordName": "www.example.com",
    "dnsZone": "edgekey.net",
    "securityType": "ENHANCED-TLS",
    "useDefaultTtl": true,
    "useDefaultMap": true,
    "ipVersionBehavior": "IPV6_IPV4_DUALSTACK",
    "productId": "DSD",
    "ttl": 21600,
    "map": "e;dscb.akamaiedge.net",
    "slotNumber": 11838,
    "comments": "Created for example.com"
}
                `)

	Init(config)

	response, err := GetEdgeHostname("www.examples.com", "edgesuite.net") 
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &EdgeHostname{}, response), true)
	assert.Equal(t, 11838, response.SlotNumber)
}

func TestGetCertificate(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.examples.com/certificate").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "certificateId": "1234",
    "commonName": "example.com",
    "serialNumber": "12:34:56:78:90:AB:CD:EF",
    "slotNumber": 8927,
    "expirationDate": "2019-10-31T23:59:59Z",
    "certificateType": "SAN",
    "validationType": "DOMAIN_VALIDATION",
    "status": "PENDING",
    "availableDomains": [
        "live.example.com",
        "secure.example.com",
        "www.example.com"
    ]
}
                `)

	Init(config)

	response, err := GetCertificate("www.examples.com", "edgesuite.net") 
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Certificate{}, response), true)
	assert.Equal(t, 8927, response.SlotNumber)
}

func TestGetEdgeHostnameById(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/edge-hostnames/3").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeHostnameId": 1,
    "recordName": "www.example.com",
    "dnsZone": "edgekey.net",
    "securityType": "ENHANCED-TLS",
    "useDefaultTtl": true,
    "useDefaultMap": true,
    "ipVersionBehavior": "IPV6_IPV4_DUALSTACK",
    "productId": "DSD",
    "ttl": 21600,
    "map": "e;dscb.akamaiedge.net",
    "slotNumber": 11838,
    "comments": "Created for example.com"
}
                `)

	Init(config)

	response, err := GetEdgeHostnameById("3") 
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &EdgeHostname{}, response), true)
	assert.Equal(t, 11838, response.SlotNumber)
}

func TestDeleteEdgeHostname(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
                Delete("/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.examples.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "action": "DELETE",
    "changeId": 1,
    "comments": "Delete for example.com",
    "status": "SUCCEEDED",
    "statusMessage": "File successfully deployed to Akamai's network",
    "statusUpdateDate": "2018-10-17T20:21:17.0Z",
    "statusUpdateEmail": "nobody@akamai.com",
    "submitDate": "2018-10-17T20:21:12.0Z",
    "submitter": "nobody",
    "submitterEmail": "nobody@akamai.com",
    "edgeHostnames": [
        {
            "edgeHostnameId": 1,
            "recordName": "www.example.com",
            "dnsZone": "edgekey.net",
            "securityType": "ENHANCED-TLS",
            "useDefaultTtl": false,
            "useDefaultMap": false,
            "ttl": 21600,
            "map": "e;b.akamaiedge.net",
            "slotNumber": 11838,
            "ipVersionBehavior": "IPV4",
            "comments": "Edited for Super Bowl"
        }
    ]
}
                `)

	Init(config)

	response, err := DeleteEdgeHostname("www.examples.com", "edgesuite.net")
	assert.NoError(t, err)
        assert.Equal(t, assert.IsType(t, &ChangeRequest{}, response), true)
	assert.Equal(t, 1, response.ChangeID)
}

func TestPatchEdgeHostname(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
                Patch("/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.examples.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "action": "EDIT",
    "changeId": 53494819,
    "edgeHostnames": [
        {
            "edgeHostnameId": 1803728,
            "recordName": "amp.nfl.com",
            "dnsZone": "edgekey.net",
            "securityType": "ENHANCED-TLS",
            "useDefaultMap": false,
            "useDefaultTtl": false,
            "ttl": 21600,
            "map": "e;b.akamaiedge.net",
            "slotNumber": 91899,
            "comments": "Edited for Super Bowl"
        }
    ],
    "comments": "Editing CNAME",
    "status": "PENDING",
    "statusUpdateDate": "2016-11-18T20:21:17.0Z",
    "statusUpdateEmail": "bbuenave@akamai.com",
    "submitDate": "2016-11-18T20:21:12.0Z",
    "submitter": "bbuenave",
    "submitterEmail": "bbuenave@akamai.com"
}
                `)

	Init(config)

	var patches = make([]Patch, 2)
	var patch Patch
	patch.Op = "replace"
	patch.Path = "/ttl"
	patch.Value = "500"

	var patch2 Patch
	patch2.Op = "replace"
	patch2.Path = "/ipVersionBehavior"
	patch2.Value = "IPV4"

	response, err := PatchEdgeHostname("www.examples.com", "edgesuite.net", patches)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ChangeRequest{}, response), true)
}

func TestListChangeRequestsByHostname (t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/dns-zones/edgesuite.net/edge-hostnames/www.examples.com/change-requests").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "changeRequests": [
        {
            "action": "EDIT",
            "changeId": 2,
            "comments": "Editing CNAME",
            "status": "PENDING",
            "statusMessage": "File uploaded and awaiting validation",
            "statusUpdateDate": "2018-10-17T20:21:17.0Z",
            "statusUpdateEmail": "nobody@akamai.com",
            "submitDate": "2018-10-17T20:21:12.0Z",
            "submitter": "nobody",
            "submitterEmail": "nobody@akamai.com",
            "edgeHostnames": [
                {
                    "edgeHostnameId": 2,
                    "recordName": "www.examples.com",
                    "dnsZone": "edgesuite.net",
                    "securityType": "STANDARD-TLS",
                    "useDefaultTtl": true,
                    "useDefaultMap": true,
                    "ipVersionBehavior": "IPV4",
                    "productId": "DSA",
                    "ttl": 21600,
                    "map": "a;g.akamai.net",
                    "serialNumber": 140,
                    "comments": "Created for site"
                }
            ]
        },
        {
            "action": "CREATE",
            "changeId": 1,
            "comments": "Created for example.com",
            "status": "PENDING",
            "statusMessage": "File uploaded and awaiting validation",
            "statusUpdateDate": "2018-10-17T20:21:17.0Z",
            "statusUpdateEmail": "nobody@akamai.com",
            "submitDate": "2018-10-17T20:21:12.0Z",
            "submitter": "nobody",
            "submitterEmail": "nobody@akamai.com",
            "edgeHostnames": [
                {
                    "edgeHostnameId": 1,
                    "recordName": "www.example.com",
                    "dnsZone": "edgekey.net",
                    "securityType": "ENHANCED-TLS",
                    "useDefaultTtl": false,
                    "useDefaultMap": false,
                    "ttl": 21600,
                    "map": "e;b.akamaiedge.net",
                    "slotNumber": 11838,
                    "ipVersionBehavior": "IPV4",
                    "comments": "Edited for Super Bowl"
                }
            ]
        }
    ]
}
                `)


	Init(config)

	response, err := ListChangeRequestsByHostname("www.examples.com", "edgesuite.net") 
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListChangeRequestsResponse{}, response), true)
	assert.Equal(t, 2, len(response.ChangeRequests))
}

func TestListChangeRequests(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/change-requests").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "changeRequests": [
        {
            "action": "EDIT",
            "changeId": 2,
            "comments": "Editing CNAME",
            "status": "PENDING",
            "statusMessage": "File uploaded and awaiting validation",
            "statusUpdateDate": "2018-10-17T20:21:17.0Z",
            "statusUpdateEmail": "nobody@akamai.com",
            "submitDate": "2018-10-17T20:21:12.0Z",
            "submitter": "nobody",
            "submitterEmail": "nobody@akamai.com",
            "edgeHostnames": [
                {
                    "edgeHostnameId": 2,
                    "recordName": "www.examples.com",
                    "dnsZone": "edgesuite.net",
                    "securityType": "STANDARD-TLS",
                    "useDefaultTtl": true,
                    "useDefaultMap": true,
                    "ipVersionBehavior": "IPV4",
                    "productId": "DSA",
                    "ttl": 21600,
                    "map": "a;g.akamai.net",
                    "serialNumber": 140,
                    "comments": "Created for site"
                }
            ]
        },
        {
            "action": "CREATE",
            "changeId": 1,
            "comments": "Created for example.com",
            "status": "PENDING",
            "statusMessage": "File uploaded and awaiting validation",
            "statusUpdateDate": "2018-10-17T20:21:17.0Z",
            "statusUpdateEmail": "nobody@akamai.com",
            "submitDate": "2018-10-17T20:21:12.0Z",
            "submitter": "nobody",
            "submitterEmail": "nobody@akamai.com",
            "edgeHostnames": [
                {
                    "edgeHostnameId": 1,
                    "recordName": "www.example.com",
                    "dnsZone": "edgekey.net",
                    "securityType": "ENHANCED-TLS",
                    "useDefaultTtl": false,
                    "useDefaultMap": false,
                    "ttl": 21600,
                    "map": "e;b.akamaiedge.net",
                    "slotNumber": 11838,
                    "ipVersionBehavior": "IPV4",
                    "comments": "Edited for Super Bowl"
                }
            ]
        }
    ]
}
                `)


	Init(config)

	response, err := ListChangeRequests()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListChangeRequestsResponse{}, response), true)
	assert.Equal(t, 2, len(response.ChangeRequests))
}

func TestGetChangeRequest(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/change-requests/1").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "action": "CREATE",
    "changeId": 1,
    "comments": "Created for example.com",
    "status": "SUCCEEDED",
    "statusMessage": "File successfully deployed to Akamai's network",
    "statusUpdateDate": "2018-10-17T20:21:17.0Z",
    "statusUpdateEmail": "nobody@akamai.com",
    "submitDate": "2018-10-17T20:21:12.0Z",
    "submitter": "nobody",
    "submitterEmail": "nobody@akamai.com",
    "edgeHostnames": [
        {
            "edgeHostnameId": 1,
            "recordName": "www.example.com",
            "dnsZone": "edgekey.net",
            "securityType": "ENHANCED-TLS",
            "useDefaultTtl": false,
            "useDefaultMap": false,
            "ttl": 21600,
            "map": "e;b.akamaiedge.net",
            "slotNumber": 11838,
            "ipVersionBehavior": "IPV4",
            "comments": "Edited for Super Bowl"
        }
    ]
}
                `)


	Init(config)

	response, err := GetChangeRequest("1")
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ChangeRequest{}, response), true)
	assert.Equal(t, 1, response.ChangeID)
}

func TestListProducts(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/hapi/v1/products/display-names").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "productDisplayNames": [
        {
            "productId": "Alta",
            "productDisplayName": "Protect & Perform"
        },
        {
            "productId": "DSD",
            "productDisplayName": "Site Delivery"
        }
    ]
}
                `)


	Init(config)

	response, err := ListProducts()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Products{}, response), true)
	assert.Equal(t, 2, len(response.ProductDisplayNames))
}
