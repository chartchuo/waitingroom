package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func ccDial(network, address string) (net.Conn, error) {
	var d net.Dialer
	newAddress, err := targetAddress(address)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, newAddress)
}

var transport = &http.Transport{
	Dial: ccDial,
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// r := d.R
	url := "http://" + r.Host + r.URL.Path

	//todo log level
	log.Println(r.RemoteAddr + " " + r.Method + " " + url)
	var req *http.Request
	req, _ = http.NewRequest(r.Method, url, r.Body)

	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}

	startTime := time.Now()

	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	diff := time.Now().Sub(startTime)
	inRespTime <- int(diff / time.Microsecond)

	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	resp.Body.Close()

}

func proxyHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	host, err := hostGet(r.Host)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "unknow host.",
		})
		return
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
		j, err := json.Marshal(client)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error: Generate client data.",
			})
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
		http.SetCookie(w, &cookie)
		c.Redirect(http.StatusTemporaryRedirect, waitRoomPath)
	}

	//todo if no queue do proxy

	//todo if not open redirec to wait

	//todo if not qtime redirect to wait

	proxyRequest(w, r)
}
