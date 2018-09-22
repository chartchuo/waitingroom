package main

import (
	"fmt"
	"math/rand"
	"time"

	cache "github.com/pmylund/go-cache"
)

const qSpanTime = time.Minute //one minute

type clientData struct {
	server       string
	arriveTime   time.Time
	qTime        time.Time
	nextAttemp   time.Time
	lastAccess   time.Time
	refreshCount int
	mac          string
}

var clientCache *cache.Cache

func init() {
	clientCache = cache.New(5*time.Minute, 10*time.Minute)

}
func newClientData(server string) clientData {
	a := time.Now()
	c := confManager.Get().ServerConfig[server]
	// s := serverdata[server]
	var q time.Time
	if a.Before(c.OpenTime) {
		r := rand.Int63n(int64(qSpanTime))
		q = c.OpenTime.Add(time.Duration(r)) //span from open to end of span
	} else if a.Before(c.OpenTime.Add(qSpanTime)) {
		d := c.OpenTime.Add(qSpanTime).Sub(a) //time remain bedore end of span
		r := rand.Int63n(int64(d))
		q = a.Add(time.Duration(r)) //span from arrive to end of span
	} else {
		q = a
	}
	fmt.Println(c.OpenTime)
	return clientData{
		server:       server,
		arriveTime:   a,
		qTime:        q,
		nextAttemp:   q,
		lastAccess:   a,
		refreshCount: 0,
	}
}
