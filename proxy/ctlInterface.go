package main

import (
	"context"
	"log"
	"time"

	"apichart.me/waitingroom/controller/ctl"
	"google.golang.org/grpc"
)

// milisec
var inRespTime = make(chan int, 10)
var avgRespTime int

func respTimePoller() {
	conn, err := grpc.Dial("controller:6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to controller: %v", err)
	}
	defer conn.Close()

	client := ctl.NewCtlClient(conn)
	stream, err := client.CtlStream(context.Background())
	if err != nil {
		log.Fatalf("ERROR connect to controller: %v", err)
	}
	defer stream.CloseSend()

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
			stat := &ctl.RequestStat{Sum: int32(sum), Count: int32(count)}
			err := stream.Send(stat)
			if err != nil {
				log.Printf("Send to Ctl error: %v\n", err)
			}
			sum = 0
			count = 0
			avgRespTime = 0
		}
	}
}
