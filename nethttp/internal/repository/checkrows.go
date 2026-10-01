package repository

import (
	"database/sql"
	"fmt"

	"github.com/maestroharyor/go-webserver/nethttp/internal/models"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanNote(s rowScanner) (*models.Note, error) {
	var note models.Note
	if err := s.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, err
	}
	return &note, nil
}

func checkRowsAffected(result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
