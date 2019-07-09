package reportsgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

        "net/http"
        "fmt"
        "net/http/httputil"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
        
        debug bool
)

// Init sets the GTM edgegrid Config
func Init(config edgegrid.Config) {
	Config = config
        debug = false
}

// Utility func to print http req
func printHttpRequest(req *http.Request, body bool) {

        if !debug {
                return
        }
        b, err := httputil.DumpRequestOut(req, body)
        if err == nil {
                 fmt.Println(string(b))
        }
}

// Utility func to print http response
func printHttpResponse(res *http.Response, body bool) {

        if !debug {
                return
        }
        b, err := httputil.DumpResponse(res, body)
        if err == nil {
                 fmt.Println(string(b))
        }
}





