package services

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ─── DTO ──────────────────────────────────────────────────────────────────────

// AcctInfo for dashboard.
type AcctInfo struct {
	Name    string
	Type    string
	Balance int64
}

// BudgetInfo for dashboard.
type BudgetInfo struct {
	CategoryName string
	Amount       int64
	Spent        int64
	Percentage   int
}

// TopCatInfo for dashboard.
type TopCatInfo struct {
	CategoryName string
	Total        int64
}

// RecentTxInfo for dashboard.
type RecentTxInfo struct {
	Date        string
	Amount      int64
	Type        string
	Description string
	AccountName string
	ToAcctName  string
}

// DashboardData holds all data for the CLI dashboard.
type DashboardData struct {
	MemberName   string
	Month        int
	Year         int
	Accounts     []AcctInfo
	TotalBalance int64
	Income       int64
	Expense      int64
	Budgets      []BudgetInfo
	TopExpense   []TopCatInfo
	RecentTx     []RecentTxInfo
}

// GetDashboard collects all data in one query-heavy call.
func GetDashboard(db *gorm.DB, memberID, memberName string, month, year int) (*DashboardData, error) {
	d := &DashboardData{
		MemberName: memberName,
		Month:      month,
		Year:       year,
	}

	// 1. Accounts
	db.Table("accounts").Select("name, type, balance").Where("member_id = ?", memberID).Order("name ASC").Find(&d.Accounts)
	for _, a := range d.Accounts {
		d.TotalBalance += a.Balance
	}

	// 2. Monthly income/expense
	type MonthRow struct {
		Income  int64
		Expense int64
	}
	var row MonthRow
	db.Table("transactions").
		Select("COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense").
		Where("member_id = ? AND EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", memberID, year, month).
		Scan(&row)
	d.Income = row.Income
	d.Expense = row.Expense

	// 3. Budget status
	db.Raw(`
		SELECT COALESCE(c.name, '') as category_name, b.amount, COALESCE(SUM(t.amount), 0) as spent
		FROM budgets b
		LEFT JOIN categories c ON c.id = b.category_id
		LEFT JOIN transactions t ON t.category_id = b.category_id AND t.member_id = b.member_id
			AND t.type = 'expense' AND EXTRACT(YEAR FROM t.date) = b.year AND EXTRACT(MONTH FROM t.date) = b.month
		WHERE b.member_id = ? AND b.month = ? AND b.year = ?
		GROUP BY c.name, b.amount
		ORDER BY c.name`, memberID, month, year).Scan(&d.Budgets)
	for i := range d.Budgets {
		if d.Budgets[i].Amount > 0 {
			pct := int((d.Budgets[i].Spent * 100) / d.Budgets[i].Amount)
			if pct > 100 {
				pct = 100
			}
			d.Budgets[i].Percentage = pct
		}
	}

	// 4. Top expense categories
	db.Table("transactions t").
		Select("COALESCE(c.name, 'Tanpa Kategori') as category_name, SUM(t.amount) as total").
		Joins("LEFT JOIN categories c ON c.id = t.category_id").
		Where("t.member_id = ? AND t.type = 'expense' AND EXTRACT(YEAR FROM t.date) = ? AND EXTRACT(MONTH FROM t.date) = ?", memberID, year, month).
		Group("c.name").Order("total DESC").Limit(5).Find(&d.TopExpense)

	// 5. Recent transactions
	type RecentRow struct {
		Date        time.Time
		Amount      int64
		Type        string
		Description string
		AccountName string
		ToAcctName  string
	}
	var raws []RecentRow
	db.Table("transactions t").
		Select("t.date, t.amount, t.type, t.description, COALESCE(a.name, '') as account_name, COALESCE(a2.name, '') as to_account_name").
		Joins("LEFT JOIN accounts a ON a.id = t.account_id").
		Joins("LEFT JOIN accounts a2 ON a2.id = t.to_account_id").
		Where("t.member_id = ? AND EXTRACT(YEAR FROM t.date) = ? AND EXTRACT(MONTH FROM t.date) = ?", memberID, year, month).
		Order("t.date DESC, t.created_at DESC").Limit(5).Find(&raws)
	for _, r := range raws {
		d.RecentTx = append(d.RecentTx, RecentTxInfo{
			Date: r.Date.Format("2006-01-02"), Amount: r.Amount, Type: r.Type,
			Description: r.Description, AccountName: r.AccountName, ToAcctName: r.ToAcctName,
		})
	}

	return d, nil
}

// FormatDashboard renders the dashboard as a formatted string.
func FormatDashboard(d *DashboardData) string {
	var b strings.Builder
	monthName := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	b.WriteString(strings.Repeat("━", 50) + "\n")
	b.WriteString(fmt.Sprintf("  📊  DASHBOARD — %s (%s %d)\n", d.MemberName, monthName[d.Month], d.Year))
	b.WriteString(strings.Repeat("━", 50) + "\n\n")

	// Saldo
	b.WriteString("  💰 SALDO\n")
	for _, a := range d.Accounts {
		b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", a.Name, formatNum(a.Balance)))
	}
	b.WriteString(fmt.Sprintf("     %-20s ─────────────\n", ""))
	b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", "TOTAL", formatNum(d.TotalBalance)))

	// Monthly
	b.WriteString("\n  📈 BULAN INI\n")
	b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", "Income", formatNum(d.Income)))
	b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", "Expense", formatNum(d.Expense)))
	sisa := d.Income - d.Expense
	sisaStr := formatNum(sisa)
	if sisa < 0 {
		sisaStr = "-" + formatNum(-sisa)
	}
	b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", "Sisa", sisaStr))

	// Budget
	if len(d.Budgets) > 0 {
		b.WriteString("\n  ⚠️  BUDGET\n")
		for _, bg := range d.Budgets {
			bar := strings.Repeat("█", bg.Percentage/10) + strings.Repeat("░", 10-bg.Percentage/10)
			warn := ""
			if bg.Percentage >= 100 {
				warn = " 🔴 OVER"
			} else if bg.Percentage >= 80 {
				warn = " ⚠️"
			}
			b.WriteString(fmt.Sprintf("     %-20s %s/%s %3d%% %s%s\n",
				bg.CategoryName, formatNum(bg.Spent), formatNum(bg.Amount), bg.Percentage, bar, warn))
		}
	}

	// Top expense
	if len(d.TopExpense) > 0 {
		b.WriteString("\n  🔥 TOP EXPENSE\n")
		for _, te := range d.TopExpense {
			b.WriteString(fmt.Sprintf("     %-20s Rp%12s\n", te.CategoryName, formatNum(te.Total)))
		}
	}

	// Recent
	if len(d.RecentTx) > 0 {
		b.WriteString("\n  📝 TRANSAKSI TERAKHIR\n")
		for _, tx := range d.RecentTx {
			symbol := ""
			switch tx.Type {
			case "income":
				symbol = "+"
			case "expense":
				symbol = "-"
			case "transfer":
				symbol = "↔"
			}
			desc := tx.Description
			if desc == "" {
				if tx.Type == "transfer" {
					desc = fmt.Sprintf("Transfer: %s→%s", tx.AccountName, tx.ToAcctName)
				}
			}
			b.WriteString(fmt.Sprintf("     %s %sRp%-12d %s\n", tx.Date, symbol, tx.Amount, desc))
		}
	}

	b.WriteString("\n" + strings.Repeat("━", 50) + "\n")
	return b.String()
}

func formatNum(n int64) string {
	s := fmt.Sprintf("%d", n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}
