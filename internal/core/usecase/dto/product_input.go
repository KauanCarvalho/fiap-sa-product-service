package dto

import "github.com/go-playground/validator/v10"

type ProductCategoryCreate struct {
	Name string `json:"name" validate:"required"`
}

type ProductImageCreate struct {
	URL string `json:"url" validate:"required,url"`
}

type ProductInputCreate struct {
	Name        string                `json:"name"        validate:"required"`
	Price       float64               `json:"price"       validate:"required,gte=0"`
	Description string                `json:"description" validate:"required"`
	Category    ProductCategoryCreate `json:"category"    validate:"required"`
	Images      []ProductImageCreate  `json:"images"      validate:"required,dive"`
}

type ProductCategoryUpdate struct {
	Name string `json:"name" validate:"required"`
}

type ProductImageUpdate struct {
	URL string `json:"url" validate:"required,url"`
}

type ProductInputUpdate struct {
	SKU         string                `json:"sku"         validate:"required"`
	Name        string                `json:"name"        validate:"required"`
	Price       float64               `json:"price"       validate:"required,gte=0"`
	Description string                `json:"description" validate:"required"`
	Category    ProductCategoryUpdate `json:"category"    validate:"required"`
	Images      []ProductImageUpdate  `json:"images"      validate:"required,dive"`
}

func ValidateProductCreate(input ProductInputCreate) error {
	return validator.New().Struct(input)
}

func ValidateProductUpdate(input ProductInputUpdate) error {
	return validator.New().Struct(input)
}
