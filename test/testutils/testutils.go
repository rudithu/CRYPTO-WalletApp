package testutils

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func MockUser() *models.User {
	return &models.User{
		ID:        1,
		Name:      "Alice",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	}
}

func MockUserModel() models.User {
	return models.User{
		ID:        101,
		Name:      "Bob",
		CreatedAt: time.Now().AddDate(-1, -1, 0),
	}
}

func MockWallets() []models.Wallet {
	return []models.Wallet{
		{
			ID:        101,
			UserId:    1,
			Currency:  "USD",
			Balance:   decimal.NewFromFloat(100.00),
			IsDefault: true,
			Type:      "primary",
		},
		{
			ID:        102,
			UserId:    1,
			Currency:  "EUR",
			Balance:   decimal.NewFromFloat(50.00),
			IsDefault: false,
			Type:      "secondary",
		},
	}
}

func MockTxns() []models.Transaction {
	return []models.Transaction{
		{
			ID:                   1,
			WalletId:             101,
			Type:                 "saving",
			Amount:               decimal.NewFromFloat(100.00),
			CreatedAt:            time.Now(),
			CounterpartyWalletId: NullInt64(0, false),
		},
	}
}

func MockCcyMapWithRate() map[string]models.CcyRateToBaseCcy {
	return map[string]models.CcyRateToBaseCcy{
		"EUR": {
			Ccy:  "EUR",
			Rate: decimal.NewFromFloat(0.5),
		},
	}
}

// Helper for sql.NullInt64 creation
func NullInt64(val int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{Int64: val, Valid: valid}
}

func WithDBMock(t *testing.T, testFunc func(db *sql.DB, mock sqlmock.Sqlmock)) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	testFunc(db, mock)
	defer db.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func MockGetUserById(mock sqlmock.Sqlmock, user models.User) {
	mock.ExpectQuery("SELECT id, name, created_at FROM users where id=\\$1").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(user.ID, user.Name, user.CreatedAt))
}

func MockGetWalletByUserIDs(mock sqlmock.Sqlmock, wallets []models.Wallet) {
	rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"})

	for _, w := range wallets {
		rows = rows.AddRow(w.ID, w.UserId, w.Balance, w.Currency, w.Type, w.IsDefault, w.CreatedAt)
	}

	mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id IN \\(\\$1\\) ORDER BY created_at DESC").
		WithArgs(wallets[0].UserId).
		WillReturnRows(rows)
}

func MockGetWalletById(mock sqlmock.Sqlmock, wallet models.Wallet) {
	mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \\$1").
		WithArgs(wallet.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
			AddRow(wallet.ID, wallet.UserId, wallet.Balance, wallet.Currency, wallet.Type, wallet.IsDefault, wallet.CreatedAt))
}

func MockGetWalletByIdNoRecord(mock sqlmock.Sqlmock, walletId int64) {
	mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \\$1").
		WithArgs(walletId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}))
}

func MockCreateTransaction(mock sqlmock.Sqlmock, txn models.Transaction) {
	mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(txn.WalletId, txn.Type, txn.Amount, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(txn.ID, txn.CreatedAt))
}

func MockCreateTransactionDBFailed(mock sqlmock.Sqlmock, txn models.Transaction) {
	mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(txn.WalletId, txn.Type, txn.Amount, sqlmock.AnyArg()).
		WillReturnError(errors.New("db failed"))
}

func MockUpdateBalanceByWalletID(mock sqlmock.Sqlmock, amount decimal.Decimal, walletId int64) {
	mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
		WithArgs(amount, walletId).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func MockUpdateBalanceByWalletIDDBFailed(mock sqlmock.Sqlmock, amount decimal.Decimal, walletId int64) {
	mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
		WithArgs(amount, walletId).
		WillReturnError(errors.New("db failed"))
}

func MockIncrementBalanceByWalletID(mock sqlmock.Sqlmock, amount decimal.Decimal, walletId int64) {
	mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1 WHERE id = \\$2").
		WithArgs(amount, walletId).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func MockGetBalance(mock sqlmock.Sqlmock, amount decimal.Decimal, walletId int64) {
	mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
		WithArgs(walletId).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(amount))
}
