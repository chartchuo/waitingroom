package main

import (
	"errors"
	"sync"
	"time"
)

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
	// c := confManager.Get()

	// log.Debug("name: " + name)
	serverdataDB[name] = ServerData{
		Status: serverStatusWaitRoom,
		// Status: serverStatusNotOpen,
		// Status:      serverStatusNormal,
		ReleaseTime: time.Now().Add(time.Minute * 2),
		OpenTime:    time.Now().Add(time.Minute * 2), //todo read from config
		MaxUsers:    10,                              //todo read from config
		// CurrentUsers: 100,
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
