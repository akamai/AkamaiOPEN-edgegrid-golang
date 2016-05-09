package edgegrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

var (
	testFile  = "testdata.json"
	timestamp = "20140321T19:34:21+0000"
	nonce     = "nonce-xx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	config    = Config{
		Host:         "https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
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

type JSONTests struct {
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

func TestMakeEdgeTimeStamp(t *testing.T) {
	actual := makeEdgeTimeStamp()
	expected := regexp.MustCompile(`^\d{4}[0-1][0-9][0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\+0000$`)
	if assert.Regexp(t, expected, actual, "Fail: Regex do not match") {
		t.Log("Pass: Regex matches")
	}
}

func TestCreateNonce(t *testing.T) {
	actual := createNonce()
	for i := 0; i < 100; i++ {
		expected := createNonce()
		assert.NotEqual(t, actual, expected, "Fail: Nonce matches")
	}
}

func TestCreateAuthHeader(t *testing.T) {
	var edgegrid JSONTests
	byt, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Test file not found, err %s", err)
	}
	url, err := url.Parse(config.Host)
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
		actual := config.createAuthHeader(req, timestamp, nonce)
		if assert.Equal(t, edge.ExpectedAuthorization, actual, fmt.Sprintf("Fail: %s", edge.Name)) {
			t.Logf("Pass: %s\n", edge.Name)
			t.Logf("Expected: %s - Actual %s", edge.ExpectedAuthorization, actual)
		}

	}
}

func TestAddRequestHeader(t *testing.T) {
	req, err := http.NewRequest("GET", config.Host, nil)
	if err != nil {
		t.Errorf("Fail: %s", err)
	}
	actual := AddRequestHeader(config, req)
	assert.NotEmpty(t, actual.Header.Get("Authorization"))
	assert.NotEmpty(t, actual.Header.Get("Content-Type"))
}

func TestInitConfigDefault(t *testing.T) {
	testSample := "sample_edgerc"
	testConfigDefault := InitConfig(testSample, "")
	assert.Equal(t, testConfigDefault.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigDefault.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.MaxBody, 131072)
	assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
}

func TestInitConfigBroken(t *testing.T) {
	testSample := "sample_edgerc"
	testConfigBroken := InitConfig(testSample, "broken")
	assert.Equal(t, testConfigBroken.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigBroken.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigBroken.MaxBody, 128*1024)
	assert.Equal(t, testConfigBroken.HeaderToSign, []string(nil))
}

func TestInitConfigUnparsable(t *testing.T) {
	testSample := "edgerc_that_doesnt_parse"
	assert.Panics(t, func() { InitConfig(testSample, "") }, "Fail: Should raise a PANIC")
}

func TestInitConfigNotFound(t *testing.T) {
	testSample := "edgerc_not_found"
	assert.Panics(t, func() { InitConfig(testSample, "") }, "Fail: Should raise a PANIC")
}

func TestInitConfigHeaders(t *testing.T) {
	testSample := "sample_edgerc"
	testConfigBroken := InitConfig(testSample, "headers")
	assert.Equal(t, testConfigBroken.HeaderToSign, []string{"X-MyThing1", "X-MyThing2"})
}

func TestInitConfigDashes(t *testing.T) {
	testSample := "sample_edgerc"
	assert.Panics(t, func() { InitConfig(testSample, "dashes") }, "Fail: Should raise a PANIC")
}
