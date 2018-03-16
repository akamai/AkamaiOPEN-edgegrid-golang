package cps

import (
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// CPSConfig struct provides package configuration and
// implements the edgegrid.ServiceConfig interface
type CPSConfig struct {
  EdgeGridConfig edgegrid.Config
  Headers map[string]string
}

// Creates a new CPS Configuration
func (c *CPSConfig) NewConfig(config edgegrid.Config) {
  c.EdgeGridConfig = config
  c.Headers["Accept"] = "application/vnd.akamai.cps.enrollments.v4+json"
}

// Package level configuration variable
var Config = CPSConfig{
  Headers: make(map[string]string),
}
