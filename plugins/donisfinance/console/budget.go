package console

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rolldone/donisgo/internal/db"
	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/spf13/cobra"
)

// BudgetCommands returns CLI commands for budget management.
func BudgetCommands() []*cobra.Command {
	return []*cobra.Command{
		setBudgetCmd(),
		statusBudgetCmd(),
		checkBudgetCmd(),
	}
}

func setBudgetCmd() *cobra.Command {
	var member, category string
	var month, year int
	var amount int64

	cmd := &cobra.Command{
		Use:   "donisfinance:budget-set",
		Short: "Set a monthly budget limit for a category",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || month == 0 || year == 0 || amount <= 0 {
				fmt.Fprintln(os.Stderr, "Error: --member, --month, --year, and --amount are required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			var catID *string
			if category != "" {
				var id string
				if err := gdb.Table("categories").Select("id").Where("name = ?", category).Scan(&id).Error; err == nil && id != "" {
					catID = &id
				}
			}

			result, err := services.SetBudget(gdb, memberID, catID, month, year, amount)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Budget set: %s %d/%d — Rp%d\n", category, month, year, result.Amount)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Category name")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	cmd.Flags().Int64VarP(&amount, "amount", "a", 0, "Budget limit in Rupiah")
	return cmd
}

func statusBudgetCmd() *cobra.Command {
	var member string
	var month, year int

	cmd := &cobra.Command{
		Use:   "donisfinance:budget-status",
		Short: "Show budget vs actual spending",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || month == 0 || year == 0 {
				fmt.Fprintln(os.Stderr, "Error: --member, --month, and --year are required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			status, err := services.GetBudgetStatus(gdb, memberID, month, year)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if len(status) == 0 {
				fmt.Println("No budgets set for this period.")
				return
			}

			monthName := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
			fmt.Printf("\n📊 Budget %s %d — %s\n", monthName[month], year, member)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("%-22s %12s %12s %12s %8s\n", "Category", "Budget", "Spent", "Remaining", "Used")
			fmt.Println("──────────────────────────────────────────────────────────────")
			for _, b := range status {
				bar := ""
				for i := 0; i < b.Percentage/10; i++ {
					bar += "█"
				}
				for i := b.Percentage / 10; i < 10; i++ {
					bar += "░"
				}
				sisa := b.Amount - b.Spent
				warn := ""
				if sisa < 0 {
					warn = " 🔴 OVER!"
				} else if sisa < int64(float64(b.Amount)*0.1) {
					warn = " ⚠️"
				}
				fmt.Printf("%-22s Rp%-9d Rp%-9d Rp%-9d %3d%% %s%s\n",
					b.CategoryName, b.Amount, b.Spent, sisa, b.Percentage, bar, warn)
			}
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	return cmd
}

func checkBudgetCmd() *cobra.Command {
	var member, category string
	var month, year int
	var amount int64

	cmd := &cobra.Command{
		Use:   "donisfinance:budget-check",
		Short: "Preview if a transaction fits the budget",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || category == "" || month == 0 || year == 0 || amount <= 0 {
				fmt.Fprintln(os.Stderr, "Error: all flags required (--member, --category, --month, --year, --amount)")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			var categoryID string
			if err := gdb.Table("categories").Select("id").Where("name = ?", category).Scan(&categoryID).Error; err != nil || categoryID == "" {
				fmt.Fprintf(os.Stderr, "Category %q not found\n", category)
				os.Exit(1)
			}

			budgetID, remaining, isOver, err := services.CheckBudgetSpending(gdb, memberID, categoryID, month, year, amount)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if budgetID == "" {
				fmt.Println("ℹ️ No budget set for this category/period.")
				return
			}

			if isOver {
				fmt.Printf("🔴 OVER BUDGET by Rp%d!\n", -remaining)
			} else {
				fmt.Printf("✅ Sisa budget: Rp%d\n", remaining)
			}
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Category name")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	cmd.Flags().Int64VarP(&amount, "amount", "a", 0, "Amount to check")
	return cmd
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseIntOrExit(s, label string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s must be a number\n", label)
		os.Exit(1)
	}
	return v
}
