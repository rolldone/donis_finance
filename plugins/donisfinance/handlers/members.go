package handlers

import (
	"net/http"

	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberHandler groups member management endpoints.
type MemberHandler struct {
	DB *gorm.DB
}

// NewMemberHandler creates a MemberHandler.
func NewMemberHandler(db *gorm.DB) *MemberHandler {
	return &MemberHandler{DB: db}
}

// CreateMemberRequest is the JSON body for creating a member.
type CreateMemberRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ListMembers handles GET /admin/plugins/donisfinance/members
func (h *MemberHandler) ListMembers(c *gin.Context) {
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	type MemberRow struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
	}

	var members []MemberRow
	if err := h.DB.Table("members").
		Select("id, name, username, email, status, created_at").
		Where("admin_id = ?", adminID).
		Order("created_at DESC").
		Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// CreateMember handles POST /admin/plugins/donisfinance/members
func (h *MemberHandler) CreateMember(c *gin.Context) {
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, username, and password are required"})
		return
	}

	result, err := services.CreateMember(h.DB, adminID, req.Name, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "member created",
		"member":  result,
	})
}

// UpdateMember handles PUT /admin/plugins/donisfinance/members/:id
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	memberID := c.Param("id")
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := services.UpdateMember(h.DB, adminID, memberID, req.Name, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member": result})
}

// DeleteMember handles DELETE /admin/plugins/donisfinance/members/:id
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	memberID := c.Param("id")
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result := h.DB.Table("members").Where("id = ? AND admin_id = ?", memberID, adminID).Delete(nil)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete member"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member deleted"})
}

// ApproveMember handles PATCH /admin/plugins/donisfinance/members/:id/approve
func (h *MemberHandler) ApproveMember(c *gin.Context) {
	memberID := c.Param("id")
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := services.ApproveMember(h.DB, adminID, memberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member approved"})
}

// RejectMember handles PATCH /admin/plugins/donisfinance/members/:id/reject
func (h *MemberHandler) RejectMember(c *gin.Context) {
	memberID := c.Param("id")
	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := services.RejectMember(h.DB, adminID, memberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member rejected"})
}
