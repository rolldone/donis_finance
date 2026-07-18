package console

import (
	"fmt"
	"os"

	"go_framework/internal/db"
	"go_framework/plugins/donisfinance/services"

	"github.com/spf13/cobra"
)

// TransactionCommands returns CLI commands for transactions and accounts.
func TransactionCommands() []*cobra.Command {
	cmds := []*cobra.Command{
		addTxCmd(),
		listTxCmd(),
		editTxCmd(),
		deleteTxCmd(),
		summaryTxCmd(),
		exportTxCmd(),
		transferTxCmd(),
		createAccountCmd(),
		listAccountCmd(),
		updateAccountCmd(),
		adjustAccountCmd(),
		importTxCmd(),
	}
	return cmds
}

func addTxCmd() *cobra.Command {
	var member, category, account, txType, description, notes, date string
	var amount int64

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-add",
		Short: "Add a transaction for a member",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || amount <= 0 || txType == "" {
				fmt.Fprintln(os.Stderr, "Error: --member, --amount, and --type are required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			// Resolve member
			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			// Resolve category by name
			var catID *string
			if category != "" {
				var id string
				if err := gdb.Table("categories").Select("id").Where("name = ?", category).Scan(&id).Error; err == nil && id != "" {
					catID = &id
				}
			}

			// Resolve account by name
			var accID *string
			if account != "" {
				var id string
				if err := gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", account, memberID).Scan(&id).Error; err == nil && id != "" {
					accID = &id
				}
			}

			result, err := services.CreateTransaction(gdb, memberID, accID, catID, nil, amount, txType, description, notes, "", date)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Transaction added: Rp%d %s (%s)\n", result.Amount, result.Type, result.Date)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().Int64VarP(&amount, "amount", "a", 0, "Amount in Rupiah")
	cmd.Flags().StringVarP(&txType, "type", "t", "", "Transaction type (income/expense/transfer)")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Category name")
	cmd.Flags().StringVarP(&account, "account", "k", "", "Account name")
	cmd.Flags().StringVarP(&description, "desc", "d", "", "Description")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Long notes / catatan")
	cmd.Flags().StringVarP(&date, "date", "D", "", "Date (YYYY-MM-DD, default today)")
	return cmd
}

func listTxCmd() *cobra.Command {
	var member, txType string
	var month, year, limit int

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-list",
		Short: "List transactions for a member",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" {
				fmt.Fprintln(os.Stderr, "Error: --member is required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			// Resolve member
			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			f := services.TransactionFilter{
				MemberID: memberID,
				Month:    month,
				Year:     year,
				Type:     txType,
				Limit:    limit,
			}

			result, err := services.ListTransactions(gdb, f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if len(result.Transactions) == 0 {
				fmt.Println("No transactions found.")
				return
			}

			fmt.Printf("%-36s %-12s %-10s %-15s %-20s %s\n", "ID", "Date", "Type", "Amount", "Category", "Description")
			fmt.Println("──────────────────────────────────────────────────────────────────────────────────────────────")
			for _, t := range result.Transactions {
				amount := t["amount"].(int64)
				ttype := t["type"].(string)
				catName := ""
				if n, ok := t["category_name"]; ok {
					catName, _ = n.(string)
				}
				desc := ""
				if d, ok := t["description"]; ok {
					desc, _ = d.(string)
				}
				date := ""
				if d, ok := t["date"]; ok {
					date = fmt.Sprintf("%v", d)
				}
				id := ""
				if i, ok := t["id"]; ok {
					id = fmt.Sprintf("%v", i)
				}
				symbol := ""
				if ttype == "income" {
					symbol = "+"
				} else {
					symbol = "-"
				}
				fmt.Printf("%-36s %-12s %-10s %sRp%-12d %-20s %s\n", id, date, ttype, symbol, amount, catName, desc)
			}
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month filter (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year filter")
	cmd.Flags().StringVarP(&txType, "type", "t", "", "Filter by type (income/expense)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Max results")
	return cmd
}

func editTxCmd() *cobra.Command {
	var id, member, category, account, txType, description, date string
	var amount int64

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-edit",
		Short: "Edit a transaction (partial update)",
		Run: func(cmd *cobra.Command, args []string) {
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			// Load existing transaction
			var tx struct {
				ID          string
				MemberID    string
				AccountID   *string
				ToAccountID *string
				CategoryID  *string
				Amount      int64
				Type        string
				Description string
				Notes       string
				Date        string
			}
			if err := gdb.Table("transactions").Where("id = ?", id).Scan(&tx).Error; err != nil {
				fmt.Fprintf(os.Stderr, "Error: transaction not found: %v\n", err)
				os.Exit(1)
			}
			if tx.ID == "" {
				fmt.Fprintln(os.Stderr, "Error: transaction not found")
				os.Exit(1)
			}

			// Build update request starting from existing values
			req := struct {
				AccountID   string `json:"account_id"`
				ToAccountID string `json:"to_account_id"`
				CategoryID  string `json:"category_id"`
				Amount      int64  `json:"amount"`
				Type        string `json:"type"`
				Description string `json:"description"`
				Notes       string `json:"notes"`
				Date        string `json:"date"`
			}{
				Amount:      tx.Amount,
				Type:        tx.Type,
				Description: tx.Description,
				Notes:       tx.Notes,
				Date:        tx.Date,
			}
			if tx.AccountID != nil {
				req.AccountID = *tx.AccountID
			}
			if tx.ToAccountID != nil {
				req.ToAccountID = *tx.ToAccountID
			}
			if tx.CategoryID != nil {
				var catName string
				gdb.Table("categories").Select("name").Where("id = ?", *tx.CategoryID).Scan(&catName)
				category = catName // for display only
				req.CategoryID = *tx.CategoryID
			}

			// Override with user-provided flags
			if amount > 0 {
				req.Amount = amount
			}
			if txType != "" {
				req.Type = txType
			}
			if description != "" {
				req.Description = description
			}
			if date != "" {
				req.Date = date
			}
			if account != "" {
				var accID string
				if err := gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", account, tx.MemberID).Scan(&accID).Error; err != nil || accID == "" {
					fmt.Fprintf(os.Stderr, "Account %q not found for this member\n", account)
					os.Exit(1)
				}
				req.AccountID = accID
			}
			if category != "" {
				var catID string
				if err := gdb.Table("categories").Select("id").Where("name = ?", category).Scan(&catID).Error; err != nil || catID == "" {
					fmt.Fprintf(os.Stderr, "Category %q not found\n", category)
					os.Exit(1)
				}
				req.CategoryID = catID
			}

			// Call service update (role=admin so it bypasses member ownership check)
			result, err := services.UpdateTransaction(gdb, id, "", "admin", req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Transaction updated: Rp%d %s (%s)\n", result.Amount, result.Type, result.Date)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "Transaction ID (required)")
	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username (not used for edit, kept for consistency)")
	cmd.Flags().StringVarP(&category, "category", "c", "", "New category name")
	cmd.Flags().Int64VarP(&amount, "amount", "a", 0, "New amount in Rupiah")
	cmd.Flags().StringVarP(&txType, "type", "t", "", "New type (income/expense/transfer)")
	cmd.Flags().StringVarP(&description, "desc", "d", "", "New description")
	cmd.Flags().StringVarP(&account, "account", "k", "", "New account name")
	cmd.Flags().StringVarP(&date, "date", "D", "", "New date (YYYY-MM-DD)")
	return cmd
}

func deleteTxCmd() *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-delete",
		Short: "Delete a transaction by ID",
		Run: func(cmd *cobra.Command, args []string) {
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			if err := services.DeleteTransaction(gdb, id); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Transaction deleted")
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "Transaction ID")
	return cmd
}

func summaryTxCmd() *cobra.Command {
	var member string
	var month, year int

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-summary",
		Short: "Monthly income/expense summary",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || month <= 0 || year <= 0 {
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

			summary, err := services.GetMonthlySummary(gdb, memberID, year, month)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			monthName := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

			fmt.Printf("\n📊 Summary %s %d — %s\n", monthName[month], year, member)
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("💚 Income:  Rp%d\n", summary.TotalIncome)
			fmt.Printf("❤️ Expense: Rp%d\n", summary.TotalExpense)
			fmt.Printf("💰 Balance: Rp%d\n", summary.TotalIncome-summary.TotalExpense)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("%-25s %10s %5s\n", "Category", "Amount", "Count")
			fmt.Println("─────────────────────────────────────")
			for _, c := range summary.CategoryBreak {
				emoji := "❤️"
				if c.Type == "income" {
					emoji = "💚"
				}
				fmt.Printf("%s %-22s Rp%-8d %3dx\n", emoji, c.CategoryName, c.Total, c.Count)
			}
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	return cmd
}

// ─── Account commands ─────────────────────────────────────────────────────────

func createAccountCmd() *cobra.Command {
	var member, name, acctType string
	var balance int64

	cmd := &cobra.Command{
		Use:   "donisfinance:account-create",
		Short: "Create an account for a member",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || name == "" {
				fmt.Fprintln(os.Stderr, "Error: --member and --name are required")
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

			result, err := services.CreateAccount(gdb, memberID, name, acctType, balance)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Account created: %s (%s) — Rp%d\n", result.Name, result.Type, result.Balance)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Account name (e.g. 'Cash', 'BCA')")
	cmd.Flags().StringVarP(&acctType, "type", "t", "cash", "Account type (cash/bank/e_wallet/savings/investment)")
	cmd.Flags().Int64VarP(&balance, "balance", "b", 0, "Initial balance")
	return cmd
}

func listAccountCmd() *cobra.Command {
	var member string

	cmd := &cobra.Command{
		Use:   "donisfinance:account-list",
		Short: "List accounts for a member",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" {
				fmt.Fprintln(os.Stderr, "Error: --member is required")
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

			accounts, err := services.ListAccounts(gdb, memberID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if len(accounts) == 0 {
				fmt.Println("No accounts found.")
				return
			}

			fmt.Printf("%-36s %-20s %-15s %15s\n", "ID", "Name", "Type", "Balance")
			fmt.Println("──────────────────────────────────────────────────────────────────────────")
			for _, a := range accounts {
				fmt.Printf("%-36s %-20s %-15s Rp%-12d\n", a.ID, a.Name, a.Type, a.Balance)
			}
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	return cmd
}

func updateAccountCmd() *cobra.Command {
	var member, account, newName, newType, reason string
	var balance int64

	cmd := &cobra.Command{
		Use:   "donisfinance:account-update",
		Short: "Update account name/type/balance (all optional, at least 1 required)",
		Long: `Update an account's name, type, and/or balance for a member.

Examples:
  console donisfinance:account-update --member donny --account "BCA" --name "BCA - 7670339836"
  console donisfinance:account-update --member donny --account "BCA" --type bank
  console donisfinance:account-update --member donny --account "BCA" --balance 44800000 --reason "Reconcile"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if member == "" || account == "" {
				return fmt.Errorf("--member and --account are required")
			}

			hasName := cmd.Flags().Changed("name")
			hasType := cmd.Flags().Changed("type")
			hasBalance := cmd.Flags().Changed("balance")

			if !hasName && !hasType && !hasBalance {
				return fmt.Errorf("at least one of --name, --type, or --balance must be provided")
			}
			if hasBalance && reason == "" {
				return fmt.Errorf("--reason is required when using --balance")
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				return fmt.Errorf("DB error: %w", err)
			}

			// Resolve member
			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				return fmt.Errorf("member %q not found", member)
			}

			// Resolve account (scoped to member)
			var accountID string
			if err := gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", account, memberID).Scan(&accountID).Error; err != nil || accountID == "" {
				return fmt.Errorf("account %q not found for member %s", account, member)
			}

			// Prepare optional balance pointer
			var balPtr *int64
			if hasBalance {
				balPtr = &balance
			}

			result, err := services.UpdateAccount(gdb, accountID, memberID, newName, newType, balPtr, reason)
			if err != nil {
				return err
			}

			changes := ""
			if hasName {
				changes += fmt.Sprintf(" name=%q", result.Name)
			}
			if hasType {
				changes += fmt.Sprintf(" type=%s", result.Type)
			}
			if hasBalance {
				changes += fmt.Sprintf(" balance=Rp%d", result.Balance)
			}
			fmt.Printf("✅ Account updated:%s\n", changes)
			return nil
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username (required)")
	cmd.Flags().StringVarP(&account, "account", "a", "", "Current account name (required)")
	cmd.Flags().StringVarP(&newName, "name", "n", "", "New account name")
	cmd.Flags().StringVarP(&newType, "type", "t", "", "New account type (cash/bank/e_wallet/savings/investment)")
	cmd.Flags().Int64VarP(&balance, "balance", "b", 0, "New balance (set saldo langsung)")
	cmd.Flags().StringVarP(&reason, "reason", "r", "", "Reason for balance change (required if --balance used)")
	return cmd
}

func adjustAccountCmd() *cobra.Command {
	var member, account, reason string
	var balance int64

	cmd := &cobra.Command{
		Use:   "donisfinance:account-adjust",
		Short: "Manually adjust an account balance (with audit trail)",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || account == "" || reason == "" {
				fmt.Fprintln(os.Stderr, "Error: --member, --account, --balance, and --reason are required")
				os.Exit(1)
			}
			if balance < 0 {
				fmt.Fprintln(os.Stderr, "Error: --balance must be >= 0")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			// Resolve member
			var memberID string
			if err := gdb.Table("members").Select("id").Where("username = ?", member).Scan(&memberID).Error; err != nil || memberID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			// Resolve account (scoped to member)
			var accountID string
			if err := gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", account, memberID).Scan(&accountID).Error; err != nil || accountID == "" {
				fmt.Fprintf(os.Stderr, "Account %q not found for member %s\n", account, member)
				os.Exit(1)
			}

			result, err := services.AdjustBalance(gdb, accountID, balance, reason)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Balance adjusted: %s → Rp%d — %s\n", result.Name, result.Balance, reason)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&account, "account", "a", "", "Account name (e.g. 'Cash', 'BCA')")
	cmd.Flags().Int64VarP(&balance, "balance", "b", 0, "New balance value in Rupiah")
	cmd.Flags().StringVarP(&reason, "reason", "r", "", "Reason for adjustment (required, e.g. 'Reconcile with bank statement')")
	return cmd
}

func exportTxCmd() *cobra.Command {
	var member, file string
	var month, year int

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-export",
		Short: "Export transactions to CSV",
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

			var memberID, memberName string
			var m struct{ ID, Name string }
			if err := gdb.Table("members").Select("id, name").Where("username = ?", member).Scan(&m).Error; err != nil || m.ID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}
			memberID = m.ID
			memberName = m.Name

			csvData, err := services.ExportTransactionsCSV(gdb, memberID, memberName, month, year)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if file == "" {
				file = fmt.Sprintf("transactions_%s_%d_%02d.csv", member, year, month)
			}

			if err := os.WriteFile(file, csvData, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Exported %d bytes to %s\n", len(csvData), file)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Output file (default: auto-name)")
	return cmd
}

func transferTxCmd() *cobra.Command {
	var member, fromAccount, toAccount, description, date string
	var amount int64

	cmd := &cobra.Command{
		Use:   "donisfinance:tx-transfer",
		Short: "Transfer money between accounts",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || fromAccount == "" || toAccount == "" || amount <= 0 {
				fmt.Fprintln(os.Stderr, "Error: --member, --from, --to, and --amount are required")
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

			// Resolve account IDs by name
			var fromID, toID string
			gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", fromAccount, memberID).Scan(&fromID)
			gdb.Table("accounts").Select("id").Where("name = ? AND member_id = ?", toAccount, memberID).Scan(&toID)

			result, err := services.TransferMoney(gdb, memberID, fromID, toID, amount, description, date)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Transferred Rp%d: %s → %s (%s)\n", result.Amount, fromAccount, toAccount, result.Description)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&fromAccount, "from", "f", "", "Source account name")
	cmd.Flags().StringVarP(&toAccount, "to", "t", "", "Destination account name")
	cmd.Flags().Int64VarP(&amount, "amount", "a", 0, "Amount to transfer")
	cmd.Flags().StringVarP(&description, "desc", "d", "", "Description")
	cmd.Flags().StringVarP(&date, "date", "D", "", "Date (YYYY-MM-DD)")
	return cmd
}
