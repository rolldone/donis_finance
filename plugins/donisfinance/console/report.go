package console

import (
	"fmt"
	"os"
	"time"

	"github.com/rolldone/donisgo/internal/db"
	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/spf13/cobra"
)

// ReportCommands returns CLI commands for reports.
func ReportCommands() []*cobra.Command {
	return []*cobra.Command{
		sendReportCmd(),
		sendBulkReportsCmd(),
	}
}

func sendReportCmd() *cobra.Command {
	var member, toEmail string
	var month, year int

	cmd := &cobra.Command{
		Use:   "donisfinance:send-report",
		Short: "Send monthly report via email to one member",
		Run: func(cmd *cobra.Command, args []string) {
			if member == "" || toEmail == "" {
				fmt.Fprintln(os.Stderr, "Error: --member and --to are required")
				os.Exit(1)
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			var m struct{ ID, Name string }
			if err := gdb.Table("members").Select("id, name").Where("username = ?", member).Scan(&m).Error; err != nil || m.ID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			if err := services.SendMonthlyReport(gdb, m.ID, m.Name, toEmail, month, year); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending report: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Report sent to %s\n", toEmail)
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().StringVarP(&toEmail, "to", "t", "", "Recipient email")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026)")
	return cmd
}

func sendBulkReportsCmd() *cobra.Command {
	var month, year int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "donisfinance:send-bulk-reports",
		Short: "Send monthly report to ALL members with email",
		Long: `Iterates all members that have an email address and sends them
the monthly financial report. Defaults to the previous month.`,
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			if month == 0 {
				month = int(now.Month()) - 1
				if month == 0 {
					month = 12
					year = now.Year() - 1
				}
			}
			if year == 0 {
				year = now.Year()
			}

			gdb, err := db.GetGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
				os.Exit(1)
			}

			type memberRow struct {
				ID    string
				Name  string
				Email string
			}
			var members []memberRow
			if err := gdb.Table("members").
				Select("id, name, email").
				Where("email IS NOT NULL AND email != ''").
				Find(&members).Error; err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching members: %v\n", err)
				os.Exit(1)
			}

			if len(members) == 0 {
				fmt.Println("No members with email found.")
				return
			}

			fmt.Printf("📬 Sending %d/%d report to %d members\n", month, year, len(members))

			ok := 0
			fail := 0
			for _, m := range members {
				if dryRun {
					fmt.Printf("  [DRY-RUN] → %s <%s>\n", m.Name, m.Email)
					ok++
					continue
				}
				if err := services.SendMonthlyReport(gdb, m.ID, m.Name, m.Email, month, year); err != nil {
					fmt.Fprintf(os.Stderr, "  ❌ %s <%s>: %v\n", m.Name, m.Email, err)
					fail++
				} else {
					fmt.Printf("  ✅ %s <%s>\n", m.Name, m.Email)
					ok++
				}
			}

			fmt.Printf("\nDone: %d sent, %d failed\n", ok, fail)
		},
	}

	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12), defaults to previous month")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (e.g. 2026), defaults to current")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Only list members, don't send")
	return cmd
}
