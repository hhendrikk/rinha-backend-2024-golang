package responses

import "rinha.backend.2024/src/internal/domains"

type TransactionResponseDto struct {
	Limit   domains.MoneyCents `json:"limite"`
	Balance domains.MoneyCents `json:"saldo"`
}

func NewTransactionResponse(limit, balance domains.MoneyCents) TransactionResponseDto {
	return TransactionResponseDto{
		Limit:   limit,
		Balance: balance,
	}
}
