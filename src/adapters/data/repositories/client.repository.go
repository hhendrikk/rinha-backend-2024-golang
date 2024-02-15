package repositories

import (
	"context"
	"database/sql"
	"errors"

	"rinha.backend.2024/src/internal/domains"
	"rinha.backend.2024/src/internal/domains/entities"
	"rinha.backend.2024/src/internal/ports"
	"rinha.backend.2024/src/pkg/database"
	"rinha.backend.2024/src/pkg/validation"
)

type ClientRepository struct {
	Repository
}

func NewClientRepository(db *sql.DB) ports.ClientRepository {
	return &ClientRepository{
		Repository: Repository{
			db: db,
		},
	}
}

func (r *ClientRepository) Get(ctx context.Context, id domains.ID) (*entities.ClientEntity, error) {
	query := `
		SELECT
			 id
			,nome
			,limite
			,saldo
		FROM
			public.clientes
		WHERE
			id = $1
	`

	var client entities.ClientEntity

	err := r.getConn().
		QueryRowContext(ctx, query, id).
		Scan(&client.ID,
			&client.Name,
			&client.Limit,
			&client.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, validation.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) UpdateAmountTx(ctx context.Context, id domains.ID, amount domains.MoneyCents) error {
	if r.tx == nil {
		return database.ErrTransactionNotStarted
	}

	query := `
		UPDATE
			public.clientes
		SET
			saldo = $1
		WHERE
			id = $2;
	`

	_, err := r.getConn().ExecContext(ctx, query, amount, id)
	if err != nil {
		return err
	}

	return err
}
