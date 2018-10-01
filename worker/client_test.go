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

func TestClientSpan(t *testing.T) {
	n := time.Now()
	for i := 0; i < 1000; i++ {
		v := spanTime(n)
		if !(v.After(n) && v.Before(n.Add(qSpanTime))) {
			t.Errorf("1 spantime not in spect range v: %v t: %v", v, n)
		}
	}

	n = time.Now().Add(time.Minute)
	for i := 0; i < 1000; i++ {
		v := spanTime(n)
		if !(v.After(n) && v.Before(n.Add(qSpanTime))) {
			t.Errorf("2 spantime not in spect range v: %v t: %v", v, n)
		}
	}

	n = time.Now().Add(qSpanTime / -2)
	for i := 0; i < 1000; i++ {
		v := spanTime(n)
		if !(v.After(n.Add(qSpanTime/2)) && v.Before(n.Add(qSpanTime))) {
			t.Errorf("3 spantime not in spect range v: %v t: %v", v, n)
		}
	}

	n = time.Now().Add(-qSpanTime)
	v := spanTime(n)
	if !v.Truncate(time.Millisecond).Equal(n.Add(qSpanTime).Truncate(time.Millisecond)) {
		t.Errorf("4 spantime not in spect range v: %v t: %v", v.Truncate(time.Millisecond), n.Truncate(time.Millisecond))
	}

	n = time.Now().Add(-qSpanTime * 2)
	v = spanTime(n)
	if !v.Truncate(time.Millisecond).Equal(n.Add(qSpanTime * 2).Truncate(time.Millisecond)) {
		t.Errorf("4 spantime not in spect range v: %v t: %v", v.Truncate(time.Millisecond), n.Truncate(time.Millisecond))
	}
}
