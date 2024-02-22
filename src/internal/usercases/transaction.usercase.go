package usercases

import (
	"context"
	"errors"
	"time"

	"rinha.backend.2024/src/internal/domains"
	"rinha.backend.2024/src/internal/domains/dtos/requests"
	"rinha.backend.2024/src/internal/domains/dtos/responses"
	"rinha.backend.2024/src/internal/domains/entities"
	"rinha.backend.2024/src/internal/ports"
	"rinha.backend.2024/src/pkg/database"
	"rinha.backend.2024/src/pkg/validation"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrClientNotFound      = errors.New("client not found")
)

type TransactionTime func() time.Time

type TransactionUsercase struct {
	clientRepository      ports.ClientRepository
	transactionRepository ports.TransactionRepository
	transactionTime       TransactionTime
	uow                   database.Worker
}

func NewTransactionUsercase(clientRepository ports.ClientRepository, transactionRepository ports.TransactionRepository, uow database.Worker, transactionTime TransactionTime) *TransactionUsercase {
	return &TransactionUsercase{
		clientRepository:      clientRepository,
		transactionRepository: transactionRepository,
		transactionTime:       transactionTime,
		uow:                   uow,
	}
}

func (u *TransactionUsercase) Execute(ctx context.Context, request requests.RegisterTransactionRequestDto) (*responses.TransactionResponseDto, error) {
	err := request.Validate()
	if err != nil {
		return nil, errors.Join(err, validation.ErrValidation)
	}

	var response responses.TransactionResponseDto

	err = u.uow.Do(
		ctx,
		func(uow database.Worker) error {
			u.clientRepository.SetTx(uow.GetTx())
			u.transactionRepository.SetTx(uow.GetTx())

			client, err := u.clientRepository.Get(ctx, request.ClientID, true)
			if err != nil {
				if errors.Is(err, validation.ErrNotFound) {
					return errors.Join(validation.ErrNotFound, ErrClientNotFound)
				}
				return err
			}

			balance, transaction, err := u.processTransaction(*client, request)
			if err != nil {
				return err
			}

			err = u.transactionRepository.SaveTx(ctx, *transaction)
			if err != nil {
				return err
			}

			err = u.clientRepository.UpdateAmountTx(ctx, balance.ID, balance.Amount)
			if err != nil {
				return err
			}

			response = responses.NewTransactionResponse(balance.Limit, balance.Amount)

			return nil
		})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *TransactionUsercase) processTransaction(client entities.ClientEntity, request requests.RegisterTransactionRequestDto) (*entities.ClientEntity, *entities.TransactionEntity, error) {
	var transactionType domains.TransactionType

	switch domains.TransactionType(request.Type) {
	case domains.CreditTransaction:
		client.Amount += domains.MoneyCents(request.Value)
		transactionType = domains.CreditTransaction

	case domains.DebitTransaction:
		client.Amount -= domains.MoneyCents(request.Value)
		transactionType = domains.DebitTransaction
		isInsufficientBalance := client.Amount < client.Limit*-1

		if isInsufficientBalance {
			return nil, nil, errors.Join(validation.ErrValidation, ErrInsufficientBalance)
		}
	}

	transaction := entities.NewTransactionEntity(
		0,
		client.ID,
		domains.MoneyCents(request.Value),
		transactionType,
		request.Description,
		u.transactionTime(),
	)

	return &client, &transaction, nil
}
