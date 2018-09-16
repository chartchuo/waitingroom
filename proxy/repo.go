package main

import (
	"errors"
	"strings"
)

type Host struct {
	Name string
}

type Target struct {
	TargetAddress string
}

type InspectField struct {
	Field string
}

//map url to host
var HostDB = map[string]string{
	"chartchuo.noip.me:80": "chartchuo.com",
	"www.pantip.com:80":    "pantip.com",
	"pantip.com:80":        "pantip.com",
	"m.pantip.com:80":      "pantip.com",
	"mockserver:80":        "mock",
}

//map host to target
var targetDB = map[string]Target{
	"chartchuo.com": {TargetAddress: "58.11.248.160:80"},
	"pantip.com":    {TargetAddress: "203.151.13.167:80"},
	"mock":          {TargetAddress: "localhost:8888"},
}

//map host+path to sqli
var sqlInjectionDB = map[string]map[string][]InspectField{
	"pantip.com": {
		"/search/es/search_tag": {
			{Field: "q"},
		},
	},
	"chartchuo.com": {
		"/users/sign_in": {
			{Field: "user[login]"},
		},
	},
}

//map host+path to sqli
var xssDB = map[string]map[string][]InspectField{
	"pantip.com": {
		"/search/es/search_tag": {
			{Field: "q"},
		},
	},
}

func HostGet(domainname string) (string, error) {
	if !strings.Contains(domainname, ":") {
		domainname = domainname + ":80"
	}
	s, ok := HostDB[domainname]
	if !ok {
		return "", errors.New("Invalid domain")
	}
	return s, nil
}

func TargetAddress(d string) (string, error) {
	h, ok := HostDB[d]
	if !ok {
		return "", errors.New("Invalid host")
	}
	t, ok := targetDB[h]
	if !ok {
		return "", errors.New("Invalid domain")
	}
	return t.TargetAddress, nil
}

func SQLInjectionFieldGet(domain string, path string) ([]InspectField, error) {
	sqli, ok := sqlInjectionDB[domain][path]
	if !ok {
		return nil, errors.New("Invalid domain")
	}
	return sqli, nil
}

func XSSFieldGet(domain string, path string) ([]InspectField, error) {
	xss, ok := xssDB[domain][path]
	if !ok {
		return nil, errors.New("Invalid domain")
	}
	return xss, nil
}
