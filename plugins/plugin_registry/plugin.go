// Package plugin_registry menyediakan mekanisme komunikasi antar plugin
// melalui interface contract. Plugin mendaftarkan implementasinya, plugin
// lain memanggil lewat delegation functions — tanpa saling import.
package plugin_registry

import (
	"context"

	"gorm.io/gorm"
)

// ─── Transaction / Catatan Finansial ──────────────────────────────────────────

// TransactionProvider allows plugins to create financial transaction records.
type TransactionProvider interface {
	RecordTransaction(ctx context.Context, db *gorm.DB, memberID string, amount int64, category string, description string) error
}

var transactionProvider TransactionProvider

// RegisterTransactionProvider registers a transaction provider implementation.
func RegisterTransactionProvider(p TransactionProvider) {
	transactionProvider = p
}

// RecordTransaction delegates to the registered provider if available.
func RecordTransaction(ctx context.Context, db *gorm.DB, memberID string, amount int64, category string, description string) error {
	if transactionProvider == nil {
		return nil
	}
	return transactionProvider.RecordTransaction(ctx, db, memberID, amount, category, description)
}

// ─── Report / Laporan ────────────────────────────────────────────────────────

// ReportProvider allows plugins to generate financial reports.
type ReportProvider interface {
	GenerateMonthlyReport(ctx context.Context, db *gorm.DB, memberID string, year, month int) (interface{}, error)
}

var reportProvider ReportProvider

// RegisterReportProvider registers a report provider implementation.
func RegisterReportProvider(p ReportProvider) {
	reportProvider = p
}

// GenerateMonthlyReport delegates to the registered provider if available.
func GenerateMonthlyReport(ctx context.Context, db *gorm.DB, memberID string, year, month int) (interface{}, error) {
	if reportProvider == nil {
		return nil, nil
	}
	return reportProvider.GenerateMonthlyReport(ctx, db, memberID, year, month)
}

// ─── Budget / Anggaran ───────────────────────────────────────────────────────

// BudgetProvider allows plugins to check and enforce budget limits.
type BudgetProvider interface {
	CheckBudget(ctx context.Context, db *gorm.DB, memberID string, category string, amount int64) (bool, error)
}

var budgetProvider BudgetProvider

// RegisterBudgetProvider registers a budget provider implementation.
func RegisterBudgetProvider(p BudgetProvider) {
	budgetProvider = p
}

// CheckBudget delegates to the registered provider if available.
func CheckBudget(ctx context.Context, db *gorm.DB, memberID string, category string, amount int64) (bool, error) {
	if budgetProvider == nil {
		return true, nil
	}
	return budgetProvider.CheckBudget(ctx, db, memberID, category, amount)
}

// ─── Export / CSV ────────────────────────────────────────────────────────────

// ExportProvider allows plugins to export financial data.
type ExportProvider interface {
	ExportTransactions(ctx context.Context, db *gorm.DB, memberID string, format string) ([]byte, error)
}

var exportProvider ExportProvider

// RegisterExportProvider registers an export provider implementation.
func RegisterExportProvider(p ExportProvider) {
	exportProvider = p
}

// ExportTransactions delegates to the registered provider if available.
func ExportTransactions(ctx context.Context, db *gorm.DB, memberID string, format string) ([]byte, error) {
	if exportProvider == nil {
		return nil, nil
	}
	return exportProvider.ExportTransactions(ctx, db, memberID, format)
}

// ─── Notification / Pemberitahuan ────────────────────────────────────────────

// NotificationProvider allows plugins to send notifications (Discord, etc.).
type NotificationProvider interface {
	SendNotification(ctx context.Context, db *gorm.DB, memberID string, title string, message string) error
}

var notificationProvider NotificationProvider

// RegisterNotificationProvider registers a notification provider implementation.
func RegisterNotificationProvider(p NotificationProvider) {
	notificationProvider = p
}

// SendNotification delegates to the registered provider if available.
func SendNotification(ctx context.Context, db *gorm.DB, memberID string, title string, message string) error {
	if notificationProvider == nil {
		return nil
	}
	return notificationProvider.SendNotification(ctx, db, memberID, title, message)
}
