package edgegrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var (
	testFile  = "testdata.json"
	baseURL   = "https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/"
	timestamp = "20140321T19:34:21+0000"
	nonce     = "nonce-xx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	base      = Base{
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        true,
		HeaderToSign: []string{
			"X-Test1",
			"X-Test2",
			"X-Test3",
		},
	}
)

type JsonTests struct {
	Tests []Test `json:"tests"`
}
type Test struct {
	Name    string `json:"testName"`
	Request struct {
		Method  string              `json:"method"`
		Path    string              `json:"path"`
		Headers []map[string]string `json:"headers"`
		Data    string              `json:"data"`
	} `json:"request"`
	ExpectedAuthorization string `json:"expectedAuthorization"`
}

func TestmakeEdgeTimeStamp(t *testing.T) {
}

func TestMakeHeader(t *testing.T) {
	var edgegrid JsonTests
	byt, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Test file not found, err %s", err)
	}
	url, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("URL is not parsable, err %s", err)
	}
	json.Unmarshal(byt, &edgegrid)
	for _, edge := range edgegrid.Tests {
		url.Path = edge.Request.Path
		req, _ := http.NewRequest(
			edge.Request.Method,
			url.String(),
			bytes.NewBuffer([]byte(edge.Request.Data)),
		)
		for _, header := range edge.Request.Headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		actual := base.createAuthHeader(req, timestamp, nonce)
		if assert.Equal(t, edge.ExpectedAuthorization, actual, fmt.Sprintf("Failed: %s", edge.Name)) {
			t.Logf("Passed: %s\n", edge.Name)
			t.Logf("Expected: %s - Actual %s", edge.ExpectedAuthorization, actual)
		}

	}
}
