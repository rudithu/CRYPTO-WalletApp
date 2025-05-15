package adapters

import (
	"database/sql"
	"time"

	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

// Converts a TransactionRequest into a Transaction model
func (r *models.TransactionRequest) ToModel() *models.Transaction {
	var counterpartyWalletId sql.NullInt64
	if r.DestinationWalletID != nil {
		counterpartyWalletId = sql.NullInt64{
			Int64: *r.DestinationWalletID,
			Valid: true,
		}
	} else {
		counterpartyWalletId = sql.NullInt64{Valid: false}
	}

	return &models.Transaction{
		WalletId:             r.WalletID,
		Type:                 r.TransactionType,
		Amount:               decimal.NewFromFloat(r.Amount),
		CounterpartyWalletId: counterpartyWalletId,
		CreatedAt:            time.Now(),
	}
}
