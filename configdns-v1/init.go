// Package configdns provides a simple wrapper around the Akamai FastDNS DNS Management API
package dns

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	log "github.com/sirupsen/logrus"
)

// Init sets the FastDNS edgegrid Config
func Init(config edgegrid.Config) {
	Config = config
	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	}
}
