// This example deletes your API client credentials.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Add the `credentialId` from the update example to the path. You can only delete inactive credentials. Sending the request on an active set will return a 400. Use the update credentials example for deactivation.
//
// **Important:** Don't use the credentials you're actively using when deleting a set of credentials. Otherwise, you'll block your access to the Akamai APIs.
//
// 3. Open a Terminal or shell instance and run "go run examples/auth-signer/delete/delete-credentials.go".
//
// A successful call returns an empty response body.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/delete-self-credential.

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

	req, err := http.NewRequest(http.MethodDelete, "/identity-management/v3/api-clients/self/credentials/123456", nil)
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
