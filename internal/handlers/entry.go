package handlers

import (
	"net/http"
	"strconv"

	"diary-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EntryHandler struct {
	db *gorm.DB
}

func NewEntryHandler(db *gorm.DB) *EntryHandler {
	return &EntryHandler{db: db}
}

func (h *EntryHandler) GetEntries(c *gin.Context) {
	userID, _ := c.Get("userID")

	var entries []models.Entry
	h.db.Where("user_id = ?", userID).Order("created_at desc").Find(&entries)
	c.JSON(http.StatusOK, entries)
}

func (h *EntryHandler) CreateEntry(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry := models.Entry{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	h.db.Create(&entry)
	c.JSON(http.StatusCreated, entry)
}

func (h *EntryHandler) GetEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("userID")

	var entry models.Entry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) UpdateEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("userID")

	var entry models.Entry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var req models.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.Title = req.Title
	entry.Content = req.Content
	h.db.Save(&entry)

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) DeleteEntry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("userID")

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Entry{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
