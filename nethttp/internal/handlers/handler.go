package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/maestroharyor/go-webserver/nethttp/internal/models"
	"github.com/maestroharyor/go-webserver/nethttp/internal/repository"
)

type NoteHandler struct {
	repo *repository.NoteRepository
}

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{repo: repo}
}

func (h *NoteHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	WriteSuccess[*struct{}](w, http.StatusOK, "Hello world", nil)
}

func (h *NoteHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	WriteSuccess[*struct{}](w, http.StatusOK, "OK", nil)
}

func (h *NoteHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := decodeJSON(w, r, &note); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if strings.TrimSpace(note.Title) == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if strings.TrimSpace(note.Content) == "" {
		WriteError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if err := h.repo.CreateNote(r.Context(), &note); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}
	WriteSuccess(w, http.StatusCreated, "Note created successfully", note)
}

func (h *NoteHandler) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := h.repo.GetNotes(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch notes")
		return
	}
	WriteSuccess(w, http.StatusOK, "Notes fetched successfully", notes)
}

func (h *NoteHandler) GetSingleNoteHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := parseID(w, r)
	if !ok {
		WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	note, err := h.repo.GetNoteById(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "Note not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch note")
		return
	}
	WriteSuccess(w, http.StatusOK, "Note fetched successfully", note)
}

func (h *NoteHandler) UpdateSingleNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var note models.Note

	fmt.Printf("Title: %s\n", note.Title)
	fmt.Printf("Content: %s\n", note.Content)

	if strings.TrimSpace(note.Title) == "" && strings.TrimSpace(note.Content) == "" {
		WriteError(w, http.StatusBadRequest, "Note title or content are required")
		return
	}

	if err := decodeJSON(w, r, &note); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	updated, err := h.repo.UpdateNoteById(r.Context(), id, &note)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "Note not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	WriteSuccess(w, http.StatusOK, "Note updated", updated)
}

func (h *NoteHandler) DeleteSingleNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	err := h.repo.DeleteNoteById(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "Note not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	WriteSuccess[*struct{}](w, http.StatusNoContent, "Note deleted", nil)

}
