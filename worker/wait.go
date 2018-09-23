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

	c.HTML(http.StatusOK, "wait.tmpl", map[string]interface{}{
		"warningText": "test test test",
		"remaintime":  remaintime,
		"refreshtime": refreshtime,
		"target":      waitRoomPath,
		"message":     msg,
	})
}

func waitHandler(c *gin.Context) {
	//todo if q time add MAC2 redirect to main site
	client, err := newClientFromC(c)
	if err != nil { //just return. error messge already response in function
		return
	}

	server := getServerData(client.Server)
	switch server.Status {
	case serverStatusNormal:
		//todo add authen mac
		client.Status = clientStatusRelease
		setClientCookie(c, client)
		c.Redirect(http.StatusTemporaryRedirect, serverEntryPath)
		return
	case serverStatusNotOpen:
		renderWaitPage(c)
		return
	case serverStatusWaitRoom:
		if client.QTime.Before(server.ReleaseTime) {
			// client.Status = clientStatusRelease
			// setClientCookie(c, client)
			// proxyRequest(c.Writer, c.Request)
			return
		} else {
			// client.Status = clientStatusWait
			// setClientCookie(c, client)
			// c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
			return
		}
	}

}

//todo
//link click not count as f5
//re queue for f5 user
