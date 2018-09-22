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

func waitHandler(c *gin.Context) {
	//random wait time
	t := time.Second * 40
	r := t
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
	c.HTML(http.StatusOK, "wait.tmpl", map[string]interface{}{
		"warningText": "test test test",
		"remaintime":  remaintime,
		"refreshtime": refreshtime,
		"target":      "/wait",
	})
}

//todo link click not count as f5
