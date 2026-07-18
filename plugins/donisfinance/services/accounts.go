package services

import (
	"fmt"
	"log"

	"go_framework/plugins/donisfinance/models"

	"gorm.io/gorm"
)

// UpdateAccount updates an existing account's name/type/balance.
// All fields except id/memberID are optional — pass zero value to skip.
// If balance is provided (non-zero) and differs from current, an audit trail is recorded.
func UpdateAccount(db *gorm.DB, id, memberID, name, acctType string, balance *int64, balanceReason string) (*AccountResult, error) {
	var a models.Account
	if err := db.First(&a, "id = ? AND member_id = ?", id, memberID).Error; err != nil {
		return nil, fmt.Errorf("account not found")
	}

	validTypes := map[string]bool{"cash": true, "bank": true, "e_wallet": true, "savings": true, "investment": true}
	if acctType != "" && !validTypes[acctType] {
		return nil, fmt.Errorf("invalid account type: must be cash, bank, e_wallet, savings, or investment")
	}

	// Build updates map (only non-empty fields)
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if acctType != "" {
		updates["type"] = acctType
	}

	// Apply name/type updates
	if len(updates) > 0 {
		if err := db.Model(&a).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update account: %w", err)
		}
	}

	// Handle balance update (with audit trail)
	if balance != nil && *balance != a.Balance {
		reason := balanceReason
		if reason == "" {
			reason = "Manual adjustment via account-update"
		}
		adj := models.BalanceAdjustment{
			AccountID:  id,
			OldBalance: a.Balance,
			NewBalance: *balance,
			Reason:     reason,
		}
		if err := db.Create(&adj).Error; err != nil {
			log.Printf("[donisfinance] warning: failed to record balance adjustment audit: %v", err)
		}
		if err := db.Model(&a).Update("balance", *balance).Error; err != nil {
			return nil, fmt.Errorf("update balance: %w", err)
		}
	}

	// Reload to get updated values
	db.First(&a, "id = ?", id)

	return &AccountResult{
		ID:       a.ID,
		MemberID: a.MemberID,
		Name:     a.Name,
		Type:     a.Type,
		Balance:  a.Balance,
	}, nil
}

// DeleteAccount removes an account by ID (checks ownership).
func DeleteAccount(db *gorm.DB, id, memberID string) error {
	r := db.Where("member_id = ?", memberID).Delete(&models.Account{}, "id = ?", id)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}
