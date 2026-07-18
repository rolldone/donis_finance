package donisfinance

import (
	"fmt"
	"log"
	"os"

	"go_framework/internal/plugins"
	"go_framework/internal/storage"
	pluginconsole "go_framework/plugins/donisfinance/console"
	pluginhandlers "go_framework/plugins/donisfinance/handlers"
	pluginmiddleware "go_framework/plugins/donisfinance/middleware"
	"go_framework/plugins/donisfinance/models"
	"go_framework/plugins/donisfinance/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Plugin Donisfinance provides user authentication and member management.
type Plugin struct {
	db    *gorm.DB
	store storage.Store
}

// New returns a new plugin instance.
func New() plugins.Plugin { return &Plugin{} }

func (p *Plugin) ID() string { return "donisfinance" }

// RegisterServices stores DB reference and runs auto-migration.
func (p *Plugin) RegisterServices(deps plugins.ServiceDeps) error {
	p.db = deps.DB
	p.store = deps.Store
	registerEventHandlers()

	// Auto-migrate our models
	if err := deps.DB.AutoMigrate(
		&models.Admin{},
		&models.Member{},
		&models.Category{},
		&models.Account{},
		&models.Transaction{},
		&models.Budget{},
		&models.Setting{},
		&models.BalanceAdjustment{},
	); err != nil {
		return err
	}
	log.Println("[donisfinance] tables migrated: admins, members, categories, accounts, transactions, budgets, settings, balance_adjustments")
	return nil
}

func (p *Plugin) RegisterMiddleware() []plugins.MiddlewareDescriptor { return nil }

func (p *Plugin) RegisterRoutes(router *gin.Engine, admin *gin.RouterGroup, api *gin.RouterGroup) error {
	// SPA serving is handled by core bootstrap (registerStoreRoutes)

	authH := pluginhandlers.NewAuthHandler(p.db)
	memberH := pluginhandlers.NewMemberHandler(p.db)
	txH := pluginhandlers.NewTransactionHandler(p.db, p.store)
	budgetH := pluginhandlers.NewBudgetHandler(p.db)
	catH := pluginhandlers.NewCategoryHandler(p.db)
	profileH := pluginhandlers.NewProfileHandler(p.db)
	settingsH := pluginhandlers.NewSettingsHandler(p.db)

	// Serve uploaded files
	if localStore, ok := p.store.(*storage.LocalStore); ok {
		root := localStore.GetRoot()
		if _, err := os.Stat(root); err == nil {
			router.GET("/uploads/*filepath", func(c *gin.Context) {
				c.File(root + c.Param("filepath"))
			})
			log.Printf("[donisfinance] serving static files from %s at /uploads", root)
		}
	}

	// --- Public (no JWT) ---
	api.POST("/admin/login", authH.AdminLogin)
	api.POST("/member/login", authH.MemberLogin)
	api.POST("/member/auth/register", authH.Register)
	api.POST("/member/auth/forgot-password", authH.ForgotPassword)
	api.POST("/member/auth/reset-password", authH.ResetPassword)
	api.GET("/admin/health", pluginhandlers.HealthHandler)

	// --- Admin-only routes (JWT required) ---
	adm := api.Group("/admin")
	adm.Use(pluginmiddleware.JWTAuth("admin"))
	{
		// Profile & password
		adm.GET("/profile", profileH.GetAdminProfile)
		adm.PUT("/profile", profileH.UpdateAdminProfile)
		adm.PUT("/password", profileH.ChangeAdminPassword)

		// Settings
		adm.GET("/settings/smtp", settingsH.GetSMTPConfig)
		adm.PUT("/settings/smtp", settingsH.SaveSMTPConfig)

		// Member management
		adm.GET("/members", memberH.ListMembers)
		adm.POST("/members", memberH.CreateMember)
		adm.PUT("/members/:id", memberH.UpdateMember)
		adm.DELETE("/members/:id", memberH.DeleteMember)
		adm.PATCH("/members/:id/approve", memberH.ApproveMember)
		adm.PATCH("/members/:id/reject", memberH.RejectMember)

		// Categories
		adm.GET("/categories", txH.ListCategories)
		adm.POST("/categories", catH.CreateCategory)
		adm.PUT("/categories/:id", catH.UpdateCategory)
		adm.DELETE("/categories/:id", catH.DeleteCategory)

		// Accounts
		adm.GET("/accounts", txH.ListAccounts)

		// Transactions
		adm.GET("/transactions", txH.ListTransactions)
		adm.GET("/transactions/summary", txH.GetSummary)
		adm.GET("/transactions/monthly", txH.GetMonthlySeries)
		adm.GET("/transactions/export", txH.ExportCSV)
		adm.PUT("/transactions/:id", txH.UpdateTransaction)
		adm.POST("/transactions/:id/attachment", txH.UploadAttachment)
		adm.GET("/transactions/:id/attachment", txH.GetAttachment)
		adm.DELETE("/transactions/:id", txH.DeleteTransaction)

		// Budgets
		adm.POST("/budgets", budgetH.SetBudget)
		adm.GET("/budgets/status", budgetH.GetBudgetStatus)
		adm.DELETE("/budgets/:id", budgetH.DeleteBudget)
	}

	// --- Member-only routes (JWT required) ---
	mem := api.Group("/member")
	mem.Use(pluginmiddleware.JWTAuth("member"))
	{
		mem.GET("/profile", profileH.GetMemberProfile)
		mem.PUT("/profile", profileH.UpdateMemberProfile)
		mem.PUT("/password", profileH.ChangeMemberPassword)

		mem.GET("/categories", txH.ListCategories)
		mem.GET("/accounts", txH.ListAccounts)
		mem.POST("/accounts", txH.CreateAccount)
		mem.PUT("/accounts/:id", txH.UpdateAccount)
		mem.DELETE("/accounts/:id", txH.DeleteAccount)
		mem.GET("/transactions", txH.ListTransactions)
		mem.GET("/transactions/summary", txH.GetSummary)
		mem.GET("/transactions/monthly", txH.GetMonthlySeries)
		mem.POST("/transactions", txH.CreateTransaction)
		mem.PUT("/transactions/:id", txH.UpdateTransaction)
		mem.DELETE("/transactions/:id", txH.DeleteTransaction)
		mem.POST("/transactions/:id/attachment", txH.UploadAttachment)
		mem.GET("/transactions/:id/attachment", txH.GetAttachment)
		mem.POST("/budgets", budgetH.SetBudget)
		mem.GET("/budgets/status", budgetH.GetBudgetStatus)
	}

	_ = admin // admin group not used directly, all routes go through /api
	return nil
}

func (p *Plugin) Seed() error {
	if p.db == nil {
		return nil
	}

	// Seed default admin if none exists
	var adminCount int64
	p.db.Table("admins").Count(&adminCount)
	if adminCount == 0 {
		_, err := services.CreateAdmin(p.db, "admin", "admin123", "admin@donis.finance")
		if err != nil {
			return err
		}
		log.Println("[donisfinance] seeded default admin: admin / admin123")
	}

	// Seed default categories if none exists
	var catCount int64
	if err := p.db.Table("categories").Count(&catCount).Error; err != nil {
		return fmt.Errorf("categories count: %w", err)
	}
	if catCount > 0 {
		return nil
	}

	categories := []map[string]interface{}{
		// Income
		{"name": "Gaji", "type": "income", "icon": "briefcase", "color": "#10b981"},
		{"name": "Bonus", "type": "income", "icon": "gift", "color": "#14b8a6"},
		{"name": "Bisnis", "type": "income", "icon": "building", "color": "#06b6d4"},
		{"name": "Investasi", "type": "income", "icon": "trending-up", "color": "#8b5cf6"},
		{"name": "Lainnya", "type": "income", "icon": "plus-circle", "color": "#6b7280"},
		// Expense
		{"name": "Makanan", "type": "expense", "icon": "shopping-cart", "color": "#ef4444"},
		{"name": "Transport", "type": "expense", "icon": "truck", "color": "#f97316"},
		{"name": "Belanja", "type": "expense", "icon": "shopping-bag", "color": "#f59e0b"},
		{"name": "Tagihan", "type": "expense", "icon": "file-text", "color": "#eab308"},
		{"name": "Hiburan", "type": "expense", "icon": "music", "color": "#22c55e"},
		{"name": "Kesehatan", "type": "expense", "icon": "heart", "color": "#ec4899"},
		{"name": "Pendidikan", "type": "expense", "icon": "book", "color": "#6366f1"},
		{"name": "Rumah Tangga", "type": "expense", "icon": "home", "color": "#a855f7"},
		{"name": "Tabungan", "type": "expense", "icon": "piggy-bank", "color": "#06b6d4"},
		{"name": "Darurat", "type": "expense", "icon": "alert-triangle", "color": "#dc2626"},
	}

	for _, c := range categories {
		if err := p.db.Table("categories").Create(c).Error; err != nil {
			return err
		}
	}
	log.Printf("[donisfinance] seeded %d default categories", len(categories))
	return nil
}

func (p *Plugin) ConsoleCommands() []*cobra.Command {
	cmds := pluginconsole.UserCommands()
	cmds = append(cmds, pluginconsole.TransactionCommands()...)
	cmds = append(cmds, pluginconsole.BudgetCommands()...)
	cmds = append(cmds, pluginconsole.DashboardCommands()...)
	cmds = append(cmds, pluginconsole.ReportCommands()...)
	return cmds
}
