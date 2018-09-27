package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"apichart.me/waitingroom/advisor/adv"
	"google.golang.org/grpc"
)

// Reponse time unit in milisec
var inRespTime = make(chan int, 10)
var avgRespTime int

func advisorPoller() {
	conn, err := grpc.Dial(confManager.Get().Advisor, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to advisor: %v", err)
	}
	defer conn.Close()

	client := adv.NewAdvServiceClient(conn)
	ctx := context.Background()

	var sum, count int
	tick := time.Tick(time.Second * 5) //update with advisor time interval
	for {
		select {
		case r := <-inRespTime:
			sum += r
			count++
			avgRespTime = sum / count
		case <-tick:
			//todo non blocking by adding go routine and channel

			stat := &adv.RequestStat{Sum: int32(sum), Count: int32(count)}
			advData, err := client.Update(ctx, stat)

			if err != nil {
				log.Errorf("Can't connect to advise: %v", err)
				continue
			}

			sum = 0
			count = 0
			avgRespTime = 0

			mockserver := serverdataDB["mock"]
			mockserver.ReleaseTime = time.Unix(0, advData.ReleaseTime)
			serverdataDB["mock"] = mockserver
			// log.Debugln(serverdataDB["mock"].ReleaseTime)

		}
	}
}
