package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

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
	w.Write([]byte("Hello, World!"))
}

func (h *NoteHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if note.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if note.Content == "" {
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
	if errors.Is(err, sql.ErrNoRows) {
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
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if note.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if note.Content == "" {
		WriteError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if err := h.repo.UpdateNoteById(r.Context(), id, &note); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	WriteSuccess(w, http.StatusOK, "Note updated", note)
}

func (h *NoteHandler) DeleteSingleNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if err := h.repo.DeleteNoteById(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	WriteSuccess[*struct{}](w, http.StatusOK, "Note deleted", nil)

}
