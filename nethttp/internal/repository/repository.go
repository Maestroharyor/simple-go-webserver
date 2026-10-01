package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/maestroharyor/go-webserver/nethttp/internal/models"
)

var ErrNotFound = errors.New("note not found")

const noteColumns = `id, title, content, created_at, updated_at`

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, `INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)`, note.Title, note.Content, now, now)

	if err != nil {
		return fmt.Errorf("create note error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get note id: %w", err)
	}

	note.ID = int(id)
	note.CreatedAt = now
	note.UpdatedAt = now

	return nil
}

func (r *NoteRepository) GetNotes(ctx context.Context) ([]*models.Note, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+noteColumns+` FROM notes`)
	if err != nil {
		return nil, fmt.Errorf("get notes error: %w", err)
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)

	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			return nil, fmt.Errorf("scan note error: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows error: %w", err)
	}

	return notes, nil
}

func (r *NoteRepository) GetNoteById(ctx context.Context, id int) (*models.Note, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+noteColumns+` from notes where id = ?`, id)
	note, err := scanNote(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get note by id error: %w", err)
	}

	return note, nil
}

func (r *NoteRepository) UpdateNoteById(ctx context.Context, id int, note *models.Note) (*models.Note, error) {
	result, err := r.db.ExecContext(ctx, `UPDATE notes SET title = ?, content = ?, updated_at = ? where id = ?`, note.Title, note.Content, time.Now().UTC(), id)
	if err != nil {
		return nil, fmt.Errorf("update note by id error: %w", err)
	}
	if err := checkRowsAffected(result); err != nil {
		return nil, err
	}

	return r.GetNoteById(ctx, id)
}

func (r *NoteRepository) DeleteNoteById(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE from notes where id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete note by id error: %w", err)
	}

	return checkRowsAffected(result)
}
