// cmd/server/main.go

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/vinaykhade/resume-builder/backend/api/protobuf/job" // replace with your proto package name
)

func main() {
	// Create a gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register the gRPC server with the protobuf job service
	pb.RegisterJobServiceServer(s, &jobService{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	// Start the server
	log.Printf("starting server on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Implement the JobService interface defined in the protobuf file
type jobService struct {
	pb.UnimplementedJobServiceServer
}
