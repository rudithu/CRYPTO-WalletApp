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

func getSourceWallet() models.Wallet {
	return models.Wallet{
		ID:        int64(201),
		UserId:    int64(101),
		Balance:   decimal.NewFromFloat(1500),
		IsDefault: true,
		Type:      "saving",
		Currency:  "SGD",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	}
}

func getTargetWallet() models.Wallet {
	return models.Wallet{
		ID:        int64(210),
		UserId:    int64(109),
		Balance:   decimal.NewFromFloat(1000),
		IsDefault: false,
		Type:      "saving",
		Currency:  "USD",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	}
}

func TestHandleTransferMoney_WalletToUser_DefaultWallet_Success(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		sourceTxnAmount := decimal.NewFromFloat(50)
		sourceWallet := getSourceWallet()

		targetUserId := int64(1)
		targetWalletId := int64(210)

		//GetWalletById
		testutils.MockGetWalletById(mock, sourceWallet)

		//GetDefaultWalletOrCurrencyByUserID
		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id = \\$1").
			WithArgs(targetUserId, sourceWallet.Currency).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
				AddRow(targetWalletId, targetUserId, decimal.NewFromFloat(100), "USD", "saving", true, time.Now()))

		// GetCcyRate
		mock.ExpectQuery("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = \\$1 AND to_ccy in \\(\\$2, \\$3\\)").
			WithArgs(models.BaseCcy, sourceWallet.Currency, "USD").
			WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}).AddRow("SGD", decimal.NewFromFloat(1.35)))

		mock.ExpectBegin()

		testutils.MockGetBalance(mock, sourceWallet.Balance, sourceWallet.ID)

		//createTransaction source
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(401), WalletId: sourceWallet.ID, Type: models.TxnTypeTransferOut,
			Amount: sourceTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: targetWalletId}, CreatedAt: time.Now()})

		//updateBalanceByWalletID source
		testutils.MockUpdateBalanceByWalletID(mock, sourceWallet.Balance.Sub(sourceTxnAmount), sourceWallet.ID)

		rate := decimal.NewFromFloat(1).Div(decimal.NewFromFloat(1.35))
		targetTxnAmount := sourceTxnAmount.Mul(rate)
		//createTransaction target
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(402), WalletId: targetWalletId, Type: models.TxnTypeTransferIn,
			Amount: targetTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: sourceWallet.ID}, CreatedAt: time.Now()})

		//updateBalanceByWalletID target
		testutils.MockIncrementBalanceByWalletID(mock, targetTxnAmount, targetWalletId)

		mock.ExpectCommit()

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_user_id": %d}`, sourceTxnAmount.String(), targetUserId)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", sourceWallet.ID), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(sourceWallet.ID, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestHandleTransferMoney_WalletToUser_SameCcy_Success(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		sourceTxnAmount := decimal.NewFromFloat(50)
		sourceWallet := getSourceWallet()

		targetUserId := int64(1)
		targetWalletId := int64(211)

		//GetWalletById
		testutils.MockGetWalletById(mock, sourceWallet)

		//GetDefaultWalletOrCurrencyByUserID
		mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id = \\$1").
			WithArgs(targetUserId, sourceWallet.Currency).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
				AddRow(int64(210), targetUserId, decimal.NewFromFloat(100), "USD", "saving", true, time.Now()).
				AddRow(targetWalletId, targetUserId, decimal.NewFromFloat(50), "SGD", "saving", false, time.Now()))

		mock.ExpectBegin()

		testutils.MockGetBalance(mock, sourceWallet.Balance, sourceWallet.ID)

		//createTransaction source
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(401), WalletId: sourceWallet.ID, Type: models.TxnTypeTransferOut,
			Amount: sourceTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: targetWalletId}, CreatedAt: time.Now()})

		//updateBalanceByWalletID source
		testutils.MockUpdateBalanceByWalletID(mock, sourceWallet.Balance.Sub(sourceTxnAmount), sourceWallet.ID)

		targetTxnAmount := sourceTxnAmount
		//createTransaction target
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(402), WalletId: targetWalletId, Type: models.TxnTypeTransferIn,
			Amount: targetTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: sourceWallet.ID}, CreatedAt: time.Now()})

		//updateBalanceByWalletID target
		testutils.MockIncrementBalanceByWalletID(mock, targetTxnAmount, targetWalletId)

		mock.ExpectCommit()

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_user_id": %d}`, sourceTxnAmount.String(), targetUserId)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", sourceWallet.ID), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(sourceWallet.ID, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestHandleTransferMoney_WalletToSameUser_NotSuccess(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		sourceTxnAmount := decimal.NewFromFloat(50)
		sourceWallet := getSourceWallet()

		//GetWalletById
		testutils.MockGetWalletById(mock, sourceWallet)

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_user_id": %d}`, sourceTxnAmount.String(), sourceWallet.UserId)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", sourceWallet.ID), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(sourceWallet.ID, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "please use destination_wallet_id to transfer for the same user")
	})
}

func TestHandleTransferMoney_WalletToWallet_Success(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		sourceTxnAmount := decimal.NewFromFloat(50)
		sourceWallet := getSourceWallet()
		targetWallet := getTargetWallet()

		//GetWalletById
		testutils.MockGetWalletById(mock, sourceWallet)

		//GetWalletById
		testutils.MockGetWalletById(mock, targetWallet)

		// GetCcyRate
		mock.ExpectQuery("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = \\$1 AND to_ccy in \\(\\$2, \\$3\\)").
			WithArgs(models.BaseCcy, sourceWallet.Currency, targetWallet.Currency).
			WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}).AddRow("SGD", decimal.NewFromFloat(1.35)))

		mock.ExpectBegin()

		testutils.MockGetBalance(mock, sourceWallet.Balance, sourceWallet.ID)

		//createTransaction source
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(401), WalletId: sourceWallet.ID, Type: models.TxnTypeTransferOut,
			Amount: sourceTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: targetWallet.ID}, CreatedAt: time.Now()})

		//updateBalanceByWalletID source
		testutils.MockUpdateBalanceByWalletID(mock, sourceWallet.Balance.Sub(sourceTxnAmount), sourceWallet.ID)

		rate := decimal.NewFromFloat(1).Div(decimal.NewFromFloat(1.35))
		targetTxnAmount := sourceTxnAmount.Mul(rate)
		//createTransaction target
		testutils.MockCreateTransaction(mock, models.Transaction{
			ID: int64(402), WalletId: targetWallet.ID, Type: models.TxnTypeTransferIn,
			Amount: targetTxnAmount, CounterpartyWalletId: sql.NullInt64{Valid: true, Int64: sourceWallet.ID}, CreatedAt: time.Now()})

		//updateBalanceByWalletID target
		testutils.MockIncrementBalanceByWalletID(mock, targetTxnAmount, targetWallet.ID)

		mock.ExpectCommit()

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_wallet_id": %d}`, sourceTxnAmount.String(), targetWallet.ID)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", sourceWallet.ID), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(sourceWallet.ID, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestHandleTransferMoney_BalanceNotEnough(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		sourceTxnAmount := decimal.NewFromFloat(2025)
		sourceWallet := getSourceWallet()
		targetWallet := getTargetWallet()

		//GetWalletById source
		testutils.MockGetWalletById(mock, sourceWallet)

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_wallet_id": %d}`, sourceTxnAmount.String(), targetWallet.ID)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", sourceWallet.ID), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(sourceWallet.ID, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "transferred is not allowed")
	})
}

func TestHandleTransferMoney_SourceWalletNotFound(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		sourceTxnAmount := decimal.NewFromFloat(10)
		walletId := int64(123)

		//GetWalletById source
		testutils.MockGetWalletByIdNoRecord(mock, int64(walletId))

		requestBody := fmt.Sprintf(`{"amount": %s, "destination_wallet_id": %d}`, sourceTxnAmount.String(), int64(64))
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", walletId), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(walletId, 10)})

		rec := httptest.NewRecorder()
		handler := handler.HandlerDB{DB: db}
		handler.HandleTransferMoney(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "source wallet not found")
	})

}

func TestHandleTransferMoney_InvalidRequestNoTarget(t *testing.T) {
	sourceTxnAmount := decimal.NewFromFloat(10)
	walletId := int64(123)

	requestBody := fmt.Sprintf(`{"amount": %s}`, sourceTxnAmount.String())
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", walletId), strings.NewReader(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(walletId, 10)})

	rec := httptest.NewRecorder()
	handler := handler.HandlerDB{DB: nil}
	handler.HandleTransferMoney(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "please specify either destination_user_id or destination_wallet_id")

}

func TestHandleTransferMoney_InvalidRequestBothTargetWalletAndUser(t *testing.T) {
	sourceTxnAmount := decimal.NewFromFloat(10)
	walletId := int64(123)

	requestBody := fmt.Sprintf(`{"amount": %s, "destination_wallet_id": 7, "destination_user_id": 8}`, sourceTxnAmount.String())
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/transfer", walletId), strings.NewReader(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(walletId, 10)})

	rec := httptest.NewRecorder()
	handler := handler.HandlerDB{DB: nil}
	handler.HandleTransferMoney(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "please specify only one of destination_user_id or destination_wallet_id, not both")
}

func TestHandleTransferMoney_InvalidSourceWalletId(t *testing.T) {
	sourceTxnAmount := decimal.NewFromFloat(10)
	walletId := "invalidId"

	requestBody := fmt.Sprintf(`{"amount": %s, "destination_wallet_id": 7}`, sourceTxnAmount.String())
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%s/transfer", walletId), strings.NewReader(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": walletId})

	rec := httptest.NewRecorder()
	handler := handler.HandlerDB{DB: nil}
	handler.HandleTransferMoney(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid wallet id")

}
