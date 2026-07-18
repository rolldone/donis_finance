package console

import (
	"fmt"
	"os"
	"strings"

	"go_framework/internal/db"
	"go_framework/plugins/donisfinance/services"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// getGormDB returns a GORM DB connection for console commands.
func getGormDB() (*gorm.DB, error) {
	return db.GetGormDB()
}

// UserCommands returns cobra commands for user management.
func UserCommands() []*cobra.Command {
	var (
		username string
		password string
		email    string
		name     string
		admin    string
	)

	createAdmin := &cobra.Command{
		Use:   "donisfinance:create-admin",
		Short: "Create a new admin user",
		Run: func(cmd *cobra.Command, args []string) {
			if username == "" || password == "" {
				fmt.Fprintln(os.Stderr, "Error: --username and --password are required")
				os.Exit(1)
			}

			gdb, err := getGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}

			result, err := services.CreateAdmin(gdb, username, password, email)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating admin: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Admin created: id=%s username=%s\n", result.ID, result.Username)
		},
	}
	createAdmin.Flags().StringVarP(&username, "username", "u", "", "Admin username")
	createAdmin.Flags().StringVarP(&password, "password", "p", "", "Admin password")
	createAdmin.Flags().StringVarP(&email, "email", "e", "", "Admin email (optional)")

	createMember := &cobra.Command{
		Use:   "donisfinance:create-member",
		Short: "Create a new member under an admin",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || username == "" || password == "" {
				fmt.Fprintln(os.Stderr, "Error: --name, --username, and --password are required")
				os.Exit(1)
			}
			if admin == "" {
				fmt.Fprintln(os.Stderr, "Error: --admin is required (admin username who owns this member)")
				os.Exit(1)
			}

			gdb, err := getGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}

			// Look up admin by username
			var adminID string
			if err := gdb.Table("admins").Select("id").Where("username = ?", admin).Scan(&adminID).Error; err != nil {
				fmt.Fprintf(os.Stderr, "Admin lookup failed: %v\n", err)
				os.Exit(1)
			}
			if adminID == "" {
				fmt.Fprintf(os.Stderr, "Error: admin not found with username: %s\n", admin)
				os.Exit(1)
			}

			result, err := services.CreateMember(gdb, adminID, name, username, password)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating member: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Member created: id=%s name=%s username=%s\n", result.ID, result.Name, result.Username)
		},
	}
	createMember.Flags().StringVarP(&name, "name", "n", "", "Member display name")
	createMember.Flags().StringVarP(&username, "username", "u", "", "Member username")
	createMember.Flags().StringVarP(&password, "password", "p", "", "Member password")
	createMember.Flags().StringVarP(&admin, "admin", "a", "", "Admin username who owns this member")

	listAdmins := &cobra.Command{
		Use:   "donisfinance:list-admins",
		Short: "List all admin users",
		Run: func(cmd *cobra.Command, args []string) {
			gdb, err := getGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}

			type AdminRow struct {
				ID       string
				Username string
				Email    string
			}
			var admins []AdminRow
			gdb.Table("admins").Select("id, username, email").Order("created_at DESC").Find(&admins)

			if len(admins) == 0 {
				fmt.Println("No admins found.")
				return
			}
			fmt.Printf("%-36s %-20s %-30s\n", "ID", "Username", "Email")
			fmt.Println(strings.Repeat("-", 86))
			for _, a := range admins {
				fmt.Printf("%-36s %-20s %-30s\n", a.ID, a.Username, a.Email)
			}
		},
	}

	listMembers := &cobra.Command{
		Use:   "donisfinance:list-members",
		Short: "List all members",
		Run: func(cmd *cobra.Command, args []string) {
			gdb, err := getGormDB()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
				os.Exit(1)
			}

			type MemberRow struct {
				ID       string
				Name     string
				Username string
			}
			var members []MemberRow
			gdb.Table("members").Select("id, name, username").Order("created_at DESC").Find(&members)

			if len(members) == 0 {
				fmt.Println("No members found.")
				return
			}
			fmt.Printf("%-36s %-20s %-20s\n", "ID", "Name", "Username")
			fmt.Println(strings.Repeat("-", 76))
			for _, m := range members {
				fmt.Printf("%-36s %-20s %-20s\n", m.ID, m.Name, m.Username)
			}
		},
	}

	return []*cobra.Command{createAdmin, createMember, listAdmins, listMembers}
}
