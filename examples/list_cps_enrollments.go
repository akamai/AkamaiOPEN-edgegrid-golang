package main

import (
  "fmt"
  "os"
  "sort"
  "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  "github.com/xelwarto/AkamaiOPEN-edgegrid-golang/cps-v2"
)

const edgerc = ".edgerc"

func main() {
  config, err := edgegrid.InitEdgeRc(os.Getenv("HOME") + "/" + edgerc, "cps")
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
  }

  e := cps.NewEnrollments()
  if err := e.Get(config); err != nil {
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
