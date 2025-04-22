package usecase_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProductUseCase_Run(t *testing.T) {
	t.Run("should create product with valid category", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputCreate{
			Name:        "X-Tudo",
			Description: "Sanduíche completo",
			Price:       29.90,
			Category: dto.ProductCategoryCreate{
				Name: "Lanche",
			},
			Images: []dto.ProductImageCreate{
				{URL: "https://example.com/xtudo.jpg"},
			},
		}

		product, err := cp.Run(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, product)

		assert.Equal(t, "X-Tudo", product.Name)
		assert.Equal(t, "lanche", product.Category.Name)
		assert.Len(t, product.Images, 1)
		assert.Contains(t, product.Images[0].URL, "xtudo.jpg")
	})

	t.Run("should return error if category does not exist", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputCreate{
			Name:        "Cachorro Quente",
			Description: "Hot dog completo",
			Price:       14.90,
			Category: dto.ProductCategoryCreate{
				Name: "Inexistente",
			},
		}

		product, err := cp.Run(ctx, input)
		require.Error(t, err)
		require.Nil(t, product)
	})

	t.Run("should return error if product with same SKU already exists", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputCreate{
			Name:        "hamburger",
			Description: "Outro X-Burguer",
			Price:       19.90,
			Category: dto.ProductCategoryCreate{
				Name: "Lanche",
			},
		}

		product, err := cp.Run(ctx, input)
		require.Error(t, err)
		require.Nil(t, product)
	})
}
