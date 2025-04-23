package usecase_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductsUseCase_Run(t *testing.T) {
	t.Run("should return all products when no filter is provided", func(t *testing.T) {
		prepareTestDatabase()

		products, total, err := gps.Run(ctx, nil)

		require.NoError(t, err)
		assert.Positive(t, total)
		assert.Len(t, products, total)
	})

	t.Run("should return paginated products", func(t *testing.T) {
		prepareTestDatabase()

		filter := &ports.ProductFilter{
			Page:     1,
			PageSize: 2,
		}

		products, total, err := gps.Run(ctx, filter)

		require.NoError(t, err)
		assert.LessOrEqual(t, len(products), 2)
		assert.Positive(t, total)
	})
}
