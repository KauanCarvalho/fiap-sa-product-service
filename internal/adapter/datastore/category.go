package datastore

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
	internalErrors "github.com/KauanCarvalho/fiap-sa-product-service/internal/shared/errors"

	"gorm.io/gorm"
)

func (ds *datastore) FindCategoryByName(ctx context.Context, name string) (*entities.Category, error) {
	category := &entities.Category{}

	err := ds.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		return nil, internalErrors.NewInternalError("Failed to find category", err)
	}

	return category, nil
}
