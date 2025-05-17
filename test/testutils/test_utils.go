package testutils

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func MockUser() *models.User {
	return &models.User{
		ID:   1,
		Name: "Alice",
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
