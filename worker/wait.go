package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	minRefreshDuration = time.Minute * 1
	maxRefreshDuration = time.Minute * 5
)

func renderWaitPage(c *gin.Context, client clientData) {
	//todo predict wait time
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

func renderErrorPage(c *gin.Context, client clientData) {
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

	server, err := getServerData(client.Server)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	switch server.Status {

	case serverStatusNormal:
		client.Status = clientStatusRelease
		client.saveCookie(c)
		c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
		return

	case serverStatusNotOpen:
		client.Status = clientStatusWait
		client.QTime = spanTime(server.OpenTime)
		client.saveCookie(c)
		renderWaitPage(c, client)
		return

	case serverStatusWaitRoom:
		if client.Status == clientStatusRelease {
			c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
			return
		}
		if client.QTime.Before(server.ReleaseTime) {
			if !client.isValid() {
				//MAC changed at client site
				client = ginContext2NewClient(c)
				client.saveCookie(c)
				c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
				return
			}
			client.Status = clientStatusRelease
			client.ReleaseTime = time.Now()
			client.saveCookie(c)
			c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
			return
		}
		client.Status = clientStatusWait
		client.saveCookie(c)
		renderWaitPage(c, client)
		return
	}

}
