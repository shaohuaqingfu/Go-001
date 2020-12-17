package service

import (
	"Week04/api/product"
	"Week04/internal/Week04/biz"
	"context"
)

type ProductService struct {
	biz *biz.ProductBiz
}

func NewProductService(productBiz *biz.ProductBiz) *ProductService {
	return &ProductService{
		productBiz,
	}
}

func (srv *ProductService) GetProductById(ctx context.Context, request *product.ProductRequest) (*product.ProductResponse, error) {
	panic("implement me")
}

func (srv *ProductService) GetProductByName(ctx context.Context, request *product.ProductRequest) (*product.ProductResponse, error) {
	panic("implement me")
}

func (srv *ProductService) GetProductByCreateUserId(ctx context.Context, request *product.ProductRequest) (*product.ProductListResponse, error) {
	panic("implement me")
}
