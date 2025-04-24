package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	"gorm.io/gorm"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/ports"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetProduct(c *gin.Context)
	GetProducts(c *gin.Context)
}

type productHandler struct {
	getProductUseCase  usecase.GetProductUseCase
	getProductsUseCase usecase.GetProductsUseCase
}

func NewProductHandler(
	getProductUseCase usecase.GetProductUseCase,
	getProductsUseCase usecase.GetProductsUseCase,
) ProductHandler {
	return &productHandler{
		getProductUseCase:  getProductUseCase,
		getProductsUseCase: getProductsUseCase,
	}
}

// Get product by SKU.
// @Summary	    Get product by SKU.
// @Description Get product by SKU.
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       sku path string true "product sku"
// @Success     200 {object} dto.ProductOutput
// @Failure     404 "No Content"
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/products/{sku} [get].
func (h *productHandler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	sku := c.Param("sku")

	product, err := h.getProductUseCase.Run(ctx, sku)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
				"",
				"",
				"Failed to get products",
			))
		}
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, mappers.ToProductDTO(*product))
}

// Get products by filter.
// @Summary	    Get products by filter.
// @Description Get products by filter.
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       page query string false "Current page"
// @Param       pageSize query string false "Page size"
// @Param       category query string false "Category name"
// @Success     200 {object} dto.ProductsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/products/ [get].
func (h *productHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"Invalid page number",
			"page",
			c.Request.URL.Query().Get("page"),
		))
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"Invalid page size",
			"pageSize",
			c.Request.URL.Query().Get("pageSize"),
		))
		return
	}

	category := c.DefaultQuery("category", "")

	products, total, err := h.getProductsUseCase.Run(ctx, &ports.ProductFilter{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Failed to get products",
		))
		return
	}

	c.JSON(http.StatusOK, dto.ProductsOutput{
		Products:    mappers.ToProductsDTO(products),
		PageSize:    pageSize,
		CurrentPage: page,
		Total:       total,
	})
}
