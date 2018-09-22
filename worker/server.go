package main

import "time"

//ServerConfig for save and load from file
type ServerConfig struct {
	OpenTime time.Time `yaml:"opentime,omitempty"`
	MaxUsers int       `yaml:"maxusers,omitempty"`
}

//ServerData dynamic server data
type ServerData struct {
	Name         string
	ReleaseTime  time.Time
	MaxUsers     int
	CurrentUsers int //todo on local proxy instant
}
