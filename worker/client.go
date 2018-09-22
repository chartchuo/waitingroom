package main

import (
	"time"
)

type clientData struct {
	id           string
	arriveTime   time.Time
	qTime        time.Time
	nextAttemp   time.Time
	lastAccess   time.Time
	refreshCount int
}

func newClientData() {

}

func (c *clientData) ff() {

}
