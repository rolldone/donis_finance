package models

import (
	"time"
)

// Admin represents a system operator / pegawai donis_finance.
type Admin struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // bcrypt hash, never exposed
	Email     string    `gorm:"" json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admins"
}
