package adapters

import (
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

func ToWalletBalanceResponse(user *models.User, wallets []models.Wallet) models.WalletBalanceResponse {
	var walletDetails []models.WalletDetail
	totalBalance := decimal.NewFromInt(0)

	for _, w := range wallets {
		walletDetails = append(walletDetails, models.WalletDetail{
			ID:       w.ID,
			Currency: w.Currency,
			Type:     w.Type,
			Balance:  w.Balance,
		})
		totalBalance = totalBalance.Add(w.Balance)
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
