package http

import (
	"category-service/internal/domain"
	"category-service/internal/usecase"
	"category-service/pkg/shared/response"
	"net/http"
	"strconv"

	customErr "category-service/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewCategoryHandler(uc usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: uc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		if err == customErr.ErrInternalError {
			response.Error(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Category created successfully", book)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	categories, err := h.usecase.GetAllCategories(&req)
	if err != nil {
		if err == customErr.ErrInternalError {
			response.Error(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.Limit,
		TotalPages:  categories.TotalPages,
		TotalItems:  int(categories.Total),
	}

	response.SuccessWithPagination(c, http.StatusOK, "Categories retrieved successfully", categories.Data, pagination)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.usecase.GetCategoryByID(uint(id))
	if err != nil {
		if err == customErr.ErrNotFound {
			response.Error(c, http.StatusNotFound, "Category not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req domain.UpdateCategoryRequest
	req.ID = uint(categoryID)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		if err == customErr.ErrInternalError {
			response.Error(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Category updated successfully", book)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.usecase.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		if err == customErr.ErrInternalError {
			response.Error(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Category deleted successfully", nil)
}
