package main

import (
	"fmt"
	"github.com/njuettner/edgegrid"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.Client{}
	baseURL := "https://xxxxxx.luna.akamaiapis.net"
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/siteshield/v1/maps", baseURL), nil)
	config := edgegrid.Config{
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

	req = edgegrid.AddRequestHeader(config, req)
	resp, _ := client.Do(req)
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
