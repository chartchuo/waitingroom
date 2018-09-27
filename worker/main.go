package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const appRunMode = "debug"

const configFile = "config/config.yml"
const waitRoomPath = "/ccwait"
const serverEntryPath = "/"

var confManager *MutexConfigManager

func main() {
	log.Println("Woker started.")

	conf, err := loadConfig(configFile)
	if err != nil {
		log.Errorf("ERROR: %v\n", err)
	}
	log.Debugf("config: %v\n", conf)

	confManager = NewMutexConfigManager(conf)
	watcher, err := WatchFile(configFile, time.Second*5, func() {
		log.Printf("Configfile Updated\n")
		conf, err := loadConfig(configFile)
		if err != nil {
			log.Errorf("ERROR: %v", err)
		} else {
			confManager.Set(conf)
			log.Debugf("config: %v\n", conf)
		}
	})
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	defer func() {
		watcher.Close()
		confManager.Close()
	}()

	go advisorPoller()

	if appRunMode == "debug" {
		serverinit() //mock data
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.ErrorLevel)
	}
	r := gin.Default()

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("tmpl/wait.tmpl", "tmpl/error.tmpl") //gin limitation: must add multiple file in one command

	r.Any("/*root", proxyHandler)
	// r.GET(waitRoomPath, waitHandler)

	// client := newClientData("test")
	// if client.isValid() {
	// 	log.Debugln("valid")
	// } else {
	// 	log.Debugln("invalid")
	// }

	r.Run(":8080")

	log.Fatal(err)
}
