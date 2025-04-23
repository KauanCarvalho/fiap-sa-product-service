package mock

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
)

type DatastoreMock struct {
	PingFn               func(ctx context.Context) error
	GetProductBySKUFn    func(ctx context.Context, sku string) (*entities.Product, error)
	GetAllProductFn      func(ctx context.Context, filter *ports.ProductFilter) ([]*entities.Product, int, error)
	FindCategoryByNameFn func(ctx context.Context, sku string) (*entities.Category, error)
	CreateProductFn      func(ctx context.Context, product *entities.Product) error
	UpdateProductFn      func(ctx context.Context, product *entities.Product) error
	DeleteProductFn      func(ctx context.Context, sku string) error
}

var ErrFunctionNotImplemented = errors.New("function not implemented")

func (m *DatastoreMock) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) GetProductBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	if m.GetProductBySKUFn != nil {
		return m.GetProductBySKUFn(ctx, sku)
	}

	return nil, ErrFunctionNotImplemented
}

func (m *DatastoreMock) GetAllProduct(
	ctx context.Context,
	filter *ports.ProductFilter,
) ([]*entities.Product, int, error) {
	if m.GetAllProductFn != nil {
		return m.GetAllProductFn(ctx, filter)
	}

	return nil, 0, ErrFunctionNotImplemented
}

func (m *DatastoreMock) FindCategoryByName(ctx context.Context, sku string) (*entities.Category, error) {
	if m.FindCategoryByNameFn != nil {
		return m.FindCategoryByNameFn(ctx, sku)
	}

	return nil, ErrFunctionNotImplemented
}

func (m *DatastoreMock) CreateProduct(ctx context.Context, product *entities.Product) error {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, product)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) UpdateProduct(ctx context.Context, product *entities.Product) error {
	if m.UpdateProductFn != nil {
		return m.UpdateProductFn(ctx, product)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) DeleteProduct(ctx context.Context, sku string) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, sku)
	}

	return ErrFunctionNotImplemented
}
