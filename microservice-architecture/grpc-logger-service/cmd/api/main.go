package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/grpc-logger-service/model"
	"google.golang.org/grpc"
)

const (
	grpcPort = "80"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	model.RegisterLogServiceServer(s, &LogServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
