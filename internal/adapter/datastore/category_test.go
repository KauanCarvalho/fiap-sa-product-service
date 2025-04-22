package datastore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFindCategoryByName(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully find category by name", func(t *testing.T) {
		category, err := ds.FindCategoryByName(ctx, "lanche")
		require.NoError(t, err, "expected no error, got %v", err)
		assert.NotNil(t, category, "expected non-nil category, got nil")
	})

	t.Run("fail to find category by name", func(t *testing.T) {
		category, err := ds.FindCategoryByName(ctx, "NonExistentCategory")
		require.ErrorIs(t, err, gorm.ErrRecordNotFound, "expected record not found error, got %v", err)
		assert.Nil(t, category, "expected nil category, got %v", category)
	})
}
