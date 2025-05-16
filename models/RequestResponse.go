package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	Amount              decimal.Decimal `json:"amount"`
	DestinationWalletID *int64          `json:"destination_wallet_id,omitempty"`
	DestinationUserID   *int64          `json:"destination_user_id,omitempty"`
}

type UserInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WalletDetail struct {
	ID           int64                    `json:"id"`
	IsDefault    bool                     `json:"is_default`
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

func (tr *TransactionRequest) ValidateRequest(txnType string) error {
	if tr.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be greater than zero")
	}

	if txnType == TxnTypeTransferOut || txnType == TxnTypeTransferIn {
		if tr.DestinationUserID != nil && tr.DestinationWalletID != nil {
			return fmt.Errorf("please specify only one of destination_user_id or destination_wallet_id, not both")
		}
		if tr.DestinationUserID == nil && tr.DestinationWalletID == nil {
			return fmt.Errorf("please specify either destination_user_id or destination_wallet_id")
		}
	}
	return nil
}
