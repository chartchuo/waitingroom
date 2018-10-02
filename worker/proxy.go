package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

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

func proxyRequest(c *gin.Context, client clientData, server ServerData) {
	targetAddress, err := host2TargetAddress(client.Server)
	if err != nil {
		log.Errorln(err)
		return
	}

	u, err := url.Parse(targetAddress)
	if err != nil {
		log.Debugln(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = transport

	start := time.Now()
	proxy.ServeHTTP(c.Writer, c.Request)
	duration := time.Since(start)

	var clientC clientChan
	clientC.clientData = client
	clientC.responseTime = int(duration / time.Microsecond)
	inRespTime <- clientC

}

func redirec2EnterWaitingRoom(c *gin.Context) {
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

	server, err := getServerData(client.Server)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	switch server.Status {

	case serverStatusNormal:
		if client.Status == clientStatusRelease {
			client.saveCookie(c)
			proxyRequest(c, client, server)
			return
		}
		if !client.isValid() {
			//expect change mac at client site
			log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
			redirec2EnterWaitingRoom(c)
		}
		client.Status = clientStatusRelease
		client.saveCookie(c)
		proxyRequest(c, client, server)
		return

	case serverStatusNotOpen:
		client.Status = clientStatusWait
		client.saveCookie(c)
		c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
		return

	case serverStatusWaitRoom:
		if client.Status == clientStatusRelease {
			if !client.isValid() {
				//expect change mac at client site
				log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
				redirec2EnterWaitingRoom(c)
			}
			client.saveCookie(c)
			proxyRequest(c, client, server)
			return
		}
		if client.QTime.Before(server.ReleaseTime) {
			if !client.isValid() {
				//expect change mac at client site
				log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
				redirec2EnterWaitingRoom(c)
			}
			client.Status = clientStatusRelease
			client.saveCookie(c)
			proxyRequest(c, client, server)
			return
		}
		client.Status = clientStatusWait
		client.saveCookie(c)
		c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
		return
	}
}

//todo if error to access backen should not response with blank page
