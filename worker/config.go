package main

import "time"

//DDOS
const IntervalTime = time.Second * 10 //IntervalTime to run leaky bucket
const BucketSize = 30                 //BucketSize total number of packet in bucket
const LeakRate = 10                   //LeakRate  request per interval
const BucketTimeout = time.Minute     //BucketTimeout time to remove from cache

const BlockTime = time.Second * 15

//Config main repository
type Config struct {
	HostDB   map[string]string `yaml:"hostdb,omitempty"`
	TargetDB map[string]string `yaml:"targetdb,omitempty"`
}
