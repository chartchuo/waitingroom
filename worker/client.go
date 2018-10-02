package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const qSpanTime = time.Minute //Spantime when enter the queue will random in open to open+spantime
const cookieName = "ccwait"
const key = "D3NRX?uVtbJEq_HHLQ5Y"

type clientStatus int

const (
	clientStatusNew clientStatus = iota
	clientStatusWait
	clientStatusRelease
)

type clientData struct {
	ID           string
	Status       clientStatus
	Server       string
	ArriveTime   time.Time
	QTime        time.Time
	NextAttemp   time.Time
	LastAccess   time.Time
	ReleaseTime  time.Time
	RefreshCount int
	MAC          string //base on ID,Status,QTime and secret
}

type clientDataCookie struct {
	ID           string
	Status       clientStatus
	Server       string
	ArriveTime   int64
	QTime        int64
	NextAttemp   int64
	LastAccess   int64
	ReleaseTime  int64
	RefreshCount int
	MAC          string //base on ID,Status,QTime and secret
}

func newClientID() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func spanTime(opentime time.Time) time.Time {
	now := time.Now().UnixNano()
	a := time.Unix(0, now)
	var q time.Time
	if a.Before(opentime) {
		r := rand.Int63n(int64(qSpanTime))
		q = opentime.Add(time.Duration(r)) //span from open to end of span
	} else if a.Before(opentime.Add(qSpanTime)) {
		d := opentime.Add(qSpanTime).Sub(a) //time remain bedore end of span
		r := rand.Int63n(int64(d))
		q = a.Add(time.Duration(r)) //span from arrive to end of span
	} else {
		q = a
	}
	return q
}

func newClientData(server string) clientData {
	now := time.Now().UnixNano()
	a := time.Unix(0, now) //workaround remove + after time from time.Now()
	// c := confManager.Get().ServerConfig[server]
	// q := spanTime(c.OpenTime)
	client := clientData{
		ID:           newClientID(),
		Status:       clientStatusNew,
		Server:       server,
		ArriveTime:   a,
		QTime:        a,
		LastAccess:   a,
		NextAttemp:   time.Unix(0, 0),
		ReleaseTime:  time.Unix(0, 0),
		RefreshCount: 0,
	}
	client.genMAC()
	return client
}

func (client *clientData) isValid() bool {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(client.ID))                           //client ID
	mac.Write([]byte(fmt.Sprint(client.Status)))           //client status
	b, _ := client.QTime.MarshalBinary()                   //Qtime
	mac.Write(b)                                           //Qtime
	mac.Write([]byte(fmt.Sprint(client.QTime.UnixNano()))) //Queue time nano
	return client.MAC == hex.EncodeToString(mac.Sum(nil))
}

func (client *clientData) genMAC() {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(client.ID))                           //client ID
	mac.Write([]byte(fmt.Sprint(client.Status)))           //client status
	b, _ := client.QTime.MarshalBinary()                   //Qtime
	mac.Write(b)                                           //Qtime
	mac.Write([]byte(fmt.Sprint(client.QTime.UnixNano()))) //Queue time nano
	client.MAC = hex.EncodeToString(mac.Sum(nil))
}
func (client *clientData) toCookie() clientDataCookie {
	return clientDataCookie{
		ID:           client.ID,
		Status:       client.Status,
		Server:       client.Server,
		ArriveTime:   client.ArriveTime.UnixNano(),
		QTime:        client.QTime.UnixNano(),
		NextAttemp:   client.NextAttemp.UnixNano(),
		LastAccess:   client.LastAccess.UnixNano(),
		ReleaseTime:  client.ReleaseTime.UnixNano(),
		RefreshCount: client.RefreshCount,
		MAC:          client.MAC,
	}
}
func (c *clientDataCookie) toClient() clientData {
	return clientData{
		ID:           c.ID,
		Status:       c.Status,
		Server:       c.Server,
		ArriveTime:   time.Unix(0, c.ArriveTime),
		QTime:        time.Unix(0, c.QTime),
		NextAttemp:   time.Unix(0, c.NextAttemp),
		LastAccess:   time.Unix(0, c.LastAccess),
		ReleaseTime:  time.Unix(0, c.ReleaseTime),
		RefreshCount: c.RefreshCount,
		MAC:          c.MAC,
	}
}

func ginContext2NewClient(c *gin.Context) clientData {
	host, _ := getHost(c.Request.Host)
	client := newClientData(host)
	return client
}

func ginContext2Client(c *gin.Context) (clientData, error) {
	r := c.Request
	host, err := getHost(r.Host)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "unknow host." + r.Host,
		})
		log.Errorln("unknow host:", r.Host)
		return clientData{}, errors.New("unknow host")
	}

	newClient := false
	client := clientData{}

	reqCookie, err := r.Cookie(cookieName)
	if err != nil {
		// log.Debugln("get cookie error: ", err)
		newClient = true
	} else {
		clientCookie := clientDataCookie{}
		j, err := base64.StdEncoding.DecodeString(reqCookie.Value)
		if err != nil {
			log.Debugln("base64 decode error: ", err)
			newClient = true
		}
		err = json.Unmarshal(j, &clientCookie)
		if err != nil {
			log.Debugln("json unmarshall error: ", err)
			newClient = true
		}

		client = clientCookie.toClient()
		if appRunMode == "debug" { //no need to verify every request in production
			if !client.isValid() {
				log.Errorln("invalid cookie mac remote ip: ", c.Request.RemoteAddr)
				log.Errorln(client)
				log.Errorln(clientCookie)
				return clientData{}, errors.New("invalid cookie mac")
			}
		}
	}

	client.Server = host

	// verify cookie if no cookie generate new and redirect to waiting room
	if newClient {
		client = newClientData(host)
	}

	// log.Debugln(client)
	return client, nil
}

func (client *clientData) saveCookie(c *gin.Context) {
	//generate new message authen
	client.LastAccess = time.Now()
	client.genMAC()

	clientCookie := client.toCookie()

	j, err := json.Marshal(clientCookie)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error: Generate client data.",
		})
		log.Fatal("Error: Generate client data")
		return
	}
	str := base64.StdEncoding.EncodeToString([]byte(j))

	maxage := 3600 //3600 one hour
	if appRunMode == "debug" {
		maxage = 600
	}
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    str,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxage,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
}
