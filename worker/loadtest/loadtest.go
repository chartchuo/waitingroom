package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func startWebserver() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	go http.ListenAndServe(":8888", nil)

}

func startLoadTest() {
	count := 0
	for {
		resp, err := myClient.Get("http://localhost:8080/") // <-- use a custom client with custom *http.Transport
		if err != nil {
			panic(fmt.Sprintf("Got error: %v", err))
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		log.Printf("Finished GET request #%v", count)
		count++
	}

}

var defaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var myClient *http.Client

func main() {

	myClient = &http.Client{Transport: defaultTransport}

	startWebserver()

	startLoadTest()

}
