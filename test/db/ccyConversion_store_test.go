package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetCcyRateToBaseCcy_Success(t *testing.T) {
	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {
		ccys := []string{"USD", "EUR"}
		expectedRates := []models.CcyRateToBaseCcy{
			{Ccy: "USD", Rate: decimal.NewFromFloat(1.2)},
			{Ccy: "EUR", Rate: decimal.NewFromFloat(1.5)},
		}

		rows := sqlmock.NewRows([]string{"to_ccy", "rate"}).
			AddRow("USD", expectedRates[0].Rate).
			AddRow("EUR", expectedRates[1].Rate)

		mock.ExpectQuery("SELECT to_ccy, rate FROM ccy_conversion").
			WithArgs("USD", "EUR").
			WillReturnRows(rows)

		rates, err := db.GetCcyRateToBaseCcy(dbTest, ccys)
		require.NoError(t, err)
		require.Len(t, rates, 2)
		require.Equal(t, expectedRates, rates)
	})
}

func TestGetCcyRate_Success(t *testing.T) {
	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {
		fromCcy := "JPY"
		toCcy := "EUR"
		baseCcy := models.BaseCcy

		fromRate := decimal.NewFromFloat(110.0) // Base -> JPY
		toRate := decimal.NewFromFloat(0.9)     // Base -> EUR

		rows := sqlmock.NewRows([]string{"to_ccy", "rate"}).
			AddRow(fromCcy, fromRate).
			AddRow(toCcy, toRate)

		mock.ExpectQuery("SELECT to_ccy, rate").
			WithArgs(baseCcy, fromCcy, toCcy).
			WillReturnRows(rows)

		result, err := db.GetCcyRate(dbTest, fromCcy, toCcy)
		require.NoError(t, err)

		expected := toRate.Div(fromRate)
		require.True(t, result.Equal(expected))
	})
}

func TestGetCcyRate_MissingRate(t *testing.T) {
	testutils.WithDBMock(t, func(dbTest *sql.DB, mock sqlmock.Sqlmock) {
		fromCcy := "JPY"
		toCcy := "EUR"

		rows := sqlmock.NewRows([]string{"to_ccy", "rate"})

		mock.ExpectQuery("SELECT to_ccy, rate").
			WithArgs(models.BaseCcy, fromCcy, toCcy).
			WillReturnRows(rows)

		result, err := db.GetCcyRate(dbTest, fromCcy, toCcy)
		require.Error(t, err)
		require.True(t, result.IsZero())

	})
}
