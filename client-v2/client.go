// Package client is a simple library for http.Client to sign Akamai OPEN Edgegrid API requests
package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	edgegrid "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
)

const (
	libraryVersion = "0.5.0"
)

var (
	edgercConfig  edgegrid.Config
	isInitialized bool
	UserAgent     = "Akamai-Open-Edgegrid-golang/" + libraryVersion + " golang/" + strings.TrimPrefix(runtime.Version(), "go")
	Client        = http.DefaultClient
)

//Init initialise our client with selected configuration which later on can be used
func Init(c edgegrid.Config) {
	edgercConfig = c
	fmt.Println("Print me host :")
	fmt.Println(c.Host)
	isInitialized = true
}

// prepareHTTPRequest creates an HTTP request accordingly to requirements of Akamai API
func prepareHTTPRequest(method, path string, body io.Reader) (*http.Request, error) {
	var (
		baseURL *url.URL
		err     error
	)

	if isInitialized != true {
		return nil, err
	}

	if strings.HasPrefix(edgercConfig.Host, "https://") {
		baseURL, err = url.Parse(edgercConfig.Host)
	} else {
		baseURL, err = url.Parse("https://" + edgercConfig.Host)
	}

	if err != nil {
		return nil, err
	}

	rel, err := url.Parse(strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}

	u := baseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", UserAgent)

	return req, nil
}

// RequestHTTP creates an HTTP request accordingly to requirements of Akamais API
// and execute it returning response object.
func RequestHTTP(method, path string, body io.Reader) (*http.Response, error) {
	var (
		err error
	)

	if isInitialized != true {
		return nil, err
	}

	req, err := prepareHTTPRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	res, err := ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// RequestJSON creates an HTTP request that can be sent to the Akamai APIs with a JSON body
// The JSON body is encoded and the Content-Type/Accept headers are set automatically.
func RequestJSON(method, path string, body interface{}) (*http.Response, error) {
	jsonBody, err := jsonhooks.Marshal(body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(jsonBody)
	req, err := prepareHTTPRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json,*/*")

	res, err := ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ExecuteRequest performs a given HTTP Request, signed with the Akamai OPEN Edgegrid
// Authorization header. An edgegrid.Response or an error is returned.
func ExecuteRequest(req *http.Request) (*http.Response, error) {
	Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req = edgegrid.AddRequestHeader(edgercConfig, req)
		return nil
	}

	req = edgegrid.AddRequestHeader(edgercConfig, req)
	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// BodyJSON unmarshals the Response.Body into a given data structure
func BodyJSON(r *http.Response, data interface{}) error {
	if data == nil {
		return errors.New("You must pass in an interface{}")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = jsonhooks.Unmarshal(body, data)

	return err
}
