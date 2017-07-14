// Package configdns provides a simple wrapper around the Akamai FastDNS DNS Management API
package dns

import "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

// Init sets the FastDNS edgegrid Config
func Init(config edgegrid.Config) {
	Config = config
}
