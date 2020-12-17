package data

import "github.com/jinzhu/gorm"

type ProductData struct {
	DB *gorm.DB
}

func NewProductData(db *gorm.DB) *ProductData {
	return &ProductData{
		DB: db,
	}
}
