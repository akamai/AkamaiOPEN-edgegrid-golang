package main

import (
  "fmt"
  "os"
  "sort"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
)

const edgerc = ".edgerc"

func main() {
  cfg_file := os.Getenv("HOME") + "/" + edgerc
  if err := edgegrid.InitServiceConfig(
    cfg_file,
    "cps",
    &cps.Config,
  ); err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
  }

  if e, err := cps.GetEnrollment("10002"); err != nil {
    fmt.Printf("Error: %v\n", err)
  } else {
    fmt.Printf("Enrollment Common Name: %s\n", e.CSR.CommonName)
    s := e.CSR.SANS
    sort.Strings(s)
    for _, san := range s {
      fmt.Printf("   %s\n", san)
    }
  }
}
