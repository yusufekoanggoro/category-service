package repository

import (
	sharedDomain "category-service/pkg/shared/domain"
	"context"
)

type CategoryRepository interface {
	SaveCategory(ctx context.Context, category *sharedDomain.Category) error
	GetAllCategories(ctx context.Context, page, limit int) ([]*sharedDomain.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*sharedDomain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
}
