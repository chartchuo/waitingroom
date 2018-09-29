package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"apichart.me/waitingroom/advisor/adv"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type clientChan struct {
	clientData
	responseTime int
}

// Reponse time unit in milisec
var inRespTime = make(chan clientChan, 10)
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

	cons := make(map[string]*cache.Cache)
	// var sum, count int
	tick := time.Tick(advInterval * time.Second) //update with advisor time interval
	for {
		select {
		case c := <-inRespTime:
			con, ok := cons[c.Server]
			if !ok {
				cons[c.Server] = cache.New(time.Minute, time.Second*advInterval)
				con, _ = cons[c.Server]
				con.Set("count", float64(0), cache.NoExpiration)
				con.Set("sum", float64(0), cache.NoExpiration)
			}
			con.Set(c.ID, 1, cache.DefaultExpiration)
			con.Increment("count", 1)
			con.Increment("sum", int64(c.responseTime))
		case <-tick:
			for host, con := range cons {
				count, ok := con.Get("count")
				if !ok {
					continue
				}
				sum, ok := con.Get("sum")
				if !ok {
					continue
				}
				stat := &adv.RequestStat{Sum: int32(sum.(float64)), Count: int32(count.(float64))}
				advData, err := client.Update(ctx, stat)

				if err != nil {
					log.Errorf("Can't connect to advise: %v", err)
					// continue
				}
				requestRateMetric.WithLabelValues(host).Set(float64(count.(float64)) / float64(advInterval))
				avgResponseTimeMetric.WithLabelValues(host).Set(sum.(float64) / count.(float64))
				concurrentUserMetric.WithLabelValues(host).Set(float64(con.ItemCount() - 2)) //exclude count and sum metric

				//reset counter to zero
				con.Set("count", float64(0), cache.NoExpiration)
				con.Set("sum", float64(0), cache.NoExpiration)

				mockserver := serverdataDB["mock"]
				mockserver.ReleaseTime = time.Unix(0, advData.ReleaseTime)
				serverdataDB["mock"] = mockserver
				log.Debugln(serverdataDB["mock"].ReleaseTime)
			}
		}
	}
}
