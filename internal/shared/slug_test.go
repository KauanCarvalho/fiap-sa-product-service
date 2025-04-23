package shared_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/shared"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	t.Run("converts string with spaces and special characters", func(t *testing.T) {
		input := "Hello World! This is a Test."
		expected := "hello-world-this-is-a-test"
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})

	t.Run("handles string with multiple spaces", func(t *testing.T) {
		input := "Hello    World"
		expected := "hello-world"
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})

	t.Run("converts to lowercase", func(t *testing.T) {
		input := "HELLO World"
		expected := "hello-world"
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})

	t.Run("removes leading and trailing hyphens", func(t *testing.T) {
		input := "  Hello World!   "
		expected := "hello-world"
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})

	t.Run("handles strings with only special characters", func(t *testing.T) {
		input := "!@#$%^&*()"
		expected := ""
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})

	t.Run("returns the same result for an already sluggified string", func(t *testing.T) {
		input := "already-slugified"
		expected := "already-slugified"
		result := shared.Slugify(input)

		assert.Equal(t, expected, result)
	})
}
