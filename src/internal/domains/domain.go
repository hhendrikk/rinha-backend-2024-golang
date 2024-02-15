package domains

type (
	MoneyCents int32
	ID         uint32
)

type TransactionType string

const (
	CreditTransaction TransactionType = "c"
	DebitTransaction  TransactionType = "d"
)
