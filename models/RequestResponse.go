package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	Amount              decimal.Decimal `json:"amount"`
	DestinationWalletID *int64          `json:"destination_wallet_id,omitempty"`
}

type UserInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WalletDetail struct {
	ID           int64                    `json:"id"`
	Type         string                   `json:"type"`
	Currency     string                   `json:"currency"`
	Balance      decimal.Decimal          `json:"balance"`
	Transactions []TransactionSummaryItem `json:"transactions,omitempty"`
}

type TransactionSummaryItem struct {
	ID                   int64           `json:"id"`
	Type                 string          `json:"type"`
	Amount               decimal.Decimal `json:"amount"`
	CounterpartyWalletID *int64          `json:"counterparty_wallet_id,omitempty"`
	Time                 time.Time       `json:"time"`
}

// Req: userID
type WalletBalanceResponse struct {
	UserInfo     UserInfo        `json:"user_info"`
	Wallets      []WalletDetail  `json:"wallets"`
	TotalBalance decimal.Decimal `json:"total_balance"`
}
