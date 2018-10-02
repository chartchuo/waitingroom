package main

import (
	"sync"
	"time"
)

type clientSession map[string]map[string]time.Time

var session clientSession
var avgSesstionTime map[string]int //second

var sessionMutex = &sync.Mutex{}

func init() {
	session = make(clientSession)
	avgSesstionTime = make(map[string]int)
}

const sessioinTimeout = time.Minute * 2
const sessionInterval = 30 //second

func (s clientSession) add(server, id string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	ss, ok := s[server]
	if !ok {
		s[server] = make(map[string]time.Time)
		ss, ok = s[server]
	}
	ss[id] = time.Now()
}

func (s clientSession) concurrent(server string) int {
	return len(s[server])
}

func (s clientSession) avgTime(server string) int {
	return avgSesstionTime[server]
}

func (s clientSession) clearSessionTimeout() {
	for {
		sessionMutex.Lock()
		n := time.Now()
		for server, m := range s {
			count := 0
			sum := 0
			for id, t := range m {
				if t.Add(sessioinTimeout).Before(n) {
					count++
					sum += int(n.Sub(t) / time.Second)
					delete(m, id)
				}
			}
			if count != 0 {
				avgSesstionTime[server] = avgSesstionTime[server] + sum/count //moving average
			}
		}
		sessionMutex.Unlock()
		time.Sleep(time.Second * sessionInterval)
	}
}

func init() {
	go session.clearSessionTimeout()
}
