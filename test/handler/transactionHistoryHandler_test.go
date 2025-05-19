package handler_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleTxHistory_Success(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		user := testutils.MockUserModel()
		wallets := testutils.MockWallets()

		wallets[0].Currency = "SGD"
		for i := range wallets {
			wallets[i].UserId = user.ID
		}

		testutils.MockGetUserById(mock, user)
		testutils.MockGetWalletByUserIDs(mock, wallets)

		//GetTransactionsByWalletIDs
		mock.ExpectQuery("SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at FROM transactions WHERE wallet_id in \\(.*\\) ORDER BY created_at DESC").
			WithArgs(wallets[0].ID, wallets[1].ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "wallet_id", "type", "amount", "counterparty_wallet_id", "created_at"}).
				AddRow(int64(201), wallets[0].ID, models.TxnTypeWithdraw, decimal.NewFromFloat(100), sql.NullInt64{Valid: false}, time.Now().AddDate(0, -1, -3)).
				AddRow(int64(202), wallets[0].ID, models.TxnTypeTransferIn, decimal.NewFromFloat(100), int64(209), time.Now().AddDate(0, 0, -30)).
				AddRow(int64(203), wallets[1].ID, models.TxnTypeDeposit, decimal.NewFromFloat(100), sql.NullInt64{Valid: false}, time.Now().AddDate(0, -2, -10)).
				AddRow(int64(204), wallets[1].ID, models.TxnTypeDeposit, decimal.NewFromFloat(60), sql.NullInt64{Valid: false}, time.Now().AddDate(0, -1, -25)))

		mock.ExpectQuery(fmt.Sprintf("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = '%s' AND to_ccy IN \\(.+\\)", models.BaseCcy)).
			WithArgs(wallets[0].Currency, wallets[1].Currency).
			WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}).
				AddRow(wallets[0].Currency, decimal.NewFromFloat(1.35)).
				AddRow(wallets[1].Currency, decimal.NewFromFloat(0.92)))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/transactions", user.ID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(user.ID, 10)})

		rec := httptest.NewRecorder()

		handler := handler.HandlerDB{DB: db}
		handler.HandleTxHistory(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.WalletBalanceResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, user.ID, response.UserInfo.ID)
		assert.Equal(t, user.Name, response.UserInfo.Name)
		// assert.Equal(t, 2, len(response.Wallets))
		// assert.Equal(t, 2, len(response.Wallets[0].Transactions))
		// assert.Equal(t, 2, len(response.Wallets[1].Transactions))
		assert.NotNil(t, response.Balance)
	})

}

func TestHandleTxHistory_NoBalance_FilteredByWalletID(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		user := testutils.MockUserModel()

		wallets := testutils.MockWallets()
		wallets[0].Currency = "SGD"
		for i := range wallets {
			wallets[i].UserId = user.ID
		}

		testutils.MockGetUserById(mock, user)
		testutils.MockGetWalletById(mock, wallets[1])

		//GetTransactionsByWalletIDs
		mock.ExpectQuery("SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at FROM transactions WHERE wallet_id in \\(.*\\) ORDER BY created_at DESC").
			WithArgs(wallets[1].ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "wallet_id", "type", "amount", "counterparty_wallet_id", "created_at"}).
				AddRow(int64(203), wallets[1].ID, models.TxnTypeDeposit, decimal.NewFromFloat(100), sql.NullInt64{Valid: false}, time.Now().AddDate(0, -2, -10)).
				AddRow(int64(204), wallets[1].ID, models.TxnTypeDeposit, decimal.NewFromFloat(60), sql.NullInt64{Valid: false}, time.Now().AddDate(0, -1, -25)))

		mock.ExpectQuery(fmt.Sprintf("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = '%s' AND to_ccy IN \\(.+\\)", models.BaseCcy)).
			WithArgs(wallets[1].Currency).
			WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/transactions?wallet_id=%d", user.ID, wallets[1].ID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(user.ID, 10)})

		rec := httptest.NewRecorder()

		handler := handler.HandlerDB{DB: db}
		handler.HandleTxHistory(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.WalletBalanceResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, user.ID, response.UserInfo.ID)
		assert.Equal(t, user.Name, response.UserInfo.Name)
		assert.Equal(t, 1, len(response.Wallets))
		assert.Equal(t, 2, len(response.Wallets[0].Transactions))

		assert.Nil(t, response.Balance)
	})
}

func TestHandleTxHistory_InvalidUserId(t *testing.T) {

	userId := "invalidId"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/wallets/transactions", userId), nil)
	req = mux.SetURLVars(req, map[string]string{"id": userId})

	handler := handler.HandlerDB{DB: nil}
	rec := httptest.NewRecorder()

	handler.HandleTxHistory(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid user id")

}

func TestHandleTxHistory_InvalidWalletId(t *testing.T) {
	walletId := "invalidId"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/101/wallets/transactions?wallet_id=%s", walletId), nil)
	req = mux.SetURLVars(req, map[string]string{"id": "101"})

	handler := handler.HandlerDB{DB: nil}
	rec := httptest.NewRecorder()

	handler.HandleTxHistory(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid wallet id")

}
