package usercases

import (
	"context"
	"errors"
	"time"

	"rinha.backend.2024/src/internal/domains/dtos/requests"
	"rinha.backend.2024/src/internal/domains/dtos/responses"
	"rinha.backend.2024/src/internal/ports"
	"rinha.backend.2024/src/pkg/timezone"
	"rinha.backend.2024/src/pkg/validation"
)

type StatementUsercase struct {
	clientRepository      ports.ClientRepository
	transactionRepository ports.TransactionRepository
	timeZone              time.Duration
}

func NewStatementUsercase(clientRepository ports.ClientRepository, transactionRepository ports.TransactionRepository, timeZone time.Duration) *StatementUsercase {
	return &StatementUsercase{
		timeZone:              timeZone,
		clientRepository:      clientRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *StatementUsercase) Execute(ctx context.Context, request requests.GetLatestStatementBalanceRequestDto) (*responses.StatementBalanceResponseDto, error) {
	client, err := s.clientRepository.Get(ctx, request.ClientID)

	if err != nil {
		if errors.Is(err, validation.ErrNotFound) {
			return nil, errors.Join(validation.ErrNotFound, ErrClientNotFound)
		}

		return nil, err
	}

	lastTransactions, err := s.transactionRepository.GetLatest(ctx, client.ID)

	if err != nil {
		return nil, err
	}

	transactions := make([]responses.StatementTransactionResponseDto, 0, len(lastTransactions))
	for _, t := range lastTransactions {
		transactions = append(transactions, *responses.NewStatementTransactionResponse(t.Value,
			string(t.Type),
			t.Description,
			timezone.ToISO8601(t.RealizedAt, s.timeZone)),
		)
	}

	balance := responses.NewBalanceResponse(client.Amount,
		timezone.NowISO8601(s.timeZone),
		client.Limit)

	return responses.NewStatementBalanceResponse(balance, transactions), nil
}
