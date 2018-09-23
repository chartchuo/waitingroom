package main

import (
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func ccDial(network, address string) (net.Conn, error) {
	var d net.Dialer
	newAddress, err := targetAddress(address)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, newAddress)
}

var transport = &http.Transport{
	Dial: ccDial,
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// r := d.R
	url := "http://" + r.Host + r.URL.Path

	log.Debugln(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	req, _ = http.NewRequest(r.Method, url, r.Body)

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

	client, err := newClientFromC(c)
	if err != nil { //just return. error messge already response in function
		return
	}

	server := getServerData(client.Server)
	switch server.Status {
	case serverStatusNormal:
		//todo authen
		client.Status = clientStatusRelease
		setClientCookie(c, client)
		proxyRequest(c.Writer, c.Request)
		return
	case serverStatusNotOpen:
		client.Status = clientStatusWait
		setClientCookie(c, client)
		c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
		return
	case serverStatusWaitRoom:
		if client.QTime.Before(server.ReleaseTime) {
			//todo authen
			client.Status = clientStatusRelease
			setClientCookie(c, client)
			proxyRequest(c.Writer, c.Request)
			return
		} else {
			client.Status = clientStatusWait
			setClientCookie(c, client)
			c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
			return
		}
	}
}
