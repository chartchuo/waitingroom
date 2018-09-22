package main

// todo
// persistent database
// health check target host
// admin ui

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const configFile = "config/config.yml"

var confManager *MutexConfigManager

func main() {
	log.Println("Proxy started.")

	go respTimePoller()

	conf, err := loadConfig(configFile)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}
	fmt.Printf("config: %v\n", conf)

	confManager = NewMutexConfigManager(conf)
	watcher, err := WatchFile(configFile, time.Second*5, func() {
		log.Printf("Configfile Updated\n")
		conf, err := loadConfig(configFile)
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else {
			confManager.Set(conf)
			fmt.Printf("config: %v\n", conf)
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

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("html/wait.tmpl")

	r.Any("/", proxyHandler)
	r.GET("/wait", waitHandler)
	r.Run(":8080")

	log.Fatal(err)
}
