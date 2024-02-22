package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"rinha.backend.2024/src/adapters/config"
	httpserver_echo "rinha.backend.2024/src/adapters/httpserver/echo"
	"rinha.backend.2024/src/pkg/database"
)

func main() {
	rinhaConfig := config.NewRinhaBackendConfig()
	connectionString := database.BuildConnectionString(
		"postgres",
		rinhaConfig.Database.Host,
		rinhaConfig.Database.Port,
		rinhaConfig.Database.User,
		rinhaConfig.Database.Pass,
		rinhaConfig.Database.DbName,
		map[string]string{
			"application_name": rinhaConfig.Database.Application,
		},
	)
	dataBase := database.NewPostgresDatabase(connectionString, database.ConfigDatabaseFunc(func(db *sql.DB) {}))

	httpserver_echo.NewHttpServer(rinhaConfig, func(app *echo.Echo) {
		httpserver_echo.NewClientGroup(app, dataBase, rinhaConfig)
	})
}
