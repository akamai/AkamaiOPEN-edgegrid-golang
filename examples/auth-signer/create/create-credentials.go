// This example creates your new API client credentials.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Open a Terminal or shell instance and run "go run examples/auth-signer/create/create-credentials.go".
//
// A successful call returns a new API client with its `credentialId`. Use this ID in both the update and delete examples.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/post-self-credentials.

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

	req, err := http.NewRequest(http.MethodPost, "/identity-management/v3/api-clients/self/credentials", nil)
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
