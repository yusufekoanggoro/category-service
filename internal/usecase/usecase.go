package usecase

import (
	"category-service/internal/domain"
	sharedDomain "category-service/pkg/shared/domain"
	"context"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest) (*sharedDomain.Category, error)
	GetAllCategories(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error)
	GetCategoryByID(ctx context.Context, id uint) (*sharedDomain.Category, error)
	UpdateCategory(ctx context.Context, req *domain.UpdateCategoryRequest) (*sharedDomain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
}
