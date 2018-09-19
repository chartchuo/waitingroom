package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"apichart.me/waitingroom/advisor/adv"
)

type server struct {
}

func (s *server) Update(c context.Context, r *adv.RequestStat) (*adv.AdvData, error) {
	log.Println(r.Sum, r.Count)
	return &adv.AdvData{}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	// ctl.RegisterCtlServer(grpcServer, &server{})
	adv.RegisterAdvServiceServer(grpcServer, &server{})

	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening on tcp://localhost:6000")
	grpcServer.Serve(l)
}
