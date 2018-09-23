package main

import (
	"time"

	log "github.com/sirupsen/logrus"
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
	MaxUsers     int
	CurrentUsers int //todo on local proxy instant only not implement cluster solution yet
}

var serverdataDB map[string]ServerData

func serverinit() {
	serverdataDB = make(map[string]ServerData)
	c := confManager.Get()
	configMock := c.ServerConfig["mock"]
	configMock.OpenTime = time.Now().Add(time.Minute)
	confManager.Set(c)
	serverdataDB["mock"] = ServerData{
		Status:       serverStatusNotOpen,
		ReleaseTime:  c.ServerConfig["mock"].OpenTime.Add(time.Minute * 2),
		MaxUsers:     100,
		CurrentUsers: 100,
	}
}

func getServerData(name string) ServerData {
	s, ok := serverdataDB[name]
	if !ok {
		log.Error("error server.go getServerData ", name)
	}
	return s
}
