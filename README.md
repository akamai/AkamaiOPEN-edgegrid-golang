# EdgeGrid for GoLang

[![Build Status](https://travis-ci.org/njuettner/edgegrid.svg?branch=master)](https://travis-ci.org/njuettner/edgegrid)
[![Coverage Status](https://coveralls.io/repos/github/njuettner/edgegrid/badge.svg?branch=master)](https://coveralls.io/github/njuettner/edgegrid?branch=master)
[![GoDoc](https://godoc.org/github.com/njuettner/edgegrid?status.svg)](https://godoc.org/github.com/njuettner/edgegrid)
[![Go Report Card](https://goreportcard.com/badge/github.com/njuettner/edgegrid)](https://goreportcard.com/report/github.com/njuettner/edgegrid)

This library implements an Authentication handler for [net/http](https://golang.org/pkg/net/http/)
that provides the [Akamai {OPEN} Edgegrid Authentication](https://developer.akamai.com/introduction/Client_Auth.html) 
scheme. For more information visit the [Akamai {OPEN} Developer Community](https://developer.akamai.com).

GET Example:

```go
  package main

  import (
    "fmt"
    "github.com/njuettner/edgegrid"
    "io/ioutil"
    "net/http"
  )

  func main() {
    client := http.Client{}

    config := edgegrid.InitConfig("~/.edgerc", "default")

    // Retrieve a list all maps belonging to an account
    req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/siteshield/v1/maps", config.Host), nil)
    req = edgegrid.AddRequestHeader(config, req)
    resp, _ := client.Do(req)

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
    "github.com/njuettner/edgegrid"
    "io/ioutil"
    "net/http"
  )

  func main() {
    client := http.Client{}

    config := edgegrid.InitConfig("~/.edgerc", "default")

    //  The ID of the report pack.
    reportPackId = "1"

    // List Audience Analytics Data Stores
    req, _ := http.NewRequest("PUT", fmt.Sprintf("https://%s/media-analytics/v1/audience-analytics/report-packs/%s", config.Host, reportPackId), nil)
    req = edgegrid.AddRequestHeader(config, req)
    resp, _ := client.Do(req)

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
    "github.com/njuettner/edgegrid"
    "io/ioutil"
    "net/http"
  )

  func main() {
    client := http.Client{}

    config := edgegrid.InitConfig("~/.edgerc", "default")
    
    // Acknowledge a map
    req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s/siteshield/v1/maps/1/acknowledge", config.Host), nil)
    req = edgegrid.AddRequestHeader(config, req)
    resp, _ := client.Do(req)

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
    "github.com/njuettner/edgegrid"
    "io/ioutil"
    "net/http"
  )

  func main() {
    client := http.Client{}

    config := edgegrid.InitConfig("~/.edgerc", "default")
    body := []byte("{\n  \"name\": \"Simple List\",\n  \"type\": \"IP\",\n  \"unique-id\": \"345_BOTLIST\",\n  \"list\": [\n    \"192.168.0.1\",\n    \"192.168.0.2\",\n  ],\n  \"sync-point\": 0\n}")
    
    // Update a Network List
    req, _ := http.NewRequest("PUT", fmt.Sprintf("https://%s/network-list/v1/network_lists/unique-id?extended=extended", config.Host), bytes.NewBuffer(body))
    req = edgegrid.AddRequestHeader(config, req)
    resp, _ := client.Do(req)

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
    "github.com/njuettner/edgegrid"
    "io/ioutil"
    "net/http"
  )

  func main() {
    client := http.Client{}
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
    
    // Retrieve a list all maps belonging to an account
    req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/siteshield/v1/maps", config.Host), nil)
    req = edgegrid.AddRequestHeader(config, req)
    resp, _ := client.Do(req)

    defer resp.Body.Close()
    byt, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(byt))
  }
```

For more examples please checkout [Akamai API Catalog]('https://developer.akamai.com/api/'). Click on the API Name you're currently interested in and use the reference.
Use also the Mock Server to test your requests (GO is supported).

## Installation

```bash
  $ go get github.com/njuettner/edgegrid
```

## Contribute

1. Fork [the repository](https://github.com/njuettner/edgegrid) to start making your changes to the **master** branch
2. Send a pull request.

## Author

[Nick Juettner](mailto:hello@juni.io)

