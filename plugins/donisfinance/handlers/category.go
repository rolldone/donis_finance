package handlers

import (
	"net/http"

	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CategoryHandler handles category CRUD endpoints.
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler creates a new handler.
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// CreateCategoryRequest is the JSON body for creating/updating a category.
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"` // income | expense
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// CreateCategory godoc
// POST /admin/plugins/donisfinance/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and type are required"})
		return
	}

	result, err := services.CreateCategory(h.db, req.Name, req.Type, req.Icon, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"category": result})
}

// UpdateCategory godoc
// PUT /admin/plugins/donisfinance/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := services.UpdateCategory(h.db, id, req.Name, req.Type, req.Icon, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": result})
}

// DeleteCategory godoc
// DELETE /admin/plugins/donisfinance/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteCategory(h.db, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
