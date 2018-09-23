package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const qSpanTime = time.Minute //one minute
const cookieName = "ccwait"
const uuidNameSpace = "ccwait"
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
		Status:       clientStatusNew,
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

func newClientFromC(c *gin.Context) (clientData, error) {
	r := c.Request
	host, err := hostGet(r.Host)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "unknow host.",
		})
		return clientData{}, errors.New("Error: unknow host")
	}

	newClient := false
	client := clientData{}

	reqCookie, err := r.Cookie(cookieName)
	if err != nil {
		newClient = true
	} else {
		client := clientData{}
		j, err := base64.StdEncoding.DecodeString(reqCookie.Value)
		if err != nil {
			newClient = true
		}
		err = json.Unmarshal(j, &client)
		if err != nil {
			newClient = true
		}
	}

	// verify cookie if no cookie generate new and redirect to waiting room
	if newClient {
		client = newClientData(host)
	}
	//ensure host consistant with URL
	if client.Server != host {
		log.Debugf("Server mismatch client.Server: %v host:%v\n", client.Server, host)
	}
	client.Server = host
	return client, nil
}

func setClientCookie(c *gin.Context, client clientData) {

	j, err := json.Marshal(client)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error: Generate client data.",
		})
		log.Fatal("Error: Generate client data")
		return
	}
	str := base64.StdEncoding.EncodeToString([]byte(j))
	// webdata.ClientID = newClientID()
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    str,
		SameSite: http.SameSiteStrictMode,
		// MaxAge:   3600, //one hour
		MaxAge: 600, //todo for testing only
	}
	http.SetCookie(c.Writer, &cookie)
}
