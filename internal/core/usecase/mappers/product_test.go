package mappers_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/mappers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToProductsDTO(t *testing.T) {
	t.Run("should map products correctly", func(t *testing.T) {
		products := []*entities.Product{
			{
				ID:          1,
				Name:        "X-Tudo",
				Description: "Complete sandwich",
				Price:       30.50,
				SKU:         "xtudo-001",
				Category:    entities.Category{Name: "Snack"},
				Images: []entities.Image{
					{URL: "https://example.com/xtudo.jpg"},
				},
			},
			{
				ID:          2,
				Name:        "X-Burger",
				Description: "Burger with extra meat",
				Price:       22.90,
				SKU:         "xburger-001",
				Category:    entities.Category{Name: "Burger"},
				Images: []entities.Image{
					{URL: "https://example.com/xburger.jpg"},
				},
			},
		}

		result := mappers.ToProductsDTO(products)

		require.Len(t, result, 2)

		assert.Equal(t, "X-Tudo", result[0].Name)
		assert.Equal(t, "Complete sandwich", result[0].Description)
		assert.InEpsilon(t, 30.50, result[0].Price, 0.01)
		assert.Equal(t, "xtudo-001", result[0].SKU)
		assert.Equal(t, "Snack", result[0].Category.Name)
		assert.Len(t, result[0].Images, 1)
		assert.Equal(t, "https://example.com/xtudo.jpg", result[0].Images[0].URL)

		assert.Equal(t, "X-Burger", result[1].Name)
		assert.Equal(t, "Burger with extra meat", result[1].Description)
		assert.InEpsilon(t, 22.90, result[1].Price, 0.01)
		assert.Equal(t, "xburger-001", result[1].SKU)
		assert.Equal(t, "Burger", result[1].Category.Name)
		assert.Len(t, result[1].Images, 1)
		assert.Equal(t, "https://example.com/xburger.jpg", result[1].Images[0].URL)
	})

	t.Run("should return empty slice if no products", func(t *testing.T) {
		var products []*entities.Product

		result := mappers.ToProductsDTO(products)

		assert.Empty(t, result)
	})

	t.Run("should skip product with ID 0", func(t *testing.T) {
		products := []*entities.Product{
			{
				ID:          0,
				Name:        "Invalid Product",
				Description: "This product should be skipped",
				Price:       10.00,
				SKU:         "invalid-001",
				Category:    entities.Category{Name: "Invalid"},
				Images: []entities.Image{
					{URL: "https://example.com/invalid.jpg"},
				},
			},
			{
				ID:          1,
				Name:        "Valid Product",
				Description: "This product should be included",
				Price:       25.00,
				SKU:         "valid-001",
				Category:    entities.Category{Name: "Valid"},
				Images: []entities.Image{
					{URL: "https://example.com/valid.jpg"},
				},
			},
		}

		result := mappers.ToProductsDTO(products)

		require.Len(t, result, 1)
		assert.Equal(t, "Valid Product", result[0].Name)
		assert.Equal(t, "This product should be included", result[0].Description)
	})
}

func TestToProductDTO(t *testing.T) {
	t.Run("should map product correctly", func(t *testing.T) {
		product := entities.Product{
			ID:          1,
			Name:        "X-Tudo",
			Description: "Complete sandwich",
			Price:       30.50,
			SKU:         "xtudo-001",
			Category:    entities.Category{Name: "Snack"},
			Images: []entities.Image{
				{URL: "https://example.com/xtudo.jpg"},
			},
		}

		result := mappers.ToProductDTO(product)

		assert.Equal(t, "X-Tudo", result.Name)
		assert.Equal(t, "Complete sandwich", result.Description)
		assert.InEpsilon(t, 30.50, result.Price, 0.01)
		assert.Equal(t, "xtudo-001", result.SKU)
		assert.Equal(t, "Snack", result.Category.Name)
		assert.Len(t, result.Images, 1)
		assert.Equal(t, "https://example.com/xtudo.jpg", result.Images[0].URL)
	})

	t.Run("should handle product with no images", func(t *testing.T) {
		product := entities.Product{
			ID:          1,
			Name:        "X-Tudo",
			Description: "Complete sandwich",
			Price:       30.50,
			SKU:         "xtudo-001",
			Category:    entities.Category{Name: "Snack"},
			Images:      []entities.Image{},
		}

		result := mappers.ToProductDTO(product)

		assert.Equal(t, "X-Tudo", result.Name)
		assert.Equal(t, "Complete sandwich", result.Description)
		assert.InEpsilon(t, 30.50, result.Price, 0.01)
		assert.Equal(t, "xtudo-001", result.SKU)
		assert.Equal(t, "Snack", result.Category.Name)
		assert.Empty(t, result.Images)
	})

	t.Run("should skip product with ID 0", func(t *testing.T) {
		product := entities.Product{
			ID:          0,
			Name:        "Invalid Product",
			Description: "This product should be skipped",
			Price:       10.00,
			SKU:         "invalid-001",
			Category:    entities.Category{Name: "Invalid"},
			Images: []entities.Image{
				{URL: "https://example.com/invalid.jpg"},
			},
		}

		result := mappers.ToProductDTO(product)

		assert.Empty(t, result.Name)
		assert.Empty(t, result.Description)
		assert.Zero(t, result.Price)
		assert.Empty(t, result.SKU)
		assert.Empty(t, result.Category.Name)
		assert.Empty(t, result.Images)
	})
}
