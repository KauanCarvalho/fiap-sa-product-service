package bdd_test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/di"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-testfixtures/testfixtures/v3"
	"go.uber.org/zap"
)

var (
	engine   *gin.Engine
	cfg      *config.Config
	recorder *httptest.ResponseRecorder
	sqlDB    *gorm.DB
	ds       domain.Datastore
	cp       usecase.CreateProductUseCase
	dp       usecase.DeleteProductUseCase
	gp       usecase.GetProductUseCase
	gps      usecase.GetProductsUseCase
	up       usecase.UpdateProductUseCase
	fixtures *testfixtures.Loader
	bodyData map[string]string
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	cfg = config.Load()

	var err error
	sqlDB, err = di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	db, dbErr := sqlDB.DB()
	if dbErr != nil {
		log.Fatalf("error when getting database connection: %v", dbErr)
	}

	fixtures, err = di.SetupFixtures(db, "../../../testdata/fixtures")
	if err != nil {
		log.Fatalf("error when initializing fixtures: %v", err)
	}

	ds = datastore.NewDatastore(sqlDB)
	cp = usecase.NewCreateProductUseCase(ds)
	dp = usecase.NewDeleteProductUseCase(ds)
	gp = usecase.NewGetProductUseCase(ds)
	gps = usecase.NewGetProductsUseCase(ds)
	up = usecase.NewUpdateProductUseCase(ds)

	engine = api.GenerateRouter(cfg, ds, gp, gps, cp, up, dp)

	code := m.Run()

	os.Exit(code)
}

func resetState() {
	recorder = httptest.NewRecorder()
	bodyData = make(map[string]string)
}

func loadTestFixtures() {
	if err := fixtures.Load(); err != nil {
		log.Fatalf("error when loading fixtures: %v", err)
	}
}

func ResetAndLoadFixtures() {
	resetState()
	loadTestFixtures()
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "Features",
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenarioProductAPI(sc)
			InitializeScenarioProductAdmin(sc)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"./features"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("tests failed")
	}
}
