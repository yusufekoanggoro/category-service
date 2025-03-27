package response

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, statusCode int, message string, data interface{}, pagination Pagination) {
	c.JSON(statusCode, gin.H{
		"status":     "success",
		"message":    message,
		"data":       data,
		"pagination": pagination,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
	})
}
