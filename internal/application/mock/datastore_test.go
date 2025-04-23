package mock_test

import (
	"context"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
	"github.com/stretchr/testify/require"
)

func TestDatastoreMock_Ping(t *testing.T) {
	t.Run("when PingFn is defined, it returns nil", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			PingFn: func(_ context.Context) error {
				return nil
			},
		}

		err := ds.Ping(ctx)
		require.NoError(t, err)
	})

	t.Run("when PingFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.Ping(ctx)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_GetProductBySKU(t *testing.T) {
	t.Run("when GetProductBySKUFn is defined, it returns the product", func(t *testing.T) {
		expectedProduct := &entities.Product{SKU: "123", Name: "Product A"}

		ds := &mock.DatastoreMock{
			GetProductBySKUFn: func(_ context.Context, _ string) (*entities.Product, error) {
				return expectedProduct, nil
			},
		}

		product, err := ds.GetProductBySKU(ctx, "123")
		require.NoError(t, err)
		require.Equal(t, expectedProduct, product)
	})

	t.Run("when GetProductBySKUFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		product, err := ds.GetProductBySKU(ctx, "123")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, product)
	})
}

func TestDatastoreMock_GetAllProduct(t *testing.T) {
	t.Run("when GetAllProductFn is defined, it returns the products", func(t *testing.T) {
		expectedProducts := []*entities.Product{
			{SKU: "123", Name: "Product A"},
			{SKU: "456", Name: "Product B"},
		}

		ds := &mock.DatastoreMock{
			GetAllProductFn: func(_ context.Context, _ *ports.ProductFilter) ([]*entities.Product, int, error) {
				return expectedProducts, len(expectedProducts), nil
			},
		}

		products, count, err := ds.GetAllProduct(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, expectedProducts, products)
		require.Equal(t, len(expectedProducts), count)
	})

	t.Run("when GetAllProductFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		products, count, err := ds.GetAllProduct(ctx, nil)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, products)
		require.Equal(t, 0, count)
	})
}

func TestDatastoreMock_FindCategoryByName(t *testing.T) {
	t.Run("when FindCategoryByNameFn is defined, it returns the category", func(t *testing.T) {
		expectedCategory := &entities.Category{Name: "Electronics"}

		ds := &mock.DatastoreMock{
			FindCategoryByNameFn: func(_ context.Context, _ string) (*entities.Category, error) {
				return expectedCategory, nil
			},
		}

		category, err := ds.FindCategoryByName(ctx, "123")
		require.NoError(t, err)
		require.Equal(t, expectedCategory, category)
	})

	t.Run("when FindCategoryByNameFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		category, err := ds.FindCategoryByName(ctx, "123")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, category)
	})
}

func TestDatastoreMock_CreateProduct(t *testing.T) {
	t.Run("when CreateProductFn is defined, it creates the product", func(t *testing.T) {
		product := &entities.Product{SKU: "123", Name: "Product A"}

		ds := &mock.DatastoreMock{
			CreateProductFn: func(_ context.Context, _ *entities.Product) error {
				return nil
			},
		}

		err := ds.CreateProduct(ctx, product)
		require.NoError(t, err)
	})

	t.Run("when CreateProductFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.CreateProduct(ctx, &entities.Product{SKU: "123", Name: "Product A"})
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_UpdateProduct(t *testing.T) {
	t.Run("when UpdateProductFn is defined, it updates the product", func(t *testing.T) {
		product := &entities.Product{SKU: "123", Name: "Product A"}

		ds := &mock.DatastoreMock{
			UpdateProductFn: func(_ context.Context, _ *entities.Product) error {
				return nil
			},
		}

		err := ds.UpdateProduct(ctx, product)
		require.NoError(t, err)
	})

	t.Run("when UpdateProductFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.UpdateProduct(ctx, &entities.Product{SKU: "123", Name: "Product A"})
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_DeleteProduct(t *testing.T) {
	t.Run("when DeleteProductFn is defined, it deletes the product", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			DeleteProductFn: func(_ context.Context, _ string) error {
				return nil
			},
		}

		err := ds.DeleteProduct(ctx, "123")
		require.NoError(t, err)
	})

	t.Run("when DeleteProductFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.DeleteProduct(ctx, "123")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}
