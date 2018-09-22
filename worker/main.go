package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const configFile = "config/config.yml"
const waitRoomPath = "/ccwait"

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

	serverinit() //todo mock data must be remove

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("html/wait.tmpl")

	r.Any("/", proxyHandler)
	r.GET(waitRoomPath, waitHandler)

	r.Run(":8080")

	log.Fatal(err)
}
