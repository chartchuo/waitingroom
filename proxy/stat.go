package main

import (
	"time"
)

// milisec
var inRespTime = make(chan int, 10)
var avgRespTime int

func respTimePoller() {
	var sum, count int
	tick := time.Tick(time.Second * 10)
	for {
		select {
		case r := <-inRespTime:
			sum += r
			count++
			avgRespTime = sum / count
		case <-tick:
			// log.Printf("sum: %v, count, %v", sum, count)
			sum = 0
			count = 0
			avgRespTime = 0
		}
	}
}
