package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maestroharyor/go-webserver/ginweb/internal/models"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pong",
		"data":    nil,
	})
}

func (h *UserHandler) HealthStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"data":    nil,
	})

}

func (h *UserHandler) CreateNoteHandler(c *gin.Context) {
	var req CreateNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create note",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notes created successfully",
		"data":    note,
	})
}

func (h *UserHandler) ListNotesHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "4"))
	if perPage < 1 {
		perPage = 1
	}

	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	//Implement pagination
	total, err := gorm.G[models.Note](h.db).Count(c.Request.Context(), "*")
	if err != nil {
		total = 0
	}

	notes, err := gorm.G[models.Note](h.db).Order("created_at DESC").Limit(perPage).Offset(offset).Find(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to fetch notes",
			"error":   err.Error(),
		})
		return
	}

	totalPages := int(total+int64(perPage)-1) / perPage
	hasNext := page < totalPages
	hasPrevious := page > 1

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notes fetched successfully",
		"data":    notes,
		"pagination": gin.H{
			"page":         page,
			"per_page":     perPage,
			"total":        total,
			"total_pages":  totalPages,
			"has_next":     hasNext,
			"has_previous": hasPrevious,
		},
	})
}

func (h *UserHandler) GetSingleNoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid note ID",
			"error":   err.Error(),
		})
		return
	}

	note, err := gorm.G[models.Note](h.db).Where("id = ?", id).First(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Note not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Note fetched successfully",
		"data":    note,
	})
}

func (h *UserHandler) UpdateSingleNoteHandler(c *gin.Context) {
	var req UpdateNoteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"data":    nil,
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid note ID",
			"error":   err.Error(),
		})
		return
	}

	row, err := gorm.G[models.Note](h.db).Where("id = ?", id).Updates(c.Request.Context(), models.Note{
		Title:   req.Title,
		Content: req.Content,
	})
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Note not found",
			"data":    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update note",
			"error":   err.Error(),
		})
		return
	}

	note, err := gorm.G[models.Note](h.db).Where("id = ?", id).First(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get note",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Note updated successfully",
		"data":    note,
	})
}

func (h *UserHandler) DeleteSingleNoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid note ID",
			"error":   err.Error(),
		})
		return
	}

	rows, error := gorm.G[models.Note](h.db).Where("id = ?", id).Delete(c.Request.Context())
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete note",
			"error":   error.Error(),
		})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Note not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Note deleted successfully",
		"data":    nil,
	})
}
