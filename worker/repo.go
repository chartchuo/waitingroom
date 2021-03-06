package main

import (
	"errors"
	"net/http"
	"strings"
)

//Config main repository
type Config struct {
	Advisor      string                  `yaml:"advisor,omitempty"`
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

func getHost(domainname string) (string, error) {
	if !strings.Contains(domainname, ":") {
		domainname = domainname + ":80"
	}
	c := confManager.Get()
	s, ok := c.HostDB[domainname]
	if !ok {
		return "", errors.New("Domain not found: " + domainname)
	}
	return s, nil
}

func host2TargetAddress(host string) (string, error) {
	c := confManager.Get()
	t, ok := c.TargetDB[host]
	if !ok {
		return "", errors.New("Target not found host: " + host)
	}
	return t, nil
}
