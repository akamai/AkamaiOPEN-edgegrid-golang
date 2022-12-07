# Akamai OPEN EdgeGrid for GoLang v1

[![Build Status](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/actions/workflows/checks.yml/badge.svg?branch=v1)](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/actions/workflows/checks.yml?query=branch%3Av1)
[![GoDoc](https://godoc.org/github.com/akamai/AkamaiOPEN-edgegrid-golang?status.svg)](https://godoc.org/github.com/akamai/AkamaiOPEN-edgegrid-golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/AkamaiOPEN-edgegrid-golang)](https://goreportcard.com/report/github.com/akamai/AkamaiOPEN-edgegrid-golang)
[![License](http://img.shields.io/:license-apache-blue.svg)](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/blob/v1/LICENSE)

This library implements an Authentication handler for [net/http](https://golang.org/pkg/net/http/)
that provides the [Akamai OPEN Edgegrid Authentication](https://developer.akamai.com/introduction/Client_Auth.html) scheme.
For more information visit the [Akamai OPEN Developer Community](https://developer.akamai.com).
This library has been released as a v1 library though future development will be on the **develop** branch for version v3.

## Announcing Akamai OPEN EdgeGrid for GoLang v3 (release v3.0.0)

The v3 version of this module is under active development and provides a subset of Akamai APIs for use in the 
Akamai Terraform Provider.

New users are encouraged to adopt v3 version it is a simpler API wrapper with little to no business logic.

Current direct users of this v1.x.x library are recommended to continue to use the v1 version as initialization 
and package structure has significantly changed in v3 and will require substantial work to migrate existing 
applications. Non-backwards compatible changes were made to improve the code quality and make the project more 
maintainable. 

## Usage of the v1 library

GET Example:

```go
  package main

import (
	"fmt"
	"io/ioutil"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	config, _ := edgegrid.Init("~/.edgerc", "default")

	// Retrieve all locations for diagnostic tools
	req, _ := client.NewRequest(config, "GET", "/diagnostic-tools/v2/ghost-locations/available", nil)
	resp, _ := client.Do(config, req)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
```

Parameter Example:

```go
  package main

import (
	"fmt"
	"io/ioutil"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	config, _ := edgegrid.Init("~/.edgerc", "default")

	// Retrieve dig information for specified location
	req, _ := client.NewRequest(config, "GET", "/diagnostic-tools/v2/ghost-locations/zurich-switzerland/dig-info", nil)

	q := req.URL.Query()
	q.Add("hostName", "developer.akamai.com")
	q.Add("queryType", "A")

	req.URL.RawQuery = q.Encode()
	resp, _ := client.Do(config, req)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
```

POST Example:

```go
  package main

  import (
    "fmt"
    "io/ioutil"
    "net/http"
    
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  )

  func main() {
    config, _ := edgegrid.Init("~/.edgerc", "default")
    
    // Acknowledge a map
    req, _ := client.NewRequest(config, "POST", "/siteshield/v1/maps/1/acknowledge", nil)
    resp, _ := client.Do(config, req)

    defer resp.Body.Close()
    byt, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(byt))
  }
```

PUT Example:

```go
  package main

  import (
    "fmt"
    "io/ioutil"
    "net/http"
    
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  )

  func main() {

    config, _ := edgegrid.Init("~/.edgerc", "default")
    body := []byte("{\n  \"name\": \"Simple List\",\n  \"type\": \"IP\",\n  \"unique-id\": \"345_BOTLIST\",\n  \"list\": [\n    \"192.168.0.1\",\n    \"192.168.0.2\",\n  ],\n  \"sync-point\": 0\n}")
    
    // Update a Network List
    req, _ := client.NewJSONRequest(config, "PUT", "/network-list/v1/network_lists/unique-id?extended=extended", body)
    resp, _ := client.Do(config, req)

    defer resp.Body.Close()
    byt, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(byt))
  }
```

Alternatively, your program can read it from config struct.

```go
  package main

  import (
    "fmt"
    "io/ioutil"
    "net/http"
    
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
  )

  func main() {
    config := edgegrid.Config{
      Host : "xxxxxx.luna.akamaiapis.net",
      ClientToken:  "xxxx-xxxxxxxxxxx-xxxxxxxxxxx",
      ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
      AccessToken:  "xxxx-xxxxxxxxxxx-xxxxxxxxxxx",
      MaxBody:      1024,
      HeaderToSign: []string{
        "X-Test1",
        "X-Test2",
        "X-Test3",
      },
      Debug:        false,
    }
    
   // Retrieve all locations for diagnostic tools
	req, _ := client.NewRequest(config, "GET", fmt.Sprintf("https://%s/diagnostic-tools/v2/ghost-locations/available",config.Host), nil)
	resp, _ := client.Do(config, req)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
```

## Contribute

1. Fork [the repository](https://github.com/akamai/AkamaiOPEN-edgegrid-golang) to start making your changes to the **v1** branch
2. Send a pull request.

## Author

[Davey Shafik](mailto:dshafik@akamai.com) - Developer Evangelist @ [Akamai Technologies](https://developer.akamai.com)
[Nick Juettner](mailto:hello@juni.io) - Software Engineer @ [Zalando SE](https://tech.zalando.com/)  

