package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go_framework/internal/auth"
	"go_framework/internal/mail"
	"go_framework/internal/uuid"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthSession stores a logged-in session (admin or member).
type AuthSession struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"` // "admin" or "member"
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AdminResult is returned after creating an admin.
type AdminResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

// MemberResult is returned after creating a member.
type MemberResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// HashPassword returns a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword compares a bcrypt hash with a plaintext password.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// LoginAdmin authenticates an admin by username/password and returns a session.
func LoginAdmin(db *gorm.DB, username, password string) (*AuthSession, error) {
	type adminRow struct {
		ID       string
		Username string
		Password string
	}
	var row adminRow
	if err := db.Table("admins").Where("username = ?", username).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}
	if !VerifyPassword(row.Password, password) {
		return nil, errors.New("invalid username or password")
	}

	token, err := auth.SignAccessToken(row.ID)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(time.Duration(auth.AccessExpirySeconds()) * time.Second)

	return &AuthSession{
		UserID:    row.ID,
		Username:  row.Username,
		Role:      "admin",
		Token:     token,
		ExpiresAt: exp,
	}, nil
}

// LoginMember authenticates a member by username/password and returns a session.
func LoginMember(db *gorm.DB, username, password string) (*AuthSession, error) {
	type memberRow struct {
		ID       string
		Username string
		Password string
		Status   string
	}
	var row memberRow
	if err := db.Table("members").Where("username = ?", username).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}
	if !VerifyPassword(row.Password, password) {
		return nil, errors.New("invalid username or password")
	}

	if row.Status != "active" {
		return nil, errors.New("akun belum diaktivasi oleh admin")
	}

	token, err := auth.SignAccessToken(row.ID)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(time.Duration(auth.AccessExpirySeconds()) * time.Second)

	return &AuthSession{
		UserID:    row.ID,
		Username:  row.Username,
		Role:      "member",
		Token:     token,
		ExpiresAt: exp,
	}, nil
}

// CreateAdmin creates a new admin user with a hashed password.
func CreateAdmin(db *gorm.DB, username, password, email string) (*AdminResult, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	admin := map[string]interface{}{
		"id":       id,
		"username": username,
		"password": hash,
		"email":    email,
	}

	if err := db.Table("admins").Create(admin).Error; err != nil {
		return nil, err
	}

	return &AdminResult{ID: id, Username: username, Email: email}, nil
}

// CreateMember creates a new member under the given admin.
func CreateMember(db *gorm.DB, adminID, name, username, password string) (*MemberResult, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	member := map[string]interface{}{
		"id":       id,
		"admin_id": adminID,
		"name":     name,
		"username": username,
		"password": hash,
	}

	if err := db.Table("members").Create(member).Error; err != nil {
		return nil, err
	}

	return &MemberResult{ID: id, Name: name, Username: username}, nil
}

// UpdateMember updates a member's name and/or username.
func UpdateMember(db *gorm.DB, adminID, memberID, name, username string) (*MemberResult, error) {
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if username != "" {
		updates["username"] = username
	}
	if len(updates) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	r := db.Table("members").Where("id = ? AND admin_id = ?", memberID, adminID).Updates(updates)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, fmt.Errorf("member not found")
	}

	var result MemberResult
	if err := db.Table("members").Select("id, name, username").Where("id = ?", memberID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── Public Registration ──────────────────────────────────────────────────────

// RegisterMember creates a new member with status "pending", assigned to the first admin.
func RegisterMember(db *gorm.DB, name, email, password string) (*MemberResult, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("name, email, and password are required")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	// Check if email is already taken
	var count int64
	db.Table("members").Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	// Find first admin to assign to
	var admin struct {
		ID string
	}
	if err := db.Table("admins").Order("created_at ASC").First(&admin).Error; err != nil {
		return nil, fmt.Errorf("no admin available, please contact administrator")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	// Generate a username from email prefix
	username := email[:len(email)-len(email[strings.LastIndex(email, "@"):])]

	member := map[string]interface{}{
		"id":       id,
		"admin_id": admin.ID,
		"name":     name,
		"username": username,
		"email":    email,
		"password": hash,
		"status":   "pending",
	}

	if err := db.Table("members").Create(member).Error; err != nil {
		return nil, err
	}

	return &MemberResult{ID: id, Name: name, Username: username}, nil
}

// ApproveMember sets a member's status to "active".
func ApproveMember(db *gorm.DB, adminID, memberID string) error {
	r := db.Table("members").Where("id = ? AND admin_id = ?", memberID, adminID).Update("status", "active")
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

// RejectMember sets a member's status to "rejected".
func RejectMember(db *gorm.DB, adminID, memberID string) error {
	r := db.Table("members").Where("id = ? AND admin_id = ?", memberID, adminID).Update("status", "rejected")
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

// ─── Forgot / Reset Password ──────────────────────────────────────────────────

// ForgotPasswordResult contains info needed to send the reset email.
type ForgotPasswordResult struct {
	MemberID    string
	MemberName  string
	MemberEmail string
	ResetToken  string
}

// ForgotPassword finds a member by email, generates a reset token, and returns the result.
func ForgotPassword(db *gorm.DB, email string) (*ForgotPasswordResult, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	var member struct {
		ID    string
		Name  string
		Email string
	}
	if err := db.Table("members").Select("id, name, email").
		Where("email = ?", email).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal whether email exists
			return nil, nil
		}
		return nil, err
	}

	if member.Email == "" {
		return nil, nil
	}

	token, err := auth.SignResetToken(member.ID)
	if err != nil {
		return nil, err
	}

	return &ForgotPasswordResult{
		MemberID:    member.ID,
		MemberName:  member.Name,
		MemberEmail: member.Email,
		ResetToken:  token,
	}, nil
}

// ResetPassword verifies a reset token and updates the member's password.
func ResetPassword(db *gorm.DB, token, newPassword string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	if len(newPassword) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	memberID, err := auth.ParseResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	r := db.Table("members").Where("id = ?", memberID).Update("password", hash)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

// ─── Reset Password Mailable (plugin-scoped) ──────────────────────────────────

// ResetPasswordMailable implements mail.Mailable for reset-password emails.
type ResetPasswordMailable struct {
	subject      string
	templateBase string
	data         map[string]interface{}
}

func (r *ResetPasswordMailable) Subject() string {
	return r.subject
}

func (r *ResetPasswordMailable) TemplateBase() string {
	return r.templateBase
}

func (r *ResetPasswordMailable) Data() map[string]interface{} {
	return r.data
}

func (r *ResetPasswordMailable) From() (string, string) {
	return "", ""
}

// SendResetPasswordEmail sends a reset-password link asynchronously.
func SendResetPasswordEmail(toEmail, toName, resetLink string) error {
	data := map[string]interface{}{
		"Name":          toName,
		"ResetLink":     resetLink,
		"ExpiryMinutes": os.Getenv("RESET_TOKEN_TTL_MIN"),
	}

	m := &ResetPasswordMailable{
		subject:      "Reset Password — donis_finance",
		templateBase: "plugins/donisfinance/templates/email/reset_password",
		data:         data,
	}

	mailer := mail.NewMailer()
	mailer.Queue(toEmail, m)
	return nil
}
