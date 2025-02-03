# Examples

This directory contains executable CRUD examples for Akamai API using the EdgeGrid Go library as an auth signer and SDK. API calls used in these examples are available to all users. But, if you find one of the write examples doesn't work for you, talk with your account's admin about your privilege level.

## Run

To run any of the files:

1. Specify the location of your `.edgerc` file and the section header for the set of credentials you'd like to use. The default is `default`.
2. For update and delete operations, replace the dummy `credentialId` with your valid `credentialId`.
   
   >**Important:** Don't use the credentials you're actively using when running the update (inactivation) and delete operations. Otherwise, you'll block your access to the Akamai APIs.

3. Open a Terminal or shell instance and run the .go file.

    ```
    $ go run examples/<subdirectory-name>/<operation-name>/<file-name>.go
    ```

## Sample files

The example in each file contains a call to one of the Identity and Access Management (IAM) API endpoints. See the [IAM API reference](https://techdocs.akamai.com/iam-api/reference/api) doc for more information on each of the calls used.

<table>
    <tr>
        <th>Operation</th>
        <th>Method</th>
        <th>Endpoint</th>
    </tr>
    <tr>
        <td>
            List your API client credentials.
            <table>
                <tr>
                    <td style="border: none;"><a href="./auth-signer/get/get-credentials.go">Run as auth signer</a></td>
                    <td style="border: none;"><a href="./sdk/get/get-credentials.go">Run as SDK</a></td>
                </tr>
            </table>
        </td>
        <td><code>GET</code></td>
        <td><code>/identity-management/v3/api-clients/self/credentials</code></td>
    </tr>
    <tr>
        <td>
            Create new API client credentials. <br> This is a <i>quick</i> client and grants you the default permissions associated with your account.
            <table>
                <tr>
                    <td style="border: none;"><a href="./auth-signer/create/create-credentials.go">Run as auth signer</a></td>
                    <td style="border: none;"><a href="./sdk/create/create-credentials.go">Run as SDK</a></td>
                </tr>
            </table>
        </td>
        <td><code>POST</code></td>
        <td><code>/identity-management/v3/api-clients/self/credentials</code></td>
    </tr>
    <tr>
        <td>
            Update your credentials by ID.
            <table>
                <tr>
                    <td style="border: none;"><a href="./auth-signer/update/update-credentials.go">Run as auth signer</a></td>
                    <td style="border: none;"><a href="./sdk/update/update-credentials.go">Run as SDK</a></td>
                </tr>
            </table>
        </td>
        <td><code>PUT</code></td>
        <td><code>/identity-management/v3/api-clients/self/credentials/{credentialId}</code></td>
    </tr>
    <tr>
        <td>
            Delete your credentials by ID.
            <table>
                <tr>
                    <td style="border: none;"><a href="./auth-signer/delete/delete-credentials.go">Run as auth signer</a></td>
                    <td style="border: none;"><a href="./sdk/delete/delete-credentials.go">Run as SDK</a></td>
                </tr>
            </table>
        </td>
        <td><code>DELETE</code></td>
        <td><code>/identity-management/v3/api-clients/self/credentials/{credentialId}</code></td>
    </tr>
</table>

Suggested chained call order:

1. Get credentials to see your base information.
2. Create a client to create a new set of credentials.
3. Update credentials to inactivate the newly created set from step 2.
4. Delete a client to delete the inactivated credentials.
5. Get credentials to verify if they're gone (the status will be `DELETED`).