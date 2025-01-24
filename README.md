# Akamai OPEN EdgeGrid for GoLang v9

![Build Status](https://github.com/akamai/akamaiOPEN-edgegrid-golang/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/AkamaiOPEN-edgegrid-golang/v9)](https://goreportcard.com/report/github.com/akamai/AkamaiOPEN-edgegrid-golang/v9)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/akamai/akamaiOPEN-edgegrid-golang)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/akamaiOPEN-edgegrid-golang?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/AkamaiOPEN-edgegrid-golang/v9)

The library implements an Authentication handler for HTTP requests using the [Akamai EdgeGrid Authentication](https://techdocs.akamai.com/developer/docs/authenticate-with-edgegrid) scheme for Go. It also currently provides Akamai REST API support for the Akamai Terraform Provider.


## Backward compatibility

This module isn't backward compatible with `v1`.

The `master` branch isn't representing `v1` anymore, it's currently representing the latest `v9`. `v1` has been moved to a dedicated `v1` branch.

## Concurrent usage

You can import the library packages alongside the `v1` library without any conflict. For example:

```go
import (
    papiv1 "github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
    papi "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
)
```

## Install

To use the library, you need to have Go 1.21+ installed on your system.


## Authentication

You can obtain the authentication credentials through an API client. Requests to the API are marked with a timestamp and a signature and are executed immediately.

1. [Create authentication credentials](https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials).

2. Place your credentials in an EdgeGrid file `~/.edgerc`, in the `[default]` section.

    ```
    [default]
    client_secret = C113nt53KR3TN6N90yVuAgICxIRwsObLi0E67/N8eRN=
    host = akab-h05tnam3wl42son7nktnlnnx-kbob3i3v.luna.akamaiapis.net
    access_token = akab-acc35t0k3nodujqunph3w7hzp7-gtm6ij
    client_token = akab-c113ntt0k3n4qtari252bfxxbsl-yvsdj
    ```

3. Import the `edgegrid` package, then use your local `.edgerc` by providing the path to your resource file and credentials' section header in the `edgegrid.New()` method.

   ```go
   package main

    import (
        "fmt"
        "io"
        "net/http"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
    )

    func main() {
        edgerc, err := edgegrid.New(
            edgegrid.WithFile("path/to/.edgerc"),
            edgegrid.WithSection("section-header"),
        )
    ...
    }
   ```

### Load from environment variables

Alternatively, you can pass your credentials to the environment variables. The library will look for the environment variables only if you use the `WithEnv()` method and set it to `true`. By default, it uses these variables:

- `AKAMAI_HOST`
- `AKAMAI_CLIENT_TOKEN`
- `AKAMAI_CLIENT_SECRET`
- `AKAMAI_ACCESS_TOKEN`
- `AKAMAI_MAX_BODY`

You can define multiple configurations by specifying the credentials' section header as an interfix that you insert between `AKAMAI_` and the credential name. The `WithEnv()` method uses the credentials' section header configured with the `WithSection()` method or it uses `default` if the credentials' section header wasn't configured.

For example, when passing `ccu` as the credentials' section header in the `WithSection()` method, the function will look for `AKAMAI_CCU_HOST`.

If the variable doesn't exist but the library was configured to search for it, the function either returns as only partially configured or it falls back to the filing data from the `.edgerc` file if you also used the `WithFile()` method.

```go
package main

    import (
        "fmt"
        "io"
        "net/http"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
    )

    func main() {
        // Load from AKAMAI_CCU_
        edgerc, err := edgegrid.New(
            edgegrid.WithEnv(true),
            edgegrid.WithSection("ccu"),
        )
    ...
    }
```

## Use as an auth signer

Import the `edgegrid` package to sign your requests. Provide the path to your `.edgerc`, your credentials' section header, and the appropriate endpoint information.

### Example

```go
package main

import (
    "fmt"
    "io"
    "net/http"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    client := http.Client{}

    req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    req.Header.Add("Accept", "application/json")
    edgerc.SignRequest(req)

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(body))
}
```

### Query string parameters

Pass the query parameters in the url after a question mark ("?") at the end of the main URL path.

```go
package main

import (
    "fmt"
    "io"
    "net/http"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    client := http.Client{}

    req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile?authGrants=true&notifications=true&actions=true", nil)
...
}
```

Alternatively, you can use one of these methods:

- [`req.URL.Query()`](https://pkg.go.dev/net/url#URL.Query). To add specific query params to an existing query.
- [`url.Values{}`](https://pkg.go.dev/net/url#Values.Add). To build a new set of params. Requires importing the `net/url` package.

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "net/url"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    client := http.Client{}

    req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    // If appending to an existing query
    q := req.URL.Query()
    q.Add("authGrants", "true")
    q.Add("notifications", "true")
    q.Add("actions", "true")

    // Or if building a new set of params
    q := url.Values{}
    q.Add("authGrants", "true")
    q.Add("notifications", "true")
    q.Add("actions", "true")

    req.URL.RawQuery = q.Encode()

    fmt.Println(req.URL.String())
    // Output:
    // /identity-management/v3/user-profile?actions=true&authGrants=true&notifications=true

    edgerc.SignRequest(req)

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(body))
}
```

### Headers

Enter request headers using the `req.Header.Add()` method.

> **Note:** You don't need to include the `Content-Type` and `Accept` headers. The authentication layer adds these values.

```go
package main

import (
    "fmt"
    "io"
    "net/http"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    client := http.Client{}

    req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    req.Header.Add("Accept", "application/json")
    edgerc.SignRequest(req)
...
}
```

### Body data

Import the `strings` package and provide the request body as an object in the `payload` property.

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    payload := strings.NewReader(`{
        "contactType": "Billing",
        "country": "USA",
        "firstName": "John",
        "lastName": "Smith",
        "preferredLanguage": "English",
        "sessionTimeOut": 30,
        "timeZone": "GMT",
        "phone": "3456788765"
      }`)

    client := http.Client{}

    req, err := http.NewRequest(http.MethodPut, "/identity-management/v3/user-profile/basic-info", payload)
...
}
```

### Debug

Enable debugging to get additional information about a request and response.

You may find using `fmt.Println()` not sufficient, as the output isn't formatted and is difficult to read.

To pretty-print the HTTP request and response, use these methods:

- [`httputil.DumpRequest()`](https://pkg.go.dev/net/http/httputil#DumpRequest). To pretty-print the request on the server side.
- [`httputil.DumpRequestOut()`](https://pkg.go.dev/net/http/httputil#DumpRequestOut). To dump the request on the client side.
- [`httputil.DumpResponse()`](https://pkg.go.dev/net/http/httputil#DumpResponse). To log the server response.

In each of these functions, the second argument is a flag indicating if the body of the request/response should also be returned.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )

    client := http.Client{}

    req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Add("Accept", "application/json")
    edgerc.SignRequest(req)

    reqDump, err := httputil.DumpRequestOut(req, false)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("REQUEST:\n%s", string(reqDump))

    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    resDump, err := httputil.DumpResponse(res, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("RESPONSE:\n%s", string(resDump))
}
```

## Use with the `session`

When making a call with the [auth signer approach](#use-as-an-auth-signer), you can add the `session` package to log additional information about the call. You can also define the `struct` for the requested object as a variable in that call, or you can omit the `struct` and define the requested object as an any type variable with the `any` keyword.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

func main() {
	edgerc, _ := edgegrid.New(
		edgegrid.WithFile("~/.edgerc"),
		edgegrid.WithSection("default"),
	)

	sess, err := session.New(
		session.WithSigner(edgerc),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	var userProfile struct {
		FirstName                string         `json:"firstName"`
		LastName                 string         `json:"lastName"`
		UserName                 string         `json:"uiUserName"`
		TimeZone                 string         `json:"timeZone"`
		Country                  string         `json:"country"`
		PreferredLanguage        string         `json:"preferredLanguage"`
		SessionTimeOut           *int           `json:"sessionTimeOut"`
	}

	req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := sess.Exec(req, &userProfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result, userProfile)
}
```

The `session` package also supports the structured logging interface from `github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log`. Thanks to this, you can adjust a logger (for example, change the logging level to `Debug`) in one of these ways:

- Apply a logger globally with the `log.SetLogger()` method to use it in all sessions. You can retrieve the logger from `context` using the `log.FromContext()` method.

  > **Note:** This method works also with the [SDK approach](#use-as-an-sdk).

  ```go
    package main

    import (
        "fmt"
        "log/slog"
        "net/http"
        "os"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
    )


    func main() {
        edgerc, _ := edgegrid.New(
            edgegrid.WithFile("~/.edgerc"),
            edgegrid.WithSection("default"),
        )

        l := log.NewSlogAdapter(log.NewSlogHandler(os.Stderr, &slog.HandlerOptions{
          Level: slog.LevelDebug,
        }))
  
        log.SetLogger(l)

        sess, _ := session.New(
            session.WithSigner(edgerc),
            session.WithHTTPTracing(true),
        )

        var userProfile struct {
            FirstName                string         `json:"firstName"`
            LastName                 string         `json:"lastName"`
            UserName                 string         `json:"uiUserName"`
            TimeZone                 string         `json:"timeZone"`
            Country                  string         `json:"country"`
            PreferredLanguage        string         `json:"preferredLanguage"`
            SessionTimeOut           *int           `json:"sessionTimeOut"`
        }

        req, _ := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)

        result, err := sess.Exec(req, &userProfile)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(result, userProfile)
    }
  ```

- Apply a logger to all operations across a session with the `session.WithLog()` method.
  
  > **Note:** This method works also with the [SDK approach](#use-as-an-sdk).

  ```go
    package main

    import (
        "fmt"
        "log/slog"
        "net/http"
        "os"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
    )


    func main() {
        edgerc, _ := edgegrid.New(
            edgegrid.WithFile("~/.edgerc"),
            edgegrid.WithSection("default"),
        )

        l := log.NewSlogAdapter(log.NewSlogHandler(os.Stderr, &slog.HandlerOptions{
          Level: slog.LevelDebug,
        }))

        sess, _ := session.New(
            session.WithSigner(edgerc),
            session.WithLog(l),
            session.WithHTTPTracing(true),
        )

        var userProfile struct {
            FirstName                string         `json:"firstName"`
            LastName                 string         `json:"lastName"`
            UserName                 string         `json:"uiUserName"`
            TimeZone                 string         `json:"timeZone"`
            Country                  string         `json:"country"`
            PreferredLanguage        string         `json:"preferredLanguage"`
            SessionTimeOut           *int           `json:"sessionTimeOut"`
        }

        req, _ := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)

        result, err := sess.Exec(req, &userProfile)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(result, userProfile)
    }
  ```

- Overwrite a default logger for a specific request with the `req.WithContext()` method.

  > **Note:** This method doesn't work with the [SDK approach](#use-as-an-sdk).

  ```go
    package main

    import (
        "fmt"
        "log/slog"
        "net/http"
        "os"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
    )


    func main() {
        edgerc, _ := edgegrid.New(
            edgegrid.WithFile("~/.edgerc"),
            edgegrid.WithSection("default"),
        )

        l := log.NewSlogAdapter(log.NewSlogHandler(os.Stderr, &slog.HandlerOptions{
          Level: slog.LevelDebug,
        }))
  
        sess, _ := session.New(
            session.WithSigner(edgerc),
            session.WithHTTPTracing(true),
        )

        var userProfile struct {
            FirstName                string         `json:"firstName"`
            LastName                 string         `json:"lastName"`
            UserName                 string         `json:"uiUserName"`
            TimeZone                 string         `json:"timeZone"`
            Country                  string         `json:"country"`
            PreferredLanguage        string         `json:"preferredLanguage"`
            SessionTimeOut           *int           `json:"sessionTimeOut"`
        }

        req, _ := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)

        req = req.WithContext(session.ContextWithOptions(req.Context(), session.WithContextLog(l)))

        result, err := sess.Exec(req, &userProfile)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(result, userProfile)
    }
  ```

If you don't provide a custom logger, the library will use its default logger included in the `session.New()` method.

You can also create your own custom logger. See the [logger README](pkg/log/README.md) for more details and an example.

You need the `session.WithHTTPTracing` option in all these ways of adding a structured logging interface. This is to log all requests, responses, and its headers. It also requires at least the `Debug` logging level.

### Custom request headers

When using the `session` package with the [auth signer approach](#use-as-an-auth-signer), you can update the context to pass custom request headers in the `session.WithContextHeaders()` method.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

func main() {
	edgerc, _ := edgegrid.New(
		edgegrid.WithFile("~/.edgerc"),
		edgegrid.WithSection("default"),
	)

	sess, _ := session.New(
		session.WithSigner(edgerc),
	)

	var userProfile struct {
		FirstName                string         `json:"firstName"`
		LastName                 string         `json:"lastName"`
		UserName                 string         `json:"uiUserName"`
		TimeZone                 string         `json:"timeZone"`
		Country                  string         `json:"country"`
		PreferredLanguage        string         `json:"preferredLanguage"`
		SessionTimeOut           *int           `json:"sessionTimeOut"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)

	customHeader := make(http.Header)
    customHeader.Set("X-Custom-Header", "some custom value")

	req = req.WithContext(session.ContextWithOptions(req.Context(), session.WithContextHeaders(customHeader)))

	result, err := sess.Exec(req, &userProfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result, userProfile)
}
```

### Retries

With the `session` package, you can use the `session.WithRetries()` option for creating new sessions with global GET retries. You can set up this option using its default configuration or customize it with these parameters:

- `RetryMax`. The maximum number of API request retries.
- `RetryWaitMin`. The minimum wait time in `time.Duration` between API requests retries.
- `RetryWaitMax`. The maximum wait time in `time.Duration` between API requests retries.
- `ExcludedEndpoints`. The list of path expressions defining endpoints which should be excluded from the retry feature.

> **Note** This option works also with the [SDK approach](#use-as-an-sdk).

```go
package main

import (
	"fmt"
	"net/http"
    "time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

func main() {
	edgerc, _ := edgegrid.New(
		edgegrid.WithFile("~/.edgerc"),
		edgegrid.WithSection("default"),
	)

    // The retries feature with default configuration.
    sess, err := session.New(
        session.WithSigner(edgerc),
        session.WithRetries(session.NewRetryConfig()),
    )

    // The retries feature with custom configuration.
	sess, err := session.New(
		session.WithSigner(edgerc),
		session.WithRetries(session.RetryConfig{
			RetryMax:     5,
			RetryWaitMax: time.Minute * 2,
			RetryWaitMin: time.Second * 3,
		}),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	var userProfile struct {
		FirstName                string         `json:"firstName"`
		LastName                 string         `json:"lastName"`
		UserName                 string         `json:"uiUserName"`
		TimeZone                 string         `json:"timeZone"`
		Country                  string         `json:"country"`
		PreferredLanguage        string         `json:"preferredLanguage"`
		SessionTimeOut           *int           `json:"sessionTimeOut"`
	}

	req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/user-profile", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := sess.Exec(req, &userProfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result, userProfile)
}
```

## Use as an SDK

Apart from signing the requests with the `edgegrid` package, you can use specific API clients and methods from this library.

### Overview of packages

To see what methods you can use within a given package, go to the `pkg` directory and select the package name. Each package contains:

- A file named after a package name or the package's abbreviated name with a main `interface`.
- The main `interface` can contain its own methods or other `interfaces` with their own methods.
- `struct` blocks distributed across the package's files, each containing a list of parameters you can pass in a given method.

| Package                                       | Description                                                                                          |
|-----------------------------------------------|------------------------------------------------------------------------------------------------------|
| [Application Security](./pkg/appsec/)         | Manage security configurations, security policies, match targets, rate policies, and firewall rules. |
| [Bot Manager](./pkg/botman/)                  | Identify, track, and respond to bot activity on your domain or in your app.                          |
| [Certificate Provisioning System](./pkg/cps/) | Manage the full life cycle of SSL certificates for your ​Akamai​ CDN applications.                   |
| [Client Lists](./pkg/clientlists/)            | Reduce harmful security attacks by allowing only trusted IP/CIDRs, locations, autonomous system numbers, and TLS fingerprints to access your services and content.|
| [Cloud Access Manager](./pkg/cloudaccess/)    | Enable cloud origin authentication and securely store and manage your cloud origin credentials as access keys. |
| [Cloudlets](./pkg/cloudlets/)                 | Solve specific business challenges using value-added apps that complement ​Akamai​'s core solutions. |
| [Cloud Wrapper](./pkg/cloudwrapper/)          | Provide your customers with a more consistent user experience by adding a custom caching layer that improves the connection between your cloud infrastructure and the Akamai platform.|
| [DataStream](./pkg/datastream/)               | Monitor activity on the ​Akamai​ platform and send live log data to a destination of your choice.    |
| [Edge DNS](./pkg/dns/)                        | Replace or augment your DNS infrastructure with a cloud-based authoritative DNS solution.            |
| [EdgeGrid](./pkg/edgegrid/)                   | Parse the Akamai `.edgerc` configuration and sign HTTP requests.            |
| [EdgeGrid Errors](./pkg/edgegriderr/)         | Parse validation errors to make them more readable.            |
| [Edge Hostnames](./pkg/hapi/)                 | Manage how requests for your site, app, or content map to Akamai edge servers.                                 |
| [EdgeWorkers](./pkg/edgeworkers/)             | Execute JavaScript functions at the edge to optimize site performance and customize web experiences. |
| [Errors](./pkg/errs/)                         | Use utilities for working with errors during JSON data unmarshalling.            |
| [Global Traffic Management](./pkg/gtm/)       | Use load balancing to manage website and mobile performance demands.                                 |
| [Identity and Access Management](./pkg/iam/)  | Create users and groups, and define policies that manage access to your Akamai applications.         |
| [Image and Video Manager](./pkg/imaging/)     | Automate image and video delivery optimizations for your website visitors.                           |
| [Log](./pkg/log/)                             | Add the structured logging interface.                           |
| [Network Lists](./pkg/networklists/)          | Automate the creation, deployment, and management of lists used in ​Akamai​ security products.       |
| [Pointer Record](./pkg/ptr/)                  | Create pointers to values of any type.  |
| [Property Manager](./pkg/papi/)               | Define rules and behaviors that govern your website delivery based on match criteria.                |
| [Session](./pkg/session/)                     | Manage the base secure HTTP client and requests for Akamai APIs.  |

### Example

To use the library as an SDK, import the `edgegrid` package to sign your requests, the `session` package, and a specific package available in this library (for example, `iam`) to use the methods included in that package. Provide the path to your `.edgerc`, your credentials' section header, and refer to a specific API client and method.

```go
package main

import (
    "context"
    "fmt"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    sess, err := session.New(
        session.WithSigner(edgerc),
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    client := iam.Client(sess)

    resp, err := client.GetUser(context.TODO(), iam.GetUserRequest{IdentityID: "A-BC-1234567"})
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(resp)
}
```

### Headers

You don't need to add any required or optional headers to a request. The authentication layer adds these values automatically for you.

### Query string parameters and body data

Depending on the method type, use `structs` of the provided API `struct` to pass query parameters and body data within a request method.

```go
package main

import (
    "context"
    "fmt"

    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
)

func main() {
    edgerc, err := edgegrid.New(
        edgegrid.WithFile("~/.edgerc"),
        edgegrid.WithSection("default"),
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    sess, err := session.New(
        session.WithSigner(edgerc),
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    client := iam.Client(sess)

    sessTimeOut := 14400

    resp, err := client.UpdateUserInfo(
        context.TODO(),
        iam.UpdateUserInfoRequest{
            IdentityID: "A-BC-123456",
            User: iam.UserBasicInfo{
                Country: "USA",
                FirstName: "John",
                LastName: "Smith",
                Phone: "3456788765",
                PreferredLanguage: "English",
                SessionTimeOut: &sessTimeOut,
                TimeZone: "GMT",
            },
        },
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(resp)
}
```

## Reporting issues

To report an issue or make a suggestion, create a new [GitHub issue](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/issues).

## License

Copyright 2025 Akamai Technologies, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use these files except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.