package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                   int64           `json:"id"`
	WalletId             int64           `json:"wallet_id"`
	Type                 string          `json:"type"`
	Amount               decimal.Decimal `json:"amount"`
	CounterpartyWalletId sql.NullInt64   `json:"counterparty_wallet_id"`
	CreatedAt            time.Time       `json:"created_at"`
}
