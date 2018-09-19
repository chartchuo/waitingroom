package main

import "net/http"

type InspecResult int

const INSPECT_OK = 0
const INSPECT_ERROR = 3
const INSPECT_SUSPICIOUS = 7
const INSPECT_ATTACK = 9

type WebInspectData struct {
	Host     string
	ClientIP string
	ClientID string //"" == new client
	R        *http.Request
	// RequestBuffered bool
	BodyBuf []byte
}
