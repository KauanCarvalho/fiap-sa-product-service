package api

import (
	"fmt"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/api/handler"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/api/middleware"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	"go.uber.org/zap"

	docs "github.com/KauanCarvalho/fiap-sa-product-service/swagger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Run() {
	// Stores.
	ds := newStores(s.db)

	// Usecases.
	getProduct := usecase.NewGetProductUseCase(ds)
	getProducts := usecase.NewGetProductsUseCase(ds)
	createProduct := usecase.NewCreateProductUseCase(ds)
	updateProduct := usecase.NewUpdateProductUseCase(ds)
	deleteProduct := usecase.NewDeleteProductUseCase(ds)

	// Web server.
	r := GenerateRouter(s.cfg, ds, getProduct, getProducts, createProduct, updateProduct, deleteProduct)

	err := r.Run(fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		zap.L().Fatal("Failed to start server",
			zap.String("port", s.cfg.Port),
			zap.Error(err),
		)
	}
}

func newStores(db *gorm.DB) domain.Datastore {
	return datastore.NewDatastore(db)
}

func GenerateRouter(
	cfg *config.Config,
	ds domain.Datastore,
	gp usecase.GetProductUseCase,
	gps usecase.GetProductsUseCase,
	cp usecase.CreateProductUseCase,
	up usecase.UpdateProductUseCase,
	dp usecase.DeleteProductUseCase,
) *gin.Engine {
	r := gin.New()

	setupMiddlewares(r, cfg)
	registerRoutes(r, cfg, ds, gp, gps, cp, up, dp)

	return r
}

func setupMiddlewares(r *gin.Engine, cfg *config.Config) {
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	r.Use(requestid.New(
		requestid.WithGenerator(func() string {
			return cfg.AppName + "-" + uuid.New().String()
		}),
	))
}

func registerRoutes(
	r *gin.Engine,
	cfg *config.Config,
	ds domain.Datastore,
	gp usecase.GetProductUseCase,
	gps usecase.GetProductsUseCase,
	cp usecase.CreateProductUseCase,
	up usecase.UpdateProductUseCase,
	dp usecase.DeleteProductUseCase,
) {
	healthCheckHandler := handler.NewHealthCheckHandler(ds)
	productHandler := handler.NewProductHandler(gp, gps)
	adminProductHandler := handler.NewProductAdminHandler(cp, up, dp)

	r.GET("/healthcheck", healthCheckHandler.Ping)

	apiV1 := r.Group("/api/v1")
	{
		products := apiV1.Group("/products")
		{
			products.GET("/:sku", productHandler.GetProduct)
			products.GET("/", productHandler.GetProducts)
		}

		admin := apiV1.Group("/admin")
		{
			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("/", adminProductHandler.Create)
				adminProducts.PUT("/:sku", adminProductHandler.Update)
				adminProducts.DELETE("/:sku", adminProductHandler.Delete)
			}
		}
	}

	if cfg.IsDevelopment() {
		docs.SwaggerInfo.BasePath = ""
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
