package repository

import (
	sharedDomain "category-service/pkg/shared/domain"
	"errors"
	"log"
	"time"

	customErr "category-service/pkg/errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories(page, limit int) ([]*sharedDomain.Category, int64, error) {
	var categories []*sharedDomain.Category
	var totalRows int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if err := r.db.Model(&sharedDomain.Category{}).Count(&totalRows).Error; err != nil {
		log.Println("GetAllCategories count error:", err)
		return nil, 0, customErr.ErrInternalError
	}

	offset := (page - 1) * limit

	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		log.Println("GetAllCategories query error:", err)
		return nil, 0, customErr.ErrInternalError
	}

	return categories, totalRows, nil
}

func (r *categoryRepository) GetCategoryByID(id uint) (*sharedDomain.Category, error) {
	var category sharedDomain.Category

	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.ErrNotFound
		}
		log.Println("GetCategoryByID error:", err)
		return nil, customErr.ErrInternalError
	}

	return &category, nil
}

func (r *categoryRepository) SaveCategory(tx *gorm.DB, category *sharedDomain.Category) error {
	if err := tx.Save(category).Error; err != nil {
		return customErr.ErrInternalError
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(tx *gorm.DB, id uint) error {
	var category sharedDomain.Category
	if err := tx.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customErr.ErrNotFound
		}
		log.Println("Find category error:", err)
		return customErr.ErrInternalError
	}

	if err := tx.Model(&category).Update("deleted_at", time.Now()).Error; err != nil {
		log.Println("Soft delete category error:", err)
		return customErr.ErrInternalError
	}

	return nil
}

func (r *categoryRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn) // Auto commit atau rollback
}
