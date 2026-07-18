package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ─── DTO ──────────────────────────────────────────────────────────────────────

// BudgetResult is returned after setting a budget.
type BudgetResult struct {
	ID         string  `json:"id"`
	MemberID   string  `json:"member_id"`
	CategoryID *string `json:"category_id,omitempty"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Amount     int64   `json:"amount"`
}

// BudgetStatus shows budget limit vs actual spending.
type BudgetStatus struct {
	BudgetResult
	CategoryName string `json:"category_name"`
	CategoryColor string `json:"category_color"`
	Spent       int64  `json:"spent"`
	Remaining   int64  `json:"remaining"`
	Percentage  int    `json:"percentage"` // 0-100
}

// ─── Operations ───────────────────────────────────────────────────────────────

// SetBudget creates or updates a budget for a member/category/month/year.
func SetBudget(db *gorm.DB, memberID string, categoryID *string, month, year int, amount int64) (*BudgetResult, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("month must be 1-12")
	}
	if year < 2000 || year > 2100 {
		return nil, fmt.Errorf("invalid year")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Verify member exists
	var mCount int64
	db.Table("members").Where("id = ?", memberID).Count(&mCount)
	if mCount == 0 {
		return nil, fmt.Errorf("member not found")
	}

	// Upsert via raw SQL to avoid GORM FK issues
	now := time.Now()
	var id string
	err := db.Raw(`
		INSERT INTO budgets (member_id, category_id, month, year, amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (member_id, category_id, month, year)
		DO UPDATE SET amount = EXCLUDED.amount, updated_at = EXCLUDED.updated_at
		RETURNING id
	`, memberID, categoryID, month, year, amount, now, now).Scan(&id).Error
	if err != nil {
		return nil, fmt.Errorf("set budget: %w", err)
	}

	return &BudgetResult{
		ID:         id,
		MemberID:   memberID,
		CategoryID: categoryID,
		Month:      month,
		Year:       year,
		Amount:     amount,
	}, nil
}

// GetBudgetStatus returns budget vs actual spending for a member/month/year.
func GetBudgetStatus(db *gorm.DB, memberID string, month, year int) ([]BudgetStatus, error) {
	var results []BudgetStatus

	rows, err := db.Table("budgets").
		Select(`budgets.id, budgets.member_id, budgets.category_id, budgets.month, budgets.year, 
		        budgets.amount,
		        COALESCE(categories.name, '') as category_name,
		        COALESCE(categories.color, '') as category_color,
		        COALESCE(SUM(transactions.amount), 0) as spent`).
		Joins("LEFT JOIN categories ON categories.id = budgets.category_id").
		Joins(`LEFT JOIN transactions ON transactions.category_id = budgets.category_id
		       AND transactions.member_id = budgets.member_id
		       AND transactions.type = 'expense'
		       AND EXTRACT(YEAR FROM transactions.date) = budgets.year
		       AND EXTRACT(MONTH FROM transactions.date) = budgets.month`).
		Where("budgets.member_id = ? AND budgets.month = ? AND budgets.year = ?", memberID, month, year).
		Group("budgets.id, budgets.member_id, budgets.category_id, budgets.month, budgets.year, budgets.amount, categories.name, categories.color").
		Order("categories.name ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s BudgetStatus
		if err := rows.Scan(&s.ID, &s.MemberID, &s.CategoryID, &s.Month, &s.Year,
			&s.Amount, &s.CategoryName, &s.CategoryColor, &s.Spent); err != nil {
			return nil, err
		}
		s.Remaining = s.Amount - s.Spent
		if s.Remaining < 0 {
			s.Remaining = 0
		}
		if s.Amount > 0 {
			pct := int((s.Spent * 100) / s.Amount)
			if pct > 100 {
				pct = 100
			}
			s.Percentage = pct
		}
		results = append(results, s)
	}

	if results == nil {
		results = []BudgetStatus{}
	}
	return results, nil
}

// CheckBudgetSpending checks if a new expense would exceed the budget.
// Returns (budgetID, remaining, isOverBudget, error).
func CheckBudgetSpending(db *gorm.DB, memberID, categoryID string, month, year int, newAmount int64) (string, int64, bool, error) {
	type BudgetRow struct {
		ID     string
		Amount int64
	}
	var row BudgetRow
	var err error

	// Split query to avoid PostgreSQL "could not determine data type of parameter $3"
	// when categoryID is empty string (NULL vs '' ambiguity in GORM binding).
	if categoryID == "" {
		err = db.Table("budgets").Select("id, amount").
			Where("member_id = ? AND category_id IS NULL AND month = ? AND year = ?",
				memberID, month, year).Scan(&row).Error
	} else {
		err = db.Table("budgets").Select("id, amount").
			Where("member_id = ? AND category_id = ? AND month = ? AND year = ?",
				memberID, categoryID, month, year).Scan(&row).Error
	}
	if err != nil || row.ID == "" {
		return "", 0, false, nil // No budget set = no warning
	}

	var spent int64
	db.Table("transactions").
		Where("member_id = ? AND category_id = ? AND type = 'expense' AND EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?",
			memberID, categoryID, year, month).
		Select("COALESCE(SUM(amount), 0)").Scan(&spent)

	remaining := row.Amount - spent - newAmount
	return row.ID, remaining, remaining < 0, nil
}

// DeleteBudget removes a budget by ID.
func DeleteBudget(db *gorm.DB, id string) error {
	r := db.Exec("DELETE FROM budgets WHERE id = ?", id)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("budget not found")
	}
	return nil
}
