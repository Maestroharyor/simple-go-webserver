package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maestroharyor/go-webserver/nethttp/internal/models"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	result, err := r.db.ExecContext(ctx, `INSERT INTO notes (title, content, created_at, updated_at) VALUES (? ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, note.Title, note.Content)

	if err != nil {
		return fmt.Errorf("create note error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get note id: %w", err)
	}

	note.ID = int(id)

	return nil
}
