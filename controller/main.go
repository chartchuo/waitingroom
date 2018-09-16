package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"apichart.me/waitingroom/controller/ctl"
)

type server struct {
}

func (s *server) CtlStream(stream ctl.Ctl_CtlStreamServer) error {
	log.Println("Started stream")
	for {
		in, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("sum: %v, count:%v\n", in.Sum, in.Count)
	}
}

func main() {
	grpcServer := grpc.NewServer()
	ctl.RegisterCtlServer(grpcServer, &server{})

	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening on tcp://localhost:6000")
	grpcServer.Serve(l)
}
