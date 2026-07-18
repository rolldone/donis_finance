package models

import "time"

// Budget represents a monthly spending limit per category for a member.
type Budget struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID   string    `gorm:"type:uuid;not null;index" json:"member_id"`
	CategoryID *string   `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Month      int       `gorm:"not null" json:"month"`
	Year       int       `gorm:"not null" json:"year"`
	Amount     int64     `gorm:"not null" json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Budget) TableName() string { return "budgets" }
