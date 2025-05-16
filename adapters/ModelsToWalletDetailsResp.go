package adapters

import (
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

func ToWalletDetailsResp(user *models.User, wallets []models.Wallet, txns []models.Transaction) models.WalletBalanceResponse {
	var walletDetails []models.WalletDetail
	totalBalance := decimal.NewFromInt(0)

	grouped := make(map[int64][]models.TransactionSummaryItem)

	if txns != nil {
		for _, tx := range txns {

			var counterId *int64
			if tx.CounterpartyWalletId.Valid {
				counterId = &tx.CounterpartyWalletId.Int64
			} else {
				counterId = nil
			}

			grouped[tx.WalletId] = append(grouped[tx.WalletId], models.TransactionSummaryItem{
				ID:                   tx.ID,
				Type:                 tx.Type,
				Amount:               tx.Amount,
				Time:                 tx.CreatedAt,
				CounterpartyWalletID: counterId,
			})
		}
	}

	for _, w := range wallets {

		walletDetails = append(walletDetails, models.WalletDetail{
			ID:           w.ID,
			IsDefault:    w.IsDefault,
			Currency:     w.Currency,
			Type:         w.Type,
			Balance:      w.Balance,
			Transactions: grouped[w.ID],
		})
		totalBalance = totalBalance.Add(w.Balance)
	}

	if walletDetails == nil {
		walletDetails = make([]models.WalletDetail, 0)
	}

	return models.WalletBalanceResponse{
		UserInfo: models.UserInfo{
			ID:   user.ID,
			Name: user.Name,
		},
		Wallets:      walletDetails,
		TotalBalance: totalBalance,
	}
}
