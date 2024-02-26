package httpserver_gin

import (
	"fmt"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"rinha.backend.2024/src/adapters/config"
)

type GroupFunc func(app *gin.Engine)

func NewHttpServer(config *config.RinhaBackendConfig, fn GroupFunc) {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()

	if config.Environment == "development" {
		app.Use(logger.SetLogger())
	}

	app.Use(gin.Recovery())
	app.Use(gzip.Gzip(gzip.BestCompression))

	fn(app)

	log.Fatal(app.Run(fmt.Sprintf(":%d", config.Port)))
}
