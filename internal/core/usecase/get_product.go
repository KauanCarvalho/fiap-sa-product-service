package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
)

type GetProductUseCase interface {
	Run(ctx context.Context, sku string) (*entities.Product, error)
}

type getProductUseCase struct {
	ds domain.Datastore
}

func NewGetProductUseCase(ds domain.Datastore) GetProductUseCase {
	return &getProductUseCase{ds: ds}
}

func (gp *getProductUseCase) Run(ctx context.Context, sku string) (*entities.Product, error) {
	product, err := gp.ds.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	return product, nil
}
