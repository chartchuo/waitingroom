package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const configFile = "config/config.yml"

func main() {
	conf, err := loadConfig(configFile)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	confManager := NewMutexConfigManager(conf)
	watcher, err := WatchFile(configFile, time.Second*5, func() {
		log.Printf("Configfile Updated\n")
		conf, err := loadConfig(configFile)
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else {
			confManager.Set(conf)
		}
	})
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	defer func() {
		watcher.Close()
		confManager.Close()
	}()

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

	r.Run(":8888")

}
