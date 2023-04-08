// internal/services/gateway.go

package services

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vinaykhade/resume-builder-backend/api/resume"
	"google.golang.org/grpc"
)

func NewGateway(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
	err := resume.RegisterResumeServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
		return err
	}

	return nil
}

func StartGatewayServer() error {
	ctx := context.Background()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := NewGateway(ctx, mux, opts)
	if err != nil {
		return err
	}

	log.Println("Starting gRPC gateway server on port 8080")

	return http.ListenAndServe(":8080", mux)
}
