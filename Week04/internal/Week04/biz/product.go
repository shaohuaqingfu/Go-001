package biz

import (
	"Week04/internal/Week04/data"
)

type ProductBiz struct {
	dao *data.ProductData
}

func NewProductBiz(productData *data.ProductData) *ProductBiz {
	return &ProductBiz{
		dao: productData,
	}
}
