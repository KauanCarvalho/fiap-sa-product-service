package datastore_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetProductBySKU(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully get product by SKU", func(t *testing.T) {
		product := &entities.Product{
			Name:        "Suco de Laranja",
			Description: "Natural e gelado",
			Price:       9.99,
			SKU:         "suco-laranja",
			Category:    entities.Category{Name: "Bebidas"},
			Images: []entities.Image{
				{URL: "http://img/suco-laranja-1"},
			},
		}
		require.NoError(t, sqlDB.Create(&product.Category).Error)
		product.CategoryID = product.Category.ID
		require.NoError(t, sqlDB.Create(&product).Error)

		found, err := ds.GetProductBySKU(ctx, "suco-laranja")
		require.NoError(t, err)
		assert.Equal(t, "Suco de Laranja", found.Name)
		assert.Len(t, found.Images, 1)
	})

	t.Run("fail to get non-existent product", func(t *testing.T) {
		found, err := ds.GetProductBySKU(ctx, "nao-existe")
		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestGetAllProduct(t *testing.T) {
	prepareTestDatabase()

	t.Run("get all products with pagination", func(t *testing.T) {
		filter := &ports.ProductFilter{Page: 1, PageSize: 10}
		products, total, err := ds.GetAllProduct(ctx, filter)

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(products), total)
	})

	t.Run("get all products with pagination with inncorrect filters", func(t *testing.T) {
		filter := &ports.ProductFilter{Page: -1, PageSize: -10}
		products, total, err := ds.GetAllProduct(ctx, filter)

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(products), total)
	})

	t.Run("filter by category", func(t *testing.T) {
		filter := &ports.ProductFilter{Page: 1, PageSize: 10, Category: "lanche"}
		products, total, err := ds.GetAllProduct(ctx, filter)

		require.NoError(t, err)
		for _, p := range products {
			assert.Equal(t, "lanche", p.Category.Name)
		}
		require.GreaterOrEqual(t, len(products), total)
	})

	t.Run("when filter is nil", func(t *testing.T) {
		products, total, err := ds.GetAllProduct(ctx, nil)

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(products), total)
	})
}

func TestCreateProduct(t *testing.T) {
	prepareTestDatabase()

	t.Run("create product successfully", func(t *testing.T) {
		category := &entities.Category{Name: "Doces"}
		require.NoError(t, sqlDB.Create(category).Error)

		product := &entities.Product{
			Name:        "Bolo de Cenoura",
			Description: "Com cobertura de chocolate",
			Price:       14.50,
			SKU:         "bolo-cenoura",
			CategoryID:  category.ID,
			Images: []entities.Image{
				{URL: "http://img/bolo-1"},
				{URL: "http://img/bolo-2"},
			},
		}

		err := ds.CreateProduct(ctx, product)
		require.NoError(t, err)
		assert.NotZero(t, product.ID)
		assert.Len(t, product.Images, 2)
	})

	t.Run("fail to create product with non-existent category", func(t *testing.T) {
		product := &entities.Product{
			Name:        "Produto Não Existente",
			Description: "Descrição não encontrada",
			Price:       15.99,
			SKU:         "produto-nao-existente",
			CategoryID:  99999,
		}

		err := ds.CreateProduct(ctx, product)
		assert.Error(t, err)
	})

	t.Run("fail to create duplicate product", func(t *testing.T) {
		category := &entities.Category{Name: "Lanches"}
		require.NoError(t, sqlDB.Create(category).Error)

		product := &entities.Product{
			Name:        "X-Salada",
			Description: "Com carne e queijo",
			Price:       12.50,
			SKU:         "x-salada",
			CategoryID:  category.ID,
		}
		require.NoError(t, sqlDB.Create(&product).Error)

		duplicateProduct := &entities.Product{
			Name:        "X-Salada",
			Description: "Com carne e queijo",
			Price:       12.50,
			SKU:         "x-salada",
			CategoryID:  category.ID,
		}

		err := ds.CreateProduct(ctx, duplicateProduct)
		require.ErrorIs(t, err, datastore.ErrExistingRecord)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("update successfully", func(t *testing.T) {
		prepareTestDatabase()

		category := &entities.Category{Name: "Refrigerante"}
		require.NoError(t, sqlDB.Create(category).Error)

		product := &entities.Product{
			Name:        "Guaraná",
			Description: "Gelado",
			Price:       6.50,
			SKU:         "guarana",
			CategoryID:  category.ID,
			Images: []entities.Image{
				{URL: "http://img/guarana-1"},
			},
		}
		require.NoError(t, sqlDB.Create(&product).Error)

		product.Name = "Guaraná Gelado"
		product.Images = []entities.Image{
			{URL: "http://img/guarana-gelado"},
		}

		err := ds.UpdateProduct(ctx, product)
		require.NoError(t, err)

		var updated entities.Product
		err = sqlDB.Preload("Images").First(&updated, product.ID).Error

		require.NoError(t, err)
		assert.Equal(t, "Guaraná Gelado", updated.Name)
		assert.Len(t, updated.Images, 1)
	})

	t.Run("fail to update product when SKU not found", func(t *testing.T) {
		prepareTestDatabase()

		product := &entities.Product{
			Name:        "Produto Inexistente",
			Description: "Produto que não existe",
			Price:       99.99,
			SKU:         "produto-nao-existe",
		}

		err := ds.UpdateProduct(ctx, product)
		assert.Error(t, err)
	})
}

func TestDeleteProduct(t *testing.T) {
	prepareTestDatabase()

	t.Run("delete product successfully", func(t *testing.T) {
		category := &entities.Category{Name: "Sobremesa-1"}
		require.NoError(t, sqlDB.Create(&category).Error)

		product := &entities.Product{
			Name:        "Pudim",
			Description: "Doce de leite",
			Price:       5.50,
			SKU:         "pudim-doce",
			CategoryID:  category.ID,
		}
		require.NoError(t, sqlDB.Create(&product).Error)

		err := ds.DeleteProduct(ctx, "pudim-doce")
		require.NoError(t, err)

		var deleted entities.Product
		err = sqlDB.Unscoped().First(&deleted, product.ID).Error
		require.NoError(t, err)
		assert.NotNil(t, deleted.DeletedAt)
	})

	t.Run("return an error when the sku is not present", func(t *testing.T) {
		err := ds.DeleteProduct(ctx, "nao-existe")
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
