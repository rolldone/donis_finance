package services

import (
	"fmt"

	"go_framework/plugins/donisfinance/models"

	"gorm.io/gorm"
)

// UpdateAccount updates an existing account's name/type.
func UpdateAccount(db *gorm.DB, id, memberID, name, acctType string) (*AccountResult, error) {
	var a models.Account
	if err := db.First(&a, "id = ? AND member_id = ?", id, memberID).Error; err != nil {
		return nil, fmt.Errorf("account not found")
	}

	if name == "" {
		return nil, fmt.Errorf("account name is required")
	}

	validTypes := map[string]bool{"cash": true, "bank": true, "e_wallet": true, "savings": true, "investment": true}
	if acctType != "" && !validTypes[acctType] {
		return nil, fmt.Errorf("invalid account type: must be cash, bank, e_wallet, savings, or investment")
	}

	updates := map[string]interface{}{
		"name": name,
	}
	if acctType != "" {
		updates["type"] = acctType
	}

	if err := db.Model(&a).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update account: %w", err)
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
