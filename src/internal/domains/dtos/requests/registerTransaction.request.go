package requests

import (
	"errors"
	"strings"

	"rinha.backend.2024/src/internal/domains"
)

var (
	ErrRequiredClientID       = errors.New("client id cannot be 0")
	ErrRequiredValue          = errors.New("value cannot be 0")
	ErrRequiredDescription    = errors.New("description cannot be empty")
	ErrLengthDescription      = errors.New("description cannot be more than 10 characters")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

type RegisterTransactionRequestDto struct {
	ClientID    domains.ID `param:"id", json:"id"`
	Value       uint64     `json:"valor"`
	Type        string     `json:"tipo"`
	Description string     `json:"descricao"`
}

func (r *RegisterTransactionRequestDto) Validate() error {
	if r.ClientID == 0 {
		return ErrRequiredClientID
	}

	if r.Value == 0 {
		return ErrRequiredValue
	}

	if r.Type != string(domains.CreditTransaction) && r.Type != string(domains.DebitTransaction) {
		return ErrInvalidTransactionType
	}

	if strings.TrimSpace(r.Description) == "" {
		return ErrRequiredDescription
	}

	if len(r.Description) > 10 {
		return ErrLengthDescription
	}

	return nil
}
