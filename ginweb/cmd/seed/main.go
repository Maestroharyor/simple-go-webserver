package main

import (
	"fmt"
	"log"
	"time"

	"github.com/maestroharyor/go-webserver/ginweb/internal/env"
	"github.com/maestroharyor/go-webserver/ginweb/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const noteCount = 100

func main() {
	db, err := gorm.Open(sqlite.Open(env.GetEnv("DATABASE_URL", "notes.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Note{}); err != nil {
		log.Fatal(err)
	}

	notes := buildNotes(noteCount, time.Now())
	if err := db.CreateInBatches(notes, 50).Error; err != nil {
		log.Fatal(err)
	}

	log.Printf("seeded %d notes", len(notes))
}

// Timestamps are spaced a minute apart so ordering by created_at is
// deterministic when testing pagination.
func buildNotes(n int, now time.Time) []models.Note {
	notes := make([]models.Note, n)
	for i := range notes {
		ts := now.Add(time.Duration(i-n) * time.Minute)
		notes[i] = models.Note{
			Title:     fmt.Sprintf("Seed note %03d", i+1),
			Content:   fmt.Sprintf("This is the content of seed note %d.", i+1),
			CreatedAt: ts,
			UpdatedAt: ts,
		}
	}
	return notes
}
