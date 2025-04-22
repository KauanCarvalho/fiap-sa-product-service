package domain

import (
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"
)

type Datastore interface {
	ports.HealthCheckRepository
	ports.CategoryRepository
	ports.ProductRepository
}
