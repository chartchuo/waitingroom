package main

import (
	"testing"
)

func TestSession(t *testing.T) {
	session.add("test1", "1")
	session.add("test1", "2")
	session.add("test1", "3")
	session.add("test2", "a")
	session.add("test3", "b")

	time, ok := session["test1"]["1"]
	if !ok {
		t.Errorf("1. error session not found t: %v ok: %v", time, ok)
	}
	time, ok = session["test1"]["2"]
	if !ok {
		t.Errorf("2. error session not found t: %v ok: %v", time, ok)
	}
	time, ok = session["test1"]["3"]
	if !ok {
		t.Errorf("3. error session not found t: %v ok: %v", time, ok)
	}
	time, ok = session["test2"]["a"]
	if !ok {
		t.Errorf("4. error session not found t: %v ok: %v", time, ok)
	}
	time, ok = session["test3"]["b"]
	if !ok {
		t.Errorf("5. error session not found t: %v ok: %v", time, ok)
	}
}
