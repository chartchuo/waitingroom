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

type serverCounter struct {
	count          int
	sum            int
	concurrentUser *cache.Cache
	p95            []int
}

// Reponse time unit in milisec
var inRespTime = make(chan clientChan, 10)
var avgRespTime int

func startAdvisor() {
	go localAdvisor()
}

const advInterval = 5 //second

func calP95(p95 []int, count int, rt int) []int {
	n := len(p95)
	if n == 0 {
		p95 = append(p95, rt)
		return p95
	}

	return p95
}

func localAdvisor() {
	conn, err := grpc.Dial(confManager.Get().Advisor, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to advisor: %v", err)
	}
	defer conn.Close()

	client := adv.NewAdvServiceClient(conn)
	ctx := context.Background()

	// cons := make(map[string]*cache.Cache)
	serverStat := make(map[string]serverCounter)
	tick := time.Tick(advInterval * time.Second) //update with advisor time interval
	for {
		select {
		case c := <-inRespTime: //receive information from worker
			_, ok := serverStat[c.Server]
			if !ok {
				serverStat[c.Server] = serverCounter{
					count:          0,
					sum:            0,
					concurrentUser: cache.New(time.Minute, time.Second*advInterval),
					p95:            make([]int, 0, 1000),
				}
			}
			serverStat[c.Server].concurrentUser.Set(c.ID, 1, cache.DefaultExpiration) //add client id cache

			//calculate 95 percentile
			p95 := calP95(serverStat[c.Server].p95, serverStat[c.Server].count, c.responseTime)

			serverStat[c.Server] = serverCounter{
				count:          serverStat[c.Server].count + 1,
				sum:            serverStat[c.Server].sum + c.responseTime,
				concurrentUser: serverStat[c.Server].concurrentUser,
				p95:            p95,
			}

		case <-tick: //calculate statistic info
			for host, counter := range serverStat {
				count := serverStat[host].count
				sum := serverStat[host].sum
				stat := &adv.RequestStat{Sum: int32(sum), Count: int32(count)}
				advData, err := client.Update(ctx, stat)
				if err != nil {
					log.Errorf("Can't connect to advise: %v %v", err, advData)
					// continue
				}
				requestRateMetric.WithLabelValues(host).Set(float64(count) / float64(advInterval))
				avgResponseTimeMetric.WithLabelValues(host).Set(float64(sum) / float64(count))
				concurrentUserMetric.WithLabelValues(host).Set(float64(counter.concurrentUser.ItemCount()))

				serverStat[host] = serverCounter{
					count:          0,
					sum:            0,
					concurrentUser: counter.concurrentUser,
				}
			}

			//mock server data tobe remove
			// mockserver := serverdataDB["mock"]
			// mockserver.ReleaseTime = time.Unix(0, advData.ReleaseTime)
			// serverdataDB["mock"] = mockserver
			// log.Debugln(serverdataDB["mock"].ReleaseTime)
		}
	}
}
