package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"rinha.backend.2024/src/adapters/config"
	httpserver_gin "rinha.backend.2024/src/adapters/httpserver/gin"
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
	dataBase := database.NewPostgresDatabase(connectionString, database.ConfigDatabaseFunc(func(db *sql.DB) {
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}))

	httpserver_gin.NewHttpServer(rinhaConfig, func(app *gin.Engine) {
		httpserver_gin.NewClientGroup(app, dataBase, rinhaConfig)
	})
}
