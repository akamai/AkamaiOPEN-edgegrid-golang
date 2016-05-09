package main

import (
	"fmt"
	"github.com/njuettner/edgegrid"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.Client{}

	config := edgegrid.InitConfig("egderc")

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/siteshield/v1/maps", config.Host), nil)
	req = edgegrid.AddRequestHeader(config, req)
	resp, _ := client.Do(req)
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
