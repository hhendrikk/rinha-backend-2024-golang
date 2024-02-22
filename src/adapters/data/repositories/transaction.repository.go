package repositories

import (
	"context"
	"database/sql"

	"rinha.backend.2024/src/internal/domains"
	"rinha.backend.2024/src/internal/domains/entities"
	"rinha.backend.2024/src/internal/ports"
	"rinha.backend.2024/src/pkg/database"
)

type TransactionRepository struct {
	Repository
}

func NewTransactionRepository(db *sql.DB) ports.TransactionRepository {
	return &TransactionRepository{
		Repository: Repository{
			db: db,
		},
	}
}

func (t *TransactionRepository) GetLatest(ctx context.Context, clientID domains.ID) ([]entities.TransactionEntity, error) {
	const LastCount = 10
	query := `
		SELECT
			id
			,tipo
			,descricao
			,valor
			,cliente_id
			,realizada_em
		FROM
			public.transacoes
		WHERE
			cliente_id = $1
		ORDER BY realizada_em DESC
		LIMIT $2
	`

	transactions := make([]entities.TransactionEntity, 0, LastCount)

	rows, err := t.getConn().QueryContext(ctx, query, clientID, LastCount)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction entities.TransactionEntity
		err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.Description,
			&transaction.Value,
			&transaction.ClientID,
			&transaction.RealizedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (t *TransactionRepository) SaveTx(ctx context.Context, transaction entities.TransactionEntity) error {
	if t.tx == nil {
		return database.ErrTransactionNotStarted
	}

	query := `
		INSERT INTO public.transacoes (
			tipo,
			descricao,
			valor,
			cliente_id,
			realizada_em
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`

	_, err := t.getConn().ExecContext(ctx, query, transaction.Type, transaction.Description, transaction.Value, transaction.ClientID, transaction.RealizedAt)

	return err
}
