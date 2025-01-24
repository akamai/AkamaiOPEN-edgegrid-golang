// This example updates the credentials from the create credentials example.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Add the `CredentialID` for the set of credentials created using the create example as a path parameter.
//
// 3. Edit the `ExpiresOn` date to today's date. The date cannot be more than two years out or it will return a 400. Optionally, you can change the `Description` value.
//
// **Important:** Don't use the credentials you're actively using when inactivating a set of credentials. Otherwise, you'll block your access to the Akamai APIs.
//
// 4. Open a Terminal or shell instance and run "go run examples/sdk/update/update-credentials.go".
//
// A successful call returns an object with modified credentials.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/put-self-credential.

package main

import (
	"context"
	"fmt"
	"time"

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

	resp, err := client.UpdateCredential(
		context.TODO(),
		iam.UpdateCredentialRequest{
			CredentialID: 123456,
			Body: iam.UpdateCredentialRequestBody{
				Description: "Update this credential",
				ExpiresOn:   time.Date(2025, 05, 6, 11, 45, 04, 0, time.UTC),
				Status:      iam.CredentialInactive,
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
