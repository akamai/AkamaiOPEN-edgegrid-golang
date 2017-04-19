package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/akamai-open/AkamaiOPEN-edgegrid-golang"
)

func random(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Intn(max-min) + min

	return random
}

type LocationsResponse struct {
	Locations []string `json:"locations"`
}

type DigResponse struct {
	Dig struct {
		Hostname    string `json:"hostname"`
		QueryType   string `json:"queryType"`
		Result      string `json:"result"`
		ErrorString string `json:"errorString"`
	} `json:"dig"`
}

func main() {
	config, err := edgegrid.Init("~/.edgerc", "default")
	config.Debug = false
	if err == nil {
		client, err := edgegrid.New(nil, config)
		if err == nil {
			fmt.Println("Requesting locations that support the diagnostic-tools API.")

			res, err := client.Get("/diagnostic-tools/v1/locations")
			if err != nil {
				log.Fatal(err.Error())
			}

			locationsResponse := LocationsResponse{}
			res.BodyJSON(&locationsResponse)

			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Printf("There are %d locations that can run dig in the Akamai Network\n", len(locationsResponse.Locations))

			if len(locationsResponse.Locations) == 0 {
				log.Fatal("No locations found")
			}

			location := locationsResponse.Locations[random(0, len(locationsResponse.Locations))-1]

			fmt.Println("We will make our call from " + location)

			fmt.Println("Running dig from " + location)

			client.Timeout = 5 * time.Minute
			res, err = client.Get("/diagnostic-tools/v1/dig?hostname=developer.akamai.com&location=" + url.QueryEscape(location) + "&queryType=A")
			if err != nil {
				log.Fatal(err.Error())
			}

			digResponse := DigResponse{}
			res.BodyJSON(&digResponse)
			fmt.Println(digResponse.Dig.Result)
		} else {
			log.Fatal(err.Error())
		}
	}
}
