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

	config := edgegrid.InitConfig("egderc")

	req = edgegrid.AddRequestHeader(config, req)
	resp, _ := client.Do(req)
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
