package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
)

type GetProductsUseCase interface {
	Run(ctx context.Context, filter *ports.ProductFilter) ([]*entities.Product, int, error)
}

type getProductsUseCase struct {
	ds domain.Datastore
}

func NewGetProductsUseCase(ds domain.Datastore) GetProductsUseCase {
	return &getProductsUseCase{ds: ds}
}

func (gp *getProductsUseCase) Run(ctx context.Context, filter *ports.ProductFilter) ([]*entities.Product, int, error) {
	products, total, err := gp.ds.GetAllProduct(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
