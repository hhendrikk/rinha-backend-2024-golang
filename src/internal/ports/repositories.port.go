package ports

import (
	"context"
	"database/sql"

	"rinha.backend.2024/src/internal/domains"
	"rinha.backend.2024/src/internal/domains/entities"
)

type Repository interface {
	SetTx(tx *sql.Tx)
}

type ClientRepository interface {
	Get(ctx context.Context, id domains.ID, lock bool) (*entities.ClientEntity, error)
	UpdateAmountTx(ctx context.Context, id domains.ID, amount domains.MoneyCents) error
	Repository
}

type TransactionRepository interface {
	GetLatest(ctx context.Context, clientID domains.ID) ([]entities.TransactionEntity, error)
	SaveTx(ctx context.Context, transaction entities.TransactionEntity) error
	Repository
}
