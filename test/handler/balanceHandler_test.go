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
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
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
	user := models.User{
		ID:        userID,
		Name:      "Alice",
		CreatedAt: time.Now().AddDate(0, -2, 0),
	}
	testutils.MockGetUserById(mock, user)

	// Step 2: Expect GetWalletByUserIDs
	testutils.MockGetWalletByUserIDs(mock, []models.Wallet{{
		ID:        walletID,
		UserId:    userID,
		Balance:   balance,
		Currency:  currency,
		IsDefault: true,
		CreatedAt: createdAt,
	}})

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
	assert.NotNil(t, response.Balance)
}

func TestHandleBalance_FilteredByWalletID(t *testing.T) {
	// Mocked data
	userID := int64(1)
	walletID := int64(1002)

	user := models.User{ID: userID, Name: "Alice", CreatedAt: time.Now().AddDate(-1, 0, 0)}

	wallet := models.Wallet{ID: walletID, UserId: userID, Balance: decimal.NewFromInt(50), Currency: "EUR", Type: "secondary", IsDefault: false}

	// Set up mocked DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Expect GetUserById
	testutils.MockGetUserById(mock, user)

	// Expect GetWalletById
	testutils.MockGetWalletById(mock, wallet)

	mock.ExpectQuery(fmt.Sprintf("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = '%s' AND to_ccy IN \\(.+\\)", models.BaseCcy)).
		WithArgs(wallet.Currency).
		WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}).AddRow(wallet.Currency, decimal.NewFromFloat(0.92)))

	// Prepare request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/balance?wallet_id=%d", userID, walletID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", userID)})

	rec := httptest.NewRecorder()
	handler := &handler.HandlerDB{DB: db}

	// Call handler
	handler.HandleBalance(rec, req)

	// Validate response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.WalletBalanceResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, userID, response.UserInfo.ID)
	assert.Len(t, response.Wallets, 1)
	assert.Equal(t, walletID, response.Wallets[0].ID)
	assert.NotNil(t, response.Balance)

	assert.NoError(t, mock.ExpectationsWereMet())
	// })

}

func TestHandleBalance_NoCurrencyRateNoBalance(t *testing.T) {
	// Mocked data
	userID := int64(1)
	walletID := int64(1002)

	user := models.User{ID: userID, Name: "Alice", CreatedAt: time.Now().AddDate(-1, 0, 0)}
	wallets := []models.Wallet{
		{ID: 1001, UserId: userID, Balance: decimal.NewFromInt(100), Currency: "USD", Type: "primary", IsDefault: true},
		{ID: walletID, UserId: userID, Balance: decimal.NewFromInt(50), Currency: "EUR", Type: "secondary", IsDefault: false},
	}

	// Set up mocked DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Expect GetUserById
	testutils.MockGetUserById(mock, user)

	// Expect GetWalletByUserIDs
	testutils.MockGetWalletByUserIDs(mock, wallets)

	mock.ExpectQuery(fmt.Sprintf("SELECT to_ccy, rate FROM ccy_conversion WHERE from_ccy = '%s' AND to_ccy IN \\(.+\\)", models.BaseCcy)).
		WithArgs(wallets[1].Currency).
		WillReturnRows(sqlmock.NewRows([]string{"to_ccy", "rate"}))

	// Prepare request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/balance", userID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", userID)})

	rec := httptest.NewRecorder()
	handler := &handler.HandlerDB{DB: db}

	// Call handler
	handler.HandleBalance(rec, req)

	// Validate response
	assert.Equal(t, http.StatusOK, rec.Code)
	var response models.WalletBalanceResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)

	require.NoError(t, err)

	assert.Equal(t, userID, response.UserInfo.ID)
	assert.Len(t, response.Wallets, 2)
	assert.Nil(t, response.Balance)

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestHandleBalance_InvalidUserID(t *testing.T) {
	// Setup a test HTTP request with an invalid user ID (e.g., "abc")
	userId := "invalidId"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/wallets/balance", userId), nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": userId,
	})

	// Create a response recorder to capture the response
	rec := httptest.NewRecorder()

	// Create a handler instance with a dummy DB (not used in this test)
	handler := &handler.HandlerDB{
		DB: nil, // DB won't be used due to early failure on user ID parse
	}

	// Call the handler
	handler.HandleBalance(rec, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid User ID")
}

func TestHandleBalance_InvalidWalletId(t *testing.T) {

	userID := int64(1)
	walletID := "invalidId"

	// Prepare request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/wallets/balance?wallet_id=%s", userID, walletID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", userID)})

	rec := httptest.NewRecorder()
	handler := &handler.HandlerDB{DB: nil}

	// Call handler
	handler.HandleBalance(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid Wallet ID")

}
