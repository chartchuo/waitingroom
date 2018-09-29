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

func startAdvisor() {
	go localAdvisor()
}

const advInterval = 5 //second

func localAdvisor() {
	conn, err := grpc.Dial(confManager.Get().Advisor, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to advisor: %v", err)
	}
	defer conn.Close()

	client := adv.NewAdvServiceClient(conn)
	ctx := context.Background()

	var sum, count int
	tick := time.Tick(advInterval * time.Second) //update with advisor time interval
	for {
		select {
		case r := <-inRespTime:
			sum += r
			count++
			avgRespTime = sum / count
		case <-tick:
			stat := &adv.RequestStat{Sum: int32(sum), Count: int32(count)}
			advData, err := client.Update(ctx, stat)

			if err != nil {
				log.Errorf("Can't connect to advise: %v", err)
				// continue
			}
			avgResponseTimeMetric.WithLabelValues("mock").Set(float64(avgRespTime))
			requestRateMetric.WithLabelValues("mock").Set(float64(count) / float64(advInterval))

			log.Debugf("sum: %v, count: %v, avg: %v\n", sum, count, avgRespTime)
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
