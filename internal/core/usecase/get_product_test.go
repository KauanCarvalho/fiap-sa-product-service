package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductUseCase_Run(t *testing.T) {
	prepareTestDatabase()

	t.Run("should return product when SKU exists", func(t *testing.T) {
		sku := "hamburger"

		product, err := gp.Run(ctx, sku)

		require.NoError(t, err)
		require.NotNil(t, product)
		assert.Equal(t, sku, product.SKU)
		assert.NotEmpty(t, product.Name)
		assert.NotNil(t, product.Category)
		assert.NotEmpty(t, product.Images)
	})

	t.Run("should return error when SKU does not exist", func(t *testing.T) {
		sku := "non-existent-sku"

		product, err := gp.Run(ctx, sku)

		require.Error(t, err)
		assert.Nil(t, product)
	})
}
