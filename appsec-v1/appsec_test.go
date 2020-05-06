package appsec

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestListConfigurationVersions(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/appsec/v1/configs/123/versions").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "totalSize": 3,
    "pageSize": 3,
    "page": 1,
    "configId": 8277,
    "configName": "TestConfig",
    "stagingExpediteRequestId": 5861,
    "productionExpediteRequestId": 6951,
    "productionActiveVersion": 9,
    "stagingActiveVersion": 8,
    "lastCreatedVersion": 9,
    "versionList": [
        {
            "version": 9,
            "versionNotes": "Membership Benefits",
            "createDate": "2013-10-07T17:58:52Z",
            "createdBy": "user1",
            "basedOn": 8,
            "production": {
                "status": "Active",
                "time": "2014-07-08T07:40:00Z"
            },
            "staging": {
                "status": "Inactive"
            }
        },
        {
            "version": 8,
            "versionNotes": "Membership Benefits",
            "createDate": "2013-10-07T17:41:52Z",
            "createdBy": "user2",
            "basedOn": 7,
            "production": {
                "status": "Inactive"
            },
            "staging": {
                "status": "Active",
                "time": "2014-07-08T07:40:00Z"
            }
        },
        {
            "version": 7,
            "versionNotes": "Membership Benefits",
            "createDate": "2013-08-07T17:41:52Z",
            "createdBy": "user3",
            "production": {
                "status": "Inactive"
            },
            "staging": {
                "status": "Inactive"
            }
        }
    ]
}
                `)

	Init(config)

	response, err := ListConfigurationVersions(123)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &VersionList{}, response), true)
	assert.Equal(t, 9, response.LastCreatedVersion)
}

func TestCloneConfigurationVersion(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/appsec/v1/configs/123/versions").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "configId": 8277,
    "configName": "TestConfig",
    "version": 124,
    "versionNotes": "Membership Benefits",
    "createDate": "2013-10-07T17:58:52Z",
    "createdBy": "user1",
    "basedOn": 123,
    "production": {
        "status": "Active",
        "time": "2014-07-08T07:40:00Z"
    },
    "staging": {
        "status": "Inactive"
    }
}
                `)

	Init(config)

	var configurationclone ConfigurationClone
	configurationclone.CreateFromVersion = 123
	configurationclone.RuleUpdate = false

	response, err := CloneConfigurationVersion(123, configurationclone)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Version{}, response), true)
	assert.Equal(t, 124, response.Version)
	assert.Equal(t, 123, response.BasedOn)
}

func TestActivateConfigurationVersion(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/appsec/v1/activations").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "dispatchCount": 1,
    "activationId": 1234,
    "action": "ACTIVATE",
    "status": "RECEIVED",
    "network": "PRODUCTION",
    "estimate": "PTM5",
    "createdBy": "user1",
    "createDate": "2013-10-07T17:41:52+00:00",
    "activationConfigs": [
        {
            "configId": 1,
            "configName": "config 1",
            "configVersion": 4,
            "previousConfigVersion": 2
        }
    ]
}
                `)

	Init(config)

	var activation Activation
	activation.Action = "ACTIVATE"
	activation.Network = "STAGING"
	activation.Note = "note"

	response, err := ActivateConfigurationVersion(activation)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ActivationResponse{}, response), true)
	assert.Equal(t, 200, response.ResponseCode)
}
