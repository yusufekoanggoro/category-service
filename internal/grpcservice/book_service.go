package grpcservice

import (
	"category-service/proto/book"
	"context"
	"log"

	"google.golang.org/grpc"
)

type BookGRPCClient struct {
	client book.BookServiceClient
}

func NewBookGRPCClient(bookServiceAddr string) *BookGRPCClient {
	conn, err := grpc.Dial(bookServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Book Service: %v", err)
	}

	client := book.NewBookServiceClient(conn)
	return &BookGRPCClient{client: client}
}

func (c *BookGRPCClient) SaveCategory(ctx context.Context, req *book.CategoryData) (*book.BookResponse, error) {
	res, err := c.client.ReceiveCategory(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}

func (c *BookGRPCClient) DeleteCategory(ctx context.Context, categoryId uint) (*book.BookResponse, error) {
	req := &book.DeleteData{Id: int64(categoryId)}
	res, err := c.client.DeleteCategory(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}
