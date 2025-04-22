package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductAdminHandler_Create(t *testing.T) {
	prepareTestDatabase()

	t.Run("should create product successfully", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Coxinha",
			"price":       7.50,
			"description": "Coxinha de frango",
			"category": map[string]string{
				"name": "lanche",
			},
			"images": []map[string]string{
				{"url": "https://example.com/coxinha.jpg"},
			},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Coxinha")
	})

	t.Run("when category is  not found", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Coxinha",
			"price":       7.50,
			"description": "Coxinha de frango",
			"category": map[string]string{
				"name": "foo",
			},
			"images": []map[string]string{
				{"url": "https://example.com/coxinha.jpg"},
			},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("when body is invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader([]byte("test")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 for missing name", func(t *testing.T) {
		body := map[string]interface{}{
			"price":       7.50,
			"description": "Sem nome",
			"category":    map[string]string{"name": "Test"},
			"images":      []map[string]string{{"url": "https://example.com/img.jpg"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 for invalid image URL", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Produto Inválido",
			"price":       10.0,
			"description": "Imagem inválida",
			"category":    map[string]string{"name": "Categoria"},
			"images":      []map[string]string{{"url": "not-a-url"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 with multiple validation errors", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "",
			"price":       -10,
			"description": "",
			"category":    map[string]string{"name": ""},
			"images":      []map[string]string{{"url": "invalid"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductAdminHandler_Update(t *testing.T) {
	prepareTestDatabase()

	t.Run("when body is invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/products/hamburger", bytes.NewReader([]byte("test")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should update product successfully", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Hambúrguer Atualizado",
			"price":       30.0,
			"description": "Hambúrguer com cheddar",
			"category":    map[string]string{"name": "Lanche"},
			"images":      []map[string]string{{"url": "https://example.com/hamburger.jpg"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/products/hamburger", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Hambúrguer Atualizado")
	})

	t.Run("should return 400 for negative price", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Produto Inválido",
			"price":       -5,
			"description": "Preço negativo",
			"category":    map[string]string{"name": "Categoria"},
			"images":      []map[string]string{{"url": "https://example.com/ok.jpg"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/products/hamburger", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 when updating non-existent product", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Qualquer",
			"price":       20.0,
			"description": "Produto não existe",
			"category":    map[string]string{"name": "Categoria"},
			"images":      []map[string]string{{"url": "https://example.com/img.jpg"}},
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/products/not-found-sku", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProductAdminHandler_Delete(t *testing.T) {
	prepareTestDatabase()

	t.Run("should delete product successfully", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/products/hamburger", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return 500 for non-existing product", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/products/nonexistent", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
