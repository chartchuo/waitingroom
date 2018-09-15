package main

// todo
// inject cookie
// inspect browser/cookie level
// operation mode 1.bypass 2.normal block 3.hard block
// global level and host level bucket
// detect base line request rate

import (
	"time"

	"github.com/pmylund/go-cache"
)

var LeakyBuckets = cache.New(BucketTimeout, IntervalTime)

func inspectHttpFlood(clientIP string) InspecResult {
	// add to bucket
	count, err := LeakyBuckets.IncrementInt(clientIP, 1)
	if err != nil {
		LeakyBuckets.Set(clientIP, int(1), cache.DefaultExpiration)
	} else {
		// Bucket full
		if count >= BucketSize {
			LeakyBuckets.Delete(clientIP)
			return INSPECT_ATTACK
		}
	}

	return INSPECT_OK
}

func Inspect(d *WebInspectData) InspecResult {

	result := inspectHttpFlood(d.ClientIP)
	if result != INSPECT_OK {
		return result
	}

	return INSPECT_OK
}

func leak() {
	for {
		time.Sleep(IntervalTime)
		buckets := LeakyBuckets.Items()
		//LeakyBuckets.Lock()
		for _, b := range buckets {
			rv, _ := b.Object.(int)
			nv := rv - LeakRate
			if nv < 0 {
				nv = 0
			}
			b.Object = nv
		}
		//LeakyBuckets.Unlock()
	}
}

func init() {
	go leak()
}
