package handler_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/di"

	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ctx       context.Context
	cfg       *config.Config
	sqlDB     *gorm.DB
	ds        domain.Datastore
	gp        usecase.GetProductUseCase
	gps       usecase.GetProductsUseCase
	cp        usecase.CreateProductUseCase
	up        usecase.UpdateProductUseCase
	dp        usecase.DeleteProductUseCase
	fixtures  *testfixtures.Loader
	ginEngine *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	ctx = context.Background()
	cfg = config.Load()

	var err error
	sqlDB, err = di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	ds = datastore.NewDatastore(sqlDB)
	gp = usecase.NewGetProductUseCase(ds)
	gps = usecase.NewGetProductsUseCase(ds)
	cp = usecase.NewCreateProductUseCase(ds)
	up = usecase.NewUpdateProductUseCase(ds)
	dp = usecase.NewDeleteProductUseCase(ds)

	ginEngine = api.GenerateRouter(cfg, ds, gp, gps, cp, up, dp)

	db, dbErr := sqlDB.DB()
	if dbErr != nil {
		log.Fatalf("error when getting database connection: %v", dbErr)
	}

	fixtures, err = di.SetupFixtures(db, "../../../../testdata/fixtures")
	if err != nil {
		log.Fatalf("error when initializing fixtures: %v", err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatalf("error when loading fixtures: %v", err)
	}
}
