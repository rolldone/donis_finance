package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/rolldone/donisgo/plugins/donisfinance/models"

	"gorm.io/gorm"
)

// CreateCategory creates a new category.
func CreateCategory(db *gorm.DB, name, catType, icon, color string) (*CategoryResult, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if catType != "income" && catType != "expense" {
		return nil, errors.New("type must be 'income' or 'expense'")
	}

	// Duplicate check: same name + type = rejected (soft check, not DB constraint)
	// This allows re-creating a category with the same name after deletion.
	var existing models.Category
	if err := db.Where("name = ? AND type = ?", name, catType).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("kategori dengan nama ini sudah ada (id: %s)", existing.ID)
	}

	cat := models.Category{
		Name: name,
		Type: catType,
		Icon: icon,
	}
	if color != "" {
		cat.Color = color
	}

	if err := db.Create(&cat).Error; err != nil {
		return nil, errors.New("failed to create category: " + err.Error())
	}

	return &CategoryResult{
		ID:    cat.ID,
		Name:  cat.Name,
		Type:  cat.Type,
		Icon:  cat.Icon,
		Color: cat.Color,
	}, nil
}

// UpdateCategory updates an existing category.
func UpdateCategory(db *gorm.DB, id, name, catType, icon, color string) (*CategoryResult, error) {
	var cat models.Category
	if err := db.First(&cat, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to fetch category: " + err.Error())
	}

	if name != "" {
		cat.Name = name
	}
	if catType == "income" || catType == "expense" {
		cat.Type = catType
	}
	if icon != "" {
		cat.Icon = icon
	}
	if color != "" {
		cat.Color = color
	}
	cat.UpdatedAt = time.Now()

	if err := db.Save(&cat).Error; err != nil {
		return nil, errors.New("failed to update category: " + err.Error())
	}

	return &CategoryResult{
		ID:    cat.ID,
		Name:  cat.Name,
		Type:  cat.Type,
		Icon:  cat.Icon,
		Color: cat.Color,
	}, nil
}

// DeleteCategory deletes a category by ID.
func DeleteCategory(db *gorm.DB, id string) error {
	result := db.Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return errors.New("failed to delete category: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
