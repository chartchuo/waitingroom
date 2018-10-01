package main

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := newClientData("test")
	client.genMAC()
	valid := client.isValid()
	if !valid {
		t.Errorf("1. validation function expect: %v got %v", true, valid)
	}

	// test unauthorize modify
	client.QTime = client.QTime.Add(time.Second)
	valid = client.isValid()
	if valid {
		t.Errorf("2. validation function expect: %v got %v", false, valid)
	}

	client.genMAC()
	valid = client.isValid()
	if !valid {
		t.Errorf("3. validation function expect: %v got %v", true, valid)
	}

}

func TestClientCookie(t *testing.T) {
	client := newClientData("test")
	client.genMAC()
	cookie := client.toCookie()
	newClient := cookie.toClient()

	if client != newClient {
		t.Errorf("convert clinet to cookie and convert back not equal client: %+v new  %+v", client, newClient)
	}
}
