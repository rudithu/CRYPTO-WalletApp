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
	IsDefault    bool                     `json:"is_default"`
	Type         string                   `json:"type"`
	Currency     string                   `json:"currency"`
	Balance      MoneyDecimal             `json:"balance"`
	Transactions []TransactionSummaryItem `json:"transactions,omitempty"`
}

type TransactionSummaryItem struct {
	ID                   int64        `json:"id"`
	Type                 string       `json:"type"`
	Amount               MoneyDecimal `json:"amount"`
	CounterpartyWalletID *int64       `json:"counterparty_wallet_id,omitempty"`
	Time                 time.Time    `json:"time"`
}

type Total struct {
	Currency string       `json:"currency"`
	Amount   MoneyDecimal `json:"amount"`
}

// Req: userID
type WalletBalanceResponse struct {
	UserInfo UserInfo       `json:"user_info"`
	Wallets  []WalletDetail `json:"wallets"`
	Balance  *Total         `json:"total,omitempty"`
}

func (tr *TransactionRequest) ValidateRequest(txnType string) error {
	if tr.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount field is mandatory and it must be greater than zero")
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

// Custom type that wraps decimal.Decimal
type MoneyDecimal struct {
	decimal.Decimal
}

// MarshalJSON limits to 2 decimal places when marshaling to JSON
func (d MoneyDecimal) MarshalJSON() ([]byte, error) {
	// Format to 2 decimal places and quote it as a string
	return []byte(fmt.Sprintf("\"%s\"", d.Decimal.StringFixed(2))), nil
}
