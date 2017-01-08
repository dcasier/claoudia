package main

import (
	"github.com/dcasier/claoudia"
	"github.com/dcasier/claoudia/metamodel"
	"fmt"
	"net/http"
	"net"
	"time"
	"encoding/json"
)

func main() {

	config := common.MustGetTlsConfiguration()

	tr := &http.Transport{
		TLSClientConfig: config,
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,	
	}
	
	c := &http.Client{
		Timeout: time.Second * 10,
		Transport: tr}
	resp, _ := c.Get("https://localhost:9443")
	
	defer resp.Body.Close()
	
	activities := new(metamodel.Activities)
	json.NewDecoder(resp.Body).Decode(activities)
	fmt.Println(activities)
	
}

