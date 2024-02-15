package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"rinha.backend.2024/src/adapters/config"
)

type GroupFunc func(app *echo.Echo)

func NewHttpServer(config *config.RinhaBackendConfig, fn GroupFunc) {
	app := echo.New()

	app.Use(middleware.Recover())
	app.Use(middleware.Gzip())

	if config.Environment == "development" {
		app.Use(middleware.Logger())
	}

	fn(app)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", config.Port)))
}
