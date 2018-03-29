package main

import (
  "fmt"
  "sort"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        false,
	}
)

func main() {
  cps.Config.NewConfig(config)
  if e, err := cps.GetEnrollments(); err != nil {
    fmt.Printf("Error: %v\n", err)
  } else {
    for _, en := range e.Enrollments {
      fmt.Printf("Enrollment Common Name: %s\n", en.CSR.CommonName)
      s := en.CSR.SANS
      sort.Strings(s)
      for _, san := range s {
        fmt.Printf("   %s\n", san)
      }
    }
  }
}
