package networklists


import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestListNetworkLists(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/network-list/v2/network-lists").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "networkLists": [
        {
            "networkListType": "networkListResponse",
            "accessControlGroup": "KSD\nwith ION 3-13H1234",
            "name": "General List",
            "elementCount": 3011,
            "readOnly": false,
            "shared": false,
            "syncPoint": 22,
            "type": "IP",
            "uniqueId": "25614_GENERALLIST",
            "links": {
                "activateInProduction": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST/environments/PRODUCTION/activate",
                    "method": "POST"
                },
                "activateInStaging": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST/environments/STAGING/activate",
                    "method": "POST"
                },
                "appendItems": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST",
                    "method": "POST"
                },
                "retrieve": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST"
                },
                "statusInProduction": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST/environments/PRODUCTION/status"
                },
                "statusInStaging": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST/environments/STAGING/status"
                },
                "update": {
                    "href": "/network-list/v2/network-lists/25614_GENERALLIST",
                    "method": "PUT"
                }
            }
        },
        {
            "networkListType": "networkListResponse",
            "name": "Ec2 Akamai Network List",
            "elementCount": 235,
            "readOnly": true,
            "shared": true,
            "syncPoint": 65,
            "type": "IP",
            "uniqueId": "1024_AMAZONELASTICCOMPUTECLOU",
            "links": {
                "activateInProduction": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/activate",
                    "method": "POST"
                },
                "activateInStaging": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/activate",
                    "method": "POST"
                },
                "appendItems": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/append",
                    "method": "POST"
                },
                "retrieve": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU"
                },
                "statusInProduction": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/status"
                },
                "statusInStaging": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/status"
                },
                "update": {
                    "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU",
                    "method": "PUT"
                }
            }
        },
        {
            "networkListType": "networkListResponse",
            "accessControlGroup": "KSD\nTest - 3-13H5523",
            "name": "GeoList_1913New",
            "elementCount": 16,
            "readOnly": false,
            "shared": false,
            "syncPoint": 2,
            "type": "GEO",
            "uniqueId": "26732_GEOLIST1913",
            "links": {
                "activateInProduction": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913/environments/PRODUCTION/activate",
                    "method": "POST"
                },
                "activateInStaging": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913/environments/STAGING/activate",
                    "method": "POST"
                },
                "appendItems": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913/append",
                    "method": "POST"
                },
                "retrieve": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913"
                },
                "statusInProduction": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913/environments/PRODUCTION/status"
                },
                "statusInStaging": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913/environments/STAGING/status"
                },
                "update": {
                    "href": "/network-list/v2/network-lists/26732_GEOLIST1913",
                    "method": "PUT"
                }
            }
        }
    ],
    "links": {
        "create": {
            "href": "/network-list/v2/network-lists/",
            "method": "POST"
        }
    }
}
                `)

	Init(config)

	response, err := ListNetworkLists()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ListNetworkListsResponse{}, response), true)
	assert.Equal(t, 3, len(response.NetworkLists))
}

func TestGetNetworkList(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/network-list/v2/network-lists/ABC123").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "name": "Ec2 Akamai Network List",
    "uniqueId": "1024_AMAZONELASTICCOMPUTECLOU",
    "syncPoint": 65,
    "type": "IP",
    "networkListType": "networkListResponse",
    "elementCount": 13,
    "readOnly": true,
    "shared": true,
    "list": [
        "13.125.0.0/16",
        "13.126.0.0/15",
        "13.210.0.0/15"
    ],
    "links": {
        "activateInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/activate",
            "method": "POST"
        },
        "activateInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/activate",
            "method": "POST"
        },
        "appendItems": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/append",
            "method": "POST"
        },
        "retrieve": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU"
        },
        "statusInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/status"
        },
        "statusInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/status"
        },
        "update": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU",
            "method": "PUT"
        }
    }
}
                `)

	Init(config)

	response, err := GetNetworkList("ABC123")
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &NetworkList{}, response), true)
	assert.Equal(t, 3, len(response.List))
}

func TestCreateNetworkList(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/network-list/v2/network-lists").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "name": "Ec2 Akamai Network List",
    "uniqueId": "1024_AMAZONELASTICCOMPUTECLOU",
    "syncPoint": 65,
    "type": "IP",
    "networkListType": "networkListResponse",
    "elementCount": 13,
    "readOnly": true,
    "shared": true,
    "list": [
        "13.125.0.0/16",
        "13.126.0.0/15",
        "13.210.0.0/15"
    ],
    "links": {
        "activateInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/activate",
            "method": "POST"
        },
        "activateInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/activate",
            "method": "POST"
        },
        "appendItems": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/append",
            "method": "POST"
        },
        "retrieve": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU"
        },
        "statusInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/status"
        },
        "statusInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/status"
        },
        "update": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU",
            "method": "PUT"
        }
    }
}
                `)

	Init(config)

	var iplist = make([]string, 3)
	iplist[0] = "192.168.2.1"
	iplist[0] = "192.168.2.2"
	iplist[0] = "192.168.2.3"
	var networklist NetworkList
	networklist.Name = "My New List"
	networklist.Type = "IP"
	networklist.Description  = "Description of changes"
	networklist.List = iplist
	
	response, err := CreateNetworkList(networklist)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &NetworkList{}, response), true)
	assert.Equal(t, 3, len(response.List))
}

func TestUpdateNetworkList(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Put("/network-list/v2/network-lists").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "name": "Ec2 Akamai Network List",
    "uniqueId": "1024_AMAZONELASTICCOMPUTECLOU",
    "syncPoint": 65,
    "type": "IP",
    "networkListType": "networkListResponse",
    "elementCount": 13,
    "readOnly": true,
    "shared": true,
    "list": [
        "13.125.0.0/16",
        "13.126.0.0/15",
        "13.210.0.0/15"
    ],
    "links": {
        "activateInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/activate",
            "method": "POST"
        },
        "activateInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/activate",
            "method": "POST"
        },
        "appendItems": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/append",
            "method": "POST"
        },
        "retrieve": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU"
        },
        "statusInProduction": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/PRODUCTION/status"
        },
        "statusInStaging": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU/environments/STAGING/status"
        },
        "update": {
            "href": "/network-list/v2/network-lists/1024_AMAZONELASTICCOMPUTECLOU",
            "method": "PUT"
        }
    }
}
                `)

	Init(config)

	var iplist = make([]string, 3)
	iplist[0] = "192.168.2.1"
	iplist[0] = "192.168.2.2"
	iplist[0] = "192.168.2.3"
	var networklist NetworkList
	networklist.Name = "My New List"
	networklist.Type = "IP"
	networklist.Description  = "Description of changes"
	networklist.List = iplist
	networklist.UniqueID = "1024_AMAZONELASTICCOMPUTECLOU"
	
	response, err := UpdateNetworkList(networklist)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &NetworkList{}, response), true)
	assert.Equal(t, 3, len(response.List))
}

func TestUpdateNetworkListNoUniqueID(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Put("/network-list/v2/network-lists").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
                `)

	Init(config)

	var iplist = make([]string, 3)
	iplist[0] = "192.168.2.1"
	iplist[0] = "192.168.2.2"
	iplist[0] = "192.168.2.3"
	var networklist NetworkList
	networklist.Name = "My New List"
	networklist.Type = "IP"
	networklist.Description  = "Description of changes"
	networklist.List = iplist
	
	_, err := UpdateNetworkList(networklist)
	assert.Error(t, err)
}

func TestDeleteNetworkList(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Delete("/network-list/v2/network-lists").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "status": 200,
    "uniqueId": "33501_TESTLIST",
    "syncPoint": 4
}
                `)

	Init(config)

	_, err := DeleteNetworkList("ABC123")
	assert.NoError(t, err)
}

func TestActivateNetworkList(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/network-list/v2/network-lists/ABC123/environments/STAGING/activate").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "activationId": 12345,
    "activationComments": "Whitelist IPs of new employees who joined this week",
    "activationStatus": "PENDING_ACTIVATION",
    "syncPoint": 5,
    "uniqueId": "25614_GENERALLIST",
    "fast": false,
    "dispatchCount": 1,
    "links": {
        "appendItems": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/append",
            "method": "POST"
        },
        "retrieve": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST"
        },
        "statusInProduction": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/environments/PRODUCTION/status"
        },
        "statusInStaging": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/environments/STAGING/status"
        },
        "syncPointHistory": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/sync-points/5/history"
        },
        "update": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST",
            "method": "PUT"
        },
        "activationDetails": {
            "href": "/network-list/v2/network-lists/activations/12345/"
        }
    }
}
                `)

	Init(config)

	var emails = make([]string, 1)
	emails[0] = "test@akamai.com"
	var activationrequest ActivationRequest
	activationrequest.Comments = "123test"
	activationrequest.Network = Staging
	activationrequest.UniqueID = "ABC123"
	activationrequest.NotificationRecipients = emails
	
	response, err := ActivateNetworkList(activationrequest)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ActivationStatus{}, response), true)
}

func TestGetActivationStatus(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/network-list/v2/network-lists/ABC123/environments/STAGING/status").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "activationId": 12345,
    "activationComments": "Whitelist IPs of new employees who joined this week",
    "activationStatus": "PENDING_ACTIVATION",
    "syncPoint": 5,
    "uniqueId": "25614_GENERALLIST",
    "fast": false,
    "dispatchCount": 1,
    "links": {
        "appendItems": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/append",
            "method": "POST"
        },
        "retrieve": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST"
        },
        "statusInProduction": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/environments/PRODUCTION/status"
        },
        "statusInStaging": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/environments/STAGING/status"
        },
        "syncPointHistory": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST/sync-points/5/history"
        },
        "update": {
            "href": "/networklist-api/rest/v2/network-lists/25614_GENERALLIST",
            "method": "PUT"
        },
        "activationDetails": {
            "href": "/network-list/v2/network-lists/activations/12345/"
        }
    }
}
                `)

	Init(config)

	response, err := GetActivationStatus("ABC123", Staging)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ActivationStatus{}, response), true)
}
