package main

// todo
// persistent database
// health check target host
// admin ui

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pmylund/go-cache"
)

var BlackList = cache.New(BlockTime, IntervalTime)
var GrayList = cache.New(BlockTime, IntervalTime)

func hardBlock(w http.ResponseWriter) {
	hj, _ := w.(http.Hijacker)
	conn, _, _ := hj.Hijack()
	conn.Close()
}

func normalBlock(w http.ResponseWriter) {
	http.Error(w, "Unauthorized. Block by Chart Web Application Firewall.", http.StatusUnauthorized)
}

func cgDial(network, address string) (net.Conn, error) {
	var d net.Dialer
	//todo validate address
	newAddress, err := TargetAddress(address)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, newAddress)
}

func newClientID() string {
	return time.Now().String() //todo gen unique id
}

func validateClientID(id string) bool {
	return true
}

func proxyRequestOld(w http.ResponseWriter, d *WebInspectData) {
	client := &http.Client{}
	r := d.R
	url := "http://" + r.Host + r.URL.Path

	log.Println(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	if d.RequestBuffered {
		req, _ = http.NewRequest(r.Method, url, bytes.NewReader(d.BodyBuf))
	} else {
		req, _ = http.NewRequest(r.Method, url, d.R.Body)
	}

	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	var transport = &http.Transport{
		Dial: cgDial,
	}

	client.Transport = transport
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func proxyRequest(w http.ResponseWriter, d *WebInspectData) {
	//client := &http.Client{}
	r := d.R
	url := "http://" + r.Host + r.URL.Path

	log.Println(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	if d.RequestBuffered {
		req, _ = http.NewRequest(r.Method, url, bytes.NewReader(d.BodyBuf))
	} else {
		req, _ = http.NewRequest(r.Method, url, d.R.Body)
	}

	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	var transport = &http.Transport{
		Dial: cgDial,
	}

	//client.Transport = transport
	//resp, err := client.Do(req)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	//retrieve client information
	var webdata WebInspectData
	newClient := false
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	webdata.ClientIP = clientIP
	//reqCookie, err := r.Cookie("CloudGuard-ClientId")
	//if err != nil {
	//	newClient = true
	//}
	//if validateClientID(reqCookie.Value){
	//	newClient = true
	//}

	//block if in black list
	_, found := BlackList.Get(clientIP)
	if found {
		log.Println(r.RemoteAddr + " BLOCK")
		normalBlock(w)
		return
	}

	host, err := HostGet(r.Host)

	if err != nil {
		log.Printf("BLOCK RemoteAddress:%v,TargetHost:%v,Msg:%v", r.RemoteAddr, r.Host, err)
		normalBlock(w)
		return
	}

	//*************
	// data prep
	//***********
	webdata.Host = host
	webdata.R = r

	// verify cookie if no cookie generate new
	if newClient {
		cookie := http.Cookie{
			Name:  "CloudGuard-ClientId",
			Value: webdata.ClientID,
		}
		if cookie.Value == "" {
			cookie.Value = newClientID()
		}
		http.SetCookie(w, &cookie)
	}

	//parse data
	if r.Method == "POST" {
		webdata.RequestBuffered = true
		webdata.BodyBuf, _ = ioutil.ReadAll(r.Body)
	}

	// Inspection
	res := Inspect(&webdata)
	if res == INSPECT_ATTACK {
		BlackList.Set(clientIP, 1, cache.DefaultExpiration)
		normalBlock(w)
		return
	}

	proxyRequest(w, &webdata)
}

func init() {
	// applog := &lumberjack.Logger{
	// 	Filename:   "log/app.log",
	// 	MaxSize:    100, // megabytes
	// 	MaxBackups: 30,
	// 	MaxAge:     1, //days
	// 	Compress:   true,
	// }

	// log.SetOutput(applog)
}
func main() {
	log.Println("Proxy started.")
	http.HandleFunc("/", mainHandler)

	// for ubuntu system can't listen port 80 workaround by listen port 8080 instead and NAT with command
	// sudo iptables -t nat -I OUTPUT -p tcp -d 127.0.0.1 --dport 80 -j REDIRECT --to-ports 8080

	err := http.ListenAndServe(":8080", nil)

	log.Fatal(err)
}
