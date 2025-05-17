package models_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultWalletOrCurrencyByUserID_Succes(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		expectedWallet := testutils.MockWallets()[0]
		ccy := "SGD"

		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
			AddRow(expectedWallet.ID, expectedWallet.UserId, expectedWallet.Balance, expectedWallet.Currency, expectedWallet.Type, expectedWallet.IsDefault, expectedWallet.CreatedAt)

		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id = \\$1 AND \\(is_default = TRUE .*\\) ORDER BY created_at DESC").
			WithArgs(expectedWallet.UserId, ccy).
			WillReturnRows(rows)

		wallets, err := models.GetDefaultWalletOrCurrencyByUserID(db, expectedWallet.UserId, ccy)

		assert.Nil(t, err)
		assert.NotNil(t, wallets)

	})

}

func TestGetDefaultWalletOrCurrencyByUserID_NotFound(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		userId := int64(201)
		ccy := "SGD"

		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id = \\$1 AND \\(is_default = TRUE .*\\) ORDER BY created_at DESC").
			WithArgs(userId, ccy).
			WillReturnError(sql.ErrNoRows)

		wallets, err := models.GetDefaultWalletOrCurrencyByUserID(db, userId, ccy)

		assert.Nil(t, err)
		assert.Nil(t, wallets)
	})
}

func TestGetDefaultWalletOrCurrencyByUserID_DBError(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		userId := int64(201)
		ccy := "SGD"

		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id = \\$1 AND \\(is_default = TRUE .*\\) ORDER BY created_at DESC").
			WithArgs(userId, ccy).
			WillReturnError(errors.New("db fail"))

		wallets, err := models.GetDefaultWalletOrCurrencyByUserID(db, userId, ccy)

		assert.NotNil(t, err)
		assert.Equal(t, "db fail", err.Error())
		assert.Nil(t, wallets)
	})
}

func TestGetWalletById_Success(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		expectedWallet := testutils.MockWallets()[0]

		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
			AddRow(expectedWallet.ID, expectedWallet.UserId, expectedWallet.Balance, expectedWallet.Currency, expectedWallet.Type, expectedWallet.IsDefault, expectedWallet.CreatedAt)

		mock.ExpectQuery(`SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \$1`).
			WithArgs(expectedWallet.ID).
			WillReturnRows(rows)

		wallet, err := models.GetWalletById(db, expectedWallet.ID)

		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, expectedWallet.ID, wallet.ID)
		assert.Equal(t, expectedWallet.UserId, wallet.UserId)
		assert.True(t, expectedWallet.Balance.Equal(wallet.Balance))
		assert.Equal(t, expectedWallet.Currency, wallet.Currency)
		assert.Equal(t, expectedWallet.Type, wallet.Type)
		assert.Equal(t, expectedWallet.IsDefault, wallet.IsDefault)

	})
}

func TestGetWalletById_NotFound(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"})

		mock.ExpectQuery(`SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		wallet, err := models.GetWalletById(db, int64(1))

		assert.Nil(t, err)
		assert.Nil(t, wallet)
	})

}

func TestGetWalletById_DBError(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db failed"))

		wallet, err := models.GetWalletById(db, int64(1))

		assert.Nil(t, wallet)
		assert.NotNil(t, err)
		assert.Equal(t, "db failed", err.Error())
	})
}

func TestGetWalletByUserIDs_Success(t *testing.T) {

	expectedWallet := testutils.MockWallets()[1]
	userIds := []int64{expectedWallet.UserId}

	rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
		AddRow(expectedWallet.ID, expectedWallet.UserId, expectedWallet.Balance, expectedWallet.Currency, expectedWallet.Type, expectedWallet.IsDefault, expectedWallet.CreatedAt)

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id IN \\(.+\\) ORDER BY created_at DESC").
			WithArgs(expectedWallet.UserId).
			WillReturnRows(rows)

		wallets, err := models.GetWalletByUserIDs(db, userIds)

		assert.Nil(t, err)
		assert.NotNil(t, wallets)

	})

}

func TestGetWalletByUserIDs_NotFound(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		userId := int64(201)
		userIds := []int64{userId}
		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id IN \\(.+\\) ORDER BY created_at DESC").
			WithArgs(201).
			WillReturnError(sql.ErrNoRows)

		wallets, err := models.GetWalletByUserIDs(db, userIds)
		assert.Nil(t, err)
		assert.Nil(t, wallets)

	})
}

func TestGetWalletByUserIDs_DBError(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		userId := int64(201)
		userIds := []int64{userId}
		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id IN \\(.+\\) ORDER BY created_at DESC").
			WithArgs(201).
			WillReturnError(errors.New("db failed"))

		wallets, err := models.GetWalletByUserIDs(db, userIds)
		assert.NotNil(t, err)
		assert.Equal(t, "db failed", err.Error())
		assert.Nil(t, wallets)
	})
}
