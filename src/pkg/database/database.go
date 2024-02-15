package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Stmt interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func BuildConnectionString(database, host, port, user, pass, dbname string, parameters map[string]string) string {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?", database, user, pass, host, port, dbname)

	for key, value := range parameters {
		connectionString += fmt.Sprintf("%s=%s&", key, value)
	}

	return connectionString
}
