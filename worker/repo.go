package main

import (
	"errors"
	"net/http"
	"strings"
)

//Config main repository
type Config struct {
	HostDB       map[string]string       `yaml:"hostdb,omitempty"`
	TargetDB     map[string]string       `yaml:"targetdb,omitempty"`
	ServerConfig map[string]ServerConfig `yaml:"serverconfig,omitempty"`
}

//InspecResult inspec result from ddos and waf
type InspecResult int

const (
	inspecOk         InspecResult = 0
	inspecError      InspecResult = 3
	inspecSuspicious InspecResult = 7
	inspedAttack     InspecResult = 7
)

//WebInspectData inspec data for ddos and waf
type WebInspectData struct {
	Host     string
	ClientIP string
	ClientID string //"" == new client
	R        *http.Request
	// RequestBuffered bool
	BodyBuf []byte
}

func hostGet(domainname string) (string, error) {
	if !strings.Contains(domainname, ":") {
		domainname = domainname + ":80"
	}
	c := confManager.Get()
	s, ok := c.HostDB[domainname]
	if !ok {
		return "", errors.New("Invalid domain")
	}
	return s, nil
}

func targetAddress(d string) (string, error) {
	c := confManager.Get()
	h, ok := c.HostDB[d]
	if !ok {
		return "", errors.New("Invalid host")
	}
	t, ok := c.TargetDB[h]
	if !ok {
		return "", errors.New("Invalid domain")
	}
	return t, nil
}
