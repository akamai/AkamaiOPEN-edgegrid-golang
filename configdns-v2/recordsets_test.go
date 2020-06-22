package dnsv2

import (
	"testing"
	//edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var recordsetBody = []byte(fmt.Sprintf(`{
    "recordsets": [{
            "name": "new.%s",
            "type": "CNAME",
            "ttl": 900,
            "rdata": ["www.example.com"]
        },
        {
            "name": "a_rec_%s",
            "type": "A",
            "ttl": 900,
            "rdata": ["10.0.0.10"]
        }]}`, dnsTestZone, dnsTestZone))

func createTestRecordsets() *Recordsets {

	rs := &Recordsets{}
	jsonhooks.Unmarshal(recordsetBody, rs)

	return rs

}

func TestListRecordsets(t *testing.T) {

	/*
			// for live testing
		        config, err := edge.Init("","")
		        if err != nil {
		                t.Fatalf("TestListRecordsets failed initializing: %s", err.Error())
		        }
	*/

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/recordsets", dnsTestZone))
	mock.
		Get(fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(fmt.Sprintf(`{
			"metadata": {
				"totalElements":2
			},
			"recordsets": [
        			{
            				"name": "new.%s",
            				"type": "CNAME",
            				"ttl": 900,
            				"rdata": ["www.example.com"]
        			},
       				{
            				"name": "a_rec_%s",
         	 	  		"type": "A",
            				"ttl": 900,
            				"rdata": ["10.0.0.10"]
       				}]}`, dnsTestZone, dnsTestZone))

	Init(config)
	listResp, err := GetRecordsets(dnsTestZone)
	assert.NoError(t, err)
	assert.Equal(t, int(len(listResp.Recordsets)), listResp.Metadata.TotalElements)

}

func TestCreateRecordsets(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/recordsets", dnsTestZone))
	mock.
		Post(fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(204).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	Init(config)
	sets := createTestRecordsets()
	err := sets.Save(dnsTestZone)
	assert.NoError(t, err)

}

func TestUpdateRecordsets(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/recordsets", dnsTestZone))
	mock.
		Put(fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(204).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	Init(config)
	sets := createTestRecordsets()
	err := sets.Update(dnsTestZone)
	assert.NoError(t, err)

}
