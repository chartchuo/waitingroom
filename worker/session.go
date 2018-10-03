package main

import (
	"sync"
	"time"
)

type clientSessionElement struct {
	arrive     time.Time
	lastAccess time.Time
}

type clientSession map[string]map[string]clientSessionElement

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
		s[server] = make(map[string]clientSessionElement)
		ss, ok = s[server]
	}

	e, ok := ss[id]
	if ok {
		ss[id] = clientSessionElement{arrive: e.arrive, lastAccess: time.Now()}
		return
	}
	ss[id] = clientSessionElement{arrive: time.Now(), lastAccess: time.Now()}
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
				if t.lastAccess.Add(sessioinTimeout).Before(n) {
					count++
					sum += int(n.Add(-sessioinTimeout).Sub(t.arrive) / time.Second)
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
