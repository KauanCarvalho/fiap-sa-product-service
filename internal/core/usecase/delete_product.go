package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
)

type DeleteProductUseCase interface {
	Run(ctx context.Context, sku string) error
}

type deleteProductUseCase struct {
	ds domain.Datastore
}

func NewDeleteProductUseCase(ds domain.Datastore) DeleteProductUseCase {
	return &deleteProductUseCase{ds: ds}
}

func (c *deleteProductUseCase) Run(ctx context.Context, sku string) error {
	return c.ds.DeleteProduct(ctx, sku)
}
