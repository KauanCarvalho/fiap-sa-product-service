package usecase_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProductUseCase_Run(t *testing.T) {
	cu := usecase.NewUpdateProductUseCase(ds)

	t.Run("should update product with valid category", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputUpdate{
			Name:        "X-Tudo Atualizado",
			Description: "Sanduíche completo com mais ingredientes",
			Price:       32.90,
			SKU:         "hamburger",
			Category: dto.ProductCategoryUpdate{
				Name: "Lanche",
			},
			Images: []dto.ProductImageUpdate{
				{URL: "https://example.com/xtudo-atualizado.jpg"},
			},
		}

		product, err := cu.Run(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, product)

		assert.Equal(t, "X-Tudo Atualizado", product.Name)
		assert.Equal(t, "lanche", product.Category.Name)
		assert.Len(t, product.Images, 1)
		assert.Contains(t, product.Images[0].URL, "xtudo-atualizado.jpg")
	})

	t.Run("should return error if category does not exist", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputUpdate{
			Name:        "Cachorro Quente Atualizado",
			Description: "Hot dog com ingredientes diferentes",
			Price:       16.90,
			SKU:         "cachorro-001",
			Category: dto.ProductCategoryUpdate{
				Name: "Inexistente",
			},
		}

		product, err := cu.Run(ctx, input)
		require.Error(t, err)
		require.Nil(t, product)
	})

	t.Run("should return error if product with SKU does not exist", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ProductInputUpdate{
			Name:        "X-Burguer Atualizado",
			Description: "X-Burguer com mais carne",
			Price:       22.90,
			SKU:         "nao-existe",
			Category: dto.ProductCategoryUpdate{
				Name: "Lanche",
			},
		}

		product, err := cu.Run(ctx, input)
		require.Error(t, err)
		require.Nil(t, product)
	})
}
