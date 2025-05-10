package datastore

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
	internalErrors "github.com/KauanCarvalho/fiap-sa-product-service/internal/shared/errors"
	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

func (ds *datastore) GetProductBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product *entities.Product

	if err := ds.db.WithContext(ctx).Model(product).
		Preload("Category").
		Preload("Images").
		Where("LOWER(sku) = LOWER(?)", sku).
		First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, internalErrors.NewInternalError("Failed to get product", err)
	}

	return product, nil
}

func (ds *datastore) GetAllProduct(ctx context.Context, filter *ports.ProductFilter) ([]*entities.Product, int, error) {
	if filter == nil {
		filter = &ports.ProductFilter{}
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize

	var products []*entities.Product
	query := ds.db.WithContext(ctx).Model(&entities.Product{}).
		Preload("Category").
		Preload("Images").
		Where("products.deleted_at IS NULL")

	if filter.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(categories.name) = LOWER(?)", filter.Category)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, internalErrors.NewInternalError("Failed to count products", err)
	}

	if err := query.
		Limit(filter.PageSize).
		Offset(offset).
		Order("products.id").
		Find(&products).Error; err != nil {
		return nil, 0, internalErrors.NewInternalError("Failed to get products", err)
	}

	return products, int(total), nil
}

func (ds *datastore) CreateProduct(ctx context.Context, product *entities.Product) error {
	err := ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			if isDuplicateSKUError(err) {
				return ErrExistingRecord
			}
			return internalErrors.NewInternalError("Failed to create product", err)
		}

		if err := tx.Preload("Category").Preload("Images").First(&product, product.ID).Error; err != nil {
			return internalErrors.NewInternalError("Failed to reload product", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func isDuplicateSKUError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "sku")
	}
	return false
}

func (ds *datastore) UpdateProduct(ctx context.Context, product *entities.Product) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingProduct := &entities.Product{SKU: product.SKU}
		if err := ds.findProductBySKU(ctx, tx, existingProduct); err != nil {
			return err
		}

		if err := ds.deleteProductImages(ctx, tx, existingProduct.ID); err != nil {
			return err
		}

		product.ID = existingProduct.ID
		if err := ds.updateProductFields(ctx, tx, product); err != nil {
			if isDuplicateSKUError(err) {
				return ErrExistingRecord
			}
		}

		if err := ds.recreateImages(ctx, tx, product); err != nil {
			return err
		}

		return ds.reloadProductWithAssociations(ctx, tx, product)
	})
}

func (ds *datastore) findProductBySKU(ctx context.Context, tx *gorm.DB, product *entities.Product) error {
	if err := tx.WithContext(ctx).Unscoped().Where("LOWER(sku) = LOWER(?)", product.SKU).First(product).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return internalErrors.NewInternalError("Product not found", err)
		}
		return err
	}
	return nil
}

func (ds *datastore) deleteProductImages(ctx context.Context, tx *gorm.DB, productID uint) error {
	if err := tx.WithContext(ctx).Where("product_id = ?", productID).Delete(&entities.Image{}).Error; err != nil {
		return internalErrors.NewInternalError("Failed to delete product images", err)
	}
	return nil
}

func (ds *datastore) updateProductFields(ctx context.Context, tx *gorm.DB, product *entities.Product) error {
	return tx.WithContext(ctx).
		Unscoped().
		Model(&entities.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category_id": product.CategoryID,
			"deleted_at":  nil,
		}).Error
}

func (ds *datastore) recreateImages(ctx context.Context, tx *gorm.DB, product *entities.Product) error {
	for _, image := range product.Images {
		image.ProductID = product.ID
		if err := tx.WithContext(ctx).Create(&image).Error; err != nil {
			return internalErrors.NewInternalError("Failed to create product image", err)
		}
	}
	return nil
}

func (ds *datastore) reloadProductWithAssociations(ctx context.Context, tx *gorm.DB, product *entities.Product) error {
	if err := tx.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		First(product, product.ID).Error; err != nil {
		return internalErrors.NewInternalError("Failed to reload product with images and category", err)
	}
	return nil
}

func (ds *datastore) DeleteProduct(ctx context.Context, sku string) error {
	result := ds.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("LOWER(sku) = LOWER(?)", sku).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return internalErrors.NewInternalError("Failed to delete product", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
