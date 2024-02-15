package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ConfigDatabaseFunc func(*sql.DB)

type PostgresDatabase struct {
	DB *sql.DB
}

func NewPostgresDatabase(connection string, fnConfig ConfigDatabaseFunc) *PostgresDatabase {
	db, err := sql.Open("pgx", connection)
	if err != nil {
		log.Fatalf("Failed to open database connection: %s \n", err)
	}

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	db.SetConnMaxIdleTime(5 * time.Minute)

	fnConfig(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to open database connection: %s \n", err)
	}

	return &PostgresDatabase{
		DB: db,
	}
}
