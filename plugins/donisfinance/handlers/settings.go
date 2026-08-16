package handlers

import (
	"log/slog"
	"net/http"

	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsHandler handles SMTP + notification configuration endpoints.
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler creates a new handler.
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// SMTPConfigResponse is the full settings response.
type SMTPConfigResponse struct {
	SMTP       services.SMTPConfig  `json:"smtp"`
	EnvSMTP    *services.SMTPConfig `json:"env_smtp"`
	Override   bool                 `json:"override"`
	NotifEmail string               `json:"notif_email"`
}

// GetSMTPConfig godoc
// GET /api/admin/settings/smtp
func (h *SettingsHandler) GetSMTPConfig(c *gin.Context) {
	cfg, err := services.GetSMTPConfig(h.db)
	if err != nil {
		slog.Error("failed to get SMTP config", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	envCfg := services.GetEnvSMTPConfig()
	var dbCfg services.SMTPConfig
	if cfg != nil {
		dbCfg = *cfg
	}

	notifEmail := services.GetNotifEmail(h.db)

	c.JSON(http.StatusOK, gin.H{
		"smtp":        dbCfg,
		"env_smtp":    envCfg,
		"override":    cfg != nil,
		"notif_email": notifEmail,
	})
}

// SaveSMTPConfig godoc
// PUT /api/admin/settings/smtp
func (h *SettingsHandler) SaveSMTPConfig(c *gin.Context) {
	var body struct {
		SMTP       services.SMTPConfig `json:"smtp"`
		NotifEmail string              `json:"notif_email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := services.SaveSMTPConfig(h.db, &body.SMTP); err != nil {
		slog.Error("failed to save SMTP config", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := services.SaveNotifEmail(h.db, body.NotifEmail); err != nil {
		slog.Error("failed to save notification email", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Re-read to get accurate state
	dbCfg, _ := services.GetSMTPConfig(h.db)
	envCfg := services.GetEnvSMTPConfig()
	notifEmail := services.GetNotifEmail(h.db)

	var respCfg services.SMTPConfig
	if dbCfg != nil {
		respCfg = *dbCfg
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Settings saved",
		"smtp":        respCfg,
		"env_smtp":    envCfg,
		"override":    dbCfg != nil,
		"notif_email": notifEmail,
	})
}
