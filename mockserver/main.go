package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	confManager := NewMutexConfigManager(loadConfig("config/config.yml"))

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
	r.GET("/config", func(c *gin.Context) {
		conf := confManager.Get()
		c.JSON(200, gin.H{
			"message": conf.Message,
		})
	})
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(302, "/destination")
	})
	r.GET("/error", func(c *gin.Context) {
		c.Err()
	})

	watcher, err := WatchFile("config/config.yml", time.Second*5, func() {
		log.Printf("Configfile Updated\n")
		conf := loadConfig("config/config.yml")
		confManager.Set(conf)
	})
	check(err)

	defer func() {
		watcher.Close()
		confManager.Close()
	}()

	r.Run(":8888")

}
