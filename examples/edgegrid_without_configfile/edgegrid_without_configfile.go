package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	client := http.Client{}
	config := edgegrid.Config{
		Host:         "xxxxxx.luna.akamaiapis.net",
		ClientToken:  "xxxx-xxxxxxxxxxx-xxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AccessToken:  "xxxx-xxxxxxxxxxx-xxxxxxxxxxx",
		MaxBody:      1024,
		HeaderToSign: []string{
			"X-Test1",
			"X-Test2",
			"X-Test3",
		},
		Debug: false,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/siteshield/v1/maps", config.Host), nil)
	req = edgegrid.AddRequestHeader(config, req)
	resp, _ := client.Do(req)
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
