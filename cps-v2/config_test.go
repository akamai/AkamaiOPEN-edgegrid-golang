package cps

import (
  "testing"
	"github.com/stretchr/testify/assert"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

var (
	config = edgegrid.Config{
		Host:         "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "testxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      131072,
		Debug:        false,
	}
  accept = "application/vnd.akamai.cps.enrollments.v4+json"
)

func TestCPS_NewConfig(t *testing.T) {
  Config.NewConfig(config)
  assert.Equal(t, config, Config.EdgeGridConfig)
  assert.Equal(t, accept, Config.Headers["Accept"])
}

func TestCPS_InitServiceConfig(t *testing.T) {
  edgegrid.InitServiceConfig("../testdata/sample_edgerc", "test", &Config)
  assert.Equal(t, config, Config.EdgeGridConfig)
  assert.Equal(t, accept, Config.Headers["Accept"])
}

func TestCPS_UpdateHeaders(t *testing.T) {
  Config.NewConfig(config)
  Config.Headers["Accept"] = "test/header"
  Config.Headers["Dummy"] = "Dummy Header"
  assert.Equal(t, "test/header", Config.Headers["Accept"])
  assert.NotEqual(t, accept, Config.Headers["Accept"])
  assert.Equal(t, "Dummy Header", Config.Headers["Dummy"])
}
