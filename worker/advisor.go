package main

import (
	"context"
	"log"
	"time"

	"apichart.me/waitingroom/advisor/adv"
	"google.golang.org/grpc"
)

// Reponse time unit in milisec
var inRespTime = make(chan int, 10)
var avgRespTime int

func respTimePoller() {
	conn, err := grpc.Dial("advisor:6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to advisor: %v", err)
	}
	defer conn.Close()

	client := adv.NewAdvServiceClient(conn)
	ctx := context.Background()

	var sum, count int
	tick := time.Tick(time.Second * 5)
	for {
		select {
		case r := <-inRespTime:
			sum += r
			count++
			avgRespTime = sum / count
		case <-tick:
			//todo non blocking by adding go routine and channel
			// log.Printf("sum: %v, count, %v", sum, count)

			// log.Println("start call grps")
			stat := &adv.RequestStat{Sum: int32(sum), Count: int32(count)}
			_, err := client.Update(ctx, stat)

			if err != nil {
				log.Printf("could not get advise: %v", err)
				continue
			}
			// log.Printf("Advise: %v", advd.ReleaseTime) //todo remove.

			sum = 0
			count = 0
			avgRespTime = 0
		}
	}
}
