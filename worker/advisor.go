package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"apichart.me/waitingroom/advisor/adv"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

const (
	advisorStatusLocal serverStatus = iota
	advisorStatusGlobal
)

var advisorState = advisorStatusLocal

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

func (c *serverCounter) getConcurrentUser() int {
	return c.concurrentUser.ItemCount()
}

// Reponse time unit in milisec
var inRespTime = make(chan clientChan, 100)
var avgRespTime int

func startAdvisor() {
	go localAdvisor()
}

const advInterval = 5 //second

// single go routine no need to aware race condition
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
					p95:            make([]int, 0, p95cap),
				}
			}
			serverStat[c.Server].concurrentUser.Set(c.ID, 1, cache.DefaultExpiration) //add client id cache

			//calculate 95 percentile
			p95 := calP95(serverStat[c.Server].p95, serverStat[c.Server].count, c.responseTime)
			// p95 := calP95(serverStat[c.Server].p95, serverStat[c.Server].count, serverStat[c.Server].count%100) //test p95 by distribyte 0-99

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
					advisorState = advisorStatusLocal
					break
				}
				requestRateMetric.WithLabelValues(host).Set(float64(count) / float64(advInterval))
				avgResponseTimeMetric.WithLabelValues(host).Set(float64(sum) / float64(count))
				concurrentUserMetric.WithLabelValues(host).Set(float64(counter.getConcurrentUser()))
				maxResponseTimeMetric.WithLabelValues(host).Set(float64(getP95Max(counter.p95)))
				p95ResponseTimeMetric.WithLabelValues(host).Set(float64(getP95(counter.p95)))
			}

			switch advisorState {
			case advisorStatusLocal:
				for host, counter := range serverStat { //local calculation
					s, err := getServerData(host)
					if err != nil {
						log.Error("Local advisor can't fund server/host" + host)
					}

					switch s.Status {
					case serverStatusNormal:
					case serverStatusNotOpen:
					case serverStatusWaitRoom:
						cu := counter.getConcurrentUser()
						if cu < s.MaxUsers/2 {
							s.ReleaseTime = s.ReleaseTime.Add(advInterval * time.Second * 2)
						} else if cu < s.MaxUsers {
							s.ReleaseTime = s.ReleaseTime.Add(advInterval * time.Second)
						}
						if s.ReleaseTime.After(time.Now()) {
							s.ReleaseTime = time.Now()
						}
						setServerData(host, s)
						log.Debugf("ADV: release time for host %v: %v\n dif from now:%v", host, s.ReleaseTime, time.Now().Sub(s.ReleaseTime))
					}

				}

			}

			for host, counter := range serverStat { //clear data for all server
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
