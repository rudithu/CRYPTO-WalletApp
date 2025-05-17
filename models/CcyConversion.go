package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type CcyConversion struct {
	FromCcy   string          `json:"from_ccy"`
	ToCcy     string          `json:"to_ccy"`
	Rate      decimal.Decimal `json:"rate"`
	CreatedAt time.Time       `json:"created_at"`
}

type CcyRateToBaseCcy struct {
	Ccy  string
	Rate decimal.Decimal
}

func GetCcyRateToBaseCcy(db *sql.DB, ccys []string) ([]CcyRateToBaseCcy, error) {
	placeholders := make([]string, len(ccys))
	args := make([]interface{}, len(ccys))

	for i, ccy := range ccys {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = ccy
	}

	query := fmt.Sprintf(`
		SELECT to_ccy, rate
		FROM ccy_conversion
		WHERE from_ccy = '%s' AND to_ccy IN (%s)
		`, BaseCcy, strings.Join(placeholders, ", "))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var ccyRates []CcyRateToBaseCcy
	for rows.Next() {
		var cr CcyRateToBaseCcy
		err = rows.Scan(&cr.Ccy, &cr.Rate)
		if err != nil {
			return nil, err
		}
		ccyRates = append(ccyRates, cr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ccyRates, nil
}

func GetCcyRate(db *sql.DB, fromCcy string, toCcy string) (decimal.Decimal, error) {

	query := `
		SELECT to_ccy, rate 
		FROM ccy_conversion
		WHERE from_ccy = $1 AND to_ccy in ($2, $3)
	`

	rows, err := db.Query(query, BaseCcy, fromCcy, toCcy)
	if err != nil {
		return decimal.Zero, err
	}
	defer rows.Close()

	var toRate decimal.Decimal
	var fromRate decimal.Decimal

	for rows.Next() {
		var rate decimal.Decimal
		var toCurrency string

		err := rows.Scan(&toCurrency, &rate)

		if err != nil {
			return decimal.Zero, err
		}

		if toCurrency == toCcy {
			toRate = rate
		} else if toCurrency == fromCcy {
			fromRate = rate
		}
	}

	if BaseCcy == toCcy {
		toRate = decimal.NewFromInt(1)
	} else if BaseCcy == fromCcy {
		fromRate = decimal.NewFromInt(1)
	}

	if toRate.IsZero() || fromRate.IsZero() {
		return decimal.Zero, fmt.Errorf("missing conversion rate for %s or %s", fromCcy, toCcy)
	}

	finalRate := toRate.Div(fromRate)

	return finalRate, nil

}
