package main

import (
	"errors"
	"sync"
	"time"
)

const serverInterval = 1 //second

//ServerConfig for save and load from file
type ServerConfig struct {
	OpenTime        time.Time `yaml:"opentime,omitempty"`
	BordingTime     int       `yaml:"bordingtime,omitempty"`
	MaxUsers        int       `yaml:"maxusers,omitempty"`
	MaxResponseTime int       `yaml:"maxresponsetime,omitempty"`
}
type serverStatus int

const (
	serverStatusNormal serverStatus = iota
	serverStatusNotOpen
	serverStatusWaitRoom
)

const defaultMaxUsers = 100

type serverCounter struct {
	count           int
	sum             int
	p95             []int
	concurrentusers int
	maxresponsetime int
}

//ServerData dynamic server data
type ServerData struct {
	Status          serverStatus
	ReleaseTime     time.Time
	OpenTime        time.Time
	BordingTime     int //minutes
	MaxUsers        int
	MaxResponseTime int //milisecs
	counter         *serverCounter
}

var serverdataDB map[string]ServerData
var serverdataMutex = &sync.Mutex{}

func init() {
	serverdataDB = make(map[string]ServerData)
}

func newServerData(name string) (ServerData, error) {
	serverdataMutex.Lock()
	defer serverdataMutex.Unlock()
	c := confManager.Get()
	open := c.ServerConfig[name].OpenTime
	max := c.ServerConfig[name].MaxUsers
	bording := c.ServerConfig[name].BordingTime
	maxrestime := c.ServerConfig[name].MaxResponseTime
	if max == 0 {
		max = defaultMaxUsers
	}
	var status serverStatus
	if open.After(time.Now()) {
		status = serverStatusNotOpen
	} else {
		status = serverStatusWaitRoom
	}

	if bording == 0 {
		bording = 60 //minutes
	}

	var release time.Time
	release = open
	if status == serverStatusNormal {
		release = time.Now()
	}

	if appRunMode == "debug" {
		open = time.Now().Add(time.Second * 10)
		status = serverStatusNotOpen
		// release = open
		// max = 10

		// open = time.Now()
		// status = serverStatusNormal
		// release = open
		// max = 10

	}
	serverdataDB[name] = ServerData{
		Status:          status,
		ReleaseTime:     release,
		OpenTime:        open,
		BordingTime:     bording,
		MaxUsers:        max,
		MaxResponseTime: maxrestime,
		counter: &serverCounter{
			p95: make([]int, 0, p95cap),
		},
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

func initServerData() {
	c := confManager.Get()
	for k := range c.ServerConfig {
		newServerData(k)
	}
}
