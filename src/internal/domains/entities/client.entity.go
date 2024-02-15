package entities

import "rinha.backend.2024/src/internal/domains"

const ClientTableName = "public.clientes"

type ClientEntity struct {
	ID     domains.ID         `json:"id"`
	Name   string             `json:"name"`
	Limit  domains.MoneyCents `json:"limit"`
	Amount domains.MoneyCents `json:"amount"`
}

func NewClientEntity(id domains.ID, name string, limit, amount domains.MoneyCents) ClientEntity {
	return ClientEntity{
		ID:     id,
		Name:   name,
		Limit:  limit,
		Amount: amount,
	}
}
