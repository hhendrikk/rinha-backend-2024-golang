package responses

import "rinha.backend.2024/src/internal/domains"

type StatementBalanceResponseDto struct {
	Balance      *BalanceResponseDto               `json:"saldo"`
	Transactions []StatementTransactionResponseDto `json:"ultimas_transacoes"`
}

type BalanceResponseDto struct {
	Total domains.MoneyCents `json:"total"`
	Date  string             `json:"data_extrato"`
	Limit domains.MoneyCents `json:"limite"`
}

type StatementTransactionResponseDto struct {
	Value       domains.MoneyCents `json:"valor"`
	Type        string             `json:"tipo"`
	Description string             `json:"descricao"`
	CreatedAt   string             `json:"realizada_em"`
}

func NewStatementBalanceResponse(balance *BalanceResponseDto, transactions []StatementTransactionResponseDto) *StatementBalanceResponseDto {
	return &StatementBalanceResponseDto{
		Balance:      balance,
		Transactions: transactions,
	}
}

func NewBalanceResponse(total domains.MoneyCents, date string, limit domains.MoneyCents) *BalanceResponseDto {
	return &BalanceResponseDto{
		Total: total,
		Date:  date,
		Limit: limit,
	}
}

func NewStatementTransactionResponse(value domains.MoneyCents, typeTransaction string, description string, date string) *StatementTransactionResponseDto {
	return &StatementTransactionResponseDto{
		Value:       value,
		Type:        typeTransaction,
		Description: description,
		CreatedAt:   date,
	}
}
