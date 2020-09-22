# Akamai OPEN EdgeGrid for GoLang v2

This is the Akamai API SDK

## Usage

GET Example:

```go
package main

import (
	"fmt"
    "github.com/apex/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func main() {
  // create a new session with the default .edgerc
  sess := session.Must(session.New())

  // get a papi client with the session
  client := papi.Client(sess)
  
  log.Debugf("calling GetGroups")
  
  // create a context with a log for this request
  ctx := session.ContextWithOptions(
     session.WithContextLog(log.Log),
  )
  
  // call a papi method, using the default thread context
  groupResp, err := client.GetGroups(ctx)
  if err != nil {
    panic(err)
  }

  for _, group := range groupResp.Groups.Items {
    log.Infof("Got group %q", group.GroupName)
  }
}

## Contribute


## Author


