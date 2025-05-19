package handler_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestHandleWithdrawMoney_Success(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		handler := &handler.HandlerDB{DB: db}

		walletId := int64(1)
		amount := decimal.NewFromInt(50)

		// Prepare mock wallet
		testutils.MockGetWalletById(mock, models.Wallet{
			ID:       walletId,
			UserId:   int64(123),
			Balance:  decimal.NewFromFloat(100.00),
			Currency: "USD", Type: "saving",
			IsDefault: false,
			CreatedAt: time.Now(),
		})

		// Expect transaction (begin)
		mock.ExpectBegin()

		testutils.MockGetBalance(mock, decimal.NewFromFloat(100.00), walletId)

		// Insert transaction
		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(walletId, models.TxnTypeWithdraw, amount, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, time.Now()))

		// Update wallet balance
		mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
			WithArgs(decimal.NewFromInt(50), walletId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect commit
		mock.ExpectCommit()

		// Request body
		requestBody := `{"amount": 50}`
		req := httptest.NewRequest(http.MethodPost, "/wallets/1/withdraw", strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rr := httptest.NewRecorder()

		handler.HandleWithdrawMoney(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}

func TestHandleWithdrawMoney_WalletNotFound(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		handler := &handler.HandlerDB{DB: db}

		walletId := int64(2)

		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE id = \\$1").
			WithArgs(walletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}))

		// Request body
		requestBody := `{"amount": 50}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/withdraw", walletId), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(walletId, 10)})

		rec := httptest.NewRecorder()

		handler.HandleWithdrawMoney(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "wallet not found")
	})

}

func TestHandleWithdrawMoney_AmountIsGreaterThanBalance(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		walletId := int64(1)
		handler := &handler.HandlerDB{DB: db}

		// Prepare mock wallet
		testutils.MockGetWalletById(mock, models.Wallet{
			ID:        walletId,
			UserId:    int64(123),
			Balance:   decimal.NewFromFloat((100)),
			Currency:  "USE",
			Type:      "saving",
			IsDefault: false,
			CreatedAt: time.Now()})

		// Request body
		requestBody := `{"amount": 500}`
		req := httptest.NewRequest(http.MethodPost, "/wallets/1/withdraw", strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		rec := httptest.NewRecorder()

		handler.HandleWithdrawMoney(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "withdrawal is not allowed")
	})
}

func TestHandleWithdrawMoney_InvalidWalletId(t *testing.T) {

	handler := &handler.HandlerDB{DB: nil}

	requestBody := `{"amount": 500}`
	walletId := "invalidID"
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%s/withdraw", walletId), strings.NewReader(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": walletId})

	rr := httptest.NewRecorder()
	handler.HandleWithdrawMoney(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

}
