package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const qSpanTime = time.Minute //one minute
const cookieName = "ccwait"
const uuidNameSpace = "ccwait"
const key = "D3NRX?uVtbJEq_HHLQ5Y"

type clientData struct {
	ID           string
	Server       string
	ArriveTime   time.Time
	QTime        time.Time
	NextAttemp   time.Time
	LastAccess   time.Time
	ReleaseTime  time.Time
	RefreshCount int
	MAC          string //base on ID,QTime and secret for pass to main system only
	Authen       string //base on ID,ReleaseTime and secret for check every request to main system
}

func newClientID() string {
	return uuid.NewV5(uuid.NamespaceURL, uuidNameSpace).String()
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
	log.Debugln(c.OpenTime)
	client := clientData{
		ID:           newClientID(),
		Server:       server,
		ArriveTime:   a,
		QTime:        q,
		NextAttemp:   q,
		LastAccess:   a,
		RefreshCount: 0,
	}
	client.genMAC()
	return client
}

func (c *clientData) isValid() bool {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(c.ID))
	b, _ := c.QTime.MarshalBinary()
	mac.Write(b)
	return c.MAC == hex.EncodeToString(mac.Sum(nil))
}

func (c *clientData) genMAC() {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(c.ID))
	b, _ := c.QTime.MarshalBinary()
	mac.Write(b)
	c.MAC = hex.EncodeToString(mac.Sum(nil))
}
