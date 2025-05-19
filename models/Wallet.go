package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int64           `json:"id"`
	UserId    int64           `json:"user_id"`
	Type      string          `json:"type"`
	IsDefault bool            `json:"is_default"`
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
}
