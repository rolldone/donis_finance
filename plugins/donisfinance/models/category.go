package models

import "time"

// Category represents a transaction category (income or expense).
type Category struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null;default:expense" json:"type"`   // income | expense
	Icon      string    `gorm:"default:''" json:"icon"`
	Color     string    `gorm:"default:'#6b7280'" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Category) TableName() string { return "categories" }
