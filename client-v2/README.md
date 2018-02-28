# Akamai Client

A golang package which helps facilitate making HTTP requests to [Akamai OPEN APIs](https://developer.akamai.com)

## Example of use
```
package main

import (
	"fmt"
	"io/ioutil"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	//TODO: This one does not handle missing sections :(
	config, _ := edgegrid.Init("~/.edgerc", "ingbv")

	client.Init(config)
	resp, _ := client.RequestHTTP("GET", "/diagnostic-tools/v1/locations", nil)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}

```