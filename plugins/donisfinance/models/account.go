package models

import "time"

// Account represents a financial account owned by a member.
type Account struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID  string    `gorm:"type:uuid;not null;index" json:"member_id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null;default:cash" json:"type"` // cash | bank | e_wallet | savings | investment
	Balance   int64     `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Member *Member `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

func (Account) TableName() string { return "accounts" }
