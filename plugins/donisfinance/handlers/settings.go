package handlers

import (
	"net/http"

	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsHandler handles SMTP configuration endpoints.
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler creates a new handler.
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// GetSMTPConfig godoc
// GET /api/admin/settings/smtp
func (h *SettingsHandler) GetSMTPConfig(c *gin.Context) {
	cfg, err := services.GetSMTPConfig(h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	envCfg := services.GetEnvSMTPConfig()
	var dbCfg services.SMTPConfig
	if cfg != nil {
		dbCfg = *cfg
	}

	c.JSON(http.StatusOK, gin.H{
		"smtp":     dbCfg,
		"env_smtp": envCfg,
		"override": cfg != nil,
	})
}

// SaveSMTPConfig godoc
// PUT /api/admin/settings/smtp
func (h *SettingsHandler) SaveSMTPConfig(c *gin.Context) {
	var cfg services.SMTPConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := services.SaveSMTPConfig(h.db, &cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Re-read to get accurate state
	dbCfg, _ := services.GetSMTPConfig(h.db)
	envCfg := services.GetEnvSMTPConfig()

	var respCfg services.SMTPConfig
	if dbCfg != nil {
		respCfg = *dbCfg
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "SMTP settings saved",
		"smtp":     respCfg,
		"env_smtp": envCfg,
		"override": dbCfg != nil,
	})
}
