package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func ccDial(network, address string) (net.Conn, error) {
	var d net.Dialer
	newAddress, err := getTargetAddress(address)
	if err != nil {
		return nil, err
	}
	d.Timeout = time.Second * 5
	return d.Dial(network, newAddress)
}

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxIdleConnsPerHost:   50,
	MaxConnsPerHost:       200,
}

func proxyRequest(c *gin.Context) {
	targetAddress, err := getTargetAddress(c.Request.Host)
	if err != nil {
		return
	}

	u, err := url.Parse(targetAddress)
	if err != nil {
		log.Debugln(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = transport
	proxy.ServeHTTP(c.Writer, c.Request)
}

func proxyRequest2(c *gin.Context) {

	w := c.Writer
	r := c.Request
	// r := d.R
	url := "http://" + r.Host + r.URL.Path

	// log.Debugln(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	req, _ = http.NewRequest(r.Method, url, r.Body)

	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	startTime := time.Now()

	resp, err := transport.RoundTrip(req)
	if err != nil {
		renderErrorPage(c)
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

func redirec2WaitingRoom(c *gin.Context) {
	host, err := getHost(c.Request.Host)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "unknow host." + c.Request.Host,
		})
		log.Errorln("unknow host:", c.Request.Host)
		return
	}
	client := newClientData(host)
	client.saveCookie(c)
	c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
}

func proxyHandler(c *gin.Context) {

	client, err := ginContext2Client(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	server := getServerData(client.Server)

	switch server.Status {

	case serverStatusNormal:
		if client.Status == clientStatusRelease {
			proxyRequest(c)
			return
		}
		if !client.isValid() {
			//expect change mac at client site
			log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
			redirec2WaitingRoom(c)
		}
		client.Status = clientStatusRelease
		client.saveCookie(c)
		proxyRequest(c)
		return

	case serverStatusNotOpen:
		redirec2WaitingRoom(c)
		return

	case serverStatusWaitRoom:
		if client.Status == clientStatusRelease {
			if !client.isValid() {
				//expect change mac at client site
				log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
				redirec2WaitingRoom(c)
			}
			proxyRequest(c)
			return
		}
		if client.QTime.Before(server.ReleaseTime) {
			if !client.isValid() {
				//expect change mac at client site
				log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
				redirec2WaitingRoom(c)
			}
			client.Status = clientStatusRelease
			client.saveCookie(c)
			proxyRequest(c)
			return
		}
		client.Status = clientStatusWait
		client.saveCookie(c)
		c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
		return
	}
}
