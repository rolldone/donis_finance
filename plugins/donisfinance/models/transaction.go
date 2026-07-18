package models

import "time"

// Transaction represents a financial income/expense/transfer record.
type Transaction struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID       string    `gorm:"type:uuid;not null;index" json:"member_id"`
	AccountID      *string   `gorm:"type:uuid;index" json:"account_id,omitempty"`
	ToAccountID    *string   `gorm:"type:uuid;index" json:"to_account_id,omitempty"`
	CategoryID     *string   `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Amount         int64     `gorm:"not null" json:"amount"`
	Type           string    `gorm:"not null" json:"type"` // income | expense | transfer
	Description    string    `gorm:"default:''" json:"description"`
	Notes          string    `gorm:"default:''" json:"notes"`
	AttachmentPath string    `gorm:"default:''" json:"attachment_path"`
	Date           string    `gorm:"type:date;not null;default:CURRENT_DATE" json:"date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Member      *Member      `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Account     *Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ToAccount   *Account     `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
	Category    *Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Transaction) TableName() string { return "transactions" }
