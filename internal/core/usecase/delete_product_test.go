package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteProductUseCase_Run(t *testing.T) {
	t.Run("should soft delete a product from category Lanche", func(t *testing.T) {
		prepareTestDatabase()

		sku := "hamburger"

		err := dpu.Run(ctx, sku)
		require.NoError(t, err)
	})

	t.Run("should return error when product does not exist", func(t *testing.T) {
		prepareTestDatabase()

		err := dpu.Run(ctx, "NON_EXISTENT_SKU")
		require.Error(t, err)
	})
}
