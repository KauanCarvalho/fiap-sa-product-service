package handler

import (
	"errors"
	"net/http"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/dto"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/shared/validation"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ProductAdminHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

type productAdminHandler struct {
	createProductUseCase usecase.CreateProductUseCase
	updateProductUseCase usecase.UpdateProductUseCase
	deleteProductUseCase usecase.DeleteProductUseCase
}

func NewProductAdminHandler(
	createProductUseCase usecase.CreateProductUseCase,
	updateProductUseCase usecase.UpdateProductUseCase,
	deleteProductUseCase usecase.DeleteProductUseCase,
) ProductAdminHandler {
	return &productAdminHandler{
		createProductUseCase: createProductUseCase,
		updateProductUseCase: updateProductUseCase,
		deleteProductUseCase: deleteProductUseCase,
	}
}

// Create product.
// @Summary	    Create product.
// @Description Create product.
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       product body useCaseDTO.ProductInputCreate true "request body"
// @Success     200 {object} dto.ProductOutput
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     404 {object} dto.APIErrorsOutput
// @Failure     409 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/admin/products/ [post].
func (h *productAdminHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input useCaseDTO.ProductInputCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	if err := useCaseDTO.ValidateProductCreate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	product, err := h.createProductUseCase.Run(ctx, input)
	if err != nil {
		if errors.Is(err, datastore.ErrExistingRecord) {
			c.JSON(http.StatusConflict, dto.SimpleAPIErrorsOutput(
				"",
				"sku",
				"SKU already exists",
			))
			return
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.SimpleAPIErrorsOutput("", "", "resource not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput("", "", "failed to create product"))
		return
	}

	c.JSON(http.StatusCreated, mappers.ToProductDTO(*product))
}

// Update product.
// @Summary	    Update product.
// @Description Update product.
// @Tags        Product
// @Accept      json
// @Produce     json
// @Param       product body useCaseDTO.ProductInputUpdate true "request body"
// @Success     200 {object} dto.ProductOutput
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     404 {object} dto.APIErrorsOutput
// @Failure     409 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/admin/products/ [put].
func (h *productAdminHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	sku := c.Param("sku")

	var input useCaseDTO.ProductInputUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	input.SKU = sku

	if err := useCaseDTO.ValidateProductUpdate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	product, err := h.updateProductUseCase.Run(ctx, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.SimpleAPIErrorsOutput("", "", "resource not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Failed to update product",
		))
		return
	}

	c.JSON(http.StatusOK, mappers.ToProductDTO(*product))
}

// Delete product.
// @Summary	    Delete product.
// @Description Delete product.
// @Tags        Product
// @Accept      json
// @Produce     json
// @Success     204 "No Content"
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/admin/products/ [delete].
func (h *productAdminHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	sku := c.Param("sku")

	err := h.deleteProductUseCase.Run(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Failed to delete product",
		))
		return
	}

	c.Status(http.StatusNoContent)
}
