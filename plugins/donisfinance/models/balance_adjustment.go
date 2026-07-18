package models

import "time"

// BalanceAdjustment records an audit trail when an account balance is manually adjusted.
type BalanceAdjustment struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AccountID  string    `gorm:"type:uuid;not null;index" json:"account_id"`
	OldBalance int64     `gorm:"not null;default:0" json:"old_balance"`
	NewBalance int64     `gorm:"not null;default:0" json:"new_balance"`
	Reason     string    `gorm:"not null;default:''" json:"reason"`
	AdjustedAt time.Time `gorm:"not null;default:now()" json:"adjusted_at"`
}

func (BalanceAdjustment) TableName() string { return "balance_adjustments" }
