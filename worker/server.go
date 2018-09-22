package main

import (
	"time"
)

//ServerConfig for save and load from file
type ServerConfig struct {
	OpenTime time.Time `yaml:"opentime,omitempty"`
	MaxUsers int       `yaml:"maxusers,omitempty"`
}

//ServerData dynamic server data
type ServerData struct {
	ReleaseTime  time.Time
	MaxUsers     int
	CurrentUsers int //todo on local proxy instant only not implement cluster solution yet
}

var serverdata map[string]ServerData

func serverinit() {
	serverdata = make(map[string]ServerData)
	c := confManager.Get()
	serverdata["mock"] = ServerData{
		ReleaseTime:  c.ServerConfig["mock"].OpenTime.Add(time.Minute),
		MaxUsers:     100,
		CurrentUsers: 100,
	}
}
