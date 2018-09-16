package main

// todo
// persistent database
// health check target host
// admin ui

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const configFile = "config/config.yml"

var confManager *MutexConfigManager

func main() {
	log.Println("Proxy started.")

	// go respTimePoller()

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

	http.HandleFunc("/", mainHandler)

	// for ubuntu system can't listen port 80 workaround by listen port 8080 instead and NAT with command
	// sudo iptables -t nat -I OUTPUT -p tcp -d 127.0.0.1 --dport 80 -j REDIRECT --to-ports 8080

	err = http.ListenAndServe(":8080", nil)

	log.Fatal(err)
}
