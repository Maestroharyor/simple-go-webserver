package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maestroharyor/go-webserver/nethttp/internal/database"
	"github.com/maestroharyor/go-webserver/nethttp/internal/env"
	"github.com/maestroharyor/go-webserver/nethttp/internal/handlers"
	"github.com/maestroharyor/go-webserver/nethttp/internal/middleware"
	"github.com/maestroharyor/go-webserver/nethttp/internal/repository"
)

func main() {

	port := fmt.Sprintf(":%s", env.GetEnv("PORT", "8080"))
	databasePath := env.GetEnv("DATABASE_PATH", "notes.db")
	db, err := database.Connect(databasePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	defer db.Close()

	serverMux := http.NewServeMux()

	noteHandler := handlers.NewNoteHandler(repository.NewNoteRepository(db))

	serverMux.HandleFunc("GET /{$}", noteHandler.RootHandler)
	serverMux.HandleFunc("GET /health", noteHandler.HealthCheckHandler)
	serverMux.HandleFunc("POST /notes", noteHandler.CreateNoteHandler)
	serverMux.HandleFunc("GET /notes", noteHandler.ListNotesHandler)
	serverMux.HandleFunc("GET /notes/{id}", noteHandler.GetSingleNoteHandler)
	serverMux.HandleFunc("PUT /notes/{id}", noteHandler.UpdateSingleNoteHandler)
	serverMux.HandleFunc("DELETE /notes/{id}", noteHandler.DeleteSingleNoteHandler)

	fmt.Println("Server listening on", port)

	if err := http.ListenAndServe(port, middleware.Logging(serverMux)); err != nil {
		log.Fatal(err)
	}
}
