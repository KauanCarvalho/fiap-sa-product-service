package main

import (
	"log"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-product-service/internal/di"

	"go.uber.org/zap"
)

// @title          Fiap SA Product Service
// @version        0.0.1
// @description    Rest API for Fiap SA Product Service
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host    fa-sa-product-service
// @schemes http
// @accept  json
// @produce json
// main is the entry point of the entire application.
func main() {
	cfg := config.Load()

	log.Println("Starting Fiap SA Product Service...")

	db, err := di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection pool created.")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck // It is not necessary to check for errors at this moment.

	zap.ReplaceGlobals(logger.With(zap.String("app", cfg.AppName), zap.String("env", cfg.AppEnv)))

	apiServer := api.NewServer(cfg, db)
	apiServer.Run()
}
