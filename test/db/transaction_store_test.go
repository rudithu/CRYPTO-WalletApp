package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionsByWalletIDs_Success(t *testing.T) {

	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {

		expectedTxn := testutils.MockTxns()[0]

		rows := sqlmock.NewRows([]string{"id", "wallet_id", "type", "amount", "counterparty_wallet_id", "created_at"}).
			AddRow(expectedTxn.ID, expectedTxn.WalletId, expectedTxn.Type, expectedTxn.Amount, expectedTxn.CounterpartyWalletId, expectedTxn.CreatedAt)

		mock.ExpectQuery("SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at FROM transactions WHERE wallet_id in \\(.+\\) ORDER BY created_at DESC").
			WithArgs(expectedTxn.WalletId).
			WillReturnRows(rows)

		txns, err := db.GetTransactionsByWalletIDs(dbTest, []int64{expectedTxn.WalletId})

		assert.Nil(t, err)
		assert.NotNil(t, txns)
		assert.Equal(t, 1, len(txns))
		assert.Equal(t, expectedTxn.ID, txns[0].ID)

	})

}

func TestGetTransactionsByWalletIDs_NotFound(t *testing.T) {
	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {
		walletId := int64(23)

		rows := sqlmock.NewRows([]string{"id", "wallet_id", "type", "amount", "counterparty_wallet_id", "created_at"})

		mock.ExpectQuery("SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at FROM transactions WHERE wallet_id in \\(.+\\) ORDER BY created_at DESC").
			WithArgs(walletId).
			WillReturnRows(rows)

		txns, err := db.GetTransactionsByWalletIDs(dbTest, []int64{walletId})

		assert.Nil(t, err)
		assert.Nil(t, txns)
	})
}

func TestGetTransactionsByWalletIDs_DBError(t *testing.T) {
	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {
		walletId := int64(23)

		mock.ExpectQuery("SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at FROM transactions WHERE wallet_id in \\(.+\\) ORDER BY created_at DESC").
			WithArgs(walletId).
			WillReturnError(errors.New("db failed"))

		txns, err := db.GetTransactionsByWalletIDs(dbTest, []int64{walletId})

		assert.NotNil(t, err)
		assert.Equal(t, "db failed", err.Error())
		assert.Nil(t, txns)
	})

}
