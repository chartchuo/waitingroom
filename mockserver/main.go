package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Mock server home",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/destination", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "destination",
		})
	})
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(302, "/destination")
	})
	r.Run(":8888")
}
