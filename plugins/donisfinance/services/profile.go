package services

import (
	"fmt"

	"go_framework/plugins/donisfinance/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MemberProfile represents the public profile of a member.
type MemberProfile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// AdminProfile represents the public profile of an admin.
type AdminProfile struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	CreatedAt string `json:"created_at"`
}

// GetMemberProfile returns a member's profile by ID.
func GetMemberProfile(db *gorm.DB, memberID string) (*MemberProfile, error) {
	var m models.Member
	if err := db.First(&m, "id = ?", memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found")
	}
	return &MemberProfile{
		ID:        m.ID,
		Name:      m.Name,
		Username:  m.Username,
		CreatedAt: m.CreatedAt.Format("2006-01-02"),
	}, nil
}

// UpdateMemberProfile updates name and/or username for a member.
func UpdateMemberProfile(db *gorm.DB, memberID, name, username string) (*MemberProfile, error) {
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

	if err := db.Model(&models.Member{}).Where("id = ?", memberID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return GetMemberProfile(db, memberID)
}

// ChangeMemberPassword changes a member's password after verifying old password.
func ChangeMemberPassword(db *gorm.DB, memberID, oldPassword, newPassword string) error {
	var m models.Member
	if err := db.First(&m, "id = ?", memberID).Error; err != nil {
		return fmt.Errorf("member not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	if len(newPassword) < 6 {
		return fmt.Errorf("new password must be at least 6 characters")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return db.Model(&m).Update("password", hash).Error
}

// GetAdminProfile returns an admin's profile by ID.
func GetAdminProfile(db *gorm.DB, adminID string) (*AdminProfile, error) {
	var a models.Admin
	if err := db.First(&a, "id = ?", adminID).Error; err != nil {
		return nil, fmt.Errorf("admin not found")
	}
	return &AdminProfile{
		ID:        a.ID,
		Username:  a.Username,
		Email:     a.Email,
		CreatedAt: a.CreatedAt.Format("2006-01-02"),
	}, nil
}

// UpdateAdminProfile updates username and/or email for an admin.
func UpdateAdminProfile(db *gorm.DB, adminID, username, email string) (*AdminProfile, error) {
	updates := map[string]interface{}{}
	if username != "" {
		updates["username"] = username
	}
	if email != "" {
		updates["email"] = email
	}
	if len(updates) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	if err := db.Model(&models.Admin{}).Where("id = ?", adminID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return GetAdminProfile(db, adminID)
}

// ChangeAdminPassword changes an admin's password after verifying old password.
func ChangeAdminPassword(db *gorm.DB, adminID, oldPassword, newPassword string) error {
	var a models.Admin
	if err := db.First(&a, "id = ?", adminID).Error; err != nil {
		return fmt.Errorf("admin not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	if len(newPassword) < 6 {
		return fmt.Errorf("new password must be at least 6 characters")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return db.Model(&a).Update("password", hash).Error
}
