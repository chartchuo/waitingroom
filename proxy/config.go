package main

import "time"

//DDOS
const IntervalTime = time.Second * 10 //interval
const BucketSize = 30
const LeakRate = 10 // request per refreshtime
const BucketTimeout = time.Minute

const BlockTime = time.Second * 15

type Config struct {
	HostDB   map[string]string `yaml:"hostdb,omitempty`
	TargetDB map[string]string `yaml:"targetdb,omitempty`
}
