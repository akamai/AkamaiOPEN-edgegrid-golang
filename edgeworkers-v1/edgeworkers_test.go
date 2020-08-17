package edgeworkers

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)


func TestListPermissionsGroups(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/groups").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "groups": [
        {
            "groupId": 4284,
            "groupName": "Group #42",
            "capabilities": [
                "VIEW",
                "VIEW_VERSION",
                "VIEW_ACTIVATION",
                "ACTIVATE"
            ]
        },
        {
            "groupId": 109795,
            "groupName": "ESSL Behavior Tests",
            "capabilities": [
                "VIEW",
                "EDIT",
                "VIEW_VERSION",
                "CREATE_VERSION",
                "VIEW_ACTIVATION",
                "ACTIVATE"
            ]
        }
    ]
}
                `)

	Init(config)

	response, err := ListPermissionsGroups()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &PermissionsList{}, response), true)
	assert.Equal(t, 2, len(response.Groups))
}

func TestGetPermissionsGroup(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/groups/109795").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "groupId": 109795,
    "groupName": "ESSL Behavior Tests",
    "capabilities": [
        "VIEW",
        "VIEW_VERSION",
        "VIEW_ACTIVATION"
    ]
}
                `)

	Init(config)

	response, err := GetPermissionsGroup(109795)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Permissions{}, response), true)
	assert.Equal(t, 109795, response.GroupID)
}

func TestListEdgeworkerIDs(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/ids").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeWorkerIds": [
        {
            "edgeWorkerId": 1,
            "name": "Ew_1",
            "accountId": "B-M-1KQK3WU",
            "groupId": 72297,
            "createdBy": "jsmith",
            "createdTime": "2018-10-15T14:49:40Z",
            "lastModifiedBy": "jsmith",
            "lastModifiedTime": "2018-10-15T15:21:15Z"
        },
        {
            "edgeWorkerId": 2,
            "name": "EdgeWorker #2",
            "accountId": "B-M-1KQK3WU",
            "groupId": 72297,
            "createdBy": "jsmith",
            "createdTime": "2018-10-15T16:54:40Z",
            "lastModifiedBy": "jsmith",
            "lastModifiedTime": "2018-10-15T16:54:40Z"
        }
    ]
}
                `)

	Init(config)

	response, err := ListEdgeWorkerIDs()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &EdgeWorkerList{}, response), true)
	assert.Equal(t, 2, len(response.EdgeWorkerIDs))
}

func TestCreateEdgeWorker(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/edgeworkers/v1/ids").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeWorkerId": 42,
    "name": "Ew_42",
    "accountId": "B-M-1KQK3WU",
    "groupId": 72297,
    "createdBy": "jsmith",
    "createdTime": "2018-10-15T14:49:40Z",
    "lastModifiedBy": "jsmith",
    "lastModifiedTime": "2018-10-15T15:21:15Z"
}
                `)

	Init(config)

	var edgeworker EdgeWorker
	edgeworker.GroupID = 72297
	edgeworker.Name = "Ew_42"

	response, err := CreateEdgeWorker(edgeworker)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &EdgeWorker{}, response), true)
}

func TestValidateEdgeWorker(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/edgeworkers/v1/validations").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "errors": [
        {
            "type": "STATIC_VALIDATION_FAILED",
            "message": "static validation failed : main.js::9:16 SyntaxError: Unexpected identifier."
        },
        {
            "type": "INVALID_MANIFEST",
            "message": "manifest file is invalid"
        }
    ]
}
                `)

	Init(config)

	var r []byte

	response, err := ValidateEdgeWorkerBundle(r)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Validation{}, response), true)
}

func TestCreateVersion(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/edgeworkers/v1/ids/42/versions").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeWorkerId": 42,
    "version": 5,
    "accountId": "B-M-1KQK3WU",
    "checksum": "de9f2c7fd25e1b3afad3e85a0bd17d9b100db4b3",
    "createdBy": "jsmith",
    "createdTime": "2018-07-05T18:17:46Z"
}
                `)

	Init(config)

	var r []byte

	response, err := CreateVersion(42, r)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Version{}, response), true)
}

func TestGetVersion(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/ids/42/versions/2.3.0").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeWorkerId": 42,
    "version": 3,
    "accountId": "B-M-1KQK3WU",
    "checksum": "de9f2c7fd25e1b3afad3e85a0bd17d9b100db4b3",
    "createdBy": "jsmith",
    "createdTime": "2018-07-09T21:51:07Z"
}
                `)

	Init(config)

	response, err := GetVersion(42, "2.3.0")
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Version{}, response), true)
}

func TestGetBundle(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/ids/42/versions/2.3.0/content").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/gzip").
		BodyString(`
                `)

	Init(config)

	response, err := GetBundle(42, "2.3.0")
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, []byte{}, response), true)
}

func TestListActivations(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/ids/42/activations").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "activations": [
        {
            "edgeWorkerId": 42,
            "version": 2,
            "activationId": 3,
            "accountId": "B-M-1KQK3WU",
            "status": "PENDING",
            "network": "PRODUCTION",
            "createdBy": "jdoe",
            "createdTime": "2018-07-09T09:03:28Z",
            "lastModifiedTime": "2018-07-09T09:04:42Z"
        },
        {
            "edgeWorkerId": 42,
            "version": 1,
            "activationId": 1,
            "accountId": "B-M-1KQK3WU",
            "status": "IN_PROGRESS",
            "network": "STAGING",
            "createdBy": "jsmith",
            "createdTime": "2018-07-09T08:13:54Z",
            "lastModifiedTime": "2018-07-09T08:35:02Z"
        },
        {
            "edgeWorkerId": 42,
            "version": 2,
            "activationId": 2,
            "accountId": "B-M-1KQK3WU",
            "status": "COMPLETE",
            "network": "PRODUCTION",
            "createdBy": "asmith",
            "createdTime": "2018-07-10T14:23:42Z",
            "lastModifiedTime": "2018-07-10T14:53:25Z"
        }
    ]
}
	`)
	Init(config)

	response, err := ListActivations(42)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ActivationList{}, response), true)
	assert.Equal(t, 3, len(response.Activations))
}

func TestCreateActivation(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post("/edgeworkers/v1/ids/42/activations").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
    "edgeWorkerId": 42,
    "version": 5,
    "activationId": 3,
    "accountId": "B-M-1KQK3WU",
    "status": "PRESUBMIT",
    "network": "PRODUCTION",
    "createdBy": "jsmith",
    "createdTime": "2019-04-05T18:17:46Z",
    "lastModifiedTime": "2019-04-05T18:17:46Z"
}
                `)

	Init(config)

	var activation Activation
	activation.Network = "PRODUCTION"
	activation.Version = 2

	response, err := CreateActivation(42, activation)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Activation{}, response), true)
}

func TestGetActivation(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get("/edgeworkers/v1/ids/42/activations/3").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
        {
            "edgeWorkerId": 42,
            "version": 2,
            "activationId": 3,
            "accountId": "B-M-1KQK3WU",
            "status": "PENDING",
            "network": "PRODUCTION",
            "createdBy": "jdoe",
            "createdTime": "2018-07-09T09:03:28Z",
            "lastModifiedTime": "2018-07-09T09:04:42Z"
        }
	`)
	Init(config)

	response, err := GetActivation(42, 3)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &Activation{}, response), true)
}
