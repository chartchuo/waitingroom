package main

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const serverInterval = 1 //second

//ServerConfig for save and load from file
type ServerConfig struct {
	OpenTime time.Time `yaml:"opentime,omitempty"`
	MaxUsers int       `yaml:"maxusers,omitempty"`
}
type serverStatus int

const (
	serverStatusNormal serverStatus = iota
	serverStatusNotOpen
	serverStatusWaitRoom
)

const defaultMaxUsers = 100

//ServerData dynamic server data
type ServerData struct {
	Status       serverStatus
	ReleaseTime  time.Time
	OpenTime     time.Time
	MaxUsers     int
	CurrentUsers int //todo on local proxy instant only not implement cluster solution yet
}

var serverdataDB map[string]ServerData
var serverdataMutex = &sync.Mutex{}

func init() {
	serverdataDB = make(map[string]ServerData)
}

// func serverinit2() {
// 	serverdataDB = make(map[string]ServerData)
// 	c := confManager.Get()
// 	configMock := c.ServerConfig["mock"]
// 	configMock.OpenTime = time.Now().Add(time.Minute)
// 	confManager.Set(c)

// 	serverdataDB["mock"] = ServerData{
// 		// Status: serverStatusWaitRoom,
// 		// Status: serverStatusNotOpen,
// 		Status:      serverStatusNormal,
// 		ReleaseTime: c.ServerConfig["mock"].OpenTime.Add(time.Minute * 2),
// 		MaxUsers:    1,
// 		// CurrentUsers: 100,
// 	}
// }

func newServerData(name string) (ServerData, error) {
	s, ok := serverdataDB[name] //check lock check
	if ok {
		return s, nil
	}
	serverdataMutex.Lock()
	defer serverdataMutex.Unlock()
	s2, ok := serverdataDB[name]
	if ok {
		return s2, nil
	}
	c := confManager.Get()
	open := c.ServerConfig[name].OpenTime
	max := c.ServerConfig[name].MaxUsers
	if max == 0 {
		max = defaultMaxUsers
	}
	var status serverStatus
	if open.After(time.Now()) {
		status = serverStatusNotOpen
	} else {
		status = serverStatusWaitRoom
	}
	var release time.Time
	switch status {
	case serverStatusNormal:
		release = time.Now()
	case serverStatusNotOpen:
		release = open
	case serverStatusWaitRoom:
		release = open
	}

	if appRunMode == "debug" {
		open = time.Now().Add(time.Second * 10)
		status = serverStatusNotOpen
		release = open
		max = 10
	}
	serverdataDB[name] = ServerData{
		Status:      status,
		ReleaseTime: release,
		OpenTime:    open,
		MaxUsers:    max,
	}
	s3, ok := serverdataDB[name]
	if ok {
		return s3, nil
	}
	return ServerData{}, errors.New("error server.go newServerData() " + name)
}

func getServerData(name string) (ServerData, error) {
	s, ok := serverdataDB[name]
	if !ok {
		return newServerData(name)
	}
	return s, nil
}

func setServerData(name string, s ServerData) {
	serverdataMutex.Lock()
	defer serverdataMutex.Unlock()
	serverdataDB[name] = s
}

func startServerJobsOpen() {
	for {
		for k, v := range serverdataDB {
			switch v.Status {
			case serverStatusNotOpen:
				if v.OpenTime.Before(time.Now()) {
					v.Status = serverStatusWaitRoom
					v.ReleaseTime = v.OpenTime
					setServerData(k, v)
					log.Debug("Open server: ", k)
				}
			}

		}
		time.Sleep(serverInterval * time.Second)
	}
}

func initServerData() {
	c := confManager.Get()
	for k := range c.ServerConfig {
		newServerData(k)
	}
}

func startServerJobs() {
	initServerData()
	go startServerJobsOpen()
}
