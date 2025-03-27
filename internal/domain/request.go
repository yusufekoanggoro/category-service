package domain

type PaginationRequest struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1,max=100"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type UpdateCategoryRequest struct {
	ID   uint    `json:"id" binding:"required"`
	Name *string `json:"name" binding:"required"`
	Bio  *string `json:"bio"`
}
