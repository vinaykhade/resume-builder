// main.go

package main

import (
	"github.com/vinaykhade/resume-builder-backend/api"
	"github.com/vinaykhade/resume-builder-backend/internal/services"
)

func main() {
	// Start gRPC server
	go services.StartGRPCServer()

	// Start gRPC-Gateway server
	go services.StartGatewayServer()

	// Start GraphQL server
	api.StartGraphQLServer()
}
