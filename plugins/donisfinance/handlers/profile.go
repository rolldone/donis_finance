package handlers

import (
	"net/http"

	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProfileHandler handles profile and password endpoints.
type ProfileHandler struct {
	db *gorm.DB
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

// GET /api/member/profile
func (h *ProfileHandler) GetMemberProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	profile, err := services.GetMemberProfile(h.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// PUT /api/member/profile
func (h *ProfileHandler) UpdateMemberProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := services.UpdateMemberProfile(h.db, userID, req.Name, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// PUT /api/member/password
func (h *ProfileHandler) ChangeMemberPassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "old_password and new_password required"})
		return
	}
	if err := services.ChangeMemberPassword(h.db, userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// GET /api/admin/profile
func (h *ProfileHandler) GetAdminProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	profile, err := services.GetAdminProfile(h.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// PUT /api/admin/profile
func (h *ProfileHandler) UpdateAdminProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := services.UpdateAdminProfile(h.db, userID, req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// PUT /api/admin/password
func (h *ProfileHandler) ChangeAdminPassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "old_password and new_password required"})
		return
	}
	if err := services.ChangeAdminPassword(h.db, userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
