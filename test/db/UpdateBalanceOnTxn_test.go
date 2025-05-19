package db_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDepositUpdate_Success(t *testing.T) {

	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {
		counterpartyWalletId := sql.NullInt64{Valid: false}

		txn := &models.Transaction{
			WalletId:             1,
			Type:                 models.TxnTypeDeposit,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: counterpartyWalletId,
		}

		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+").
			WithArgs(txn.Amount, txn.WalletId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := db.DepositUpdate(sqlDB, txn)
		assert.Nil(t, err)
	})
}

func TestDepositUpdate_FailCreateTxn(t *testing.T) {
	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {

		txn := testutils.MockTxns()[0]
		txn.Type = models.TxnTypeDeposit

		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).
			WillReturnError(errors.New("update failed"))

		mock.ExpectRollback()

		err := db.DepositUpdate(sqlDB, &txn)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to create incoming-transaction: update failed", err.Error())

	})
}

func TestWithdrawUpdate_Success(t *testing.T) {

	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {
		counterpartyWalletId := sql.NullInt64{Valid: false}
		txn := &models.Transaction{
			WalletId:             1,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: counterpartyWalletId,
		}

		initialBalance := decimal.NewFromFloat(212.00)

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs(txn.WalletId).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(initialBalance))

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = \\$1").
			WithArgs(initialBalance.Sub(txn.Amount), txn.WalletId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := db.WithdrawUpdate(sqlDB, txn)
		assert.Nil(t, err)
	})
}

func TestWithdrawUpdate_FailToUpdateBalance(t *testing.T) {
	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {
		counterpartyWalletId := sql.NullInt64{Valid: false}
		txn := &models.Transaction{
			WalletId:             1,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: counterpartyWalletId,
		}

		initialBalance := decimal.NewFromFloat(212.00)

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs(txn.WalletId).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(initialBalance))

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = \\$1").
			WithArgs(initialBalance.Sub(txn.Amount), txn.WalletId).
			WillReturnError(errors.New("db failed"))

		mock.ExpectRollback()

		err := db.WithdrawUpdate(sqlDB, txn)

		assert.NotNil(t, err)
		assert.Equal(t, "failed to update outgoing-balance: db failed", err.Error())

	})
}

func TestTransferUpdate_Success(t *testing.T) {
	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {

		initialBalance := decimal.NewFromFloat(212.00)

		txnOut := &models.Transaction{
			WalletId:             101,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: 102},
		}
		txnIn := &models.Transaction{
			WalletId:             102,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: 101},
		}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs(txnOut.WalletId).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(initialBalance))

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txnOut.WalletId, txnOut.Type, txnOut.Amount, txnOut.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = \\$1").
			WithArgs(initialBalance.Sub(txnOut.Amount), txnOut.WalletId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txnIn.WalletId, txnIn.Type, txnIn.Amount, txnIn.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+").
			WithArgs(txnIn.Amount, txnIn.WalletId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := db.TransferUpdate(sqlDB, txnOut, txnIn)
		assert.Nil(t, err)
	})
}

func TestTransferUpdate_FailToUpdateWithdrawBalance(t *testing.T) {
	testutils.WithDBMock(t, func(sqlDB *sql.DB, mock sqlmock.Sqlmock) {
		initialBalance := decimal.NewFromFloat(212.00)

		txnOut := &models.Transaction{
			WalletId:             101,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: 102},
		}
		txnIn := &models.Transaction{
			WalletId:             102,
			Type:                 models.TxnTypeWithdraw,
			Amount:               decimal.NewFromFloat(100.0),
			CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: 101},
		}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs(txnOut.WalletId).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(initialBalance))

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(txnOut.WalletId, txnOut.Type, txnOut.Amount, txnOut.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(123, time.Now()))

		mock.ExpectExec("UPDATE wallets SET balance = \\$1").
			WithArgs(initialBalance.Sub(txnOut.Amount), txnOut.WalletId).WillReturnError(errors.New("db failed"))

		mock.ExpectRollback()

		err := db.TransferUpdate(sqlDB, txnOut, txnIn)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to update outgoing-balance: db failed", err.Error())
	})
}
