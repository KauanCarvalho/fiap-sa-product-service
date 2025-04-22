package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/shared"
)

type CreateProductUseCase interface {
	Run(ctx context.Context, input dto.ProductInputCreate) (*entities.Product, error)
}

type createProductsUseCase struct {
	ds domain.Datastore
}

func NewCreateProductUseCase(ds domain.Datastore) CreateProductUseCase {
	return &createProductsUseCase{ds: ds}
}

func (c *createProductsUseCase) Run(ctx context.Context, input dto.ProductInputCreate) (*entities.Product, error) {
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
		SKU:         shared.Slugify(input.Name),
		CategoryID:  category.ID,
		Images:      images,
	}

	err = c.ds.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
