// This example deletes your API client credentials.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Add the `CredentialID` from the update example to the path. You can only delete inactive credentials. Sending the request on an active set will return a 400. Use the update credentials example for deactivation.
//
// **Important:** Don't use your actual credentials for this operation. Otherwise, you'll block your access to the Akamai APIs.
//
// 3. Open a Terminal or shell instance and run "go run examples/sdk/delete/delete-credentials.go".
//
// A successful call returns an empty response body.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/delete-self-credential.

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

	resp := client.DeleteCredential(context.TODO(), iam.DeleteCredentialRequest{CredentialID: 123456})

	fmt.Println(resp)
}
