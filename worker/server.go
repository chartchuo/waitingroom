package main

import (
	"errors"
	"sync"
	"time"
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

type serverCounter struct {
	count        int
	sum          int
	p95          []int
	currentUsers int
}

//ServerData dynamic server data
type ServerData struct {
	Status      serverStatus
	ReleaseTime time.Time
	OpenTime    time.Time
	Bording     int //minute
	MaxUsers    int
	counter     *serverCounter
}

var serverdataDB map[string]ServerData
var serverdataMutex = &sync.Mutex{}

func init() {
	serverdataDB = make(map[string]ServerData)
}

func newServerData(name string) (ServerData, error) {
	// s, ok := serverdataDB[name] //check lock check
	// if ok {
	// 	return s, nil
	// }
	serverdataMutex.Lock()
	defer serverdataMutex.Unlock()
	// s2, ok := serverdataDB[name]
	// if ok {
	// 	return s2, nil
	// }
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
		Status:      status,
		ReleaseTime: release,
		OpenTime:    open,
		MaxUsers:    max,
		counter: &serverCounter{
			count:        0,
			sum:          0,
			p95:          make([]int, 0, p95cap),
			currentUsers: 0,
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
