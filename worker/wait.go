package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	minRefreshDuration = time.Minute * 1
	maxRefreshDuration = time.Minute * 5
)

func renderWaitPage(c *gin.Context) {
	//todo mock random wait time
	t := time.Second * 40
	r := t / 2
	if r < minRefreshDuration {
		r = minRefreshDuration
	} else if r > maxRefreshDuration {
		r = maxRefreshDuration
	}
	if r > t {
		t = r
	}
	remaintime := int(t / time.Millisecond)
	refreshtime := int(r / time.Millisecond)
	msg := ""
	warn := ""

	c.HTML(http.StatusOK, "wait.tmpl", map[string]interface{}{
		"warningText": warn,
		"remaintime":  remaintime,
		"refreshtime": refreshtime,
		"target":      waitRoomPath,
		"message":     msg,
	})
}

func renderErrorPage(c *gin.Context) {
	t := time.Second * 40
	r := t / 2
	if r < minRefreshDuration {
		r = minRefreshDuration
	} else if r > maxRefreshDuration {
		r = maxRefreshDuration
	}
	if r > t {
		t = r
	}
	remaintime := int(t / time.Millisecond)
	refreshtime := int(r / time.Millisecond)
	msg := ""
	warn := ""

	c.HTML(http.StatusOK, "error.tmpl", map[string]interface{}{
		"warningText": warn,
		"remaintime":  remaintime,
		"refreshtime": refreshtime,
		"target":      waitRoomPath,
		"message":     msg,
	})
}

func waitHandler(c *gin.Context) {

	client, err := ginContext2Client(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	server := getServerData(client.Server)

	switch server.Status {

	case serverStatusNormal:
		client.Status = clientStatusRelease
		client.saveCookie(c)
		c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
		return

	case serverStatusNotOpen:
		client.saveCookie(c)
		renderWaitPage(c)
		return

	case serverStatusWaitRoom:
		if client.Status == clientStatusRelease {
			c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
			return
		}
		if client.QTime.Before(server.ReleaseTime) {
			if !client.isValid() {
				//expect change mac at client site
				log.Infoln("invalid MAC detect remote ip: ", c.Request.RemoteAddr)
				host, err := hostGet(c.Request.Host)
				if err != nil {
					c.JSON(200, gin.H{
						"message": "unknow host." + c.Request.Host,
					})
					log.Errorln("unknow host:", c.Request.Host)
					return
				}
				client = newClientData(host)
				client.saveCookie(c)
				c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
			}
			client.Status = clientStatusRelease
			client.saveCookie(c)
			c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
			return
		}
		client.Status = clientStatusWait
		client.saveCookie(c)
		renderWaitPage(c)
		return
	}

}

//todo
//link click not count as f5
//re queue for f5 user
