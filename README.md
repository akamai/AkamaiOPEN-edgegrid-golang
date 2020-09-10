# Akamai OPEN EdgeGrid for GoLang v2

This is the Akamai API SDK

## Usage

GET Example:

```go
  package main

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func main() {
  sess := session.Must(session.New())

  groupResp, err := papi.API(sess).GetGroups(context.Background())
  if err != nil {
    panic(err)
  }

  for _, group := range groupResp.Groups.Items {
    fmt.Printf("%s", group.GroupName)
  }
}

## Contribute


## Author


