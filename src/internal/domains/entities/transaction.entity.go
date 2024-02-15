package entities

import (
	"time"

	"rinha.backend.2024/src/internal/domains"
)

const TransactionTableName = "public.transacoes"

type TransactionEntity struct {
	ID          domains.ID              `json:"id"`
	ClientID    domains.ID              `json:"client_id"`
	Value       domains.MoneyCents      `json:"value"`
	Type        domains.TransactionType `json:"type"`
	Description string                  `json:"description"`
	RealizedAt  time.Time               `json:"realized_at"`
}

func NewTransactionEntity(
	id domains.ID,
	clientID domains.ID,
	value domains.MoneyCents,
	transactionType domains.TransactionType,
	description string,
	realizedAt time.Time,
) TransactionEntity {
	return TransactionEntity{
		ID:          id,
		ClientID:    clientID,
		Value:       value,
		Type:        transactionType,
		Description: description,
		RealizedAt:  realizedAt,
	}
}
