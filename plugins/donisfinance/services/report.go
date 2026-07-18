package services

import (
	"fmt"

	"go_framework/internal/mail"

	"gorm.io/gorm"
)

// ReportMailable implements mail.Mailable for monthly reports.
type ReportMailable struct {
	subject string
	tplBase string
	data    map[string]interface{}
	from    string
}

func (m *ReportMailable) Subject() string                     { return m.subject }
func (m *ReportMailable) TemplateBase() string                 { return m.tplBase }
func (m *ReportMailable) Data() map[string]interface{}         { return m.data }
func (m *ReportMailable) From() (string, string)               { return m.from, "" }

// SendMonthlyReport generates and sends a monthly financial report via email.
func SendMonthlyReport(db *gorm.DB, memberID, memberName, toEmail string, month, year int) error {
	data, err := GetDashboard(db, memberID, memberName, month, year)
	if err != nil {
		return fmt.Errorf("get dashboard: %w", err)
	}

	report := FormatDashboard(data)
	monthName := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	subject := fmt.Sprintf("Laporan Keuangan %s — %s %d", memberName, monthName[month], year)

	mailer := mail.NewMailer()
	return mailer.Send(toEmail, &ReportMailable{
		subject: subject,
		tplBase: "plugins/donisfinance/templates/email/report",
		data: map[string]interface{}{
			"Name":      memberName,
			"MonthName": monthName[month],
			"Year":      year,
			"Dashboard": report,
		},
	})
}
