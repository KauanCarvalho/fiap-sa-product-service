package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"
)

type UpdateProductUseCase interface {
	Run(ctx context.Context, input dto.ProductInputUpdate) (*entities.Product, error)
}

type updateProductUseCase struct {
	ds domain.Datastore
}

func NewUpdateProductUseCase(ds domain.Datastore) UpdateProductUseCase {
	return &updateProductUseCase{ds: ds}
}

func (c *updateProductUseCase) Run(ctx context.Context, input dto.ProductInputUpdate) (*entities.Product, error) {
	images := make([]entities.Image, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, entities.Image{
			URL: image.URL,
		})
	}

	category, err := c.ds.FindCategoryByName(ctx, input.Category.Name)
	if err != nil {
		return nil, err
	}

	product := &entities.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SKU:         input.SKU,
		CategoryID:  category.ID,
		Images:      images,
	}

	err = c.ds.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
