package requests

import "rinha.backend.2024/src/internal/domains"

type GetLatestStatementBalanceRequestDto struct {
	ClientID domains.ID `param:"id"`
}
