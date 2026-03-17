package api

import (
	"log/slog"
	"net/http"
	"store-service/internal/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductAPI struct {
	ProductService product.ProductInterface
}

// @Summary Search products
// @Description Search for products with optional filtering
// @Tags product
// @Accept json
// @Produce json
// @Param q query string false "Search keyword"
// @Param limit query string false "Number of items per page" default(30)
// @Param offset query string false "Offset for pagination" default(0)
// @Success 200 {array} product.Product
// @Failure 500
// @Router /api/v1/product [get]
func (api ProductAPI) SearchHandler(context *gin.Context) {
	ctx := context.Request.Context()
	keyword := context.DefaultQuery("q", "")
	limit := context.DefaultQuery("limit", "30")
	offset := context.DefaultQuery("offset", "0")

	productResult, err := api.ProductService.GetProducts(ctx, keyword, limit, offset)

	if err != nil {
		slog.ErrorContext(ctx, "ProductService.GetProducts failed",
			"log_type", "error",
			"error_code", "PRODUCT_SEARCH_FAILED",
			"error_message", err.Error(),
			"user_id", 0,
			slog.Any("request", map[string]any{"keyword": keyword, "limit": limit, "offset": offset}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, productResult)
}

// @Summary Get product by ID
// @Description Get detailed information about a specific product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} product.Product
// @Failure 400 {object} string
// @Failure 500
// @Router /api/v1/product/{id} [get]
func (api ProductAPI) GetProductHandler(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "id is not integer",
		})
		return
	}
	ctx := context.Request.Context()
	product, err := api.ProductService.GetProductByID(ctx, id)

	if err != nil {
		slog.ErrorContext(ctx, "ProductService.GetProductByID failed",
			"log_type", "error",
			"error_code", "PRODUCT_QUERY_FAILED",
			"error_message", err.Error(),
			"user_id", 0,
			slog.Any("request", map[string]any{"product_id": id}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, product)
}
