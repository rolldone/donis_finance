package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"go_framework/plugins/donisfinance/models"

	"gorm.io/gorm"
)

// ─── DTO ──────────────────────────────────────────────────────────────────────

// TransactionResult is returned after creating a transaction.
type TransactionResult struct {
	ID             string  `json:"id"`
	MemberID       string  `json:"member_id"`
	AccountID      *string `json:"account_id,omitempty"`
	ToAccountID    *string `json:"to_account_id,omitempty"`
	CategoryID     *string `json:"category_id,omitempty"`
	Amount         int64   `json:"amount"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	Notes          string  `json:"notes"`
	AttachmentPath string  `json:"attachment_path"`
	Date           string  `json:"date"`
	BudgetWarning  string  `json:"budget_warning,omitempty"`
}

// TransactionFilter for listing transactions.
type TransactionFilter struct {
	MemberID  string
	Month     int
	Year      int
	Type      string // optional: income | expense
	Search    string // optional: search by description
	SortBy    string // optional: date | amount | description
	SortOrder string // optional: asc | desc
	Limit     int
	Offset    int
}

// MonthlySummary groups transactions by category.
type MonthlySummary struct {
	TotalIncome   int64             `json:"total_income"`
	TotalExpense  int64             `json:"total_expense"`
	CategoryBreak []CategorySummary `json:"category_breakdown"`
}

// CategorySummary per category.
type CategorySummary struct {
	CategoryName  string `json:"category_name"`
	CategoryColor string `json:"category_color"`
	Type          string `json:"type"`
	Total         int64  `json:"total"`
	Count         int    `json:"count"`
}

// ─── Operations ───────────────────────────────────────────────────────────────

// CreateTransaction adds a new transaction record.
func CreateTransaction(db *gorm.DB, memberID string, accountID, categoryID, toAccountID *string, amount int64, txType, description, notes, attachmentPath, date string) (*TransactionResult, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if txType != "income" && txType != "expense" && txType != "transfer" {
		return nil, fmt.Errorf("type must be income, expense, or transfer")
	}

	// For transfers, require to_account_id
	if txType == "transfer" {
		if toAccountID == nil || *toAccountID == "" {
			return nil, fmt.Errorf("to_account_id is required for transfer")
		}
		if accountID != nil && toAccountID != nil && *accountID == *toAccountID {
			return nil, fmt.Errorf("source and destination accounts must be different")
		}
	}

	// If categoryID provided, verify it exists and matches type
	if categoryID != nil && *categoryID != "" {
		var cat models.Category
		if err := db.First(&cat, "id = ?", *categoryID).Error; err != nil {
			return nil, fmt.Errorf("category not found")
		}
		if txType != cat.Type && txType != "transfer" {
			return nil, fmt.Errorf("category type %s does not match transaction type %s", cat.Type, txType)
		}
	}

	// Validate date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
		}
	}

	tx := models.Transaction{
		MemberID:       memberID,
		AccountID:      accountID,
		ToAccountID:    toAccountID,
		CategoryID:     categoryID,
		Amount:         amount,
		Type:           txType,
		Description:    description,
		Notes:          notes,
		AttachmentPath: attachmentPath,
		Date:           date,
	}

	if err := db.Create(&tx).Error; err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	// Auto-update account balance
	if txType == "transfer" {
		// Deduct from source
		if accountID != nil && *accountID != "" {
			db.Model(&models.Account{}).Where("id = ?", *accountID).
				Update("balance", gorm.Expr("balance - ?", amount))
		}
		// Add to destination
		if toAccountID != nil && *toAccountID != "" {
			db.Model(&models.Account{}).Where("id = ?", *toAccountID).
				Update("balance", gorm.Expr("balance + ?", amount))
		}
	} else if accountID != nil && *accountID != "" {
		var delta int64
		if txType == "income" {
			delta = amount
		} else {
			delta = -amount
		}
		db.Model(&models.Account{}).Where("id = ?", *accountID).
			Update("balance", gorm.Expr("balance + ?", delta))
	}

	// Check budget if expense
	result := &TransactionResult{
		ID:             tx.ID,
		MemberID:       tx.MemberID,
		AccountID:      tx.AccountID,
		ToAccountID:    tx.ToAccountID,
		CategoryID:     tx.CategoryID,
		Amount:         tx.Amount,
		Type:           tx.Type,
		Description:    tx.Description,
		Notes:          tx.Notes,
		AttachmentPath: tx.AttachmentPath,
		Date:           tx.Date,
	}

	if txType == "expense" && categoryID != nil && *categoryID != "" {
		month, _ := time.Parse("2006-01-02", tx.Date)
		budgetID, remaining, over, _ := CheckBudgetSpending(db, tx.MemberID, *categoryID, int(month.Month()), month.Year(), 0)
		if budgetID != "" {
			if over {
				result.BudgetWarning = "OVER BUDGET! sisa Rp0"
			} else if remaining < amount {
				result.BudgetWarning = fmt.Sprintf("⚠️ Sisa budget: Rp%d", remaining)
			}
		}
	}

	return result, nil
}

// ListTransactionsResult includes results and total count.
type ListTransactionsResult struct {
	Transactions []map[string]interface{}
	Total        int64
}

// buildSortOrder returns a safe ORDER BY clause from sort params.
func buildSortOrder(sortBy, sortOrder string) string {
	// Whitelist sort columns to prevent SQL injection
	switch sortBy {
	case "date":
		sortBy = "transactions.date"
	case "amount":
		sortBy = "transactions.amount"
	case "description":
		sortBy = "transactions.description"
	default:
		return "transactions.date DESC, transactions.created_at DESC"
	}

	if sortOrder == "asc" {
		return sortBy + " ASC"
	}
	return sortBy + " DESC"
}

// ListTransactions returns filtered transactions with joins and total count.
func ListTransactions(db *gorm.DB, f TransactionFilter) (*ListTransactionsResult, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}

	q := db.Table("transactions").
		Select(`transactions.id, transactions.member_id, transactions.account_id, 
		        transactions.category_id, transactions.to_account_id,
		        transactions.amount, transactions.type, 
		        transactions.description, transactions.notes, transactions.attachment_path,
		        transactions.date, transactions.created_at,
		        COALESCE(categories.name, '') as category_name,
		        COALESCE(categories.color, '') as category_color,
		        COALESCE(categories.icon, '') as category_icon,
		        COALESCE(accounts.name, '') as account_name,
		        COALESCE(to_accounts.name, '') as to_account_name`).
		Joins("LEFT JOIN categories ON categories.id = transactions.category_id").
		Joins("LEFT JOIN accounts ON accounts.id = transactions.account_id").
		Joins("LEFT JOIN accounts as to_accounts ON to_accounts.id = transactions.to_account_id")

	if f.MemberID != "" {
		q = q.Where("transactions.member_id = ?", f.MemberID)
	}
	if f.Type != "" {
		q = q.Where("transactions.type = ?", f.Type)
	}
	if f.Search != "" {
		q = q.Where("LOWER(transactions.description) LIKE ?", "%"+f.Search+"%")
	}
	if f.Year > 0 {
		if f.Month > 0 {
			q = q.Where("EXTRACT(YEAR FROM transactions.date) = ? AND EXTRACT(MONTH FROM transactions.date) = ?", f.Year, f.Month)
		} else {
			q = q.Where("EXTRACT(YEAR FROM transactions.date) = ?", f.Year)
		}
	}

	// Build order clause
	orderClause := buildSortOrder(f.SortBy, f.SortOrder)

	// Count total matching rows
	var total int64
	countQ := db.Table("transactions").
		Joins("LEFT JOIN categories ON categories.id = transactions.category_id").
		Joins("LEFT JOIN accounts ON accounts.id = transactions.account_id").
		Joins("LEFT JOIN accounts as to_accounts ON to_accounts.id = transactions.to_account_id")

	if f.MemberID != "" {
		countQ = countQ.Where("transactions.member_id = ?", f.MemberID)
	}
	if f.Type != "" {
		countQ = countQ.Where("transactions.type = ?", f.Type)
	}
	if f.Search != "" {
		countQ = countQ.Where("LOWER(transactions.description) LIKE ?", "%"+f.Search+"%")
	}
	if f.Year > 0 {
		if f.Month > 0 {
			countQ = countQ.Where("EXTRACT(YEAR FROM transactions.date) = ? AND EXTRACT(MONTH FROM transactions.date) = ?", f.Year, f.Month)
		} else {
			countQ = countQ.Where("EXTRACT(YEAR FROM transactions.date) = ?", f.Year)
		}
	}
	countQ.Count(&total)

	var results []map[string]interface{}
	if err := q.Order(orderClause).
		Limit(f.Limit).Offset(f.Offset).Find(&results).Error; err != nil {
		return nil, err
	}
	return &ListTransactionsResult{Transactions: results, Total: total}, nil
}

// UpdateTransaction updates an existing transaction with balance re-calculation.
func UpdateTransaction(db *gorm.DB, id, memberID, role string, req struct {
	AccountID   string `json:"account_id"`
	ToAccountID string `json:"to_account_id"`
	CategoryID  string `json:"category_id"`
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	Date        string `json:"date"`
}) (*TransactionResult, error) {
	var tx models.Transaction
	if err := db.First(&tx, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("transaction not found")
	}

	if role != "admin" && tx.MemberID != memberID {
		return nil, fmt.Errorf("not your transaction")
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if req.Type != "income" && req.Type != "expense" && req.Type != "transfer" {
		return nil, fmt.Errorf("type must be income, expense, or transfer")
	}

	// Balance: reverse old effect
	reverseBalance := func(tType string, amount int64, acctID, toAcctID *string) {
		switch tType {
		case "income":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance - ?", amount))
			}
		case "expense":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance + ?", amount))
			}
		case "transfer":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance + ?", amount))
			}
			if toAcctID != nil && *toAcctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *toAcctID).
					Update("balance", gorm.Expr("balance - ?", amount))
			}
		}
	}

	reverseBalance(tx.Type, tx.Amount, tx.AccountID, tx.ToAccountID)

	// Apply new values
	var newAcctID, newToAcctID, newCatID *string
	if req.AccountID != "" {
		newAcctID = &req.AccountID
	}
	if req.ToAccountID != "" {
		newToAcctID = &req.ToAccountID
	}
	if req.CategoryID != "" {
		newCatID = &req.CategoryID
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	} else if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	updates := map[string]interface{}{
		"account_id":    newAcctID,
		"to_account_id": newToAcctID,
		"category_id":   newCatID,
		"amount":        req.Amount,
		"type":          req.Type,
		"description":   req.Description,
		"notes":         req.Notes,
		"date":          req.Date,
	}

	if err := db.Model(&tx).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update transaction: %w", err)
	}

	// Re-apply balance with new values
	applyBalance := func(tType string, amount int64, acctID, toAcctID *string) {
		switch tType {
		case "income":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance + ?", amount))
			}
		case "expense":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance - ?", amount))
			}
		case "transfer":
			if acctID != nil && *acctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *acctID).
					Update("balance", gorm.Expr("balance - ?", amount))
			}
			if toAcctID != nil && *toAcctID != "" {
				db.Model(&models.Account{}).Where("id = ?", *toAcctID).
					Update("balance", gorm.Expr("balance + ?", amount))
			}
		}
	}

	applyBalance(req.Type, req.Amount, newAcctID, newToAcctID)

	return &TransactionResult{
		ID:          tx.ID,
		MemberID:    tx.MemberID,
		AccountID:   newAcctID,
		ToAccountID: newToAcctID,
		CategoryID:  newCatID,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		Notes:       req.Notes,
		Date:        req.Date,
	}, nil
}

// DeleteTransaction removes a transaction by ID (admin only).
func DeleteTransaction(db *gorm.DB, id string) error {
	r := db.Delete(&models.Transaction{}, "id = ?", id)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

// GetMonthlySummary returns income/expense totals grouped by category.
func GetMonthlySummary(db *gorm.DB, memberID string, year, month int) (*MonthlySummary, error) {
	var results []CategorySummary
	q := db.Table("transactions").
		Select(`categories.name as category_name, 
		        categories.color as category_color,
		        transactions.type,
		        SUM(transactions.amount) as total,
		        COUNT(*) as count`).
		Joins("LEFT JOIN categories ON categories.id = transactions.category_id").
		Where("EXTRACT(YEAR FROM transactions.date) = ? AND EXTRACT(MONTH FROM transactions.date) = ?", year, month).
		Where("transactions.type != 'transfer'") // exclude transfers from summary

	if memberID != "" {
		q = q.Where("transactions.member_id = ?", memberID)
	}

	if err := q.Group("categories.name, categories.color, transactions.type").
		Order("total DESC").
		Find(&results).Error; err != nil {
		return nil, err
	}

	summary := &MonthlySummary{}
	for _, r := range results {
		if r.Type == "income" {
			summary.TotalIncome += r.Total
		} else if r.Type == "expense" {
			summary.TotalExpense += r.Total
		}
		summary.CategoryBreak = append(summary.CategoryBreak, r)
	}

	return summary, nil
}

// MonthlySeriesPoint represents income/expense for a single month.
type MonthlySeriesPoint struct {
	Year    int   `json:"year"`
	Month   int   `json:"month"`
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
}

// GetMonthlySeries returns income/expense totals grouped by month for the last N months.
func GetMonthlySeries(db *gorm.DB, memberID string, months int) ([]MonthlySeriesPoint, error) {
	if months <= 0 || months > 24 {
		months = 6
	}

	var results []MonthlySeriesPoint

	q := db.Table("transactions").
		Select(`EXTRACT(YEAR FROM date) as year,
		        EXTRACT(MONTH FROM date) as month,
		        COALESCE(SUM(CASE WHEN type='income' THEN amount ELSE 0 END), 0) as income,
		        COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0) as expense`)

	if memberID != "" {
		q = q.Where("member_id = ?", memberID)
	}

	// Filter to last N months
	now := time.Now()
	startDate := now.AddDate(0, -months+1, 0)
	q = q.Where("date >= ? AND date <= ?", startDate.Format("2006-01-02"), now.Format("2006-01-02"))

	if err := q.Group("year, month").
		Order("year ASC, month ASC").
		Find(&results).Error; err != nil {
		return nil, err
	}

	// Fill in missing months with zeroes
	pointMap := make(map[string]*MonthlySeriesPoint)
	for i := range results {
		key := fmt.Sprintf("%d-%02d", results[i].Year, results[i].Month)
		pointMap[key] = &results[i]
	}

	var series []MonthlySeriesPoint
	for i := months - 1; i >= 0; i-- {
		d := now.AddDate(0, -i, 0)
		y, m := d.Year(), int(d.Month())
		key := fmt.Sprintf("%d-%02d", y, m)
		if p, ok := pointMap[key]; ok {
			series = append(series, *p)
		} else {
			series = append(series, MonthlySeriesPoint{Year: y, Month: m, Income: 0, Expense: 0})
		}
	}

	return series, nil
}

// ─── Accounts ─────────────────────────────────────────────────────────────────

// AccountResult for API/CLI responses.
type AccountResult struct {
	ID       string `json:"id"`
	MemberID string `json:"member_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Balance  int64  `json:"balance"`
}

// CreateAccount creates a new account for a member.
func CreateAccount(db *gorm.DB, memberID, name, acctType string, initialBalance int64) (*AccountResult, error) {
	if name == "" {
		return nil, fmt.Errorf("account name is required")
	}
	validTypes := map[string]bool{"cash": true, "bank": true, "e_wallet": true, "savings": true, "investment": true}
	if !validTypes[acctType] {
		return nil, fmt.Errorf("invalid account type: must be cash, bank, e_wallet, savings, or investment")
	}

	a := models.Account{
		MemberID: memberID,
		Name:     name,
		Type:     acctType,
		Balance:  initialBalance,
	}
	if err := db.Create(&a).Error; err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	return &AccountResult{
		ID:       a.ID,
		MemberID: a.MemberID,
		Name:     a.Name,
		Type:     a.Type,
		Balance:  a.Balance,
	}, nil
}

// ListAccounts returns accounts for a member (or all if memberID is empty).
func ListAccounts(db *gorm.DB, memberID string) ([]AccountResult, error) {
	q := db.Table("accounts").Select("id, member_id, name, type, balance")
	if memberID != "" {
		q = q.Where("member_id = ?", memberID)
	}
	var results []AccountResult
	if err := q.Order("name ASC").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// ─── Transfer ─────────────────────────────────────────────────────────────────

// TransferMoney moves money between two accounts of the same member.
func TransferMoney(db *gorm.DB, memberID string, fromAccountID, toAccountID string, amount int64, description, date string) (*TransactionResult, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if fromAccountID == "" || toAccountID == "" {
		return nil, fmt.Errorf("from_account and to_account are required")
	}
	if fromAccountID == toAccountID {
		return nil, fmt.Errorf("cannot transfer to the same account")
	}

	// Verify both accounts exist and belong to member
	var fromCount, toCount int64
	db.Table("accounts").Where("id = ? AND member_id = ?", fromAccountID, memberID).Count(&fromCount)
	db.Table("accounts").Where("id = ? AND member_id = ?", toAccountID, memberID).Count(&toCount)
	if fromCount == 0 || toCount == 0 {
		return nil, fmt.Errorf("one or both accounts not found")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Create transfer transaction with to_account_id
	tx := models.Transaction{
		MemberID:    memberID,
		AccountID:   &fromAccountID,
		ToAccountID: &toAccountID,
		Amount:      amount,
		Type:        "transfer",
		Description: description,
		Date:        date,
	}

	if err := db.Create(&tx).Error; err != nil {
		return nil, fmt.Errorf("create transfer: %w", err)
	}

	// Debit from source, credit to destination
	db.Model(&models.Account{}).Where("id = ?", fromAccountID).
		Update("balance", gorm.Expr("balance - ?", amount))
	db.Model(&models.Account{}).Where("id = ?", toAccountID).
		Update("balance", gorm.Expr("balance + ?", amount))

	return &TransactionResult{
		ID:          tx.ID,
		MemberID:    tx.MemberID,
		AccountID:   tx.AccountID,
		ToAccountID: tx.ToAccountID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Description: tx.Description,
		Date:        tx.Date,
	}, nil
}

// ─── Categories ───────────────────────────────────────────────────────────────

// CategoryResult for API responses.
type CategoryResult struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// ListCategories returns all categories, optionally filtered by type.
func ListCategories(db *gorm.DB, catType string) ([]CategoryResult, error) {
	q := db.Table("categories").Select("id, name, type, icon, color")
	if catType == "income" || catType == "expense" {
		q = q.Where("type = ?", catType)
	}
	var results []CategoryResult
	if err := q.Order("type ASC, name ASC").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// ─── Export ───────────────────────────────────────────────────────────────────

// ExportTransactionsCSV exports transactions as CSV bytes.
func ExportTransactionsCSV(db *gorm.DB, memberID, memberName string, month, year int) ([]byte, error) {
	f := TransactionFilter{
		MemberID: memberID,
		Month:    month,
		Year:     year,
		Limit:    99999,
	}
	result, err := ListTransactions(db, f)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Header
	w.Write([]string{"Date", "Type", "Category", "Amount", "Description", "Notes"})

	for _, r := range result.Transactions {
		date, _ := r["date"].(string)
		if date == "" {
			if t, ok := r["date"].(time.Time); ok {
				date = t.Format("2006-01-02")
			}
		}
		ttype, _ := r["type"].(string)
		catName, _ := r["category_name"].(string)
		amount, _ := r["amount"].(int64)
		desc, _ := r["description"].(string)
		notes, _ := r["notes"].(string)

		w.Write([]string{
			date,
			ttype,
			catName,
			fmt.Sprintf("%d", amount),
			desc,
			notes,
		})
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
