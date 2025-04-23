package datastore

import (
	"errors"
	"time"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"

	"gorm.io/gorm"
)

const DefaultConnectionTimeout = 5 * time.Second

var ErrExistingRecord = errors.New("record already exists")

type datastore struct {
	db *gorm.DB
}

func NewDatastore(db *gorm.DB) domain.Datastore {
	return &datastore{db: db}
}
