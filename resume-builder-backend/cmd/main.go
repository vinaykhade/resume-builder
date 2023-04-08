// cmd/main.go

// That's it for step 4. The next step is to define your gRPC services and register them with the server.

package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	"google.golang.org/grpc"
	"log"

	"github.com/vinaykhade/resume-builder-backend/internal/models"
	"github.com/vinaykhade/resume-builder-backend/internal/services"
)

func main() {

	// Connect to database
	db, err := services.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate database schema
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register your gRPC services here

	// Start the server
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down gRPC server...")
	server.GracefulStop()
	log.Println("gRPC server stopped.")
}
