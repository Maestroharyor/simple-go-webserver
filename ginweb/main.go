package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maestroharyor/go-webserver/ginweb/internal/env"
	"github.com/maestroharyor/go-webserver/ginweb/internal/handlers"
	"github.com/maestroharyor/go-webserver/ginweb/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open(env.GetEnv("DATABASE_URL", "notes.db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Note{})
	userHandler := handlers.NewUserHandler(db)

	r := gin.Default()

	r.GET("/", handlers.RootHandler)
	r.POST("/notes", userHandler.CreateNoteHandler)
	r.GET("/notes", userHandler.ListNotesHandler)
	r.GET("/notes/:id", userHandler.GetSingleNoteHandler)
	r.PUT("/notes/:id", userHandler.UpdateSingleNoteHandler)
	r.DELETE("/notes/:id", userHandler.DeleteSingleNoteHandler)

	r.Run()

}
