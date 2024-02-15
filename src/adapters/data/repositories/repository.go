package repositories

import (
	"database/sql"

	"rinha.backend.2024/src/pkg/database"
)

type Repository struct {
	db *sql.DB
	tx *sql.Tx
}

func (r *Repository) getConn() database.Stmt {
	if r.tx == nil {
		return r.db
	}
	return r.tx
}

func (r *Repository) SetTx(tx *sql.Tx) {
	r.tx = tx
}
