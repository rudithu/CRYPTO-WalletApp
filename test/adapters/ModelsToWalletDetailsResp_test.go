package adapters

import (
	"database/sql"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/rudithu/CRYPTO-WalletApp/adapters"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
)

func TestToWalletDetailsResp(t *testing.T) {

	// Call the function
	resp := adapters.ToWalletDetailsResp(testutils.MockUser(), testutils.MockWallets(), testutils.MockTxns(), testutils.MockCcyMapWithRate())

	// Assertions
	assert.Equal(t, testutils.MockUser().ID, resp.UserInfo.ID)
	assert.Equal(t, 2, len(resp.Wallets))

	assert.NotNil(t, resp.Wallets[0].Transactions)
	assert.Equal(t, "USD", resp.Wallets[0].Currency)
	assert.Equal(t, decimal.NewFromFloat(100.00), resp.Wallets[0].Balance)
	assert.NotNil(t, resp.Balance)

	// Total balance should be USD 100 + (EUR 50 / 0.5) = 100 + 100 = 200
	expectedTotal := decimal.NewFromFloat(200.00)
	assert.True(t, resp.Balance.Amount.Equal(expectedTotal))
	assert.Equal(t, models.BaseCcy, resp.Balance.Currency)
}

func TestToWalletDetailsResp_NoTotalBalance(t *testing.T) {
	// Call the function
	resp := adapters.ToWalletDetailsResp(testutils.MockUser(), testutils.MockWallets(), testutils.MockTxns(), map[string]models.CcyRateToBaseCcy{})

	// Assertions
	assert.Equal(t, testutils.MockUser().ID, resp.UserInfo.ID)
	assert.Equal(t, 2, len(resp.Wallets))

	assert.Nil(t, resp.Wallets[1].Transactions)
	assert.Equal(t, "EUR", resp.Wallets[1].Currency)
	assert.Equal(t, decimal.NewFromFloat(50.00), resp.Wallets[1].Balance)
	assert.Nil(t, resp.Balance)
}

func NullInt64(val int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{
		Int64: val,
		Valid: valid,
	}
}
