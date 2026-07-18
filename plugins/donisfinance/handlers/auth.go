package handlers

import (
	"fmt"
	"net/http"
	"os"

	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler groups auth endpoints.
type AuthHandler struct {
	DB *gorm.DB
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// AdminLoginRequest is the JSON body for admin login.
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// MemberLoginRequest is the JSON body for member login.
type MemberLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin handles POST /admin/plugins/donisfinance/auth/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	session, err := services.LoginAdmin(h.DB, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      session.Token,
		"expires_at": session.ExpiresAt,
		"user_id":    session.UserID,
		"username":   session.Username,
		"role":       session.Role,
	})
}

// MemberLogin handles POST /api/member/auth/login
func (h *AuthHandler) MemberLogin(c *gin.Context) {
	var req MemberLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	session, err := services.LoginMember(h.DB, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      session.Token,
		"expires_at": session.ExpiresAt,
		"user_id":    session.UserID,
		"username":   session.Username,
		"role":       session.Role,
	})
}

// ─── Public Register ──────────────────────────────────────────────────────────

// RegisterRequest is the JSON body for public registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /api/member/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, and password are required"})
		return
	}

	result, err := services.RegisterMember(h.DB, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil. Tunggu aktivasi oleh admin.",
		"member":  result,
	})
}

// ─── Forgot / Reset Password ──────────────────────────────────────────────────

// ForgotPasswordRequest is the JSON body for forgot password.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// ForgotPassword handles POST /api/member/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	result, err := services.ForgotPassword(h.DB, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Always return success to avoid email enumeration
	if result == nil || result.ResetToken == "" {
		c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, link reset akan dikirim"})
		return
	}

	// Send email with reset link
	resetLink := fmt.Sprintf("%s/member/auth/reset-password?token=%s", getAppURL(), result.ResetToken)
	_ = services.SendResetPasswordEmail(result.MemberEmail, result.MemberName, resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, link reset akan dikirim"})
}

// ResetPasswordRequest is the JSON body for reset password.
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ResetPassword handles POST /api/member/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and password are required"})
		return
	}

	if err := services.ResetPassword(h.DB, req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset. Silakan login."})
}

// getAppURL returns the base URL of the frontend from env or falls back.
func getAppURL() string {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8202"
	}
	return appURL
}
