package dns

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	log "github.com/sirupsen/logrus"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	Config edgegrid.Config
)

// GetZone retrieves a DNS Zone for a given hostname
func GetZone(hostname string) (*Zone, error) {
	zone := NewZone(hostname)
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-dns/v1/zones/"+hostname,
		nil,
	)
	if err != nil {
		return nil, err
	}

	log.Debugf("Request being sent: %s", req)
	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}
	log.Debugf("Response received: %s", res)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, &ZoneNotFoundError{zoneName: hostname}
	} else {
		err = client.BodyJSON(res, &zone)
		if err != nil {
			return nil, err
		}

		return &zone, nil
	}
}
