package main

import (
	"time"

	log "github.com/sirupsen/logrus"

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

// Reponse time unit in milisec
var inRespTime = make(chan clientChan, 100)
var avgRespTime int

func startAdvisor() {
	go localAdvisor()
}

const advInterval = 1      //second
const advResetInterval = 5 //second

// single go routine no need to aware race condition
func localAdvisor() {
	conn, err := grpc.Dial(confManager.Get().Advisor, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR connect to advisor: %v", err)
	}
	defer conn.Close()

	// client := adv.NewAdvServiceClient(conn)
	// ctx := context.Background()

	tick := time.Tick(advInterval * time.Second) //update with advisor time interval
	resettick := time.Tick(advResetInterval * time.Second)
	for {
		select {
		case c := <-inRespTime: //receive information from worker
			session.add(c.Server, c.ID)
			ss, err := getServerData(c.Server)
			if err != nil {
				log.Error("Local advisor can't fund server/host" + c.Server)
			}
			s := ss.counter
			//calculate 95 percentile
			p95 := calP95(s.p95, s.count, c.responseTime)

			s.count++
			s.sum += c.responseTime
			s.p95 = p95
			if session.concurrent(c.Server) >= ss.MaxUsers {
				ss.counter.concurrentusers = session.concurrent(c.Server)
				ss.Status = serverStatusWaitRoom
				ss.ReleaseTime = time.Now()
				setServerData(c.Server, ss)
			}

		case <-resettick:
			//update statistic and communicate to global advisor
			for host := range serverdataDB {
				s, err := getServerData(host)
				if err != nil {
					log.Error("Local advisor can't fund server/host" + host)
				}
				counter := s.counter

				//update data
				count := counter.count
				sum := counter.sum
				counter.maxresponsetime = getP95(counter.p95)
				counter.concurrentusers = session.concurrent(host)

				//update prometheus metrics
				requestRateMetric.WithLabelValues(host).Set(float64(count) / float64(advResetInterval))
				avgResponseTimeMetric.WithLabelValues(host).Set(float64(sum) / float64(count))
				concurrentUserMetric.WithLabelValues(host).Set(float64(counter.concurrentusers))
				maxResponseTimeMetric.WithLabelValues(host).Set(float64(getP95Max(counter.p95)))
				p95ResponseTimeMetric.WithLabelValues(host).Set(float64(counter.maxresponsetime))
				avgSessionTimeMetric.WithLabelValues(host).Set(float64(session.avgTime(host)))

				//communicate to advisor
				// stat := &adv.RequestStat{Sum: int32(sum), Count: int32(count)}
				// advData, err := client.Update(ctx, stat)
				// if err != nil {
				// 	log.Errorf("Can't connect to advise: %v %v", err, advData)
				// 	advisorState = advisorStatusLocal
				// 	break
				// }
			}

			for host := range serverdataDB {
				s, err := getServerData(host)
				if err != nil {
					log.Error("Local advisor can't fund server/host" + host)
				}
				counter := s.counter

				//update data
				counter.maxresponsetime = getP95(counter.p95)
				counter.concurrentusers = session.concurrent(host)

				//reset counter
				counter.count = 0
				counter.sum = 0
				counter.p95 = make([]int, 0, p95cap)
			}

		case <-tick: //calculate statistic info
			switch advisorState {
			case advisorStatusLocal:
				for host := range serverdataDB { //local calculation
					s, err := getServerData(host)
					if err != nil {
						log.Error("Local advisor can't fund server/host" + host)
					}

					switch s.Status {
					case serverStatusNormal:
						if s.counter.concurrentusers > s.MaxUsers {
							s.Status = serverStatusWaitRoom
							s.ReleaseTime = time.Now()
							setServerData(host, s)
						}
					case serverStatusNotOpen:
						n := time.Now()
						if s.OpenTime.Before(n) {
							s.Status = serverStatusWaitRoom
							s.ReleaseTime = s.OpenTime
							setServerData(host, s)
							log.Debugf("Open server: %v", host)
						}
					case serverStatusWaitRoom:
						cu := s.counter.concurrentusers
						ff := 4.0 //max fast forward 4x
						if cu != 0 {
							ff = float64(s.MaxUsers) / float64(cu)
							if ff > 4 {
								ff = 4
							}
						}
						if cu < s.MaxUsers {
							s.ReleaseTime = s.ReleaseTime.Add(advInterval * time.Second * time.Duration(ff))
						}
						if s.ReleaseTime.After(time.Now()) {
							s.ReleaseTime = time.Now()
							if cu < s.MaxUsers {
								s.Status = serverStatusNormal
							}
						}
						// log.Debugf("release: %v", time.Now().Sub(s.ReleaseTime))
						setServerData(host, s)
					}
				}
			}
		}
	}
}
