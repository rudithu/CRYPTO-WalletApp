package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	TransactionType     string  `json:"transaction_type"`
	Amount              float64 `json:"amount"`
	WalletID            int64   `json:"wallet_id"`
	DestinationWalletID *int64  `json:"destination_wallet_id,omitempty"`
}

type UserInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WalletDetail struct {
	ID       int64           `json:"id"`
	Currency string          `json:"currency"`
	Type     string          `json:"type"`
	Balance  decimal.Decimal `json:"balance`
}

type TransactionSummaryItem struct {
	ID                   int64     `json:"id"`
	Type                 string    `json:"type"`
	Amount               float64   `json:"amount"`
	CounterpartyWalletID *int64    `json:"counterparty_wallet_id,omitempty"`
	Time                 time.Time `json:"time"`
}

// Req: userID
type WalletBalanceResponse struct {
	UserInfo     UserInfo        `json:"user_info"`
	Wallets      []WalletDetail  `json:"wallets`
	TotalBalance decimal.Decimal `json:"total_balance`
}

// Req userID, wallet_type
type WalletTransactionResponse struct {
	UserInfo     UserInfo                 `json:"user_info"`
	WalletInfo   WalletDetail             `json:"wallet_info`
	Transactions []TransactionSummaryItem `json:"transactions"`
}

func validateTransactionRequest(req TransactionRequest) error {
	validTypes := map[string]bool{
		"withdraw":     true,
		"deposit":      true,
		"transfer-out": true,
	}

	if !validTypes[req.TransactionType] {
		return errors.New("invalid transaction_type: must be 'withdraw', 'deposit', or 'transfer-out'")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if req.WalletID <= 0 {
		return errors.New("wallet_id must be a positive number")
	}

	if req.TransactionType == "transfer-out" && req.DestinationWalletID == nil {
		return errors.New("destination_wallet_id is required for transfer-out")
	}

	return nil
}
