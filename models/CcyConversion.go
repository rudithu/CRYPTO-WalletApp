package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type CcyConversion struct {
	FromCcy   string          `json:"from_ccy"`
	ToCcy     string          `json:"to_ccy"`
	Rate      decimal.Decimal `json:"rate"`
	CreatedAt time.Time       `json:"created_at"`
}

type CcyRateToBaseCcy struct {
	Ccy  string
	Rate decimal.Decimal
}
