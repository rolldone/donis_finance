package console

import (
	"fmt"
	"os"

	"github.com/rolldone/donisgo/internal/db"
	"github.com/rolldone/donisgo/plugins/donisfinance/services"

	"github.com/spf13/cobra"
)

// DashboardCommands returns CLI commands for the dashboard.
func DashboardCommands() []*cobra.Command {
	return []*cobra.Command{
		dashboardCmd(),
	}
}

func dashboardCmd() *cobra.Command {
	var member string
	var month, year int

	cmd := &cobra.Command{
		Use:   "donisfinance:dashboard",
		Short: "Show financial dashboard (saldo, budget, rekap)",
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

			var m struct{ ID, Name string }
			if err := gdb.Table("members").Select("id, name").Where("username = ?", member).Scan(&m).Error; err != nil || m.ID == "" {
				fmt.Fprintf(os.Stderr, "Member %q not found\n", member)
				os.Exit(1)
			}

			data, err := services.GetDashboard(gdb, m.ID, m.Name, month, year)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Print(services.FormatDashboard(data))
		},
	}

	cmd.Flags().StringVarP(&member, "member", "m", "", "Member username")
	cmd.Flags().IntVarP(&month, "month", "M", 0, "Month (1-12, default: current)")
	cmd.Flags().IntVarP(&year, "year", "Y", 0, "Year (default: current)")
	return cmd
}
