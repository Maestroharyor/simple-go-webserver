package handlers

import (
	"net/http"

	"github.com/maestroharyor/go-webserver/nethttp/internal/repository"
)

type NoteHandler struct {
	repo *repository.NoteRepository
}

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{repo: repo}
}

func (h *NoteHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (h *NoteHandler) NotesRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func (h *NoteHandler) NotesSingleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}
