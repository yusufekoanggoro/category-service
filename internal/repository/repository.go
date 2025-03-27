package repository

import (
	sharedDomain "category-service/pkg/shared/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	SaveCategory(tx *gorm.DB, category *sharedDomain.Category) error
	GetAllCategories(page, limit int) ([]*sharedDomain.Category, int64, error)
	GetCategoryByID(id uint) (*sharedDomain.Category, error)
	DeleteCategory(tx *gorm.DB, id uint) error

	Transaction(fn func(tx *gorm.DB) error) error
}
