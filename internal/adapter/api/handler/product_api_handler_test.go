package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductHandler_GetProduct(t *testing.T) {
	prepareTestDatabase()

	t.Run("should return product successfully", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/hamburger", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})

	t.Run("should return 404 when product not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/invalid-sku", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProductHandler_GetProducts(t *testing.T) {
	prepareTestDatabase()

	t.Run("when param page is invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/?page=abc", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})

	t.Run("when param page is invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/?pageSize=abc", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})

	t.Run("successfully get products", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})
}
