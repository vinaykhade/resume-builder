// api/server.go

package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vinaykhade/resume-builder-backend/api/graph/generated"
	"github.com/vinaykhade/resume-builder-backend/internal/services"
)

func StartGraphQLServer() error {
	resumeService := services.NewResumeService(NewDatabase())

	resolver := NewResolver(resumeService)

	cfg := generated.Config{
		Resolvers: resolver,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	r := mux.NewRouter()

	r.HandleFunc("/", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", srv)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(r)

	log.Println("Starting GraphQL server on port 8080")

	return http.ListenAndServe(":8080", handler)
}
