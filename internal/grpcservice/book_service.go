package grpcservice

import (
	"category-service/proto/category"
	"context"
	"log"

	"google.golang.org/grpc"
)

type BookGRPCClient struct {
	client category.CategoryServiceClient
}

func NewBookGRPCClient(bookServiceAddr string) *BookGRPCClient {
	conn, err := grpc.Dial(bookServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Book Service: %v", err)
	}

	client := category.NewCategoryServiceClient(conn)
	return &BookGRPCClient{client: client}
}

func (c *BookGRPCClient) SaveCategory(ctx context.Context, req *category.SaveCategoryRequest) (*category.SaveCategoryResponse, error) {
	res, err := c.client.SaveCategory(ctx, req)
	if err != nil {
		return &category.SaveCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}

func (c *BookGRPCClient) DeleteCategory(ctx context.Context, categoryId uint) (*category.DeleteCategoryResponse, error) {
	req := &category.DeleteCategoryRequest{CategoryID: int64(categoryId)}
	res, err := c.client.DeleteCategory(ctx, req)
	if err != nil {
		return &category.DeleteCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}
