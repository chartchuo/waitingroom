package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const cookieClientIDname = "ccuuid"
const uuidNameSpace = "ccuuid"

func validateClientID(id string) bool {
	_, err := uuid.FromString(id)
	if err != nil {
		return false
	}
	return true
}

func unknowHost(w http.ResponseWriter) {
	http.Error(w, "Unknow Host.", http.StatusUnauthorized)
}

func ccDial(network, address string) (net.Conn, error) {
	var d net.Dialer
	newAddress, err := targetAddress(address)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, newAddress)
}

func newClientID() string {
	return uuid.NewV5(uuid.NamespaceURL, uuidNameSpace).String()
}

var transport = &http.Transport{
	Dial: ccDial,
}

func proxyRequest(w http.ResponseWriter, d *WebInspectData) {
	r := d.R
	url := "http://" + r.Host + r.URL.Path

	//todo log level
	log.Println(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	req, _ = http.NewRequest(r.Method, url, d.R.Body)

	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	startTime := time.Now()

	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	diff := time.Now().Sub(startTime)
	inRespTime <- int(diff / time.Microsecond)

	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	resp.Body.Close()

}

func proxyHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	//retrieve client information
	var webdata WebInspectData
	newClient := false
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	webdata.ClientIP = clientIP

	reqCookie, err := r.Cookie(cookieClientIDname)
	if err != nil {
		newClient = true
	} else if validateClientID(reqCookie.Value) {
		newClient = true
	}

	// verify cookie if no cookie generate new
	if newClient {
		webdata.ClientID = newClientID()
		cookie := http.Cookie{
			Name:  cookieClientIDname,
			Value: webdata.ClientID,
		}
		http.SetCookie(w, &cookie)
	}

	host, err := hostGet(r.Host)
	if err != nil {
		unknowHost(w)
		return
	}
	webdata.Host = host
	webdata.R = r

	proxyRequest(w, &webdata)
}
