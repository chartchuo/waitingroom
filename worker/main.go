package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	startAdvisor()

	if appRunMode == "debug" {
		// serverinit() //mock data
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.ErrorLevel)
	}
	r := gin.Default()

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("tmpl/wait.tmpl", "tmpl/error.tmpl") //gin limitation: must add multiple file in one command

	// r.GET(waitRoomPath, waitHandler)
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// r.Any("/", proxyHandler)

	r.Any("/*path", func(c *gin.Context) {
		path := c.Param("path")
		if path == waitRoomPath {
			waitHandler(c)
			return
		} else if path == "/metrics" {
			h := gin.WrapH(promhttp.Handler())
			h(c)
			return
		}
		proxyHandler(c)
	})

	r.Run(":8080")

	log.Fatal(err)
}
