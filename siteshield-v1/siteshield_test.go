package siteshield

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestGetMap(t *testing.T) {

	defer gock.Off()

	mapid := "1"

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Get(fmt.Sprintf("/siteshield/v1/maps/%s", mapid)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
        "acknowledged": false,
        "contacts": [
            "test@akamai.com",
            "test2@akamai.com"
        ],
        "currentCidrs": [
            "131.103.136.0/24", "131.103.137.0/24",
            "165.254.127.0/24", "165.254.137.0/24", "184.25.254.0/24"
        ],
        "proposedCidrs": [
            "107.14.42.0/24", "117.103.188.0/24", "195.59.54.0/24",
            "209.211.216.0/24", "216.246.75.0/24"
        ],
        "ruleName": "a;s36.akamai.net",
        "type": "Production",
        "service": "S",
        "shared": false,
        "acknowledgeRequiredBy": 1392154239000,
        "previouslyAcknowledgedOn": 1392154239000
}
                `)

	Init(config)

	response, err := GetMap(mapid)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &SiteShieldMapResponse{}, response), true)
	assert.Equal(t, 5, len(response.CurrentCidrs))
	assert.Equal(t, 5, len(response.ProposedCidrs))
}

func TestListServices(t *testing.T) {

	mapid := "1"

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net")
	mock.
		Post(fmt.Sprintf("/siteshield/v1/maps/%s/acknowledge", mapid)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`
{
        "acknowledged": true,
        "contacts": [
            "test@akamai.com",
            "test2@akamai.com"
        ],
        "currentCidrs": [
            "131.103.136.0/24", "131.103.137.0/24",
            "165.254.127.0/24", "165.254.137.0/24", "184.25.254.0/24"
        ],
        "ruleName": "a;s36.akamai.net",
        "type": "Production",
        "service": "S",
        "shared": false,
        "acknowledgeRequiredBy": 1392154239000,
        "previouslyAcknowledgedOn": 1392154239000
}
                `)

	Init(config)

	response, err := Acknowledge(mapid)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &SiteShieldMapResponse{}, response), true)
	assert.Equal(t, 5, len(response.CurrentCidrs))
}
