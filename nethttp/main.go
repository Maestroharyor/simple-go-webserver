package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maestroharyor/go-webserver/nethttp/internal/database"
	"github.com/maestroharyor/go-webserver/nethttp/internal/handlers"
	"github.com/maestroharyor/go-webserver/nethttp/internal/middleware"
	"github.com/maestroharyor/go-webserver/nethttp/internal/repository"
)

func main() {

	db, err := database.Connect("notes")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	defer db.Close()

	serverMux := http.NewServeMux()

	noteHandler := handlers.NewNoteHandler(repository.NewNoteRepository(db))

	serverMux.HandleFunc("GET /{$}", noteHandler.RootHandler)
	serverMux.HandleFunc("POST /notes", noteHandler.CreateNoteHandler)
	serverMux.HandleFunc("GET /notes", noteHandler.ListNotesHandler)
	serverMux.HandleFunc("GET /notes/{id}", noteHandler.GetSingleNoteHandler)
	serverMux.HandleFunc("PUT /notes/{id}", noteHandler.UpdateSingleNoteHandler)
	serverMux.HandleFunc("DELETE /notes/{id}", noteHandler.DeleteSingleNoteHandler)

	addr := ":8080"
	fmt.Println("Server listening on", addr)
	if err := http.ListenAndServe(addr, middleware.Logging(serverMux)); err != nil {
		log.Fatal(err)
	}
}
