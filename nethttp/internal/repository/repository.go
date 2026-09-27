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
	result, err := r.db.ExecContext(ctx, `INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, note.Title, note.Content)

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

func (r *NoteRepository) GetNotes(ctx context.Context) ([]*models.Note, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM notes`)
	if err != nil {
		return nil, fmt.Errorf("get notes error: %w", err)
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan notes error: %w", err)
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (r *NoteRepository) GetNoteById(ctx context.Context, id int) (*models.Note, error) {
	row := r.db.QueryRowContext(ctx, `SELECT * from notes where id = ?`, id)
	var note models.Note

	if err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, fmt.Errorf("get notes by if error: %w", err)
	}

	return &note, nil
}

func (r *NoteRepository) UpdateNoteById(ctx context.Context, id int, note *models.Note) error {
	result, err := r.db.ExecContext(ctx, `UPDATE notes SET title = ?, content = ? where id = ?`, note.Title, note.Content, id)
	if err != nil {
		return fmt.Errorf("update note by id error: %w", err)
	}

	return checkRowsAffected(result)
}

func (r *NoteRepository) DeleteNoteById(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE from notes where id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete note by id error: %w", err)
	}

	return checkRowsAffected(result)
}
