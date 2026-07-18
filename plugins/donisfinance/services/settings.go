package services

import (
	"os"

	"go_framework/internal/mail"
	"go_framework/plugins/donisfinance/models"

	"gorm.io/gorm"
)

// ─── SMTP Config ──────────────────────────────────────────────────────────────

// SMTPConfig represents the SMTP settings stored in DB.
type SMTPConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pass        string `json:"pass"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	UseTLS      bool   `json:"use_tls"`
	UseStartTLS bool   `json:"use_starttls"`
	SkipVerify  bool   `json:"skip_verify"`
}

// SMTPConfigKey constants for DB storage.
const (
	SMTPKeyHost        = "smtp_host"
	SMTPKeyPort        = "smtp_port"
	SMTPKeyUser        = "smtp_user"
	SMTPKeyPass        = "smtp_pass"
	SMTPKeyFromEmail   = "smtp_from_email"
	SMTPKeyFromName    = "smtp_from_name"
	SMTPKeyUseTLS      = "smtp_use_tls"
	SMTPKeyUseStartTLS = "smtp_use_starttls"
	SMTPKeySkipVerify  = "smtp_skip_verify"
)

// GetSMTPConfig loads SMTP config from DB. Returns nil if no config stored.
func GetSMTPConfig(db *gorm.DB) (*SMTPConfig, error) {
	var rows []models.Setting
	if err := db.Model(&models.Setting{}).Where("key LIKE 'smtp_%'").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	cfg := &SMTPConfig{}
	for _, r := range rows {
		switch r.Key {
		case SMTPKeyHost:
			cfg.Host = r.Value
		case SMTPKeyPort:
			cfg.Port = r.Value
		case SMTPKeyUser:
			cfg.User = r.Value
		case SMTPKeyPass:
			cfg.Pass = r.Value
		case SMTPKeyFromEmail:
			cfg.FromEmail = r.Value
		case SMTPKeyFromName:
			cfg.FromName = r.Value
		case SMTPKeyUseTLS:
			cfg.UseTLS = r.Value == "true"
		case SMTPKeyUseStartTLS:
			cfg.UseStartTLS = r.Value == "true"
		case SMTPKeySkipVerify:
			cfg.SkipVerify = r.Value == "true"
		}
	}
	return cfg, nil
}

// SaveSMTPConfig saves SMTP config to DB.
// If all core fields (host, port, user, pass, from_email, from_name) are empty,
// the existing settings are removed from DB (clear override).
func SaveSMTPConfig(db *gorm.DB, cfg *SMTPConfig) error {
	// Check if this is a "clear" request — all core fields are empty
	isClear := cfg.Host == "" && cfg.Port == "" && cfg.User == "" && cfg.Pass == "" &&
		cfg.FromEmail == "" && cfg.FromName == ""

	if isClear {
		return db.Where("key LIKE 'smtp_%'").Delete(&models.Setting{}).Error
	}

	pairs := []models.Setting{
		{Key: SMTPKeyHost, Value: cfg.Host},
		{Key: SMTPKeyPort, Value: cfg.Port},
		{Key: SMTPKeyUser, Value: cfg.User},
		{Key: SMTPKeyPass, Value: cfg.Pass},
		{Key: SMTPKeyFromEmail, Value: cfg.FromEmail},
		{Key: SMTPKeyFromName, Value: cfg.FromName},
	}
	if cfg.UseTLS {
		pairs = append(pairs, models.Setting{Key: SMTPKeyUseTLS, Value: "true"})
	} else {
		pairs = append(pairs, models.Setting{Key: SMTPKeyUseTLS, Value: "false"})
	}
	if cfg.UseStartTLS {
		pairs = append(pairs, models.Setting{Key: SMTPKeyUseStartTLS, Value: "true"})
	} else {
		pairs = append(pairs, models.Setting{Key: SMTPKeyUseStartTLS, Value: "false"})
	}
	if cfg.SkipVerify {
		pairs = append(pairs, models.Setting{Key: SMTPKeySkipVerify, Value: "true"})
	} else {
		pairs = append(pairs, models.Setting{Key: SMTPKeySkipVerify, Value: "false"})
	}

	for _, s := range pairs {
		if err := db.Model(&models.Setting{}).Where("key = ?", s.Key).Assign(s).FirstOrCreate(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

// ToMailConfig converts service SMTPConfig to mail.SMTPConfig.
func (c *SMTPConfig) ToMailConfig() *mail.SMTPConfig {
	if c == nil {
		return nil
	}
	return &mail.SMTPConfig{
		Host:        c.Host,
		Port:        c.Port,
		User:        c.User,
		Pass:        c.Pass,
		FromEmail:   c.FromEmail,
		FromName:    c.FromName,
		UseTLS:      c.UseTLS,
		UseStartTLS: c.UseStartTLS,
		SkipVerify:  c.SkipVerify,
	}
}

// GetEnvSMTPConfig reads SMTP config from environment variables (fallback).
func GetEnvSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:        os.Getenv("SMTP_HOST"),
		Port:        os.Getenv("SMTP_PORT"),
		User:        os.Getenv("SMTP_USER"),
		Pass:        os.Getenv("SMTP_PASS"),
		FromEmail:   os.Getenv("SMTP_FROM_EMAIL"),
		FromName:    os.Getenv("SMTP_FROM_NAME"),
		UseTLS:      os.Getenv("SMTP_USE_TLS") == "true",
		UseStartTLS: os.Getenv("SMTP_STARTTLS") == "true",
		SkipVerify:  os.Getenv("SMTP_TLS_SKIP_VERIFY") == "true",
	}
}
