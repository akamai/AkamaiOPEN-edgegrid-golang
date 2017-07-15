package dns

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
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
		"/edgegrid-dns/v1/zones/"+hostname,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, fmt.Errorf(errorMap[ErrZoneNotFound], hostname)
	} else {
		err = client.BodyJSON(res, &zone)
		if err != nil {
			return nil, err
		}

		return &zone, nil
	}
}
