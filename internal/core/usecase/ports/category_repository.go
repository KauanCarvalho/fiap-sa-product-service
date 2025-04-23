package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain/entities"
)

type CategoryRepository interface {
	FindCategoryByName(ctx context.Context, sku string) (*entities.Category, error)
}
