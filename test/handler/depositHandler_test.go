package handler_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHandleDepositMoney_Succes(t *testing.T) {

	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		handler := handler.HandlerDB{DB: db}

		txn := testutils.MockTxns()[0]
		txn.Type = models.TxnTypeDeposit

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO transactions \\(wallet_id, type, amount, counterparty_wallet_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, created_at").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(txn.ID, txn.CreatedAt))

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1 WHERE id = \\$2").
			WithArgs(txn.Amount, txn.WalletId).WillReturnResult(sqlmock.NewResult(10, 1))

		mock.ExpectCommit()

		requestBody := fmt.Sprintf(`{"amount": %s}`, txn.Amount.String())

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/deposit", txn.WalletId), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(txn.WalletId, 10)})

		rr := httptest.NewRecorder()
		handler.HandleDepositMoney(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

	})

}

func TestHandleDepositMoney_WalletNotExist(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		handler := handler.HandlerDB{DB: db}

		txn := testutils.MockTxns()[0]
		txn.Type = models.TxnTypeDeposit

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO transactions \\(wallet_id, type, amount, counterparty_wallet_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, created_at").
			WithArgs(txn.WalletId, txn.Type, txn.Amount, txn.CounterpartyWalletId).WillReturnError(errors.New("wallet id not exist"))
		mock.ExpectRollback()

		requestBody := fmt.Sprintf(`{"amount": %s}`, txn.Amount.String())

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%d/deposit", txn.WalletId), strings.NewReader(requestBody))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(txn.WalletId, 10)})

		rr := httptest.NewRecorder()
		handler.HandleDepositMoney(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

	})
}

func TestHandleDepositMoney_InvalidWalletId(t *testing.T) {

	handler := handler.HandlerDB{DB: nil}
	walledId := "invalidId"

	requestBody := `{"amount": 15000}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/wallets/%s/deposit", walledId), strings.NewReader(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": walledId})

	rr := httptest.NewRecorder()
	handler.HandleDepositMoney(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}
