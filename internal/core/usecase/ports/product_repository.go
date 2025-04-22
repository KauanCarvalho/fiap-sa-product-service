package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
)

type ProductFilter struct {
	Category string
	Page     int
	PageSize int
}

type ProductRepository interface {
	GetProductBySKU(ctx context.Context, sku string) (*entities.Product, error)
	GetAllProduct(ctx context.Context, filter *ProductFilter) ([]*entities.Product, int, error)
	CreateProduct(ctx context.Context, product *entities.Product) error
	UpdateProduct(ctx context.Context, product *entities.Product) error
	DeleteProduct(ctx context.Context, sku string) error
}
