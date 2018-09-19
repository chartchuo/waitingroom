package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"apichart.me/waitingroom/advisor/adv"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

func (s *server) Update(c context.Context, r *adv.RequestStat) (*adv.AdvData, error) {
<<<<<<< HEAD
	log.Println(r.Sum, r.Count)
=======
	log.Println(r)
	log.Println(r.Count)
>>>>>>> dc9ce660c91b819f172fdb5a8e090a2fb8ce61a4
	return &adv.AdvData{}, nil
}

func main() {
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
