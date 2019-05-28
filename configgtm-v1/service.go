package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

        log "github.com/sirupsen/logrus"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
)

// Init sets the GTM edgegrid Config
func Init(config edgegrid.Config) {
	Config = config

        if config.Debug {
                log.SetLevel(log.DebugLevel)
        }
 
        log.Debugf("Log debug level set to  %v", config.Debug)

}
