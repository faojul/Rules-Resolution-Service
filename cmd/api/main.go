package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rules-resolution-service/internal/handler"
	"rules-resolution-service/internal/repository"
	"rules-resolution-service/internal/service"
	"rules-resolution-service/pkg/db"
)

func main() {
    conn := db.NewPostgres()

    repo := repository.NewPostgresOverrideRepository(conn)
    resolver := service.NewResolver(repo)
    overrideService := service.NewOverrideService(repo)

    resolveHandler := handler.NewResolveHandler(resolver)
    overrideHandler := handler.NewOverrideHandler(overrideService)
    explainHandler := handler.NewExplainHandler(resolver)
    conflictHandler := handler.NewConflictHandler(overrideService)

    r := chi.NewRouter()
    r.Post("/api/resolve", resolveHandler.Resolve)
    r.Post("/api/bulk-resolve", resolveHandler.BulkResolve)
    r.Post("/api/resolve/explain", explainHandler.Explain)
    r.Get("/api/overrides/{id}", overrideHandler.GetByID)
    r.Post("/api/overrides", overrideHandler.Create)
    r.Get("/api/overrides", overrideHandler.List)
    r.Put("/api/overrides/{id}", overrideHandler.Update)
    r.Patch("/api/overrides/{id}/status", overrideHandler.UpdateStatus)
    r.Get("/api/overrides/conflicts", conflictHandler.GetConflicts)

    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", r)
}