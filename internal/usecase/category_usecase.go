package usecase

import (
	"category-service/internal/domain"
	"category-service/internal/grpcservice"
	"category-service/internal/repository"
	sharedDomain "category-service/pkg/shared/domain"
	"context"

	protoCategory "category-service/proto/category"
)

type categoryUsecase struct {
	repo       repository.CategoryRepository
	bookClient *grpcservice.BookGRPCClient
}

func NewAuthorUsecase(
	repo repository.CategoryRepository,
	bookClient *grpcservice.BookGRPCClient,
) CategoryUsecase {
	return &categoryUsecase{repo: repo, bookClient: bookClient}
}

func (uc *categoryUsecase) CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest) (*sharedDomain.Category, error) {
	var category *sharedDomain.Category
	newCategory := &sharedDomain.Category{
		Name: req.Name,
	}

	err := uc.repo.SaveCategory(ctx, newCategory)
	if err != nil {
		return nil, err
	}
	category = newCategory

	grpcRequest := &protoCategory.SaveCategoryRequest{CategoryID: int64(category.ID), Name: category.Name}
	_, err = uc.bookClient.SaveCategory(ctx, grpcRequest)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *categoryUsecase) GetAllCategories(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error) {
	categories, totalRows, err := uc.repo.GetAllCategories(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	paginatedResponse := &domain.PaginatedResponse{
		Data:       categories,
		Total:      totalRows,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((totalRows + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	return paginatedResponse, nil
}

func (uc *categoryUsecase) GetCategoryByID(ctx context.Context, id uint) (*sharedDomain.Category, error) {
	category, err := uc.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) UpdateCategory(ctx context.Context, req *domain.UpdateCategoryRequest) (*sharedDomain.Category, error) {
	var category *sharedDomain.Category
	existingCategory, err := uc.repo.GetCategoryByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingCategory.Name = *req.Name
	}

	err = uc.repo.SaveCategory(ctx, existingCategory)
	if err != nil {
		return nil, err
	}
	category = existingCategory

	grpcRequest := &protoCategory.SaveCategoryRequest{CategoryID: int64(category.ID), Name: category.Name}
	_, err = uc.bookClient.SaveCategory(ctx, grpcRequest)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *categoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	err := uc.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	_, err = uc.bookClient.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
