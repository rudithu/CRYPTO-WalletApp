package adapters

import "github.com/rudithu/CRYPTO-WalletApp/models"

func ToWalletTransactionResponse(user *models.User, wallet models.Wallet, txns []models.Transaction) models.WalletTransactionResponse {
	var txnItems []models.TransactionSummaryItem

	for _, t := range txns {
		var counterparty *int64
		if t.CounterpartyWalletId.Valid {
			id := t.CounterpartyWalletId.Int64
			counterparty = &id
		}

		txnItems = append(txnItems, models.TransactionSummaryItem{
			ID:                   t.ID,
			Type:                 t.Type,
			Amount:               t.Amount.InexactFloat64(),
			CounterpartyWalletID: counterparty,
			Time:                 t.CreatedAt,
		})
	}

	return models.WalletTransactionResponse{
		UserInfo: models.UserInfo{
			ID:   user.ID,
			Name: user.Name,
		},
		WalletInfo: models.WalletDetail{
			ID:       wallet.ID,
			Currency: wallet.Currency,
			Type:     wallet.Type,
			Balance:  wallet.Balance,
		},
		Transactions: txnItems,
	}
}
