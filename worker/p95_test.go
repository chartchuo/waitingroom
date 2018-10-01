package main

import "testing"

func TestP95(t *testing.T) {
	data := make([]int, 0, p95cap)
	var v int

	t.Log("Test round 1 100 element 1-100")
	for i := 1; i <= 100; i++ {
		data = calP95(data, i, i)
	}
	t.Log(data)
	v = getP95Max(data)
	if v != 100 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 100)
	}
	v = getP95(data)
	if v != 95 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 95)
	}

	t.Log("Test round 2 more 100 element 1-100")
	for i := 1; i <= 100; i++ {
		data = calP95(data, i+100, i)
	}
	t.Log(data)
	v = getP95Max(data)
	if v != 100 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 100)
	}
	v = getP95(data)
	if v != 95 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 95)
	}

	t.Log("Test round 3 more 100 element 1-100")
	for i := 1; i <= 100; i++ {
		data = calP95(data, i+200, i)
	}
	t.Log(data)
	v = getP95Max(data)
	if v != 100 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 100)
	}
	v = getP95(data)
	if v != 95 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, 95)
	}

}

func TestP95limit(t *testing.T) {
	data := make([]int, 0, p95cap)
	var v int

	t.Log("Test round 1 100 element 1-1000")
	for i := 1; i <= p95cap; i++ {
		data = calP95(data, i, i)
	}
	v = getP95Max(data)
	if v != p95cap {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, p95cap)
	}
	v = getP95(data)
	if v != p95cap*95/100 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, p95cap*95/100)
	}
}
func TestP95overflow(t *testing.T) {
	data := make([]int, 0, p95cap)
	var v int

	t.Log("Test round 1 100 element 1-1000")
	for i := 1; i <= p95cap+1; i++ {
		data = calP95(data, i, i)
	}
	v = getP95Max(data)
	if v != p95cap+1 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, p95cap+1)
	}
	v = getP95(data)
	if v != p95cap*95/100+1 {
		t.Errorf("P95 max was incorrect, got: %d, want: %d.", v, p95cap*95/100+1)
	}
}
