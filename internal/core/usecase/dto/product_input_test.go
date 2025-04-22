package dto_test

import (
	"strings"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProductCreate(t *testing.T) {
	t.Run("invalid create - missing fields", func(t *testing.T) {
		input := dto.ProductInputCreate{}

		err := dto.ValidateProductCreate(input)
		require.Error(t, err)

		validationErrors, ok := err.(validator.ValidationErrors)
		require.True(t, ok)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Name")
		assert.Contains(t, fieldNames, "Price")
		assert.Contains(t, fieldNames, "Description")
		assert.Contains(t, fieldNames, "Category.Name")
		assert.Contains(t, fieldNames, "Images")
	})

	t.Run("invalid create - name and description", func(t *testing.T) {
		input := dto.ProductInputCreate{
			Name:        "",
			Price:       10,
			Description: "",
			Category:    dto.ProductCategoryCreate{Name: ""},
			Images:      []dto.ProductImageCreate{{URL: "https://example.com/image.jpg"}},
		}

		err := dto.ValidateProductCreate(input)
		require.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Name")
		assert.Contains(t, fieldNames, "Description")
		assert.Contains(t, fieldNames, "Category.Name")
	})

	t.Run("invalid create - invalid URL", func(t *testing.T) {
		input := dto.ProductInputCreate{
			Name:        "Valid Name",
			Price:       100.5,
			Description: "A valid description here",
			Category:    dto.ProductCategoryCreate{Name: "Tech"},
			Images:      []dto.ProductImageCreate{{URL: "invalid-url"}}, // URL inválida
		}

		err := dto.ValidateProductCreate(input)
		require.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Images[0].URL")
	})
}

func TestValidateProductUpdate(t *testing.T) {
	t.Run("invalid update - missing fields", func(t *testing.T) {
		input := dto.ProductInputUpdate{}

		err := dto.ValidateProductUpdate(input)
		require.Error(t, err)

		validationErrors, ok := err.(validator.ValidationErrors)
		require.True(t, ok)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "SKU")
		assert.Contains(t, fieldNames, "Name")
		assert.Contains(t, fieldNames, "Price")
		assert.Contains(t, fieldNames, "Description")
		assert.Contains(t, fieldNames, "Category.Name")
		assert.Contains(t, fieldNames, "Images")
	})

	t.Run("invalid update - name and description", func(t *testing.T) {
		input := dto.ProductInputUpdate{
			SKU:         "",
			Name:        "",
			Description: "",
			Price:       10,
			Category:    dto.ProductCategoryUpdate{Name: ""},
			Images:      []dto.ProductImageUpdate{{URL: "https://example.com/image.jpg"}},
		}

		err := dto.ValidateProductUpdate(input)
		require.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Name")
		assert.Contains(t, fieldNames, "Description")
		assert.Contains(t, fieldNames, "Category.Name")
	})

	t.Run("invalid update - invalid URL", func(t *testing.T) {
		input := dto.ProductInputUpdate{
			SKU:         "sku-xyz",
			Name:        "Valid Name",
			Description: "A valid description here",
			Price:       100.5,
			Category:    dto.ProductCategoryUpdate{Name: "Tech"},
			Images:      []dto.ProductImageUpdate{{URL: "invalid-url"}},
		}

		err := dto.ValidateProductUpdate(input)
		require.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Images[0].URL")
	})
}
