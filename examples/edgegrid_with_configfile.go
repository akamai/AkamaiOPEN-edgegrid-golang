package main

import (
	"fmt"
	"github.com/akamai-open/AkamaiOPEN-edgegrid-golang"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.Client{}

	config := edgegrid.InitConfig("~/.egderc", "default")

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/diagnostic-tools/v1/locations", config.Host), nil)
	req = edgegrid.AddRequestHeader(config, req)
	resp, _ := client.Do(req)
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
