package handler_test

import (
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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBalance_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	handler := &handler.HandlerDB{DB: db}

	// Mock input
	userID := int64(101)
	walletID := int64(1)
	currency := "USD"
	balance := decimal.NewFromFloat(100.00)
	createdAt := time.Now()

	// Step 1: Expect GetUserById
	mock.ExpectQuery("SELECT id, name, created_at FROM users where id=\\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(userID, "Alice", time.Now().AddDate(0, -2, 0)))

	// Step 2: Expect GetWalletByUserIDs
	mock.ExpectQuery("SELECT id, user_id, balance, currency, type, is_default, created_at FROM wallets WHERE user_id IN \\(\\$1\\) ORDER BY created_at DESC").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "type", "is_default", "created_at"}).
			AddRow(walletID, userID, balance, currency, "primary", true, createdAt))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/balance", userID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(userID, 10)})

	rec := httptest.NewRecorder()
	handler.HandleBalance(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Optional: Check if JSON body contains expected values
	var response models.WalletBalanceResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Alice", response.UserInfo.Name)
}
