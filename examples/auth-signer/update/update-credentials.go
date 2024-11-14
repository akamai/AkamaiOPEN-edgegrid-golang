// This example updates the credentials from the create credentials example.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Add the `credentialId` for the set of credentials created using the create example as a path parameter.
//
// 3. Edit the `expiresOn` date to today's date. The date cannot be more than two years out or it will return a 400. Optionally, you can change the `description` value.
//
// **Important:** Don't use your actual credentials when inactivating them. Otherwise, you'll block your access to the Akamai APIs.
//
// 4. Open a Terminal or shell instance and run "go run examples/auth-signer/update/update-credentials.go".
//
// A successful call returns an object with modified credentials.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/put-self-credential.

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
        "description": "Update this credential",
        "expiresOn": "2025-12-10T23:06:59.000Z",
        "status": "INACTIVE"
      }`)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPut, "/identity-management/v3/api-clients/self/credentials/123456", payload)
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
