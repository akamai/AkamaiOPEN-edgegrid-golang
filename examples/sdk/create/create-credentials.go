// This example creates your new API client credentials.
//
// To run this example:
//
// 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
//
// The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
//
// 2. Open a Terminal or shell instance and run "go run examples/sdk/create/create-credentials.go".
//
// A successful call returns a new API client with its `CredentialID`. Use this ID in both the update and delete examples.
//
// For more information on the call used in this example, see https://techdocs.akamai.com/iam-api/reference/post-self-credentials.

package main

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
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

	resp, err := client.CreateCredential(context.TODO(), iam.CreateCredentialRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
