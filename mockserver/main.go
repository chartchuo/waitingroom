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
	r.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "login",
		})
	})
	r.GET("/logout", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "logout",
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
	r.GET("/error", func(c *gin.Context) {
		c.Err()
	})
	r.Run(":8888")
}
