package main

import (
	"context"

	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"apichart.me/waitingroom/advisor/adv"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

func (s *server) Update(c context.Context, r *adv.RequestStat) (*adv.AdvData, error) {
	log.Debugf("Sum: %v, Count: %v\n", r.Sum, r.Count)
	return &adv.AdvData{}, nil
}

const appRunMode = "debug"

func main() {
	if appRunMode == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	adv.RegisterAdvServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
