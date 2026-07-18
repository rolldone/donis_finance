package handlers

import (
	"net/http"
	"strconv"

	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BudgetHandler handles budget-related endpoints.
type BudgetHandler struct {
	db *gorm.DB
}

// NewBudgetHandler creates a new handler.
func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

// SetBudget godoc
// POST /admin/plugins/donisfinance/budgets
// POST /api/member/budgets
func (h *BudgetHandler) SetBudget(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("user_type")

	var req struct {
		MemberID   string `json:"member_id"`
		CategoryID string `json:"category_id"`
		Month      int    `json:"month"`
		Year       int    `json:"year"`
		Amount     int64  `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	memberID := userID
	if role == "admin" && req.MemberID != "" {
		memberID = req.MemberID
	}

	var catID *string
	if req.CategoryID != "" {
		catID = &req.CategoryID
	}

	result, err := services.SetBudget(h.db, memberID, catID, req.Month, req.Year, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"budget": result})
}

// GetBudgetStatus godoc
// GET /admin/plugins/donisfinance/budgets/status?member_id=x&month=7&year=2026
// GET /api/member/budgets/status?month=7&year=2026
func (h *BudgetHandler) GetBudgetStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("user_type")

	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	if month < 1 || month > 12 || year < 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month (1-12) and year are required"})
		return
	}

	memberID := userID
	if role == "admin" {
		if m := c.Query("member_id"); m != "" {
			memberID = m
		}
	}

	results, err := services.GetBudgetStatus(h.db, memberID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budgets": results})
}

// DeleteBudget godoc
// DELETE /admin/plugins/donisfinance/budgets/:id
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteBudget(h.db, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
