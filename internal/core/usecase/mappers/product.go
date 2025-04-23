package mappers

import (
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
)

func ToProductsDTO(products []*entities.Product) []dto.ProductOutput {
	productOutputs := make([]dto.ProductOutput, 0, len(products))
	for _, product := range products {
		if product.ID == 0 {
			continue
		}

		productOutputs = append(productOutputs, ToProductDTO(*product))
	}

	return productOutputs
}

func ToProductDTO(product entities.Product) dto.ProductOutput {
	if product.ID == 0 {
		return dto.ProductOutput{}
	}

	imagesOutputs := make([]dto.ProductImageOutput, 0, len(product.Images))
	for _, image := range product.Images {
		imagesOutputs = append(imagesOutputs, dto.ProductImageOutput{
			URL: image.URL,
		})
	}

	return dto.ProductOutput{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		SKU:         product.SKU,
		Category: dto.ProductCategoryOutput{
			Name: product.Category.Name,
		},
		Images: imagesOutputs,
	}
}
