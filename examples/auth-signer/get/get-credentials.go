// This example returns a list of your API client credentials.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Open a Terminal or shell instance and run "go run examples/auth-signer/get/get-credentials.go".
//
// A successful call returns your credentials grouped by `credentialId`.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/get-self-credentials.

package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
)

func main() {
	edgerc, err := edgegrid.New(
		edgegrid.WithFile("~/.edgerc"),
		edgegrid.WithSection("default"),
	)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "/identity-management/v3/api-clients/self/credentials", nil)
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
