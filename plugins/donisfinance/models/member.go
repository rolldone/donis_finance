package models

import (
	"time"
)

// MemberStatus represents the account status.
type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusPending  MemberStatus = "pending"
	MemberStatusRejected MemberStatus = "rejected"
)

// Member represents a financial data owner (istri, anak, keluarga).
type Member struct {
	ID        string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string       `gorm:"not null" json:"name"`
	Username  string       `gorm:"uniqueIndex;not null" json:"username"`
	Password  string       `gorm:"not null" json:"-"` // bcrypt hash
	Email     string       `gorm:"" json:"email,omitempty"`
	Status    MemberStatus `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// Members belong to an admin
	AdminID string `gorm:"type:uuid;not null;index" json:"admin_id"`
}

func (Member) TableName() string {
	return "members"
}
